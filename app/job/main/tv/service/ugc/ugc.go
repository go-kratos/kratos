package ugc

import (
	"context"
	"strconv"
	"sync"

	"go-common/app/job/main/tv/conf"
	appDao "go-common/app/job/main/tv/dao/app"
	arcdao "go-common/app/job/main/tv/dao/archive"
	"go-common/app/job/main/tv/dao/cms"
	"go-common/app/job/main/tv/dao/ftp"
	"go-common/app/job/main/tv/dao/lic"
	playdao "go-common/app/job/main/tv/dao/playurl"
	ugcdao "go-common/app/job/main/tv/dao/ugc"
	updao "go-common/app/job/main/tv/dao/upper"
	arccli "go-common/app/service/main/archive/api"
	archive "go-common/app/service/main/archive/api/gorpc"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/robfig/cron"
)

const _chanSize = 10240

var ctx = context.TODO()

// Service is show service.
type Service struct {
	c *conf.Config
	// dao
	dao        *ugcdao.Dao
	playurlDao *playdao.Dao
	licDao     *lic.Dao
	ftpDao     *ftp.Dao
	appDao     *appDao.Dao
	arcDao     *arcdao.Dao
	upDao      *updao.Dao
	cmsDao     *cms.Dao
	// logic
	daoClosed bool
	// waiter
	waiter        *sync.WaitGroup
	consumerLimit chan struct{}
	// rpc
	arcClient arccli.ArchiveClient
	arcRPC    *archive.Service2
	// databus
	archiveNotifySub *databus.Databus
	ugcSub           *databus.Databus
	// memory data
	ugcTypesRel map[int32]*conf.UgcType
	ugcTypesCat map[int32]int32
	arcTypes    map[int32]*arccli.Tp // map for arc types
	pgcTypes    map[string]int       // filter pgc types data
	activeUps   map[int64]int        // store all the trusted uppers
	// cron
	cron *cron.Cron
	// channels
	modArcCh, audAidCh     chan []int64
	repCidCh, reshelfAidCh chan int64
}

// New inits the ugc service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                c,
		dao:              ugcdao.New(c),
		playurlDao:       playdao.New(c),
		licDao:           lic.New(c),
		ftpDao:           ftp.New(c),
		arcDao:           arcdao.New(c),
		appDao:           appDao.New(c),
		upDao:            updao.New(c),
		cmsDao:           cms.New(c),
		waiter:           new(sync.WaitGroup),
		archiveNotifySub: databus.New(c.ArchiveNotifySub),
		ugcSub:           databus.New(c.UgcSub),
		cron:             cron.New(),
		arcTypes:         make(map[int32]*arccli.Tp),
		ugcTypesRel:      make(map[int32]*conf.UgcType),
		pgcTypes:         make(map[string]int),
		activeUps:        make(map[int64]int),
		consumerLimit:    make(chan struct{}, c.UgcSync.Cfg.ThreadLimit),
		modArcCh:         make(chan []int64, _chanSize),
		audAidCh:         make(chan []int64, _chanSize),
		reshelfAidCh:     make(chan int64, _chanSize),
		repCidCh:         make(chan int64, _chanSize),
		arcRPC:           archive.New2(c.ArchiveRPC),
	}
	// transform cfg to map, in order to filter pgc types archive
	var err error
	if s.arcClient, err = arccli.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	for _, v := range s.c.Cfg.PgcTypes {
		s.pgcTypes[v] = 1
	}
	for k, v := range s.c.Cfg.UgcZones { // transform cfg map
		s.ugcTypesRel[atoi(k)] = v
	}
	if err := s.cron.AddFunc(s.c.Redis.CronUGC, s.ZoneIdx); err != nil { // load Zone Idx & types
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.UgcSync.Frequency.TypesCron, s.loadTypes); err != nil {
		panic(err)
	}
	s.cron.Start()
	s.loadTypes()           // load types
	s.loadTids()            // load ugc idx relationship
	s.refreshUp(ctx, false) // init upper list
	go s.ZoneIdx()
	s.waiter.Add(1)
	go s.syncUpproc()    // sync modified uppers' info to license owner
	go s.refreshUpproc() // refresh upper info
	s.waiter.Add(1)
	go s.manualproc() // manual import videos
	s.waiter.Add(1)
	go s.upImportproc()    // import upper history data
	go s.fullRefreshproc() // full refrsh video data
	s.waiter.Add(1)
	go s.modArcproc() // sync modified archive data
	s.waiter.Add(1)
	go s.delArcproc() // sync deleted archive data
	s.waiter.Add(1)
	go s.delVideoproc() // sync deleted video data
	s.waiter.Add(1)
	go s.delUpproc()      // treat deleted uppers
	go s.seaUgcContproc() // uploads ugc search content to sug's FTP
	s.waiter.Add(1)
	go s.arcConsumeproc() // archive Notify-T databus
	s.waiter.Add(1)
	go s.consumeVideo() // consume video databus, report cid info
	s.waiter.Add(1)
	go s.repCidproc() // consume channel and report cid to playurl
	s.waiter.Add(1)
	go s.audCidproc() // consume audit aid data
	s.waiter.Add(1)
	go s.reshelfArcproc() // reshelf the cms invalid arcs when they have at least one video that can play now
	return
}

// Close close the services
func (s *Service) Close() {
	if s.dao != nil {
		s.daoClosed = true
		log.Info("Dao Closed!")
	}
	log.Info("Close ArcNotifySub!")
	s.archiveNotifySub.Close()
	log.Info("Close ugcSub!")
	s.ugcSub.Close()
	log.Info("Wait Sync!")
	s.waiter.Wait()
	log.Info("DB Closed Physically!")
	s.dao.DB.Close()
}

// transform string to int
func atoi(number string) int32 {
	res, _ := strconv.Atoi(number)
	return int32(res)
}
