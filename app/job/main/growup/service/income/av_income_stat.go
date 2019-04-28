package income

import (
	"bytes"
	"context"
	"strconv"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// AvIncomeStatSvr av income stat service
type AvIncomeStatSvr struct {
	batchSize int
	dao       *incomeD.Dao
}

// NewAvIncomeStatSvr new av income service
func NewAvIncomeStatSvr(dao *incomeD.Dao, batchSize int) (svr *AvIncomeStatSvr) {
	return &AvIncomeStatSvr{
		batchSize: batchSize,
		dao:       dao,
	}
}

// AvIncomeStat get av income stat
func (p *AvIncomeStatSvr) AvIncomeStat(c context.Context, limit int64) (m map[int64]*model.AvIncomeStat, err error) {
	m = make(map[int64]*model.AvIncomeStat)
	var id int64
	for {
		var am map[int64]*model.AvIncomeStat
		am, id, err = p.dao.AvIncomeStat(c, id, limit)
		if err != nil {
			return
		}
		if len(am) == 0 {
			break
		}
		for avID, stat := range am {
			m[avID] = stat
		}
	}
	return
}

// BatchInsertAvIncomeStat batch insert av income statistics
func (p *AvIncomeStatSvr) BatchInsertAvIncomeStat(c context.Context, as map[int64]*model.AvIncomeStat) (err error) {
	var (
		buff    = make([]*model.AvIncomeStat, p.batchSize)
		buffEnd = 0
	)
	for _, a := range as {
		if a.DataState == 0 {
			continue
		}
		buff[buffEnd] = a
		buffEnd++

		if buffEnd >= p.batchSize {
			values := avIncomeStatValues(buff[:buffEnd])
			buffEnd = 0
			_, err = p.dao.InsertAvIncomeStat(c, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := avIncomeStatValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.dao.InsertAvIncomeStat(c, values)
	}
	return
}

func avIncomeStatValues(as []*model.AvIncomeStat) (values string) {
	var buf bytes.Buffer
	for _, a := range as {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(a.AvID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(a.IsOriginal))
		buf.WriteByte(',')
		buf.WriteString("'" + a.UploadTime.Time().Format(_layoutSec) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TotalIncome, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}
