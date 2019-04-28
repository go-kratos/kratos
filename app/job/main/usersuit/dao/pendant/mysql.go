package pendant

import (
	"context"
	"database/sql"

	"go-common/app/job/main/usersuit/model"
	xsql "go-common/library/database/sql"
)

const (
	_upEquipSQL        = "UPDATE user_pendant_equip SET pid = 0 AND expires = 0  WHERE mid = ?"
	_upEquipExpiresSQL = "UPDATE user_pendant_equip SET expires = ?  WHERE mid = ?"
	_selEquipSQL       = "SELECT mid FROM user_pendant_equip WHERE expires =< ?"
	_selEquipMIDSQL    = "SELECT mid,pid,expires FROM user_pendant_equip WHERE mid = ?"
	_selGidPidSQL      = "SELECT gid FROM pendant_group_ref WHERE pid = ?"
)

// UpEquipMID  update equip empty by mid
func (d *Dao) UpEquipMID(c context.Context, mid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upEquipSQL, mid); err != nil {
		return
	}
	return res.RowsAffected()
}

// UpEquipExpires  update equip expires by mid
func (d *Dao) UpEquipExpires(c context.Context, mid, expires int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upEquipExpiresSQL, expires, mid); err != nil {
		return
	}
	return res.RowsAffected()
}

// PendantEquipMID get user equip pendant by mid.
func (d *Dao) PendantEquipMID(c context.Context, mid int64) (pe *model.PendantEquip, err error) {
	row := d.db.QueryRow(c, _selEquipMIDSQL, mid)
	pe = new(model.PendantEquip)
	if err = row.Scan(&pe.Mid, &pe.Pid, &pe.Expires); err != nil {
		return
	}
	return
}

// ExpireEquipPendant  get expire equip pendant
func (d *Dao) ExpireEquipPendant(c context.Context, expires int64) (res []int64, err error) {
	var (
		row *xsql.Rows
		mid int64
	)
	if row, err = d.db.Query(c, _selEquipSQL, expires); err != nil {
		return
	}
	defer row.Close()
	for row.Next() {
		if err = row.Scan(&mid); err != nil {
			return
		}
		res = append(res, mid)
	}
	return
}

// PendantEquipGidPid get gid of its equip pendant by pid.
func (d *Dao) PendantEquipGidPid(c context.Context, pid int64) (gid int64, err error) {
	row := d.db.QueryRow(c, _selGidPidSQL, pid)
	if err = row.Scan(&gid); err != nil {
		return
	}
	return
}
