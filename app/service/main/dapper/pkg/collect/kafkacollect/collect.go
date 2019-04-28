package kafkacollect

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"

	"go-common/app/service/main/dapper/model"
	"go-common/app/service/main/dapper/pkg/collect"
	"go-common/app/service/main/dapper/pkg/process"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

var (
	collectCount    = prom.New().WithCounter("dapper_kafka_collect_count", []string{"name"})
	collectErrCount = prom.New().WithCounter("dapper_kafka_collect_err_count", []string{"name"})
)

// Option set option
type Option func(*option)

type option struct {
	group string
	topic string
	addrs []string
}

func (o option) saramaConfig() *sarama.Config {
	return nil
}

var defaultOption = option{
	group: "default",
}

//func NewConsumer(addrs []string, config *Config) (Consumer, error)

// New kafka collect
func New(topic string, addrs []string, options ...Option) (collect.Collecter, error) {
	log.V(10).Info("new kafkacollect topic %s addrs: %v", topic, addrs)
	if len(addrs) == 0 {
		return nil, fmt.Errorf("kafka addrs required")
	}
	opt := defaultOption
	for _, fn := range options {
		fn(&opt)
	}
	opt.addrs = addrs
	opt.topic = topic
	clt := &kafkaCollect{opt: opt}
	return clt, nil
}

type kafkaCollect struct {
	wg            sync.WaitGroup
	opt           option
	ps            []process.Processer
	consumers     []*consumer
	client        sarama.Client
	offsetManager sarama.OffsetManager
	baseConsumer  sarama.Consumer
}

func (k *kafkaCollect) RegisterProcess(p process.Processer) {
	k.ps = append(k.ps, p)
}

func (k *kafkaCollect) Start() error {
	var err error
	if k.client, err = sarama.NewClient(k.opt.addrs, k.opt.saramaConfig()); err != nil {
		return fmt.Errorf("new kafka client error: %s", err)
	}
	if k.offsetManager, err = sarama.NewOffsetManagerFromClient(k.opt.group, k.client); err != nil {
		return fmt.Errorf("new offset manager error: %s", err)
	}
	if k.baseConsumer, err = sarama.NewConsumerFromClient(k.client); err != nil {
		return fmt.Errorf("new kafka consumer error: %s", err)
	}
	log.Info("kafkacollect consumer from topic: %s addrs: %s", k.opt.topic, k.opt.topic)
	return k.start()
}

func (k *kafkaCollect) handler(protoSpan *model.ProtoSpan) {
	var err error
	for _, p := range k.ps {
		if err = p.Process(context.Background(), protoSpan); err != nil {
			log.Error("process span error: %s, discard", err)
		}
	}
}

func (k *kafkaCollect) start() error {
	ps, err := k.client.Partitions(k.opt.topic)
	if err != nil {
		return fmt.Errorf("get partitions error: %s", err)
	}
	for _, p := range ps {
		var pom sarama.PartitionOffsetManager
		if pom, err = k.offsetManager.ManagePartition(k.opt.topic, p); err != nil {
			return fmt.Errorf("new manage partition error: %s", err)
		}
		offset, _ := pom.NextOffset()
		if offset == -1 {
			offset = sarama.OffsetOldest
		}
		var c sarama.PartitionConsumer
		log.V(10).Info("partitions %d start offset %d", p, offset)
		if c, err = k.baseConsumer.ConsumePartition(k.opt.topic, p, offset); err != nil {
			return fmt.Errorf("new consume partition error: %s", err)
		}
		log.V(10).Info("start partition consumer partition: %d, offset: %d", p, offset)
		consumer := newConsumer(k, c, pom)
		k.consumers = append(k.consumers, consumer)
		k.wg.Add(1)
		go consumer.start()
	}
	return nil
}

func (k *kafkaCollect) Close() error {
	for _, c := range k.consumers {
		if err := c.close(); err != nil {
			log.Warn("close consumer error: %s", err)
		}
	}
	k.wg.Wait()
	return nil
}

func newConsumer(k *kafkaCollect, c sarama.PartitionConsumer, pom sarama.PartitionOffsetManager) *consumer {
	return &consumer{kafkaCollect: k, consumer: c, pom: pom, closeCh: make(chan struct{}, 1)}
}

type consumer struct {
	*kafkaCollect
	pom      sarama.PartitionOffsetManager
	consumer sarama.PartitionConsumer
	closeCh  chan struct{}
}

func (c *consumer) close() error {
	c.closeCh <- struct{}{}
	c.pom.Close()
	return c.consumer.Close()
}

func (c *consumer) start() {
	defer c.wg.Done()
	var err error
	var value []byte
	for {
		select {
		case msg := <-c.consumer.Messages():
			collectCount.Incr("count")
			c.pom.MarkOffset(msg.Offset+1, "")
			log.V(10).Info("receive message from kafka topic: %s key: %s content: %s", msg.Key, msg.Topic, msg.Value)
			protoSpan := new(model.ProtoSpan)
			if err = json.Unmarshal(msg.Value, protoSpan); err != nil {
				collectErrCount.Incr("count_error")
				log.Error("unmarshal span from kafka error: %s, value: %v", err, value)
				continue
			}
			c.handler(protoSpan)
		case <-c.closeCh:
			log.V(10).Info("receive closed return")
			return
		}
	}
}
