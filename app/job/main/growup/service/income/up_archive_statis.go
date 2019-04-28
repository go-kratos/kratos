package income

import (
	"context"
	"strconv"
	"strings"

	model "go-common/app/job/main/growup/model/income"

	xtime "go-common/library/time"
)

func (s *UpIncomeSvr) handleUpArchStatis(c context.Context, upArchStatisCh chan map[int64]*model.UpArchStatis, archiveCh chan []*model.ArchiveIncome) (err error) {
	defer close(upArchStatisCh)
	upArchMap := make(map[int64]*model.UpArchStatis)
	for income := range archiveCh {
		s.calUpArchStatis(income, upArchMap)
	}
	upArchStatisCh <- upArchMap
	return
}

func (s *UpIncomeSvr) calUpArchStatis(incomes []*model.ArchiveIncome, upArch map[int64]*model.UpArchStatis) {
	for _, income := range incomes {
		if _, ok := upArch[income.MID]; !ok {
			upArch[income.MID] = addUpArchStatis(income)
		}
		updateUpArchStatis(income.AID, income.Date, upArch[income.MID])
	}
}

func addUpArchStatis(income *model.ArchiveIncome) *model.UpArchStatis {
	return &model.UpArchStatis{
		MID:         income.MID,
		WeeklyDate:  xtime.Time(startWeeklyDate.Unix()),
		WeeklyAIDs:  "",
		MonthlyDate: xtime.Time(startMonthlyDate.Unix()),
		MonthlyAIDs: "",
	}
}

func updateUpArchStatis(aid int64, date xtime.Time, statis *model.UpArchStatis) {
	idStr := strconv.FormatInt(aid, 10)
	if date >= statis.WeeklyDate {
		if statis.WeeklyAIDs == "" {
			statis.WeeklyAIDs = idStr
		} else if !isExist(idStr, statis.WeeklyAIDs) {
			statis.WeeklyAIDs += "," + idStr
		}
	}
	if date >= statis.MonthlyDate {
		if statis.MonthlyAIDs == "" {
			statis.MonthlyAIDs = idStr
		} else if !isExist(idStr, statis.MonthlyAIDs) {
			statis.MonthlyAIDs += "," + idStr
		}
	}
}

func isExist(id string, old string) bool {
	oldSli := strings.Split(old, ",")
	for _, str := range oldSli {
		if id == str {
			return true
		}
	}
	return false
}
