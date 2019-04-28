package income

import (
	"bytes"
	"context"
	"strconv"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// BgmIncomeSvr bgm income service
type BgmIncomeSvr struct {
	batchSize int
	dao       *incomeD.Dao
}

// NewBgmIncomeSvr new bgm income service
func NewBgmIncomeSvr(dao *incomeD.Dao, batchSize int) (svr *BgmIncomeSvr) {
	return &BgmIncomeSvr{
		batchSize: batchSize,
		dao:       dao,
	}
}

// BatchInsertBgmIncome batch insert bgm income
func (p *BgmIncomeSvr) BatchInsertBgmIncome(ctx context.Context, bm map[int64]map[int64]map[int64]*model.BgmIncome) (err error) {
	var (
		buff    = make([]*model.BgmIncome, p.batchSize)
		buffEnd = 0
	)

	for _, sm := range bm {
		for _, bs := range sm {
			for _, b := range bs {
				buff[buffEnd] = b
				buffEnd++

				if buffEnd >= p.batchSize {
					values := bgmIncomeValues(buff[:buffEnd])
					buffEnd = 0
					_, err = p.dao.InsertBgmIncome(ctx, values)
					if err != nil {
						return
					}
				}
			}
		}
	}
	if buffEnd > 0 {
		values := bgmIncomeValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.dao.InsertBgmIncome(ctx, values)
	}
	return
}

func bgmIncomeValues(bs []*model.BgmIncome) (values string) {
	var buf bytes.Buffer
	for _, b := range bs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(b.AID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.SID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.CID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.Income, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.TotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.TaxMoney, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + b.Date.Time().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.BaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.DailyTotalIncome, 10))
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
