package income

import (
	"context"
	"fmt"
)

const (
	_inAvIncomeSQL = "INSERT INTO av_income(av_id,mid,tag_id,is_original,upload_time,play_count,total_income,income,tax_money,date,base_income) VALUES %s ON DUPLICATE KEY UPDATE av_id=VALUES(av_id),mid=VALUES(mid),tag_id=VALUES(tag_id),is_original=VALUES(is_original),upload_time=VALUES(upload_time),play_count=VALUES(play_count),total_income=VALUES(total_income),income=VALUES(income),tax_money=VALUES(tax_money),date=VALUES(date),base_income=VALUES(base_income)"
)

// InsertAvIncome batch insert av income
func (d *Dao) InsertAvIncome(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inAvIncomeSQL, values))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
