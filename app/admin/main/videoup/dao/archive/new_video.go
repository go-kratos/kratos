package archive

import (
	"context"
	"fmt"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_videoByCid         = "SELECT vr.id,vr.aid,vr.title AS eptitle,vr.description,v.filename,v.src_type,vr.cid,v.duration,v.filesize,v.resolutions,vr.index_order,vr.ctime,vr.mtime,v.status,v.playurl,v.attribute,v.failcode AS failinfo,v.xcode_state,v.weblink FROM archive_video_relation AS vr LEFT JOIN video AS v ON vr.cid = v.id WHERE vr.cid = ?"
	_inRelationSQL      = "INSERT IGNORE INTO archive_video_relation (id,aid,cid,title,description,index_order,ctime,mtime) VALUES (?,?,?,?,?,?,?,?)"
	_upRelationSQL      = "UPDATE archive_video_relation SET title=?,description=? WHERE id=?"
	_upRelationOrderSQL = "UPDATE archive_video_relation SET index_order=? WHERE id=?"
	_upRelationStateSQL = "UPDATE archive_video_relation SET state=? WHERE id=?"
	_upVideoLinkSQL     = "UPDATE video SET weblink=? WHERE id=?"
	_upVideoStatusSQL   = "UPDATE video SET status=? WHERE id=?"
	_upVideoAttrSQL     = "UPDATE video SET attribute=attribute&(~(1<<?))|(?<<?) WHERE id=?"
	_slPlayurl          = "SELECT playurl FROM video WHERE id=? LIMIT 1"
	_newVideoIDSQL      = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,avr.ctime,avr.mtime FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE avr.id=? LIMIT 1`
	_newVideoIDsSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,avr.ctime,avr.mtime FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE avr.id in (%s)`
	_newVideosAIDSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,v.status,avr.ctime,avr.mtime FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE aid=? and state != -100 ORDER BY index_order ASC`
	_newVideoCntSQL = `SELECT COUNT(*) FROM archive_video_relation WHERE aid=? AND state!=-100`
	_slSrcTypeSQL   = "SELECT `id`, `src_type` FROM `video` WHERE `id` IN (%s)"
	_slVIDSQL       = "SELECT ar.id FROM archive_video_relation AS ar, video AS v WHERE ar.cid = v.id AND ar.aid=? AND v.filename=?;"
	_videoInfo      = `SELECT vr.id, vr.aid, vr.title AS eptitle, vr.description, vr.cid, vr.ctime AS epctime, v.filename, v.xcode_state, v.playurl,
	a.ctime, a.author, a.title, a.tag, a.content, a.cover, a.typeid, a.mid, a.copyright,
	coalesce(addit.source, '') source, coalesce(addit.dynamic, '') dynamic, coalesce(addit.desc_format_id, 0) desc_format_id, coalesce(addit.description, '') description
	FROM archive_video_relation AS vr JOIN archive AS a ON vr.aid = a.id
	LEFT OUTER JOIN video AS v ON vr.cid = v.id
	LEFT OUTER JOIN archive_addit AS addit ON vr.aid = addit.aid
	WHERE vr.aid = ? AND vr.cid=? LIMIT 1`
	_videoRelated = `SELECT v.filename,v.status,vr.aid,vr.index_order,a.title,a.ctime FROM archive_video_relation AS vr LEFT JOIN video AS v ON vr.cid = v.id JOIN archive AS a ON vr.aid = a.id WHERE vr.aid = ?`
)

// VideoByCID get video by cid
func (d *Dao) VideoByCID(c context.Context, cid int64) (v *archive.Video, err error) {
	row := d.rddb.QueryRow(c, _videoByCid, cid)
	v = &archive.Video{}
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

// TxAddRelation insert archive_video_relation.
func (d *Dao) TxAddRelation(tx *sql.Tx, v *archive.Video) (rows int64, err error) {
	res, err := tx.Exec(_inRelationSQL, v.ID, v.Aid, v.Cid, v.Title, v.Desc, v.Index, v.CTime, v.MTime)
	if err != nil {
		log.Error("d.inRelation.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpRelation update title and desc on archive_video_relation by vid.
func (d *Dao) TxUpRelation(tx *sql.Tx, vid int64, title, desc string) (rows int64, err error) {
	res, err := tx.Exec(_upRelationSQL, title, desc, vid)
	if err != nil {
		log.Error("d.upRelation.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpRelationOrder update index_order on archive_video_relation by vid.
func (d *Dao) TxUpRelationOrder(tx *sql.Tx, vid int64, index int) (rows int64, err error) {
	res, err := tx.Exec(_upRelationOrderSQL, index, vid)
	if err != nil {
		log.Error("d.upRelationOrder.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpRelationState update state on archive_video_relation  by vid.
func (d *Dao) TxUpRelationState(tx *sql.Tx, vid int64, state int16) (rows int64, err error) {
	res, err := tx.Exec(_upRelationStateSQL, state, vid)
	if err != nil {
		log.Error("d.upRelationState.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpWebLink update weblink on video by cid.
func (d *Dao) TxUpWebLink(tx *sql.Tx, cid int64, weblink string) (rows int64, err error) {
	res, err := tx.Exec(_upVideoLinkSQL, weblink, cid)
	if err != nil {
		log.Error("d.upVideoLink.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpStatus update status on video by cid.
func (d *Dao) TxUpStatus(tx *sql.Tx, cid int64, status int16) (rows int64, err error) {
	res, err := tx.Exec(_upVideoStatusSQL, status, cid)
	if err != nil {
		log.Error("d.upVideoStatus.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpAttr update attribute on video by cid.
func (d *Dao) TxUpAttr(tx *sql.Tx, cid int64, bit uint, val int32) (rows int64, err error) {
	res, err := tx.Exec(_upVideoAttrSQL, bit, val, bit, cid)
	if err != nil {
		log.Error("d.upVideoAttr.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// VideoPlayurl get video play url
func (d *Dao) VideoPlayurl(c context.Context, cid int64) (playurl string, err error) {
	row := d.rddb.QueryRow(c, _slPlayurl, cid)
	if err = row.Scan(&playurl); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// NewVideoByID Video get video info by id.
func (d *Dao) NewVideoByID(c context.Context, id int64) (v *archive.Video, err error) {
	row := d.rddb.QueryRow(c, _newVideoIDSQL, id)
	v = &archive.Video{}
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

// NewVideoByIDs Video get video info by ids. NOTE: NOT USED
func (d *Dao) NewVideoByIDs(c context.Context, id []int64) (vs []*archive.Video, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_newVideoIDsSQL, xstr.JoinInts(id)))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &archive.Video{}
		if err = rows.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
			&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &v.Status, &v.CTime, &v.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		vs = append(vs, v)
	}
	return
}

// NewVideosByAid Video get video info by aid.
func (d *Dao) NewVideosByAid(c context.Context, aid int64) (vs []*archive.Video, err error) {
	rows, err := d.rddb.Query(c, _newVideosAIDSQL, aid)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var avrState, vState int16
		v := &archive.Video{}
		if err = rows.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
			&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &avrState, &vState, &v.CTime, &v.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		// 2 state map to 1
		if avrState == archive.VideoStatusDelete {
			v.Status = archive.VideoStatusDelete
		} else {
			v.Status = vState
		}
		vs = append(vs, v)
	}
	return
}

// NewVideoCount get all video duration by aid. NOTE: NOT USED
func (d *Dao) NewVideoCount(c context.Context, aid int64) (count int, err error) {
	row := d.rddb.QueryRow(c, _newVideoCntSQL, aid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

//VideoSrcTypeByIDs video src_type and id map
func (d *Dao) VideoSrcTypeByIDs(c context.Context, ids []int64) (st map[int64]string, err error) {
	st = map[int64]string{}
	idStr := xstr.JoinInts(ids)
	rows, err := d.db.Query(c, fmt.Sprintf(_slSrcTypeSQL, idStr))
	if err != nil {
		log.Error("VideoSrcTypeByIDs d.db.Query (ids(%v)) error(%v)", idStr, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id      int64
			srcType string
		)
		if err = rows.Scan(&id, &srcType); err != nil {
			log.Error("VideoSrcTypeByIDs rows.Scan (ids(%v)) error(%v)", idStr, err)
			return
		}
		st[id] = srcType
	}
	return
}

//VIDByAIDFilename 根据filename查询视频的vid
func (d *Dao) VIDByAIDFilename(c context.Context, aid int64, filename string) (vid int64, err error) {
	row := d.db.QueryRow(c, _slVIDSQL, aid, filename)
	if err = row.Scan(&vid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("VideoRelationIDByFilename row.Scan err(%v) aid(%d) filename(%s)", err, aid, filename)
		}
	}

	return
}

//VideoInfo video info
func (d *Dao) VideoInfo(c context.Context, aid int64, cid int64) (v *archive.VideoInfo, err error) {
	var (
		descFormatID int64
		formatDesc   string
	)
	v = &archive.VideoInfo{}
	row := d.rddb.QueryRow(c, _videoInfo, aid, cid)
	if err = row.Scan(&v.ID, &v.AID, &v.Eptitle, &v.Description, &v.CID, &v.Epctime, &v.Filename, &v.XcodeState, &v.Playurl,
		&v.Ctime, &v.Author, &v.Title, &v.Tag, &v.Content, &v.Cover, &v.Typeid, &v.MID, &v.Copyright,
		&v.Source, &v.Dynamic, &descFormatID, &formatDesc); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			v = nil
		} else {
			log.Error("VideoInfo row.Scan error(%v) aid(%d) cid(%d)", err, aid, cid)
		}
		return
	}

	if descFormatID > 0 {
		v.Content = formatDesc
	}
	return
}

//VideoRelated related videos
func (d *Dao) VideoRelated(c context.Context, aid int64) (vs []*archive.RelationVideo, err error) {
	var rows *sql.Rows
	vs = []*archive.RelationVideo{}
	if rows, err = d.rddb.Query(c, _videoRelated, aid); err != nil {
		log.Error("VideoRelated d.rddb.Query error(%v) aid(%d)", err, aid)
		return
	}

	defer rows.Close()
	for rows.Next() {
		v := &archive.RelationVideo{}
		if err = rows.Scan(&v.Filename, &v.Status, &v.AID, &v.IndexOrder, &v.Title, &v.Ctime); err != nil {
			log.Error("VideoRelated rows.Scan error(%v) aid(%d)", err, aid)
			return
		}

		vs = append(vs, v)
	}
	return
}
