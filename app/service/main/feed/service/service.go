package service

import (
	"context"

	artmdl "go-common/app/interface/openplatform/article/model"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	account "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/feed/conf"
	"go-common/app/service/main/feed/dao"
	feedmdl "go-common/app/service/main/feed/model"
)

// Service struct info.
type Service struct {
	c       *conf.Config
	dao     *dao.Dao
	arcRPC  ArcRPC
	accRPC  AccRPC
	artRPC  ArtRPC
	bangumi Bangumi
	missch  chan func()
}

//go:generate mockgen -source=./service.go -destination=mock_test.go -package=service
type ArcRPC interface {
	Archive3(c context.Context, arg *archive.ArgAid2) (res *api.Arc, err error)
	Archives3(c context.Context, arg *archive.ArgAids2) (res map[int64]*api.Arc, err error)
	UpsPassed2(c context.Context, arg *archive.ArgUpsArcs2) (res map[int64][]*archive.AidPubTime, err error)
}

type AccRPC interface {
	Attentions3(c context.Context, arg *account.ArgMid) (res []int64, err error)
}
type ArtRPC interface {
	UpsArtMetas(c context.Context, arg *artmdl.ArgUpsArts) (res map[int64][]*artmdl.Meta, err error)
	ArticleMetas(c context.Context, arg *artmdl.ArgAids) (res map[int64]*artmdl.Meta, err error)
}

type Bangumi interface {
	BangumiPull(c context.Context, mid int64, ip string) (seasonIDS []int64, err error)
	BangumiSeasons(c context.Context, seasonIDs []int64, ip string) (psm map[int64]*feedmdl.Bangumi, err error)
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	d := dao.New(c)
	s = &Service{
		c:       c,
		dao:     d,
		bangumi: d,
		arcRPC:  arcrpc.New2(c.ArchiveRPC),
		accRPC:  accrpc.New3(c.AccountRPC),
		artRPC:  artrpc.New(c.ArticleRPC),
		missch:  make(chan func(), 1000),
	}
	go s.cacheproc()
	return
}

func (s *Service) addCache(fn func()) {
	select {
	case s.missch <- fn:
	default:
		dao.PromError("cache队列已满", "cacheproc chan full!!!")
	}
}

func (s *Service) cacheproc() {
	for i := 0; i < 10; i++ {
		go func() {
			for {
				fn := <-s.missch
				fn()
			}
		}()
	}
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
