package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/admin/main/point/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_allPointConf    = "SELECT id,app_id,point,operator,change_type,ctime,mtime FROM point_conf"
	_getPointConf    = "SELECT id,app_id,point,operator,change_type,ctime,mtime FROM point_conf WHERE id=?"
	_addPointConf    = "INSERT INTO point_conf (app_id,point,operator,change_type) VALUES (?,?,?,?)"
	_updatePointConf = "UPDATE point_conf SET point=?,operator=?,change_type=? where ID=?"
	_allAppInfo      = "SELECT `id`,`name`,`app_key`,`purge_url` FROM `point_app_info`;"
)

// PointConfList .
func (d *Dao) PointConfList(c context.Context) (res []*model.PointConf, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allPointConf); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.PointConf{}
		if err = rows.Scan(&r.ID, &r.AppID, &r.Point, &r.Operator, &r.ChangeType, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// PointCoinInfo .
func (d *Dao) PointCoinInfo(c context.Context, id int64) (r *model.PointConf, err error) {
	row := d.db.QueryRow(c, _getPointConf, id)
	r = new(model.PointConf)
	if err = row.Scan(&r.ID, &r.AppID, &r.Point, &r.Operator, &r.ChangeType, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// PointCoinAdd .
func (d *Dao) PointCoinAdd(c context.Context, pc *model.PointConf) (id int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _addPointConf, pc.AppID, pc.Point, pc.Operator, pc.ChangeType); err != nil {
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// PointCoinEdit .
func (d *Dao) PointCoinEdit(c context.Context, mp *model.PointConf) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updatePointConf, mp.Point, mp.Operator, mp.ChangeType, mp.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AllAppInfo all appinfo.
func (d *Dao) AllAppInfo(c context.Context) (res []*model.AppInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allAppInfo); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.AppInfo{}
		if err = rows.Scan(&r.ID, &r.Name, &r.AppKey, &r.PurgeURL); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}
