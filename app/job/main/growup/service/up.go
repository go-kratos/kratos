package service

import (
	"bytes"
	"context"
	"strconv"

	"go-common/app/job/main/growup/model"
	"go-common/app/job/main/growup/model/income"
)

func (s *Service) getAllUps(c context.Context, limit int64) (m map[int64]*model.UpInfoVideo, err error) {
	var id int64
	m = make(map[int64]*model.UpInfoVideo)
	for {
		var us map[int64]*model.UpInfoVideo
		id, us, err = s.dao.UpInfoVideo(c, id, limit)
		if err != nil {
			return
		}
		for k, v := range us {
			m[k] = v
		}
		if len(us) < _dbLimit {
			break
		}
	}
	return
}

// SyncUpAccount sync up_account to up_tag_year
func (s *Service) SyncUpAccount(c context.Context) (err error) {
	var id int64
	for {
		var um map[int64]*income.UpAccount
		um, id, err = s.income.UpAccounts(c, id, 2000)
		if err != nil {
			return
		}
		if len(um) == 0 {
			break
		}
		_, err = s.tag.InsertUpYearAccount(c, assembleUpYear(um))
		if err != nil {
			return
		}
	}
	return
}

func assembleUpYear(ups map[int64]*income.UpAccount) (vals string) {
	var buf bytes.Buffer
	for mid, info := range ups {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(mid, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(info.TotalIncome, 10))
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
