package dao

import (
	"context"
	"fmt"
	"go-common/app/job/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_selSalaryMaxID   = "SELECT IFNULL(MAX(id),0) id FROM vip_view_coupon_salary_log_%s;"
	_selOldSalaryList = "SELECT `mid`,`coupon_count`,`state`,`type` FROM `vip_view_coupon_salary_log_%s` WHERE id>? AND id <=?;"
)

// SalaryLogMaxID select salary log max id.
func (d *Dao) SalaryLogMaxID(c context.Context, dv string) (maxID int, err error) {
	var row = d.oldDb.QueryRow(c, fmt.Sprintf(_selSalaryMaxID, dv))
	if err = row.Scan(&maxID); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//SelOldSalaryList sel old salary list
func (d *Dao) SelOldSalaryList(c context.Context, id, endID int, dv string) (res []*model.OldSalaryLog, err error) {
	var rows *sql.Rows
	if rows, err = d.oldDb.Query(c, fmt.Sprintf(_selOldSalaryList, dv), id, endID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.OldSalaryLog)
		if err = rows.Scan(&r.Mid, &r.CouponCount, &r.State, &r.Type); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
