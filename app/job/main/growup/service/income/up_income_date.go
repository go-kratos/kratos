package income

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"time"

	model "go-common/app/job/main/growup/model/income"

	"go-common/library/log"
	xtime "go-common/library/time"
)

var (
	_upIncomeWeekly  = "up_income_weekly"
	_upIncomeMonthly = "up_income_monthly"
)

func (s *UpIncomeSvr) handleUpIncomeWeeklyAndMonthly(
	c context.Context,
	date time.Time,
	upAvStatisCh chan map[int64]*model.UpArchStatis,
	upCmStatisCh chan map[int64]*model.UpArchStatis,
	upBgmStatisCh chan map[int64]*model.UpArchStatis,
	upSliceCh chan []*model.UpIncome) (weeklyMap, monthlyMap map[int64]*model.UpIncome, err error) {
	weeklyMap, monthlyMap, err = s.GetUpIncomeWeeklyAndMonthly(c, date)
	if err != nil {
		log.Error("s.GetUpIncomeWeeklyAndMonthly error(%v)", err)
		return
	}
	upAvStatis := <-upAvStatisCh
	upCmStatis := <-upCmStatisCh
	upBgmStatis := <-upBgmStatisCh
	s.calUpIncomeWeeklyAndMonthly(weeklyMap, monthlyMap, upAvStatis, upCmStatis, upBgmStatis, upSliceCh)
	return
}

// GetUpIncomeWeeklyAndMonthly get up_income_weekly and up_income_monthly
func (s *UpIncomeSvr) GetUpIncomeWeeklyAndMonthly(c context.Context, date time.Time) (weeklyMap map[int64]*model.UpIncome, monthlyMap map[int64]*model.UpIncome, err error) {
	upIncomeWeekly, err := s.GetUpIncomeTable(c, startWeeklyDate, _upIncomeWeekly)
	if err != nil {
		log.Error("s.GetUpIncomeTable error(%v)", err)
		return
	}

	upIncomeMonthly, err := s.GetUpIncomeTable(c, startMonthlyDate, _upIncomeMonthly)
	if err != nil {
		log.Error("s.GetUpIncomeTable error(%v)", err)
		return
	}

	weeklyMap = make(map[int64]*model.UpIncome)
	monthlyMap = make(map[int64]*model.UpIncome)
	for _, weeklyIncome := range upIncomeWeekly {
		weeklyMap[weeklyIncome.MID] = weeklyIncome
	}

	for _, monthlyIncome := range upIncomeMonthly {
		monthlyMap[monthlyIncome.MID] = monthlyIncome
	}

	return
}

// GetUpIncomeTable get up income table
func (s *UpIncomeSvr) GetUpIncomeTable(c context.Context, date time.Time, table string) (upIncomes []*model.UpIncome, err error) {
	var id int64
	for {
		upIncome, err1 := s.dao.GetUpIncomeTable(c, table, date.Format(_layout), id, _limitSize)
		if err1 != nil {
			err = err1
			return
		}
		upIncomes = append(upIncomes, upIncome...)
		if len(upIncome) < _limitSize {
			break
		}
		id = upIncome[len(upIncome)-1].ID
	}
	return
}

func (s *UpIncomeSvr) calUpIncomeWeeklyAndMonthly(weeklyMap, monthlyMap map[int64]*model.UpIncome,
	upAvStatis, upCmStatis, upBgmStatis map[int64]*model.UpArchStatis, upSliceCh chan []*model.UpIncome) {
	for upIncome := range upSliceCh {
		s.calUpIncome(upIncome, weeklyMap, monthlyMap, upAvStatis, upCmStatis, upBgmStatis)
	}
}

func (s *UpIncomeSvr) calUpIncome(upIncome []*model.UpIncome, weeklyMap, monthlyMap map[int64]*model.UpIncome,
	upAvStatis, upCmStatis, upBgmStatis map[int64]*model.UpArchStatis) {
	var weeklyAvCount, monthlyAvCount int
	var weeklyCmCount, monthlyCmCount int
	var weeklyBgmCount, monthlyBgmCount int
	for _, income := range upIncome {
		weeklyAvCount, monthlyAvCount = 0, 0
		weeklyCmCount, monthlyCmCount = 0, 0
		weeklyBgmCount, monthlyBgmCount = 0, 0
		if statis, ok := upAvStatis[income.MID]; ok {
			weeklyAvCount = len(strings.Split(statis.WeeklyAIDs, ","))
			monthlyAvCount = len(strings.Split(statis.MonthlyAIDs, ","))
		}

		if statis, ok := upCmStatis[income.MID]; ok {
			weeklyCmCount = len(strings.Split(statis.WeeklyAIDs, ","))
			monthlyCmCount = len(strings.Split(statis.MonthlyAIDs, ","))
		}

		if statis, ok := upBgmStatis[income.MID]; ok {
			weeklyBgmCount = len(strings.Split(statis.WeeklyAIDs, ","))
			monthlyBgmCount = len(strings.Split(statis.MonthlyAIDs, ","))
		}

		if weeklyIncome, ok := weeklyMap[income.MID]; ok {
			updateUpIncome(weeklyIncome, income, weeklyAvCount, weeklyCmCount, weeklyBgmCount)
		} else {
			weeklyMap[income.MID] = addUpIncome(income, startWeeklyDate, weeklyAvCount, weeklyCmCount, weeklyBgmCount)
		}

		if weeklyIncome, ok := monthlyMap[income.MID]; ok {
			updateUpIncome(weeklyIncome, income, monthlyAvCount, monthlyCmCount, monthlyBgmCount)
		} else {
			monthlyMap[income.MID] = addUpIncome(income, startMonthlyDate, monthlyAvCount, monthlyCmCount, monthlyBgmCount)
		}
	}
}

