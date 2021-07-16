package data

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-kratos/kratos/v2/log"
)

var _ Queue = (*Kafka)(nil)

type Message struct {
	Value    []byte
	Topic    string
}

type Queue interface {
	Send(topic string, value string) error
	SubscribeChan(topic string, chanSize int) (Consumer, error)
	Subscribe(topic string, handler Handler) error
	close()
}
type Handler func(context.Context, *Message)

type Consumer interface {
	Receive() <-chan *Message
}

type consumer struct {
	message chan *Message
}

type Kafka struct {
	consumer sarama.Consumer
	producer sarama.SyncProducer
	log      log.Logger
}

func NewKafka(address []string, logger log.Logger) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// producer
	producer, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		_ = logger.Log(log.LevelFatal, "producer_test create producer error ", err.Error())
		return nil, err
	}

	// consumer
	consumer, err := sarama.NewConsumer(address, nil)
	if err != nil {
		_ = logger.Log(log.LevelFatal, "producer_test create producer error ", err.Error())
		return nil, err
	}
	return &Kafka{consumer: consumer, producer: producer, log: logger}, nil
}

func (k *Kafka) Send(topic string, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}
	_, _, err := k.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func (k *Kafka) Subscribe(topic string,handler Handler) error {
	partitionList, err := k.consumer.Partitions(topic)
	if err != nil {
		_ = k.log.Log(log.LevelInfo, fmt.Sprintf("fail to get list of partition: %v", err))
		return err
	}
	for partition := range partitionList {
		pc, err := k.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			_ = k.log.Log(log.LevelInfo, fmt.Sprintf("failed to start consumer for partition %d,err:%v\n", partition, err))
			return err
		}
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				handler(context.Background(),&Message{
					Value:    msg.Value,
					Topic:    msg.Topic,
				})
			}
			defer pc.AsyncClose()
		}(pc)
	}
	return nil
}

func (k *Kafka) SubscribeChan(topic string, chanSize int) (Consumer, error) {
	c := &consumer{message: make(chan *Message, chanSize)}
	partitionList, err := k.consumer.Partitions(topic)
	if err != nil {
		_ = k.log.Log(log.LevelInfo, fmt.Sprintf("fail to get list of partition: %v", err))
		return nil, err
	}
	for partition := range partitionList {
		pc, err := k.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			_ = k.log.Log(log.LevelInfo, fmt.Sprintf("failed to start consumer for partition %d,err:%v\n", partition, err))
			return nil, err
		}
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				c.message <- &Message{
					Value:    msg.Value,
					Topic:    msg.Topic,
				}
			}
			defer pc.AsyncClose()
		}(pc)
	}
	return c, nil
}

func (c *consumer) Receive() <-chan *Message {
	return c.message
}

func (k *Kafka) close() {
	_ = k.producer.Close()
	_ = k.consumer.Close()
}
