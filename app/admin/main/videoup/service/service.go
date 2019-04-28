package service

import (
	"context"
	"go-common/app/admin/main/videoup/dao/monitor"
	"math"
	"sync"
	"time"

	"go-common/app/admin/main/videoup/conf"
	arcdao "go-common/app/admin/main/videoup/dao/archive"
	datadao "go-common/app/admin/main/videoup/dao/data"
	busdao "go-common/app/admin/main/videoup/dao/databus"
	mngdao "go-common/app/admin/main/videoup/dao/manager"
	musicdao "go-common/app/admin/main/videoup/dao/music"
	overseadao "go-common/app/admin/main/videoup/dao/oversea"
	searchdao "go-common/app/admin/main/videoup/dao/search"
	staffdao "go-common/app/admin/main/videoup/dao/staff"
	tagDao "go-common/app/admin/main/videoup/dao/tag"
	taskdao "go-common/app/admin/main/videoup/dao/task"
	trackdao "go-common/app/admin/main/videoup/dao/track"
	arcmdl "go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/manager"
	mngmdl "go-common/app/admin/main/videoup/model/manager"
	msgmdl "go-common/app/admin/main/videoup/model/message"
	accApi "go-common/app/service/main/account/api"
	upsrpc "go-common/app/service/main/up/api/v1"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/jinzhu/gorm"
	"go-common/library/net/http/blademaster/middleware/permit"
)

// Service is service.
type Service struct {
	c         *conf.Config
	arc       *arcdao.Dao
	busCache  *busdao.Dao
	mng       *mngdao.Dao
	oversea   *overseadao.Dao
	track     *trackdao.Dao
	music     *musicdao.Dao
	tag       *tagDao.Dao
	DB        *gorm.DB
	search    *searchdao.Dao
	staff     *staffdao.Dao
	overseaDB *gorm.DB
	data      *datadao.Dao
	monitor   *monitor.Dao
	task      *taskdao.Dao
	// acc rpc
	accRPC accApi.AccountClient
	upsRPC upsrpc.UpClient
	// databus
	videoupPub  *databus.Databus
	upCreditPub *databus.Databus
	// cache:  upper
	adtTpsCache       map[int16]struct{}
	thrTpsCache       map[int16]int
	thrMin, thrMax    int
	upperCache        map[int8]map[int64]struct{} //TODO 这个缓存需要从up服务里取
	allUpGroupCache   map[int64]*manager.UpGroup  //UP主分组列表
	fansCache         int64
	roundTpsCache     map[int16]struct{}
	typeCache         map[int16]*arcmdl.Type
	typeCache2        map[int16][]int64 // 记录每个一级分区下的二级分区
	flowsCache        map[int64]string
	porderConfigCache map[int64]*arcmdl.PorderConfig
	twConCache        map[int8]map[int64]*arcmdl.WCItem
	// error chan
	msgCh chan *msgmdl.Videoup
	// wait
	wg     sync.WaitGroup
	closed bool
	stop   chan struct{}
	auth   *permit.Permit
}

// New is videoup-admin service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		arc:      arcdao.New(c),
		mng:      mngdao.New(c),
		oversea:  overseadao.New(c),
		music:    musicdao.New(c),
		track:    trackdao.New(c),
		busCache: busdao.New(c),
		search:   searchdao.New(c),
		staff:    staffdao.New(c),
		tag:      tagDao.New(c),
		data:     datadao.New(c),
		monitor:  monitor.New(c),
		task:     taskdao.New(c),
		// pub
		videoupPub:  databus.New(c.VideoupPub),
		upCreditPub: databus.New(c.UpCreditPub),
		// chan
		msgCh: make(chan *msgmdl.Videoup, c.ChanSize),
		// map
		stop: make(chan struct{}),
		auth: permit.New(c.Auth),
	}
	var err error
	if s.accRPC, err = accApi.NewClient(c.AccountRPC); err != nil {
		panic(err)
	}
	if s.upsRPC, err = upsrpc.NewClient(c.UpsRPC); err != nil {
		panic(err)
	}
	s.DB = s.music.DB
	s.overseaDB = s.oversea.OverseaDB
	// load cache.
	s.loadType()
	s.loadUpGroups()
	s.loadUpper()
	s.loadConf()
	go s.cacheproc()
	go s.msgproc()
	s.wg.Add(1) // NOTE: sync  add wait group
	go s.multSyncProc()
	return s
}

func (s *Service) isWhite(mid int64) bool {
	if ups, ok := s.upperCache[mngmdl.UpperTypeWhite]; ok {
		_, is := ups[mid]
		return is
	}
	return false
}

//PGCWhite whether mid in pgcwhite list
func (s *Service) PGCWhite(mid int64) bool {
	if ups, ok := s.upperCache[mngmdl.UpperTypePGCWhite]; ok {
		_, is := ups[mid]
		return is
	}
	return false
}

