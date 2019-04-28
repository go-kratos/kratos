package dao

import (
	"context"

	model "go-common/app/admin/main/macross/model/manager"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// load cache(get all).
	_usersSQL = `SELECT user.id,user.system,user.name,user.rid,role.name,user.ctime,user.mtime FROM user,role WHERE user.rid=role.id`
	_userSQL  = `SELECT user.id,user.system,user.name,user.rid,role.name,user.ctime,user.mtime FROM user,role WHERE user.rid=role.id AND user.id=?`
	// user.
	_inUserSQL  = `INSERT INTO user (system,name,rid) VALUES(?,?,?)`
	_upUserSQL  = `UPDATE user SET name=?,rid=? WHERE id=?`
	_delUserSQL = `DELETE FROM user WHERE id=?`
)

// Users select all user from db.
func (d *Dao) Users(c context.Context) (res map[string]map[string]*model.User, err error) {
	rows, err := d.db.Query(c, _usersSQL)
	if err != nil {
		log.Error("UserAll d.db.Query(%d) error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[string]map[string]*model.User)
	for rows.Next() {
		var (
			users map[string]*model.User
			ok    bool
		)
		user := &model.User{}
		if err = rows.Scan(&user.UserID, &user.System, &user.UserName, &user.RoleID, &user.RoleName, &user.CTime, &user.MTime); err != nil {
			log.Error("Users rows.Scan error(%v)", err)
			return
		}
		if users, ok = res[user.System]; !ok {
			users = make(map[string]*model.User)
			res[user.System] = users
		}
		users[user.UserName] = user
	}
	return
}

// User get user.
func (d *Dao) User(c context.Context, userID int64) (re *model.User, err error) {
	row := d.db.QueryRow(c, _userSQL, userID)
	re = &model.User{}
	if err = row.Scan(&re.UserID, &re.System, &re.UserName, &re.RoleID, &re.RoleName, &re.CTime, &re.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("User d.db.QueryRow(%d) error(%v)", userID, err)
		}
	}
	return
}

// AddUser insert user.
func (d *Dao) AddUser(c context.Context, roleID int64, system, userName string) (rows int64, err error) {
	res, err := d.db.Exec(c, _inUserSQL, system, userName, roleID)
	if err != nil {
		log.Error("AddUser d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// UpUser update user.
func (d *Dao) UpUser(c context.Context, userID, roleID int64, userName string) (rows int64, err error) {
	res, err := d.db.Exec(c, _upUserSQL, userName, roleID, userID)
	if err != nil {
		log.Error("UpUser d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// DelUser del user.
func (d *Dao) DelUser(c context.Context, userID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delUserSQL, userID)
	if err != nil {
		log.Error("DelUser d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
