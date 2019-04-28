package archive

import (
	"context"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
	"time"
)

const (
	_statsPointSQL = "SELECT id,type,content,ctime,mtime FROM archive_report_sum WHERE mtime>=? AND mtime<? AND type=?"
)

// StatsPoints get archive_report_sum data by type and time
func (d *Dao) StatsPoints(c context.Context, stime, etime time.Time, typeInt int8) (points []*archive.StatsPoint, err error) {
	rows, err := d.rddb.Query(c, _statsPointSQL, stime, etime, typeInt)
	if err != nil {
		log.Error("d.StatsPoints.Query(%v,%v,%v) error(%v)", stime, etime, typeInt, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		point := &archive.StatsPoint{}
		if err = rows.Scan(&point.ID, &point.Type, &point.Content, &point.Ctime, &point.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		points = append(points, point)
	}
	return
}
