package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go-common/app/interface/main/mcn/tool/worker"
	accrpc "go-common/app/service/main/account/api"
	"go-common/app/service/main/videoup/conf"
	"go-common/app/service/main/videoup/dao/agent"
	"go-common/app/service/main/videoup/dao/archive"
	"go-common/app/service/main/videoup/dao/bgm"
	busdao "go-common/app/service/main/videoup/dao/databus"
	"go-common/app/service/main/videoup/dao/dede"
	"go-common/app/service/main/videoup/dao/manager"
	"go-common/app/service/main/videoup/dao/monitor"
	"go-common/app/service/main/videoup/dao/msg"
	"go-common/app/service/main/videoup/dao/relation"
	"go-common/app/service/main/videoup/dao/ups"
	arcmdl "go-common/app/service/main/videoup/model/archive"
	ddmdl "go-common/app/service/main/videoup/model/dede"
	msgmdl "go-common/app/service/main/videoup/model/message"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

// Service is service.
type Service struct {
	c        *conf.Config
	arc      *archive.Dao
	dede     *dede.Dao
	mng      *manager.Dao
	busCache *busdao.Dao
	monitor  *monitor.Dao
	agent    *agent.Dao
	bgm      *bgm.Dao
	relation *relation.Dao
	ups      *ups.Dao
	// databus
	videoupPub    *databus.Databus
	videoupPGCPub *databus.Databus
	// cache: type、upper
	typeCache       map[int16]*arcmdl.Type
	upperCache      map[int8]map[int64]struct{}
	monitorMap      map[string]*arcmdl.Alert
	locker          *sync.Mutex
	flowsCache      []*arcmdl.Flow
	forbidMidsCache map[int64]string
	specialUpsCache []*arcmdl.Up
	whiteMidsCache  map[int64]int64
	grayCheckUps    map[int64]int64
	// error chan
	padCh   chan *ddmdl.PadInfo
	msgCh   chan *msgmdl.Videoup
	asyncCh chan func()
	// wait
	wg      sync.WaitGroup
	closed  bool
	stop    chan struct{}
	promSub *prom.Prom
	promP   *prom.Prom
	promErr *prom.Prom

	veditor *arcmdl.VideosEditor
	worker  *worker.Pool
	msg     *msg.Dao
	msgMap  map[int]*arcmdl.MSG

	accRPC accrpc.AccountClient
}

// New is go-common/app/service/videoup service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		arc:          archive.New(c),
		dede:         dede.New(c),
		mng:          manager.New(c),
		busCache:     busdao.New(c),
		monitor:      monitor.New(c),
		agent:        agent.New(c),
		bgm:          bgm.New(c),
		ups:          ups.New(c),
		relation:     relation.New(c),
		monitorMap:   make(map[string]*arcmdl.Alert),
		locker:       &sync.Mutex{},
		grayCheckUps: make(map[int64]int64),
		// pub
		videoupPub:    databus.New(c.VideoupPub),
		videoupPGCPub: databus.New(c.VideoupPGCPub),
		// chan
		padCh:   make(chan *ddmdl.PadInfo, c.ChanSize),
		msgCh:   make(chan *msgmdl.Videoup, c.ChanSize),
		asyncCh: make(chan func(), c.ChanSize),
		stop:    make(chan struct{}),
		promSub: prom.BusinessInfoCount,
		promP:   prom.BusinessInfoCount,
		promErr: prom.BusinessErrCount,
		veditor: arcmdl.NewEditor(c.FailThreshold),
		worker:  worker.New(nil),
		msg:     msg.New(c),
	}
	var err error
	if s.accRPC, err = accrpc.NewClient(c.AccRPC); err != nil {
		panic(err)
	}
	// load cache.
	s.loadType()
	s.loadUpper()
	s.loadFlows()
	s.loadUpSpecial()
	s.loadWhiteMids()
	s.loadForbidMids()
	s.setMsgTypeMap()
	s.loadGrayCheckUps()
	go s.cacheproc()
	go s.padproc()
	go s.msgproc()
	go s.asyncproc()
	go s.monitorConsume()
	return s
}

// loadGrayCheckUps .
func (s *Service) loadGrayCheckUps() {
	if s.c.GrayGroup == 0 {
		return
	}
	checkUps, err := s.ups.GrayCheckUps(context.TODO(), s.c.GrayGroup)
	if err != nil {
		return
	}
	s.grayCheckUps = checkUps
}

func (s *Service) checkGrayMid(mid int64) (ok bool) {
	if s.c.GrayGroup == 0 {
		ok = true
	} else {
		_, ok = s.grayCheckUps[mid]
	}
	log.Info("mid(%d) checkGrayMid(%v) and gray_group is (%d)", mid, ok, s.c.GrayGroup)
	return
}

// AllowType is typeid in typeinfo
func (s *Service) AllowType(c context.Context, typeID int16) (ok bool) {
	tp := s.typeCache[typeID]
	ok = (tp != nil) && (tp.PID > 0)
	return
}

// Types get all types
func (s *Service) Types(c context.Context) (types map[int16]*arcmdl.Type) {
	types = s.typeCache
	return
}

