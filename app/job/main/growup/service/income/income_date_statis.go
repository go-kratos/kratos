package income

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"

	"go-common/library/log"
	xtime "go-common/library/time"
)

var (
	_avIncomeDailyStatis   = "av_income_daily_statis"
	_avIncomeWeeklyStatis  = "av_income_weekly_statis"
	_avIncomeMonthlyStatis = "av_income_monthly_statis"

	_cmIncomeDailyStatis   = "column_income_daily_statis"
	_cmIncomeWeeklyStatis  = "column_income_weekly_statis"
	_cmIncomeMonthlyStatis = "column_income_monthly_statis"

	_bgmIncomeDailyStatis   = "bgm_income_daily_statis"
	_bgmIncomeWeeklyStatis  = "bgm_income_weekly_statis"
	_bgmIncomeMonthlyStatis = "bgm_income_monthly_statis"
)

// SectionEntries section entries
type SectionEntries struct {
	avDaily   []*model.DateStatis
	avWeekly  []*model.DateStatis
	avMonthly []*model.DateStatis
}

// DateStatis income date statistics
type DateStatis struct {
	dao *incomeD.Dao
}

// NewDateStatis new income date statistics service
func NewDateStatis(dao *incomeD.Dao) *DateStatis {
	return &DateStatis{dao: dao}
}

func initIncomeSections(income, tagID int64, date xtime.Time) []*model.DateStatis {
	incomeSections := make([]*model.DateStatis, 12)
	incomeSections[0] = initIncomeSection(0, 1, 0, income, tagID, date)
	incomeSections[1] = initIncomeSection(1, 5, 1, income, tagID, date)
	incomeSections[2] = initIncomeSection(5, 10, 2, income, tagID, date)
	incomeSections[3] = initIncomeSection(10, 30, 3, income, tagID, date)
	incomeSections[4] = initIncomeSection(30, 50, 4, income, tagID, date)
	incomeSections[5] = initIncomeSection(50, 100, 5, income, tagID, date)
	incomeSections[6] = initIncomeSection(100, 200, 6, income, tagID, date)
	incomeSections[7] = initIncomeSection(200, 500, 7, income, tagID, date)
	incomeSections[8] = initIncomeSection(500, 1000, 8, income, tagID, date)
	incomeSections[9] = initIncomeSection(1000, 3000, 9, income, tagID, date)
	incomeSections[10] = initIncomeSection(3000, 5000, 10, income, tagID, date)
	incomeSections[11] = initIncomeSection(5000, math.MaxInt32, 11, income, tagID, date)
	return incomeSections
}

func initIncomeSection(min, max, section, income, tagID int64, date xtime.Time) *model.DateStatis {
	var tips string
	if max == math.MaxInt32 {
		tips = fmt.Sprintf("\"%d+\"", min)
	} else {
		tips = fmt.Sprintf("\"%d~%d\"", min, max)
	}
	return &model.DateStatis{
		MinIncome:    min,
		MaxIncome:    max,
		MoneySection: section,
		MoneyTips:    tips,
		Income:       income,
		CategoryID:   tagID,
		CDate:        date,
	}
}

func (s *DateStatis) handleDateStatis(c context.Context, archiveCh chan []*model.ArchiveIncome, date time.Time, table string) (incomeSections []*model.DateStatis, err error) {
	// delete
	if table != "" {
		_, err = s.dao.DelIncomeStatisTable(c, table, date.Format(_layout))
		if err != nil {
			log.Error("s.dao.DelIncomeStatisTable error(%v)", err)
			return
		}
	}
	// add
	incomeSections = s.handleArchives(c, archiveCh, date)
	return
}

// handleArchives handle archive_income_daily_statis, archive_income_weekly_statis, archive_income_monthly_statis
func (s *DateStatis) handleArchives(c context.Context, archiveCh chan []*model.ArchiveIncome, date time.Time) (incomeSections []*model.DateStatis) {
	archTagMap := make(map[int64]map[int64]int64) // key TagID, value map[int64]int64 -> key aid, value income
	tagIncomeMap := make(map[int64]int64)         // key TagID, value TagID total income
	for archive := range archiveCh {
		handleArchive(archive, archTagMap, tagIncomeMap, date)
	}
	incomeSections = make([]*model.DateStatis, 0)
	for tagID, avMap := range archTagMap {
		incomeSection := countIncomeDailyStatis(avMap, tagIncomeMap[tagID], tagID, date)
		incomeSections = append(incomeSections, incomeSection...)
	}
	return
}

func handleArchive(archives []*model.ArchiveIncome, archTagMap map[int64]map[int64]int64, tagIncomeMap map[int64]int64, startDate time.Time) {
	if archives == nil {
		return
	}
	if archTagMap == nil {
		archTagMap = make(map[int64]map[int64]int64)
	}
	if tagIncomeMap == nil {
		tagIncomeMap = make(map[int64]int64)
	}

	for _, archive := range archives {
		if !startDate.After(archive.Date.Time()) {
			tagIncomeMap[archive.TagID] += archive.Income
			if _, ok := archTagMap[archive.TagID]; !ok {
				archTagMap[archive.TagID] = make(map[int64]int64)
			}
			archTagMap[archive.TagID][archive.AID] += archive.Income
		}
	}
}

