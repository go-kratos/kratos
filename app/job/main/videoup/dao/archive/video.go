package archive

import (
	"context"
	"database/sql"

	"go-common/app/job/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// insert // NOTE: will delete???
	// _inFilenameSQL = `INSERT INTO archive_video (filename,filesize,xcode_state) VALUES(?,?,?)
	// 					ON DUPLICATE KEY UPDATE filesize=?,xcode_state=?`
	// update
	_upXStateSQL     = "UPDATE archive_video SET xcode_state=? WHERE filename=?"
	_upStatusSQL     = "UPDATE archive_video SET status=? WHERE filename=?"
	_upPlayurlSQL    = "UPDATE archive_video SET playurl=? WHERE filename=?"
	_upVDuraSQL      = "UPDATE archive_video SET duration=? WHERE filename=?"
	_upResolutionSQL = "UPDATE archive_video SET resolutions=? WHERE filename=?"
	_upFilesizeSQL   = "UPDATE archive_video SET filesize=? WHERE filename=?"
	_upFailCodeSQL   = "UPDATE archive_video SET failinfo=? WHERE filename=?"
	// select
	_videoSQL = `SELECT id,filename,cid,aid,eptitle,description,src_type,duration,filesize,resolutions,playurl,failinfo,
						index_order,attribute,xcode_state,status,ctime,mtime FROM archive_video WHERE filename=?`
	_videoByAidSQL = `SELECT id,filename,cid,aid,eptitle,description,src_type,duration,filesize,resolutions,playurl,failinfo,
						index_order,attribute,xcode_state,status,ctime,mtime FROM archive_video WHERE filename=? AND aid=?`
	_videosSQL = `SELECT id,filename,cid,aid,eptitle,description,src_type,duration,filesize,resolutions,playurl,failinfo,
						index_order,attribute,xcode_state,status,ctime,mtime FROM archive_video WHERE aid=? ORDER BY index_order`
	_videoCntSQL     = `SELECT COUNT(*) FROM archive_video WHERE aid=? AND status!=-100`
	_sumDuraSQL      = `SELECT SUM(duration) FROM archive_video WHERE aid=? AND (status=0 || status=10000)`
	_vdoBvcCntSQL    = `SELECT COUNT(*) FROM archive_video WHERE cid=? AND (status=0 || status=10000) AND xcode_state=6`
	_vdoAidBvcCntSQL = `SELECT COUNT(*) FROM archive_video av LEFT JOIN archive a ON av.aid=a.id WHERE 
					av.cid=? AND (av.status=0 || av.status=10000) AND av.xcode_state=6 AND (a.state>=0 || a.state=-6)`
)

// TxUpXcodeState update video state.
func (d *Dao) TxUpXcodeState(tx *xsql.Tx, filename string, xState int8) (rows int64, err error) {
	res, err := tx.Exec(_upXStateSQL, xState, filename)
	if err != nil {
		log.Error("tx.Exec(%d, %s) error(%v)", xState, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpStatus update video status.
func (d *Dao) TxUpStatus(tx *xsql.Tx, filename string, status int16) (rows int64, err error) {
	res, err := tx.Exec(_upStatusSQL, status, filename)
	if err != nil {
		log.Error("tx.Exec(%d, %s) error(%v)", status, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpPlayurl update video playurl and duration.
func (d *Dao) TxUpPlayurl(tx *xsql.Tx, filename, playurl string) (rows int64, err error) {
	res, err := tx.Exec(_upPlayurlSQL, playurl, filename)
	if err != nil {
		log.Error("tx.Exec(%s, %s) error(%v)", playurl, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpVideoDuration update video playurl and duration.
func (d *Dao) TxUpVideoDuration(tx *xsql.Tx, filename string, duration int64) (rows int64, err error) {
	res, err := tx.Exec(_upVDuraSQL, duration, filename)
	if err != nil {
		log.Error("tx.Exec(%d, %s) error(%v)", duration, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpFilesize update video filesize.
func (d *Dao) TxUpFilesize(tx *xsql.Tx, filename string, filesize int64) (rows int64, err error) {
	res, err := tx.Exec(_upFilesizeSQL, filesize, filename)
	if err != nil {
		log.Error("tx.Exec(%d, %s) error(%v)", filesize, filename, err)
	}
	return res.RowsAffected()
}

// TxUpResolutions update video resolutions.
func (d *Dao) TxUpResolutions(tx *xsql.Tx, filename, resolutions string) (rows int64, err error) {
	res, err := tx.Exec(_upResolutionSQL, resolutions, filename)
	if err != nil {
		log.Error("tx.Exec(%s, %s) error(%v)", resolutions, filename, err)
		return
	}
	return res.RowsAffected()
}

// TxUpFailCode update video fail info.
func (d *Dao) TxUpFailCode(tx *xsql.Tx, filename string, fileCode int8) (rows int64, err error) {
	res, err := tx.Exec(_upFailCodeSQL, fileCode, filename)
	if err != nil {
		log.Error("tx.Exec(%s, %d) error(%v)", filename, fileCode, err)
		return
	}
	return res.RowsAffected()
}

// Video get video info by filename. NOTE Deprecated
func (d *Dao) Video(c context.Context, filename string) (v *archive.Video, err error) {
	row := d.db.QueryRow(c, _videoSQL, filename)
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

// VideoByAid get video info by filename. and aid.  NOTE Deprecated
func (d *Dao) VideoByAid(c context.Context, filename string, aid int64) (v *archive.Video, err error) {
	row := d.db.QueryRow(c, _videoByAidSQL, filename, aid)
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

// Videos get videos info by aid. NOTE Deprecated
func (d *Dao) Videos(c context.Context, aid int64) (vs []*archive.Video, err error) {
	rows, err := d.db.Query(c, _videosSQL, aid)
	if err != nil {
		log.Error("d.db.Query(%d) error(%v)", aid, err)
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

// VideoCount get all video duration by aid. NOTE Deprecated
func (d *Dao) VideoCount(c context.Context, aid int64) (count int, err error) {
	row := d.db.QueryRow(c, _videoCntSQL, aid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// SumDuration get all video duration by aid. NOTE Deprecated
func (d *Dao) SumDuration(c context.Context, aid int64) (sumDura int64, err error) {
	var (
		r   = &sql.NullInt64{}
		row = d.db.QueryRow(c, _sumDuraSQL, aid)
	)
	if err = row.Scan(r); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	sumDura = r.Int64
	return
}

// VideoCountCapable get all video duration by aid. NOTE Deprecated
func (d *Dao) VideoCountCapable(c context.Context, cid int64) (count int, err error) {
	row := d.db.QueryRow(c, _vdoBvcCntSQL, cid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// VdoWithArcCntCapable get all video duration by aid. NOTE Deprecated
func (d *Dao) VdoWithArcCntCapable(c context.Context, cid int64) (count int, err error) {
	row := d.db.QueryRow(c, _vdoAidBvcCntSQL, cid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
