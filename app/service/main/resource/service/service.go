package service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	arcrpc "go-common/app/service/main/archive/api/gorpc"
	arcmdl "go-common/app/service/main/archive/model/archive"
	locrpc "go-common/app/service/main/location/rpc/client"
	pb "go-common/app/service/main/resource/api/v1"
	"go-common/app/service/main/resource/conf"
	"go-common/app/service/main/resource/dao/abtest"
	"go-common/app/service/main/resource/dao/ads"
	"go-common/app/service/main/resource/dao/alarm"
	"go-common/app/service/main/resource/dao/cpm"
	"go-common/app/service/main/resource/dao/manager"
	"go-common/app/service/main/resource/dao/resource"
	"go-common/app/service/main/resource/dao/show"
	"go-common/app/service/main/resource/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_updateAct    = "update"
	_archiveTable = "archive"
)

var (
	_emptyResources = make(map[int]*model.Resource)
)

// Service define resource service
type Service struct {
	c        *conf.Config
	cpm      *cpm.Dao
	abtest   *abtest.Dao
	res      *resource.Dao
	ads      *ads.Dao
	alarmDao *alarm.Dao
	show     *show.Dao
	manager  *manager.Dao
	// location rpc
	locationRPC *locrpc.Service
	// web cache
	resCache            []*model.Resource           // => resource
	asgCache            []*model.Assignment         // => assignments
	resCacheMap         map[int]*model.Resource     // resID => resource
	asgCacheMap         map[int][]*model.Assignment // resID => [ => assignments]
	defBannerCache      *model.Assignment
	videoAdsAPPCache    map[int8]map[int8]map[int8]map[string]*model.VideoAD // plat => [ => adsType] => [ => adsTarget] => aid(seasonId or typeId)
	missch              chan interface{}
	typeList            map[string]string
	resArchiveWarnCache map[int64][]*model.ResWarnInfo
	resURLWarnCache     map[string][]*model.ResWarnInfo
	posCache            map[int][]int
	// app cache
	bannerCache         map[int8]map[int][]*model.Banner
	categoryBannerCache map[int8]map[int][]*model.Banner
	bannerHashCache     map[int8]string
	bannerLimitCache    map[int]int
	indexIcon           map[int][]*model.IndexIcon
	playIcon            *model.PlayerIcon
	cardCache           map[int8]*model.Head
	sideBarCache        []*model.SideBar
	sideBarLimitCache   map[int64][]*model.SideBarLimit
	// live
	cmtbox map[int64]*model.Cmtbox
	// abtest
	abTestCache map[string]*model.AbTest
	// PasterAIDCache
	PasterAIDCache []map[int64]int64
	// rpc
	arcRPC *arcrpc.Service2
	// database
	archiveSub      *databus.Databus
	closeSub        bool
	closeMonitorURL bool
	// waiter
	waiter sync.WaitGroup
	// lock
	abTestLock sync.Mutex
	//pgc special cards
	specialCache map[int64]*pb.SpecialReply
	//pgc relate cards
	relateCache map[int64]*model.Relate
	//pgc id relate relate card id
	relatePgcMapCache map[int64]int64
	// audit
	auditCache map[string][]int
}

// New return service object
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                   c,
		cpm:                 cpm.New(c),
		abtest:              abtest.New(c),
		res:                 resource.New(c),
		ads:                 ads.New(c),
		alarmDao:            alarm.New(c),
		show:                show.New(c),
		manager:             manager.New(c),
		locationRPC:         locrpc.New(c.LocationRPC),
		resCacheMap:         make(map[int]*model.Resource),
		asgCacheMap:         make(map[int][]*model.Assignment),
		videoAdsAPPCache:    make(map[int8]map[int8]map[int8]map[string]*model.VideoAD),
		missch:              make(chan interface{}, 10240),
		typeList:            make(map[string]string),
		resArchiveWarnCache: make(map[int64][]*model.ResWarnInfo),
		resURLWarnCache:     make(map[string][]*model.ResWarnInfo),
		bannerHashCache:     make(map[int8]string),
		indexIcon:           make(map[int][]*model.IndexIcon),
		bannerCache:         make(map[int8]map[int][]*model.Banner),
		categoryBannerCache: make(map[int8]map[int][]*model.Banner),
		posCache:            make(map[int][]int),
		cardCache:           make(map[int8]*model.Head),
		cmtbox:              make(map[int64]*model.Cmtbox),
		sideBarLimitCache:   make(map[int64][]*model.SideBarLimit),
		abTestCache:         make(map[string]*model.AbTest),
		arcRPC:              arcrpc.New2(c.ArchiveRPC),
		archiveSub:          databus.New(c.ArchiveSub),
		specialCache:        make(map[int64]*pb.SpecialReply),
		relateCache:         make(map[int64]*model.Relate),
		relatePgcMapCache:   make(map[int64]int64),
		auditCache:          make(map[string][]int),
	}
	if err := s.loadRes(); err != nil {
		panic(err)
	}
	if err := s.loadVideoAds(); err != nil {
		panic(err)
	}
	if err := s.loadBannerCahce(); err != nil {
		panic(err)
	}
	s.loadTypeList()
	s.loadPlayIcon()
	s.loadCmtbox()
	s.loadSpecialCache()
	s.loadRelateCache()
	s.loadAudit()
	go s.loadproc()
	go s.loadCmtboxproc()
	go s.loadCardCache()
	go s.loadSideBarCache()
	go s.cacheproc()
	if s.c.MonitorArchive {
		s.waiter.Add(1)
		go s.arcConsume()
	}
	if s.c.MonitorURL {
		go s.checkResURL()
	}
	return
}

