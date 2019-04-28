package income

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/income"
	"go-common/library/log"
)

const (
	_avIncomeStatSQL   = "SELECT id,av_id,mid,tag_id,is_original,upload_time,total_income,ctime FROM av_income_statis WHERE id > ? ORDER BY id LIMIT ?"
	_inAvIncomeStatSQL = "INSERT INTO av_income_statis(av_id,mid,tag_id,is_original,upload_time,total_income) VALUES %s ON DUPLICATE KEY UPDATE av_id=VALUES(av_id),mid=VALUES(mid),tag_id=VALUES(tag_id),is_original=VALUES(is_original),upload_time=VALUES(upload_time),total_income=VALUES(total_income)"
)

// AvIncomeStat key: av_id
func (d *Dao) AvIncomeStat(c context.Context, id int64, limit int64) (m map[int64]*model.AvIncomeStat, last int64, err error) {
	rows, err := d.db.Query(c, _avIncomeStatSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query AvIncomeStat error(%v)", err)
		return
	}
	defer rows.Close()
	m = make(map[int64]*model.AvIncomeStat)
	for rows.Next() {
		a := &model.AvIncomeStat{}
		err = rows.Scan(&last, &a.AvID, &a.MID, &a.TagID, &a.IsOriginal, &a.UploadTime, &a.TotalIncome, &a.CTime)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		m[a.AvID] = a
	}
	return
}

// InsertAvIncomeStat batch insert av income stat
func (d *Dao) InsertAvIncomeStat(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inAvIncomeStatSQL, values))
	if err != nil {
		log.Error("d.db.Exec InsertAvIncomeStat error(%v)", err)
		return
	}
	return res.RowsAffected()
}
