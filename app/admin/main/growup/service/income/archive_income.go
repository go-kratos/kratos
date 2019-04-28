package income

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	model "go-common/app/admin/main/growup/model/income"
	"go-common/library/log"
	"go-common/library/xstr"
)

// ArchiveStatis archive income statis
func (s *Service) ArchiveStatis(c context.Context, categoryID []int64, typ, groupType int, fromTime, toTime int64) (data interface{}, err error) {
	table := setArchiveTableByGroup(typ, groupType)
	from := getDateByGroup(groupType, time.Unix(fromTime, 0))
	to := getDateByGroup(groupType, time.Unix(toTime, 0))
	query := formatArchiveQuery(categoryID, from, to)

	if typ == _lottery {
		data, err = s.lotteryStatis(c, categoryID, from, addDayByGroup(groupType, to).AddDate(0, 0, -1), groupType)
		if err != nil {
			log.Error("s.lotteryStatis error(%v)", err)
		}
		return
	}
	avs, err := s.GetArchiveStatis(c, table, query)
	if err != nil {
		log.Error("s.GetArchiveStatis error(%v)", err)
		return
	}

	data = archiveStatis(avs, from, to, groupType)
	return
}

func archiveStatis(avs []*model.ArchiveStatis, from, to time.Time, groupType int) interface{} {
	avsMap := make(map[string]*model.ArchiveStatis)
	ctgyMap := make(map[string]bool)
	for _, av := range avs {
		date := formatDateByGroup(av.CDate.Time(), groupType)
		ctgykey := date + strconv.FormatInt(av.CategroyID, 10)
		if val, ok := avsMap[date]; ok {
			val.Avs += av.Avs
			if !ctgyMap[ctgykey] {
				val.Income += av.Income
				ctgyMap[ctgykey] = true
			}
		} else {
			avsMap[date] = &model.ArchiveStatis{
				Avs:    av.Avs,
				Income: av.Income,
			}
			ctgyMap[ctgykey] = true
		}
	}
	return parseArchiveStatis(avsMap, from, to, groupType)
}

func parseArchiveStatis(avsMap map[string]*model.ArchiveStatis, from, to time.Time, groupType int) interface{} {
	income, counts, xAxis := []string{}, []int64{}, []string{}
	// get result by date
	to = to.AddDate(0, 0, 1)
	for from.Before(to) {
		dateStr := formatDateByGroup(from, groupType)
		xAxis = append(xAxis, dateStr)
		if val, ok := avsMap[dateStr]; ok {
			income = append(income, fmt.Sprintf("%.2f", float64(val.Income)/float64(100)))
			counts = append(counts, val.Avs)
		} else {
			income = append(income, "0")
			counts = append(counts, int64(0))
		}
		from = addDayByGroup(groupType, from)
	}

	return map[string]interface{}{
		"counts":  counts,
		"incomes": income,
		"xaxis":   xAxis,
	}
}

// ArchiveSection get av/column income section
func (s *Service) ArchiveSection(c context.Context, categoryID []int64, typ, groupType int, fromTime, toTime int64) (data interface{}, err error) {
	table := setArchiveTableByGroup(typ, groupType)
	from := getDateByGroup(groupType, time.Unix(fromTime, 0))
	to := getDateByGroup(groupType, time.Unix(toTime, 0))
	query := formatArchiveQuery(categoryID, from, to)

	avs, err := s.GetArchiveStatis(c, table, query)
	if err != nil {
		log.Error("s.GetArchiveStatis error(%v)", err)
		return
	}
	data = archiveSection(avs, from, to, groupType)
	return
}

func archiveSection(avs []*model.ArchiveStatis, from, to time.Time, groupType int) interface{} {
	ret := make([]map[string]interface{}, 0)
	avsMap := make(map[string][]int64)
	for _, av := range avs {
		date := formatDateByGroup(av.CDate.Time(), groupType)
		if val, ok := avsMap[date]; ok {
			val[av.MoneySection] += av.Avs
		} else {
			avsMap[date] = make([]int64, 12)
			avsMap[date][av.MoneySection] = av.Avs
			ret = append(ret, map[string]interface{}{
				"date_format": date,
				"sections":    avsMap[date],
			})
		}
	}
	return ret
}

