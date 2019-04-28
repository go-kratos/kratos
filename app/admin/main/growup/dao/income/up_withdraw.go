package income

import (
	"context"
	"fmt"

	model "go-common/app/admin/main/growup/model/income"

	"go-common/library/log"
)

const (
	// select
	_upWithdrawSQL = "SELECT id, mid, withdraw_income, date_version, mtime FROM up_income_withdraw WHERE id > ? AND state = 2 AND %s is_deleted = 0 LIMIT ?"
)

// ListUpWithdraw list up_income_withdraw by query
func (d *Dao) ListUpWithdraw(c context.Context, id int64, query string, limit int) (upWithdraw []*model.UpIncomeWithdraw, err error) {
	upWithdraw = make([]*model.UpIncomeWithdraw, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_upWithdrawSQL, query), id, limit)
	if err != nil {
		log.Error("GetUpWithdraw d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		w := &model.UpIncomeWithdraw{}
		err = rows.Scan(&w.ID, &w.MID, &w.WithdrawIncome, &w.DateVersion, &w.MTime)
		if err != nil {
			log.Error("GetUpWithdraw rows scan error(%v)", err)
			return
		}
		upWithdraw = append(upWithdraw, w)
	}
	err = rows.Err()
	return
}
