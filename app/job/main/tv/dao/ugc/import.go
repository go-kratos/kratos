package ugc

import (
	"context"
	"fmt"
	"time"

	ugcmdl "go-common/app/job/main/tv/model/ugc"
	"go-common/app/service/main/archive/api"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_import     = "SELECT mid FROM ugc_uploader WHERE toinit = 1 AND retry < UNIX_TIMESTAMP(now()) AND deleted = 0 LIMIT "
	_postponeUp = "UPDATE ugc_uploader SET retry = ? WHERE mid = ? AND deleted = 0"
	_finishUp   = "UPDATE ugc_uploader SET toinit = 0 WHERE mid = ? AND deleted = 0"
	_filterAids = "SELECT aid FROM ugc_archive WHERE aid IN (%s) AND deleted = 0"
	_importArc  = "REPLACE INTO ugc_archive(aid, videos, mid, typeid, title, cover, content, duration, copyright, pubtime, state) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
)

// TxImportArc imports an arc
func (d *Dao) TxImportArc(tx *sql.Tx, arc *api.Arc) (err error) {
	if _, err = tx.Exec(_importArc, arc.Aid,
		arc.Videos, arc.Author.Mid, arc.TypeID, arc.Title, arc.Pic, arc.Desc, arc.Duration,
		arc.Copyright, arc.PubDate, arc.State); err != nil {
		log.Error("_importArc, failed to update: (%v), Error: %v", arc, err)
	}
	return
}

// Import picks the uppers to init with the RPC data
func (d *Dao) Import(c context.Context) (res []*ugcmdl.Upper, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(c, _import+fmt.Sprintf("%d", d.conf.UgcSync.Batch.ImportNum)); err != nil {
		log.Error("d.Import.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &ugcmdl.Upper{}
		if err = rows.Scan(&r.MID); err != nil {
			log.Error("Manual row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.Import.Query error(%v)", err)
	}
	return
}

// PpUpper means postpone upper init operation due to some error happened
func (d *Dao) PpUpper(c context.Context, mid int64) (err error) {
	var delay = time.Now().Unix() + int64(d.conf.UgcSync.Frequency.ErrorWait)
	if _, err = d.DB.Exec(c, _postponeUp, delay, mid); err != nil {
		log.Error("PpUpper, failed to delay: (%v,%v), Error: %v", delay, mid, err)
	}
	return
}

// FinishUpper updates the upper's to_init status to 0 means we finish the import operation
func (d *Dao) FinishUpper(c context.Context, mid int64) (err error) {
	if _, err = d.DB.Exec(c, _finishUp, mid); err != nil {
		log.Error("FinishUpper, failed to Update: (%v,%v), Error: %v", _finishUp, mid, err)
	}
	return
}

// FilterExist filters the existing archives and remove them from the res, to have only non-existing data to insert
func (d *Dao) FilterExist(c context.Context, res *map[int64]*api.Arc, aids []int64) (err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(c, fmt.Sprintf(_filterAids, xstr.JoinInts(aids))); err != nil {
		if err == sql.ErrNoRows {
			err = nil // if non of them exist, it's good, we do nothing
			return
		}
		log.Error("d._filterAids.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aidEx int64
		if err = rows.Scan(&aidEx); err != nil {
			log.Error("Manual row.Scan() error(%v)", err)
			return
		}
		delete(*res, aidEx) // remove existing data from the map
	}
	if err = rows.Err(); err != nil {
		log.Error("d.FilterExist.Query error(%v)", err)
	}
	return
}