func (s *DateStatis) handleDateUp(c context.Context, upStatisCh chan []*model.UpIncome, date time.Time) (upSections, upAvSections, upCmSections, upBgmSections []*model.DateStatis, err error) {
	_, err = s.dao.DelIncomeStatisTable(c, "up_income_daily_statis", date.Format(_layout))
	if err != nil {
		log.Error("s.dao.DelIncomeStatisTable error(%v)", err)
		return
	}

	_, err = s.dao.DelIncomeStatisTable(c, "up_av_daily_statis", date.Format(_layout))
	if err != nil {
		log.Error("s.dao.DelIncomeStatisTable error(%v)", err)
		return
	}

	_, err = s.dao.DelIncomeStatisTable(c, "up_column_daily_statis", date.Format(_layout))
	if err != nil {
		log.Error("s.dao.DelIncomeStatisTable error(%v)", err)
		return
	}

	_, err = s.dao.DelIncomeStatisTable(c, "up_bgm_daily_statis", date.Format(_layout))
	if err != nil {
		log.Error("s.dao.DelIncomeStatisTable error(%v)", err)
		return
	}

	upMap := make(map[int64]int64)
	upAvMap := make(map[int64]int64)
	upCmMap := make(map[int64]int64)
	upBgmMap := make(map[int64]int64)
	var upTotal, avTotal, cmTotal, bgmTotal int64
	for up := range upStatisCh {
		up, av, cm, bgm := handleUp(up, upMap, upAvMap, upCmMap, upBgmMap, date)
		upTotal += up
		avTotal += av
		cmTotal += cm
		bgmTotal += bgm
	}
	upSections = countIncomeDailyStatis(upMap, upTotal, 0, date)
	upAvSections = countIncomeDailyStatis(upAvMap, avTotal, 0, date)
	upCmSections = countIncomeDailyStatis(upCmMap, cmTotal, 0, date)
	upBgmSections = countIncomeDailyStatis(upBgmMap, bgmTotal, 0, date)
	return
}

func handleUp(upIncomes []*model.UpIncome, upMap, upAvMap, upCmMap, upBgmMap map[int64]int64, startDate time.Time) (income, avIncome, cmIncome, bgmIncome int64) {
	if len(upIncomes) == 0 {
		return
	}
	for _, upIncome := range upIncomes {
		if startDate.Equal(upIncome.Date.Time()) {
			income += upIncome.Income
			avIncome += upIncome.AvIncome
			cmIncome += upIncome.ColumnIncome
			bgmIncome += upIncome.BgmIncome
			upMap[upIncome.MID] += upIncome.Income
			if upIncome.AvIncome > 0 {
				upAvMap[upIncome.MID] += upIncome.AvIncome
			}
			if upIncome.ColumnIncome > 0 {
				upCmMap[upIncome.MID] += upIncome.ColumnIncome
			}
			if upIncome.BgmIncome > 0 {
				upBgmMap[upIncome.MID] += upIncome.BgmIncome
			}
		}
	}
	return
}

func (s *DateStatis) incomeDateStatisInsert(c context.Context, incomeSection []*model.DateStatis, table string) (rows int64, err error) {
	var buf bytes.Buffer
	for _, row := range incomeSection {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.Count, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MoneySection, 10))
		buf.WriteByte(',')
		buf.WriteString(row.MoneyTips)
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Income, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CategoryID, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + row.CDate.Time().Format(_layout) + "'")
		buf.WriteString(")")
		buf.WriteByte(',')
	}

	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals := buf.String()
	buf.Reset()
	rows, err = s.dao.InsertIncomeStatisTable(c, table, vals)
	return
}

func countIncomeDailyStatis(incomes map[int64]int64, totalIncome, tagID int64, date time.Time) (incomeSections []*model.DateStatis) {
	if len(incomes) == 0 {
		return
	}
	incomeSections = initIncomeSections(totalIncome, tagID, xtime.Time(date.Unix()))
	for _, income := range incomes {
		for _, section := range incomeSections {
			min, max := section.MinIncome*100, section.MaxIncome*100
			if income >= min && income < max {
				section.Count++
			}
		}
	}
	return
}

func (s *DateStatis) upIncomeDailyStatisInsert(c context.Context, upIncomeSection []*model.DateStatis, table string) (rows int64, err error) {
	var buf bytes.Buffer
	for _, row := range upIncomeSection {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.Count, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MoneySection, 10))
		buf.WriteByte(',')
		buf.WriteString(row.MoneyTips)
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Income, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + row.CDate.Time().Format(_layout) + "'")
		buf.WriteString(")")
		buf.WriteByte(',')
	}

	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals := buf.String()
	buf.Reset()
	rows, err = s.dao.InsertUpIncomeDailyStatis(c, table, vals)
	return
}
