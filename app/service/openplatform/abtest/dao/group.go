package dao

import (
	"context"
	xsql "database/sql"
	"fmt"

	"go-common/app/service/openplatform/abtest/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_addGroup    = "INSERT INTO abtest_group (`name`,`desc`,`is_deleted`)VALUES(?, ?, 0)"
	_listGroup   = "SELECT `id`, `name`, `desc` FROM abtest_group WHERE `is_deleted` = 0 ORDER BY id"
	_updateGroup = "UPDATE abtest_group set `name` = ?, `desc` = ? where id = ?"
	_deleteGroup = "UPDATE abtest_group set `is_deleted` = 1 WHERE id = ?"
)

//AddGroup add a new group
func (d *Dao) AddGroup(c context.Context, g model.Group) (i int, err error) {
	var (
		res xsql.Result
	)
	if res, err = d.db.Exec(c, _addGroup, g.Name, g.Desc); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.group|AddGroup] d.db.Exec err: %v", err))
		return
	}
	return intConv(res.LastInsertId())
}

//UpdateGroup update group by id
func (d *Dao) UpdateGroup(c context.Context, g model.Group) (i int, err error) {
	var (
		res xsql.Result
	)
	if res, err = d.db.Exec(c, _updateGroup, g.Name, g.Desc, g.ID); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.group|UpdateGroup] d.db.Exec err: %v", err))
		return
	}
	return intConv(res.RowsAffected())
}

//DeleteGroup delete the group by id
func (d *Dao) DeleteGroup(c context.Context, id int) (i int, err error) {
	var (
		res xsql.Result
	)
	if res, err = d.db.Exec(c, _deleteGroup, id); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.group|DeleteGroup] d.db.Exec err: %v", err))
		return
	}
	return intConv(res.RowsAffected())
}

//ListGroup list all groups
func (d *Dao) ListGroup(c context.Context) (res []*model.Group, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _listGroup); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.group|ListGroup] d.db.Query err: %v", err))
		return
	}

	defer rows.Close()
	for rows.Next() {
		g := &model.Group{}
		if err = rows.Scan(&g.ID, &g.Name, &g.Desc); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[dao.group|ListGroup] d.db.Query err: %v", err))
			return
		}
		res = append(res, g)
	}
	return
}

func intConv(i int64, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	return int(i), nil
}
