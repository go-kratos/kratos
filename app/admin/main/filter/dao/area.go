package dao

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/filter/model"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_areaGroupTotal  = `SELECT count(*) FROM filter_area_group WHERE is_delete=0`
	_areaGroupList   = `SELECT id,name,ctime,mtime FROM filter_area_group WHERE is_delete=0 LIMIT ?,?`
	_areaGroup       = `SELECT id,name,ctime,mtime FROM filter_area_group WHERE id=? AND is_delete=0 LIMIT 1`
	_areaGroupByName = `SELECT id,name,ctime,mtime FROM filter_area_group WHERE name=? AND is_delete=0 LIMIT 1`
	_insertAreaGroup = `INSERT INTO filter_area_group (name) VALUES (?)`

	_insertAreaGroupLog = `INSERT INTO filter_area_group_log (groupid,adid,ad_name,comment,state) VALUES (?,?,?,?,?)`

	_areaTotal  = `SELECT count(*) FROM filter_area_type WHERE is_delete=0 AND groupid=?`
	_areaList   = `SELECT id,name,showname,groupid,common_flag,ctime,mtime FROM filter_area_type WHERE is_delete=0 AND groupid=? LIMIT ?,?`
	_area       = `SELECT id,name,showname,groupid,common_flag,ctime,mtime FROM filter_area_type WHERE id=? AND is_delete=0`
	_areaByName = `SELECT id,name,showname,groupid,common_flag,ctime,mtime FROM filter_area_type WHERE name=? AND is_delete=0`
	_insertArea = `INSERT INTO filter_area_type (name,showname,common_flag,groupid) VALUES (?,?,?,?)`
	_updateArea = `UPDATE filter_area_type SET name=?,showname=?,common_flag=?,groupid=? WHERE id=?`

	_areaLog       = `SELECT id,adid,ad_name,comment,state,ctime FROM filter_area_type_log WHERE areaid=?`
	_insertAreaLog = `INSERT INTO filter_area_type_log (areaid,adid,ad_name,comment,state) VALUES (?,?,?,?,?)`
)

// AreaGroupTotal .
func (d *Dao) AreaGroupTotal(ctx context.Context) (total int, err error) {
	row := d.mysql.QueryRow(ctx, _areaGroupTotal)
	if err = row.Scan(&total); err != nil {
		return
	}
	return
}

// AreaGroupList .
func (d *Dao) AreaGroupList(ctx context.Context, pn, ps int) (list []*model.AreaGroup, err error) {
	list = make([]*model.AreaGroup, 0)
	var (
		rows   *xsql.Rows
		start  = (pn - 1) * ps
		offset = ps
	)
	if rows, err = d.mysql.Query(ctx, _areaGroupList, start, offset); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ag := &model.AreaGroup{}
		if err = rows.Scan(&ag.ID, &ag.Name, &ag.Ctime, &ag.Mtime); err != nil {
			err = errors.WithStack(err)
			return
		}
		list = append(list, ag)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AreaGroup .
func (d *Dao) AreaGroup(ctx context.Context, groupID int) (ag *model.AreaGroup, err error) {
	var (
		row *xsql.Row
	)
	ag = &model.AreaGroup{}
	row = d.mysql.QueryRow(ctx, _areaGroup, groupID)
	if err = row.Scan(&ag.ID, &ag.Name, &ag.Ctime, &ag.Mtime); err != nil {
		if err == xsql.ErrNoRows {
			ag = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// AreaGroupByName .
func (d *Dao) AreaGroupByName(ctx context.Context, name string) (ag *model.AreaGroup, err error) {
	var (
		row *xsql.Row
	)
	ag = &model.AreaGroup{}
	row = d.mysql.QueryRow(ctx, _areaGroupByName, name)
	if err = row.Scan(&ag.ID, &ag.Name, &ag.Ctime, &ag.Mtime); err != nil {
		if err == xsql.ErrNoRows {
			ag = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// TxInsertAreaGroupLog .
func (d *Dao) TxInsertAreaGroupLog(ctx context.Context, tx *xsql.Tx, groupID int64, log *model.AreaGroupLog) (err error) {
	if _, err = tx.Exec(_insertAreaGroupLog, groupID, log.AdID, log.AdName, log.Comment, log.State); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// TxInsertAreaGroup .
func (d *Dao) TxInsertAreaGroup(ctx context.Context, tx *xsql.Tx, areaGroup *model.AreaGroup) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = tx.Exec(_insertAreaGroup, areaGroup.Name); err != nil {
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// AreaLog .
func (d *Dao) AreaLog(ctx context.Context, id int) (list []*model.AreaLog, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.mysql.Query(ctx, _areaLog, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		log := &model.AreaLog{}
		if err = rows.Scan(&log.ID, &log.AdID, &log.AdName, &log.Comment, &log.State, &log.Ctime); err != nil {
			err = errors.WithStack(err)
			return
		}
		list = append(list, log)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// TxInsertAreaLog .
func (d *Dao) TxInsertAreaLog(ctx context.Context, tx *xsql.Tx, areaid int64, log *model.AreaLog) (err error) {
	if _, err = tx.Exec(_insertAreaLog, areaid, log.AdID, log.AdName, log.Comment, log.State); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// AreaTotal .
func (d *Dao) AreaTotal(ctx context.Context, groupID int) (total int, err error) {
	row := d.mysql.QueryRow(ctx, _areaTotal, groupID)
	if err = row.Scan(&total); err != nil {
		return
	}
	return
}

// AreaList .
func (d *Dao) AreaList(ctx context.Context, groupID int, pn, ps int) (list []*model.Area, err error) {
	list = make([]*model.Area, 0)
	var (
		rows   *xsql.Rows
		start  = (pn - 1) * ps
		offset = ps
	)
	if rows, err = d.mysql.Query(ctx, _areaList, groupID, start, offset); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		area := &model.Area{}
		if err = rows.Scan(&area.ID, &area.Name, &area.ShowName, &area.GroupID, &area.CommonFlag, &area.Ctime, &area.Mtime); err != nil {
			err = errors.WithStack(err)
			return
		}
		list = append(list, area)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// Area .
func (d *Dao) Area(ctx context.Context, id int) (area *model.Area, err error) {
	var (
		row *xsql.Row
	)
	area = &model.Area{}
	row = d.mysql.QueryRow(ctx, _area, id)
	if err = row.Scan(&area.ID, &area.Name, &area.ShowName, &area.GroupID, &area.CommonFlag, &area.Ctime, &area.Mtime); err != nil {
		if err == xsql.ErrNoRows {
			area = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// AreaByName .
func (d *Dao) AreaByName(ctx context.Context, name string) (area *model.Area, err error) {
	var (
		row *xsql.Row
	)
	area = &model.Area{}
	row = d.mysql.QueryRow(ctx, _areaByName, name)
	if err = row.Scan(&area.ID, &area.Name, &area.ShowName, &area.GroupID, &area.CommonFlag, &area.Ctime, &area.Mtime); err != nil {
		if err == xsql.ErrNoRows {
			area = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// TxInsertArea .
func (d *Dao) TxInsertArea(ctx context.Context, tx *xsql.Tx, area *model.Area) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = tx.Exec(_insertArea, area.Name, area.ShowName, area.CommonFlag, area.GroupID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// TxUpdateArea .
func (d *Dao) TxUpdateArea(ctx context.Context, tx *xsql.Tx, area *model.Area) (err error) {
	if _, err = tx.Exec(_updateArea, area.Name, area.ShowName, area.CommonFlag, area.GroupID, area.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
