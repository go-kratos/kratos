package v2

import (
	"context"
	"time"

	"go-common/library/log"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
)

// 获取分区入口
func (s *IndexService) getAreaEntrance(ctx context.Context) (res []*v2pb.MAreaEntrance) {
	moduleInfoMap := s.GetAllModuleInfoMapFromCache(ctx)
	listMap := s.getAreaEntranceListMapFromCache(ctx)
	res = make([]*v2pb.MAreaEntrance, 0)
	for _, m := range moduleInfoMap[_entranceType] {
		if l, ok := listMap[m.Id]; ok {
			res = append(res, &v2pb.MAreaEntrance{
				ModuleInfo: m,
				List:       l,
			})
		}
	}

	return
}

// load from cache
func (s *IndexService) getAreaEntranceListMapFromCache(ctx context.Context) (res map[int64][]*v2pb.PicItem) {
	// load
	i := s.areaEntranceListMap.Load()
	// assert
	res, ok := i.(map[int64][]*v2pb.PicItem)
	if ok {
		return
	}
	// 回源&log
	res = s.getAreaEntranceListMap(ctx)
	log.Warn("[getAreaEntranceListMapFromCache]memory cache miss!! i:%+v; res:%+v", i, res)
	return
}

// getAreaEntranceListMap raw
func (s *IndexService) getAreaEntranceListMap(ctx context.Context) (listMap map[int64][]*v2pb.PicItem) {
	moduleIds := s.getIdsFromModuleMap(ctx, []int64{_entranceType})
	if len(moduleIds) <= 0 {
		return
	}
	areaResult, err := s.roomDao.GetAreaEntrance(ctx, moduleIds)
	if err != nil {
		log.Error("[loadAreaEntranceCache]roomDao.GetAreaEntrance get data error: %+v, data: %+v", err, areaResult)
		return
	}
	if len(areaResult) > 0 {
		listMap = make(map[int64][]*v2pb.PicItem)
		for moduleId, i := range areaResult {
			if i != nil && i.List != nil {
				for _, ii := range i.List {
					listMap[moduleId] = append(listMap[moduleId], &v2pb.PicItem{
						Id:    ii.Id,
						Pic:   ii.Pic,
						Link:  ii.Link,
						Title: ii.Title,
					})
				}
			}
		}
	}

	return
}

// ticker
func (s *IndexService) areaEntranceProc() {
	for {
		time.Sleep(time.Minute * 1)
		s.loadAreaEntranceCache()
	}
}
func (s *IndexService) loadAreaEntranceCache() {
	areaEntranceListMap := s.getAreaEntranceListMap(context.TODO())
	if len(areaEntranceListMap) > 0 {
		s.areaEntranceListMap.Store(areaEntranceListMap)
		log.Info("[loadAreaEntranceCache]load data success!")
	}
	return
}
