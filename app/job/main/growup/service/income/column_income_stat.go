package income

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// ColumnIncomeStatSvr column income statistics service
type ColumnIncomeStatSvr struct {
	batchSize int
	dao       *incomeD.Dao
}

// NewColumnIncomeStatSvr new column income stat service
func NewColumnIncomeStatSvr(dao *incomeD.Dao, batchSize int) (svr *ColumnIncomeStatSvr) {
	return &ColumnIncomeStatSvr{
		batchSize: batchSize,
		dao:       dao,
	}
}

// ColumnIncomeStat column income statistisc
func (p *ColumnIncomeStatSvr) ColumnIncomeStat(c context.Context, limit int64) (m map[int64]*model.ColumnIncomeStat, err error) {
	m = make(map[int64]*model.ColumnIncomeStat)
	var id int64
	for {
		var cm map[int64]*model.ColumnIncomeStat
		cm, id, err = p.dao.ColumnIncomeStat(c, id, limit)
		if err != nil {
			return
		}
		if len(cm) == 0 {
			break
		}
		for sid, stat := range cm {
			m[sid] = stat
		}
	}
	return
}

// BatchInsertColumnIncomeStat batch insert column income stat
func (p *ColumnIncomeStatSvr) BatchInsertColumnIncomeStat(ctx context.Context, cs map[int64]*model.ColumnIncomeStat) (err error) {
	var (
		buff    = make([]*model.ColumnIncomeStat, p.batchSize)
		buffEnd = 0
	)
	for _, c := range cs {
		if c.DataState == 0 {
			continue
		}
		buff[buffEnd] = c
		buffEnd++

		if buffEnd >= p.batchSize {
			values := columnIncomeStatValues(buff[:buffEnd])
			buffEnd = 0
			_, err = p.dao.InsertColumnIncomeStat(ctx, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := columnIncomeStatValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.dao.InsertColumnIncomeStat(ctx, values)
	}
	return
}

func columnIncomeStatValues(cs []*model.ColumnIncomeStat) (values string) {
	var buf bytes.Buffer
	for _, c := range cs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(c.ArticleID, 10))
		buf.WriteByte(',')
		buf.WriteString("\"" + strings.Replace(c.Title, "\"", "\\\"", -1) + "\"")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + c.UploadTime.Time().Format(_layoutSec) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(c.TotalIncome, 10))
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
