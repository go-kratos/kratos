package upcrmservice

import (
	"context"
	"time"

	"go-common/app/admin/main/up/dao/upcrm"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/library/log"
)

//YesterdayRes struct
type YesterdayRes struct {
	ActivityUps int64 `json:"activity_ups"`
	IncrUps     int64 `json:"incr_ups"`
	TotalUps    int64 `json:"total_ups"`
	Date        int64 `json:"date"` // 截止日期
}

//TrendRes struct
type TrendRes struct {
	Date int64 `json:"date"` // 截止日期
	Ups  int64 `json:"ups"`
}

//TrendDetail  struct
type TrendDetail struct {
	ActivityUps int64 `json:"activity_ups"`
	IncrUps     int64 `json:"incr_ups"`
	TotalUps    int64 `json:"total_ups"`
	Date        int64 `json:"date"`
}

//QueryYesterday query yesterday data
func (s *Service) QueryYesterday(ctx context.Context, date time.Time) (item *YesterdayRes, err error) {
	var lastday, e = s.crmdb.GetUpStatLastDate(date)
	if e != nil {
		err = e
		log.Error("err get up stat last day, err=%+v", err)
		return
	}
	histories, err := s.crmdb.QueryYesterday(lastday)
	if err != nil {
		return nil, err
	}
	item = &YesterdayRes{}
	item.Date = lastday.Unix()
	for _, v := range histories {
		if v.Type == upcrmmodel.ActivityType {
			item.ActivityUps += v.Value1
		}
		if v.Type == upcrmmodel.IncrType {
			item.IncrUps += v.Value1
		}
		if v.Type == upcrmmodel.TotalType {
			item.TotalUps += v.Value1
		}
	}
	return
}

//QueryTrend  query trend
func (s *Service) QueryTrend(ctx context.Context, statType int, currentDate time.Time, days int) (result []*TrendRes, err error) {
	histories, err := s.crmdb.QueryTrend(statType, currentDate, days)
	if err != nil {
		return nil, err
	}
	var items = make(map[string]*TrendRes)
	var length = len(histories)
	var endDate = currentDate
	if length > 0 {
		endDate = histories[0].GenerateDate.Time()
	}

	for _, v := range histories {
		key := v.GenerateDate.Time().Format(upcrm.ISO8601DATE)
		if old, ok := items[key]; ok {
			new := &TrendRes{
				Date: old.Date,
				Ups:  old.Ups + v.Value1,
			}
			items[key] = new
		} else {
			items[key] = &TrendRes{
				Date: v.GenerateDate.Time().Unix(),
				Ups:  v.Value1,
			}
		}
	}

	// 这里的范围是[endDate - 6, endDate]，包含endDate
	for date := endDate.AddDate(0, 0, 1-days); !date.After(endDate); date = date.AddDate(0, 0, 1) {
		var dateStr = date.Format(upcrmmodel.TimeFmtDate)
		var d *TrendRes
		var ok bool
		if d, ok = items[dateStr]; !ok {
			d = &TrendRes{
				Date: date.Unix(),
				Ups:  0,
			}
		}
		result = append(result, d)
	}
	return
}

//QueryDetail query detail info
func (s *Service) QueryDetail(ctx context.Context, currentDate time.Time, days int) (result []*TrendDetail, err error) {
	var endDate, e = s.crmdb.GetUpStatLastDate(currentDate)
	if e != nil {
		err = e
		log.Error("err get up stat last day, err=%+v", err)
		return
	}
	var startDate = endDate.AddDate(0, 0, 1-days)
	histories, err := s.crmdb.QueryDetail(startDate, endDate)
	if err != nil {
		return nil, err
	}
	var items = make(map[string]*TrendDetail)
	for _, v := range histories {
		key := v.GenerateDate.Time().Format(upcrm.ISO8601DATE)
		if old, ok := items[key]; ok {
			if v.Type == upcrmmodel.ActivityType {
				old.ActivityUps += v.Value1
			}
			if v.Type == upcrmmodel.IncrType {
				old.IncrUps += v.Value1
			}
			if v.Type == upcrmmodel.TotalType {
				old.TotalUps += v.Value1
			}
			items[key] = old
		} else {
			new := &TrendDetail{
				Date: v.GenerateDate.Time().Unix(),
			}
			if v.Type == upcrmmodel.ActivityType {
				new.ActivityUps = v.Value1
			}
			if v.Type == upcrmmodel.IncrType {
				new.IncrUps = v.Value1
			}
			if v.Type == upcrmmodel.TotalType {
				new.TotalUps = v.Value1
			}
			items[key] = new
		}
	}

	// 这里的范围是[endDate - 6, endDate]，包含endDate
	for date := endDate.AddDate(0, 0, 1-days); !date.After(endDate); date = date.AddDate(0, 0, 1) {
		var dateStr = date.Format(upcrmmodel.TimeFmtDate)
		var d *TrendDetail
		var ok bool
		if d, ok = items[dateStr]; !ok {
			d = &TrendDetail{
				Date: date.Unix(),
			}
		}
		result = append(result, d)
	}
	return
}
