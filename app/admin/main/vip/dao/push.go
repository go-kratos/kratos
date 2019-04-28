package dao

import (
	"context"
	xsql "database/sql"
	"fmt"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

// GetPushData get push data by id
func (d *Dao) GetPushData(c context.Context, id int64) (r *model.VipPushData, err error) {
	row := d.db.QueryRow(c, _selPushDataByIDSQL, id)
	r = new(model.VipPushData)
	if err = row.Scan(&r.ID, &r.DisableType, &r.GroupName, &r.Title, &r.Content, &r.PushTotalCount, &r.PushedCount, &r.ProgressStatus, &r.Status, &r.Platform, &r.LinkType, &r.LinkURL, &r.ErrorCode, &r.ExpiredDayStart, &r.ExpiredDayEnd, &r.EffectStartDate, &r.EffectEndDate, &r.PushStartTime, &r.PushEndTime, &r.Operator); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("scan_error")
	}
	return
}

// AddPushData add push data
func (d *Dao) AddPushData(c context.Context, r *model.VipPushData) (id int64, err error) {
	var result xsql.Result
	if result, err = d.db.Exec(c, _addPushDataSQL, r.GroupName, r.Title, r.Content, r.PushTotalCount, r.PushedCount, r.ProgressStatus, r.Status, r.Platform, r.LinkType, r.LinkURL, r.ExpiredDayStart, r.ExpiredDayEnd, r.EffectStartDate, r.EffectEndDate, r.PushStartTime, r.PushEndTime, r.Operator); err != nil {
		err = errors.WithStack(err)
		return
	}

	if id, err = result.LastInsertId(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpdatePushData update push data
func (d *Dao) UpdatePushData(c context.Context, r *model.VipPushData) (eff int64, err error) {
	var result xsql.Result
	if result, err = d.db.Exec(c, _updatePushDataSQL, r.GroupName, r.Title, r.Content, r.PushTotalCount, r.ProgressStatus, r.Platform, r.LinkType, r.LinkURL, r.ExpiredDayStart, r.ExpiredDayEnd, r.EffectStartDate, r.EffectEndDate, r.PushStartTime, r.PushEndTime, r.Operator, r.ID); err != nil {
		err = errors.WithStack(err)
		return
	}

	if eff, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// PushDataCount sel push data count
func (d *Dao) PushDataCount(c context.Context, arg *model.ArgPushData) (count int64, err error) {
	row := d.db.QueryRow(c, _selPushDataCountSQL+d.convertPushDataSQL(arg))

	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// PushDatas sel push datas
func (d *Dao) PushDatas(c context.Context, arg *model.ArgPushData) (res []*model.VipPushData, err error) {
	var rows *sql.Rows
	sql := _selPushDataSQL + d.convertPushDataSQL(arg)

	if arg.PN == 0 {
		arg.PN = 1
	}

	if arg.PS == 0 || arg.PS > 100 {
		arg.PS = _defps
	}

	sql += fmt.Sprintf(" ORDER BY id DESC LIMIT %v,%v", (arg.PN-1)*arg.PS, arg.PS)

	if rows, err = d.db.Query(c, sql); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := new(model.VipPushData)

		if err = rows.Scan(&r.ID, &r.DisableType, &r.GroupName, &r.Title, &r.Content, &r.PushTotalCount, &r.PushedCount, &r.ProgressStatus, &r.Status, &r.Platform, &r.LinkType, &r.LinkURL, &r.ErrorCode, &r.ExpiredDayStart, &r.ExpiredDayEnd, &r.EffectStartDate, &r.EffectEndDate, &r.PushStartTime, &r.PushEndTime, &r.Operator); err != nil {
			err = errors.WithStack(err)
		}
		res = append(res, r)
	}

	err = rows.Err()
	return
}

func (d *Dao) convertPushDataSQL(arg *model.ArgPushData) string {
	sql := " "
	if arg.ProgressStatus != 0 {
		sql += fmt.Sprintf(" AND progress_status=%v", arg.ProgressStatus)
	}
	if arg.Status != 0 {
		sql += fmt.Sprintf(" AND status=%v", arg.Status)
	}
	return sql
}

// DelPushData .
func (d *Dao) DelPushData(c context.Context, id int64) (err error) {
	if _, err = d.db.Exec(c, _delPushDataSQL, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// DisablePushData .
func (d *Dao) DisablePushData(c context.Context, res *model.VipPushData) (err error) {
	if _, err = d.db.Exec(c, _disablePushDataSQL, res.ProgressStatus, res.PushTotalCount, res.EffectEndDate, res.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
