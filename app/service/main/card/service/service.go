package service

import (
	"context"

	"go-common/app/service/main/card/conf"
	"go-common/app/service/main/card/dao"
	"go-common/app/service/main/card/model"
	viprpc "go-common/app/service/main/vip/rpc/client"
)

// Service struct
type Service struct {
	c            *conf.Config
	dao          *dao.Dao
	cardmap      map[int64]*model.Card
	cardgidmap   map[int64][]*model.Card
	cardhots     []*model.Card
	cardgroupmap map[int64]*model.CardGroup
	// vip rpc
	vipRPC *viprpc.Service
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
		// cache chan
		cardmap:      make(map[int64]*model.Card),
		cardgidmap:   make(map[int64][]*model.Card),
		cardgroupmap: make(map[int64]*model.CardGroup),
		cardhots:     []*model.Card{},
		vipRPC:       viprpc.New(c.RPCClient2.Vip),
	}
	if err := s.loadGroup(); err != nil {
		panic(err)
	}
	if err := s.loadCard(); err != nil {
		panic(err)
	}
	go s.loadcardproc()
	go s.loadcardgroupproc()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
