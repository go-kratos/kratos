package income

import (
	"context"
	"fmt"
)

const (
	_inBGMIncomeSQL = "INSERT INTO bgm_income(aid,sid,mid,cid,income,total_income,tax_money,date,base_income,daily_total_income) VALUES %s ON DUPLICATE KEY UPDATE income=VALUES(income),total_income=VALUES(total_income),tax_money=VALUES(tax_money),base_income=VALUES(base_income),daily_total_income=VALUES(daily_total_income)"
)

// InsertBgmIncome batch insert bgm income
func (d *Dao) InsertBgmIncome(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inBGMIncomeSQL, values))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
