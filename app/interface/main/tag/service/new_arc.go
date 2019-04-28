package service

import (
	"context"

	"go-common/app/service/main/archive/api"
	rpcModel "go-common/app/service/main/tag/model"
	"go-common/library/log"
)

// newArcs .
func (s *Service) newArcs(c context.Context, tid int64, start, end int) (as []int64, count int, err error) {
	if as, count, err = s.dao.NewArcsCache(c, tid, start, end); err != nil {
		log.Error("d.NewArcsCache(%d,%d,%d) err(%v)", tid, start, end, err)
	}
	if len(as) > 0 {
		return
	}
	var aids []int64
	if aids, err = s.resOidsByTid(c, tid, rpcModel.ResTypeArchive); err != nil {
		return
	}
	res, err := s.batchArchives(c, aids)
	if err != nil {
		log.Error("d.BatchArchives() err(%v)", err)
		return
	}
	var arcs []*api.Arc
	for k, v := range res {
		if v.IsNormal() {
			as = append(as, k)
			arcs = append(arcs, v)
		}
	}
	count = len(as)
	if count == 0 || count < start {
		as = []int64{}
	} else if end+1 > count {
		as = as[start:]
	} else {
		as = as[start : end+1]
	}
	s.dao.AddNewArcsCache(c, tid, arcs...)
	return
}

func (s *Service) tagPridAids(c context.Context, tid, prid int64, start, end int) (as []int64, count int, err error) {
	var aids []int64
	if as, count, err = s.dao.ZrangeTagPridArc(c, tid, prid, start, end); err != nil {
		log.Error("d.ZrangeTagPridArc(%d,%d,%d,%d) err(%v)", tid, prid, start, end, err)
		return
	}
	if len(as) > 0 {
		return
	}
	if aids, err = s.resOidsByTid(c, tid, rpcModel.ResTypeArchive); err != nil {
		return
	}
	var res map[int64]*api.Arc
	if res, err = s.batchArchives(c, aids); err != nil {
		log.Error("d.BatchArchives() err(%v)", err)
		return
	}
	pridArcMap := make(map[int64][]*api.Arc, 100)
	for _, aid := range aids {
		if arc, ok := res[aid]; ok {
			if arc.IsNormal() {
				if tempPrid, ok := s.pridMap[int64(arc.TypeID)]; ok {
					pridArcMap[tempPrid] = append(pridArcMap[tempPrid], arc)
					if tempPrid == prid {
						as = append(as, aid)
					}
				}
			}
		}
	}
	count = len(as)
	if count == 0 || count < start {
		as = []int64{}
	} else if end+1 > count {
		as = as[start:]
	} else {
		as = as[start : end+1]
	}
	if len(as) == 0 {
		var arc = &api.Arc{Aid: -1}
		pridArcMap[prid] = []*api.Arc{arc}
	}
	for prid, arcs := range pridArcMap {
		s.dao.AddTagPridCache(context.Background(), []int64{tid}, prid, arcs...)
	}
	return
}

//  过审过来的绑定 .
func (s *Service) auditBindTag(c context.Context, aid int64, arc *api.Arc, tids []int64) (err error) {
	if err = s.dao.AddNewArcCache(c, arc, tids...); err != nil {
		log.Error("s.dao.AddNewArcCache(aid:%d) err(%v)", aid, err)
	}
	// tags of one region
	if err = s.dao.AddRegionNewArcCache(c, arc.TypeID, arc, tids...); err != nil {
		log.Error("s.dao.AddRegionNewArcCache(%d) err(%v)", aid, err)
	}
	// tag 一级分区视频
	if prid, ok := s.pridMap[int64(arc.TypeID)]; ok {
		if err = s.dao.AddTagPridArcCache(c, tids, prid, arc); err != nil {
			log.Error("s.dao.AddTagPridArcCache(%v,%d,%v) err(%v)", tids, prid, arc, err)
		}
	}
	return
}

//  添加视频tags相关 . TagsArcRem  remTagsArc
func (s *Service) remTagsArc(c context.Context, aid int64, arc *api.Arc, tids []int64) (err error) {
	// tags of all region
	if err = s.dao.RemoveNewArcsCache(c, aid, tids...); err != nil {
		log.Error("d.RemoveNewArcsCache(tid:%d) err(%v)", tids, err)
	}
	// tags of one region
	if err = s.dao.RemoveRegionNewArcCache(c, arc.TypeID, arc, tids...); err != nil {
		log.Error("d.RemoveRegionNewArcCache(tid:%d) err(%v)", tids, err)
	}
	// tag 一级分区视频
	if prid, ok := s.pridMap[int64(arc.TypeID)]; ok {
		if err = s.dao.RemoveTagPridArcCache(c, tids, prid, arc.Aid); err != nil {
			log.Error("d.RemoveTagPridArcCache(%v,%d,%d) err(%v)", tids, prid, arc.Aid, err)
		}
	}
	return
}
