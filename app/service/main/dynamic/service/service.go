package service

import (
	"context"
	"fmt"
	"time"

	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/dynamic/conf"
	"go-common/app/service/main/dynamic/dao"
	"go-common/library/cache"
	"go-common/library/log"
)

// Service service.
type Service struct {
	dao *dao.Dao
	c   *conf.Config
	// new dynamic arcs.
	hotRidTids    map[int32][]int64
	regionTotal   map[int32]int
	regionArcs    map[int32][]int64
	regionTagArcs map[string][]int64
	// live
	live int
	// rpc
	arcRPC *arcrpc.Service2
	// cache
	cache *cache.Cache
}

// New service new.
func New(c *conf.Config) *Service {
	s := &Service{
		dao: dao.New(c),
		c:   c,
		// rpc
		arcRPC: arcrpc.New2(c.ArchiveRPC),
		// new dynamic arcs
		hotRidTids:    make(map[int32][]int64),
		regionTotal:   make(map[int32]int),
		regionArcs:    make(map[int32][]int64),
		regionTagArcs: make(map[string][]int64),
		cache:         cache.New(1, 1024),
	}
	go s.regionproc()
	go s.tagproc()
	return s
}

func regionTagKey(rid int32, tagID int64) string {
	return fmt.Sprintf("%d_%d", rid, tagID)
}

// regionproc is a routine for pull region dynamic into cache.
func (s *Service) regionproc() {
	var (
		c           = context.TODO()
		res         map[int32][]int64
		rids        []int32
		cacheRegion map[int32][]int64
		err         error
	)
	for {
		// load hot tags from tag api.
		regionTotal := make(map[int32]int)
		regionArcs := make(map[int32][]int64)
		if res, err = s.dao.Hot(c); err != nil {
			log.Error("dao.Hot() error(%v)", err)
			time.Sleep(time.Second)
			continue
		}
		if len(res) > 0 {
			s.hotRidTids = res
		}
		if rids, err = s.dao.Rids(c); err != nil {
			log.Error("dao.Rids() error(%v)", err)
			time.Sleep(time.Second)
			continue
		}
		//get region cache
		cacheRegion = s.dao.RegionCache(c)
		// init dynamic arcs from bigdata.
		for rid := range s.hotRidTids {
			// init region dynamic arcs.
			if aids, total, err := s.dao.RegionArcs(c, rid, ""); err != nil || len(aids) < conf.Conf.Rule.MinRegionCount {
				regionTotal[rid] = s.regionTotal[rid]
				if len(s.regionArcs[rid]) >= conf.Conf.Rule.MinRegionCount {
					regionArcs[rid] = s.regionArcs[rid]
				} else if cacheRegion != nil && len(cacheRegion[rid]) >= conf.Conf.Rule.MinRegionCount {
					regionArcs[rid] = cacheRegion[rid]
				} else {
					dao.PromError("热门动态数据错误", "dynamic data error rid(%d)  bigdata(%d)  memory(%d)", rid, len(aids), len(s.regionArcs[rid]))
					if len(aids) > 0 {
						regionArcs[rid] = aids
					} else {
						regionArcs[rid] = s.regionArcs[rid]
					}
				}
			} else {
				regionTotal[rid] = total
				regionArcs[rid] = aids
			}
		}
		for _, rid := range rids {
			if aids, total, err := s.dao.RegionArcs(c, rid, ""); err != nil || len(aids) < conf.Conf.Rule.MinRegionCount {
				regionTotal[rid] = s.regionTotal[rid]
				if len(s.regionArcs[rid]) >= conf.Conf.Rule.MinRegionCount {
					regionArcs[rid] = s.regionArcs[rid]
				} else if cacheRegion != nil && len(cacheRegion[rid]) >= conf.Conf.Rule.MinRegionCount {
					regionArcs[rid] = cacheRegion[rid]
				} else {
					dao.PromError("分区动态数据错误", "dynamic data error rid(%d)  bigdata(%d)  memory(%d)", rid, len(aids), len(s.regionArcs[rid]))
					if len(aids) > 0 {
						regionArcs[rid] = aids
					} else {
						regionArcs[rid] = s.regionArcs[rid]
					}
				}
			} else {
				regionTotal[rid] = total
				regionArcs[rid] = aids
			}
		}
		if count, err := s.dao.Live(c); err != nil {
			log.Error("s.dao.Live() error(%v)", err)
		} else {
			s.live = count
		}
		s.regionTotal = regionTotal
		s.regionArcs = regionArcs
		if regionNeedCache(s.regionArcs) {
			s.cache.Save(func() {
				s.dao.SetRegionCache(context.TODO(), s.regionArcs)
			})
		}
		time.Sleep(time.Duration(s.c.Rule.TickRegion))
	}
}

// tagproc is a routine for pull tag dynamic into cache.
func (s *Service) tagproc() {
	var (
		c            = context.TODO()
		cacheTag     map[string][]int64
		tagNeedCache bool
	)
	for {
		//get tag cache
		cacheTag = s.dao.TagCache(c)
		// load hot tags from tag api.
		regionTagArcs := make(map[string][]int64)
		// init dynamic arcs from bigdata.
		for rid, tids := range s.hotRidTids {
			// init region tag dynamic arcs.
			for _, tid := range tids {
				k := regionTagKey(rid, tid)
				if aids, err := s.dao.RegionTagArcs(c, rid, tid, ""); err != nil || len(aids) == 0 {
					if len(s.regionTagArcs[k]) == 0 {
						if cacheTag != nil && len(cacheTag[k]) > 0 {
							regionTagArcs[k] = cacheTag[k]
						}
						tagNeedCache = false || tagNeedCache
					} else {
						regionTagArcs[k] = s.regionTagArcs[k]
						tagNeedCache = true
					}
				} else {
					regionTagArcs[k] = aids
					tagNeedCache = true
				}
			}
		}
		s.regionTagArcs = regionTagArcs
		if tagNeedCache {
			s.cache.Save(func() {
				s.dao.SetTagCache(context.TODO(), s.regionTagArcs)
			})
		}
		time.Sleep(time.Duration(s.c.Rule.TickTag))
	}
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}

func archivesLog(name string, aids []int64) {
	if aidLen := len(aids); aidLen >= 50 {
		log.Info("s.archives3 func(%s) len(%d), arg(%v)", name, aidLen, aids)
	}
}
