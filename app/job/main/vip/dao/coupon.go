package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_SelSalaryVideoCoupon = "SELECT `id`,`mid`,`coupon_count`,`coupon_type`,`state`,`type`,`ver` FROM `vip_view_coupon_salary_log_%s` WHERE `mid` = ?;"
	_addSalaryLogSQL      = "INSERT INTO `vip_view_coupon_salary_log_%s`(`mid`,`coupon_count`,`coupon_type`,`state`,`type`)VALUES(?,?,?,?,?);"
)

//SalaryVideoCouponList select salary video coupon list.
func (d *Dao) SalaryVideoCouponList(c context.Context, mid int64, dv string) (res []*model.VideoCouponSalaryLog, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_SelSalaryVideoCoupon, dv), mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VideoCouponSalaryLog)
		if err = rows.Scan(&r.ID, &r.Mid, &r.CouponCount, &r.CouponType, &r.State, &r.Type, &r.Ver); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//AddSalaryLog add salary log.
func (d *Dao) AddSalaryLog(c context.Context, l *model.VideoCouponSalaryLog, dv string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_addSalaryLogSQL, dv), l.Mid, l.CouponCount, l.CouponType, l.State, l.Type); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
