package income

import (
	"context"
	"fmt"

	"go-common/library/log"

	model "go-common/app/job/main/growup/model/income"
)

const (
	_upIncomeStatSQL = "SELECT id,mid,total_income,av_total_income,column_total_income,bgm_total_income FROM up_income_statis WHERE id > ? ORDER BY id LIMIT ?"

	_inUpIncomeStatSQL = "INSERT INTO up_income_statis(mid,total_income,av_total_income,column_total_income,bgm_total_income) VALUES %s ON DUPLICATE KEY UPDATE mid=VALUES(mid),total_income=VALUES(total_income),av_total_income=VALUES(av_total_income),column_total_income=VALUES(column_total_income),bgm_total_income=VALUES(bgm_total_income)"

	_fixUpIncomeStatSQL = "INSERT INTO up_income_statis(mid,av_total_income,column_total_income,bgm_total_income) VALUES %s ON DUPLICATE KEY UPDATE mid=VALUES(mid),av_total_income=VALUES(av_total_income),column_total_income=VALUES(column_total_income),bgm_total_income=VALUES(bgm_total_income)"
)

// UpIncomeStat return m key: mid, value: total_income
func (d *Dao) UpIncomeStat(c context.Context, id int64, limit int64) (m map[int64]*model.UpIncomeStat, last int64, err error) {
	rows, err := d.db.Query(c, _upIncomeStatSQL, id, limit)
	if err != nil {
		log.Error("UpIncomeStat Query (%s, %d, %d) error(%v)", _upIncomeStatSQL, id, limit, err)
		return
	}

	defer rows.Close()
	m = make(map[int64]*model.UpIncomeStat)
	for rows.Next() {
		u := &model.UpIncomeStat{}
		err = rows.Scan(&last, &u.MID, &u.TotalIncome, &u.AvTotalIncome, &u.ColumnTotalIncome, &u.BgmTotalIncome)
		if err != nil {
			log.Error("UpIncomeStat rows scan error(%v)", err)
			return
		}
		m[u.MID] = u
	}
	return
}

// InsertUpIncomeStat batch insert up income stat
func (d *Dao) InsertUpIncomeStat(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpIncomeStatSQL, values))
	if err != nil {
		log.Error("d.db.Exec InsertUpIncomeStat error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// FixInsertUpIncomeStat fix insert up income stat
func (d *Dao) FixInsertUpIncomeStat(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_fixUpIncomeStatSQL, values))
	if err != nil {
		log.Error("d.db.Exec InsertUpIncomeStat error(%v)", err)
		return
	}
	return res.RowsAffected()
}
