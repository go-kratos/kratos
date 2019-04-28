package income

import (
	"bytes"
	"context"
	"strconv"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// ColumnIncomeSvr column income service
type ColumnIncomeSvr struct {
	batchSize int
	dao       *incomeD.Dao
}

// NewColumnIncomeSvr new income service
func NewColumnIncomeSvr(dao *incomeD.Dao, batchSize int) (svr *ColumnIncomeSvr) {
	return &ColumnIncomeSvr{
		batchSize: batchSize,
		dao:       dao,
	}
}

// BatchInsertColumnIncome batch insert column income
func (p *ColumnIncomeSvr) BatchInsertColumnIncome(ctx context.Context, cm map[int64][]*model.ColumnIncome) (err error) {
	var (
		buff    = make([]*model.ColumnIncome, p.batchSize)
		buffEnd = 0
	)

	for _, cs := range cm {
		for _, c := range cs {
			buff[buffEnd] = c
			buffEnd++

			if buffEnd >= p.batchSize {
				values := columnIncomeValues(buff[:buffEnd])
				buffEnd = 0
				_, err = p.dao.InsertColumnIncome(ctx, values)
				if err != nil {
					return
				}
			}
		}
	}
	if buffEnd > 0 {
		values := columnIncomeValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.dao.InsertColumnIncome(ctx, values)
	}
	return
}

func columnIncomeValues(cs []*model.ColumnIncome) (values string) {
	var buf bytes.Buffer
	for _, c := range cs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(c.ArticleID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + c.UploadTime.Time().Format(_layoutSec) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.ViewCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.Income, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.TotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.TaxMoney, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + c.Date.Time().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.BaseIncome, 10))
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
