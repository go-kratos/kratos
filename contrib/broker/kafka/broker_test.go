package kafka_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/contrib/broker/kafka/v2"
	"github.com/go-kratos/kratos/v2/broker"
	"github.com/segmentio/kafka-go/snappy"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	var (
		ctx         = context.Background()
		addr        = []string{"127.0.0.1:9092"}
		group       = "kratos_test"
		topic       = "kratos1"
		kafkaBroker = kafka.NewBroker(addr,
			group,
			topic,
			kafka.WithCommitIgnoreError(true),
			kafka.WithCompressionCodec(snappy.NewCompressionCodec()),
		)
	)

	go func() {
		for i := 0; i < 10; i++ {
			err := kafkaBroker.Publish(ctx, &broker.Message{
				Body: []byte(fmt.Sprintf("publish-%d_%d", time.Now().UnixNano(), i)),
			})
			assert.NoError(t, err)
		}
	}()
	for {
		kafkaBroker.Consume(ctx, func(ctx context.Context, m *broker.Message) error {
			fmt.Println("receive msg: ", string(m.Body))
			// return errors.New("do not ack")
			return nil
		})
	}
}
