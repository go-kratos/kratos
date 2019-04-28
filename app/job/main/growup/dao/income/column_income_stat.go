package income

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/income"
	"go-common/library/log"
)

const (
	_columnIncomeStatSQL   = "SELECT id,aid,mid,upload_time,total_income,title FROM column_income_statis WHERE id > ? ORDER BY id LIMIT ?"
	_inColumnIncomeStatSQL = "INSERT INTO column_income_statis(aid,title,mid,tag_id,upload_time,total_income) VALUES %s ON DUPLICATE KEY UPDATE aid=VALUES(aid),title=VALUES(title),mid=VALUES(mid),upload_time=VALUES(upload_time),total_income=VALUES(total_income)"
)

// ColumnIncomeStat key: av_id
func (d *Dao) ColumnIncomeStat(c context.Context, id int64, limit int64) (m map[int64]*model.ColumnIncomeStat, last int64, err error) {
	rows, err := d.db.Query(c, _columnIncomeStatSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query AvIncomeStat error(%v)", err)
		return
	}
	defer rows.Close()
	m = make(map[int64]*model.ColumnIncomeStat)
	for rows.Next() {
		c := &model.ColumnIncomeStat{}
		err = rows.Scan(&last, &c.ArticleID, &c.MID, &c.UploadTime, &c.TotalIncome, &c.Title)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		m[c.ArticleID] = c
	}
	return
}

// InsertColumnIncomeStat batch insert column income stat
func (d *Dao) InsertColumnIncomeStat(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inColumnIncomeStatSQL, values))
	if err != nil {
		log.Error("d.db.Exec InsertColumnIncomeStat error(%v)", err)
		return
	}
	return res.RowsAffected()
}
