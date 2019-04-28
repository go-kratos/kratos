package dao

import (
	"context"
	"database/sql"
	"time"

	pushmdl "go-common/app/service/main/push/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_delCallbacksSQL   = `DELETE FROM push_callbacks where ctime <= ? limit ?`
	_reportLastIDSQL   = `SELECT MAX(id) from push_reports`
	_reportsByRangeSQL = `SELECT id,app_id,platform_id,mid,buvid,device_token,build,time_zone,notify_switch,device_brand,device_model,os_version,extra FROM push_reports WHERE id>? and id<? and dtime=0`
	// for 全量推送
	_reportsTaskAllByRangeSQL = `SELECT platform_id,device_token,build FROM push_reports WHERE id>? and id<=? and app_id=? and dtime=0 and notify_switch=1`
)

// BeginTx begin transaction.
func (d *Dao) BeginTx(c context.Context) (*xsql.Tx, error) {
	return d.db.Begin(c)
}

// DelCallbacks deletes callbacks.
func (d *Dao) DelCallbacks(c context.Context, t time.Time, limit int) (rows int64, err error) {
	res, err := d.delCallbacksStmt.Exec(c, t, limit)
	if err != nil {
		log.Error("d.DelCallbacks(%v) error(%v)", t, err)
		PromError("mysql:DelCallbacks")
		return
	}
	rows, err = res.RowsAffected()
	return
}

// ReportLastID gets the latest ID of report database record.
func (d *Dao) ReportLastID(c context.Context) (id int64, err error) {
	if err = d.reportLastIDStmt.QueryRow(c).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return
		}
		log.Error("d.ReportLastID() error(%v)", err)
		PromError("mysql:ReportLastID")
	}
	return
}

// ReportsByRange gets reports by id range.
func (d *Dao) ReportsByRange(c context.Context, min, max int64) (rs []*pushmdl.Report, err error) {
	rows, err := d.reportsByRangeStmt.Query(c, min, max)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &pushmdl.Report{}
		if err = rows.Scan(&r.ID, &r.APPID, &r.PlatformID, &r.Mid, &r.Buvid, &r.DeviceToken,
			&r.Build, &r.TimeZone, &r.NotifySwitch, &r.DeviceBrand, &r.DeviceModel, &r.OSVersion, &r.Extra); err != nil {
			log.Error("d.ReportsByRange Scan() error(%v)", err)
			PromError("mysql:ReportsByRange")
			return
		}
		rs = append(rs, r)
	}
	return
}

// ReportsTaskAll gets reports by range
func (d *Dao) ReportsTaskAll(c context.Context, min, max, app int64) (rows *xsql.Rows, err error) {
	if rows, err = d.db.Query(c, _reportsTaskAllByRangeSQL, min, max, app); err != nil {
		log.Error("ReportsTaskAll load reports start(%d) end(%d) error(%v)", min, max, err)
	}
	return
}
