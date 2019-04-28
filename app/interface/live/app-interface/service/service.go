package service

import (
	"context"
	"go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
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

// Test ...
func (s *Service) Test(c context.Context) (err error) {
	// srv := &v2.IndexService{}
	// res, err := srv.GetIndexV2TagList(c, &liveUserV1.UserSettingGetTagReq{})
	// res, err := srv.GetIndexV2SeaPatrol(c, &liveUserV1.NoteGetReq{})
	// fmt.Printf("%#v  \n", res)
	return
}