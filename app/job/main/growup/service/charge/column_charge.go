package charge

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	model "go-common/app/job/main/growup/model/charge"
	task "go-common/app/job/main/growup/service"

	"go-common/library/log"
	xtime "go-common/library/time"
	"golang.org/x/sync/errgroup"
)

func (s *Service) columnCharges(c context.Context, date time.Time, ch chan []*model.Column) (err error) {
	defer func() {
		close(ch)
	}()
	var id int64
	for {
		var charges []*model.Column
		charges, err = s.dao.ColumnCharge(c, date, id, _limitSize, "column_daily_charge")
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

func (s *Service) handleColumn(c context.Context, date time.Time,
	dailyChannel chan []*model.Column) (weeklyMap, monthlyMap map[int64]*model.Column, statisMap map[int64]*model.ColumnStatis, err error) {
	var eg errgroup.Group
	weeklyMap = make(map[int64]*model.Column)
	monthlyMap = make(map[int64]*model.Column)
	statisMap = make(map[int64]*model.ColumnStatis)

	eg.Go(func() (err error) {
		weekly, err := s.getColumnCharge(c, getStartWeeklyDate(date), "column_weekly_charge")
		if err != nil {
			log.Error("s.GetColumnCharge(column_weekly_charge) error(%v)", err)
			return
		}
		for _, cm := range weekly {
			weeklyMap[cm.AID] = cm
		}
		return
	})

	eg.Go(func() (err error) {
		monthly, err := s.getColumnCharge(c, getStartMonthlyDate(date), "column_monthly_charge")
		if err != nil {
			log.Error("s.GetColumnCharge(column_monthly_charge) error(%v)", err)
			return
		}
		for _, cm := range monthly {
			monthlyMap[cm.AID] = cm
		}
		return
	})

	eg.Go(func() (err error) {
		statis, err := s.getCmStatisMap(c)
		if err != nil {
			log.Error("s.getCmStatisMap error(%v)", err)
		}
		for _, cm := range statis {
			statisMap[cm.AID] = cm
		}
		return
	})

	if err = eg.Wait(); err != nil {
		log.Error("handleColumn eg.Wait error(%v)", err)
		return
	}

	s.calColumns(date, weeklyMap, monthlyMap, statisMap, dailyChannel)
	return
}

func (s *Service) calColumns(date time.Time, weeklyMap, monthlyMap map[int64]*model.Column, statisMap map[int64]*model.ColumnStatis, dailyChannel chan []*model.Column) {
	for daily := range dailyChannel {
		s.calColumn(date, daily, weeklyMap, monthlyMap, statisMap)
	}
}

func (s *Service) calColumn(date time.Time, daily []*model.Column, weeklyMap, monthlyMap map[int64]*model.Column, statisMap map[int64]*model.ColumnStatis) {
	for _, charge := range daily {
		if weekly, ok := weeklyMap[charge.AID]; ok {
			updateColumnCharge(weekly, charge)
		} else {
			weeklyMap[charge.AID] = addColumnCharge(charge, startWeeklyDate)
		}
		if monthly, ok := monthlyMap[charge.AID]; ok {
			updateColumnCharge(monthly, charge)
		} else {
			monthlyMap[charge.AID] = addColumnCharge(charge, startMonthlyDate)
		}
		if statis, ok := statisMap[charge.AID]; ok {
			updateColumnStatis(statis, charge)
		} else {
			statisMap[charge.AID] = addColumnStatis(charge)
		}
	}
}

func addColumnCharge(daily *model.Column, fixDate time.Time) *model.Column {
	return &model.Column{
		AID:        daily.AID,
		MID:        daily.MID,
		Title:      daily.Title,
		TagID:      daily.TagID,
		Words:      daily.Words,
		UploadTime: daily.UploadTime,
		IncCharge:  daily.IncCharge,
		Date:       xtime.Time(fixDate.Unix()),
		DBState:    _dbInsert,
	}
}

func updateColumnCharge(origin, daily *model.Column) {
	origin.IncCharge += daily.IncCharge
	origin.DBState = _dbUpdate
}

func addColumnStatis(daily *model.Column) *model.ColumnStatis {
	return &model.ColumnStatis{
		AID:         daily.AID,
		MID:         daily.MID,
		Title:       daily.Title,
		TagID:       daily.TagID,
		UploadTime:  daily.UploadTime,
		TotalCharge: daily.IncCharge,
		DBState:     _dbInsert,
	}
}

func updateColumnStatis(statis *model.ColumnStatis, daily *model.Column) {
	statis.TotalCharge += daily.IncCharge
	statis.DBState = _dbUpdate
}

func (s *Service) getColumnCharge(c context.Context, date time.Time, table string) (cms []*model.Column, err error) {
	cms = make([]*model.Column, 0)
	var id int64
	for {
		var cm []*model.Column
		cm, err = s.dao.ColumnCharge(c, date, id, _limitSize, table)
		if err != nil {
			return
		}
		cms = append(cms, cm...)
		if len(cm) < _limitSize {
			break
		}
		id = cm[len(cm)-1].ID
	}
	return
}

func (s *Service) getCmStatisMap(c context.Context) (cms []*model.ColumnStatis, err error) {
	cms = make([]*model.ColumnStatis, 0)
	var id int64
	for {
		var cm []*model.ColumnStatis
		cm, err = s.dao.CmStatis(c, id, _limitSize)
		if err != nil {
			return
		}
		cms = append(cms, cm...)
		if len(cm) < _limitSize {
			break
		}
		id = cm[len(cm)-1].ID
	}
	return
}

func (s *Service) cmDBStore(c context.Context, table string, cmMap map[int64]*model.Column) error {
	insert, update := make([]*model.Column, _batchSize), make([]*model.Column, _batchSize)
	insertIndex, updateIndex := 0, 0
	for _, charge := range cmMap {
		if charge.DBState == _dbInsert {
			insert[insertIndex] = charge
			insertIndex++
		} else if charge.DBState == _dbUpdate {
			update[updateIndex] = charge
			updateIndex++
		}

		if insertIndex >= _batchSize {
			_, err := s.cmBatchInsert(c, insert[:insertIndex], table)
			if err != nil {
				log.Error("s.cmBatchInsert error(%v)", err)
				return err
			}
			insertIndex = 0
		}

		if updateIndex >= _batchSize {
			_, err := s.cmBatchInsert(c, update[:updateIndex], table)
			if err != nil {
				log.Error("s.cmBatchInsert error(%v)", err)
				return err
			}
			updateIndex = 0
		}
	}

	if insertIndex > 0 {
		_, err := s.cmBatchInsert(c, insert[:insertIndex], table)
		if err != nil {
			log.Error("s.cmBatchInsert error(%v)", err)
			return err
		}
	}

	if updateIndex > 0 {
		_, err := s.cmBatchInsert(c, update[:updateIndex], table)
		if err != nil {
			log.Error("s.cmBatchInsert error(%v)", err)
			return err
		}
	}
	return nil
}

func assembleCmCharge(charges []*model.Column) (vals string) {
	var buf bytes.Buffer
	for _, row := range charges {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.AID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.IncCharge, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + row.Date.Time().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.UploadTime, 10))
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

func (s *Service) cmBatchInsert(c context.Context, charges []*model.Column, table string) (rows int64, err error) {
	vals := assembleCmCharge(charges)
	rows, err = s.dao.InsertCmChargeTable(c, vals, table)
	return
}

func (s *Service) cmStatisDBStore(c context.Context, statisMap map[int64]*model.ColumnStatis) error {
	insert, update := make([]*model.ColumnStatis, _batchSize), make([]*model.ColumnStatis, _batchSize)
	insertIndex, updateIndex := 0, 0
	for _, charge := range statisMap {
		if charge.DBState == _dbInsert {
			insert[insertIndex] = charge
			insertIndex++
		} else if charge.DBState == _dbUpdate {
			update[updateIndex] = charge
			updateIndex++
		}
		if insertIndex >= _batchSize {
			_, err := s.cmStatisBatchInsert(c, insert[:insertIndex])
			if err != nil {
				log.Error("s.cmStatisBatchInsert error(%v)", err)
				return err
			}
			insertIndex = 0
		}

		if updateIndex >= _batchSize {
			_, err := s.cmStatisBatchInsert(c, update[:updateIndex])
			if err != nil {
				log.Error("s.cmStatisBatchInsert error(%v)", err)
				return err
			}
			updateIndex = 0
		}
	}

	if insertIndex > 0 {
		_, err := s.cmStatisBatchInsert(c, insert[:insertIndex])
		if err != nil {
			log.Error("s.cmStatisBatchInsert error(%v)", err)
			return err
		}
	}

	if updateIndex > 0 {
		_, err := s.cmStatisBatchInsert(c, update[:updateIndex])
		if err != nil {
			log.Error("s.cmStatisBatchInsert error(%v)", err)
			return err
		}
	}
	return nil
}

func assembleCmStatis(statis []*model.ColumnStatis) (vals string) {
	var buf bytes.Buffer
	for _, row := range statis {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.AID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TotalCharge, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + time.Unix(row.UploadTime, 0).Format(_layoutSec) + "'")
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

func (s *Service) cmStatisBatchInsert(c context.Context, statis []*model.ColumnStatis) (rows int64, err error) {
	vals := assembleCmStatis(statis)
	rows, err = s.dao.InsertCmStatisBatch(c, vals)
	return
}

// CheckTaskColumn  check column count by date
func (s *Service) CheckTaskColumn(c context.Context, date string) (err error) {
	defer func() {
		task.GetTaskService().SetTaskStatus(c, task.TaskCmCharge, date, err)
	}()

	count, err := s.dao.CountCmDailyCharge(c, date)
	if count == 0 {
		err = fmt.Errorf("date(%s) column_daily_charge = 0", date)
	}
	return
}
