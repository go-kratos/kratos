package dao

import (
	"context"
	xsql "database/sql"
	"time"

	"go-common/app/job/main/point/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_insertPoint           = "INSERT INTO point_info(mid,point_balance,ver) VALUES(?,?,?);"
	_updatePoint           = "UPDATE point_info SET point_balance=?,ver=? WHERE mid=? AND ver=?;"
	_insertPointHistory    = "INSERT INTO point_change_history(mid,point,order_id,change_type,change_time,relation_id,point_balance,remark,operator) VALUES(?,?,?,?,?,?,?,?,?);"
	_checkHistoryCount     = "SELECT COUNT(1) FROM  point_change_history WHERE mid = ? AND order_id = ?;"
	_insertPointSQL        = "INSERT INTO point_info(mid,point_balance,ver) VALUES(?,?,?)"
	_updatePointSQL        = "UPDATE point_info SET point_balance = ?,ver=? WHERE mid=? AND ver=?"
	_pointInfoSQL          = "SELECT mid,point_balance,ver FROM point_info WHERE mid=?"
	_InsertPointHistorySQL = "INSERT INTO point_change_history(mid,point,order_id,change_type,change_time,relation_id,point_balance,remark,operator) VALUES(?,?,?,?,?,?,?,?,?)"
	_midByMtime            = "SELECT mid, point_balance FROM point_info where mtime > ?;"
	_lastOneHistory        = "SELECT `point_balance` FROM `point_change_history` WHERE mid =?  ORDER  BY id DESC LIMIT 1;"
	_fixUpdatePointSQL     = "UPDATE point_info SET point_balance = ? WHERE mid=?"
)

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
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

//AddPoint addPoint
func (d *Dao) AddPoint(c context.Context, p *model.VipPoint) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _insertPoint, &p.Mid, &p.PointBalance, &p.Ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//UpdatePoint UpdatePoint row
func (d *Dao) UpdatePoint(c context.Context, p *model.VipPoint, oldver int64) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updatePoint, &p.PointBalance, &p.Ver, &p.Mid, oldver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//AddPointHistory add point history
func (d *Dao) AddPointHistory(c context.Context, ph *model.VipPointChangeHistory) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _insertPointHistory, &ph.Mid, &ph.Point, &ph.OrderID, &ph.ChangeType, &ph.ChangeTime, &ph.RelationID, &ph.PointBalance, &ph.Remark, &ph.Operator); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//HistoryCount check if have repeat record.
func (d *Dao) HistoryCount(c context.Context, mid int64, orderID string) (count int64, err error) {
	row := d.db.QueryRow(c, _checkHistoryCount, mid, orderID)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

//MidsByMtime point mids by mtime.
func (d *Dao) MidsByMtime(c context.Context, mtime time.Time) (pis []*model.PointInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _midByMtime, mtime); err != nil {
		err = errors.WithStack(err)
		return
	}
	for rows.Next() {
		pi := new(model.PointInfo)
		if err = rows.Scan(&pi.Mid, &pi.PointBalance); err != nil {
			pis = nil
			err = errors.WithStack(err)
			return
		}
		pis = append(pis, pi)
	}
	return
}

//LastOneHistory last one history.
func (d *Dao) LastOneHistory(c context.Context, mid int64) (point int64, err error) {
	row := d.db.QueryRow(c, _lastOneHistory, mid)
	if err = row.Scan(&point); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

//FixPointInfo fix point data .
func (d *Dao) FixPointInfo(c context.Context, mid int64, pointBalance int64) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _fixUpdatePointSQL, pointBalance, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
