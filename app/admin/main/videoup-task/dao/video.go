package dao

import (
	"context"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_videoVID = `SELECT vr.id FROM archive_video_relation vr LEFT JOIN video v ON vr.cid=v.id
LEFT JOIN archive a ON vr.aid=a.id
WHERE vr.aid=? AND vr.cid=? AND vr.state != -100 AND v.status != -100 AND a.state != -100`
	_video = `SELECT vr.id, vr.aid, vr.cid, ar.mid, ar.copyright, ar.typeid, v.status, v.attribute, v.xcode_state, vr.title, vr.description, v.filename,
coalesce(ad.tid, 0) tid, coalesce(ad.reason, '') reason, coalesce(ao.remark, '') note
FROM archive_video_relation vr LEFT JOIN video v ON vr.cid=v.id
LEFT JOIN archive ar ON vr.aid=ar.id
LEFT JOIN archive_video_audit ad ON vr.id = ad.vid
LEFT JOIN archive_video_oper ao ON vr.id = ao.vid
WHERE vr.aid=? AND vr.cid=? AND ao.content NOT LIKE '%一审任务质检TAG: [%'
ORDER BY ao.id DESC LIMIT 0,1`
	_videoByCid = `SELECT vr.id,vr.aid,vr.title AS eptitle,vr.description,v.filename,v.src_type,vr.cid,v.duration,v.filesize,
v.resolutions,vr.index_order,vr.ctime,vr.mtime,v.status,v.playurl,v.attribute,v.failcode AS failinfo,v.xcode_state,v.weblink
 FROM archive_video_relation AS vr LEFT JOIN video AS v ON vr.cid = v.id WHERE vr.cid = ?`
	_newVideoIDSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
 avr.index_order,v.attribute,v.xcode_state,avr.state,avr.ctime,avr.mtime FROM archive_video_relation avr JOIN video v on avr.cid = v.id
 WHERE avr.id=? LIMIT 1`
	_videoAttributeSQL = `SELECT attribute FROM video WHERE id = ?`
)

// VideoAttribute get attr
func (d *Dao) VideoAttribute(ctx context.Context, cid int64) (attr int32, err error) {
	if err = d.arcDB.QueryRow(ctx, _videoAttributeSQL, cid).Scan(&attr); err != nil {
		if err == sql.ErrNoRows {
			attr = 0
			err = nil
		} else {
			PromeErr("arcdb: scan", "GetVVideoAttributeID row.Scan error(%v), cid(%d)", err, cid)
		}
	}
	return
}

//GetVID get vid
func (d *Dao) GetVID(ctx context.Context, aid int64, cid int64) (vid int64, err error) {
	if err = d.arcReadDB.QueryRow(ctx, _videoVID, aid, cid).Scan(&vid); err != nil {
		if err == sql.ErrNoRows {
			vid = 0
			err = nil
		} else {
			PromeErr("arcReaddb: scan", "GetVID row.Scan error(%v) aid(%d), cid(%d)", err, aid, cid)
		}
	}
	return
}

//Video get video by aid & cid
func (d *Dao) Video(ctx context.Context, aid int64, cid int64) (v *model.Video, err error) {
	v = &model.Video{}
	if err = d.arcReadDB.QueryRow(ctx, _video, aid, cid).Scan(&v.ID, &v.AID, &v.CID, &v.MID, &v.Copyright, &v.TypeID, &v.Status,
		&v.Attribute, &v.XcodeState, &v.Title, &v.Description, &v.Filename,
		&v.TagID, &v.Reason, &v.Note); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			v = nil
		} else {
			PromeErr("arcReaddb: scan", "Video row.Scan error(%v) aid(%d), cid(%d)", err, aid, cid)
		}
	}
	return
}

// ArcVideoByCID get video by cid
func (d *Dao) ArcVideoByCID(c context.Context, cid int64) (v *model.ArcVideo, err error) {
	row := d.arcDB.QueryRow(c, _videoByCid, cid)
	v = &model.ArcVideo{}
	if err = row.Scan(&v.ID, &v.Aid, &v.Title, &v.Desc, &v.Filename, &v.SrcType, &v.Cid, &v.Duration, &v.Filesize, &v.Resolutions, &v.Index, &v.CTime, &v.MTime, &v.Status, &v.Playurl, &v.Attribute, &v.FailCode, &v.XcodeState, &v.WebLink); err != nil {
		if err == sql.ErrNoRows {
			v = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// NewVideoByID .
func (d *Dao) NewVideoByID(c context.Context, id int64) (v *model.ArcVideo, err error) {
	row := d.arcDB.QueryRow(c, _newVideoIDSQL, id)
	v = &model.ArcVideo{}
	if err = row.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
		&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &v.Status, &v.CTime, &v.MTime); err != nil {
		if err == sql.ErrNoRows {
			v = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
