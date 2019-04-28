package service

import (
	"context"
	"runtime"

	"go-common/app/job/main/mcn/conf"
	"go-common/app/job/main/mcn/dao"
	accgrpc "go-common/app/service/main/account/api"
	"go-common/library/sync/pipeline/fanout"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// rpc
	accGRPC accgrpc.AccountClient
	worker  *fanout.Fanout
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		worker: fanout.New("cache", fanout.Worker(runtime.NumCPU()), fanout.Buffer(1024)),
	}
	var err error
	if s.accGRPC, err = accgrpc.NewClient(c.GRPCClient.Account); err != nil {
		panic(errors.WithMessage(err, "Failed to dial account service"))
	}
	if err := s.initEmailTemplate(); err != nil {
		panic(err)
	}
	t := cron.New()
	t.AddFunc(c.Property.UpMcnSignStateCron, s.UpMcnSignStateCron)
	t.AddFunc(c.Property.UpMcnUpStateCron, s.UpMcnUpStateCron)
	t.AddFunc(c.Property.UpExpirePayCron, s.UpExpirePayCron)
	//t.AddFunc(c.Property.UpMcnDataSummaryCron, s.UpMcnDataSummaryCron)
	t.AddFunc(c.Property.McnRecommendCron, s.McnRecommendCron)
	t.AddFunc(c.Property.DealFailRecommendCron, s.DealFailRecommendCron)
	t.AddFunc(c.Property.CheckMcnSignUpDueCron, s.CheckDateDueCron)
	t.Start()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
	s.worker.Close()
}
