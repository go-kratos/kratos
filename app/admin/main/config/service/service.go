package service

import (
	"sync"

	"go-common/app/admin/main/config/conf"
	"go-common/app/admin/main/config/dao"
	confrpc "go-common/app/infra/config/rpc/client"

	"go-common/app/admin/main/config/model"

	"github.com/jinzhu/gorm"
)

// Service service
type Service struct {
	c *conf.Config

	// rpcconf config service Rpc
	confSvr *confrpc.Service2
	dao     *dao.Dao

	cLock sync.RWMutex
	cache map[string]*model.CacheData
	//apm gorm
	DBApm *gorm.DB
	//db gorm
	DB *gorm.DB
}

// New new a service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		confSvr: confrpc.New2(c.ConfSvr),
		dao:     dao.New(c),
	}
	s.cache = make(map[string]*model.CacheData)
	s.DBApm = s.dao.DBApm
	s.DB = s.dao.DB
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