// ArchiveDetail archive detail (av column)
func (s *Service) ArchiveDetail(c context.Context, mid int64, typ, groupType int, fromTime, toTime int64) (archives []*model.ArchiveIncome, err error) {
	archives = make([]*model.ArchiveIncome, 0)
	from := getDateByGroup(groupType, time.Unix(fromTime, 0))
	to := getDateByGroup(groupType, time.Unix(toTime, 0))
	to = addDayByGroup(groupType, to).AddDate(0, 0, -1)
	if typ == _video || typ == _up {
		var avs []*model.ArchiveIncome
		avs, err = s.archiveDetail(c, _video, groupType, mid, from, to)
		if err != nil {
			log.Error("s.archiveDetail error(%v)", err)
			return
		}
		archives = append(archives, avs...)
	}

	if typ == _column || typ == _up {
		var columns []*model.ArchiveIncome
		columns, err = s.archiveDetail(c, _column, groupType, mid, from, to)
		if err != nil {
			log.Error("s.archiveDetail error(%v)", err)
			return
		}
		archives = append(archives, columns...)
	}

	if typ == _bgm || typ == _up {
		var bgms []*model.ArchiveIncome
		bgms, err = s.archiveDetail(c, _bgm, groupType, mid, from, to)
		if err != nil {
			log.Error("s.archiveDetail error(%v)", err)
			return
		}
		archives = append(archives, bgms...)
	}
	return
}

func (s *Service) archiveDetail(c context.Context, typ, groupType int, mid int64, from, to time.Time) (archives []*model.ArchiveIncome, err error) {
	archives = make([]*model.ArchiveIncome, 0)
	query := fmt.Sprintf("mid = %d", mid)
	origins, err := s.GetArchiveIncome(c, typ, query, from.Format(_layout), to.Format(_layout))
	if err != nil {
		log.Error("s.GetArchiveIncome error(%v)", err)
		return
	}

	var black map[int64]struct{}
	black, err = s.dao.GetAvBlackListByMID(c, mid, typ)
	if err != nil {
		log.Error("s.dao.GetAvBlackListByMID error(%v)", err)
		return
	}
	archives = calArchiveDetail(origins, black, groupType)
	return
}

func calArchiveDetail(archives []*model.ArchiveIncome, blackMap map[int64]struct{}, groupType int) []*model.ArchiveIncome {
	avsMap := make(map[string]*model.ArchiveIncome)
	for _, av := range archives {
		if _, ok := blackMap[av.AvID]; ok {
			continue
		}
		date := formatDateByGroup(av.Date.Time(), groupType)
		key := date + strconv.FormatInt(av.AvID, 10)
		if val, ok := avsMap[key]; ok {
			val.Income += av.Income
		} else {
			av.DateFormat = date
			avsMap[key] = av
		}
	}
	list := make([]*model.ArchiveIncome, 0)
	for _, av := range avsMap {
		list = append(list, av)
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].DateFormat == list[j].DateFormat {
			return list[i].Income > list[j].Income
		}
		return list[i].DateFormat > list[j].DateFormat
	})
	return list
}

// ArchiveTop archive_income top
func (s *Service) ArchiveTop(c context.Context, aIDs []int64, typ int, groupType int, fromTime, toTime int64, from, limit int) (data []*model.ArchiveIncome, total int, err error) {
	query := ""
	if len(aIDs) != 0 {
		switch typ {
		case _video, _lottery:
			query = fmt.Sprintf("av_id IN (%s)", xstr.JoinInts(aIDs))
		case _column:
			query = fmt.Sprintf("aid IN (%s)", xstr.JoinInts(aIDs))
		case _bgm:
			query = fmt.Sprintf("sid IN (%s)", xstr.JoinInts(aIDs))
		}
	}
	if query == "" && typ != _lottery {
		query = fmt.Sprintf("income >= %d", _leastAvIncome)
	}
	avs, err := s.GetArchiveIncome(c, typ, query, time.Unix(fromTime, 0).Format(_layout), time.Unix(toTime, 0).Format(_layout))
	if err != nil {
		log.Error("s.GetArchiveIncome error(%v)", err)
		return
	}
	if typ == _lottery {
		typ = _video
	}
	avBMap, err := s.GetAvBlackListByAvIds(c, avs, typ)
	if err != nil {
		log.Error("s.GetAvBlackListByAvIds error(%v)", err)
		return
	}
	data, total = archiveTop(avs, avBMap, from, limit)

	upInfo, err := s.GetUpInfoByAIDs(c, data)
	if err != nil {
		log.Error("s.GetUpInfoByAIDs error(%v)", err)
		return
	}

	for i := 0; i < len(data); i++ {
		data[i].Nickname = upInfo[data[i].MID]
	}
	return
}

