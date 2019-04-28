package service

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"go-common/app/job/main/push/conf"
	"go-common/app/job/main/push/dao"
	pushrpc "go-common/app/service/main/push/api/grpc/v1"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/cache"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_max   = 1024
	_retry = 3
)

// Service .
type Service struct {
	c           *conf.Config
	dao         *dao.Dao
	waiter      sync.WaitGroup
	addTaskWg   sync.WaitGroup
	cache       *cache.Cache
	pushRPC     pushrpc.PushClient
	reportSub   *databus.Databus // consumer for new reports
	callbackSub *databus.Databus // consumer for callback
	reportCh    chan []*pushmdl.Report
	callbackCh  chan []*pushmdl.Callback
	addTaskCh   chan *pushmdl.Task
	reportCnt   int64
	callbackCnt int64
	closedCnt   int64
	closed      bool
}

// New creates a Service instance.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		dao:         dao.New(c),
		cache:       cache.New(1, 102400),
		reportSub:   databus.New(c.ReportSub),
		callbackSub: databus.New(c.CallbackSub),
		reportCh:    make(chan []*pushmdl.Report, 1024),
		callbackCh:  make(chan []*pushmdl.Callback, 1024),
		addTaskCh:   make(chan *pushmdl.Task, 10240),
	}
	var err error
	if s.pushRPC, err = pushrpc.NewClient(c.PushRPC); err != nil {
		panic(err)
	}
	if env.DeployEnv == env.DeployEnvProd {
		go s.delInvalidReportsproc() // 主动删除无效token
	}
	for i := 0; i < s.c.Job.ReportShard; i++ {
		s.waiter.Add(1)
		go s.reportproc()
	}
	for i := 0; i < s.c.Job.CallbackShard; i++ {
		s.waiter.Add(1)
		go s.callbackproc()
	}
	if s.c.Job.PretreatTask {
		for i := 0; i < s.c.Job.PretreatmentTaskShard; i++ {
			s.waiter.Add(1)
			go s.pretreatTaskproc() // 预处理任务，将任务转化成按平台分的token任务
		}
	}
	s.addTaskWg.Add(1)
	go s.addTaskproc()
	s.waiter.Add(1)
	go s.consumeReport()
	s.waiter.Add(1)
	go s.consumeCallback()
	go s.checkConsumer()
	// 删除过期的数据
	go s.delCallbacksproc()
	go s.delTasksproc()
	// 定期更新token缓存
	go s.refreshTokensproc()
	// data platform
	s.waiter.Add(1)
	go s.dpQueryproc()
	s.waiter.Add(1)
	go s.dpFileproc()
	return
}

// consumeReport consumes report.
func (s *Service) consumeReport() {
	defer s.waiter.Done()
	reports := make([]*pushmdl.Report, _max)
	ticker := time.NewTicker(time.Duration(s.c.Job.ReportTicker))
	for {
		select {
		case msg, ok := <-s.reportSub.Messages():
			if !ok {
				log.Info("databus: push-job report consumer exit!")
				if len(reports) > 0 {
					s.reportCh <- reports
				}
				if !atomic.CompareAndSwapInt64(&s.closedCnt, 0, 1) {
					close(s.reportCh)
				}
				return
			}
			s.reportCnt++
			msg.Commit()
			m := &pushmdl.Report{}
			if err := json.Unmarshal(msg.Value, m); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
				dao.PromError("service:解析计数databus消息")
				continue
			}
			log.Info("consumeReport key(%s) partition(%d) offset(%d) msg(%+v)", msg.Key, msg.Partition, msg.Offset, m)
			reports = append(reports, m)
			if len(reports) < _max {
				continue
			}
		case <-ticker.C:
		}
		if len(reports) > 0 {
			temp := make([]*pushmdl.Report, len(reports))
			copy(temp, reports)
			reports = []*pushmdl.Report{}
			s.reportCh <- temp
		}
	}
}

// checkConsumer checks consumer state.
func (s *Service) checkConsumer() {
	if env.DeployEnv != env.DeployEnvProd {
		return
	}
	var c1, c2 int64
	for {
		time.Sleep(5 * time.Minute)

		if s.reportCnt-c1 == 0 {
			msg := "push-job report did not consume within 5 minute"
			s.dao.SendWechat(msg)
			log.Warn(msg)
		}
		c1 = s.reportCnt

		if s.callbackCnt-c2 == 0 {
			msg := "push-job callback did not consume within 5 minute"
			s.dao.SendWechat(msg)
			log.Warn(msg)
		}
		c2 = s.callbackCnt
	}
}

// Ping reports the heath of services.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close releases resources which owned by the Service instance.
func (s *Service) Close() {
	s.closed = true
	s.reportSub.Close()
	s.callbackSub.Close()
	s.dao.Close()
	s.waiter.Wait()
	close(s.addTaskCh)
	s.addTaskWg.Wait()
}