// loadproc is a routine load ads to cache
func (s *Service) loadproc() {
	for {
		time.Sleep(time.Duration(conf.Conf.Reload.Ad))
		s.loadRes()
		s.loadVideoAds()
		s.loadBannerCahce()
		s.loadTypeList()
		s.loadPlayIcon()
		s.loadCardCache()
		s.loadSideBarCache()
		s.loadRelateCache()
		s.loadSpecialCache()
		s.loadAudit()
	}
}

// loadCmtboxproc is a routine load cmtbox to cache
func (s *Service) loadCmtboxproc() {
	for {
		time.Sleep(time.Second * 10)
		s.loadCmtbox()
	}
}

func (s *Service) loadTypeList() (err error) {
	var (
		tmpTypeList  map[int16]*arcmdl.ArcType
		tmpTypeList2 = make(map[string]string)
	)
	if tmpTypeList, err = s.arcRPC.Types2(context.TODO()); err != nil || len(tmpTypeList) == 0 {
		log.Error("s.arcRPC.Types2() error(%v) or typelist len is zero", err)
		return
	}
	for tid, typeInfo := range tmpTypeList {
		tidStr := strconv.Itoa(int(tid))
		pidStr := strconv.Itoa(int(typeInfo.Pid))
		tmpTypeList2[tidStr] = pidStr
	}
	s.typeList = tmpTypeList2
	return
}

func (s *Service) loadPlayIcon() {
	var (
		pi  *model.PlayerIcon
		err error
	)
	if pi, err = s.res.PlayerIcon(context.TODO()); err != nil {
		log.Error("s.res.PlayerIcon() error(%v)", err)
		return
	}
	s.playIcon = pi
}

func (s *Service) loadCmtbox() {
	var (
		cmtbox map[int64]*model.Cmtbox
		err    error
	)
	if cmtbox, err = s.res.Cmtbox(context.TODO()); err != nil {
		log.Error("s.res.Cmtbox() error(%v)", err)
		return
	}
	s.cmtbox = cmtbox
}

func (s *Service) loadAudit() {
	var (
		at  map[string][]int
		err error
	)
	if at, err = s.show.Audit(context.TODO()); err != nil {
		log.Error("s.show.Audit error(%v)", err)
		return
	}
	s.auditCache = at
}

func (s *Service) checkResURL() {
	for {
		time.Sleep(time.Duration(conf.Conf.Reload.Ad))
		if s.closeMonitorURL {
			return
		}
		for url, resURL := range s.resURLWarnCache {
			s.alarmDao.CheckURL(url, resURL)
			time.Sleep(time.Duration(conf.Conf.SpLimit))
		}
	}
}

// arcConsume consumer archive
func (s *Service) arcConsume() {
	defer s.waiter.Done()
	var (
		msgs = s.archiveSub.Messages()
		err  error
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.archiveSub.Message closed", err)
			return
		}
		if s.closeSub {
			return
		}
		msg.Commit()
		m := &model.Message{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if m.Table == _archiveTable {
			s.arcChan(m.Action, m.New, m.Old)
		}
	}
}

// Ping ping service
func (s *Service) Ping(c context.Context) (err error) {
	return s.res.Ping(c)
}

// Close close service
func (s *Service) Close() {
	if s.c.MonitorArchive {
		s.closeSub = true
		time.Sleep(2 * time.Second)
		s.archiveSub.Close()
		s.waiter.Wait()
	}
	s.res.Close()
}

// Monitor for monitorURL
func (s *Service) Monitor(c context.Context) {
	s.closeMonitorURL = true
}

func (s *Service) addCache(d interface{}) {
	// asynchronous add rules to redis
	select {
	case s.missch <- d:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for add rules into redis.
func (s *Service) cacheproc() {
	for {
		d := <-s.missch
		switch d.(type) {
		case map[string]map[int64]int64:
			v := d.(map[string]map[int64]int64)
			if err := s.ads.AddBuvidCount(context.TODO(), v); err != nil {
				log.Error("s.ads.AddBuvidCount(%v) error(%+v)", v, err)
			}
		default:
			log.Warn("cacheproc can't process the type")
		}
	}
}
