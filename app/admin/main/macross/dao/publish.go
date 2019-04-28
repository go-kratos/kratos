package dao

import (
	"bytes"
	"context"
	"fmt"

	"go-common/app/admin/main/macross/model/publish"
	"go-common/library/log"
)

const (
	_logSharding = 10
	// dashborad
	_inDashboradSQL     = `INSERT INTO dashboard (name,label,commit_info,out_url,coverage_url,text_size_arm64,res_size,extra) VALUES(?,?,?,?,?,?,?,?)`
	_inDashboradLogsSQL = `INSERT INTO dashboard_log_%02d (dashboard_id,level,msg) VALUES`
)

func (d *Dao) hitLogs(id int64) int64 {
	return id % _logSharding
}

// Dashborad insert dashboard.
func (d *Dao) Dashborad(c context.Context, dashboard *publish.Dashboard) (rows int64, err error) {
	res, err := d.db.Exec(c, _inDashboradSQL, dashboard.Name, dashboard.Label, dashboard.Commit, dashboard.OutURL, dashboard.CoverageURL, dashboard.TextSizeArm64, dashboard.ResSize, dashboard.Extra)
	if err != nil {
		log.Error("Dashborad() d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.LastInsertId()
	return
}

// DashboradLogs insert dashboard log.
func (d *Dao) DashboradLogs(c context.Context, id int64, logs []*publish.Log) (rows int64, err error) {
	var (
		buffer   bytes.Buffer
		insertTp string
	)
	insertTp = "(%d,'%s','%s'),"
	buffer.WriteString(fmt.Sprintf(_inDashboradLogsSQL, d.hitLogs(id)))
	for _, v := range logs {
		buffer.WriteString(fmt.Sprintf(insertTp, id, v.Level, v.Msg))
	}
	buffer.Truncate(buffer.Len() - 1)
	res, err := d.db.Exec(c, buffer.String())
	if err != nil {
		log.Error("DashboradLogs d.db.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
