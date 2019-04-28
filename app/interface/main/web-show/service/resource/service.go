package resource

import (
	"context"
	"time"

	"go-common/app/interface/main/web-show/conf"
	"go-common/app/interface/main/web-show/dao/ad"
	"go-common/app/interface/main/web-show/dao/bangumi"
	"go-common/app/interface/main/web-show/dao/data"
	resdao "go-common/app/interface/main/web-show/dao/resource"
	rsmdl "go-common/app/interface/main/web-show/model/resource"
	accrc "go-common/app/service/main/account/rpc/client"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	locrpc "go-common/app/service/main/location/rpc/client"
	recrpc "go-common/app/service/main/resource/rpc/client"
	"go-common/library/log"
)

// Service define web-show service
type Service struct {
	resdao     *resdao.Dao
	bangumiDao *bangumi.Dao
	accRPC     *accrc.Service3
	arcRPC     *arcrpc.Service2
	recrpc     *recrpc.Service
	adDao      *ad.Dao
	dataDao    *data.Dao
	locRPC     *locrpc.Service
	// cache
	asgCache       map[int][]*rsmdl.Assignment // resID => assignments
	urlMonitor     map[int]map[string]string   // pf=>map[rs.name=>url]
	videoCache     map[int64][][]*rsmdl.VideoAD
	posCache       map[string]*rsmdl.Position // resID=>srcIDs
	defBannerCache *rsmdl.Assignment
	adsCache       map[int]*rsmdl.Assignment
}

// New return service object
func New(c *conf.Config) *Service {
	s := &Service{
		adDao:      ad.New(c),
		resdao:     resdao.New(c),
		bangumiDao: bangumi.New(c),
		dataDao:    data.New(c),
		asgCache:   make(map[int][]*rsmdl.Assignment),
		videoCache: make(map[int64][][]*rsmdl.VideoAD),
		posCache:   make(map[string]*rsmdl.Position),
		// crm
		adsCache: make(map[int]*rsmdl.Assignment),
	}
	s.arcRPC = arcrpc.New2(c.RPCClient2.Archive)
	s.accRPC = accrc.New3(c.RPCClient2.Account)
	s.recrpc = recrpc.New(c.RPCClient2.Resource)
	s.locRPC = locrpc.New(c.LocationRPC)
	s.init()
	return s
}

func (s *Service) init() (err error) {
	if err = s.loadRes(); err != nil {
		log.Error("adService.Load, err (%v)", err)
	}
	if err = s.loadVideoAd(); err != nil {
		log.Error("adService.LoadVideo, err (%v)", err)
	}
	// if err = s.loadAds(); err != nil {
	// 	log.Error("s.loadAds err(%v)", err)
	// }
	go s.checkproc()
	go s.loadproc()
	return
}

// loadproc is a routine load ads to cache
func (s *Service) loadproc() {
	for {
		s.loadRes()
		s.loadVideoAd()
		//s.loadAds()
		time.Sleep(time.Duration(conf.Conf.Reload.Ad))
	}
}

// checkpro a routine check diff
func (s *Service) checkproc() {
	for {
		s.checkDiff()
		time.Sleep(time.Duration(conf.Conf.Reload.Ad))
	}
}

// Close close service
func (s *Service) Close() {
	s.resdao.Close()
}

// Ping ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.resdao.Ping(c); err != nil {
		log.Error("s.resDap.Ping err(%v)", err)
		return
	}
	if err = s.adDao.Ping(c); err != nil {
		log.Error("s.adDao.Ping err(%v)", err)
	}
	return
}
