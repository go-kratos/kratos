package service

import (
	"context"
	"time"

	"go-common/app/admin/main/point/conf"
	"go-common/app/admin/main/point/dao"
	"go-common/app/admin/main/point/model"
	pointrpc "go-common/app/service/main/point/rpc/client"
	"go-common/library/log"
)

// Service struct
type Service struct {
	c        *conf.Config
	dao      *dao.Dao
	pointRPC *pointrpc.Service
	appMap   map[int64]string
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      dao.New(c),
		pointRPC: pointrpc.New(c.RPC.Point),
		appMap:   make(map[int64]string),
	}
	go s.loadappinfoproc()
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

func (s *Service) loadappinfoproc() {

	for {
		s.loadAppInfo()
		time.Sleep(time.Minute * 2)
	}
}

func (s *Service) loadAppInfo() {
	var (
		res []*model.AppInfo
		err error
	)
	if res, err = s.dao.AllAppInfo(context.TODO()); err != nil {
		log.Error("loadAppInfo AllAppInfo error(%v)", err)
		return
	}
	aMap := make(map[int64]string, len(res))
	for _, v := range res {
		aMap[v.ID] = v.Name
	}
	s.appMap = aMap
}
