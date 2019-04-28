package income

import (
	"context"
	"fmt"
)

const (
	_inColumnIncomeSQL = "INSERT INTO column_income(aid,mid,tag_id,upload_time,view_count,income,total_income,tax_money,date,base_income) VALUES %s ON DUPLICATE KEY UPDATE income=VALUES(income),total_income=VALUES(total_income),base_income=VALUES(base_income)"
)

// InsertColumnIncome batch insert column income
func (d *Dao) InsertColumnIncome(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inColumnIncomeSQL, values))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
