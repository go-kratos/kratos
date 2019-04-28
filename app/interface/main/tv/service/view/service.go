package view

import (
	"context"

	"go-common/app/interface/main/tv/conf"
	"go-common/app/interface/main/tv/dao/account"
	"go-common/app/interface/main/tv/dao/archive"
	"go-common/app/interface/main/tv/dao/cms"
	"go-common/app/interface/main/tv/dao/favorite"
	"go-common/app/interface/main/tv/dao/upper"
	"go-common/library/stat/prom"
)

// Service .
type Service struct {
	conf *conf.Config
	// dao
	arcDao *archive.Dao
	accDao *account.Dao
	cmsDao *cms.Dao
	upDao  *upper.Dao
	favDao *favorite.Dao
	// prom
	pHit       *prom.Prom
	pMiss      *prom.Prom
	prom       *prom.Prom
	emptyArcCh chan int64
}

var ctx = context.TODO()

// New .
func New(c *conf.Config) *Service {
	srv := &Service{
		conf:       c,
		arcDao:     archive.New(c),
		accDao:     account.New(c),
		cmsDao:     cms.New(c),
		upDao:      upper.New(c),
		favDao:     favorite.New(c),
		pHit:       prom.CacheHit,
		pMiss:      prom.CacheMiss,
		prom:       prom.BusinessInfoCount,
		emptyArcCh: make(chan int64, c.Cfg.EmptyArc.ChanSize),
	}
	go srv.emptyArcproc()
	return srv
}
