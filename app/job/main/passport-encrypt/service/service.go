package service

import (
	"context"
	"sync"

	"go-common/app/job/main/passport-encrypt/conf"
	"go-common/app/job/main/passport-encrypt/dao"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	// table and duration
	_asoAccountTable = "aso_account"
	_insertAction    = "insert"
	_updateAction    = "update"
	_deleteAction    = "delete"
)

// Service service.
type Service struct {
	c *conf.Config
	d *dao.Dao
	// aso binlog databus
	dsAsoBinLogSub *databus.Databus
	merges         []chan *message
	done           chan []*message
	// proc
	head, last *message
	mu         sync.Mutex
}

type message struct {
	next   *message
	data   *databus.Message
	object interface{}
	done   bool
}

// New new a service instance.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		d:              dao.New(c),
		dsAsoBinLogSub: databus.New(c.DataBus.AsoBinLogSub),
		merges:         make([]chan *message, c.Group.AsoBinLog.Num),
		done:           make(chan []*message, c.Group.AsoBinLog.Chan),
	}
	if c.DataSwitch.Full {
		go s.fullMigration(c.StepGroup.Group1.Start, c.StepGroup.Group1.End, c.StepGroup.Group1.Inc, c.StepGroup.Group1.Limit, "group1")
		go s.fullMigration(c.StepGroup.Group2.Start, c.StepGroup.Group2.End, c.StepGroup.Group2.Inc, c.StepGroup.Group2.Limit, "group2")
		go s.fullMigration(c.StepGroup.Group3.Start, c.StepGroup.Group3.End, c.StepGroup.Group3.Inc, c.StepGroup.Group3.Limit, "group3")
		go s.fullMigration(c.StepGroup.Group4.Start, c.StepGroup.Group4.End, c.StepGroup.Group4.Inc, c.StepGroup.Group4.Limit, "group4")
	}
	if c.DataSwitch.Inc {
		go s.asobinlogcommitproc()
		for i := 0; i < c.Group.AsoBinLog.Num; i++ {
			ch := make(chan *message, c.Group.AsoBinLog.Chan)
			s.merges[i] = ch
			go s.asobinlogmergeproc(ch)
		}
		go s.asobinlogconsumeproc()
	}

	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Close close service, including databus and outer service.
func (s *Service) Close() (err error) {
	if err = s.dsAsoBinLogSub.Close(); err != nil {
		log.Error("srv.asoBinLog.Close() error(%v)", err)
	}
	return
}
