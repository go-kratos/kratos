package share

import (
	"context"

	"go-common/app/interface/main/web-goblin/conf"
	"go-common/app/interface/main/web-goblin/dao/share"
	accrpc "go-common/app/service/main/account/rpc/client"
	suitrpc "go-common/app/service/main/usersuit/rpc/client"
	"go-common/library/cache"
	"go-common/library/log"
)

// Service service struct.
type Service struct {
	c   *conf.Config
	dao *share.Dao
	// cache proc
	cache    *cache.Cache
	suit     *suitrpc.Service2
	accRPC   *accrpc.Service3
	Pendants map[int64]int64
}

// New new service.
func New(c *conf.Config) *Service {
	s := &Service{
		c:        c,
		dao:      share.New(c),
		cache:    cache.New(1, 1024),
		suit:     suitrpc.New(c.SuitRPC),
		accRPC:   accrpc.New3(c.AccountRPC),
		Pendants: make(map[int64]int64),
	}
	s.loadPendant()
	return s
}

// Ping ping service.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.Ping(c); err != nil {
		log.Error("s.dao.Ping error(%v)", err)
	}
	return
}
