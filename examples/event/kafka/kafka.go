package kafka

import (
	"context"
	"github.com/go-kratos/kratos/examples/event/event"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"log"
)

var _ event.Sender = (*kafkaSender)(nil)
var _ event.Receiver = (*kafkaReceiver)(nil)
var _ event.Message = (*Message)(nil)

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

func NewMessage(key string, value []byte, header map[string]string) event.Message {
	return &Message{
		key:    key,
		value:  value,
		header: header,
	}
}

type kafkaSender struct {
	writer *kafka.Writer
	topic  string
}

func (s *kafkaSender) Send(ctx context.Context, message event.Message) error {
	var h []kafka.Header
	if len(message.Header()) > 0 {
		for k, v := range message.Header() {
			h = append(h, protocol.Header{
				Key:   k,
				Value: []byte(v),
			})
		}
	}
	err := s.writer.WriteMessages(ctx, kafka.Message{
		Key:     []byte(message.Key()),
		Value:   message.Value(),
		Headers: h,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *kafkaSender) Close() error {
	err := s.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaSender(address []string, topic string) (event.Sender, error) {

	w := &kafka.Writer{
		Topic:    topic,
		Addr:     kafka.TCP(address...),
		Balancer: &kafka.LeastBytes{},
	}
	return &kafkaSender{writer: w, topic: topic}, nil
}

type kafkaReceiver struct {
	reader *kafka.Reader
	topic  string
}

func (k *kafkaReceiver) Receive(ctx context.Context, handler event.Handler) error {
	go func() {
		for {
			m, err := k.reader.FetchMessage(context.Background())
			if err != nil {
				break
			}
			h := make(map[string]string)
			if len(m.Headers) > 0 {
				for _, header := range m.Headers {
					h[header.Key] = string(header.Value)
				}
			}
			err = handler(context.Background(), &Message{
				key:    string(m.Key),
				value:  m.Value,
				header: h,
			})
			if err != nil {
				log.Fatal("message handling exception:", err)
			}
			if err := k.reader.CommitMessages(ctx, m); err != nil {
				log.Fatal("failed to commit messages:", err)
			}
		}
	}()
	return nil
}

func (k *kafkaReceiver) Close() error {
	err := k.reader.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaReceiver(address []string, topic string) (event.Receiver, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  address,
		GroupID:  "group-a",
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &kafkaReceiver{reader: r, topic: topic}, nil
}