func (s *Service) isBlack(mid int64) bool {
	if ups, ok := s.upperCache[mngmdl.UpperTypeBlack]; ok {
		_, is := ups[mid]
		return is
	}
	return false
}

func (s *Service) getAllUPGroups(mid int64) (gs []int64) {
	gs = []int64{}
	for tp, item := range s.upperCache {
		if _, exist := item[mid]; !exist {
			continue
		}
		gs = append(gs, int64(tp))
	}
	return
}

func (s *Service) isAuditType(tpID int16) bool {
	_, isAt := s.adtTpsCache[tpID]
	return isAt
}

func (s *Service) isRoundType(tpID int16) bool {
	_, in := s.roundTpsCache[tpID]
	return in
}

func (s *Service) isTypeID(tpID int16) bool {
	_, in := s.typeCache[tpID]
	return in
}

func (s *Service) loadType() {
	// TODO : audit types
	// threshold
	thr, err := s.arc.ThresholdConf(context.TODO())
	if err != nil {
		log.Error("s.arc.ThresholdConf error(%v)", err)
		return
	}
	s.thrTpsCache = thr
	var min, max = math.MaxInt32, 0
	for _, t := range thr {
		if min > t {
			min = t
		}
		if max < t {
			max = t
		}
	}
	s.thrMin = min
	s.thrMax = max
}

func (s *Service) loadUpper() {
	upm, err := s.upSpecial(context.Background())
	if err != nil {
		log.Error("s.upSpecial error(%v)", err)
		return
	}
	s.upperCache = upm
}

// loadUpGroups 加载所有UP分组列表
func (s *Service) loadUpGroups() {
	groups, err := s.mng.UpGroups(context.TODO())
	if err != nil {
		log.Error("s.mng.UpGroups() error(%v)", err)
		return
	}
	s.allUpGroupCache = groups
}

func (s *Service) loadConf() {
	var (
		fans          int64
		err           error
		auditTypes    map[int16]struct{}
		roundTypes    map[int16]struct{}
		flows         map[int64]string
		tpm           map[int16]*arcmdl.Type
		porderConfigs map[int64]*arcmdl.PorderConfig
		twConCache    map[int8]map[int64]*arcmdl.WCItem
		tpm2          map[int16][]int64
	)
	if fans, err = s.arc.FansConf(context.TODO()); err != nil {
		log.Error("s.arc.FansConf error(%v)", err)
		return
	}
	s.fansCache = fans
	if auditTypes, err = s.arc.AuditTypesConf(context.TODO()); err != nil {
		log.Error("s.arc.AuditTypesConf error(%v)", err)
		return
	}
	s.adtTpsCache = auditTypes
	if roundTypes, err = s.arc.RoundTypeConf(context.TODO()); err != nil {
		log.Error("s.arc.RoundTypeConf error(%v)", err)
		return
	}
	s.roundTpsCache = roundTypes
	if flows, err = s.arc.Flows(context.TODO()); err != nil {
		log.Error("s.arc.Flows error(%v)", err)
		return
	}
	s.flowsCache = flows
	if tpm, err = s.arc.TypeMapping(context.TODO()); err != nil {
		log.Error("s.arc.TypeMapping error(%v)", err)
		return
	}
	s.typeCache = tpm

	tpm2 = make(map[int16][]int64)
	for id, tmod := range tpm {
		if tmod.PID == 0 {
			if _, ok := tpm2[id]; !ok {
				tpm2[id] = []int64{}
			}
			continue
		}
		arrid, ok := tpm2[tmod.PID]
		if !ok {
			tpm2[tmod.PID] = []int64{int64(id)}
		} else {
			tpm2[tmod.PID] = append(arrid, int64(id))
		}
	}
	s.typeCache2 = tpm2

	if porderConfigs, err = s.arc.PorderConfig(context.TODO()); err != nil {
		log.Error("s.arc.PorderConfig error(%v)", err)
		return
	}
	s.porderConfigCache = porderConfigs

	if twConCache, err = s.weightConf(context.TODO()); err != nil {
		log.Error("s.weightConf error(%v)", err)
		return
	}
	s.twConCache = twConCache
}

func (s *Service) cacheproc() {
	for {
		time.Sleep(3 * time.Minute)
		s.loadType()
		s.loadUpper()
		s.loadUpGroups()
		s.loadConf()
		s.lockVideo()
		go s.MonitorNotifyResult(context.TODO())
	}
}

// Close  consumer close.
func (s *Service) Close() {
	s.arc.Close()
	s.mng.Close()
	s.music.Close()
	s.busCache.Close()
	time.Sleep(1 * time.Second)
	close(s.stop)
	close(s.msgCh)
	s.closed = true
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
