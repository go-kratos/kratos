package service

import (
	"context"
	"sync"

	"go-common/app/job/main/workflow/conf"
	"go-common/app/job/main/workflow/dao"
	"go-common/app/job/main/workflow/model"
	"go-common/library/sync/pipeline/fanout"
)

// Service struct of service.
type Service struct {
	c            *conf.Config
	dao          *dao.Dao
	wg           *sync.WaitGroup
	closeCh      chan struct{}
	businessAttr []*model.BusinessAttr
	// cache
	cache *fanout.Fanout
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		dao:     dao.New(c),
		wg:      &sync.WaitGroup{},
		closeCh: make(chan struct{}),
		cache:   fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
	}
	var err error
	if s.businessAttr, err = s.dao.BusinessAttr(context.Background()); err != nil {
		panic(err)
	}
	//s.wg.Add(1)
	//go s.expireproc(context.Background())
	go s.queueproc(context.Background(), _feedbackDealType)
	go s.taskExpireproc(context.Background(), _feedbackDealType)
	go s.repairQueueproc(context.Background(), _feedbackDealType)
	// push
	go s.notifyproc(context.Background())

	// 单条申诉过期
	go s.singleExpireproc()
	// 整体申诉过期
	go s.overallExpireproc()
	// 释放用户未评价反馈
	go s.releaseExpireproc()
	// 刷新权重值
	go s.refreshWeightproc()
	// 进任务池
	go s.enterPoolproc()
	return
}

// Ping check service health.
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}

// Close related backend.
func (s *Service) Close() (err error) {
	err = s.dao.Close()
	close(s.closeCh)
	s.wg.Wait()

	return
}
