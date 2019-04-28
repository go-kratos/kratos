package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/databus"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/Shopify/sarama"
	scluster "github.com/bsm/sarama-cluster"
)

var (
	commitInterval = 20 * time.Millisecond
)

// Client kafka client
type Client struct {
	sarama.Client
	config *sarama.Config
	addrs  []string
	topic  string
	group  string
}

// NewClient returns a new kafka client instance
func NewClient(addrs []string, topic, group string) (c *Client, err error) {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Offsets.CommitInterval = commitInterval
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	c = &Client{
		config: config,
		addrs:  addrs,
		topic:  topic,
		group:  group,
	}
	c.Client, err = sarama.NewClient(addrs, config)
	return
}

// Close close kafka client
func (c *Client) Close() {
	c.Client.Close()
}

// Topics return the list of topics for the given kafka cluster
func (c *Client) Topics() []string {
	topics, err := c.Client.Topics()
	if err != nil {
		panic(err)
	}
	return topics
}

// Partitions returns the list of parititons for the given topic, as registered with SetMetadata
func (c *Client) Partitions() []int32 {
	ps, err := c.Client.Partitions(c.topic)
	if err != nil {
		panic(err)
	}
	return ps
}

// OffsetNew return newest offset for given topic
func (c *Client) OffsetNew() (info map[int32]int64, err error) {
	var (
		offset int64
	)
	ps, err := c.Client.Partitions(c.topic)
	if err != nil {
		return
	}
	info = make(map[int32]int64)
	for _, p := range ps {
		offset, err = c.Client.GetOffset(c.topic, p, sarama.OffsetNewest)
		if err != nil {
			return
		}
		info[p] = offset
	}
	return
}

// OffsetOld return newest oldset for given topic
func (c *Client) OffsetOld() (info map[int32]int64, err error) {
	var (
		offset int64
	)
	ps, err := c.Client.Partitions(c.topic)
	if err != nil {
		return
	}
	info = make(map[int32]int64)
	for _, p := range ps {
		offset, err = c.Client.GetOffset(c.topic, p, sarama.OffsetOldest)
		if err != nil {
			return
		}
		info[p] = offset
	}
	return
}

// SetOffset commit offset of partition to kafka
func (c *Client) SetOffset(partition int32, offset int64) (err error) {
	om, err := sarama.NewOffsetManagerFromClient(c.group, c.Client)
	if err != nil {
		return
	}
	defer om.Close()
	pm, err := om.ManagePartition(c.topic, partition)
	if err != nil {
		return
	}
	pm.MarkOffset(offset, "")
	if err = pm.Close(); err != nil {
		return
	}
	time.Sleep(10 * commitInterval)
	// verify
	marked, err := c.OffsetMarked()
	log.Info("partititon:%d, before:%d, after:%d\n", partition, offset, marked[partition])
	return
}

// ResetOffset commit offset of partition to kafka
func (c *Client) ResetOffset(partition int32, offset int64) (err error) {
	om, err := sarama.NewOffsetManagerFromClient(c.group, c.Client)
	if err != nil {
		return
	}
	defer om.Close()
	pm, err := om.ManagePartition(c.topic, partition)
	if err != nil {
		return
	}
	pm.ResetOffset(offset, "")
	if err = pm.Close(); err != nil {
		return
	}
	time.Sleep(10 * commitInterval)
	// verify
	marked, err := c.OffsetMarked()
	log.Info("partititon:%d, before:%d, after:%d\n", partition, offset, marked[partition])
	return
}

// SeekBegin commit newest offset to kafka
func (c *Client) SeekBegin() (err error) {
	offset, err := c.OffsetOld()
	if err != nil {
		return
	}
	for partition, offset := range offset {
		err = c.ResetOffset(partition, offset)
		if err != nil {
			log.Info("partititon:%d, offset:%d, topic:%s, group:%s\n", partition, offset, c.topic, c.group)
			return
		}
	}
	return
}

// SeekEnd commit newest offset to kafka
func (c *Client) SeekEnd() (err error) {
	offset, err := c.OffsetNew()
	if err != nil {
		return
	}
	for partition, offset := range offset {
		err = c.SetOffset(partition, offset)
		if err != nil {
			log.Info("partititon:%d, offset:%d, topic:%s, group:%s\n", partition, offset, c.topic, c.group)
			return
		}
	}
	return
}

// OffsetMarked fetch marked offset
func (c *Client) OffsetMarked() (marked map[int32]int64, err error) {
	req := &sarama.OffsetFetchRequest{
		Version:       1,
		ConsumerGroup: c.group,
	}
	partitions := c.Partitions()
	marked = make(map[int32]int64)
	broker, err := c.Client.Coordinator(c.group)
	if err != nil {
		return
	}
	defer broker.Close()
	for _, partition := range partitions {
		req.AddPartition(c.topic, partition)
	}
	resp, err := broker.FetchOffset(req)
	if err != nil {
		return
	}
	for _, p := range partitions {
		block := resp.GetBlock(c.topic, p)
		if block == nil {
			err = sarama.ErrIncompleteResponse
			return
		}
		if block.Err == sarama.ErrNoError {
			marked[p] = block.Offset
		} else {
			err = block.Err
			log.Info("block error(%v)", err)
			return
		}
	}
	return
}

