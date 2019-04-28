package ugc

import (
	"context"
	"fmt"
	"time"

	ugcmdl "go-common/app/job/main/tv/model/ugc"
	arccli "go-common/app/service/main/archive/api"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_manual       = "SELECT id,aid FROM ugc_archive WHERE manual = 1 AND retry < UNIX_TIMESTAMP(now()) AND deleted = 0 LIMIT "
	_postpone     = "UPDATE ugc_archive SET retry = ? WHERE aid = ? AND deleted = 0"
	_importFinish = "UPDATE ugc_archive SET manual = 0 WHERE aid = ? AND deleted = 0"
	_manualArc    = "UPDATE ugc_archive SET videos = ?, mid = ?, typeid = ?, title = ?, cover = ?, content = ?, duration = ?, " +
		"copyright = ?, pubtime = ?, state = ?, submit = ? WHERE aid = ? AND deleted = 0"
	_autoArc = "REPLACE INTO ugc_archive (videos, mid, typeid, title, cover, content, duration, copyright, pubtime, state, submit, aid) VALUES " +
		"(?,?,?,?,?,?,?,?,?,?,?,?)"
	_importVideo = "REPLACE INTO ugc_video (aid,cid,eptitle,index_order,duration,description) VALUES (?,?,?,?,?,?)"
)

// TxMnlArc updates the db with data from API
func (d *Dao) TxMnlArc(tx *sql.Tx, arc *ugcmdl.Archive) (err error) {
	if _, err = tx.Exec(_manualArc,
		arc.Videos, arc.MID, arc.TypeID, arc.Title, arc.Cover, arc.Content, arc.Duration,
		arc.Copyright, arc.Pubtime, arc.State, _needSubmit, arc.AID); err != nil {
		log.Error("_importArc, failed to update: (%v), Error: %v", arc, err)
	}
	return
}

// TxAutoArc imports the db an arc
func (d *Dao) TxAutoArc(tx *sql.Tx, arc *ugcmdl.Archive) (err error) {
	if _, err = tx.Exec(_autoArc,
		arc.Videos, arc.MID, arc.TypeID, arc.Title, arc.Cover, arc.Content, arc.Duration,
		arc.Copyright, arc.Pubtime, arc.State, _needSubmit, arc.AID); err != nil {
		log.Error("TxAutoArc, failed to update: (%v), Error: %v", arc, err)
	}
	return
}

// TxMnlVideos updates the db with data from API, if the
func (d *Dao) TxMnlVideos(tx *sql.Tx, view *arccli.ViewReply) (err error) {
	for _, v := range view.Pages {
		if _, err = tx.Exec(_importVideo, view.Arc.Aid, v.Cid, v.Part, v.Page, v.Duration, v.Desc); err != nil {
			log.Error("_importArc, failed to insert: (%v), Error: %v", v, err)
			return
		}
	}
	return
}

// TxMnlStatus updates the aid's manual status to 0
func (d *Dao) TxMnlStatus(tx *sql.Tx, aid int64) (err error) {
	if _, err = tx.Exec(_importFinish, aid); err != nil {
		log.Error("_importFinish, failed to update: (%v), Error: %v", aid, err)
	}
	return
}

// Manual picks the archives that added manually
func (d *Dao) Manual(c context.Context) (res []*ugcmdl.Archive, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(c, _manual+fmt.Sprintf("%d", d.conf.UgcSync.Batch.ManualNum)); err != nil {
		log.Error("d.Import.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &ugcmdl.Archive{}
		if err = rows.Scan(&r.ID, &r.AID); err != nil {
			log.Error("Manual row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.Manual.Query error(%v)", err)
	}
	return
}

// Ppmnl means postpone manual operation due to some error happened
func (d *Dao) Ppmnl(c context.Context, aid int64) (err error) {
	var delay = time.Now().Unix() + int64(d.conf.UgcSync.Frequency.ErrorWait)
	if _, err = d.DB.Exec(c, _postpone, delay, aid); err != nil {
		log.Error("Ppmnl, failed to delay: (%v,%v), Error: %v", delay, aid, err)
	}
	return
}
