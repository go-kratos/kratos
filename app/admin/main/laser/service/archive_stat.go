package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/admin/main/laser/model"
)

const (
	// ALLNAME is specific name video_audit overview
	ALLNAME = "ALL全体总览"
	// ALLUID is specific uid for video_audit overview
	ALLUID = -1
)

// ArchiveRecheck is stat recheck flow data node.
func (s *Service) ArchiveRecheck(c context.Context, typeIDS []int64, unames string, startDate int64, endDate int64) (recheckViews []*model.StatView, err error) {
	start := time.Unix(startDate, 0)
	end := time.Unix(endDate, 0)
	var uids []int64
	var res map[string]int64
	if len(unames) != 0 {
		res, err = s.dao.GetUIDByNames(c, unames)
		if err != nil {
			return
		}
		for _, uid := range res {
			uids = append(uids, uid)
		}
	}
	statTypes := []int64{model.TotalArchive, model.TotalOper, model.ReCheck, model.Lock,
		model.ThreeLimit, model.FirstCheck, model.SecondCheck, model.ThirdCheck, model.NoRankArchive, model.NoIndexArchive, model.NoRecommendArchive, model.NoPushArchive,
		model.FirstCheckOper, model.FirstCheckTime, model.SecondCheckOper,
		model.SecondCheckTime, model.ThirdCheckOper, model.ThirdCheckTime}

	var statViews []*model.StatView
	for !start.After(end) {
		statViews, err = s.dailyStatArchiveRecheck(c, model.ArchiveRecheck, typeIDS, statTypes, uids, start)
		if err != nil {
			return
		}
		recheckViews = append(recheckViews, statViews...)
		start = start.AddDate(0, 0, 1)
	}
	return
}

// UserRecheck is stat user recheck data node.
func (s *Service) UserRecheck(c context.Context, typeIDS []int64, unames string, startDate int64, endDate int64) (
	recheckViews []*model.StatView, err error) {
	start := time.Unix(startDate, 0)
	end := time.Unix(endDate, 0)
	var uids []int64
	var res map[string]int64
	if len(unames) != 0 {
		res, err = s.dao.GetUIDByNames(c, unames)
		if err != nil {
			return
		}
		for _, uid := range res {
			uids = append(uids, uid)
		}
	}
	statTypes := []int64{model.TotalOperFrequency, model.FirstCheckOper, model.SecondCheckOper, model.ThirdCheckOper,
		model.FirstCheckOper, model.FirstCheckTime, model.SecondCheckOper,
		model.SecondCheckTime, model.ThirdCheckOper, model.ThirdCheckTime}
	var statViews []*model.StatView
	for !start.After(end) {
		statViews, err = s.dailyStatArchiveRecheck(c, model.ArchiveRecheck, typeIDS, statTypes, uids, start)
		if err != nil {
			return
		}
		recheckViews = append(recheckViews, statViews...)
		start = start.AddDate(0, 0, 1)
	}
	return
}

func (s *Service) dailyStatArchiveRecheck(c context.Context, business int, typeIDS []int64, statTypes []int64, uids []int64, statDate time.Time) (statViews []*model.StatView, err error) {
	mediateView, err := s.dailyArchiveStat(c, business, typeIDS, statTypes, uids, statDate)
	if err != nil || len(mediateView) == 0 {
		return
	}
	statViews = makeUpArchiveRecheck(mediateView)
	return
}

