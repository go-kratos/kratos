package service

import (
	"context"
	"math"
	"sync"
	"time"

	"go-common/app/job/main/videoup/conf"
	"go-common/app/job/main/videoup/dao/activity"
	"go-common/app/job/main/videoup/dao/archive"
	"go-common/app/job/main/videoup/dao/bvc"
	"go-common/app/job/main/videoup/dao/manager"
	"go-common/app/job/main/videoup/dao/message"
	"go-common/app/job/main/videoup/dao/monitor"
	"go-common/app/job/main/videoup/dao/redis"
	mngmdl "go-common/app/job/main/videoup/model/manager"
	accApi "go-common/app/service/main/account/api"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"

	"github.com/pkg/errors"
)

// Service is service.
type Service struct {
	c *conf.Config
	// wait group
	wg sync.WaitGroup
	// acc rpc
	accRPC accApi.AccountClient
	// dao
	arc      *archive.Dao
	mng      *manager.Dao
	msg      *message.Dao
	redis    *redis.Dao
	monitor  *monitor.Dao
	activity *activity.Dao
	bvc      *bvc.Dao
	// databus sub
	bvc2VuSub     *databus.Databus
	videoupSub    *databus.Databus
	arcResultSub  *databus.Databus
	videoshotSub2 *databus.Databus
	// videoupSub 幂等判断
	videoupSubIdempotent map[int32]int64
	statSub              *databus.Databus
	// databus pub
	videoupPub *databus.Databus
	blogPub    *databus.Databus
	// cache: type, upper
	sfTpsCache        map[int16]int16
	TypeMap           map[int16]string
	adtTpsCache       map[int16]struct{}
	thrTpsCache       map[int16]int
	thrMin, thrMax    int
	upperCache        map[int8]map[int64]struct{}
	fansCache         int
	roundTpsCache     map[int16]struct{}
	roundDelayCache   int64
	delayRoundMinTime time.Time
	specialUp         map[int64]struct{}
	// monitor
	bvc2VuMo      int64
	bvc2VuDelayMo int64
	videoupMo     int64
	arcResultMo   int64
	statMo        int64

	//prom moni
	promDatabus *prom.Prom
	promRetry   *prom.Prom
	//统计差值
	promVideoS *prom.Prom
	promVideoE *prom.Prom
	promPanic  *prom.Prom
	// closed
	closed bool
}

// New is videoup service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		// dao
		arc:      archive.New(c),
		mng:      manager.New(c),
		msg:      message.New(c),
		redis:    redis.New(c),
		monitor:  monitor.New(c),
		activity: activity.New(c),
		bvc:      bvc.New(c),

		bvc2VuSub:            databus.New(c.Bvc2VuSub),
		videoupSubIdempotent: make(map[int32]int64),
		videoupSub:           databus.New(c.VideoupSub),
		arcResultSub:         databus.New(c.ArcResultSub),
		statSub:              databus.New(c.StatSub),
		// databus pub
		videoupPub: databus.New(c.VideoupPub),
		blogPub:    databus.New(c.BlogPub),

		promDatabus: prom.BusinessInfoCount,
		promVideoS:  prom.CacheHit,
		promVideoE:  prom.CacheMiss,
		promPanic:   prom.CacheMiss,
		promRetry:   prom.BusinessErrCount,
	}
	var err error
	if s.accRPC, err = accApi.NewClient(c.AccRPC); err != nil {
		panic(err)
	}
	s.specialUp = make(map[int64]struct{}, len(c.SpecialUp))
	for _, mid := range c.SpecialUp {
		s.specialUp[mid] = struct{}{}
	}
	// load cache
	s.loadType()
	s.loadUpper()
	s.loadConf()
	s.wg.Add(1)
	go s.bvc2VuConsumer()
	s.wg.Add(1)
	go s.videoupConsumer()
	s.wg.Add(1)
	go s.statConsumer()
	s.wg.Add(1)
	go s.arcResultConsumer()
	if env.DeployEnv == env.DeployEnvProd {
		s.videoshotSub2 = databus.New(c.VideoshotSub2)
		s.wg.Add(1)
		go s.videoshotSHConsumer()
	}
	s.wg.Add(1)
	go s.retryproc()
	s.wg.Add(1)
	go s.QueueProc()
	s.wg.Add(1)
	go s.delayproc()
	s.wg.Add(1)
	go s.roundproc()
	go s.cacheproc()
	go s.monitorConsume()
	go s.edithistoryproc()
	return s
}

// Ping ping service.
func (s *Service) Ping(c context.Context) (err error) {
	return s.arc.Ping(c)
}

//Rescue runtime panic rescue
func (s *Service) Rescue(data interface{}) {
	r := recover()
	if r != nil {
		r = errors.WithStack(r.(error))
		log.Error("Runtime error caught: %+v and data is %+v", r, data)
		s.promPanic.Incr("panic")
	}
}

func (s *Service) edithistoryproc() {
	for {
		time.Sleep(nextDay(5))
		for {
			rows, _ := s.delArcEditHistory(100)
			time.Sleep(1 * time.Second)
			if rows == 0 {
				break
			}
		}
		for {
			rows, _ := s.delArcVideoEditHistory(100)
			time.Sleep(1 * time.Second)
			if rows == 0 {
				break
			}
		}
	}
}