// Diff offset
func Diff(cluster, topic, group string) (records []*databus.Record, err error) {
	client, err := NewClient(conf.Conf.Kafka[cluster].Brokers, topic, group)
	if err != nil {
		log.Error("service.NewClient() error(%v)\n", err)
		return
	}
	defer client.Close()
	marked, err := client.OffsetMarked()
	if err != nil {
		log.Error("client.OffsetMarked() error(%v)\n", err)
		return
	}
	new, err := client.OffsetNew()
	if err != nil {
		log.Error("client.OffsetNew() error(%v)\n", err)
		return
	}
	records = make([]*databus.Record, 0, len(new))
	for partition, offset := range new {
		r := &databus.Record{
			Partition: partition,
			New:       offset,
		}
		if tmp, ok := marked[partition]; ok {
			r.Diff = offset - tmp
		}
		records = append(records, r)
	}
	return
}

// CreateTopic to kafka
func CreateTopic(addrs []string, topic string, partitions int32, factor int16) (err error) {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true

	cli, err := sarama.NewClient(addrs, config)
	if err != nil {
		return
	}
	defer cli.Close()

	var topics []string
	topics, err = cli.Topics()
	if err != nil {
		log.Error("CreateTopic get all topics on cluster err(%v)", err)
		return
	}
	for _, t := range topics {
		if t == topic {
			return
		}
	}

	broker, err := cli.Controller()
	if err != nil {
		return
	}
	// retention := "-1"
	req := &sarama.CreateTopicsRequest{
		TopicDetails: map[string]*sarama.TopicDetail{
			topic: {
				NumPartitions:     partitions,
				ReplicationFactor: factor,
				// ConfigEntries: map[string]*string{
				// 	"retention.ms": &retention,
				// },
			},
		},
		Timeout: time.Second,
	}
	var res *sarama.CreateTopicsResponse
	if res, err = broker.CreateTopics(req); err != nil {
		log.Info("CreateTopic CreateTopics error(%v) res(%v)", err, res)
		return
	}
	if !(res.TopicErrors[topic].Err == sarama.ErrNoError || res.TopicErrors[topic].Err == sarama.ErrTopicAlreadyExists) {
		log.Error("CreateTopic CreateTopics kfafka topic create error(%v) topic(%v)", res.TopicErrors[topic].Err, topic)
		err = res.TopicErrors[topic].Err
	}
	log.Info("CreateTopic CreateTopics kfafka topic create info(%v) topic(%v)", res.TopicErrors[topic].Err, topic)
	return
}

// FetchMessage fetch topic message by timestamp or offset.
func FetchMessage(c context.Context, cluster, topic, group, key string, start, end int64, limit int) (res []*databus.Message, err error) {
	res = make([]*databus.Message, 0)
	kc, ok := conf.Conf.Kafka[cluster]
	if !ok {
		err = ecode.NothingFound
		return
	}
	cfg := scluster.NewConfig()
	cfg.Version = sarama.V1_0_0_0
	cfg.ClientID = fmt.Sprintf("%s-%s", group, topic)
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Return.Errors = true
	consumer, err := scluster.NewConsumer(kc.Brokers, group, []string{topic}, cfg)
	if err != nil {
		log.Error("fetchMsg NewConsumer(%v,%s,%s) error(%v)", kc.Brokers, group, topic, err)
		return
	}
	defer consumer.Close()
	bkey := []byte(key)
	for {
		select {
		case msg := <-consumer.Messages():
			consumer.MarkPartitionOffset(topic, int32(msg.Partition), msg.Offset, "")
			if key != "" && bytes.Equal(bkey, msg.Key) {
				continue
			}
			if start > 0 && msg.Timestamp.UnixNano()/int64(time.Millisecond) < start {
				continue
			}
			r := &databus.Message{}
			r.Key = string(msg.Key)
			r.Value = string(msg.Value)
			r.Topic = topic
			r.Partition = msg.Partition
			r.Offset = msg.Offset
			r.Timestamp = msg.Timestamp.Unix()
			res = append(res, r)
			if len(res) >= limit {
				return
			}
			if end > 0 && msg.Timestamp.UnixNano()/int64(time.Millisecond) > end {
				return
			}
		case err = <-consumer.Errors():
			log.Error("fetchMsg error(%v)", err)
			return
		case <-time.After(time.Second * 10):
			err = ecode.Deadline
			return
		case <-c.Done():
			return
		}
	}
}

// NewOffset commit new offset to kafka
func (c *Client) NewOffset(partition int32, offset int64) (err error) {
	err = c.ResetOffset(partition, offset)
	if err != nil {
		log.Info("partititon:%d, offset:%d, topic:%s, group:%s\n", partition, offset, c.topic, c.group)
		return
	}
	return
}

// OffsetNewTime return newest offset for given topic
func (c *Client) OffsetNewTime(time int64) (info map[int32]int64, err error) {
	var (
		offset  int64
		message string
	)
	ps, err := c.Client.Partitions(c.topic)
	if err != nil {
		return
	}
	info = make(map[int32]int64)
	for _, p := range ps {
		offset, err = c.Client.GetOffset(c.topic, p, time)
		message += fmt.Sprintf("partititon:%v, offset:%v, topic:%v, group:%v, time:%v, err:%v\n", p, offset, c.topic, c.group, time, err)
		if err != nil {
			log.Info(message)
			return
		}
		info[p] = offset
	}
	log.Info(message)
	return
}

// NewTime commit new time offset to kafka
func (c *Client) NewTime(time int64) (err error) {
	offset, err := c.OffsetNewTime(time)
	if err != nil {
		return
	}
	for partition, offset := range offset {
		err = c.ResetOffset(partition, offset)
		if err != nil {
			log.Info("partititon:%d, offset:%d, topic:%s, group:%s\n", partition, offset, c.topic, c.group)
			return
		}
	}
	return
}
