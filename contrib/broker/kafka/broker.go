package kafka

import (
	"context"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/broker"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"github.com/segmentio/kafka-go/protocol"
)

type kafkaBroker struct {
	brokers           []string
	group             string
	topic             string
	compressionCodec  compress.Codec
	offset            int64
	minBytes          int
	maxBytes          int
	maxWait           time.Duration
	commitInterval    time.Duration
	queueCapacity     int
	commitIgnoreError bool

	kafkaReader *kafka.Reader
	kafkaWriter *kafka.Writer

	channel            chan kafka.Message
	consumers          int
	consumerWaitGroup  sync.WaitGroup
	processors         int
	processorWaitGroup sync.WaitGroup
}

func NewBroker(brokers []string, group string, topic string, opts ...Option) broker.Broker {
	kb := &kafkaBroker{
		brokers:    brokers,
		group:      group,
		topic:      topic,
		consumers:  1,
		processors: 1,
		channel:    make(chan kafka.Message),
	}
	for _, o := range opts {
		o(kb)
	}
	kb.kafkaReader = kb.newKafkaReader()
	kb.kafkaWriter = kb.newKafkaWriter()
	return kb
}

// Close implements broker.Broker
func (b *kafkaBroker) Close() error {
	if b.kafkaWriter != nil {
		b.kafkaWriter.Close()
	}
	if b.kafkaReader != nil {
		b.kafkaReader.Close()
	}
	return nil
}

// Publish implements broker.Broker
func (b *kafkaBroker) Publish(ctx context.Context, msg *broker.Message) error {
	m := kafka.Message{
		Key:     []byte(strconv.FormatInt(time.Now().UnixNano(), 10)),
		Value:   msg.Body,
		Headers: parseKafkaHeader(msg.Header),
	}
	return b.kafkaWriter.WriteMessages(ctx, m)
}

// Consume implements broker.Broker
func (b *kafkaBroker) Consume(ctx context.Context, handler broker.ConsumeHandler) {
	for i := 0; i < b.processors; i++ {
		b.consumerWaitGroup.Add(1)

		go func() {
			b.consumerWaitGroup.Done()

			for {
				msg, err := b.kafkaReader.FetchMessage(ctx)
				if err == io.EOF || err == io.ErrClosedPipe {
					return
				}
				if err != nil {
					continue
				}
				b.channel <- msg
			}
		}()
	}

	b.consumeMessage(ctx, handler)
}

// consumeMessage implements broker.Broker
func (b *kafkaBroker) consumeMessage(ctx context.Context, handler broker.ConsumeHandler) {
	for i := 0; i < b.processors; i++ {
		b.processorWaitGroup.Add(1)

		go func() {
			b.processorWaitGroup.Done()

			for msg := range b.channel {
				err := handler(ctx, &broker.Message{
					Header: parseBrokerHeader(msg.Headers),
					Body:   msg.Value,
				})
				if err == nil || b.commitIgnoreError {
					b.kafkaReader.CommitMessages(ctx, msg)
					return
				}
				log.Errorf("consuming '%s' error: %v", string(msg.Value), err)
			}
		}()
	}
}

func (b *kafkaBroker) newKafkaReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        b.brokers,
		GroupID:        b.group,
		Topic:          b.topic,
		StartOffset:    b.offset,
		MinBytes:       b.minBytes,
		MaxBytes:       b.maxBytes,
		MaxWait:        b.maxWait,
		CommitInterval: b.commitInterval,
		QueueCapacity:  b.queueCapacity,
	})
}

func (b *kafkaBroker) newKafkaWriter() *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:          b.brokers,
		Topic:            b.topic,
		Balancer:         &kafka.LeastBytes{},
		CompressionCodec: b.compressionCodec,
	})
}

func parseKafkaHeader(in []broker.Header) []protocol.Header {
	var headers = make([]protocol.Header, len(in))
	for i, h := range in {
		headers[i] = protocol.Header{
			Key:   h.Key,
			Value: h.Value,
		}
	}
	return headers
}

func parseBrokerHeader(in []protocol.Header) []broker.Header {
	var headers = make([]broker.Header, len(in))
	for i, h := range in {
		headers[i] = broker.Header{
			Key:   h.Key,
			Value: h.Value,
		}
	}
	return headers
}
