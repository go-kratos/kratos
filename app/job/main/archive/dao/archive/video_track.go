package archive

import (
	"context"

	"go-common/library/log"
)

const (
	_inVideoHis = "INSERT INTO archive_video_track(aid,filename,status,xcode_state,remark,ctime,mtime) VALUES(?,?,?,?,?,?,?)"
)

// InVideoHis insert video track history
func (d *Dao) InVideoHis(c context.Context, aid int64, filename string, status int16, xcodeState int8, remark string, ctime, mtime string) (rows int64, err error) {
	rs, err := d.db.Exec(c, _inVideoHis, aid, filename, status, xcodeState, remark, ctime, mtime)
	if err != nil {
		log.Error("d.inVideoHisStmt.Exec(%d, %s, %d, %d, %s, %s, %s) error(%v)", aid, filename, status, xcodeState, remark, ctime, mtime, err)
		return
	}
	rows, err = rs.RowsAffected()
	return
}
