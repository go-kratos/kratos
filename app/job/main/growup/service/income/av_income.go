package income

import (
	"bytes"
	"context"
	"strconv"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// AvIncomeSvr Av income service
type AvIncomeSvr struct {
	batchSize int
	dao       *incomeD.Dao
}

// NewAvIncomeSvr new av income service
func NewAvIncomeSvr(dao *incomeD.Dao, batchSize int) (svr *AvIncomeSvr) {
	return &AvIncomeSvr{
		batchSize: batchSize,
		dao:       dao,
	}
}

// BatchInsertAvIncome batch insert av income
func (p *AvIncomeSvr) BatchInsertAvIncome(c context.Context, am map[int64][]*model.AvIncome) (err error) {
	var (
		buff    = make([]*model.AvIncome, p.batchSize)
		buffEnd = 0
	)

	for _, as := range am {
		for _, a := range as {
			buff[buffEnd] = a
			buffEnd++

			if buffEnd >= p.batchSize {
				values := avIncomeValues(buff[:buffEnd])
				buffEnd = 0
				_, err = p.dao.InsertAvIncome(c, values)
				if err != nil {
					return
				}
			}
		}
	}
	if buffEnd > 0 {
		values := avIncomeValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.dao.InsertAvIncome(c, values)
	}
	return
}

func avIncomeValues(as []*model.AvIncome) (values string) {
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
		buf.WriteString(strconv.FormatInt(a.PlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.Income, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TaxMoney, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + a.Date.Time().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.BaseIncome, 10))
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
