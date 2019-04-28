package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/reply/conf"
	"go-common/app/admin/main/reply/dao"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	accrpc "go-common/app/service/main/account/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	rlrpc "go-common/app/service/main/relation/rpc/client"
	thumbup "go-common/app/service/main/thumbup/api"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/sync/pipeline/fanout"
)

// Service is a service.
type Service struct {
	conf          *conf.Config
	dao           *dao.Dao
	httpClient    *bm.Client
	accSrv        accrpc.AccountClient
	arcSrv        *arcrpc.Service2
	articleSrv    *artrpc.Service
	cache         *fanout.Fanout
	thumbupClient thumbup.ThumbupClient

	relationSvc *rlrpc.Service

	// 特殊admin，针对大忽悠事件消息通知
	ads map[string]struct{}
	// 特殊稿件，针对大忽悠时间 不让删评论
	oids map[int64]int32

	// mark folded or unmark folded worker
	marker *fanout.Fanout
	// del cache or add cache worker
	cacheOperater *fanout.Fanout
}

// New new a service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:          c,
		dao:           dao.New(c),
		httpClient:    bm.NewClient(c.HTTPClient),
		arcSrv:        arcrpc.New2(c.RPCClient2.Archive),
		articleSrv:    artrpc.New(c.RPCClient2.Article),
		cache:         fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024*10)),
		relationSvc:   rlrpc.New(c.RPCClient2.Relation),
		marker:        fanout.New("marker", fanout.Worker(1), fanout.Buffer(1024*10)),
		cacheOperater: fanout.New("cacheOp", fanout.Worker(1), fanout.Buffer(1024*10)),
	}
	accSrv, err := accrpc.NewClient(c.AccountClient)
	if err != nil {
		panic(err)
	}
	s.accSrv = accSrv
	cc, err := warden.NewConn(fmt.Sprintf("discovery://default/%s", thumbup.AppID))
	if err != nil {
		panic(err)
	}
	s.thumbupClient = thumbup.NewThumbupClient(cc)
	ads := make(map[string]struct{})
	for _, ID := range c.Reply.AdminName {
		ads[ID] = struct{}{}
	}
	s.ads = ads
	oids := make(map[int64]int32)
	for i, oid := range c.Reply.Oids {
		oids[oid] = c.Reply.Tps[i]
	}
	s.oids = oids
	return
}

// Ping check service is ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}
