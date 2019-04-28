package service

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

// RegionTopArcs3 top region archives order by pubtime.
func (s *Service) RegionTopArcs3(c context.Context, reid int16, pn, ps int) (as []*api.Arc, err error) {
	var (
		start = (pn - 1) * ps
		end   = start + ps - 1
		aids  []int64
		am    map[int64]*api.Arc
	)
	if aids, err = s.arc.RegionTopArcsCache(c, reid, start, end); err != nil {
		log.Error("s.arc.RegionTopArcsCache(%v, %d, %d) error(%v)", reid, start, end, err)
		return
	}
	if am, err = s.arc.Archives3(c, aids); err != nil {
		log.Error("s.Archives(%v) error(%v)", aids, err)
		return
	}
	as = make([]*api.Arc, 0, len(aids))
	for _, aid := range aids {
		if a, ok := am[aid]; ok {
			as = append(as, a)
		}
	}
	return
}

// RegionAllArcs3 get left 7 days all region arcs
func (s *Service) RegionAllArcs3(c context.Context, pn, ps int) (as *archive.RankArchives3, err error) {
	var count int
	as = &archive.RankArchives3{}
	if count, err = s.arc.RegionAllCountCache(c); err != nil {
		err = nil
		log.Error("s.arc.RegionAllCountCache error(%v)", err)
		return
	}
	as.Count = count
	var (
		aids  []int64
		start = (pn - 1) * ps
		end   = start + ps - 1
		am    map[int64]*api.Arc
	)
	if aids, err = s.arc.RegionTopArcsCache(c, 0, start, end); err != nil {
		return
	}
	if am, err = s.arc.Archives3(c, aids); err != nil {
		return
	}
	as.Archives = make([]*api.Arc, 0)
	for _, aid := range aids {
		if a, ok := am[aid]; ok {
			as.Archives = append(as.Archives, a)
		}
	}
	return
}

// RegionArcs3 region archives.
func (s *Service) RegionArcs3(c context.Context, rid int16, pn, ps int) (as []*api.Arc, count int, err error) {
	if count, err = s.arc.RegionCountCache(c, rid); err != nil {
		log.Error("s.arc.RegionCountCache(%d) error(%v)", rid, err)
		err = nil
	}
	var (
		aids  []int64
		start = (pn - 1) * ps
		end   = start + ps - 1
		am    map[int64]*api.Arc
	)
	if aids, err = s.arc.RegionArcsCache(c, rid, start, end); err != nil {
		log.Error("s.arc.RegionArcsCache(rid:%d, start:%d, end:%d) error(%v)", rid, start, end, err)
		return
	}
	if am, err = s.arc.Archives3(c, aids); err != nil {
		log.Error("s.arc.Archives(%v) error(%v)", aids, err)
		return
	}
	as = make([]*api.Arc, 0, len(aids))
	for _, aid := range aids {
		if a, ok := am[aid]; ok {
			as = append(as, a)
		}
	}
	return
}

// RegionOriginArcs3 region origin archives.
func (s *Service) RegionOriginArcs3(c context.Context, rid int16, pn, ps int) (as []*api.Arc, count int, err error) {
	if count, err = s.arc.RegionOriginCountCache(c, rid); err != nil {
		log.Error("s.arc.RegionOriginCountCache(%d) error(%v)", rid, err)
		err = nil
	}
	var (
		aids  []int64
		start = (pn - 1) * ps
		end   = start + ps - 1
		am    map[int64]*api.Arc
	)
	if aids, err = s.arc.RegionOriginArcsCache(c, rid, start, end); err != nil {
		log.Error("s.arc.RegionOriginArcsCache(rid:%d, start:%d, end:%d) error(%v)", rid, start, end, err)
		return
	}
	if am, err = s.arc.Archives3(c, aids); err != nil {
		log.Error("s.arc.Archives(%v) error(%v)", aids, err)
		return
	}
	as = make([]*api.Arc, 0, len(aids))
	for _, aid := range aids {
		if a, ok := am[aid]; ok {
			as = append(as, a)
		}
	}
	return
}
