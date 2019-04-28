package archive

import (
	"context"
	"time"

	"go-common/library/log"
)

const (
	_inVideoShotSQL = "INSERT INTO archive_video_shot (id,count,ctime,mtime) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE count=?,mtime=? "
)

// AddVideoShot add a videoshot into mysql.
func (d *Dao) AddVideoShot(c context.Context, cid int64, count int) (rows int64, err error) {
	var now = time.Now()
	res, err := d.db.Exec(c, _inVideoShotSQL, cid, count, now, now, count, now)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
