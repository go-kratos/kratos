package archive

import (
	"context"
	"database/sql"
	"time"

	"go-common/library/log"
)

const (
	_hisCntSQL          = "SELECT COUNT(*) FROM archive_edit_history WHERE aid=?"
	_delEditHisSQL      = "DELETE FROM archive_edit_history WHERE mtime < ? LIMIT ?"
	_delVideoEditHisSQL = "DELETE FROM archive_video_edit_history WHERE mtime < ? LIMIT ?"
)

// HistoryCount get a archive history count.
func (d *Dao) HistoryCount(c context.Context, aid int64) (count int, err error) {
	row := d.db.QueryRow(c, _hisCntSQL, aid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// DelArcEditHistoryBefore  delete archive_edit_history before t.
func (d *Dao) DelArcEditHistoryBefore(c context.Context, t time.Time, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delEditHisSQL, t, limit)
	if err != nil {
		log.Error("db.Exec(%s, %s) error(%v)", _delEditHisSQL, t, err)
		return
	}
	return res.RowsAffected()
}

// DelArcVideoEditHistoryBefore  delete archive_video_edit_history before t.
func (d *Dao) DelArcVideoEditHistoryBefore(c context.Context, t time.Time, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delVideoEditHisSQL, t, limit)
	if err != nil {
		log.Error("db.Exec(%s, %s) error(%v)", _delVideoEditHisSQL, t, err)
		return
	}
	return res.RowsAffected()
}
