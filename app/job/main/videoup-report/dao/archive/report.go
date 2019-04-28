package archive

import (
	"context"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_reportAddByTypeIDSQL  = "INSERT INTO archive_report_sum(content,ctime,mtime,type) VALUE(?,?,?,?)"
	_reportGetByTypeIDSQL  = "SELECT id,content,ctime,mtime,type FROM archive_report_sum WHERE type=? AND mtime>=? AND mtime<=?"
	_reportLastByTypeIDSQL = "SELECT id,content,ctime,mtime,type FROM archive_report_sum WHERE type=? ORDER BY id DESC LIMIT 1"
)

// ReportLast get last inserted report
func (d *Dao) ReportLast(c context.Context, typeid int8) (report *archive.Report, err error) {
	row := d.db.QueryRow(c, _reportLastByTypeIDSQL, typeid)
	report = &archive.Report{}
	if err = row.Scan(&report.ID, &report.Content, &report.CTime, &report.MTime, &report.TypeID); err != nil {
		if err == sql.ErrNoRows {
			report = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}

// ReportAdd report add of typeid
func (d *Dao) ReportAdd(c context.Context, typeid int8, content string, ctime, mtime time.Time) (lastID int64, err error) {
	res, err := d.db.Exec(c, _reportAddByTypeIDSQL, content, ctime, mtime, typeid)
	if err != nil {
		log.Error("d.TaskTookAddStmt.Exec error(%v)", err)
		return
	}
	lastID, err = res.LastInsertId()
	return
}

// Reports report get of typeid
func (d *Dao) Reports(c context.Context, typeid int8, stime, etime time.Time) (reports []*archive.Report, err error) {
	rows, err := d.db.Query(c, _reportGetByTypeIDSQL, typeid, stime, etime)
	if err != nil {
		log.Error("d.Reports.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		report := &archive.Report{}
		if err = rows.Scan(&report.ID, &report.Content, &report.CTime, &report.MTime, &report.TypeID); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		reports = append(reports, report)
	}
	return
}
