package income

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/income"
	"go-common/library/log"
)

const (
	_bgmIncomeStatSQL   = "SELECT id,sid,total_income FROM bgm_income_statis WHERE id > ? ORDER BY id LIMIT ?"
	_inBgmIncomeStatSQL = "INSERT INTO bgm_income_statis(sid,total_income) VALUES %s ON DUPLICATE KEY UPDATE sid=VALUES(sid),total_income=VALUES(total_income)"
)

// BgmIncomeStat key: sid
func (d *Dao) BgmIncomeStat(c context.Context, id int64, limit int64) (m map[int64]*model.BgmIncomeStat, last int64, err error) {
	rows, err := d.db.Query(c, _bgmIncomeStatSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query BgmIncomeStat error(%v)", err)
		return
	}
	defer rows.Close()
	m = make(map[int64]*model.BgmIncomeStat)
	for rows.Next() {
		b := &model.BgmIncomeStat{}
		err = rows.Scan(&last, &b.SID, &b.TotalIncome)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		m[b.SID] = b
	}
	return
}

// InsertBgmIncomeStat batch insert bgm income stat
func (d *Dao) InsertBgmIncomeStat(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inBgmIncomeStatSQL, values))
	if err != nil {
		log.Error("d.db.Exec InsertBgmIncomeStat error(%v)", err)
		return
	}
	return res.RowsAffected()
}
