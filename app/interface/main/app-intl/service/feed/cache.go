package feed

import (
	"context"
	"hash/crc32"
	"math/rand"
	"time"

	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-intl/model"
	"go-common/library/log"
)

// indexCache is.
func (s *Service) indexCache(c context.Context, mid int64, count int) (rs []*ai.Item, err error) {
	var (
		pos, nextPos int
	)
	cache := s.rcmdCache
	if len(cache) < count {
		return
	}
	if pos, err = s.rcmd.PositionCache(c, mid); err != nil {
		return
	}
	rs = make([]*ai.Item, 0, count)
	if pos < len(cache)-count-1 {
		nextPos = pos + count
		rs = append(rs, cache[pos:nextPos]...)
	} else if pos < len(cache)-1 {
		nextPos = count - (len(cache) - pos)
		rs = append(rs, cache[pos:]...)
		rs = append(rs, cache[:nextPos]...)
	} else {
		nextPos = count - 1
		rs = append(rs, cache[:nextPos]...)
	}
	s.addCache(func() {
		s.rcmd.AddPositionCache(context.Background(), mid, nextPos)
	})
	return
}

// recommendCache is.
func (s *Service) recommendCache(count int) (rs []*ai.Item) {
	cache := s.rcmdCache
	index := len(cache)
	if count > 0 && count < index {
		index = count
	}
	rs = make([]*ai.Item, 0, index)
	for _, idx := range rand.Perm(len(cache))[:index] {
		rs = append(rs, cache[idx])
	}
	return
}

// group is.
func (s *Service) group(mid int64, buvid string) (group int) {
	if mid == 0 && buvid == "" {
		group = -1
		return
	}
	if mid != 0 {
		if v, ok := s.groupCache[mid]; ok {
			group = v
			return
		}
		group = int(mid % 20)
		return
	}
	group = int(crc32.ChecksumIEEE([]byte(buvid)) % 20)
	return
}

// loadRcmdCache is.
func (s *Service) loadRcmdCache() {
	is, err := s.rcmd.RcmdCache(context.Background())
	if err != nil {
		log.Error("%+v", err)
	}
	if len(is) >= 50 {
		for _, i := range is {
			i.Goto = model.GotoAv
		}
		s.rcmdCache = is
		return
	}
	aids, err := s.rcmd.Hots(context.Background())
	if err != nil {
		log.Error("%+v", err)
	}
	if len(aids) == 0 {
		if aids, err = s.rcmd.RcmdAidsCache(context.Background()); err != nil {
			return
		}
	}
	if len(aids) < 50 && len(s.rcmdCache) != 0 {
		return
	}
	s.rcmdCache = s.fromAids(aids)
	s.addCache(func() {
		s.rcmd.AddRcmdAidsCache(context.Background(), aids)
	})
}

// fromAids is.
func (s *Service) fromAids(aids []int64) (is []*ai.Item) {
	is = make([]*ai.Item, 0, len(aids))
	for _, aid := range aids {
		i := &ai.Item{
			ID:   aid,
			Goto: model.GotoAv,
		}
		is = append(is, i)
	}
	return
}

// rcmdproc is.
func (s *Service) rcmdproc() {
	for {
		time.Sleep(s.tick)
		s.loadRcmdCache()
	}
}

// loadRankCache is.
func (s *Service) loadRankCache() {
	rank, err := s.rank.AllRank(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.rankCache = rank
}

// rankproc is.
func (s *Service) rankproc() {
	for {
		time.Sleep(s.tick)
		s.loadRankCache()
	}
}

// loadUpCardCache is.
func (s *Service) loadUpCardCache() {
	follow, err := s.card.Follow(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.followCache = follow
}

// upCardproc is.
func (s *Service) upCardproc() {
	for {
		time.Sleep(s.tick)
		s.loadUpCardCache()
	}
}

// loadGroupCache is.
func (s *Service) loadGroupCache() {
	group, err := s.rcmd.Group(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.groupCache = group
}

// groupproc is.
func (s *Service) groupproc() {
	for {
		time.Sleep(s.tick)
		s.loadGroupCache()
	}
}
