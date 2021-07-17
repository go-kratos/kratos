package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-kratos/kratos/v2/log"
)

type Message interface {
	Topic() string
	Value() []byte
}

type kafkaMessage struct {
	topic string
	value []byte
}

type Handler func(context.Context, Message)

func (m *kafkaMessage) Topic() string {
	return m.topic
}

func (m *kafkaMessage) Value() []byte {
	return m.value
}

func (m *kafkaMessage) Ack() error {
	return nil
}

func (m *kafkaMessage) NAck() error {
	return nil
}

type Event interface {
	Publish(topic string, value string) error
	Subscribe(topic string, handler Handler) error
	SubscribeChan(topic string) (<-chan Message, error)
	Ack(message Message) error
	NAck(message Message) error
	Close() error
}

type Kafka struct {
	consumer sarama.Consumer
	producer sarama.SyncProducer
	logger   log.Logger
}

type options struct {
	logger log.Logger
	addr   []string
}

type Option func(*options)

func WithKafkaAddr(addr []string) Option {
	return func(opts *options) {
		opts.addr = addr
	}
}

func WithKafkaLogger(logger log.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewKafkaEvent(opts ...Option) (Event, error) {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	if len(options.addr) == 0 {
		return nil, errors.New("addr is empty")
	}
	if options.logger == nil {
		options.logger = log.DefaultLogger
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// producer
	producer, err := sarama.NewSyncProducer(options.addr, config)
	if err != nil {
		_ = options.logger.Log(log.LevelFatal, "producer_test create producer error ", err.Error())
		return nil, err
	}

	// consumer
	consumer, err := sarama.NewConsumer(options.addr, nil)
	if err != nil {
		_ = options.logger.Log(log.LevelFatal, "producer_test create producer error ", err.Error())
		return nil, err
	}
	return &Kafka{consumer: consumer, producer: producer, logger: options.logger}, nil
}

func (k *Kafka) Publish(topic string, value string) error {
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

func (k *Kafka) Subscribe(topic string, handler Handler) error {
	partitionList, err := k.consumer.Partitions(topic)
	if err != nil {
		_ = k.logger.Log(log.LevelInfo, fmt.Sprintf("fail to get list of partition: %v", err))
		return err
	}
	for partition := range partitionList {
		pc, err := k.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			_ = k.logger.Log(log.LevelInfo, fmt.Sprintf("failed to start consumer for partition %d,err:%v\n", partition, err))
			return err
		}
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				handler(context.Background(), &kafkaMessage{
					topic: msg.Topic,
					value: msg.Value,
				})
			}
			defer pc.AsyncClose()
		}(pc)
	}
	return nil
}

func (k *Kafka) SubscribeChan(topic string) (<-chan Message, error) {
	message := make(chan Message, 256)
	partitionList, err := k.consumer.Partitions(topic)
	if err != nil {
		_ = k.logger.Log(log.LevelInfo, fmt.Sprintf("fail to get list of partition: %v", err))
		return nil, err
	}
	for partition := range partitionList {
		pc, err := k.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			_ = k.logger.Log(log.LevelInfo, fmt.Sprintf("failed to start consumer for partition %d,err:%v\n", partition, err))
			return nil, err
		}
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				message <- &kafkaMessage{
					topic: msg.Topic,
					value: msg.Value,
				}
			}
			defer pc.AsyncClose()
		}(pc)
	}
	return message, nil
}

func (k *Kafka) Ack(message Message) error {
	k.logger.Log(log.LevelInfo, "Ack", message.Topic())
	msg := message.(*kafkaMessage)
	return msg.Ack()
}

func (k *Kafka) NAck(message Message) error {
	k.logger.Log(log.LevelInfo, "NAck", message.Topic())
	msg := message.(*kafkaMessage)
	return msg.NAck()
}

func (k *Kafka) Close() error {
	err := k.consumer.Close()
	if err != nil {
		return err
	}
	err = k.producer.Close()
	if err != nil {
		return err
	}
	return nil
}
