package archive

import (
	"context"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
)

const (
	_taskTooksByHalfHourSQL = "SELECT id,m50,m60,m80,m90,type,ctime,mtime FROM task_dispatch_took WHERE type=2 AND ctime>=? AND ctime<=? ORDER BY ctime ASC"
)

// TaskTooksByHalfHour get TaskTooks by half hour
func (d *Dao) TaskTooksByHalfHour(c context.Context, stime time.Time, etime time.Time) (tooks []*archive.TaskTook, err error) {
	rows, err := d.rddb.Query(c, _taskTooksByHalfHourSQL, stime, etime)
	if err != nil {
		log.Error("d.TaskTooksByHalfHour.Query(%v,%v) error(%v)", stime, etime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		took := &archive.TaskTook{}
		if err = rows.Scan(&took.ID, &took.M50, &took.M60, &took.M80, &took.M90, &took.TypeID, &took.Ctime, &took.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		tooks = append(tooks, took)
	}
	return
}
