package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_couponCodeSQL       = "SELECT id,batch_token,state,code,mid,coupon_type,ver FROM coupon_code WHERE code = ?;"
	_countCouponCountSQL = "SELECT COUNT(1) FROM coupon_code WHERE mid = ? AND batch_token = ?;"
	_updateCodeStateSQL  = "UPDATE coupon_code SET state = ?, mid = ?,coupon_token = ?, ver = ver + 1 WHERE code = ? AND ver = ? AND state = 1;"
)

// CouponCode get open info by code.
func (d *Dao) CouponCode(c context.Context, code string) (res *model.CouponCode, err error) {
	res = new(model.CouponCode)
	if err = d.db.QueryRow(c, _couponCodeSQL, code).
		Scan(&res.ID, &res.BatchToken, &res.State, &res.Code, &res.Mid, &res.CouponType, &res.Ver); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao coupon code(%s)", code)
	}
	return
}

// CountCodeByMid get count code by mid.
func (d *Dao) CountCodeByMid(c context.Context, mid int64, batckToken string) (count int64, err error) {
	if err = d.db.QueryRow(c, _countCouponCountSQL, mid, batckToken).Scan(&count); err != nil {
		err = errors.Wrapf(err, "dao count code")
	}
	return
}

// TxUpdateCodeState update code state.
func (d *Dao) TxUpdateCodeState(tx *sql.Tx, a *model.CouponCode) (aff int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updateCodeStateSQL, a.State, a.Mid, a.CouponToken, a.Code, a.Ver); err != nil {
		err = errors.Wrapf(err, "dao update code state(%+v)", a)
		return
	}
	return res.RowsAffected()
}
