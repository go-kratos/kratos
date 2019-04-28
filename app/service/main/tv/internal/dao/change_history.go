package dao

import (
	"context"
	"database/sql"

	"go-common/app/service/main/tv/internal/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_getUserChangeHistoryByID            = "SELECT `id`, `mid`, `change_type`, `change_time`, `order_no`, `days`, `operator_id`, `remark`, `ctime`, `mtime` FROM `tv_user_change_history` WHERE `id`=?"
	_getUserChangeHistorysByMid          = "SELECT `id`, `mid`, `change_type`, `change_time`, `order_no`, `days`, `operator_id`, `remark`, `ctime`, `mtime` FROM `tv_user_change_history` WHERE `mid`=? ORDER BY `ctime` DESC LIMIT ?,?"
	_countUserChangeHistoryByMid         = "SELECT count(*) FROM `tv_user_change_history` WHERE `mid`=?"
	_getUserChangeHistorysByMidAndCtime  = "SELECT `id`, `mid`, `change_type`, `change_time`, `order_no`, `days`, `operator_id`, `remark`, `ctime`, `mtime` FROM `tv_user_change_history` WHERE `mid`=? AND `ctime`>=? AND `ctime`<? ORDER BY `ctime` DESC LIMIT ?,?"
	_countUserChangeHistoryByMidAndCtime = "SELECT count(*) FROM `tv_user_change_history` WHERE `mid`=? AND `ctime`>=? AND `ctime`<?"

	_insertUserChangeHistory = "INSERT INTO tv_user_change_history (`mid`, `change_type`, `change_time`, `order_no`, `days`, `operator_id`, `remark`) VALUES (?,?,?,?,?,?,?)"
)

// UserChangeHistoryByID quires one row from tv_user_change_history.
func (d *Dao) UserChangeHistoryByID(c context.Context, id int32) (uch *model.UserChangeHistory, err error) {
	row := d.db.QueryRow(c, _getUserChangeHistoryByID, id)
	uch = &model.UserChangeHistory{}
	err = row.Scan(&uch.ID, &uch.Mid, &uch.ChangeType, &uch.ChangeTime, &uch.OrderNo, &uch.Days, &uch.OperatorId, &uch.Remark, &uch.Ctime, &uch.Mtime)
	if err != nil {
		log.Error("rows.Scan(%s) error(%v)", _getUserChangeHistoryByID, err)
		err = errors.WithStack(err)
		return nil, err
	}
	return uch, nil
}

// UserChangeHistorysByMid quires rows from tv_user_change_history.
func (d *Dao) UserChangeHistorysByMid(c context.Context, mid int64, pn, ps int32) (res []*model.UserChangeHistory, total int, err error) {
	res = make([]*model.UserChangeHistory, 0)
	totalRow := d.db.QueryRow(c, _countUserChangeHistoryByMid, mid)
	if err = totalRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		err = errors.WithStack(err)
		return
	}
	rows, err := d.db.Query(c, _getUserChangeHistorysByMid, mid, (pn-1)*ps, ps)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _getUserChangeHistorysByMid, err)
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		uch := &model.UserChangeHistory{}
		if err = rows.Scan(&uch.ID, &uch.Mid, &uch.ChangeType, &uch.ChangeTime, &uch.OrderNo, &uch.Days, &uch.OperatorId, &uch.Remark, &uch.Ctime, &uch.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			err = errors.WithStack(err)
			return
		}
		res = append(res, uch)
	}
	return
}

// UserChangeHistorysByMidAndCtime quires rows from tv_user_change_history.
func (d *Dao) UserChangeHistorysByMidAndCtime(c context.Context, mid int64, from, to xtime.Time, pn, ps int32) (res []*model.UserChangeHistory, total int, err error) {
	res = make([]*model.UserChangeHistory, 0)
	totalRow := d.db.QueryRow(c, _countUserChangeHistoryByMidAndCtime, mid, from, to)
	if err = totalRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		err = errors.WithStack(err)
		return
	}
	rows, err := d.db.Query(c, _getUserChangeHistorysByMidAndCtime, mid, from, to, (pn-1)*ps, ps)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _getUserChangeHistorysByMidAndCtime, err)
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		uch := &model.UserChangeHistory{}
		if err = rows.Scan(&uch.ID, &uch.Mid, &uch.ChangeType, &uch.ChangeTime, &uch.OrderNo, &uch.Days, &uch.OperatorId, &uch.Remark, &uch.Ctime, &uch.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			err = errors.WithStack(err)
			return
		}
		res = append(res, uch)
	}
	return
}

// TxInsertUserChangeHistory insert one row into tv_user_change_history.
func (d *Dao) TxInsertUserChangeHistory(ctx context.Context, tx *xsql.Tx, uch *model.UserChangeHistory) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = tx.Exec(_insertUserChangeHistory, uch.Mid, uch.ChangeType, uch.ChangeTime, uch.OrderNo, uch.Days, uch.OperatorId, uch.Remark); err != nil {
		log.Error("tx.Exec(%s) error(%v)", _insertUserChangeHistory, err)
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		log.Error("res.LastInsertId(%s) error(%v)", _insertUserChangeHistory, err)
		err = errors.WithStack(err)
		return
	}
	return
}
