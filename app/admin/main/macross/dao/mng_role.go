package dao

import (
	"context"

	model "go-common/app/admin/main/macross/model/manager"
	"go-common/library/log"
)

const (
	// load cache(get all).
	_rolesSQL = `SELECT id,system,name,ctime,mtime FROM role`
	// role.
	_inRoleSQL           = `INSERT INTO role (system,name) VALUES(?,?)`
	_upRoleSQL           = `UPDATE role SET name=? WHERE id=?`
	_delRoleSQL          = `DELETE FROM role WHERE id=?`
	_cleanRelationByRole = "DELETE FROM auth_relation WHERE rid=?"
)

// Roles select all role from db.
func (d *Dao) Roles(c context.Context) (res map[string]map[int64]*model.Role, err error) {
	rows, err := d.db.Query(c, _rolesSQL)
	if err != nil {
		log.Error("Roles d.db.Query(%d) error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[string]map[int64]*model.Role)
	for rows.Next() {
		var (
			roles map[int64]*model.Role
			ok    bool
		)
		role := &model.Role{}
		if err = rows.Scan(&role.RoleID, &role.System, &role.RoleName, &role.CTime, &role.MTime); err != nil {
			log.Error("Roles rows.Scan error(%v)", err)
			return
		}
		if roles, ok = res[role.System]; !ok {
			roles = make(map[int64]*model.Role)
			res[role.System] = roles
		}
		roles[role.RoleID] = role
	}
	return
}

// AddRole insert role.
func (d *Dao) AddRole(c context.Context, system, roleName string) (rows int64, err error) {
	res, err := d.db.Exec(c, _inRoleSQL, system, roleName)
	if err != nil {
		log.Error("AddRole d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// UpRole update role.
func (d *Dao) UpRole(c context.Context, roleName string, roleID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _upRoleSQL, roleName, roleID)
	if err != nil {
		log.Error("UpRole d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// DelRole del role.
func (d *Dao) DelRole(c context.Context, roleID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delRoleSQL, roleID)
	if err != nil {
		log.Error("DelRole d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// CleanAuthRelationByRole del all auth relation by role.
func (d *Dao) CleanAuthRelationByRole(c context.Context, roleID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _cleanRelationByRole, roleID)
	if err != nil {
		log.Error("CleanAuthRelationByRole d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
