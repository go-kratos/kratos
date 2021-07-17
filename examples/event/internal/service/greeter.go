package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/examples/event/internal/data"
	"github.com/go-kratos/kratos/v2/errors"

	v1 "github.com/go-kratos/kratos/examples/event/api/helloworld/v1"
	"github.com/go-kratos/kratos/examples/event/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc   *biz.GreeterUsecase
	log  *log.Helper
	data *data.Data
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger, d *data.Data) *GreeterService {
	service := &GreeterService{uc: uc, log: log.NewHelper(logger), data: d}
	err := d.Event.Subscribe("foo", func(ctx context.Context, message data.Message) {
		_ = logger.Log(log.LevelInfo, "topic", message.Topic(), "value", string(message.Value()))
		_ = d.Event.Ack(message)
	})
	if err != nil {
		return nil
	}
	go service.SubscribeChan()
	return service
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	s.log.WithContext(ctx).Infof("SayHello Received: %v", in.GetName())
	err := s.data.Event.Publish("foo", fmt.Sprintf("hello %s", in.GetName()))
	if err != nil {
		return nil, errors.New(500, "SEND_MESSAGE_ERROR", fmt.Sprintf("send message error:%s", err.Error()))
	}
	err = s.data.Event.Publish("bar", fmt.Sprintf("hello %s", in.GetName()))
	if err != nil {
		return nil, errors.New(500, "SEND_MESSAGE_ERROR", fmt.Sprintf("send message error:%s", err.Error()))
	}
	if in.GetName() == "error" {
		return nil, v1.ErrorUserNotFound("user not found: %s", in.GetName())
	}
	return &v1.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *GreeterService) SubscribeChan() {
	subscribeChan, err := s.data.Event.SubscribeChan("bar")
	if err != nil {
		return
	}
	for message := range subscribeChan {
		s.log.Infow("topic", message.Topic(), "value", string(message.Value()))
		_ = s.data.Event.NAck(message)
	}
}