// Until next day x hours
func nextDay(hour int) time.Duration {
	n := time.Now().Add(24 * time.Hour)
	d := time.Date(n.Year(), n.Month(), n.Day(), hour, 0, 0, 0, n.Location())
	return time.Until(d)
}

// Close  consumer close.
func (s *Service) Close() {
	s.bvc2VuSub.Close()
	s.videoupSub.Close()
	s.arcResultSub.Close()
	s.statSub.Close()
	if env.DeployEnv == env.DeployEnvProd {
		s.videoshotSub2.Close()
	}
	s.closed = true
	s.wg.Wait()
	s.redis.Close()
}

func (s *Service) isMission(c context.Context, aid int64) bool {
	if addit, _ := s.arc.Addit(c, aid); addit != nil && addit.MissionID > 0 {
		return true
	}
	return false
}

func (s *Service) isWhite(mid int64) bool {
	if ups, ok := s.upperCache[mngmdl.UpperTypeWhite]; ok {
		_, isWhite := ups[mid]
		return isWhite
	}
	return false
}

func (s *Service) isBlack(mid int64) bool {
	if ups, ok := s.upperCache[mngmdl.UpperTypeBlack]; ok {
		_, isBlack := ups[mid]
		return isBlack
	}
	return false
}

func (s *Service) isAuditType(tpID int16) bool {
	_, isAt := s.adtTpsCache[tpID]
	return isAt
}

func (s *Service) loadType() {
	tpMap, err := s.arc.TypeMapping(context.TODO())
	if err != nil {
		log.Error("s.arc.TypeMapping error(%v)", err)
		return
	}
	s.sfTpsCache = tpMap
	log.Info("s.sfTpsCache Data is (%+v)", s.sfTpsCache)
	tpNaming, err := s.arc.TypeNaming(context.TODO())
	if err != nil {
		log.Error("s.arc.TypeNaming error(%v)", err)
		return
	}
	s.TypeMap = tpNaming
	log.Info("s.TypeMap Data is (%+v)", s.TypeMap)
	// audit types
	adt, err := s.arc.AuditTypesConf(context.TODO())
	if err != nil {
		log.Error("s.arc.AuditTypesConf error(%v)", err)
		return
	}
	s.adtTpsCache = adt
	log.Info("s.adtTpsCache Data is (%+v)", s.adtTpsCache)
	// threshold
	thr, err := s.arc.ThresholdConf(context.TODO())
	if err != nil {
		log.Error("s.arc.ThresholdConf error(%v)", err)
		return
	}
	s.thrTpsCache = thr
	log.Info("s.thrTpsCache Data is (%+v)", s.thrTpsCache)
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
	var (
		c = context.TODO()
	)

	upm, err := s.mng.Uppers(c)
	if err != nil {
		log.Error("s.mng.Uppers error(%v)", err)
		return
	}
	s.upperCache = upm
}

func (s *Service) isRoundType(tpID int16) bool {
	_, in := s.roundTpsCache[tpID]
	return in
}

func (s *Service) loadConf() {
	var (
		fans       int64
		days       int64
		err        error
		roundTypes map[int16]struct{}
	)
	if fans, err = s.arc.FansConf(context.TODO()); err != nil {
		log.Error("s.arc.FansConf error(%v)", err)
		return
	}
	s.fansCache = int(fans)
	if roundTypes, err = s.arc.RoundTypeConf(context.TODO()); err != nil {
		log.Error("s.arc.RoundTypeConf error(%v)", err)
		return
	}
	s.roundTpsCache = roundTypes
	if days, err = s.arc.RoundEndConf(context.TODO()); err != nil {
		log.Error("s.arc.RoundEndConf")
		return
	}
	s.roundDelayCache = days
}

func (s *Service) cacheproc() {
	for {
		time.Sleep(1 * time.Minute)
		s.loadType()
		s.loadUpper()
		s.loadConf()
	}
}

func (s *Service) monitorConsume() {
	if env.DeployEnv != env.DeployEnvProd {
		return
	}
	var bvc2Vu, videoup, arcResult, stat, bvc2VuDelay int64
	for {
		time.Sleep(1 * time.Minute)
		if s.bvc2VuMo-bvc2Vu == 0 {
			s.monitor.Send(context.TODO(), "video-job bvc2Video did not consume within a minute")
		}
		if s.videoupMo-videoup == 0 {
			s.monitor.Send(context.TODO(), "video-job videoup did not consume within a minute")
		}
		if s.arcResultMo-arcResult == 0 {
			s.monitor.Send(context.TODO(), "video-job arcResult did not consume within a minute")
		}
		if s.statMo-stat == 0 {
			s.monitor.Send(context.TODO(), "video-job stat did not consume within a minute")
		}
		if s.bvc2VuDelayMo-bvc2VuDelay > 0 {
			s.monitor.Send(context.TODO(), "video-job bvc2videoup consume delayed.")
		}
		bvc2Vu = s.bvc2VuMo
		videoup = s.videoupMo
		arcResult = s.arcResultMo
		stat = s.statMo
		bvc2VuDelay = s.bvc2VuDelayMo
	}
}
