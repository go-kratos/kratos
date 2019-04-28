package ugc

import (
	"context"
	"time"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_deletedUp   = "SELECT mid FROM ugc_uploader WHERE toinit = 2 AND deleted = 1 AND retry < unix_timestamp(now()) LIMIT 1"
	_finishDelUp = "UPDATE ugc_uploader SET toinit = 0 WHERE mid = ? AND deleted = 1"
	_ppDelUp     = "UPDATE ugc_uploader SET retry = ? WHERE mid = ? AND deleted = 1"
	_upArcs      = "SELECT aid FROM ugc_archive WHERE mid = ? AND deleted = 0 LIMIT 50"
	_upCountArc  = "SELECT count(1) FROM ugc_archive WHERE mid = ? AND deleted = 0"
)

// DeletedUp picks the deleted uppers, toinit = 2 and deleted = 1
func (d *Dao) DeletedUp(c context.Context) (mid int64, err error) {
	err = d.DB.QueryRow(c, _deletedUp).Scan(&mid)
	return
}

// FinishDelUp updates the submit toinit from 2 to 0
func (d *Dao) FinishDelUp(c context.Context, mid int64) (err error) {
	if _, err = d.DB.Exec(c, _finishDelUp, mid); err != nil {
		log.Error("FinishDelUp Error: %v", mid, err)
	}
	return
}

// PpDelUp postpones the upper's videos submit in 30 mins
func (d *Dao) PpDelUp(c context.Context, mid int64) (err error) {
	var delay = time.Now().Unix() + int64(d.conf.UgcSync.Frequency.ErrorWait)
	if _, err = d.DB.Exec(c, _ppDelUp, delay, mid); err != nil {
		log.Error("PostponeArc, failed to delay: (%v,%v), Error: %v", delay, mid, err)
	}
	return
}

// CountUpArcs counts the upper's archives
func (d *Dao) CountUpArcs(c context.Context, mid int64) (count int64, err error) {
	if err = d.DB.QueryRow(c, _upCountArc, mid).Scan(&count); err != nil {
		log.Error("d.CountUpArcs.Query error(%v)", err)
	}
	return
}

// UpArcs picks 50 arcs of the upper
func (d *Dao) UpArcs(c context.Context, mid int64) (aids []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(c, _upArcs, mid); err != nil { // get the qualified aid to sync
		log.Error("d.UpArcs.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("ParseVideos row.Scan() error(%v)", err)
			return
		}
		aids = append(aids, aid)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.UpArcs.Query error(%v)", err)
	}
	return
}
