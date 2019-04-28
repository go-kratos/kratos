package service

import (
	"context"
	"time"

	"go-common/app/job/main/credit-timer/conf"
	"go-common/app/job/main/credit-timer/dao"
	"go-common/library/log"
)

// Service struct of service.
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	go s.loadConfproc()
	go s.caseproc()
	go s.juryproc()
	go s.voteproc()
	go s.kpiproc()
	return
}

func (s *Service) loadConfproc() {
	for {
		s.loadConf(context.TODO())
		time.Sleep(time.Duration(s.c.Judge.ConfTimer))
	}
}
func (s *Service) caseproc() {
	for {
		s.caseProc(context.TODO())
		time.Sleep(time.Duration(s.c.Judge.CaseTimer))
	}
}

func (s *Service) juryproc() {
	for {
		s.juryProc(context.TODO())
		time.Sleep(time.Duration(s.c.Judge.JuryTimer))
	}
}

func (s *Service) voteproc() {
	for {
		s.voteProc(context.TODO())
		time.Sleep(time.Duration(s.c.Judge.VoteTimer))
	}
}

func (s *Service) kpiproc() {
	var err error
	for {
		d := time.Now().AddDate(0, 0, 1)
		ts := time.Until(time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 1, 0, time.Local))
		time.Sleep(ts)
		for {
			err = s.kpiPointProc(context.TODO())
			if err != nil {
				log.Error("kpiPointProc err(%v)", err)
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
		log.Info("KPIPointproc err(%v)", err)
		for {
			err = s.KPIProc(context.TODO())
			if err != nil {
				log.Error("kpiProc err(%v)", err)
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
		log.Info("kpiproc err(%v)", err)
	}
}

// Close kafka consumer close.
func (s *Service) Close() (err error) {
	return
}

// Ping check service health.
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}
