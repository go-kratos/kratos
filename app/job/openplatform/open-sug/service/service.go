package service

import (
	"context"

	"go-common/app/job/openplatform/open-sug/conf"
	"go-common/app/job/openplatform/open-sug/dao"
	"go-common/library/queue/databus"
	"sync"
	"time"
)

// Service struct
type Service struct {
	c            *conf.Config
	dao          *dao.Dao
	pgcSub       *databus.Databus
	pgcMsgCnt    int64
	wg           *sync.WaitGroup
	seasonMsgCnt int64
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		dao:          dao.New(c),
		wg:           new(sync.WaitGroup),
		pgcSub:       databus.New(c.PgcSub),
		pgcMsgCnt:    0,
		seasonMsgCnt: 0,
	}
	s.existsOrCreate(c.ElasticSearch.Season)
	go s.pgcConsumePROC()
	go s.fetchData()

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

func (s *Service) fetchData() {
	for {
		s.buildData()
		s.Refresh()
		s.Filter()
		time.Sleep(time.Hour * 24)
	}
}

// Refresh ...
func (s *Service) Refresh() {
	s.dao.ItemSalesMin = make(map[string]int)
	s.dao.ItemSalesMax = make(map[string]int)
	s.dao.ItemWishMax = make(map[string]int)
	s.dao.ItemWishMin = make(map[string]int)
	s.dao.ItemCommentMax = make(map[string]int)
	s.dao.ItemCommentMin = make(map[string]int)
}
