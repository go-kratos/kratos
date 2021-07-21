package kafka

import (
	"context"
	"github.com/go-kratos/kratos/examples/event/event"

	"github.com/Shopify/sarama"
	"github.com/go-kratos/kratos/v2/log"
)

var _ event.Sender = (*kafkaSender)(nil)
var _ event.Receiver = (*kafkaReceiver)(nil)
var _ event.Message = (*Message)(nil)

type Option func(*options)

type options struct {
	logger *log.Helper
}

// WithLogger with logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = log.NewHelper(logger)
	}
}

type Message struct {
	key    string
	value  []byte
	header map[string]string
}

func (m *Message) Key() string {
	return m.key
}
func (m *Message) Value() []byte {
	return m.value
}
func (m *Message) Header() map[string]string {
	return m.header
}
func NewKafkaClient(address []string) (sarama.Client, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.DefaultVersion
	client, err := sarama.NewClient(address, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMessage(key string, value []byte, header map[string]string) event.Message {
	return &Message{
		key:    key,
		value:  value,
		header: header,
	}
}

type kafkaSender struct {
	producer sarama.SyncProducer
	topic    string
	options  *options
}

func (s *kafkaSender) Send(ctx context.Context, message event.Message) error {
	msg := &sarama.ProducerMessage{
		Topic: s.topic,
		Key:   sarama.StringEncoder(message.Key()),
		Value: sarama.ByteEncoder(message.Value()),
	}
	if len(message.Header()) > 0 {
		msg.Headers = header2RecordHeader(message.Header())
	}
	_, _, err := s.producer.SendMessage(msg)
	if err != nil {
		s.options.logger.Errorw(err.Error())
		return err
	}
	return nil
}

func (s *kafkaSender) Close() error {
	err := s.producer.Close()
	if err != nil {
		return err
	}
	return err
}

func NewKafkaSender(client sarama.Client, topic string, opts ...Option) (event.Sender, error) {
	options := &options{
		logger: log.NewHelper(log.DefaultLogger),
	}
	for _, o := range opts {
		o(options)
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &kafkaSender{options: options, producer: producer, topic: topic}, nil
}

type kafkaReceiver struct {
	consumer sarama.Consumer
	topic    string
	options  *options
}

func (k *kafkaReceiver) Receive(ctx context.Context, handler event.Handler) error {
	partitionList, err := k.consumer.Partitions(k.topic)
	if err != nil {
		k.options.logger.Errorw("err", err.Error())
		return err
	}
	for partition := range partitionList {
		pc, err := k.consumer.ConsumePartition(k.topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			k.options.logger.Errorw("err", err.Error())
		}
		go func() {
			for msg := range pc.Messages() {
				err := handler(context.Background(), &Message{
					key:    string(msg.Key),
					value:  msg.Value,
					header: recordHeader2Header(msg.Headers),
				})
				if err != nil {
					k.options.logger.Errorw("err", err.Error())
					// do msg nack
				}
				// do msg ack
			}
		}()
	}
	return nil
}

func (k *kafkaReceiver) Close() error {
	err := k.consumer.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaReceiver(client sarama.Client, topic string, opts ...Option) (event.Receiver, error) {
	options := &options{
		logger: log.NewHelper(log.DefaultLogger),
	}
	for _, o := range opts {
		o(options)
	}
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &kafkaReceiver{options: options, consumer: consumer, topic: topic}, nil
}

func header2RecordHeader(kvs map[string]string) []sarama.RecordHeader {
	var h []sarama.RecordHeader
	for key, value := range kvs {
		h = append(h, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(value),
		})
	}
	return h
}

func recordHeader2Header(kvs []*sarama.RecordHeader) map[string]string {
	h := make(map[string]string)
	for _, value := range kvs {
		h[string(value.Key)] = string(value.Value)
	}
	return h
}
