package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_countOrderByAppidSQL            = "SELECT COUNT(1) FROM vip_pay_order WHERE  mid = ? AND app_id = ? AND status=2;"
	_countAssociateGrantSQL          = "SELECT COUNT(1) FROM vip_order_associate_grant WHERE mid = ? AND app_id = ?;"
	_countAssociateByOrderNoSQL      = "SELECT COUNT(1) FROM vip_order_associate_grant WHERE out_trade_no = ? AND app_id = ?;"
	_countAssociateByMidAndMonthsSQL = "SELECT COUNT(1) FROM vip_order_associate_grant WHERE mid = ? AND app_id = ? AND months = ?;"
	_insertAssociateGrantOrderSQL    = "INSERT INTO vip_order_associate_grant(app_id,mid,months,out_open_id,out_trade_no,state,ctime)VALUES(?,?,?,?,?,?,?);"
	_updateAssociateGrantStateSQL    = "UPDATE vip_order_associate_grant SET state = ? WHERE out_trade_no = ? AND app_id = ?"
	_associateGrantCountInfoSQL      = "SELECT id,app_id,mid,months,current_count,ctime,mtime FROM vip_order_associate_grant_count WHERE mid = ? AND app_id = ? AND months = ?;"
	_addAssociateGrantCountSQL       = "INSERT IGNORE INTO vip_order_associate_grant_count(app_id,mid,months,current_count)VALUES(?,?,?,?);"
	_updateAssociateGrantCountSQL    = "UPDATE vip_order_associate_grant_count SET current_count = current_count + 1 WHERE mid = ? AND months = ? AND app_id = ? AND current_count = ?;"
)

// CountAssociateOrder associate order count.
func (d *Dao) CountAssociateOrder(c context.Context, mid int64, appID int64) (count int64, err error) {
	if err = d.olddb.QueryRow(c, _countOrderByAppidSQL, mid, appID).Scan(&count); err != nil {
		err = errors.Wrapf(err, "dao associate count orders(%d,%d)", mid, appID)
	}
	return
}

// CountAssociateGrants associate grant count.
func (d *Dao) CountAssociateGrants(c context.Context, mid int64, appID int64) (count int64, err error) {
	if err = d.olddb.QueryRow(c, _countAssociateGrantSQL, mid, appID).Scan(&count); err != nil {
		err = errors.Wrapf(err, "dao associate count grant(%d,%d)", mid, appID)
	}
	return
}

// CountGrantOrderByOutTradeNo grant order by out_trade_no count.
func (d *Dao) CountGrantOrderByOutTradeNo(c context.Context, outTradeNo string, appID int64) (count int64, err error) {
	if err = d.olddb.QueryRow(c, _countAssociateByOrderNoSQL, outTradeNo, appID).Scan(&count); err != nil {
		err = errors.Wrapf(err, "dao associate count grant by outno (%s,%d)", outTradeNo, appID)
	}
	return
}

// CountAssociateByMidAndMonths count grant order by mid and months .
func (d *Dao) CountAssociateByMidAndMonths(c context.Context, mid int64, appID int64, months int32) (count int64, err error) {
	if err = d.olddb.QueryRow(c, _countAssociateByMidAndMonthsSQL, mid, appID, months).Scan(&count); err != nil {
		err = errors.Wrapf(err, "dao associate count grant by mid and months (%d,%d,%d)", mid, appID, months)
	}
	return
}

// TxInsertAssociateGrantOrder insert associate grant order.
func (d *Dao) TxInsertAssociateGrantOrder(tx *sql.Tx, oa *model.VipOrderAssociateGrant) (aff int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_insertAssociateGrantOrderSQL, oa.AppID, oa.Mid, oa.Months, oa.OutOpenID, oa.OutTradeNO, oa.State, oa.Ctime); err != nil {
		err = errors.Wrapf(err, "dao insert grant(%+v)", oa)
		return
	}
	return res.RowsAffected()
}

// TxUpdateAssociateGrantState update associate grant state.
func (d *Dao) TxUpdateAssociateGrantState(tx *sql.Tx, oa *model.VipOrderAssociateGrant) (aff int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updateAssociateGrantStateSQL, oa.State, oa.OutTradeNO, oa.AppID); err != nil {
		err = errors.Wrapf(err, "dao update grant(%+v)", oa)
		return
	}
	return res.RowsAffected()
}

// AssociateGrantCountInfo associate grant count info.
func (d *Dao) AssociateGrantCountInfo(c context.Context, mid int64, appID int64, months int32) (res *model.VipAssociateGrantCount, err error) {
	res = new(model.VipAssociateGrantCount)
	if err = d.olddb.QueryRow(c, _associateGrantCountInfoSQL, mid, appID, months).
		Scan(&res.ID, &res.AppID, &res.Mid, &res.Months, &res.CurrentCount, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao associate openinfo openid(%d,%d,%d)", mid, appID, months)
	}
	return
}

// AddAssociateGrantCount add associate grant count.
func (d *Dao) AddAssociateGrantCount(c context.Context, arg *model.VipAssociateGrantCount) (err error) {
	if _, err = d.olddb.Exec(c, _addAssociateGrantCountSQL, arg.AppID, arg.Mid, arg.Months, arg.CurrentCount); err != nil {
		err = errors.Wrapf(err, "dao update bind state(%+v)", arg)
	}
	return
}

// UpdateAssociateGrantCount update associate grant count.
func (d *Dao) UpdateAssociateGrantCount(c context.Context, arg *model.VipAssociateGrantCount) (err error) {
	if _, err = d.olddb.Exec(c, _updateAssociateGrantCountSQL, arg.Mid, arg.Months, arg.AppID, arg.CurrentCount); err != nil {
		err = errors.Wrapf(err, "dao update associate grant count(%+v)", arg)
	}
	return
}
