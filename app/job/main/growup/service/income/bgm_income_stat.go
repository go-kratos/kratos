package income

import (
	"bytes"
	"context"
	"strconv"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// BgmIncomeStatSvr bgm income stat service
type BgmIncomeStatSvr struct {
	batchSize int
	dao       *incomeD.Dao
}

// NewBgmIncomeStatSvr new bgm income stat svr
func NewBgmIncomeStatSvr(dao *incomeD.Dao, batchSize int) (svr *BgmIncomeStatSvr) {
	return &BgmIncomeStatSvr{
		batchSize: batchSize,
		dao:       dao,
	}
}

// BgmIncomeStat get bgm income stat
func (p *BgmIncomeStatSvr) BgmIncomeStat(c context.Context, limit int64) (m map[int64]*model.BgmIncomeStat, err error) {
	m = make(map[int64]*model.BgmIncomeStat)
	var id int64
	for {
		var bm map[int64]*model.BgmIncomeStat
		bm, id, err = p.dao.BgmIncomeStat(c, id, limit)
		if err != nil {
			return
		}
		if len(bm) == 0 {
			break
		}
		for sid, stat := range bm {
			m[sid] = stat
		}
	}
	return
}

// BatchInsertBgmIncomeStat batch insert bgm income stat
func (p *BgmIncomeStatSvr) BatchInsertBgmIncomeStat(c context.Context, bs map[int64]*model.BgmIncomeStat) (err error) {
	var (
		buff    = make([]*model.BgmIncomeStat, p.batchSize)
		buffEnd = 0
	)
	for _, b := range bs {
		if b.DataState == 0 {
			continue
		}
		buff[buffEnd] = b
		buffEnd++

		if buffEnd >= p.batchSize {
			values := bgmIncomeStatValues(buff[:buffEnd])
			buffEnd = 0
			_, err = p.dao.InsertBgmIncomeStat(c, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := bgmIncomeStatValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.dao.InsertBgmIncomeStat(c, values)
	}
	return
}

func bgmIncomeStatValues(bs []*model.BgmIncomeStat) (values string) {
	var buf bytes.Buffer
	for _, b := range bs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(b.SID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.TotalIncome, 10))
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
