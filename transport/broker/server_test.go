package broker_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/go-kratos/kratos/contrib/broker/kafka/v2"
	"github.com/go-kratos/kratos/v2/broker"
	brktransport "github.com/go-kratos/kratos/v2/transport/broker"
	"github.com/segmentio/kafka-go/snappy"
	"github.com/stretchr/testify/assert"
)

func TestQueueServer(t *testing.T) {
	var (
		ctx         = context.Background()
		addr        = []string{"127.0.0.1:9092"}
		group       = "kratos_test"
		topic       = "kratos1"
		kafkaBroker = kafka.NewBroker(addr,
			group,
			topic,
			kafka.WithConsumes(3),
			kafka.WithProcessors(16),
			kafka.WithCompressionCodec(snappy.NewCompressionCodec()),
		)
		srv = brktransport.NewServer(kafkaBroker,
			brktransport.WithConsumeHandler(Consume),
		)
	)
	go func() {
		for i := 0; i < 10; i++ {
			err := kafkaBroker.Publish(ctx, &broker.Message{
				Body: []byte("publish-" + strconv.Itoa(int(time.Now().UnixNano()))),
			})
			assert.NoError(t, err)
		}
	}()
	err := srv.Start(ctx)
	assert.NoError(t, err)
}

func Consume(ctx context.Context, msg *broker.Message) error {
	fmt.Println("receive msg: ", string(msg.Body))
	// return errors.New("do not ack")
	return nil
}