func makeUpArchiveRecheck(mediateView map[int64]map[int]int64) (statViews []*model.StatView) {
	items := make(map[int64][]*model.StatItem)
	for k1, v1 := range mediateView {
		var recheckItems []*model.StatItem
		denominatorValue, ok1 := v1[model.FirstCheckOper]
		numeratorValue, ok2 := v1[model.FirstCheckTime]
		if ok1 && ok2 {
			if denominatorValue == 0 {
				recheckItems = append(recheckItems, &model.StatItem{
					DataCode: model.FirstAvgTime,
					Value:    0,
				})
			} else {
				recheckItems = append(recheckItems, &model.StatItem{
					DataCode: model.FirstAvgTime,
					Value:    numeratorValue / denominatorValue,
				})
			}
		}
		denominatorValue, ok1 = v1[model.SecondCheckOper]
		numeratorValue, ok2 = v1[model.SecondCheckTime]
		if ok1 && ok2 {
			if denominatorValue == 0 {
				recheckItems = append(recheckItems, &model.StatItem{
					DataCode: model.SecondAvgTime,
					Value:    0,
				})
			} else {
				recheckItems = append(recheckItems, &model.StatItem{
					DataCode: model.SecondAvgTime,
					Value:    numeratorValue / denominatorValue,
				})
			}
		}
		denominatorValue, ok1 = v1[model.ThirdCheckOper]
		numeratorValue, ok2 = v1[model.ThirdCheckTime]
		if ok1 && ok2 {
			if denominatorValue == 0 {
				recheckItems = append(recheckItems, &model.StatItem{
					DataCode: model.ThirdAvgTime,
					Value:    0,
				})
			} else {
				recheckItems = append(recheckItems, &model.StatItem{
					DataCode: model.ThirdAvgTime,
					Value:    numeratorValue / denominatorValue,
				})
			}
		}

		for k2, v2 := range v1 {
			recheckItems = append(recheckItems, &model.StatItem{DataCode: k2, Value: v2})
		}
		items[k1] = recheckItems
	}

	for k, v := range items {
		statViews = append(statViews, &model.StatView{Date: k, Stats: v})
	}
	return
}

func (s *Service) dailyArchiveStat(c context.Context, business int, typeIDS []int64, statTypes []int64, uids []int64, statDate time.Time) (mediateView map[int64]map[int]int64, err error) {
	statNodes, err := s.dao.StatArchiveStat(c, business, typeIDS, uids, statTypes, statDate)
	if err != nil || len(statNodes) == 0 {
		return
	}
	mediateView = make(map[int64]map[int]int64)
	for _, node := range statNodes {
		k1 := node.StatDate.Time().Unix()
		k2 := node.StatType
		newValue := node.StatValue
		if v1, ok := mediateView[k1]; ok {
			if v2, ok := v1[k2]; ok {
				mediateView[k1][k2] = v2 + newValue
			} else {
				mediateView[k1][k2] = newValue
			}
		} else {
			mediateView[k1] = map[int]int64{k2: newValue}
		}
	}
	return
}

// TagRecheck is stat archive tag recheck.
func (s *Service) TagRecheck(c context.Context, startDate int64, endDate int64, unames string) (tagViews []*model.StatView, err error) {
	start := time.Unix(startDate, 0)
	end := time.Unix(endDate, 0)
	var uids []int64
	var uname2uid map[string]int64
	if len(unames) != 0 {
		uname2uid, err = s.dao.GetUIDByNames(c, unames)
		if err != nil {
			return
		}
		for _, uid := range uname2uid {
			uids = append(uids, uid)
		}
	}
	statTypes := []int64{model.TagRecheckTotalTime, model.TagRecheckTotalCount, model.TagChangeCount, model.TagRecheckTotalCount, model.TagRecheckTotalTime}
	var statViews []*model.StatView
	for !start.After(end) {
		statViews, err = s.dailyStatTagRecheck(c, model.TagRecheck, statTypes, uids, start)
		if err != nil {
			return
		}
		tagViews = append(tagViews, statViews...)
		start = start.AddDate(0, 0, 1)
	}
	return
}

func (s *Service) dailyStatTagRecheck(c context.Context, business int, statTypes []int64, uids []int64, statDate time.Time) (statViews []*model.StatView, err error) {
	mediateView, err := s.dailyArchiveStat(c, business, []int64{}, statTypes, uids, statDate)
	if err != nil || len(mediateView) == 0 {
		return
	}
	statViews = makeUpTagRecheck(mediateView)
	return
}

func makeUpTagRecheck(mediateView map[int64]map[int]int64) (statViews []*model.StatView) {
	items := make(map[int64][]*model.StatItem)
	for k1, v1 := range mediateView {
		var recheckItems []*model.StatItem
		denominatorValue, ok1 := v1[model.TagRecheckTotalCount]
		numeratorValue, ok2 := v1[model.TagRecheckTotalTime]
		if ok1 && ok2 {
			if denominatorValue == 0 {
				recheckItems = append(recheckItems, &model.StatItem{
					DataCode: model.TagRecheckAvgTime,
					Value:    0,
				})
			} else {
				recheckItems = append(recheckItems, &model.StatItem{
					DataCode: model.TagRecheckAvgTime,
					Value:    numeratorValue / denominatorValue,
				})
			}
		}

		for k2, v2 := range v1 {
			recheckItems = append(recheckItems, &model.StatItem{DataCode: k2, Value: v2})
		}
		items[k1] = recheckItems
	}

	for k, v := range items {
		statViews = append(statViews, &model.StatView{Date: k, Stats: v})
	}
	return
}

