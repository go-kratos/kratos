package service

import (
	"context"
	"sync"

	"go-common/app/job/main/passport-game-local/conf"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	// table and duration
	_asoAccountTable = "aso_account"
)

// Service service.
type Service struct {
	c *conf.Config
	// aso encrypt trans pub databus
	dsAsoEncryptTransPub *databus.Databus
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
		c:                    c,
		dsAsoEncryptTransPub: databus.New(c.DataBus.EncryptTransPub),
		dsAsoBinLogSub:       databus.New(c.DataBus.AsoBinLogSub),
		merges:               make([]chan *message, c.Group.AsoBinLog.Num),
		done:                 make(chan []*message, c.Group.AsoBinLog.Chan),
	}
	go s.asobinlogcommitproc()
	for i := 0; i < c.Group.AsoBinLog.Num; i++ {
		ch := make(chan *message, c.Group.AsoBinLog.Chan)
		s.merges[i] = ch
		go s.asobinlogmergeproc(ch)
	}
	go s.asobinlogconsumeproc()
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
	if err = s.dsAsoEncryptTransPub.Close(); err != nil {
		log.Error("srv.dsAsoEncryptTransPub.Close() error(%v)", err)
	}
	return
}
