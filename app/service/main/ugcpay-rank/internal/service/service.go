package service

import (
	"context"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"go-common/app/service/main/ugcpay-rank/internal/conf"
	"go-common/app/service/main/ugcpay-rank/internal/dao"
	"go-common/app/service/main/ugcpay-rank/internal/model"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"
)

// Service struct
type Service struct {
	Dao              *dao.Dao
	ElecUserSettings atomic.Value
	Asyncer          *fanout.Fanout
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		Dao:     dao.New(),
		Asyncer: fanout.New("async_worker", fanout.Worker(10), fanout.Buffer(10240)),
	}
	s.ElecUserSettings.Store(&sync.Map{})
	go s.reloadproc()
	return
}

func (s *Service) reloadproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("reloadproc panic: %+v\n%s", x, debug.Stack())
			s.reloadproc()
		}
	}()
	log.Info("reloadproc start on every %s", time.Duration(conf.Conf.Biz.ReloadDuration))
	ticker := time.NewTicker(time.Duration(conf.Conf.Biz.ReloadDuration))
	for {
		log.Info("reloadproc reload")
		var (
			limit    = 1000
			id       = 0
			m        = map[int64]model.ElecUserSetting{-1: -1}
			err      error
			count    = 0
			settings = &sync.Map{}
		)
		for len(m) > 0 {
			if m, id, err = s.Dao.RawElecUserSettings(_ctx, id, limit); err != nil {
				log.Error("s.Dao.RawElecUserSettings err: %+v", err)
				break
			}
			for mid, setting := range m {
				count++
				settings.Store(mid, setting)
			}
		}
		s.ElecUserSettings.Store(settings)
		log.Info("reloadproc reload end, count: %d", count)
		<-ticker.C
	}
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.Dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.Dao.Close()
}
