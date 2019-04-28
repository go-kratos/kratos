package income

import (
	"context"
	"fmt"
	"sort"
	"time"

	model "go-common/app/admin/main/growup/model/income"
	"go-common/app/admin/main/growup/service"

	"go-common/library/log"
	"go-common/library/xstr"
)

// UpIncomeList up income list
func (s *Service) UpIncomeList(c context.Context, mids []int64, typ, groupType int, fromTime, toTime, minIncome, maxIncome int64, from, limit int) (upIncome []*model.UpIncome, total int, err error) {
	table := setUpTableByGroup(groupType)
	_, incomeType := getUpInfoByType(typ)
	fromDate := getDateByGroup(groupType, time.Unix(fromTime, 0))
	toDate := getDateByGroup(groupType, time.Unix(toTime, 0))
	typeField := getUpFieldByType(typ)
	query := formatUpQuery(mids, fromDate, toDate, incomeType)
	if maxIncome != 0 || minIncome != 0 {
		query = fmt.Sprintf("%s AND %s >= %d AND %s <= %d", query, incomeType, minIncome, incomeType, maxIncome)
	}
	total, err = s.dao.UpIncomeCount(c, table, query)
	if err != nil {
		log.Error("s.dao.UpIncomeCount error(%v)", err)
		return
	}

	upIncome, err = s.dao.GetUpIncomeBySort(c, table, typeField, incomeType, query, from, limit)
	if err != nil {
		log.Error("s.dao.GetUpIncomeBySort error(%v)", err)
		return
	}
	if len(upIncome) == 0 {
		return
	}
	for _, up := range upIncome {
		mids = append(mids, up.MID)
	}

	nicknames, err := s.dao.ListUpInfo(c, mids)
	if err != nil {
		log.Error("s.dao.ListUpInfo error(%v)", err)
		return
	}

	var breachType []int64
	switch typ {
	case _video:
		breachType = []int64{0}
	case _column:
		breachType = []int64{2}
	case _bgm:
		breachType = []int64{3}
	case _up:
		breachType = []int64{0, 1, 2, 3}
	}
	breachs, err := s.dao.GetAvBreachByMIDs(c, mids, breachType)
	if err != nil {
		log.Error("s.dao.GetAvBreachByMIDs error(%v)", err)
		return
	}
	upIncomeList(upIncome, nicknames, breachs, typ, groupType)
	return
}

func upIncomeList(upIncome []*model.UpIncome, nicknames map[int64]string, breachs []*model.AvBreach, typ, groupType int) {
	midDateBreach := make(map[int64]map[string]int64) // map[mid][date] = money
	for _, b := range breachs {
		dateFormat := formatDateByGroup(b.CDate.Time(), groupType)
		if _, ok := midDateBreach[b.MID]; !ok {
			midDateBreach[b.MID] = make(map[string]int64)
		}
		midDateBreach[b.MID][dateFormat] += b.Money
	}

	for _, up := range upIncome {
		up.Nickname = nicknames[up.MID]
		up.DateFormat = formatDateByGroup(up.Date.Time(), groupType)
		up.ExtraIncome = up.Income - up.BaseIncome
		if _, ok := midDateBreach[up.MID]; ok {
			up.Breach = midDateBreach[up.MID][up.DateFormat]
		}
		switch typ {
		case _video:
			up.Count = up.AvCount
		case _column:
			up.Count = up.ColumnCount
		case _bgm:
			up.Count = up.BgmCount
		case _up:
			up.Count = up.AvCount + up.ColumnCount + up.BgmCount
		}
	}
}

// UpIncomeListExport up income list
func (s *Service) UpIncomeListExport(c context.Context, mids []int64, typ, groupType int, fromTime, toTime, minIncome, maxIncome int64, from, limit int) (res []byte, err error) {
	ups, _, err := s.UpIncomeList(c, mids, typ, groupType, fromTime, toTime, minIncome, maxIncome, from, limit)
	if err != nil {
		log.Error("s.UpIncomeList error(%v)", err)
		return
	}
	records := formatUpIncome(ups)
	res, err = service.FormatCSV(records)
	if err != nil {
		log.Error("FormatCSV error(%v)")
	}
	return
}

