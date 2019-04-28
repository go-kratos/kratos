package audit

import (
	"context"
	"fmt"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"
)

const (
	_PgcCID       = "SELECT `id` FROM `tv_content` WHERE `cid` = ? AND `is_deleted` = 0"
	_UgcCID       = "SELECT `id` FROM `ugc_video` WHERE `cid` = ? AND `deleted` = 0"
	_transcodePGC = "UPDATE `tv_content` SET `transcoded` = ?, mark_time = NOW() WHERE `id` IN (%s) "
	_transcodeUGC = "UPDATE `ugc_video` SET `transcoded` = ? WHERE `id` IN (%s)"
	_applyPGC     = "UPDATE `tv_content` SET `apply_time` = ? WHERE `id` IN (%s)"
)

// checkCID picks one ugc or pgc video data with its cid
func (d *Dao) checkCID(c context.Context, cid int64, query string) (ids []int64, err error) {
	rows, err := d.db.Query(c, query, cid)
	if err != nil {
		log.Error("checkCID d.db.Query (%s) cid %d, error(%v)", query, cid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var li int64
		if err = rows.Scan(&li); err != nil {
			log.Error("checkCID Query (%s) row.Scan error(%v)", query, err)
			return
		}
		ids = append(ids, li)
	}
	if err = rows.Err(); err != nil {
		log.Error("checkCID Query (%s) rows Err cid %d, err %v", query, cid, err)
		return
	}
	if len(ids) == 0 {
		err = ecode.NothingFound
	}
	return
}

// UgcCID picks one ugc video data with its cid
func (d *Dao) UgcCID(c context.Context, cid int64) (ids []int64, err error) {
	return d.checkCID(c, cid, _UgcCID)
}

// PgcCID picks one ugc video data with its cid
func (d *Dao) PgcCID(c context.Context, cid int64) (ids []int64, err error) {
	return d.checkCID(c, cid, _PgcCID)
}

func (d *Dao) updateCIDs(c context.Context, query string, value interface{}) (err error) {
	if _, err = d.db.Exec(c, query, value); err != nil {
		log.Error("updateCIDs, Query %s, d.db.Exec.error(%v)", query, err)
	}
	return
}

// PgcTranscode updates the transcoded status of an ep data
func (d *Dao) PgcTranscode(c context.Context, ids []int64, action int64) (err error) {
	return d.updateCIDs(c, fmt.Sprintf(_transcodePGC, xstr.JoinInts(ids)), action)
}

// UgcTranscode updates the transcoded status of an ep data
func (d *Dao) UgcTranscode(c context.Context, ids []int64, action int64) (err error) {
	return d.updateCIDs(c, fmt.Sprintf(_transcodeUGC, xstr.JoinInts(ids)), action)
}

// ApplyPGC saves pgc apply_time; only PGC needs this
func (d *Dao) ApplyPGC(c context.Context, ids []int64, aTime int64) (err error) {
	return d.updateCIDs(c, fmt.Sprintf(_applyPGC, xstr.JoinInts(ids)), time.Time(aTime))
}
