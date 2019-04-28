package archive

import (
	"context"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inArcReportSQL = "INSERT IGNORE INTO archive_report (mid,aid,type,reason,pics,ctime,mtime) VALUES(?,?,?,?,?,?,?)"
	_arcReportSQL   = "SELECT aid,mid,type,reason,pics,ctime,mtime FROM archive_report WHERE aid=? AND mid=? LIMIT 1"
)

// AddArcReport insert archive_report.
func (d *Dao) AddArcReport(c context.Context, aa *archive.ArcReport) (id int64, err error) {
	res, err := d.db.Exec(c, _inArcReportSQL, aa.Mid, aa.Aid, aa.Type, aa.Reason, aa.Pics, aa.CTime, aa.MTime)
	if err != nil {
		log.Error("_inArcReport.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// ArcReport get archive_report by aid and mid.
func (d *Dao) ArcReport(c context.Context, aid, mid int64) (aa *archive.ArcReport, err error) {
	row := d.rddb.QueryRow(c, _arcReportSQL, aid, mid)
	aa = &archive.ArcReport{}
	if err = row.Scan(&aa.Aid, &aa.Mid, &aa.Type, &aa.Reason, &aa.Pics, &aa.CTime, &aa.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			aa = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
