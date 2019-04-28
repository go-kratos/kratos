package archive

import (
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inAuditSQL = "INSERT INTO archive_video_audit (vid,aid,tid,oname,note,reason,passtime) VALUES (?,?,?,?,?,?,now()) ON DUPLICATE KEY UPDATE tid=?,oname=?,note=?,reason=?,passtime=now()"
)

// TxAddAudit insert video audit
func (d *Dao) TxAddAudit(tx *sql.Tx, aid, vid, tagID int64, oname, note, reason string) (rows int64, err error) {
	res, err := tx.Exec(_inAuditSQL, vid, aid, tagID, oname, note, reason, tagID, oname, note, reason)
	if err != nil {
		log.Error("d.TxAddAudit.tx.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
