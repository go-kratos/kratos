package wall

import (
	"context"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	walldao "go-common/app/interface/main/app-wall/dao/wall"
	"go-common/app/interface/main/app-wall/model/wall"
	log "go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

type Service struct {
	c         *conf.Config
	client    *httpx.Client
	dao       *walldao.Dao
	tick      time.Duration
	cache     []*wall.Wall
	testCache []*wall.Wall
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		client: httpx.NewClient(c.HTTPClient),
		dao:    walldao.New(c),
		tick:   time.Duration(c.Tick),
	}
	s.load()
	go s.loadproc()
	return
}

// GetWall All
func (s *Service) Wall() (res []*wall.Wall) {
	res = s.cache
	return
}

// load WallAll
func (s *Service) load() {
	res, err := s.dao.WallAll(context.TODO())
	if err != nil {
		log.Error("s.dao.wallAll error(%v)", err)
		return
	}
	s.cache = res
	s.testCache = res
	log.Info("loadWallsCache success")
}

// cacheproc load cache
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.load()
	}
}
