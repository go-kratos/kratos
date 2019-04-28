package archive

import (
	"context"
	"database/sql"

	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// insert
	_inAuditSQL = `INSERT IGNORE INTO archive_video_audit (vid,aid,tid,oname,reason) VALUES(?,?,0,'videoup-job',?) 
					ON DUPLICATE KEY UPDATE oname='videoup-job',reason=?`
	_seRsnSQL = "SELECT reason FROM archive_video_audit WHERE vid=?"
)

// TxAddAudit add video audit by vid.
func (d *Dao) TxAddAudit(tx *xsql.Tx, vid, aid int64, reason string) (rows int64, err error) {
	row, err := tx.Exec(_inAuditSQL, vid, aid, reason, reason)
	if err != nil {
		log.Error("tx.Exec(%d, %d, %s) error(%v)", vid, aid, reason, err)
		return
	}
	return row.RowsAffected()
}

// Reason get a archive video reject reason by vid.
func (d *Dao) Reason(c context.Context, vid int64) (reason string, err error) {
	row := d.db.QueryRow(c, _seRsnSQL, vid)
	if err = row.Scan(&reason); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