// Recheck123 is stat 123 recheck.
func (s *Service) Recheck123(c context.Context, startDate int64, endDate int64, typeIDS []int64) (recheckView []*model.StatView, err error) {
	start := time.Unix(startDate, 0)
	end := time.Unix(endDate, 0)
	emptyUids := []int64{}
	statTypes := []int64{model.FIRST_RECHECK_IN, model.FIRST_RECHECK_OUT, model.SECOND_RECHECK_IN, model.SECOND_RECHECK_OUT, model.THIRD_RECHECK_IN, model.THIRD_RECHECK_OUT}
	var statViews []*model.StatView
	for !start.After(end) {
		statViews, err = s.dailyStatArchiveStreamStat(c, model.Recheck123, typeIDS, emptyUids, statTypes, start)
		if err != nil {
			return
		}
		recheckView = append(recheckView, statViews...)
		start = start.AddDate(0, 0, 1)
	}
	return
}

func (s *Service) dailyStatArchiveStreamStat(c context.Context, business int, typeIDS []int64, uids []int64, statTypes []int64, statDate time.Time) (statViews []*model.StatView, err error) {
	statNodes, err := s.dao.StatArchiveStatStream(c, model.Recheck123, typeIDS, uids, statTypes, statDate)
	if err != nil || len(statNodes) == 0 {
		return
	}
	mediateView := make(map[int64]map[int]int64)
	for _, v := range statNodes {
		k1 := v.StatDate.Time().Unix()
		k2 := v.StatType
		newValue := v.StatValue
		if v1, ok := mediateView[k1]; ok {
			if v2, ok := v1[k2]; ok {
				mediateView[k1][k2] = newValue + v2
			} else {
				mediateView[k1][k2] = newValue
			}
		} else {
			mediateView[k1] = map[int]int64{k2: newValue}
		}
	}
	for k1, v1 := range mediateView {
		var statItems []*model.StatItem
		for k2, v2 := range v1 {
			statItems = append(statItems, &model.StatItem{DataCode: k2, Value: v2})
		}
		statViews = append(statViews, &model.StatView{Date: k1, Stats: statItems})
	}
	return
}

func wrap(cargoMap map[int64]*model.CargoItem) (views []*model.CargoView) {
	// 适配返回的JSON结构,减少前端工作量, cargoMap 2 views.
	mediateView := make(map[string]map[int]*model.CargoItem)
	for k, v := range cargoMap {
		statDate := time.Unix(k, 0)
		k1 := statDate.Format("2006-01-02")
		k2 := statDate.Hour()
		if value1, ok := mediateView[k1]; ok {
			if value2, ok := value1[k2]; ok {
				value2.AuditValue = value2.AuditValue + v.AuditValue
				value2.ReceiveValue = value2.ReceiveValue + v.ReceiveValue
				mediateView[k1][k2] = value2
			} else {
				mediateView[k1][k2] = &model.CargoItem{
					ReceiveValue: v.ReceiveValue,
					AuditValue:   v.AuditValue,
				}
			}
		} else {
			mediateView[k1] = map[int]*model.CargoItem{
				k2: {
					ReceiveValue: v.ReceiveValue,
					AuditValue:   v.AuditValue,
				},
			}
		}
	}

	for k, v := range mediateView {
		views = append(views, &model.CargoView{
			Date: k,
			Data: v,
		})
	}
	return
}

// CsvAuditCargo is download archive cargo audit data by csv file type.
func (s *Service) CsvAuditCargo(c context.Context, startDate int64, endDate int64, unames string) (res []byte, err error) {
	wrappers, lineWidth, err := s.AuditorCargoList(c, startDate, endDate, unames)
	if err != nil {
		return
	}
	data := formatAuditCargo(wrappers, lineWidth)
	return FormatCSV(data)
}

