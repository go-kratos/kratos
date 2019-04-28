package income

import (
	"context"
	"time"

	model "go-common/app/job/main/growup/model/income"
)

func (s *DateStatis) getArchiveByDate(c context.Context, archiveCh chan []*model.ArchiveIncome, startDate, endDate time.Time, typ, limit int) (err error) {
	defer close(archiveCh)
	var aid, table string
	switch typ {
	case _video:
		aid, table = "av_id", "av_income"
	case _column:
		aid, table = "aid", "column_income"
	case _bgm:
		aid, table = "sid", "bgm_income"
	default:
		return
	}
	endDate = endDate.AddDate(0, 0, 1)
	for startDate.Before(endDate) {
		err = s.getArchiveIncome(c, archiveCh, aid, table, startDate.Format(_layout), limit)
		if err != nil {
			return
		}
		startDate = startDate.AddDate(0, 0, 1)
	}
	return
}

func (s *DateStatis) getArchiveIncome(c context.Context, archiveCh chan []*model.ArchiveIncome, aid, table, date string, limit int) (err error) {
	var id int64
	for {
		var archive []*model.ArchiveIncome
		if aid == "sid" {
			archive, err = s.dao.GetBgmIncomeByDate(c, date, id, limit)
		} else {
			archive, err = s.dao.GetArchiveByDate(c, aid, table, date, id, limit)
		}
		if err != nil {
			return
		}
		archiveCh <- archive
		if len(archive) < limit {
			break
		}
		id = archive[len(archive)-1].ID
	}
	return
}
