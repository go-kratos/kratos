package service

import (
	"context"
	"time"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

// RegionTopCount top region count when one day.
func (s *Service) RegionTopCount(c context.Context, reids []int16) (res map[int16]int, err error) {
	var (
		t   = time.Now()
		min = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
		max = t.Unix()
	)
	res, err = s.arc.RegionTopCountCache(c, reids, min, max)
	return
}

// DelRegionArc delete a archive cache by aid.
func (s *Service) DelRegionArc(c context.Context, aid int64, rid int16) (err error) {
	if rid == 0 {
		var a *api.Arc
		if a, err = s.arc.Archive3(c, aid); err != nil {
			log.Error("s.arc.Archive(%d) error(%v)", aid, err)
			return
		}
		rid = int16(a.TypeID)
	}
	if err = s.arc.DelRegionArcCache(c, rid, s.ridToReid[rid], aid); err != nil {
		log.Error("s.arc.DelRegionArcCache(%d) error(%v)", aid, err)
	}
	return
}

// AddRegionArc add a archive cache by aid.
func (s *Service) AddRegionArc(c context.Context, aid int64) (err error) {
	var a *api.Arc
	if a, err = s.arc.Archive3(c, aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", aid, err)
		return
	}
	if !a.IsNormal() {
		return
	}
	var ra = &api.RegionArc{Aid: aid, Attribute: a.Attribute, Copyright: int8(a.Copyright), PubDate: a.PubDate}
	if err = s.arc.AddRegionArcCache(c, int16(a.TypeID), s.ridToReid[int16(a.TypeID)], ra); err != nil {
		log.Error("s.arc.AddRegionArcCache(%v) error(%v)", ra, err)
	}
	return
}

// AddRegionArcs load region arcs into redis.
func (s *Service) AddRegionArcs(c context.Context, rid int16) (err error) {
	// get all rids
	var (
		rids  []int16
		types map[int16]*archive.ArcType
	)
	if rid > 0 {
		rids = []int16{rid}
	} else {
		if types, err = s.arc.Types(context.Background()); err != nil {
			log.Error("s.arc.Types() error(%v)", err)
			return
		}
		rids = make([]int16, 0)
		for _, t := range types {
			if t.Pid != 0 {
				rids = append(rids, t.ID)
			}
		}
	}
NEXT:
	for _, rid := range rids {
		var (
			start  = 0
			length = 5000 // 100 everytime
		)
		for {
			var ras []*api.RegionArc
			// set all arcs of rid into redis.
			if ras, err = s.arc.RegionArcs(context.Background(), rid, start, length); err != nil {
				log.Error("s.arc.RegionArcs() error(%v)", err)
				return
			}
			if len(ras) == 0 {
				break NEXT
			}
			if err = s.arc.AddRegionArcCache(context.Background(), rid, s.ridToReid[rid], ras...); err != nil {
				log.Error("s.arc.AddRegionArcCache(%d) error(%v)", rid, err)
				return
			}
			start += length
			log.Info("init rid(%d) name(%s) now(%d)", rid, types[rid], start)
		}
	}
	return
}
