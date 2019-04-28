package show

import (
	"context"
	"time"

	"go-common/app/job/main/app/conf"
	showdao "go-common/app/job/main/app/dao/show"
	"go-common/app/job/main/app/model"
	"go-common/library/log"
)

// Service is show service.
type Service struct {
	c    *conf.Config
	dao  *showdao.Dao
	tick time.Duration
}

// New new a show service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:    c,
		dao:  showdao.New(c),
		tick: time.Duration(c.Tick),
	}
	if model.EnvRun() {
		s.pub(time.Now())
		go s.loadproc()
	}
	return
}

// pub publish show data by timer.
func (s *Service) pub(now time.Time) {
	c := context.Background()
	ps, err := s.dao.PTime(c, now)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(ps) == 0 {
		return
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	for _, p := range ps {
		if err = s.dao.Pub(tx, p); err != nil {
			tx.Rollback()
			log.Error("%+v", err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("%+v", err)
		return
	}
	log.Info("show publish success plats(%v)", ps)
}

// cacheproc load all cache.
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.pub(time.Now())
	}
}

func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.PingDB(c)
	return
}