// AuditorCargoList is query archive audit cargo by uname respectively with stat_date condition.
func (s *Service) AuditorCargoList(c context.Context, startDate int64, endDate int64, unames string) (wrappers []*model.CargoViewWrapper, lineWidth int, err error) {
	start := time.Unix(startDate, 0)
	end := time.Unix(endDate, 0)
	var uids []int64
	var name2uid map[string]int64
	uid2name := make(map[int64]string)
	if len(unames) != 0 {
		name2uid, err = s.dao.GetUIDByNames(c, unames)
		if err != nil {
			return
		}
		for name, uid := range name2uid {
			uids = append(uids, uid)
			uid2name[uid] = name
		}
	}

	var items, itemsBlock []*model.CargoDetail
	for !start.After(end) {
		itemsBlock, err = s.dao.QueryArchiveCargo(c, start, uids)
		if err != nil {
			return
		}
		items = append(items, itemsBlock...)
		start = start.Add(time.Hour * 1)
	}
	if len(items) == 0 {
		return
	}

	mediateViews := make(map[int64]map[int64]*model.CargoItem)
	uidMap := make(map[int64]bool)
	for _, v := range items {
		k1 := v.UID
		k2 := v.StatDate.Time().Unix()
		uidMap[k1] = true

		if v1, ok := mediateViews[k1]; ok {
			if v2, ok := v1[k2]; ok {
				v2.ReceiveValue = v2.ReceiveValue + v.ReceiveValue
				v2.AuditValue = v2.AuditValue + v.AuditValue
				mediateViews[k1][k2] = v2
			} else {
				lineWidth = lineWidth + 1
				mediateViews[k1][k2] = &model.CargoItem{
					ReceiveValue: v.ReceiveValue,
					AuditValue:   v.AuditValue,
				}
			}
		} else {
			lineWidth = lineWidth + 1
			mediateViews[k1] = map[int64]*model.CargoItem{
				k2: {
					ReceiveValue: v.ReceiveValue,
					AuditValue:   v.AuditValue,
				},
			}
		}
	}

	if len(unames) == 0 {
		for uid := range uidMap {
			uids = append(uids, uid)
		}
		uid2name, err = s.dao.GetUNamesByUids(c, uids)
		if err != nil {
			return
		}
	}

	for k, v := range mediateViews {
		cargoViews := wrap(v)
		for _, v := range cargoViews {
			wrappers = append(wrappers, &model.CargoViewWrapper{
				Username:  uid2name[k],
				CargoView: v,
			})
		}
	}
	return
}

// CsvRandomVideoAudit is download random video audit statistic data by csv file type.
func (s *Service) CsvRandomVideoAudit(c context.Context, startDate int64, endDate int64, unames string, typeIDS []int64) (res []byte, err error) {
	statViewExts, lineWidth, err := s.RandomVideo(c, startDate, endDate, typeIDS, unames)
	if err != nil {
		return
	}
	sort.Slice(statViewExts, func(i, j int) bool {
		return statViewExts[i].Date > statViewExts[j].Date
	})
	data := formatVideoAuditStat(statViewExts, lineWidth)
	return FormatCSV(data)
}

// CsvFixedVideoAudit is download fixed video audit statistic data by csv file type.
func (s *Service) CsvFixedVideoAudit(c context.Context, startDate int64, endDate int64, unames string, typeIDS []int64) (res []byte, err error) {
	statViewExts, lineWidth, err := s.FixedVideo(c, startDate, endDate, typeIDS, unames)
	if err != nil {
		return
	}
	sort.Slice(statViewExts, func(i, j int) bool {
		return statViewExts[i].Date > statViewExts[j].Date
	})
	data := formatVideoAuditStat(statViewExts, lineWidth)
	return FormatCSV(data)
}

// RandomVideo is stat random video type.
func (s *Service) RandomVideo(c context.Context, startDate int64, endDate int64, typeIDS []int64, uname string) (statViewExts []*model.StatViewExt, lineWidth int, err error) {
	start := time.Unix(startDate, 0)
	end := time.Unix(endDate, 0)
	var viewExts []*model.StatViewExt
	var width int
	for !start.After(end) {
		viewExts, width, err = s.videoAudit(c, model.RandomVideoAudit, start, typeIDS, uname)
		if err != nil {
			return
		}
		lineWidth = lineWidth + width
		statViewExts = append(statViewExts, viewExts...)
		start = start.AddDate(0, 0, 1)
	}
	return
}

// FixedVideo is stat fixed video type.
func (s *Service) FixedVideo(c context.Context, startDate int64, endDate int64, typeIDS []int64, uname string) (statViewExts []*model.StatViewExt, lineWidth int, err error) {
	start := time.Unix(startDate, 0)
	end := time.Unix(endDate, 0)
	var viewExts []*model.StatViewExt
	var width int
	for !start.After(end) {
		viewExts, width, err = s.videoAudit(c, model.FixedVideoAudit, start, typeIDS, uname)
		if err != nil {
			return
		}
		lineWidth = lineWidth + width
		statViewExts = append(statViewExts, viewExts...)
		start = start.AddDate(0, 0, 1)
	}
	return
}

