package broker_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/broker"
	brktransport "github.com/go-kratos/kratos/v2/transport/broker"
	"github.com/stretchr/testify/assert"
)

func TestQueueServer(t *testing.T) {
	var (
		wg    sync.WaitGroup
		ctx   = context.Background()
		topic = "kratos1"
		brk   = newMemoryBroker(topic)
		srv   = brktransport.NewServer(brk, brktransport.WithConsumeHandler(Consume))
	)
	wg.Add(1)
	go func() {
		wg.Done()
		for i := 0; i < 10; i++ {
			err := brk.Publish(ctx, &broker.Message{
				Body: []byte("publish-" + strconv.Itoa(int(time.Now().UnixNano()))),
			})
			assert.NoError(t, err)
		}
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		err := srv.Start(ctx)
		assert.NoError(t, err)
	}()
}

func Consume(ctx context.Context, msg *broker.Message) error {
	fmt.Println("receive msg: ", string(msg.Body))
	return nil
}

type memoryBroker struct{}

// Close implements broker.Broker
func (*memoryBroker) Close() error {
	return nil
}

// Consume implements broker.Broker
func (*memoryBroker) Consume(context.Context, broker.ConsumeHandler) {

}

// Publish implements broker.Broker
func (*memoryBroker) Publish(context.Context, *broker.Message) error {
	return nil
}

func newMemoryBroker(topic string) broker.Broker {
	return &memoryBroker{}
}
