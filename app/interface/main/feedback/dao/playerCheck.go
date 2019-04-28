package dao

import (
	"context"
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_inPlayCheckSQL = `INSERT INTO campus_network_check_record (platform,check_time,isp,region,school,mid,ip,ip_change_times,cdn,connect_speed,io_speed,aid,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
)

// InPlayCheck insert player check info into db
func (d *Dao) InPlayCheck(c context.Context, platform, isp, ipChangeTimes int, mid, checkTime, aid, connectSpeed, ioSpeed int64, region, school, ip, cdn string) (rows int64, err error) {
	var (
		now        = time.Now()
		checkTime2 = xtime.Time(checkTime).Time()
	)
	res, err := d.dbMs.Exec(c, _inPlayCheckSQL, platform, checkTime2, isp, region, school, mid, ip, ipChangeTimes, cdn, connectSpeed, ioSpeed, aid, now, now)
	if err != nil {
		log.Error("d.InPlayCheck.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
