package dao

import (
	"context"

	model "go-common/app/admin/main/macross/model/manager"
	"go-common/library/log"
)

const (
	// load cache(get all).
	_authRelationSQL = `SELECT rid,auth_id FROM auth_relation`
	_authsSQL        = `SELECT id,system,name,flag,ctime,mtime FROM auth`
	// auth.
	_inAuthSQL  = `INSERT INTO auth (system,name,flag) VALUES(?,?,?)`
	_upAuthSQL  = `UPDATE auth SET name=? WHERE id=?`
	_delAuthSQL = `DELETE FROM auth WHERE id=?`
	// auth_relation.
	_inAuthRelationSQL   = `INSERT INTO auth_relation (rid,auth_id) VALUES(?,?)`
	_delAuthRelationSQL  = `DELETE FROM auth_relation WHERE rid=? AND auth_id=?`
	_cleanRelationByAuth = "DELETE FROM auth_relation WHERE auth_id=?"
)

// Auths select all auth from db.
func (d *Dao) Auths(c context.Context) (res map[string]map[int64]*model.Auth, err error) {
	rows, err := d.db.Query(c, _authsSQL)
	if err != nil {
		log.Error("Auths d.db.Query(%d) error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[string]map[int64]*model.Auth)
	for rows.Next() {
		var (
			auths map[int64]*model.Auth
			ok    bool
		)
		auth := &model.Auth{}
		if err = rows.Scan(&auth.AuthID, &auth.System, &auth.AuthName, &auth.AuthFlag, &auth.CTime, &auth.MTime); err != nil {
			log.Error("Auths rows.Scan error(%v)", err)
			return
		}
		if auths, ok = res[auth.System]; !ok {
			auths = make(map[int64]*model.Auth)
			res[auth.System] = auths
		}
		auths[auth.AuthID] = auth
	}
	return
}

// AddAuth insert auth.
func (d *Dao) AddAuth(c context.Context, system, authName, authFlag string) (rows int64, err error) {
	res, err := d.db.Exec(c, _inAuthSQL, system, authName, authFlag)
	if err != nil {
		log.Error("AddAuth d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// UpAuth update auth.
func (d *Dao) UpAuth(c context.Context, authName string, authID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _upAuthSQL, authName, authID)
	if err != nil {
		log.Error("UpAuth d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// DelAuth del auth.
func (d *Dao) DelAuth(c context.Context, authID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delAuthSQL, authID)
	if err != nil {
		log.Error("DelAuth d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// CleanAuthRelationByAuth del all auth relation by auth.
func (d *Dao) CleanAuthRelationByAuth(c context.Context, authID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _cleanRelationByAuth, authID)
	if err != nil {
		log.Error("CleanAuthRelationByAuth d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// AuthRelation select all auth_relation from db.
func (d *Dao) AuthRelation(c context.Context) (res map[int64][]int64, err error) {
	rows, err := d.db.Query(c, _authRelationSQL)
	if err != nil {
		log.Error("Roles d.db.Query(%d) error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64][]int64)
	for rows.Next() {
		var (
			roleID, authID int64
		)
		if err = rows.Scan(&roleID, &authID); err != nil {
			log.Error("Roles rows.Scan error(%v)", err)
			return
		}
		res[roleID] = append(res[roleID], authID)
	}
	return
}

// AddAuthRelation insert auth_relation.
func (d *Dao) AddAuthRelation(c context.Context, roleID, authID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _inAuthRelationSQL, roleID, authID)
	if err != nil {
		log.Error("AddAuthRelation d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// DelAuthRelation del auth_relation.
func (d *Dao) DelAuthRelation(c context.Context, roleID, authID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delAuthRelationSQL, roleID, authID)
	if err != nil {
		log.Error("DelAuthRelation d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
