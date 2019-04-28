package service

import (
	"context"
	"go-common/library/log"
	"time"

	"go-common/app/admin/main/tv/conf"
	"go-common/app/admin/main/tv/dao"
	"go-common/app/admin/main/tv/model"
	acccli "go-common/app/service/main/account/api"
	arccli "go-common/app/service/main/archive/api"
	httpx "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

var ctx = context.Background()

// Service biz service def.
type Service struct {
	c           *conf.Config
	dao         *dao.Dao
	DB, DBShow  *gorm.DB
	accClient   acccli.AccountClient
	arcClient   arccli.ArchiveClient
	SupCats     []*model.ParentCat
	supCatMap   *model.SupCats
	IntervLimit int
	arcPTids    map[int32][]int32 // archive parent type ids
	ArcTypes    map[int32]*arccli.Tp
	avaiTps     *model.AvailTps
	snsInfo     map[int64]*model.TVEpSeason
	snsCats     map[int][]int64
	abnCids     []*model.AbnorCids // abnormal cids
	pgcCatName  map[int]string     // pgc category name
	labelTps    map[int][]*model.TpLabel
	client      *httpx.Client
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		dao:         dao.New(c),
		IntervLimit: c.Cfg.IntervLimit,
		SupCats:     make([]*model.ParentCat, 0),
		ArcTypes:    make(map[int32]*arccli.Tp),
		arcPTids:    make(map[int32][]int32),
		avaiTps:     &model.AvailTps{},
		snsInfo:     make(map[int64]*model.TVEpSeason),
		snsCats:     make(map[int][]int64),
		pgcCatName:  make(map[int]string),
		labelTps:    make(map[int][]*model.TpLabel),
		client:      httpx.NewClient(conf.Conf.HTTPClient),
	}
	s.DB = s.dao.DB
	s.DBShow = s.dao.DBShow
	var err error
	if s.accClient, err = acccli.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	if s.arcClient, err = arccli.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	for k, v := range c.Cfg.PgcNames {
		s.pgcCatName[atoi(k)] = v
	}
	s.loadData()
	go s.loadDataproc()
	s.loadSns(context.Background()) // load season info
	go s.loadSnsproc()
	s.loadAbnCids() // load abnormal cids
	go s.loadAbnCidsproc()
	go s.refLabelproc() // refresh ugc + pgc labels
	go s.checkPanel()
	return s
}

func (s *Service) loadDataproc() {
	for {
		time.Sleep(time.Duration(s.c.Cfg.SupportCat.ReloadFre))
		s.loadData()
	}
}

func (s *Service) refLabelproc() {
	for {
		s.ugcLabels()
		s.pgcLabels()
		time.Sleep(time.Duration(s.c.Cfg.RefLabel.Fre))
	}
}

func (s *Service) checkPanel() {
	for {
		time.Sleep(time.Duration(3600) * time.Second)
		log.Info("check panel info start!")
		s.checkRemotePanel(ctx)
		log.Info("check panel info end!")
	}
}

func (s *Service) loadData() {
	s.loadTypes() // load ugc types
	s.loadTps()   // load passed tps and all tps for cms type list
	s.loadCats()  // load support categorys ( pgc & ugc)
	s.loadLabel() // load pgc label types
}

// Wait wait all closed.
func (s *Service) Wait() {
}

// Close close all dao.
func (s *Service) Close() {
}
