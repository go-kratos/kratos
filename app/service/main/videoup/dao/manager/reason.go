package manager

import (
	"context"
	xsql "database/sql"

	"go-common/library/log"
)

const (
	_reasonSQL = "SELECT reason.tag_id as tag_id from reason_log left join reason on reason_log.reason_id=reason.id where reason_log.type=1 AND reason_log.oid=? order by reason_log.id desc limit 1;"
)

// ArcReason get a archive reason tag
func (d *Dao) ArcReason(c context.Context, aid int64) (tagID int64, err error) {
	var (
		row    = d.managerDB.QueryRow(c, _reasonSQL, aid)
		tagIDI xsql.NullInt64
	)
	if err := row.Scan(&tagIDI); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("ArcReason row.Scan error(%v)", err)
		}
	}
	log.Info("ArcReason retrun(%v)", tagIDI)
	tagID = tagIDI.Int64
	return
}