func archiveTop(avs []*model.ArchiveIncome, avBlack map[int64]struct{}, from, limit int) ([]*model.ArchiveIncome, int) {
	nAvs := make([]*model.ArchiveIncome, 0)
	for _, av := range avs {
		if _, ok := avBlack[av.AvID]; ok {
			continue
		}
		av.DateFormat = av.Date.Time().Format(_layout)
		nAvs = append(nAvs, av)
	}
	sort.Slice(nAvs, func(i, j int) bool {
		if nAvs[i].Date == nAvs[j].Date {
			return nAvs[i].Income > nAvs[j].Income
		}
		return nAvs[i].Date > nAvs[j].Date
	})

	if limit+from > len(nAvs) {
		limit = len(nAvs)
	}
	total := len(nAvs)
	return nAvs[from:limit], total
}

// BgmDetail bgm detail
func (s *Service) BgmDetail(c context.Context, sid int64, fromTime, toTime int64, from, limit int) (avs []*model.ArchiveIncome, total int, err error) {
	avs = make([]*model.ArchiveIncome, 0)
	fromDate := time.Unix(fromTime, 0).Format(_layout)
	toDate := time.Unix(toTime, 0).Format(_layout)
	avMap, err := s.dao.GetAvByBgm(c, sid, fromDate, toDate)
	if err != nil {
		log.Error("s.dao.GetAvByBgm error(%v)", err)
		return
	}
	if len(avMap) == 0 {
		return
	}
	avIDs := make([]int64, 0, len(avMap))
	for avID := range avMap {
		avIDs = append(avIDs, avID)
	}
	avs, err = s.GetArchiveIncome(c, _video, fmt.Sprintf("av_id in (%s)", xstr.JoinInts(avIDs)), fromDate, toDate)
	if err != nil {
		log.Error("s.GetArchiveIncome error(%v)", err)
		return
	}
	if limit > len(avs) {
		limit = len(avs)
	}
	total = len(avs)
	avs = avs[from:limit]
	return
}

// GetArchiveStatis get up income
func (s *Service) GetArchiveStatis(c context.Context, table, query string) (avs []*model.ArchiveStatis, err error) {
	offset, size := 0, 2000
	for {
		av, err := s.dao.GetArchiveStatis(c, table, query, offset, size)
		if err != nil {
			return nil, err
		}
		avs = append(avs, av...)
		if len(av) < size {
			break
		}
		offset += len(av)
	}
	return
}

// GetArchiveIncome get archive income
func (s *Service) GetArchiveIncome(c context.Context, typ int, query string, from, to string) (archs []*model.ArchiveIncome, err error) {
	var id int64
	limit := 2000
	for {
		var arch []*model.ArchiveIncome
		arch, err = s.dao.GetArchiveIncome(c, id, query, from, to, limit, typ)
		if err != nil {
			return
		}
		archs = append(archs, arch...)
		if len(arch) < limit {
			break
		}
		id = arch[len(arch)-1].ID
	}
	if typ == _bgm {
		bgms := make(map[string]*model.ArchiveIncome)
		for _, bgm := range archs {
			key := bgm.Date.Time().Format(_layout) + strconv.FormatInt(bgm.AvID, 10)
			if b, ok := bgms[key]; !ok {
				bgms[key] = bgm
			} else {
				b.Income += bgm.Income
				b.TotalIncome += bgm.TotalIncome
				b.TaxMoney += bgm.TaxMoney
			}
			bgms[key].Avs++
		}
		archs = make([]*model.ArchiveIncome, 0, len(bgms))
		for _, b := range bgms {
			archs = append(archs, b)
		}
	}
	return
}

func formatDateByGroup(date time.Time, groupType int) string {
	str := ""
	if groupType == _groupWeek {
		date = getStartWeekDate(date)
		str = date.Format(_layout) + "~" + date.AddDate(0, 0, 6).Format(_layout)
	} else if groupType == _groupMonth {
		date = getStartMonthDate(date)
		str = date.Format(_layoutMonth)
	} else {
		str = date.Format(_layout)
	}
	return str
}

func formatArchiveQuery(categoryID []int64, from, to time.Time) string {
	query := "cdate >= '" + from.Format(_layout) + "'"
	query += " AND "
	query += "cdate <= '" + to.Format(_layout) + "'"
	if len(categoryID) != 0 {
		query += " AND "
		query += "category_id in (" + xstr.JoinInts(categoryID) + ")"
	}
	return query
}
