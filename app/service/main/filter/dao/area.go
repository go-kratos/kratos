package dao

import (
	"context"

	"go-common/app/service/main/filter/model"
	xsql "go-common/library/database/sql"
	"go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_areaList      = `SELECT id,name,showname,groupid,common_flag,ctime,mtime FROM filter_area_type WHERE is_delete=0`
	_lastTimestamp = `SELECT mtime FROM filter_area ORDER BY mtime desc LIMIT 1`
)

// AreaList .
func (d *Dao) AreaList(ctx context.Context) (list []*model.Area, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.mysql.Query(ctx, _areaList); err != nil {
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

// AreaLastTime .
func (d *Dao) AreaLastTime(ctx context.Context) (res int64, err error) {
	var mtime time.Time
	if err = d.mysql.QueryRow(ctx, _lastTimestamp).Scan(&mtime); err != nil {
		err = errors.WithStack(err)
		return
	}
	res = mtime.Time().Unix()
	return
}
