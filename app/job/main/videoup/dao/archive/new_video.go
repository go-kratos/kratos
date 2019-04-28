package archive

import (
	"context"
	"database/sql"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"go-common/app/job/main/videoup/model/archive"

	farm "github.com/dgryski/go-farm"
)

const (
	_upVideoXStateSQL     = "UPDATE video SET xcode_state=? WHERE hash64=? AND filename=?"
	_upVideoStatusSQL     = "UPDATE video SET status=? WHERE hash64=? AND filename=?"
	_upVideoPlayurlSQL    = "UPDATE video SET playurl=? WHERE hash64=? AND filename=?"
	_upVideoDuraSQL       = "UPDATE video SET duration=? WHERE hash64=? AND filename=?"
	_upVideoFilesizeSQL   = "UPDATE video SET filesize=? WHERE hash64=? AND filename=?"
	_upVideoResolutionSQL = "UPDATE video SET resolutions=?,dimensions=? WHERE hash64=? AND filename=?"
	_upVideoFailCodeSQL   = "UPDATE video SET failcode=? WHERE hash64=? AND filename=?"

	_upRelationStatusSQL = "UPDATE archive_video_relation SET state=? WHERE cid=?"

	_newVideoByFnSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,v.status,avr.ctime,avr.mtime FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE hash64=? AND filename=?`
	_newVideoByAidFnSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,v.status,avr.ctime,avr.mtime FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE aid=? AND hash64=? AND filename=?`
	_newVideosSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,v.status,avr.ctime,avr.mtime FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE aid=? ORDER BY index_order`
	_newVideoCntSQL  = `SELECT COUNT(*) FROM archive_video_relation WHERE aid=? AND state!=-100`
	_newSumDuraSQL   = `SELECT SUM(duration) FROM archive_video_relation avr JOIN video v on avr.cid = v.id WHERE aid=? AND avr.state=0 AND (v.status=0 || v.status=10000)`
	_newVdoBvcCntSQL = `SELECT COUNT(*) FROM archive_video_relation avr JOIN video v on avr.cid = v.id WHERE cid=? AND avr.state=0 AND (v.status=0 || v.status=10000) AND v.xcode_state=6`
	_validAidByCid   = `SELECT DISTINCT aid FROM archive_video_relation avr JOIN video v on avr.cid = v.id WHERE cid=? AND avr.state=0 AND (v.status=0 || v.status=10000) AND v.xcode_state=6`
)

// TxUpVideoXState update video xcodestate.
func (d *Dao) TxUpVideoXState(tx *xsql.Tx, filename string, xState int8) (rows int64, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	res, err := tx.Exec(_upVideoXStateSQL, xState, hash64, filename)
	if err != nil {
		log.Error("tx.upVideoXState.Exec(%d, %s) error(%v)", xState, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpVideoStatus update video status.
func (d *Dao) TxUpVideoStatus(tx *xsql.Tx, filename string, status int16) (rows int64, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	res, err := tx.Exec(_upVideoStatusSQL, status, hash64, filename)
	if err != nil {
		log.Error("tx.upVideoStatus.Exec(%d, %s) error(%v)", status, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpVideoPlayurl update video playurl and duration.
func (d *Dao) TxUpVideoPlayurl(tx *xsql.Tx, filename, playurl string) (rows int64, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	res, err := tx.Exec(_upVideoPlayurlSQL, playurl, hash64, filename)
	if err != nil {
		log.Error("tx.upVideoPlayurl.Exec(%s, %s) error(%v)", playurl, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpVDuration update video playurl and duration.
func (d *Dao) TxUpVDuration(tx *xsql.Tx, filename string, duration int64) (rows int64, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	res, err := tx.Exec(_upVideoDuraSQL, duration, hash64, filename)
	if err != nil {
		log.Error("tx.upVideoDura.Exec(%d, %s) error(%v)", duration, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpVideoFilesize update video filesize.
func (d *Dao) TxUpVideoFilesize(tx *xsql.Tx, filename string, filesize int64) (rows int64, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	res, err := tx.Exec(_upVideoFilesizeSQL, filesize, hash64, filename)
	if err != nil {
		log.Error("tx.upVideoFilesize.Exec(%d, %s) error(%v)", filesize, filename, err)
	}
	return res.RowsAffected()
}

// TxUpVideoResolutionsAndDimensions update video resolutions and dimensions.
func (d *Dao) TxUpVideoResolutionsAndDimensions(tx *xsql.Tx, filename, resolutions, dimensions string) (rows int64, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	res, err := tx.Exec(_upVideoResolutionSQL, resolutions, dimensions, hash64, filename)
	if err != nil {
		log.Error("tx.TxUpVideoResolutionsAndDimensions.Exec(%s,%s, %s) error(%v)", resolutions, dimensions, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpVideoFailCode update video fail info.
func (d *Dao) TxUpVideoFailCode(tx *xsql.Tx, filename string, fileCode int8) (rows int64, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	res, err := tx.Exec(_upVideoFailCodeSQL, fileCode, hash64, filename)
	if err != nil {
		log.Error("tx.upVideoFailCode.Exec(%s, %d) error(%v)", filename, fileCode, err)
		return
	}
	return res.RowsAffected()
}

// TxUpRelationStatus update video status.
func (d *Dao) TxUpRelationStatus(tx *xsql.Tx, cid int64, status int8) (rows int64, err error) {
	res, err := tx.Exec(_upRelationStatusSQL, status, cid)
	if err != nil {
		log.Error("tx.TxUpRelationStatus.Exec(%d, %d) error(%v)", status, cid, err)
		return
	}
	return res.RowsAffected()
}

// NewVideo get video info by filename.
func (d *Dao) NewVideo(c context.Context, filename string) (v *archive.Video, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	row := d.db.QueryRow(c, _newVideoByFnSQL, hash64, filename)
	v = &archive.Video{}
	var avrState, vState int16
	if err = row.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
		&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &avrState, &vState, &v.CTime, &v.MTime); err != nil {
		if err == sql.ErrNoRows {
			v = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	// 2 state map to 1
	if avrState == archive.VideoStatusDelete {
		v.Status = archive.VideoStatusDelete
	} else {
		v.Status = vState
	}
	return
}

// NewVideoByAid get video info by filename. and aid
func (d *Dao) NewVideoByAid(c context.Context, filename string, aid int64) (v *archive.Video, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	row := d.db.QueryRow(c, _newVideoByAidFnSQL, aid, hash64, filename)
	v = &archive.Video{}
	var avrState, vState int16
	if err = row.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
		&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &avrState, &vState, &v.CTime, &v.MTime); err != nil {
		if err == sql.ErrNoRows {
			v = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	// 2 state map to 1
	if avrState == archive.VideoStatusDelete {
		v.Status = archive.VideoStatusDelete
	} else {
		v.Status = vState
	}
	return
}

// NewVideos get videos info by aid.
func (d *Dao) NewVideos(c context.Context, aid int64) (vs []*archive.Video, err error) {
	rows, err := d.db.Query(c, _newVideosSQL, aid)
	if err != nil {
		log.Error("d.db.Query(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &archive.Video{}
		var avrState, vState int16
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

// NewVideoCount get all video duration by aid.
func (d *Dao) NewVideoCount(c context.Context, aid int64) (count int, err error) {
	row := d.db.QueryRow(c, _newVideoCntSQL, aid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// NewSumDuration get all video duration by aid.
func (d *Dao) NewSumDuration(c context.Context, aid int64) (sumDura int64, err error) {
	var (
		r   = &sql.NullInt64{}
		row = d.db.QueryRow(c, _newSumDuraSQL, aid)
	)
	if err = row.Scan(r); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	sumDura = r.Int64
	return
}

// NewVideoCountCapable get all video duration by aid.
func (d *Dao) NewVideoCountCapable(c context.Context, cid int64) (count int, err error) {
	row := d.db.QueryRow(c, _newVdoBvcCntSQL, cid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// ValidAidByCid get all video duration by aid.
func (d *Dao) ValidAidByCid(c context.Context, cid int64) (aids []int64, err error) {
	rows, err := d.db.Query(c, _validAidByCid, cid)
	if err != nil {
		log.Error("d.db.Query(%d) error(%v)", cid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		aids = append(aids, aid)
	}
	return
}
