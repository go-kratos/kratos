package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/service/main/push/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_batch = 1000

	_reportSQL        = `SELECT id,app_id,platform_id,mid,buvid,device_token,build,time_zone,notify_switch,device_brand,device_model,os_version,extra,dtime FROM push_reports WHERE token_hash=?`
	_reportByIDSQL    = `SELECT id,app_id,platform_id,mid,buvid,device_token,build,time_zone,notify_switch,device_brand,device_model,os_version,extra FROM push_reports WHERE id=?`
	_reportsSQL       = `SELECT id,app_id,platform_id,mid,buvid,device_token,build,time_zone,notify_switch,device_brand,device_model,os_version,extra FROM push_reports WHERE token_hash IN(%s) AND dtime=0`
	_lastReportIDSQL  = `SELECT id FROM push_reports ORDER BY id DESC LIMIT 1`
	_addReportSQL     = `INSERT INTO push_reports (app_id,platform_id,mid,buvid,device_token,token_hash,build,time_zone,notify_switch,device_brand,device_model,os_version,extra) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_updateReportSQL  = `UPDATE push_reports SET app_id=?,platform_id=?,mid=?,buvid=?,build=?,time_zone=?,notify_switch=?,device_brand=?,device_model=?,os_version=?,extra=?,mtime=?,dtime=0 WHERE id=?`
	_delReportSQL     = `UPDATE push_reports SET dtime=? WHERE token_hash=?`
	_reportsByMidSQL  = `SELECT id,app_id,platform_id,mid,buvid,device_token,build,time_zone,notify_switch,device_brand,device_model,os_version,extra FROM push_reports WHERE mid=? AND dtime=0`
	_reportsByMidsSQL = `SELECT id,app_id,platform_id,mid,buvid,device_token,build,time_zone,notify_switch,device_brand,device_model,os_version,extra FROM push_reports WHERE mid IN (%s) AND dtime=0`
	_reportsByIDSQL   = `SELECT id,app_id,platform_id,mid,buvid,device_token,build,time_zone,notify_switch,device_brand,device_model,os_version,extra FROM push_reports WHERE id IN (%s) AND dtime=0`
)

// Report gets report by device_token.
func (d *Dao) Report(c context.Context, dt string) (r *model.Report, err error) {
	r = &model.Report{}
	th := model.HashToken(dt)
	if err = d.reportStmt.QueryRow(c, th).Scan(&r.ID, &r.APPID, &r.PlatformID, &r.Mid, &r.Buvid,
		&r.DeviceToken, &r.Build, &r.TimeZone, &r.NotifySwitch, &r.DeviceBrand, &r.DeviceModel, &r.OSVersion, &r.Extra, &r.Dtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		log.Error("d.reportStmt.QueryRow(%s) error(%v)", dt, err)
		PromError("mysql:获取上报")
	}
	return
}

// ReportByID gets report by id.
func (d *Dao) ReportByID(c context.Context, id int64) (r *model.Report, err error) {
	r = &model.Report{}
	if err = d.reportByIDStmt.QueryRow(c, id).Scan(&r.ID, &r.APPID, &r.PlatformID, &r.Mid, &r.Buvid,
		&r.DeviceToken, &r.Build, &r.TimeZone, &r.NotifySwitch, &r.DeviceBrand, &r.DeviceModel, &r.OSVersion, &r.Extra); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		log.Error("d.reportByIDStmt.QueryRow(%d) error(%v)", id, err)
		PromError("mysql:获取上报")
	}
	return
}

// Reports gets reports by device_token.
func (d *Dao) Reports(c context.Context, dts []string) (rs []*model.Report, err error) {
	var ths []int64
	for _, dt := range dts {
		ths = append(ths, model.HashToken(dt))
	}
	s := fmt.Sprintf(_reportsSQL, xstr.JoinInts(ths))
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, s); err != nil {
		log.Error("d.reports Query() error(%v)", err)
		PromError("mysql:通过tokens获取上报Query")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Report{}
		if err = rows.Scan(&r.ID, &r.APPID, &r.PlatformID, &r.Mid, &r.Buvid, &r.DeviceToken,
			&r.Build, &r.TimeZone, &r.NotifySwitch, &r.DeviceBrand, &r.DeviceModel, &r.OSVersion, &r.Extra); err != nil {
			log.Error("d.Reports Scan() error(%v)", err)
			PromError("mysql:通过tokens获取上报Scan")
			return
		}
		rs = append(rs, r)
	}
	return
}

// AddReport adds report.
func (d *Dao) AddReport(ctx context.Context, r *model.Report) (id int64, err error) {
	th := model.HashToken(r.DeviceToken)
	res, err := d.addReportStmt.Exec(ctx, r.APPID, r.PlatformID, r.Mid, r.Buvid, r.DeviceToken, th, r.Build, r.TimeZone, r.NotifySwitch, r.DeviceBrand, r.DeviceModel, r.OSVersion, r.Extra)
	if err != nil {
		log.Error("d.AddReport(%+v) error(%v)", r, err)
		PromError("mysql:AddReport")
		return
	}
	PromInfo("mysql:AddReport")
	return res.LastInsertId()
}

// UpdateReport update report.
func (d *Dao) UpdateReport(ctx context.Context, r *model.Report) (err error) {
	if _, err = d.updateReportStmt.Exec(ctx, r.APPID, r.PlatformID, r.Mid, r.Buvid, r.Build, r.TimeZone, r.NotifySwitch, r.DeviceBrand, r.DeviceModel, r.OSVersion, r.Extra, time.Now(), r.ID); err != nil {
		log.Error("UpdateReport(%+v) error(%v)", r, err)
		PromError("mysql:UpdateReport")
		return
	}
	PromInfo("mysql:UpdateReport")
	return
}

// DelReport delete report.
func (d *Dao) DelReport(c context.Context, dt string) (rows int64, err error) {
	th := model.HashToken(dt)
	now := time.Now().Unix()
	var res sql.Result
	if res, err = d.delReportStmt.Exec(c, now, th); err != nil {
		log.Error("d.delReportStmt.Exec(%s,%d) error(%v)", dt, th, err)
		PromError("mysql:删除上报")
		return
	}
	return res.RowsAffected()
}

// ReportsByMid gets reports by mid.
func (d *Dao) ReportsByMid(c context.Context, mid int64) (res []*model.Report, err error) {
	var rows *xsql.Rows
	if rows, err = d.reportsByMidStmt.Query(c, mid); err != nil {
		log.Error("d.reportsByMid Query() error(%v)", err)
		PromError("mysql:通过mid获取上报Query")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Report{}
		if err = rows.Scan(&r.ID, &r.APPID, &r.PlatformID, &r.Mid, &r.Buvid, &r.DeviceToken,
			&r.Build, &r.TimeZone, &r.NotifySwitch, &r.DeviceBrand, &r.DeviceModel, &r.OSVersion, &r.Extra); err != nil {
			log.Error("d.ReportsByMid Scan() error(%v)", err)
			PromError("mysql:通过mid获取上报Scan")
			return
		}
		res = append(res, r)
	}
	PromInfo("mysql:通过mid获取上报")
	return
}

// ReportsByMids gets reports by mids.
func (d *Dao) ReportsByMids(c context.Context, mids []int64) (reports map[int64][]*model.Report, err error) {
	reports = make(map[int64][]*model.Report)
	for {
		midsLen := len(mids)
		if midsLen == 0 {
			return
		}
		var part []int64
		if midsLen >= _batch {
			part = mids[:_batch]
		} else {
			part = mids[:]
		}
		mids = mids[len(part):]
		s := fmt.Sprintf(_reportsByMidsSQL, xstr.JoinInts(part))
		var rows *xsql.Rows
		if rows, err = d.db.Query(c, s); err != nil {
			log.Error("d.reportsByMids Query() error(%v)", err)
			PromError("mysql:通过mids获取上报Query")
			return
		}
		defer rows.Close()
		for rows.Next() {
			r := &model.Report{}
			if err = rows.Scan(&r.ID, &r.APPID, &r.PlatformID, &r.Mid, &r.Buvid, &r.DeviceToken,
				&r.Build, &r.TimeZone, &r.NotifySwitch, &r.DeviceBrand, &r.DeviceModel, &r.OSVersion, &r.Extra); err != nil {
				log.Error("d.ReportsByMids Scan() error(%v)", err)
				PromError("mysql:通过mids获取上报Scan")
				return
			}
			reports[r.Mid] = append(reports[r.Mid], r)
		}
		PromInfo("mysql:通过mids获取上报")
	}
}

// ReportsByID gets reports by mids.
func (d *Dao) ReportsByID(c context.Context, ids []int64) (reports []*model.Report, err error) {
	s := fmt.Sprintf(_reportsByIDSQL, xstr.JoinInts(ids))
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, s); err != nil {
		log.Error("d.reportsByID Query() error(%v)", err)
		PromError("mysql:通过id获取上报Query")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Report{}
		if err = rows.Scan(&r.ID, &r.APPID, &r.PlatformID, &r.Mid, &r.Buvid, &r.DeviceToken,
			&r.Build, &r.TimeZone, &r.NotifySwitch, &r.DeviceBrand, &r.DeviceModel, &r.OSVersion, &r.Extra); err != nil {
			log.Error("d.ReportsByID Scan() error(%v)", err)
			PromError("mysql:通过id获取上报Scan")
			return
		}
		reports = append(reports, r)
	}
	PromInfo("mysql:通过id获取上报")
	return
}
