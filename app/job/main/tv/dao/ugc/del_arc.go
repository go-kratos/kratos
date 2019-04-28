package ugc

import (
	"context"

	"go-common/library/log"
	"time"
)

const (
	_deletedArc   = "SELECT aid FROM ugc_archive WHERE submit = 1 AND deleted = 1 AND retry < unix_timestamp(now()) LIMIT 1"
	_finishDelArc = "UPDATE ugc_archive SET submit = 0 WHERE aid = ? AND deleted = 1"
	_ppDelArc     = "UPDATE ugc_archive SET retry = ? WHERE aid = ? AND deleted = 1"
)

// DeletedArc picks the deleted archive to sync
func (d *Dao) DeletedArc(c context.Context) (aid int64, err error) {
	err = d.DB.QueryRow(c, _deletedArc).Scan(&aid)
	return
}

// FinishDelArc updates the submit status from 1 to 0
func (d *Dao) FinishDelArc(c context.Context, aid int64) (err error) {
	if _, err = d.DB.Exec(c, _finishDelArc, aid); err != nil {
		log.Error("FinishVideos Error: %v", aid, err)
	}
	return
}

// PpDelArc postpones the archive's submit in 30 mins
func (d *Dao) PpDelArc(c context.Context, aid int64) (err error) {
	var delay = time.Now().Unix() + int64(d.conf.UgcSync.Frequency.ErrorWait)
	if _, err = d.DB.Exec(c, _ppDelArc, delay, aid); err != nil {
		log.Error("PostponeArc, failed to delay: (%v,%v), Error: %v", delay, aid, err)
	}
	return
}