func addUpIncome(daily *model.UpIncome, fixDate time.Time, avCount, cmCount, bgmCount int) *model.UpIncome {
	return &model.UpIncome{
		MID:           daily.MID,
		AvCount:       int64(avCount),
		PlayCount:     daily.PlayCount,
		AvIncome:      daily.AvIncome,
		AvBaseIncome:  daily.AvBaseIncome,
		AvTax:         daily.AvTax,
		AvTotalIncome: daily.AvTotalIncome,

		ColumnCount:       int64(cmCount),
		ColumnIncome:      daily.ColumnIncome,
		ColumnBaseIncome:  daily.ColumnBaseIncome,
		ColumnTax:         daily.ColumnTax,
		ColumnTotalIncome: daily.ColumnTotalIncome,

		BgmCount:       int64(bgmCount),
		BgmIncome:      daily.BgmIncome,
		BgmBaseIncome:  daily.BgmBaseIncome,
		BgmTax:         daily.BgmTax,
		BgmTotalIncome: daily.BgmTotalIncome,

		AudioIncome: daily.AudioIncome,

		TaxMoney:    daily.TaxMoney,
		Income:      daily.Income,
		BaseIncome:  daily.BaseIncome,
		TotalIncome: daily.TotalIncome,
		Date:        xtime.Time(fixDate.Unix()),
		DBState:     _dbInsert,
	}
}

func updateUpIncome(origin, daily *model.UpIncome, avCount, cmCount, bgmCount int) {
	origin.AvCount = int64(avCount)
	origin.PlayCount += daily.PlayCount
	origin.AvIncome += daily.AvIncome
	origin.AvBaseIncome += daily.AvBaseIncome
	origin.AvTax += daily.AvTax
	origin.AvTotalIncome = daily.AvTotalIncome

	origin.ColumnCount = int64(cmCount)
	origin.ColumnIncome += daily.ColumnIncome
	origin.ColumnBaseIncome += daily.ColumnBaseIncome
	origin.ColumnTax += daily.ColumnTax
	origin.ColumnTotalIncome = daily.ColumnTotalIncome

	origin.BgmCount = int64(bgmCount)
	origin.BgmIncome += daily.BgmIncome
	origin.BgmBaseIncome += daily.BgmBaseIncome
	origin.BgmTax += daily.BgmTax
	origin.BgmTotalIncome = daily.BgmTotalIncome

	origin.AudioIncome += daily.AudioIncome

	origin.TaxMoney += daily.TaxMoney
	origin.Income += daily.Income
	origin.BaseIncome += daily.BaseIncome
	origin.TotalIncome = daily.TotalIncome
	origin.DBState = _dbUpdate
}

// UpIncomeDBStore insert up_income
func (s *UpIncomeSvr) UpIncomeDBStore(c context.Context, weeklyMap, monthlyMap map[int64]*model.UpIncome) (err error) {
	err = s.UpIncomeDBStoreBatch(c, _upIncomeWeekly, weeklyMap)
	if err != nil {
		log.Error("s.UpIncomeDBStoreBatch up_income_weekly error(%v)", err)
		return
	}

	err = s.UpIncomeDBStoreBatch(c, _upIncomeMonthly, monthlyMap)
	if err != nil {
		log.Error("s.UpIncomeDBStoreBatch up_income_monthly error(%v)", err)
		return
	}
	return
}

// UpIncomeDBStoreBatch up income db batch store
func (s *UpIncomeSvr) UpIncomeDBStoreBatch(c context.Context, table string, upIncomeMap map[int64]*model.UpIncome) error {
	insert, update := make([]*model.UpIncome, batchSize), make([]*model.UpIncome, batchSize)
	insertIndex, updateIndex := 0, 0
	for _, income := range upIncomeMap {
		if income.DBState == _dbInsert {
			insert[insertIndex] = income
			insertIndex++
		} else if income.DBState == _dbUpdate {
			update[updateIndex] = income
			updateIndex++
		}

		if insertIndex >= batchSize {
			_, err := s.upIncomeBatchInsert(c, table, insert[:insertIndex])
			if err != nil {
				log.Error("s.upIncomeBatchInsert error(%v)", err)
				return err
			}
			insertIndex = 0
		}

		if updateIndex >= batchSize {
			_, err := s.upIncomeBatchInsert(c, table, update[:updateIndex])
			if err != nil {
				log.Error("s.upIncomeBatchInsert error(%v)", err)
				return err
			}
			updateIndex = 0
		}
	}

	if insertIndex > 0 {
		_, err := s.upIncomeBatchInsert(c, table, insert[:insertIndex])
		if err != nil {
			log.Error("s.upIncomeBatchInsert error(%v)", err)
			return err
		}
	}

	if updateIndex > 0 {
		_, err := s.upIncomeBatchInsert(c, table, update[:updateIndex])
		if err != nil {
			log.Error("s.upIncomeBatchInsert error(%v)", err)
			return err
		}
	}

	return nil
}

func (s *UpIncomeSvr) upIncomeBatchInsert(c context.Context, table string, us []*model.UpIncome) (rows int64, err error) {
	var buf bytes.Buffer
	for _, u := range us {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(u.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.PlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AudioIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.TaxMoney, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.Income, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.TotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvBaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvTax, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnBaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnTax, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmBaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmTax, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + u.Date.Time().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvTotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnTotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmTotalIncome, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values := buf.String()
	buf.Reset()
	rows, err = s.dao.InsertUpIncomeTable(c, table, values)
	return
}
