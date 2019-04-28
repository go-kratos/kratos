package dao

import (
	"context"
	"database/sql"

	"go-common/app/job/main/passport-user/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addUserTelDuplicateSQL        = "INSERT INTO user_tel_duplicate (mid,tel,cid,tel_bind_time,ts) VALUES (?,?,?,?,?)"
	_addUserEmailDuplicateSQL      = "INSERT INTO user_email_duplicate (mid,email,verified,email_bind_time,ts) VALUES (?,?,?,?,?)"
	_getUserTelDuplicateSQL        = "SELECT id,mid,tel,cid,tel_bind_time,status,ts FROM user_tel_duplicate WHERE status = 0 order by ts"
	_getUserEmailDuplicateSQL      = "SELECT id,mid,email,verified,email_bind_time,status,ts FROM user_email_duplicate WHERE status = 0 order by ts"
	_updateTelDuplicateStatusSQL   = "UPDATE user_tel_duplicate SET status = 1 WHERE id = ?"
	_updateEmailDuplicateStatusSQL = "UPDATE user_email_duplicate SET status = 1 WHERE id = ?"
)

// AddUserTelDuplicate add user tel duplicate.
func (d *Dao) AddUserTelDuplicate(c context.Context, a *model.UserTelDuplicate) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _addUserTelDuplicateSQL, a.Mid, a.Tel, a.Cid, a.TelBindTime, a.Timestamp); err != nil {
		log.Error("fail to add user tel duplicate, userTelDuplicate(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// AddUserEmailDuplicate add user email duplicate.
func (d *Dao) AddUserEmailDuplicate(c context.Context, a *model.UserEmailDuplicate) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _addUserEmailDuplicateSQL, a.Mid, a.Email, a.Verified, a.EmailBindTime, a.Timestamp); err != nil {
		log.Error("fail to add user email duplicate, userEmailDuplicate(%+v) dao.userDB.Exec() error(%+v)", a, err)
		return
	}
	return res.RowsAffected()
}

// UserTelDuplicate get user tel duplicate.
func (d *Dao) UserTelDuplicate(c context.Context) (res []*model.UserTelDuplicate, err error) {
	var rows *xsql.Rows
	if rows, err = d.userDB.Query(c, _getUserTelDuplicateSQL); err != nil {
		log.Error("fail to get UserTelDuplicate, dao.userDB.Query(%s) error(%v)", _getUserTelDuplicateSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.UserTelDuplicate)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Tel, &r.Cid, &r.TelBindTime, &r.Status, &r.Timestamp); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// UserEmailDuplicate get user email duplicate.
func (d *Dao) UserEmailDuplicate(c context.Context) (res []*model.UserEmailDuplicate, err error) {
	var rows *xsql.Rows
	if rows, err = d.userDB.Query(c, _getUserEmailDuplicateSQL); err != nil {
		log.Error("fail to get UserEmailDuplicate, dao.userDB.Query(%s) error(%v)", _getUserEmailDuplicateSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.UserEmailDuplicate)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Email, &r.Verified, &r.EmailBindTime, &r.Status, &r.Timestamp); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// UpdateUserTelDuplicateStatus update user tel duplicate status.
func (d *Dao) UpdateUserTelDuplicateStatus(c context.Context, id int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _updateTelDuplicateStatusSQL, id); err != nil {
		log.Error("fail to update tel duplicate status, id(%d) dao.userDB.Exec() error(%+v)", id, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserEmailDuplicateStatus update user email duplicate status.
func (d *Dao) UpdateUserEmailDuplicateStatus(c context.Context, id int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.userDB.Exec(c, _updateEmailDuplicateStatusSQL, id); err != nil {
		log.Error("fail to update email duplicate status, id(%d) dao.userDB.Exec() error(%+v)", id, err)
		return
	}
	return res.RowsAffected()
}
