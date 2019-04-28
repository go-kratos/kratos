package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/service/main/point/model"
	"go-common/library/database/sql"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	// point sql
	_pointInfoSQL          = "SELECT mid,point_balance,ver FROM point_info WHERE mid=?"
	_pointHistorySQL       = "SELECT id,mid,point,order_id,change_type,change_time,relation_id,point_balance,remark,operator FROM point_change_history WHERE mid = ? AND id < ?  ORDER BY id DESC LIMIT ?"
	_pointHistoryCountSQL  = "SELECT COUNT(1) FROM point_change_history WHERE mid = ?"
	_updatePointSQL        = "UPDATE point_info SET point_balance = ?,ver=? WHERE mid=? AND ver=?"
	_insertPointSQL        = "INSERT INTO point_info(mid,point_balance,ver) VALUES(?,?,?)"
	_InsertPointHistorySQL = "INSERT INTO point_change_history(mid,point,order_id,change_type,change_time,relation_id,point_balance,remark,operator) VALUES(?,?,?,?,?,?,?,?,?)"
	_pointHistoryCheckSQL  = "SELECT id FROM point_change_history WHERE order_id=?"
	_selPointHistorySQL    = "SELECT id,mid,point,order_id,change_type,change_time,relation_id,point_balance,remark,operator FROM point_change_history WHERE mid = ? AND change_time>=? AND change_time <= ? "
	_allPointConfig        = "SELECT `id`,`app_id`,`point`,`operator`,`ctime`,`mtime` FROM `point_conf`;"
	//TODO compatible vip point old gateway
	_oldPointHistorySQL = "SELECT id,mid,point,order_id,change_type,change_time,relation_id,point_balance,remark,operator FROM point_change_history WHERE mid = ? ORDER BY id DESC LIMIT ?,?;"
)

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}

// PointInfo .
func (d *Dao) PointInfo(c context.Context, mid int64) (pi *model.PointInfo, err error) {
	row := d.db.QueryRow(c, _pointInfoSQL, mid)
	pi = new(model.PointInfo)
	if err = row.Scan(&pi.Mid, &pi.PointBalance, &pi.Ver); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			pi = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

//TxPointInfo .
func (d *Dao) TxPointInfo(c context.Context, tx *sql.Tx, mid int64) (pi *model.PointInfo, err error) {
	row := tx.QueryRow(_pointInfoSQL, mid)
	pi = new(model.PointInfo)
	if err = row.Scan(&pi.Mid, &pi.PointBalance, &pi.Ver); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			pi = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

//PointHistory point history
func (d *Dao) PointHistory(c context.Context, mid int64, cursor int, ps int) (phs []*model.PointHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _pointHistorySQL, mid, cursor, ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	for rows.Next() {
		ph := new(model.PointHistory)
		if err = rows.Scan(&ph.ID, &ph.Mid, &ph.Point, &ph.OrderID, &ph.ChangeType, &ph.ChangeTime, &ph.RelationID, &ph.PointBalance, &ph.Remark, &ph.Operator); err != nil {
			phs = nil
			err = errors.WithStack(err)
			return
		}
		phs = append(phs, ph)
	}
	return
}

//PointHistoryCount point history
func (d *Dao) PointHistoryCount(c context.Context, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, _pointHistoryCountSQL, mid)
	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//UpdatePointInfo .
func (d *Dao) UpdatePointInfo(c context.Context, tx *sql.Tx, pi *model.PointInfo, ver int64) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updatePointSQL, pi.PointBalance, pi.Ver, pi.Mid, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//InsertPoint .
func (d *Dao) InsertPoint(c context.Context, tx *sql.Tx, pi *model.PointInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_insertPointSQL, pi.Mid, pi.PointBalance, pi.Ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//InsertPointHistory .
func (d *Dao) InsertPointHistory(c context.Context, tx *sql.Tx, ph *model.PointHistory) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_InsertPointHistorySQL, ph.Mid, ph.Point, ph.OrderID, ph.ChangeType, ph.ChangeTime, ph.RelationID, ph.PointBalance, ph.Remark, ph.Operator); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//SelPointHistory .
func (d *Dao) SelPointHistory(c context.Context, mid int64, startDate, endDate xtime.Time) (phs []*model.PointHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selPointHistorySQL, mid, startDate, endDate); err != nil {
		err = errors.WithStack(err)
		return
	}
	for rows.Next() {
		ph := new(model.PointHistory)
		if err = rows.Scan(&ph.ID, &ph.Mid, &ph.Point, &ph.OrderID, &ph.ChangeType, &ph.ChangeTime, &ph.RelationID, &ph.PointBalance, &ph.Remark, &ph.Operator); err != nil {
			phs = nil
			err = errors.WithStack(err)
		}
		phs = append(phs, ph)
	}
	return
}

// ExistPointOrder check orderID is uniq or not
func (d *Dao) ExistPointOrder(c context.Context, orID string) (id int, err error) {
	row := d.db.QueryRow(c, _pointHistoryCheckSQL, orID)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// AllPointConfig all point config.
func (d *Dao) AllPointConfig(c context.Context) (res []*model.VipPointConf, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allPointConfig); err != nil {
		err = errors.WithStack(err)
		return
	}
	for rows.Next() {
		vf := new(model.VipPointConf)
		if err = rows.Scan(&vf.ID, &vf.AppID, &vf.Point, &vf.Operator, &vf.Ctime, &vf.Mtime); err != nil {
			res = nil
			err = errors.WithStack(err)
		}
		res = append(res, vf)
	}
	return
}

//OldPointHistory point history.
func (d *Dao) OldPointHistory(c context.Context, mid int64, start int, ps int) (phs []*model.OldPointHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _oldPointHistorySQL, mid, start, ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	for rows.Next() {
		ph := new(model.PointHistory)
		if err = rows.Scan(&ph.ID, &ph.Mid, &ph.Point, &ph.OrderID, &ph.ChangeType, &ph.ChangeTime, &ph.RelationID, &ph.PointBalance, &ph.Remark, &ph.Operator); err != nil {
			phs = nil
			err = errors.WithStack(err)
			return
		}
		oph := new(model.OldPointHistory)
		oph.ID = ph.ID
		oph.Mid = ph.Mid
		oph.Point = ph.Point
		oph.OrderID = ph.OrderID
		oph.ChangeType = ph.ChangeType
		oph.ChangeTime = ph.ChangeTime.Time().Unix()
		oph.RelationID = ph.RelationID
		oph.PointBalance = ph.PointBalance
		oph.Remark = ph.Remark
		oph.Operator = ph.Operator
		phs = append(phs, oph)
	}
	return
}
