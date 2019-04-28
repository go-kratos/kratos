package ugc

import (
	"context"
	"fmt"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

const (
	_videoCond = " AND (v.transcoded = 1 OR v.cid <= %d) " +
		"AND v.retry < unix_timestamp(now()) " +
		"AND v.deleted = 0 "
	_parseVideos = "SELECT id,cid,index_order,eptitle,duration,description FROM ugc_video v " +
		"WHERE v.aid = ? AND v.submit = 1 " + _videoCond
	_postArc      = "UPDATE ugc_video SET retry = ? WHERE cid = ? AND deleted = 0"
	_finishVideos = "UPDATE ugc_video SET submit = 0 WHERE cid = ? AND aid = ? AND deleted = 0"
	_finishArc    = "UPDATE ugc_archive SET submit = 0 WHERE aid = ? AND deleted = 0"
	_parseArc     = "SELECT id,aid,mid,typeid,videos,title,cover,content,duration,copyright,pubtime,ctime,mtime,state," +
		"manual,valid,submit,retry,result,deleted FROM ugc_archive WHERE aid = ?"
	_shouldAudit = "SELECT COUNT(1) as cnt FROM ugc_video v WHERE v.aid = ? AND v.submit = 1 " + _videoCond
	_videoSubmit = "SELECT cid FROM ugc_video v WHERE v.aid = ? AND v.submit = 0 " + _videoCond + " LIMIT 1"
)

// PpVideos postpones the archive's videos submit in 30 mins
func (d *Dao) PpVideos(c context.Context, cids []int64) (err error) {
	var delay = time.Now().Unix() + int64(d.conf.UgcSync.Frequency.ErrorWait)
	for _, v := range cids {
		if _, err = d.DB.Exec(c, _postArc, delay, v); err != nil {
			log.Error("PostponeArc, failed to delay: (%v,%v), Error: %v", delay, v, err)
			return
		}
	}
	return
}

// FinishVideos updates the submit status from 1 to 0
func (d *Dao) FinishVideos(c context.Context, videos []*ugcmdl.SimpleVideo, aid int64) (err error) {
	for _, v := range videos {
		if _, err = d.DB.Exec(c, _finishVideos, v.CID, aid); err != nil { // avoid updating the cid under another archive
			log.Error("FinishVideos Error: %v", v.CID, err)
			return
		}
	}
	if _, err = d.DB.Exec(c, _finishArc, aid); err != nil {
		log.Error("FinishVideos Error: %v", aid, err)
	}
	return
}

// ParseArc parses one archive data
func (d *Dao) ParseArc(c context.Context, aid int64) (res *ugcmdl.Archive, err error) {
	res = &ugcmdl.Archive{}
	if err = d.DB.QueryRow(c, _parseArc, aid).Scan(&res.ID, &res.AID, &res.MID, &res.TypeID, &res.Videos, &res.Title,
		&res.Cover, &res.Content, &res.Duration, &res.Copyright, &res.Pubtime, &res.Ctime, &res.Mtime, &res.State,
		&res.Manual, &res.Valid, &res.Submit, &res.Retry, &res.Result, &res.Deleted); err != nil { // get the qualified aid to sync
		log.Warn("d.ParseArc.Query error(%v)", err)
	}
	return
}

// ShouldAudit distinguishes whether the archive should ask for audit or not
func (d *Dao) ShouldAudit(c context.Context, aid int64) (res bool, err error) {
	var cnt int
	if err = d.DB.QueryRow(c, fmt.Sprintf(_shouldAudit, d.criCID), aid).Scan(&cnt); err != nil {
		log.Error("d.ShouldAudit Aid %d Err %v", aid, err)
		return
	}
	res = cnt > 0
	return
}

// VideoSubmit tells whether the archive already has some video submitted
func (d *Dao) VideoSubmit(c context.Context, aid int64) (cid int64, err error) {
	if err = d.DB.QueryRow(c, fmt.Sprintf(_videoSubmit, d.criCID), aid).Scan(&cid); err != nil {
		log.Warn("d.videoSubmit Aid %d, Err %v", aid, err)
	}
	return
}

// ParseVideos picks 20 videos of one qualified archive
func (d *Dao) ParseVideos(c context.Context, aid int64, ps int) (res [][]*ugcmdl.SimpleVideo, err error) {
	var (
		rows   *sql.Rows
		videos []*ugcmdl.SimpleVideo
	)
	if rows, err = d.DB.Query(c, fmt.Sprintf(_parseVideos, d.criCID), aid); err != nil {
		log.Error("d._parseVideos.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = ugcmdl.SimpleVideo{}
		if err = rows.Scan(&r.ID, &r.CID, &r.IndexOrder, &r.Eptitle, &r.Duration, &r.Description); err != nil {
			log.Error("ParseVideos row.Scan() error(%v)", err)
			return
		}
		videos = append(videos, &r)
		if len(videos) >= ps {
			var videoPce = append([]*ugcmdl.SimpleVideo{}, videos...)
			videos = []*ugcmdl.SimpleVideo{}
			res = append(res, videoPce)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("d._parseVideos.Query error(%v)", err)
		return
	}
	if len(videos) > 0 {
		res = append(res, videos)
	}
	return
}