// UpIncomeStatis up income statis
func (s *Service) UpIncomeStatis(c context.Context, mids []int64, typ, groupType int, fromTime, toTime int64) (data interface{}, err error) {
	from := getDateByGroup(groupType, time.Unix(fromTime, 0))
	to := getDateByGroup(groupType, time.Unix(toTime, 0))
	if groupType == _groupDay && len(mids) == 0 {
		return s.upStatisDaily(c, typ, from, to)
	}
	return s.upIncomeStatisDate(c, mids, typ, groupType, from, to)
}

func (s *Service) upStatisDaily(c context.Context, typ int, from, to time.Time) (data interface{}, err error) {
	table, _ := getUpInfoByType(typ)
	statis, err := s.dao.GetUpDailyStatis(c, table, from.Format(_layout), to.Format(_layout))
	if err != nil {
		log.Error("s.dao.GetIncomeDailyStatis error(%v)", err)
		return
	}
	dateMap := make(map[string]*model.UpStatisRsp)
	for _, sta := range statis {
		date := sta.Date.Time().Format(_layout)
		if val, ok := dateMap[date]; ok {
			val.Ups += sta.Ups
		} else {
			dateMap[date] = &model.UpStatisRsp{
				Income: sta.Income,
				Ups:    sta.Ups,
			}
		}
	}
	data = calUpStatis(dateMap, 1, from, to)
	return
}

func (s *Service) upIncomeStatisDate(c context.Context, mids []int64, typ, groupType int, from, to time.Time) (data interface{}, err error) {
	table := setUpTableByGroup(groupType)
	_, incomeType := getUpInfoByType(typ)
	query := formatUpQuery(mids, from, to, incomeType)
	upIncome, err := s.GetUpIncome(c, table, incomeType, query)
	if err != nil {
		log.Error("s.GetUpIncome error(%v)", err)
		return
	}

	sort.Slice(upIncome, func(i, j int) bool {
		return upIncome[i].Date < upIncome[j].Date
	})

	dateMap := make(map[string]*model.UpStatisRsp)
	for _, up := range upIncome {
		date := formatDateByGroup(up.Date.Time(), groupType)
		if val, ok := dateMap[date]; ok {
			val.Income += up.Income
			val.Ups++
		} else {
			dateMap[date] = &model.UpStatisRsp{
				Income: up.Income,
				Ups:    1,
			}
		}
	}
	data = calUpStatis(dateMap, groupType, from, to)
	return
}

func calUpStatis(dateMap map[string]*model.UpStatisRsp, groupType int, from, to time.Time) interface{} {
	incomes, ups, xAxis := []string{}, []int{}, []string{}
	to = to.AddDate(0, 0, 1)
	for from.Before(to) {
		dateStr := formatDateByGroup(from, groupType)
		xAxis = append(xAxis, dateStr)
		if val, ok := dateMap[dateStr]; ok {
			incomes = append(incomes, fmt.Sprintf("%.2f", float64(val.Income)/float64(100)))
			ups = append(ups, val.Ups)
		} else {
			incomes = append(incomes, "0")
			ups = append(ups, 0)
		}

		from = addDayByGroup(groupType, from)
	}

	return map[string]interface{}{
		"counts":  ups,
		"incomes": incomes,
		"xaxis":   xAxis,
	}
}

// GetUpIncome get
func (s *Service) GetUpIncome(c context.Context, table, incomeType, query string) (upIncomes []*model.UpIncome, err error) {
	var id int64
	limit := 2000
	for {
		upIncome, err := s.dao.GetUpIncome(c, table, incomeType, query, id, limit)
		if err != nil {
			return upIncomes, err
		}
		upIncomes = append(upIncomes, upIncome...)
		if len(upIncome) < limit {
			break
		}
		id = upIncome[len(upIncome)-1].ID
	}
	return
}

func formatUpQuery(mids []int64, fromTime, toTime time.Time, incomeType string) string {
	query := "date >= '" + fromTime.Format("2006-01-02") + "'"
	query += " AND "
	query += "date <= '" + toTime.Format("2006-01-02") + "'"
	if len(mids) != 0 {
		query += " AND "
		query += "mid in (" + xstr.JoinInts(mids) + ")"
	}
	query += fmt.Sprintf(" AND %s > 0", incomeType)
	return query
}
