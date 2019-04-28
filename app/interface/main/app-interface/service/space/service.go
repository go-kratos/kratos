package space

import (
	"context"
	"runtime"
	"time"

	"go-common/app/interface/main/app-interface/conf"
	accdao "go-common/app/interface/main/app-interface/dao/account"
	arcdao "go-common/app/interface/main/app-interface/dao/archive"
	artdao "go-common/app/interface/main/app-interface/dao/article"
	audiodao "go-common/app/interface/main/app-interface/dao/audio"
	bgmdao "go-common/app/interface/main/app-interface/dao/bangumi"
	bplusdao "go-common/app/interface/main/app-interface/dao/bplus"
	coindao "go-common/app/interface/main/app-interface/dao/coin"
	commdao "go-common/app/interface/main/app-interface/dao/community"
	elecdao "go-common/app/interface/main/app-interface/dao/elec"
	favdao "go-common/app/interface/main/app-interface/dao/favorite"
	livedao "go-common/app/interface/main/app-interface/dao/live"
	memberdao "go-common/app/interface/main/app-interface/dao/member"
	paydao "go-common/app/interface/main/app-interface/dao/pay"
	reldao "go-common/app/interface/main/app-interface/dao/relation"
	srchdao "go-common/app/interface/main/app-interface/dao/search"
	shopdao "go-common/app/interface/main/app-interface/dao/shop"
	spcdao "go-common/app/interface/main/app-interface/dao/space"
	tagdao "go-common/app/interface/main/app-interface/dao/tag"
	thumbupdao "go-common/app/interface/main/app-interface/dao/thumbup"
	"go-common/library/log"
)

// Service is space service
type Service struct {
	c          *conf.Config
	arcDao     *arcdao.Dao
	spcDao     *spcdao.Dao
	accDao     *accdao.Dao
	coinDao    *coindao.Dao
	commDao    *commdao.Dao
	srchDao    *srchdao.Dao
	favDao     *favdao.Dao
	bgmDao     *bgmdao.Dao
	tagDao     *tagdao.Dao
	liveDao    *livedao.Dao
	elecDao    *elecdao.Dao
	artDao     *artdao.Dao
	audioDao   *audiodao.Dao
	relDao     *reldao.Dao
	bplusDao   *bplusdao.Dao
	shopDao    *shopdao.Dao
	thumbupDao *thumbupdao.Dao
	payDao     *paydao.Dao
	memberDao  *memberdao.Dao
	// chan
	mCh       chan func()
	tick      time.Duration
	BlackList map[int64]struct{}
}

// New new space
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		arcDao:     arcdao.New(c),
		spcDao:     spcdao.New(c),
		accDao:     accdao.New(c),
		coinDao:    coindao.New(c),
		commDao:    commdao.New(c),
		srchDao:    srchdao.New(c),
		favDao:     favdao.New(c),
		bgmDao:     bgmdao.New(c),
		tagDao:     tagdao.New(c),
		liveDao:    livedao.New(c),
		elecDao:    elecdao.New(c),
		artDao:     artdao.New(c),
		audioDao:   audiodao.New(c),
		relDao:     reldao.New(c),
		bplusDao:   bplusdao.New(c),
		shopDao:    shopdao.New(c),
		thumbupDao: thumbupdao.New(c),
		payDao:     paydao.New(c),
		memberDao:  memberdao.New(c),
		// mc proc
		mCh:       make(chan func(), 1024),
		tick:      time.Duration(c.Tick),
		BlackList: make(map[int64]struct{}),
	}
	// video db
	for i := 0; i < runtime.NumCPU(); i++ {
		go s.cacheproc()
	}
	if c != nil && c.Space != nil {
		for _, mid := range c.Space.ForbidMid {
			s.BlackList[mid] = struct{}{}
		}
	}
	s.loadBlacklist()
	go s.blacklistproc()
	return
}

// addCache add archive to mc or redis
func (s *Service) addCache(f func()) {
	select {
	case s.mCh <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc write memcache and stat redis use goroutine
func (s *Service) cacheproc() {
	for {
		f := <-s.mCh
		f()
	}
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// loadBlacklist
func (s *Service) loadBlacklist() {
	list, err := s.spcDao.Blacklist(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.BlackList = list
}

func (s *Service) blacklistproc() {
	for {
		time.Sleep(s.tick)
		s.loadBlacklist()
	}
}
