package web

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go-common/app/job/main/web-goblin/conf"
	mdlweb "go-common/app/job/main/web-goblin/dao/web"
	"go-common/app/job/main/web-goblin/model/web"
	dymdl "go-common/app/service/main/dynamic/model"
	dyrpc "go-common/app/service/main/dynamic/rpc/client"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service struct .
type Service struct {
	c                *conf.Config
	dao              *mdlweb.Dao
	dy               *dyrpc.Service
	waiter           *sync.WaitGroup
	archiveNotifySub *databus.Databus
}

// New init .
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                c,
		dao:              mdlweb.New(c),
		dy:               dyrpc.New(c.DynamicRPC),
		waiter:           new(sync.WaitGroup),
		archiveNotifySub: databus.New(c.ArchiveNotifySub),
	}
	go s.broadcastDy()
	s.waiter.Add(1)
	go s.allSearch()
	return s
}

// Ping Service .
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service .
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) broadcastDy() {
	var (
		dynamics map[string]int
		err      error
		b        []byte
	)
	for {
		if dynamics, err = s.dy.RegionTotal(context.Background(), &dymdl.ArgRegionTotal{RealIP: ""}); err != nil {
			mdlweb.PromError("RegionTotal接口错误", "s.dy.RegionTotal error(%v)", err)
			time.Sleep(time.Second)
			continue
		}
		if b, err = json.Marshal(dynamics); err != nil {
			log.Error("broadcastDy json.Marshal error(%v)", err)
			return
		}
		if err = s.dao.PushAll(context.Background(), string(b), ""); err != nil {
			log.Error("s.dao.PushAll(%+v) error(%v)", dynamics, err)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second * 5)
	}
}

func (s *Service) allSearch() {
	var (
		err error
		ctx = context.Background()
	)

	defer s.waiter.Done()
	for {
		msg, ok := <-s.archiveNotifySub.Messages()
		if !ok {
			log.Error("web-goblin-job databus Consumer exit")
			return
		}
		res := &web.ArcMsg{}
		if err = json.Unmarshal(msg.Value, res); err != nil {
			msg.Commit()
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if res.Table != _archive {
			continue
		}
		s.UgcIncrement(ctx, res)
		msg.Commit()
		log.Info("consume allSearch ugc key:%s partition:%d offset:%d)", msg.Key, msg.Partition, msg.Offset)
		time.Sleep(10 * time.Millisecond)
	}
}
