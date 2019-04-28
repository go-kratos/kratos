package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
	"go-common/library/xstr"
)

const (
	_usernameRoleSQL       = `SELECT u.id, u.username, coalesce(r.role,0) role FROM user u LEFT JOIN auth_role r ON u.id = r.uid WHERE u.id IN (%s)`
	_usernameDepartmentSQL = `SELECT u.id, u.username, coalesce(d.name,'') department FROM user u LEFT JOIN user_department d ON u.department_id = d.id WHERE u.id IN (%s)`
	_usernameSQL           = `SELECT id,username FROM user WHERE id IN (%s)`
)

//GetUsernameAndRole batch get username & role
func (d *Dao) GetUsernameAndRole(ctx context.Context, uids []int64) (list map[int64]*model.UserRole, err error) {
	var (
		rows *sql.Rows
	)
	list = map[int64]*model.UserRole{}
	uidStr := xstr.JoinInts(uids)
	if rows, err = d.mngDB.Query(ctx, fmt.Sprintf(_usernameRoleSQL, uidStr)); err != nil {
		PromeErr("mngdb: query", "GetUsernameAndRole d.mngDB.Query error(%v) uids(%s)", err, uidStr)
		return
	}
	defer rows.Close()

	for rows.Next() {
		u := new(model.UserRole)
		if err = rows.Scan(&u.UID, &u.Name, &u.Role); err != nil {
			PromeErr("mngdb: scan", "GetUsernameAndRole rows.Scan error(%v) uids(%s)", err, uidStr)
			return
		}

		list[u.UID] = u
	}

	return
}

//GetUsernameAndDepartment batch get username & department
func (d *Dao) GetUsernameAndDepartment(ctx context.Context, uids []int64) (list map[int64]*model.UserDepart, err error) {
	var (
		rows *sql.Rows
	)
	list = map[int64]*model.UserDepart{}
	uidStr := xstr.JoinInts(uids)
	if rows, err = d.mngDB.Query(ctx, fmt.Sprintf(_usernameDepartmentSQL, uidStr)); err != nil {
		PromeErr("mngdb: query", "GetUsernameAndDepartment d.mngDB.Query error(%v) uids(%s)", err, uidStr)
		return
	}
	defer rows.Close()

	for rows.Next() {
		u := new(model.UserDepart)
		if err = rows.Scan(&u.UID, &u.Name, &u.Department); err != nil {
			PromeErr("mngdb: scan", "GetUsernameAndDepartment rows.Scan error(%v) uids(%s)", err, uidStr)
			return
		}

		list[u.UID] = u
	}

	return
}

//GetUsername get username
func (d *Dao) GetUsername(ctx context.Context, uids []int64) (list map[int64]string, err error) {
	var (
		rows *sql.Rows
		uid  int64
		name string
	)
	list = map[int64]string{}
	uidStr := xstr.JoinInts(uids)
	if rows, err = d.mngDB.Query(ctx, fmt.Sprintf(_usernameSQL, uidStr)); err != nil {
		PromeErr("mngdb: query", "GetUsername d.mngDB.Query error(%v) uids(%s)", err, uidStr)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&uid, &name); err != nil {
			PromeErr("mngdb: scan", "GetUsername rows.Scan error(%v) uids(%s)", err, uidStr)
			return
		}

		list[uid] = name
	}

	return
}
