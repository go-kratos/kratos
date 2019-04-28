package archive

import (
	"context"

	"go-common/library/log"
)

const (
	_seNoteSQL = "SELECT note FROM archive_video_audit WHERE vid=?;"
)

//VideoAuditNote get note by vid
func (d *Dao) VideoAuditNote(c context.Context, vid int64) (note string, err error) {
	if err = d.db.QueryRow(c, _seNoteSQL, vid).Scan(&note); err != nil {
		log.Error("VideoAuditNote db.row.scan error(%v) vid(%d)", err, vid)
	}
	return
}
