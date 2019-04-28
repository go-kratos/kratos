package data

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/up/model/datamodel"
	"go-common/app/admin/main/up/service/cache"
	"go-common/app/interface/main/creative/model/tag"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// ViewerBase get up viewer base data.
func (s *Service) ViewerBase(c context.Context, mid int64) (res *datamodel.ViewerBaseInfo, err error) {
	dt := getDateLastSunday()

	if res, err = s.data.ViewerBase(c, mid, dt); err != nil {
		log.Error("get viewer base error, err=%v", err)
		return
	}
	return
}

// ViewerArea get up viewer area data.
func (s *Service) ViewerArea(c context.Context, mid int64) (res *datamodel.ViewerAreaInfo, err error) {
	dt := getDateLastSunday()

	res, err = s.data.ViewerArea(c, mid, dt)
	if err != nil {
		log.Error("get viewer area error, err=%v", err)
		return
	}
	return
}

// CacheTrend get trend from mc.
func (s *Service) CacheTrend(c context.Context, mid int64) (res *datamodel.ViewerTrendInfo, err error) {
	dt := getDateLastSunday()

	if res, err = s.viewerTrend(c, mid, dt); err != nil {
		return
	}

	return
}

//GetTags get tag
func (s *Service) GetTags(c context.Context, ids ...int64) (result map[int64]*tag.Meta) {
	var tagMetaMap, leftIDs = cache.GetTagCache(ids...)
	if len(leftIDs) > 0 {
		if tlist, err := s.dtag.TagList(c, leftIDs); err != nil {
			log.Error("trend s.dtag.TagList err(%v)", err)
		} else {
			for _, v := range tlist {
				tagMetaMap[v.TagID] = v
				cache.AddTagCache(v)
			}
		}
	}
	result = tagMetaMap
	return
}

// ViewerTrend get up viewer trend data.
func (s *Service) viewerTrend(c context.Context, mid int64, dt time.Time) (res *datamodel.ViewerTrendInfo, err error) {
	ut, err := s.data.ViewerTrend(c, mid, dt)
	if err != nil || ut == nil {
		log.Error("trend s.data.ViewerTrend err(%v)", err)
		return
	}
	f := []string{"fan", "guest"}
	skeys := make([]int, 0) //for tag sort.
	tgs := make([]int64, 0) // for request tag name.
	res = &datamodel.ViewerTrendInfo{}
	var dataMap = make(map[string]*datamodel.ViewerTypeTagInfo)
	for _, fk := range f {
		td := ut[fk]
		vt := &datamodel.ViewerTypeTagInfo{}
		if td == nil {
			vt.Type = nil
			vt.Tag = nil
			dataMap[fk] = vt
			continue
		}
		typeMap := make(map[int]*datamodel.ViewerTypeData) //return type map to user.
		//deal type for type name.
		if td.Ty != nil {
			for k, v := range td.Ty {
				var tagInfo, ok = cache.VideoUpTypeCache[k]
				var name = ""
				if ok {
					name = tagInfo.Name
				}
				typeMap[k] = &datamodel.ViewerTypeData{
					Tid:  k,
					Name: name,
					Play: v,
				}
			}
		} else {
			typeMap = nil
		}
		// deal tag for tag name.
		if td.Tag != nil {
			for k, v := range td.Tag {
				tgs = append(tgs, v)
				skeys = append(skeys, k)
			}
			var tagMetaMap, leftIDs = cache.GetTagCache(tgs...)
			if len(leftIDs) > 0 {
				var tlist []*tag.Meta
				if tlist, err = s.dtag.TagList(c, leftIDs); err != nil {
					log.Error("trend s.dtag.TagList err(%v)", err)
				} else {
					for _, v := range tlist {
						tagMetaMap[v.TagID] = v
						cache.AddTagCache(v)
					}
				}
			}

			for _, k := range skeys {
				var tagID = td.Tag[k]
				var tagData = &datamodel.ViewerTagData{
					Idx:   k,
					TagID: int(tagID),
					Name:  "",
				}
				if tagMeta, ok := tagMetaMap[tagID]; ok {
					tagData.Name = tagMeta.TagName
				}
				vt.Tag = append(vt.Tag, tagData)
			}
		}
		for _, v := range typeMap {
			vt.Type = append(vt.Type, v)
		}

		dataMap[fk] = vt
	}
	res.Fans = dataMap["fan"]
	res.Guest = dataMap["guest"]
	return
}

//ViewData view data for up
type ViewData struct {
	Base  *datamodel.ViewerBaseInfo  `json:"viewer_base"`
	Trend *datamodel.ViewerTrendInfo `json:"viewer_trend"`
	Area  *datamodel.ViewerAreaInfo  `json:"viewer_area"`
}

//GetUpViewInfo get view info by arg
func (s *Service) GetUpViewInfo(c context.Context, arg *datamodel.GetUpViewInfoArg) (result *ViewData, err error) {
	if arg == nil {
		log.Error("arg is nil")
		return
	}
	res, err := s.GetViewData(c, arg.Mid)
	result = &res
	return
}

//GetViewData get all view data
func (s *Service) GetViewData(c context.Context, mid int64) (result ViewData, err error) {
	var group, ctx = errgroup.WithContext(c)
	group.Go(func() error {
		var e error
		result.Base, e = s.ViewerBase(ctx, mid)
		if e != nil {
			log.Error("get base view err happen, err=%v", e)
		}
		return nil
	})
	group.Go(func() error {
		var e error
		result.Area, e = s.ViewerArea(ctx, mid)
		if e != nil {
			log.Error("get area view err happen, err=%v", e)
		}
		return nil
	})
	group.Go(func() error {
		var e error
		result.Trend, e = s.CacheTrend(ctx, mid)
		if e != nil {
			log.Error("get trend view err happen, err=%v", e)
		}
		return nil
	})

	err = group.Wait()
	if err != nil {
		log.Error("get view data error, err=%+v", err)
	}
	return
}

//GetFansSummary get fan summary
func (s *Service) GetFansSummary(c context.Context, arg *datamodel.GetFansSummaryArg) (result *datamodel.FansSummaryResult, err error) {
	result = &datamodel.FansSummaryResult{}
	var fanInfo, e = s.data.UpFansAnalysis(c, arg.Mid, datamodel.Thirty)
	if e != nil || fanInfo == nil {
		err = e
		log.Error("get fan analysis err happen, err=%v", e)
		return
	}
	result.FanSummary = fanInfo.Summary
	return
}

//GetRelationFansDay get fan day history
func (s *Service) GetRelationFansDay(c context.Context, arg *datamodel.GetRelationFansHistoryArg) (result *datamodel.GetRelationFansHistoryResult, err error) {

	result = &datamodel.GetRelationFansHistoryResult{}
	if arg.DataType == 0 {
		arg.DataType = datamodel.DataType30Day
	}
	switch arg.DataType {
	case datamodel.DataType30Day:
		result.RelationFanHistoryData, err = s.data.RelationFansDay(c, arg.Mid)
	case datamodel.DataTypeMonth:
		result.RelationFanHistoryData, err = s.data.RelationFansMonth(c, arg.Mid)
	default:
		err = fmt.Errorf("invalid data type(%d)", arg.DataType)
	}
	if err != nil {
		log.Error("get relation fans day fail, err=%v", err)
		return
	}
	return
}
