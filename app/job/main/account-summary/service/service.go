package service

import (
	"context"

	"go-common/app/job/main/account-summary/conf"
	"go-common/app/job/main/account-summary/dao"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao

	MemberBinLog           *databus.Databus
	BlockBinLog            *databus.Databus
	PassportBinLog         *databus.Databus
	RelationBinLog         *databus.Databus
	AccountSummaryProducer *databus.Databus
}

// New init
func New(c *conf.Config) *Service {
	s := &Service{
		c:                      c,
		dao:                    dao.New(c),
		RelationBinLog:         databus.New(c.RelationBinLog),
		MemberBinLog:           databus.New(c.MemberBinLog),
		BlockBinLog:            databus.New(c.BlockBinLog),
		PassportBinLog:         databus.New(c.PassportBinLog),
		AccountSummaryProducer: databus.New(c.AccountSummaryProducer),
	}
	s.Main()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// Main is
func (s *Service) Main() {
	subproc := func() {
		worker := s.c.AccountSummary.SubProcessWorker
		if worker <= 0 {
			worker = 1
		}
		log.Info("Starting sub process with %d workers", worker)
		for i := uint64(0); i < worker; i++ {
			go s.memberBinLogproc(context.Background())
			go s.blockBinLogproc(context.Background())
			go s.passportBinLogproc(context.Background())
			go s.relationBinLogproc(context.Background())
		}
	}

	syncrange := func() {
		start := s.c.AccountSummary.SyncRangeStart
		if start <= 0 {
			start = 1
		}
		end := s.c.AccountSummary.SyncRangeEnd
		if end <= 0 {
			end = 1
		}
		worker := s.c.AccountSummary.SyncRangeWorker
		if worker <= 0 {
			worker = 1
		}
		go s.syncRangeproc(context.Background(), start, end, worker)
	}

	// initial := func() {
	// go s.initialproc(context.Background())
	// }

	if !s.c.FeatureGate.DisableSubProcess {
		subproc()
	}

	if s.c.FeatureGate.SyncRange {
		syncrange()
	}

	// if s.c.FeatureGate.Initial {
	// 	initial()
	// }
}