func (s *Service) loadType() {
	tpm, err := s.arc.TypeMapping(context.TODO())
	if err != nil {
		log.Error("s.arc.TypeMapping error(%v)", err)
		return
	}
	s.typeCache = tpm
}

func (s *Service) loadUpper() {
	upm, err := s.mng.Uppers(context.TODO())
	if err != nil {
		log.Error("s.mng.Uppers error(%v)", err)
		return
	}
	s.upperCache = upm
}

func (s *Service) loadFlows() {
	flows, err := s.arc.Flows(context.TODO())
	if err != nil {
		log.Error("s.arc.Flows error(%v)", err)
		return
	}
	s.flowsCache = flows
}

func (s *Service) loadUpSpecial() {
	sups, err := s.mng.UpSpecial(context.TODO())
	if err != nil {
		log.Error("s.arc.loadUpSpecial error(%v)", err)
		return
	}
	s.specialUpsCache = sups
}

//ArcTag .
func (s *Service) ArcTag(c context.Context, aid int64) (tagID int64, err error) {
	tagID, err = s.mng.ArcReason(c, aid)
	return
}

func (s *Service) loadWhiteMids() {
	whiteMids, err := s.arc.WhiteMids(context.TODO())
	if err != nil {
		log.Error("s.arc.WhiteMids error(%v)", err)
		return
	}
	s.whiteMidsCache = whiteMids
}

func (s *Service) loadForbidMids() {
	forbidMidList, err := s.arc.ForbidMids(context.TODO())
	if err != nil {
		log.Error("s.arc.ForbidMids error(%v)", err)
		return
	}
	log.Info("forbidMidList (%+v)", forbidMidList)
	forbidMids := make(map[int64]string)
	for mid, up := range forbidMidList {
		sumForbid := &arcmdl.ForbidAttr{RankV: int32(0), RecommendV: int32(0), ShowV: int32(0), SearchV: int32(0), PushBlogV: int32(0)}
		for _, forbidJSON := range up {
			lineForbid := &arcmdl.ForbidAttr{}
			if err = json.Unmarshal([]byte(forbidJSON), lineForbid); err != nil {
				log.Error("forbidMidList(%+v)  mid(%s) line json.Unmarshal(%+v) error(%v)", forbidMidList, mid, forbidJSON, err)
				continue
			}
			//config forbid
			lineForbid.Convert()
			sumForbid.Convert()
			ok := int32(1)
			//rank
			if lineForbid.Rank.Main == ok {
				sumForbid.Rank.Main = 1
			}
			if lineForbid.Rank.RecentArc == ok {
				sumForbid.Rank.RecentArc = 1
			}
			if lineForbid.Rank.AllArc == ok {
				sumForbid.Rank.AllArc = 1
			}
			//Recommend
			if lineForbid.Recommend.Main == ok {
				sumForbid.Recommend.Main = 1
			}
			//Search
			if lineForbid.SearchV == ok {
				sumForbid.SearchV = 1
			}
			//PushBlog
			if lineForbid.PushBlogV == ok {
				sumForbid.PushBlogV = 1
			}
			//Dynamic
			if lineForbid.Dynamic.Main == ok {
				sumForbid.Dynamic.Main = 1
			}
			//show
			if lineForbid.Show.Main == ok {
				sumForbid.Show.Main = 1
			}
			if lineForbid.Show.Mobile == ok {
				sumForbid.Show.Mobile = 1
			}
			if lineForbid.Show.Web == ok {
				sumForbid.Show.Web = 1
			}
			if lineForbid.Show.Oversea == ok {
				sumForbid.Show.Oversea = 1
			}
			if lineForbid.Show.Online == ok {
				sumForbid.Show.Online = 1
			}
			sumForbid.Reverse()
		}
		forbidStr, err := json.Marshal(sumForbid)
		if err != nil {
			log.Error("s.loadForbidMids Marshal(%+v) error (%s)", sumForbid, err)
			continue
		}
		forbidMids[mid] = string(forbidStr)
	}
	s.forbidMidsCache = forbidMids
}

func (s *Service) cacheproc() {
	for {
		time.Sleep(5 * time.Minute)
		s.loadType()
		s.loadUpper()
		s.loadFlows()
		s.loadUpSpecial()
		s.loadWhiteMids()
		s.loadForbidMids()
		s.loadGrayCheckUps()
	}
}

func (s *Service) asyncproc() {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		f, ok := <-s.asyncCh
		if !ok {
			return
		}
		f()
	}
}

// Close  consumer close.
func (s *Service) Close() {
	s.worker.Close()
	s.worker.Wait()
	s.arc.Close()
	s.mng.Close()
	s.busCache.Close()
	close(s.asyncCh)
	time.Sleep(1 * time.Second)
	close(s.stop)
	close(s.padCh)
	close(s.msgCh)
	s.closed = true
	s.veditor.Close()
	s.wg.Wait()

}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.arc.Ping(c); err != nil {
		return
	}
	if err = s.mng.Ping(c); err != nil {
		return
	}
	return s.busCache.Ping(c)
}
