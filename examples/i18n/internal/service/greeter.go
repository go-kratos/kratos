package service

import (
	"context"
	"github.com/go-kratos/kratos/examples/i18n/internal/pkg/middleware/localize"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/go-kratos/kratos/examples/i18n/api/helloworld/v1"
	"github.com/go-kratos/kratos/examples/i18n/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc  *biz.GreeterUsecase
	log *log.Helper
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{uc: uc, log: log.NewHelper(logger)}
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	s.log.WithContext(ctx).Infof("SayHello Received: %v", in.GetName())

	if in.GetName() == "error" {
		return nil, v1.ErrorUserNotFound("user not found: %s", in.GetName())
	}
	localizer := localize.FromContext(ctx)
	helloMsg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			Description: "sayhello",
			ID:    "sayHello",
			One:   "Hello {{.Name}}",
			Other: "Hello {{.Name}}",
		},
		TemplateData: map[string]interface{}{
			"Name": in.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: helloMsg}, nil
}
