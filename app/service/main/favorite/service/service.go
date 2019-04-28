package service

import (
	"context"
	"time"

	accapi "go-common/app/service/main/account/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/favorite/conf"
	"go-common/app/service/main/favorite/dao"
	fltmdl "go-common/app/service/main/filter/model/rpc"
	fltrpc "go-common/app/service/main/filter/rpc/client"
	rankrpc "go-common/app/service/main/rank/api/gorpc"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

// Service define fav service
type Service struct {
	conf        *conf.Config
	missch      chan interface{}
	cleanCDTime int64
	httpClient  *httpx.Client
	favDao      *dao.Dao
	// cache chan
	cache *fanout.Fanout
	// prom
	prom *prom.Prom
	// rpc
	rankRPC   *rankrpc.Service
	filterRPC *fltrpc.Service
	accClient accapi.AccountClient
	arcRPC    *arcrpc.Service2
}

// New return fav service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:        c,
		missch:      make(chan interface{}, 1024),
		cleanCDTime: int64(time.Duration(c.Fav.CleanCDTime) / time.Second),
		httpClient:  httpx.NewClient(c.HTTPClient),
		// dao
		favDao: dao.New(c),
		// cache
		cache: fanout.New("cache"),
		// prom
		prom: prom.New().WithTimer("fav_add_video", []string{"method"}),
		// rpc
		filterRPC: fltrpc.New(c.RPCClient2.Filter),
		rankRPC:   rankrpc.New(c.RPCClient2.Rank),
		arcRPC:    arcrpc.New2(c.RPCClient2.Archive),
	}
	var err error
	if s.accClient, err = accapi.NewClient(c.RPCClient2.Account); err != nil {
		panic(err)
	}
	if s.conf.Fav.MaxParallelSize == 0 {
		s.conf.Fav.MaxParallelSize = s.conf.Fav.DefaultFolderLimit
	}
	return
}

// Ping check service health
func (s *Service) Ping(c context.Context) (err error) {
	return s.favDao.Ping(c)
}

// Close close service
func (s *Service) Close() {
	s.favDao.Close()
}

// PromError stat and log.
func (s *Service) PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}

// Filter filter folder name
func (s *Service) filter(c context.Context, name string) (err error) {
	var (
		res *fltmdl.FilterRes
	)
	arg := &fltmdl.ArgFilter{
		Area:    "open_medialist",
		Message: name,
		TypeID:  0,
		ID:      0,
	}
	if res, err = s.filterRPC.Filter(c, arg); err != nil {
		log.Error("s.filterRPC.Filter(%s) error(%v)", name, err)
		return
	}
	if res.Level >= 20 {
		err = ecode.FavFolderBanned
	} else if res.Level > 0 {
		err = ecode.FavHitSensitive
	}
	return
}

func (s *Service) checkUser(c context.Context, mid int64) (err error) {
	profileReply, err := s.accClient.Profile3(c, &accapi.MidReq{Mid: mid})
	if err != nil {
		log.Error("s.accClient.Profile3(%d) error(%v)", mid, err)
		return nil
	}
	profile := profileReply.Profile
	if profile.Identification == 0 && profile.TelStatus == 0 {
		err = ecode.UserCheckNoPhone
		return
	}
	if profile.Identification == 0 && profile.TelStatus == 2 {
		err = ecode.UserCheckInvalidPhone
		return
	}
	if profile.EmailStatus == 0 && profile.TelStatus == 0 {
		err = ecode.UserInactive
		return
	}
	if profile.Silence == 1 {
		err = ecode.UserDisabled
	}
	return
}

func (s *Service) checkRealname(c context.Context, mid int64) (err error) {
	profileReply, err := s.accClient.Profile3(c, &accapi.MidReq{Mid: mid})
	if err != nil {
		log.Error("s.accClient.Profile3(%d) error(%v)", mid, err)
		return nil
	}
	profile := profileReply.Profile
	if profile.Identification == 0 && profile.TelStatus == 0 {
		err = ecode.UserCheckNoPhone
		return
	}
	if profile.Identification == 0 && profile.TelStatus == 2 {
		err = ecode.UserCheckInvalidPhone
	}
	return
}
