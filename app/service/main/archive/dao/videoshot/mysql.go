package videoshot

import (
	"context"
	"database/sql"

	"go-common/app/service/main/archive/model/videoshot"
	"go-common/library/log"
)

const (
	_inSQL  = "INSERT INTO archive_video_shot (id,count,ctime,mtime) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE count=?,mtime=? "
	_getSQL = "SELECT id,count,ctime,mtime FROM archive_video_shot WHERE id=?"
)

// videoshot get a videoshot by id.
func (d *Dao) videoshot(c context.Context, cid int64) (shot *videoshot.Videoshot, err error) {
	d.infoProm.Incr("videoshot")
	row := d.getStmt.QueryRow(c, cid)
	shot = &videoshot.Videoshot{}
	if err = row.Scan(&shot.Cid, &shot.Count, &shot.CTime, &shot.MTime); err != nil {
		if err == sql.ErrNoRows {
			shot = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// addVideoshot add a videoshot into mysql.
func (d *Dao) addVideoshot(c context.Context, shot *videoshot.Videoshot) (cid int64, err error) {
	res, err := d.inStmt.Exec(c, shot.Cid, shot.Count, shot.CTime, shot.MTime, shot.Count, shot.MTime)
	if err != nil {
		log.Error("inStmt.Exec error(%v)", err)
		return
	}
	cid, err = res.LastInsertId()
	return
}
