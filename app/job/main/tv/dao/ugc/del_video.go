package ugc

import (
	"context"
	"time"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_deletedVideos  = "SELECT cid FROM ugc_video WHERE submit = 1 AND deleted = 1 AND retry < unix_timestamp(now()) LIMIT 5"
	_finishDelVideo = "UPDATE ugc_video SET submit = 0 WHERE cid = ? AND deleted = 1"
	_ppDelVideos    = "UPDATE ugc_video SET retry = ? WHERE cid = ? AND deleted = 1"
)

// DeletedVideos picks the deleted videos to sync
func (d *Dao) DeletedVideos(c context.Context) (delIds []int, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(c, _deletedVideos); err != nil { // get the qualified aid to sync
		return
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		if err = rows.Scan(&cid); err != nil {
			log.Error("ParseVideos row.Scan() error(%v)", err)
			return
		}
		delIds = append(delIds, cid)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.deletedVideos.Query error(%v)", err)
	}
	return
}

// FinishDelVideos updates the submit status from 1 to 0
func (d *Dao) FinishDelVideos(c context.Context, delIds []int) (err error) {
	for _, v := range delIds {
		if _, err = d.DB.Exec(c, _finishDelVideo, v); err != nil {
			log.Error("FinishDelVideos Error: %v", v, err)
			return
		}
	}
	return
}

// PpDelVideos postpones the archive's videos submit in 30 mins
func (d *Dao) PpDelVideos(c context.Context, delIds []int) (err error) {
	var delay = time.Now().Unix() + int64(d.conf.UgcSync.Frequency.ErrorWait)
	for _, v := range delIds {
		if _, err = d.DB.Exec(c, _ppDelVideos, delay, v); err != nil {
			log.Error("PpDelVideos, failed to delay: (%v,%v), Error: %v", delay, v, err)
			return
		}
	}
	return
}
