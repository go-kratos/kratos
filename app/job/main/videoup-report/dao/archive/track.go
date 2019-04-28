package archive

import (
	"context"

	"go-common/library/log"
)

const (
	_inArcHis = "INSERT INTO archive_track(aid,state,round,attribute,remark,ctime,mtime) VALUES(?,?,?,?,?,?,?)"
)

// AddTrack insert archive track history
func (d *Dao) AddTrack(c context.Context, aid int64, state int, round int8, attr int32, remark string, ctime, mtime string) (rows int64, err error) {
	rs, err := d.db.Exec(c, _inArcHis, aid, state, round, attr, remark, ctime, mtime)
	if err != nil {
		log.Error("d.inArcHisStmt.Exec(%d, %d, %d, %s, %s, %s) error(%v)", aid, state, round, attr, remark, ctime, mtime, err)
		return
	}
	rows, err = rs.RowsAffected()
	return
}
