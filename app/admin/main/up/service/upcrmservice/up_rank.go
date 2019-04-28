package upcrmservice

import (
	"context"
	"sort"
	"time"

	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/app/admin/main/up/util/mathutil"
	"go-common/library/log"
	xtime "go-common/library/time"
)

var (
	//AllRankTyps all rank types
	AllRankTyps = []int{
		upcrmmodel.UpRankTypeFans30day1k,
		upcrmmodel.UpRankTypeFans30day1w,
		upcrmmodel.UpRankTypePlay30day1k,
		upcrmmodel.UpRankTypePlay30day1w,
		upcrmmodel.UpRankTypePlay30day10k,
		upcrmmodel.UpRankTypeFans30dayIncreaseCount,
		upcrmmodel.UpRankTypeFans30dayIncreasePercent,
	}
)

type sortRankFunc func(p1, p2 *upcrmmodel.UpRankInfo) bool

type upRankSorter struct {
	datas []*upcrmmodel.UpRankInfo
	by    sortRankFunc // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *upRankSorter) Len() int {
	return len(s.datas)
}

// Swap is part of sort.Interface.
func (s *upRankSorter) Swap(i, j int) {
	s.datas[i], s.datas[j] = s.datas[j], s.datas[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *upRankSorter) Less(i, j int) bool {
	return s.by(s.datas[i], s.datas[j])
}

func sortRankInfo(planets []*upcrmmodel.UpRankInfo, sortfunc sortRankFunc) {
	ps := &upRankSorter{
		datas: planets,
		by:    sortfunc, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

func sortByValueAsc(p1, p2 *upcrmmodel.UpRankInfo) bool {
	if p1.Value < p2.Value {
		return true
	}
	if p1.Value == p2.Value {
		return p1.Value2 < p2.Value2
	}
	return false
}

func sortByValueDesc(p1, p2 *upcrmmodel.UpRankInfo) bool {
	if p1.Value > p2.Value {
		return true
	}
	if p1.Value == p2.Value {
		return p1.Value2 > p2.Value2
	}
	return false
}

// 从数据库中更新数据
func (s *Service) refreshUpRankDate(date time.Time) {
	rankData, err := s.crmdb.QueryUpRankAll(date)
	if err != nil {
		log.Error("refresh from db fail, err=%+v", err)
		return
	}
	var typeMap = map[int][]*upcrmmodel.UpRankInfo{}
	var upInfoMap = map[int64]*upcrmmodel.InfoQueryResult{}
	for _, v := range rankData {
		var rankType = int(v.Type)
		var rankInfo = &upcrmmodel.UpRankInfo{}
		rankInfo.Mid = v.Mid
		rankInfo.CopyFromUpRank(&v)
		typeMap[rankType] = append(typeMap[rankType], rankInfo)
		upInfoMap[v.Mid] = nil
	}
	log.Info("refresh from db, get for date=%v, len=%d", date, len(rankData))
	// interestring code
	var mids []int64
	for k := range upInfoMap {
		mids = append(mids, k)
	}

	// 查询并合并
	var infoData, e = s.upBaseInfoQueryBatch(s.crmdb.QueryUpBaseInfoBatchByMid, mids...)
	err = e
	if err != nil {
		log.Error("get from base info fail, err=%+v", err)
		return
	}
	for _, v := range infoData {
		upInfoMap[v.Mid] = v
	}

	for _, list := range typeMap {
		for _, info := range list {
			var data, _ = upInfoMap[info.Mid]
			if data == nil {
				continue
			}
			info.InfoQueryResult = *data
			switch info.RankType {
			case upcrmmodel.UpRankTypeFans30day1k,
				upcrmmodel.UpRankTypeFans30day1w,
				upcrmmodel.UpRankTypePlay30day1k,
				upcrmmodel.UpRankTypePlay30day1w,
				upcrmmodel.UpRankTypePlay30day10k,
				upcrmmodel.UpRankTypeFans30dayIncreaseCount,
				upcrmmodel.UpRankTypeFans30dayIncreasePercent:
				info.CompleteTime = info.FirstUpTime + xtime.Time(info.Value)
			}
		}
	}

	// 排序
	for k, list := range typeMap {
		var sortFunc sortRankFunc = sortByValueAsc
		switch k {
		case upcrmmodel.UpRankTypeFans30dayIncreaseCount,
			upcrmmodel.UpRankTypeFans30dayIncreasePercent:
			sortFunc = sortByValueDesc
		}
		sortRankInfo(list, sortFunc)
		for i, v := range list {
			v.Rank = i + 1
		}
		log.Info("cache rank type=%d, len=%d, date=%v", k, len(list), date)
	}
	s.uprankCache = typeMap
	s.lastCacheDate = date
}

//UpRankQueryList query up rank list
func (s *Service) UpRankQueryList(c context.Context, arg *upcrmmodel.UpRankQueryArgs) (result upcrmmodel.UpRankQueryResult, err error) {
	if arg.Page == 0 {
		arg.Page = 1
	}

	if arg.Size <= 0 || arg.Size > 50 {
		arg.Size = 20
	}
	// 1.从内存中读数据
	var rankList, ok = s.uprankCache[arg.Type]
	if !ok || rankList == nil {
		log.Warn("no available rank data in cache, type=%d", arg.Type)
		return
	}

	var startIndex = (arg.Page - 1) * arg.Size
	var endIndex = startIndex + arg.Size
	var totalCount = len(rankList)
	endIndex = mathutil.Min(endIndex, totalCount)
	if startIndex < endIndex {
		result.Result = rankList[startIndex:endIndex]
	}
	result.Date = xtime.Time(s.lastCacheDate.Unix())
	result.TotalCount = totalCount
	result.Size = arg.Size
	result.Page = arg.Page
	log.Info("get rank list for type=%d", arg.Type)
	// 2.如果没有，从数据库中刷数据
	return
}
