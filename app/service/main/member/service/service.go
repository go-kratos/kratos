package service

import (
	"context"

	"go-common/app/service/main/member/conf"
	mbDao "go-common/app/service/main/member/dao"
	"go-common/app/service/main/member/model"
	"go-common/app/service/main/member/service/block"
	"go-common/app/service/main/member/service/crypto"
	"go-common/library/sync/pipeline/fanout"
)

// Service struct of service.
type Service struct {
	c               *conf.Config
	mbDao           *mbDao.Dao
	cache           *fanout.Fanout
	officials       map[int64]*model.OfficialInfo
	block           *block.Service
	realnameCryptor *crypto.Realname
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:               c,
		mbDao:           mbDao.New(c),
		cache:           fanout.New("cache", fanout.Worker(1), fanout.Buffer(10240)),
		realnameCryptor: crypto.NewRealname(conf.RsaPub(), conf.RsaPriv()),
	}
	s.block = block.New(c, s.mbDao.BlockImpl())
	if err := s.loadOfficial(); err != nil {
		panic(err)
	}
	go s.loadOfficialproc()
	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.mbDao.Ping(c)
}

// Close dao.
func (s *Service) Close() {
	s.mbDao.Close()
	s.block.Close()
}

// BlockImpl is
func (s *Service) BlockImpl() *block.Service {
	return s.block
}
