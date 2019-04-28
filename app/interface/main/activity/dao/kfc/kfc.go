package kfc

import (
	"context"
	"database/sql"
	"net/url"
	"strconv"

	"go-common/app/interface/main/activity/model/kfc"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// KfcCodeUsed kfcDao const.
var (
	KfcCodeUsed         = 1
	KfcCodeNotGiveOut   = 0
	_kfcCouponSQL       = "select `id`,`mid`,`ctime`,`mtime`,`coupon_code`,`desc`,`state`,`delete_time` from bnj_kfc_coupon where id = ?"
	_kfcCodeSQL         = "select `id`,`mid`,`ctime`,`mtime`,`coupon_code`,`desc`,`state`,`delete_time` from bnj_kfc_coupon where coupon_code = ?"
	_kfcCodeGiveOuteSQL = "update `bnj_kfc_coupon` set `state` = ? where `id` = ? and `state` = ?"
	_kfcDeliverSQL      = "update `bnj_kfc_coupon` set mid = ? where id = ? and mid = 0"
)

// RawKfcCoupon get coupon .
func (d *Dao) RawKfcCoupon(c context.Context, id int64) (res *kfc.BnjKfcCoupon, err error) {
	res = &kfc.BnjKfcCoupon{}
	row := d.db.QueryRow(c, _kfcCouponSQL, id)
	if err = row.Scan(&res.ID, &res.Mid, &res.Ctime, &res.Mtime, &res.CouponCode, &res.Desc, &res.State, &res.DeleteTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "RawKfcCoupon:row.Scan()")
		}
	}
	return
}

// RawKfcCode .
func (d *Dao) RawKfcCode(c context.Context, code string) (res *kfc.BnjKfcCoupon, err error) {
	res = &kfc.BnjKfcCoupon{}
	row := d.db.QueryRow(c, _kfcCodeSQL, code)
	if err = row.Scan(&res.ID, &res.Mid, &res.Ctime, &res.Mtime, &res.CouponCode, &res.Desc, &res.State, &res.DeleteTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "RawKfcCode:row.Scan()")
		}
	}
	return
}

// KfcCodeGiveOut .
func (d *Dao) KfcCodeGiveOut(c context.Context, id int64) (res int64, err error) {
	var (
		sqlRes sql.Result
	)
	if sqlRes, err = d.db.Exec(c, _kfcCodeGiveOuteSQL, KfcCodeUsed, id, KfcCodeNotGiveOut); err != nil {
		err = errors.Wrap(err, "d.db.Exec()")
		return
	}
	return sqlRes.RowsAffected()
}

// KfcDeliver .
func (d *Dao) KfcDeliver(c context.Context, id, mid int64) (res int64, err error) {
	var (
		sqlRes sql.Result
	)
	if sqlRes, err = d.db.Exec(c, _kfcDeliverSQL, mid, id); err != nil {
		err = errors.Wrap(err, "d.db.Exec()")
		return
	}
	return sqlRes.RowsAffected()
}

// KfcWinner .
func (d *Dao) KfcWinner(c context.Context, id int64) (uid int64, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(id, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			UID int64 `json:"uid"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.kfcWinnerURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		err = errors.Wrapf(err, "d.client.Get(%s)", d.kfcWinnerURL+"?"+params.Encode())
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		return
	}
	uid = res.Data.UID
	return
}
