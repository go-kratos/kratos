package service

import (
	"go-common/app/interface/live/lottery-interface/internal/conf"
	risk "go-common/app/service/live/live_riskcontrol/api/grpc/v1"
	storm "go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/library/log/infoc"
)

// Service struct
type Service struct {
	c                 *conf.Config
	Infoc             *infoc.Infoc
	StormClient       storm.StormClient
	IsForbiddenClient risk.IsForbiddenClient
}

// New init
func New(c *conf.Config) (s *Service) {
	sc, err := storm.NewClient(c.LongClient)
	if err != nil {
		panic(err)
	}
	isForbiddenClient, err := risk.NewClient(c.ShortClient)
	if err != nil {
		panic(err)
	}
	s = &Service{
		c:                 c,
		Infoc:             infoc.New(c.Infoc),
		StormClient:       sc.StormClient,
		IsForbiddenClient: isForbiddenClient,
	}
	return s
}

// Close Service
func (s *Service) Close() {
	s.Infoc.Close()
}

// ServiceInstance instance
var ServiceInstance *Service

// Init init
func Init(c *conf.Config) {
	ServiceInstance = New(c)
}