func (s *Service) videoAudit(c context.Context, business int, statDate time.Time, typeIDS []int64, unames string) (viewExts []*model.StatViewExt, lineWidth int, err error) {
	var uids []int64
	var res map[string]int64
	needAll := true
	if len(unames) != 0 {
		needAll = false
		res, err = s.dao.GetUIDByNames(c, unames)
		if err != nil {
			return
		}
		for _, uid := range res {
			uids = append(uids, uid)
		}
	}
	statNodes, err := s.dao.StatArchiveStat(c, business, typeIDS, uids, []int64{}, statDate)
	if err != nil || len(statNodes) == 0 {
		return
	}
	return s.statNode2ViewExt(c, statNodes, needAll)
}

func (s *Service) statNode2ViewExt(c context.Context, statNodes []*model.StatNode, needAll bool) (statViewsExts []*model.StatViewExt, lineWidth int, err error) {
	mediateViews := make(map[int64]map[int64]map[int]int64)
	uidMap := make(map[int64]bool)
	var uids []int64
	for _, v := range statNodes {
		k1 := v.StatDate.Time().Unix()
		k2 := v.UID
		k3 := v.StatType
		newValue := v.StatValue

		uidMap[k2] = true
		if v1, ok := mediateViews[k1]; ok {
			if needAll {
				if allV2, ok := v1[ALLUID]; ok {
					if allV3, ok := allV2[k3]; ok {
						mediateViews[k1][ALLUID][k3] = allV3 + newValue
					} else {
						mediateViews[k1][ALLUID][k3] = newValue
					}
				} else {
					lineWidth = lineWidth + 1
					mediateViews[k1][ALLUID] = map[int]int64{k3: newValue}
				}
			}

			if v2, ok := v1[k2]; ok {
				if v3, ok := v2[k3]; ok {
					mediateViews[k1][k2][k3] = v3 + newValue
				} else {
					mediateViews[k1][k2][k3] = newValue
				}
			} else {
				lineWidth = lineWidth + 1
				mediateViews[k1][k2] = map[int]int64{k3: newValue}
			}
		} else {
			lineWidth = lineWidth + 1
			mediateViews[k1] = map[int64]map[int]int64{k2: {k3: newValue}}
			if needAll {
				lineWidth = lineWidth + 1
				mediateViews[k1] = map[int64]map[int]int64{ALLUID: {k3: newValue}}
			}

		}
	}

	//fetch uid map uname
	for uid := range uidMap {
		uids = append(uids, uid)
	}
	uid2name, err := s.dao.GetUNamesByUids(c, uids)
	if err != nil {
		return
	}
	if needAll {
		uid2name[ALLUID] = ALLNAME
	}

	for k1, v1 := range mediateViews {
		for k2, v2 := range v1 {
			var numeratorValue int64
			var denominatorValue int64
			for k3, v3 := range v2 {
				if k3 == model.WaitAuditDuration {
					numeratorValue = v3
				}
				if k3 == model.WaitAuditOper {
					denominatorValue = v3
				}
			}
			if denominatorValue == 0 {
				mediateViews[k1][k2][model.WaitAuditAvgTime] = 0
			} else {
				mediateViews[k1][k2][model.WaitAuditAvgTime] = numeratorValue / denominatorValue
			}
		}
	}

	for k1, v1 := range mediateViews {
		var wraps []*model.StatItemExt
		for k2, v2 := range v1 {
			var statItems []*model.StatItem
			for k3, v3 := range v2 {
				statItems = append(statItems, &model.StatItem{
					DataCode: k3,
					Value:    v3,
				})
			}
			wraps = append(wraps, &model.StatItemExt{
				Uname: uid2name[k2],
				Stats: statItems,
			})
		}
		// uname排序
		sort.Slice(wraps, func(i, j int) bool {
			if wraps[i].Uname == ALLNAME {
				return true
			}
			if wraps[j].Uname == ALLNAME {
				return false
			}
			return wraps[i].Uname < wraps[j].Uname
		})

		statViewsExts = append(statViewsExts, &model.StatViewExt{
			Date:  k1,
			Wraps: wraps,
		})
	}
	return
}
