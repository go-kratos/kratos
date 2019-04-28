package archive

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/activity"
	"go-common/app/interface/main/creative/dao/appeal"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/coin"
	gD "go-common/app/interface/main/creative/dao/game"
	"go-common/app/interface/main/creative/dao/order"
	"go-common/app/interface/main/creative/dao/search"
	"go-common/app/interface/main/creative/dao/tag"
	"go-common/app/interface/main/creative/dao/template"
	actmdl "go-common/app/interface/main/creative/model/activity"
	"go-common/app/interface/main/creative/model/game"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

//Service struct
type Service struct {
	c          *conf.Config
	arc        *archive.Dao
	acc        *account.Dao
	sear       *search.Dao
	act        *activity.Dao
	tpl        *template.Dao
	coin       *coin.Dao
	order      *order.Dao
	ap         *appeal.Dao
	tag        *tag.Dao
	game       *gD.Dao
	p          *service.Public
	prom       *prom.Prom
	missch     chan func()
	pCacheHit  *prom.Prom
	pCacheMiss *prom.Prom
	// cache
	orderUps map[int64]int64
	gameMap  map[int64]*game.ListItem
	ArcTip   string
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos, p *service.Public) *Service {
	s := &Service{
		c:          c,
		arc:        rpcdaos.Arc,
		acc:        rpcdaos.Acc,
		sear:       search.New(c),
		act:        activity.New(c),
		tpl:        template.New(c),
		tag:        tag.New(c),
		coin:       coin.New(c),
		order:      order.New(c),
		ap:         appeal.New(c),
		game:       gD.New(c),
		p:          p,
		prom:       prom.BusinessInfoCount,
		missch:     make(chan func(), 1024),
		pCacheHit:  prom.CacheHit,
		pCacheMiss: prom.CacheMiss,
		ArcTip:     c.Host.ArcTip,
	}
	s.loadOrderUps()
	s.loadAllGameMap()
	go s.loadproc()
	go s.cacheproc()
	return s
}

// TopAct  fn
func (s *Service) TopAct() (ret []*actmdl.Activity) {
	return s.p.TopActCache
}
func (s *Service) loadOrderUps() {
	orderUps, err := s.order.Ups(context.TODO())
	if err != nil {
		return
	}
	s.orderUps = orderUps
}

func (s *Service) loadAllGameMap() {
	list, err := s.game.List(context.TODO(), "", "")
	if err != nil || list == nil || len(list) == 0 {
		return
	}
	s.gameMap = make(map[int64]*game.ListItem)
	for _, v := range list {
		s.gameMap[v.GameBaseID] = v
	}
	log.Info("s.loadAllGameMap: s.gameMapLen(%d)", len(s.gameMap))
}

// loadproc
func (s *Service) loadproc() {
	for {
		time.Sleep(5 * time.Minute)
		s.loadOrderUps()
		s.loadAllGameMap()
	}
}

// AllowOrderUps 检查用户商单信息
func (s *Service) AllowOrderUps(mid int64) (ok bool) {
	_, ok = s.orderUps[mid]
	return
}

// AddCache add to chan for cache
func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for execute closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.arc.Ping(c); err != nil {
		log.Error("s.archive.Dao.PingDb err(%v)", err)
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.arc.Close()
}
