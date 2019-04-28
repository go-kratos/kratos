package feed

import (
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/app/conf"
	feeddao "go-common/app/job/main/app/dao/feed"
	"go-common/app/job/main/app/model/feed"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// Service is show service.
type Service struct {
	c       *conf.Config
	dao     *feeddao.Dao
	tick    time.Duration
	cacheCh chan func()
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		dao:     feeddao.New(c),
		tick:    time.Duration(c.Tick),
		cacheCh: make(chan func(), 1024),
	}
	s.loadRcmdCache()
	go s.loadproc()
	go s.cacheproc()
	return
}

func (s *Service) addCache(f func()) {
	select {
	case s.cacheCh <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

func (s *Service) cacheproc() {
	for {
		f, ok := <-s.cacheCh
		if !ok {
			log.Warn("cache proc exit")
			return
		}
		f()
	}
}

func (s *Service) loadRcmdCache() {
	var (
		c    = context.Background()
		now  = time.Now()
		aids []int64
		is   []*feed.RcmdItem
		err  error
	)
	if aids, err = s.dao.Hots(c); err != nil {
		log.Error("%+v", err)
		return
	}
	is = s.fromAids(c, aids, now)
	if len(is) == 0 {
		return
	}
	s.addCache(func() {
		s.dao.UpRcmdCache(c, is)
	})
}

// cacheproc load all cache.
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.loadRcmdCache()
	}
}

func (s *Service) fromAids(c context.Context, aids []int64, now time.Time) (is []*feed.RcmdItem) {
	if len(aids) == 0 {
		return
	}
	const _count = 50
	var (
		am    map[int64]*api.Arc
		shard int
		err   error
	)
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		am, err = s.dao.Archives(ctx, aids, "")
		return
	})
	if len(aids) < _count {
		shard = 1
	} else {
		shard = len(aids) / _count
		if len(aids)%(shard*_count) != 0 {
			shard++
		}
	}
	aidss := make([][]int64, shard)
	for i, aid := range aids {
		aidss[i%shard] = append(aidss[i%shard], aid)
	}
	tagms := make([]map[string][]*feed.Tag, len(aidss))
	for i, aids := range aidss {
		if len(aids) == 0 {
			continue
		}
		g.Go(func() (err error) {
			if tagms[i], err = s.dao.Tags(ctx, aids, now); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(am) == 0 {
		return
	}
	tagm := make(map[string][]*feed.Tag, len(aids))
	for _, tm := range tagms {
		for aid, tag := range tm {
			tagm[aid] = tag
		}
	}
	is = make([]*feed.RcmdItem, 0, len(am))
	for _, aid := range aids {
		i := &feed.RcmdItem{}
		if a, ok := am[aid]; ok {
			i.ID = a.Aid
			i.Archive = a
		}
		if ts, ok := tagm[strconv.FormatInt(aid, 10)]; ok {
			if len(ts) != 0 {
				i.Tid = ts[0].ID
				i.Tag = ts[0]
			}
		}
		is = append(is, i)
	}
	return
}
