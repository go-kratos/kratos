package ugc

import (
	"context"
	"fmt"

	"go-common/app/job/main/tv/model/ugc"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_delArc       = "UPDATE ugc_archive SET deleted = 1, submit = 1 WHERE aid = ? AND deleted = 0"
	_delVideos    = "UPDATE ugc_video SET deleted = 1, submit = 1 WHERE aid = ? AND deleted = 0"
	_delVideo     = "UPDATE ugc_video SET deleted = 1, submit = 1 WHERE cid = ? AND deleted = 0"
	_delVideoCids = "UPDATE ugc_video SET deleted = 1, submit = 1 WHERE cid IN (%s) AND deleted = 0"
	_checkVideos  = "SELECT cid FROM ugc_video WHERE aid = ? AND deleted = 0 LIMIT 1"
)

// TxDelArc deletes an arc
func (d *Dao) TxDelArc(tx *sql.Tx, aid int64) (err error) {
	if _, err = tx.Exec(_delArc, aid); err != nil {
		log.Error("TxDelArc, failed to update: (%v), Error: %v", aid, err)
	}
	return
}

// DelVideos delete one archive all videos
func (d *Dao) DelVideos(ctx context.Context, aid int64) (err error) {
	if _, err = d.DB.Exec(ctx, _delVideos, aid); err != nil {
		log.Error("DelVideos, failed to update: (%v), Error: %v", aid, err)
		return
	}
	log.Info("Aid %d is deleted, delete its videos", aid)
	return
}

// TxDelVideos deletes the videos of an arc
func (d *Dao) TxDelVideos(tx *sql.Tx, aid int64) (err error) {
	if _, err = tx.Exec(_delVideos, aid); err != nil {
		log.Error("TxDelVideos, failed to update: (%v), Error: %v", aid, err)
	}
	return
}

// TxDelVideo deletes a video
func (d *Dao) TxDelVideo(tx *sql.Tx, cid int64) (err error) {
	if _, err = tx.Exec(_delVideo, cid); err != nil {
		log.Error("TxDelVideo, failed to update: (%v), Error: %v", cid, err)
	}
	return
}

// DelVideoArc deletes some videos of an archive, if the archive is empty, also delete it
func (d *Dao) DelVideoArc(ctx context.Context, req *ugc.DelVideos) (arcValid bool, err error) {
	var cid int64
	arcValid = true
	if _, err = d.DB.Exec(ctx, fmt.Sprintf(_delVideoCids, xstr.JoinInts(req.CIDs))); err != nil {
		log.Error("DelVideos Cids %v, Aid %d, Error: %v", req.CIDs, req.AID, err)
		return
	}
	if err = d.DB.QueryRow(ctx, _checkVideos, req.AID).Scan(&cid); err != nil { // if no active videos, delete the arc
		if err == sql.ErrNoRows {
			err = nil
			arcValid = false
			if _, err = d.DB.Exec(ctx, _delArc, req.AID); err != nil {
				log.Error("DelVideos DelArc Cids %v, Aid %d, Error: %v", req.CIDs, req.AID, err)
				return
			}
			log.Info("DelArc Aid %d Because No Active Video", req.AID)
		} else {
			log.Error("DelVideos Cids %v, Aid %d, Error: %v", req.CIDs, req.AID, err)
			return
		}
	}
	return
}
