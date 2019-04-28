package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_jointlysWillEffectSQL = "SELECT id,title,content,operator,start_time,end_time,link,is_hot,ctime,mtime FROM vip_jointly WHERE start_time > ? AND deleted = 0 ORDER BY mtime DESC;"
	_jointlysEffectSQL     = "SELECT id,title,content,operator,start_time,end_time,link,is_hot,ctime,mtime FROM vip_jointly WHERE start_time < ? AND end_time > ? AND deleted = 0 ORDER BY mtime DESC;"
	_jointlysLoseEffectSQL = "SELECT id,title,content,operator,start_time,end_time,link,is_hot,ctime,mtime FROM vip_jointly WHERE end_time < ? AND deleted = 0 ORDER BY mtime DESC;"
	_addjointlySQL         = "INSERT INTO vip_jointly(title,content,operator,start_time,end_time,link,is_hot)VALUES(?,?,?,?,?,?,?);"
	_updateJointlySQL      = "UPDATE vip_jointly SET title = ?,content = ?,operator = ?,link = ?,is_hot = ?,start_time = ?,end_time = ? WHERE id = ?;"
	_deleteJointlySQL      = "UPDATE vip_jointly SET deleted = 1 WHERE id = ?;"
)

// AddJointly add jointly.
func (d *Dao) AddJointly(c context.Context, j *model.Jointly) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _addjointlySQL, j.Title, j.Content, j.Operator, j.StartTime, j.EndTime, j.Link, j.IsHot); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// JointlysByState jointly by state.
func (d *Dao) JointlysByState(c context.Context, state int8, now int64) (res []*model.Jointly, err error) {
	var rows *sql.Rows
	switch state {
	case model.WillEffect:
		rows, err = d.db.Query(c, _jointlysWillEffectSQL, now)
	case model.Effect:
		rows, err = d.db.Query(c, _jointlysEffectSQL, now, now)
	case model.LoseEffect:
		rows, err = d.db.Query(c, _jointlysLoseEffectSQL, now)
	default:
		return
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Jointly)
		if err = rows.Scan(&r.ID, &r.Title, &r.Content, &r.Operator, &r.StartTime, &r.EndTime, &r.Link, &r.IsHot, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// UpdateJointly update jointly.
func (d *Dao) UpdateJointly(c context.Context, j *model.Jointly) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateJointlySQL, j.Title, j.Content, j.Operator, j.Link, j.IsHot, j.StartTime, j.EndTime, j.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DeleteJointly delete jointly.
func (d *Dao) DeleteJointly(c context.Context, id int64) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _deleteJointlySQL, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}
