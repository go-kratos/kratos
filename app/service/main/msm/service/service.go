package service

import (
	"container/list"
	"sync"
	"sync/atomic"

	confrpc "go-common/app/infra/config/rpc/client"
	"go-common/app/service/main/msm/conf"
	"go-common/app/service/main/msm/dao"
	"go-common/app/service/main/msm/model"
)

const (
	_maxVerNum = 100
)

// Service service
type Service struct {
	c *conf.Config

	// rpcconf config service Rpc
	confSvr *confrpc.Service2
	dao     *dao.Dao

	// ecode
	lock     sync.RWMutex
	version  *model.Version
	codes    atomic.Value
	scopeMap map[int64]map[int64]*model.Scope
	msmScope map[int64]*model.Scope

	// langs
	langsLock    sync.RWMutex
	langsVersion *model.Version
	langsCodes   atomic.Value
}

// New new a service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		confSvr:      confrpc.New2(c.ConfSvr),
		dao:          dao.New(c),
		version:      &model.Version{List: list.New(), Map: make(map[int64]*list.Element)},
		scopeMap:     make(map[int64]map[int64]*model.Scope),
		msmScope:     make(map[int64]*model.Scope),
		langsVersion: &model.Version{List: list.New(), Map: make(map[int64]*list.Element)},
	}
	if err := s.all(); err != nil {
		panic(err)
	}
	if err := s.allLang(); err != nil {
		panic(err)
	}
	if err := s.updateScope(); err != nil {
		panic(err)
	}
	if err := s.updateMsmScope(); err != nil {
		panic(err)
	}
	go s.updateLangproc()
	go s.updateproc()
	go s.updateScopeproc()
	go s.updateMsmScopeproc()
	return
}

// Ping check server ok.
func (s *Service) Ping() (err error) {
	return
}

// Close close resource
func (s *Service) Close() {
	s.dao.Close()
}
