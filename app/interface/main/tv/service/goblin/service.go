package goblin

import (
	"context"

	"go-common/app/interface/main/tv/conf"
	"go-common/app/interface/main/tv/dao/account"
	"go-common/app/interface/main/tv/dao/archive"
	gobDao "go-common/app/interface/main/tv/dao/goblin"
	"go-common/app/interface/main/tv/model"
	"go-common/app/interface/main/tv/model/goblin"
	tvapi "go-common/app/service/main/tv/api"
)

// Service .
type Service struct {
	conf      *conf.Config
	dao       *gobDao.Dao
	accDao    *account.Dao
	arcDao    *archive.Dao
	ChlSplash map[string]string // channel's splash data
	Hotword   []*model.Hotword  // search hotword data
	VipQns    map[string]int    // playurl qualities for vips
	labels    *goblin.IndexLabels
	tvCilent  tvapi.TVServiceClient
}

var ctx = context.TODO()

// New .
func New(c *conf.Config) *Service {
	srv := &Service{
		conf:      c,
		dao:       gobDao.New(c),
		ChlSplash: make(map[string]string),
		VipQns:    make(map[string]int),
		accDao:    account.New(c),
		arcDao:    archive.New(c),
		labels:    &goblin.IndexLabels{},
	}
	var err error
	if srv.tvCilent, err = tvapi.NewClient(c.TvVipClient); err != nil {
		panic(err)
	}
	for _, v := range c.Cfg.VipQns {
		srv.VipQns[v] = 1
	}
	go srv.loadSph()         // splash
	go srv.loadHotword()     // hotword
	go srv.loadSphproc()     // splash proc
	go srv.loadHotwordproc() // hotword proc
	srv.prepareLabels()      // prepare index labels
	go srv.labelsproc()
	return srv
}
