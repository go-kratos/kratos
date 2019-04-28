package service

import (
	"context"

	"go-common/app/admin/main/esports/conf"
	"go-common/app/admin/main/esports/dao"
	accclient "go-common/app/service/main/account/api"
	accwarden "go-common/app/service/main/account/api"
	arcclient "go-common/app/service/main/archive/api"
)

// Service biz service def.
type Service struct {
	c         *conf.Config
	dao       *dao.Dao
	arcClient arcclient.ArchiveClient
	accClient accwarden.AccountClient
}

const (
	_notDeleted = 0
	_deleted    = 1
	_online     = 1
	_downLine   = 0
	_statusOn   = 0
	_statusAll  = -1
)

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	var err error
	if s.arcClient, err = arcclient.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	if s.accClient, err = accclient.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	return s
}

// Ping .
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

func unique(ids []int64) (outs []int64) {
	idMap := make(map[int64]int64, len(ids))
	for _, v := range ids {
		if _, ok := idMap[v]; ok {
			continue
		} else {
			idMap[v] = v
		}
		outs = append(outs, v)
	}
	return
}
