package service

import (
	"context"
	"fmt"

	v1 "github.com/go-kratos/kratos/examples/queue/api/helloworld/v1"
	"github.com/go-kratos/kratos/examples/queue/internal/biz"
	"github.com/go-kratos/kratos/examples/queue/internal/data"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc    *biz.GreeterUsecase
	log   *log.Helper
	kafka data.Queue
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger, d *data.Data) *GreeterService {
	service := &GreeterService{uc: uc, log: log.NewHelper(logger), kafka: d.Kafka}
	// kafka subscribe topic
	err := service.kafka.Subscribe("hello", func(ctx context.Context, msg *data.Message) {
		service.log.Infow("type", "subscribe", "topic", msg.Topic, "value", string(msg.Value))
	})
	if err != nil {
		return nil
	}
	// kafka Subscribe topic channel
	go service.kafkaQueue()
	return service
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	s.log.WithContext(ctx).Infof("SayHello Received: %v", in.GetName())
	err := s.kafka.Send("hello", fmt.Sprintf("hello %s", in.GetName()))
	if err != nil {
		return nil, errors.New(500, "MQ_ERROR", "send message error")
	}
	err = s.kafka.Send("hello_chan", fmt.Sprintf("hello %s", in.GetName()))
	if err != nil {
		return nil, errors.New(500, "MQ_ERROR", "send message error")
	}
	return &v1.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *GreeterService) kafkaQueue() {
	consumer, err := s.kafka.SubscribeChan("hello_chan", 256)
	if err != nil {
		return
	}
	for msg := range consumer.Receive() {
		s.log.Infow("type", "subscribe channel", "topic", msg.Topic, "value", string(msg.Value))
	}
}
