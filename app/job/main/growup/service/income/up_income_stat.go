package income

import (
	"bytes"
	"context"
	"strconv"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// UpIncomeStatSvr up_income_stat svr
type UpIncomeStatSvr struct {
	batchSize int
	dao       *incomeD.Dao
}

// NewUpIncomeStatSvr new server
func NewUpIncomeStatSvr(dao *incomeD.Dao, batchSize int) (svr *UpIncomeStatSvr) {
	return &UpIncomeStatSvr{
		batchSize: batchSize,
		dao:       dao,
	}
}

// UpIncomeStat return stats, key: mid, value: total_income
func (p *UpIncomeStatSvr) UpIncomeStat(c context.Context, limit int64) (m map[int64]*model.UpIncomeStat, err error) {
	var id int64
	m = make(map[int64]*model.UpIncomeStat)
	for {
		var um map[int64]*model.UpIncomeStat
		um, id, err = p.dao.UpIncomeStat(c, id, limit)
		if err != nil {
			return
		}
		if len(um) == 0 {
			break
		}
		for mid, u := range um {
			if u.IsDeleted == 0 {
				m[mid] = u
			}
		}
	}
	return
}

// BatchInsertUpIncomeStat insert up_income_statis batch
func (p *UpIncomeStatSvr) BatchInsertUpIncomeStat(c context.Context, us map[int64]*model.UpIncomeStat) (err error) {
	var (
		buff    = make([]*model.UpIncomeStat, batchSize)
		buffEnd = 0
	)
	for _, u := range us {
		if u.DataState == 0 {
			continue
		}
		buff[buffEnd] = u
		buffEnd++
		if buffEnd >= p.batchSize {
			values := upIncomeStatValues(buff[:buffEnd])
			buffEnd = 0
			_, err = p.dao.InsertUpIncomeStat(c, values)
			if err != nil {
				return
			}
		}
	}

	if buffEnd > 0 {
		values := upIncomeStatValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.dao.InsertUpIncomeStat(c, values)
	}
	return
}

func upIncomeStatValues(us []*model.UpIncomeStat) (values string) {
	var buf bytes.Buffer
	for _, u := range us {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(u.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.TotalIncome, 10))
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
	values = buf.String()
	buf.Reset()
	return
}
