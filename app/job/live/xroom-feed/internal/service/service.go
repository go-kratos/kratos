package service

import (
	"context"
	"sync/atomic"
	"time"

	"go-common/library/log"
	"go-common/library/net/rpc/liverpc"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"

	"go-common/app/job/live/xroom-feed/internal/dao"
	daoAnchor "go-common/app/service/live/dao-anchor/api/grpc/v1"
	roomClient "go-common/app/service/live/room/api/liverpc"

	"go-common/library/conf/paladin"
)

// Service service.
type Service struct {
	ac          *paladin.Map
	dao         *dao.Dao
	daoAnchor   daoAnchor.DaoAnchorClient
	roomService *roomClient.Client
	ruleConf    atomic.Value

	indexBlackList atomic.Value
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  ac,
		dao: dao.New(),
	}
	wdConf := new(warden.ClientConfig)
	wdConf.Timeout = xtime.Duration(time.Second * 10)
	err := s.ac.Get("daoAnchorClient").UnmarshalTOML(wdConf)
	if err != nil {
		log.Error("[service]get daoAnchorClient warden error:%+v", err)
		wdConf.Dial = xtime.Duration(time.Millisecond * 100)
		wdConf.Timeout = xtime.Duration(time.Second * 10)
	}
	conn, err := daoAnchor.NewClient(wdConf) // 目前传空，如果需要配置
	if err != nil {
		panic(err)
	}
	s.daoAnchor = conn

	roomClientConf := new(liverpc.ClientConfig)
	rerr := s.ac.Get("roomClient").UnmarshalTOML(roomClientConf)
	if rerr != nil {
		log.Error("[service]get roomClient conf error:%+v", rerr)
		roomClientConf.ConnTimeout = xtime.Duration(time.Millisecond * 50)
		roomClientConf.AppID = "live.room"
	}
	s.roomService = roomClient.New(roomClientConf)

	s.loadConfFromDb()
	s.loadBlackList()
	go s.reloadConfFromDb()
	go s.reloadRecList()
	go s.blackListProc()
	return s
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
