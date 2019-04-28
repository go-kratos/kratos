package charge

import (
	"bytes"
	"context"
	"strconv"
	"time"

	dao "go-common/app/job/main/growup/dao/charge"
	model "go-common/app/job/main/growup/model/charge"

	"go-common/library/log"
	xtime "go-common/library/time"

	"golang.org/x/sync/errgroup"
)

func (s *Service) handleAvCharge(c context.Context, date time.Time,
	dailyChannel chan []*model.AvCharge) (weeklyChargeMap, monthlyChargeMap map[int64]*model.AvCharge, chargeStatisMap map[int64]*model.AvChargeStatis, err error) {
	var eg errgroup.Group
	weeklyChargeMap = make(map[int64]*model.AvCharge)
	monthlyChargeMap = make(map[int64]*model.AvCharge)

	eg.Go(func() (err error) {
		avWeeklyCharge, err := s.GetAvCharge(c, getStartWeeklyDate(date), s.dao.AvWeeklyCharge)
		if err != nil {
			log.Error("s.GetAvCharge(av_weekly_charge) error(%v)", err)
			return
		}
		for _, weeklyCharge := range avWeeklyCharge {
			weeklyChargeMap[weeklyCharge.AvID] = weeklyCharge
		}
		return
	})

	eg.Go(func() (err error) {
		avMonthlyCharge, err := s.GetAvCharge(c, getStartMonthlyDate(date), s.dao.AvMonthlyCharge)
		if err != nil {
			log.Error("s.GetAvCharge(av_monthly_charge) error(%v)", err)
			return
		}
		for _, monthlyCharge := range avMonthlyCharge {
			monthlyChargeMap[monthlyCharge.AvID] = monthlyCharge
		}
		return
	})

	eg.Go(func() (err error) {
		chargeStatisMap, err = s.GetAvChargeStatisMap(c)
		if err != nil {
			log.Error("s.GetAvChargeStatisMap error(%v)", err)
		}
		return
	})

	if err = eg.Wait(); err != nil {
		log.Error("HandleAvCharge eg.Wait error(%v)", err)
		return
	}

	s.calAvCharges(date, weeklyChargeMap, monthlyChargeMap, chargeStatisMap, dailyChannel)
	return
}

func (s *Service) avDailyCharges(c context.Context, date time.Time, ch chan []*model.AvCharge) (err error) {
	defer func() {
		close(ch)
	}()
	var id int64
	for {
		var charges []*model.AvCharge
		charges, err = s.dao.AvDailyCharge(c, date, id, _limitSize)
		if err != nil {
			return
		}
		ch <- charges
		if len(charges) < _limitSize {
			break
		}
		id = charges[len(charges)-1].ID
	}
	return
}

// GetAvCharge get av charge
func (s *Service) GetAvCharge(c context.Context, date time.Time, fn dao.IAvCharge) (avCharges []*model.AvCharge, err error) {
	var id int64
	for {
		var avCharge []*model.AvCharge
		avCharge, err = fn(c, date, id, _limitSize)
		if err != nil {
			return
		}
		avCharges = append(avCharges, avCharge...)
		if len(avCharge) < _limitSize {
			break
		}
		id = avCharge[len(avCharge)-1].ID
	}
	return
}

func (s *Service) calAvCharges(date time.Time, weeklyChargeMap, monthlyChargeMap map[int64]*model.AvCharge, chargeStatisMap map[int64]*model.AvChargeStatis, dailyChannel chan []*model.AvCharge) {
	for avDailyCharge := range dailyChannel {
		s.calAvCharge(date, avDailyCharge, weeklyChargeMap, monthlyChargeMap, chargeStatisMap)
	}
}

func (s *Service) calAvCharge(date time.Time, avDailyCharge []*model.AvCharge, weeklyChargeMap, monthlyChargeMap map[int64]*model.AvCharge, chargeStatisMap map[int64]*model.AvChargeStatis) {
	for _, dailyCharge := range avDailyCharge {
		if weeklyCharge, ok := weeklyChargeMap[dailyCharge.AvID]; ok {
			updateAvCharge(weeklyCharge, dailyCharge)
		} else {
			weeklyChargeMap[dailyCharge.AvID] = addAvCharge(dailyCharge, startWeeklyDate)
		}

		if monthlyCharge, ok := monthlyChargeMap[dailyCharge.AvID]; ok {
			updateAvCharge(monthlyCharge, dailyCharge)
		} else {
			monthlyChargeMap[dailyCharge.AvID] = addAvCharge(dailyCharge, startMonthlyDate)
		}
		s.CalAvChargeStatis(dailyCharge, chargeStatisMap)
	}
}

