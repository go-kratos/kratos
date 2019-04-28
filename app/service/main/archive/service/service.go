package service

import (
	"context"
	"runtime"
	"time"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/conf"
	arcdao "go-common/app/service/main/archive/dao/archive"
	shareDao "go-common/app/service/main/archive/dao/share"
	shotDao "go-common/app/service/main/archive/dao/videoshot"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

var (
	_emptyArchives3 = make(map[int64]*api.Arc)
)

// Service is service.
type Service struct {
	c *conf.Config
	// dao
	arc   *arcdao.Dao
	share *shareDao.Dao
	shot  *shotDao.Dao
	// acc rpc
	// acc *accrpc.Service2
	// types
	allTypes  map[int16]*archive.ArcType
	ridToReid map[int16]int16
	bnjList   map[int64]struct{}
	// cache chan
	cacheCh chan func()
	// prom
	hitProm  *prom.Prom
	missProm *prom.Prom
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		// dao
		arc:   arcdao.New(c),
		share: shareDao.New(c),
		shot:  shotDao.New(c),
		// acc rpc
		// acc: accrpc.New2(c.AccountRPC),
		// types
		allTypes:  make(map[int16]*archive.ArcType),
		ridToReid: make(map[int16]int16),
		// cache chan
		cacheCh: make(chan func(), 1024),
		// prom
		hitProm:  prom.CacheHit,
		missProm: prom.CacheMiss,
		bnjList:  make(map[int64]struct{}),
	}
	s.loadBnjList()
	s.loadTypes()
	go s.loadproc()
	for i := 0; i < runtime.NumCPU(); i++ {
		go s.cacheproc()
	}
	return
}

// AllTypes return all types
func (s *Service) AllTypes(c context.Context) (types map[int16]*archive.ArcType) {
	types = s.allTypes
	return
}

// CacheUpdate job update/del/add archive cache
func (s *Service) CacheUpdate(c context.Context, aid int64, tp string, oldMid int64) (err error) {
	if err = s.arc.UpArchiveCache(c, aid); err != nil {
		log.Error("s.arc.UpArchiveCache(%d) error(%v)", aid, err)
	}
	if err = s.arc.InitStatCache3(c, aid); err != nil {
		log.Error("s.arc.InitStatCache3(%d) error(%v)", aid, err)
	}
	if oldMid != 0 {
		if err = s.DelUpperPassedCache(c, aid, oldMid); err != nil {
			log.Error("s.DelUpperPassedCache(%d, %d) error(%v)", aid, oldMid)
		}
	}
	switch tp {
	case archive.CacheAdd:
		if err = s.AddUpperPassedCache(c, aid); err != nil {
			log.Error("s.AddUpperPassedCache(%d) error(%v)", aid, err)
		}
		if err = s.AddRegionArc(c, aid); err != nil {
			log.Error("s.AddRegionArc(%d) error(%v)", aid, err)
		}
	case archive.CacheUpdate:
		// NOTE: nothing todo
	case archive.CacheDelete:
		if err = s.DelUpperPassedCache(c, aid, 0); err != nil {
			log.Error("s.DelUpperPassedCache(%d) error(%v)", aid, err)
		}
		if err = s.DelRegionArc(c, aid, 0); err != nil {
			log.Error("s.DelRegionArc(%d) error(%v)", aid, err)
		}
	default:
		// NOTE: nothing todo
	}
	return
}

// FieldCacheUpdate job update field cache
func (s *Service) FieldCacheUpdate(c context.Context, aid int64, oldType, nwType int16) (err error) {
	if nwType != 0 {
		if err = s.AddRegionArc(c, aid); err != nil {
			log.Error("s.AddRegionArc(%d) error(%v)", aid, err)
		}
	}
	if oldType != 0 {
		if err = s.DelRegionArc(c, aid, oldType); err != nil {
			log.Error("s.DelRegionArc(%d,%d) error(%d)", aid, oldType, err)
		}
	}
	return
}

// Ping ping success.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.arc.Ping(c)
	return
}

// Close resource.
func (s *Service) Close() {
	s.arc.Close()
}

func (s *Service) loadBnjList() {
	bnjList := make(map[int64]struct{})
	for _, aid := range s.c.BnjList {
		bnjList[aid] = struct{}{}
	}
	s.bnjList = bnjList
}

func (s *Service) loadTypes() {
	var (
		ridToReid = make(map[int16]int16)
		types     map[int16]*archive.ArcType
		err       error
	)
	if types, err = s.arc.Types(context.TODO()); err != nil {
		log.Error("s.arc.Types error(%v)", err)
		return
	}
	for _, t := range types {
		if t.Pid != 0 {
			ridToReid[t.ID] = t.Pid
		}
	}
	s.allTypes = types
	s.ridToReid = ridToReid
}

func (s *Service) loadproc() {
	for {
		time.Sleep(time.Duration(s.c.Tick))
		s.loadTypes()
		s.loadBnjList()
	}
}

func (s *Service) addCache(f func()) {
	select {
	case s.cacheCh <- f:
	default:
		log.Warn("s.cacheCh is full")
	}
}

func (s *Service) cacheproc() {
	for {
		f, ok := <-s.cacheCh
		if !ok {
			return
		}
		f()
	}
}