func addAvCharge(daily *model.AvCharge, fixDate time.Time) *model.AvCharge {
	return &model.AvCharge{
		AvID:           daily.AvID,
		MID:            daily.MID,
		TagID:          daily.TagID,
		IsOriginal:     daily.IsOriginal,
		DanmakuCount:   daily.DanmakuCount,
		CommentCount:   daily.CommentCount,
		CollectCount:   daily.CollectCount,
		CoinCount:      daily.CoinCount,
		ShareCount:     daily.ShareCount,
		ElecPayCount:   daily.ElecPayCount,
		TotalPlayCount: daily.TotalPlayCount,
		WebPlayCount:   daily.WebPlayCount,
		AppPlayCount:   daily.AppPlayCount,
		H5PlayCount:    daily.H5PlayCount,
		LvUnknown:      daily.LvUnknown,
		Lv0:            daily.Lv0,
		Lv1:            daily.Lv1,
		Lv2:            daily.Lv2,
		Lv3:            daily.Lv3,
		Lv4:            daily.Lv4,
		Lv5:            daily.Lv5,
		Lv6:            daily.Lv6,
		VScore:         daily.VScore,
		IncCharge:      daily.IncCharge,
		TotalCharge:    daily.IncCharge,
		Date:           xtime.Time(fixDate.Unix()),
		UploadTime:     daily.UploadTime,
		DBState:        _dbInsert,
	}
}

func updateAvCharge(origin, daily *model.AvCharge) {
	origin.DanmakuCount += daily.DanmakuCount
	origin.CommentCount += daily.CommentCount
	origin.CollectCount += daily.CollectCount
	origin.CoinCount += daily.CoinCount
	origin.ShareCount += daily.ShareCount
	origin.ElecPayCount += daily.ElecPayCount
	origin.TotalPlayCount += daily.TotalPlayCount
	origin.WebPlayCount += daily.WebPlayCount
	origin.AppPlayCount += daily.AppPlayCount
	origin.H5PlayCount += daily.H5PlayCount
	origin.LvUnknown += daily.LvUnknown
	origin.Lv0 += daily.Lv0
	origin.Lv1 += daily.Lv1
	origin.Lv2 += daily.Lv2
	origin.Lv3 += daily.Lv3
	origin.Lv4 += daily.Lv4
	origin.Lv5 += daily.Lv5
	origin.Lv6 += daily.Lv6
	origin.VScore += daily.VScore
	origin.IncCharge += daily.IncCharge
	origin.TotalCharge += daily.IncCharge
	origin.DBState = _dbUpdate
}

// AvChargeDBStore store data
func (s *Service) AvChargeDBStore(c context.Context, table string, avChargeMap map[int64]*model.AvCharge) error {
	insert, update := make([]*model.AvCharge, _batchSize), make([]*model.AvCharge, _batchSize)
	insertIndex, updateIndex := 0, 0
	for _, charge := range avChargeMap {
		if charge.DBState == _dbInsert {
			insert[insertIndex] = charge
			insertIndex++
		} else if charge.DBState == _dbUpdate {
			update[updateIndex] = charge
			updateIndex++
		}

		if insertIndex >= _batchSize {
			_, err := s.avChargeBatchInsert(c, insert[:insertIndex], table)
			if err != nil {
				log.Error("s.avChargeBatchInsert error(%v)", err)
				return err
			}
			insertIndex = 0
		}

		if updateIndex >= _batchSize {
			_, err := s.avChargeBatchInsert(c, update[:updateIndex], table)
			if err != nil {
				log.Error("s.avChargeBatchInsert error(%v)", err)
				return err
			}
			updateIndex = 0
		}
	}

	if insertIndex > 0 {
		_, err := s.avChargeBatchInsert(c, insert[:insertIndex], table)
		if err != nil {
			log.Error("s.avChargeBatchInsert error(%v)", err)
			return err
		}
	}

	if updateIndex > 0 {
		_, err := s.avChargeBatchInsert(c, update[:updateIndex], table)
		if err != nil {
			log.Error("s.avChargeBatchInsert error(%v)", err)
			return err
		}
	}

	return nil
}

func assembleAvCharge(avCharge []*model.AvCharge) (vals string) {
	var buf bytes.Buffer
	for _, row := range avCharge {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.AvID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(row.IsOriginal))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.DanmakuCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CommentCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CollectCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CoinCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.ShareCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.ElecPayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TotalPlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.WebPlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.AppPlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.H5PlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.LvUnknown, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv0, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv1, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv2, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv3, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv4, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv5, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv6, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.VScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.IncCharge, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TotalCharge, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + row.Date.Time().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString("'" + row.UploadTime.Time().Format(_layoutSec) + "'")
		buf.WriteString(")")
		buf.WriteByte(',')
	}

	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals = buf.String()
	buf.Reset()
	return
}

func (s *Service) avChargeBatchInsert(c context.Context, avCharge []*model.AvCharge, table string) (rows int64, err error) {
	vals := assembleAvCharge(avCharge)
	rows, err = s.dao.InsertAvChargeTable(c, vals, table)
	return
}
