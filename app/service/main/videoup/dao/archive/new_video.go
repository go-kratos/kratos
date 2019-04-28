package archive

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	farm "github.com/dgryski/go-farm"
)

const (
	_inVideoCidSQL = `INSERT IGNORE INTO video (id,filename,src_type,resolutions,playurl,status,xcode_state,duration,filesize,attribute,failcode,hash64) 
					VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`
	_inNewVideoSQL = `INSERT INTO video (filename,src_type,resolutions,playurl,status,xcode_state,duration,filesize,attribute,failcode,hash64)
	VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	_inVideoRelationSQL = "INSERT IGNORE INTO archive_video_relation (id,aid,cid,title,description,index_order,ctime) VALUES (?,?,?,?,?,?,?)"
	_upVideoRelationSQL = "UPDATE archive_video_relation SET title=?,description=?,index_order=? ,state=? WHERE aid=? and cid=?"
	_upRelationStateSQL = "UPDATE archive_video_relation SET state=? WHERE aid=? AND cid=?"
	_upVideoStatusSQL   = "UPDATE video SET status=? WHERE id=?"
	_upNewVideoSQL      = "UPDATE video SET src_type=?,status=?,xcode_state=? WHERE id=?"
	_newVideoFnSQL      = "SELECT id,filename,src_type,resolutions,playurl,status,xcode_state,duration,filesize,attribute,failcode,ctime,mtime,dimensions FROM video WHERE hash64=? AND filename=?"
	_newVideoByFnSQL    = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,avr.ctime,avr.mtime,v.dimensions FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE hash64=? AND filename=?`
	_newVideoDataCidsFnSQL = "SELECT id,filename FROM video WHERE hash64 in (%s) AND filename in (%s)"
	_newsimpleArcVideoSQL  = `SELECT cid,title,index_order,state,mtime FROM archive_video_relation WHERE aid=?`
	_newVideosSQL          = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,v.status,avr.ctime,avr.mtime,v.dimensions FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE aid=? ORDER BY index_order`
	_newvideoCidSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,v.status,avr.ctime,avr.mtime,v.dimensions FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE cid=? ORDER BY id LIMIT 1`
	_newVideosCidSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,v.status,avr.ctime,avr.mtime,v.dimensions FROM archive_video_relation avr JOIN video v on avr.cid = v.id
		WHERE cid IN (%s)`
	_newVideosFnSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
		avr.index_order,v.attribute,v.xcode_state,avr.state,v.status,avr.ctime,avr.mtime,v.dimensions FROM archive_video_relation avr JOIN video v on avr.cid = v.id
        WHERE hash64 in (%s) AND filename in (%s)`
	_newVidReasonSQL     = `SELECT ava.vid,ava.reason FROM archive_video_audit ava LEFT JOIN archive_video_relation avr ON ava.vid=avr.id WHERE ava.aid=? AND avr.state!=-100`
	_newVideosTimeoutSQL = `SELECT  id ,filename,ctime,mtime from video WHERE hash64 in (%s) AND filename in (%s)`
)

// TxAddVideoCid  insert video to get cid.
func (d *Dao) TxAddVideoCid(tx *sql.Tx, v *archive.Video) (cid int64, err error) {
	hash64 := int64(farm.Hash64([]byte(v.Filename)))
	res, err := tx.Exec(_inVideoCidSQL, v.Cid, v.Filename, v.SrcType, v.Resolutions, v.Playurl, v.Status, v.XcodeState, v.Duration, v.Filesize, v.Attribute, v.FailCode, hash64)
	if err != nil {
		log.Error("d.inVideoCid.Exec error(%v)", err)
		return
	}
	cid, err = res.LastInsertId()
	return
}

// AddNewVideo insert new video.
func (d *Dao) AddNewVideo(c context.Context, v *archive.Video) (cid int64, err error) {
	hash64 := int64(farm.Hash64([]byte(v.Filename)))
	res, err := d.db.Exec(c, _inNewVideoSQL, v.Filename, v.SrcType, v.Resolutions, v.Playurl, v.Status, v.XcodeState, v.Duration, v.Filesize, v.Attribute, v.FailCode, hash64)
	if err != nil {
		log.Error("d.inNewVideo.Exec error(%v)", err)
		return
	}
	cid, err = res.LastInsertId()
	return
}

// TxAddNewVideo insert new video.
func (d *Dao) TxAddNewVideo(tx *sql.Tx, v *archive.Video) (cid int64, err error) {
	hash64 := int64(farm.Hash64([]byte(v.Filename)))
	res, err := tx.Exec(_inNewVideoSQL, v.Filename, v.SrcType, v.Resolutions, v.Playurl, v.Status, v.XcodeState, v.Duration, v.Filesize, v.Attribute, v.FailCode, hash64)
	if err != nil {
		log.Error("tx.inNewVideo.Exec error(%v)", err)
		return
	}
	cid, err = res.LastInsertId()
	return
}

// TxAddVideoRelation insert archive_video_relation to get vid.
func (d *Dao) TxAddVideoRelation(tx *sql.Tx, v *archive.Video) (vid int64, err error) {
	res, err := tx.Exec(_inVideoRelationSQL, v.ID, v.Aid, v.Cid, v.Title, v.Desc, v.Index, v.CTime)
	if err != nil {
		log.Error("d.inVideoRelation.Exec error(%v)", err)
		return
	}
	vid, err = res.LastInsertId()
	return
}

// TxUpVideoRelation update archive_video_relation info by aid and cid.
func (d *Dao) TxUpVideoRelation(tx *sql.Tx, v *archive.Video) (rows int64, err error) {
	res, err := tx.Exec(_upVideoRelationSQL, v.Title, v.Desc, v.Index, archive.VideoStatusOpen, v.Aid, v.Cid)
	if err != nil {
		log.Error("d.upVideoRelation.Exec(%v) error(%v)", v, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpRelationState update archive_video_relation state by aid and cid.
func (d *Dao) TxUpRelationState(tx *sql.Tx, aid, cid int64, state int16) (rows int64, err error) {
	res, err := tx.Exec(_upRelationStateSQL, state, aid, cid)
	if err != nil {
		log.Error("d.upRelationState.Exec(%d,%d,%d) error(%v)", aid, cid, state, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVdoStatus update video state by cid.
func (d *Dao) TxUpVdoStatus(tx *sql.Tx, cid int64, status int16) (rows int64, err error) {
	res, err := tx.Exec(_upVideoStatusSQL, status, cid)
	if err != nil {
		log.Error("d.upVideoStatus.Exec(%d,%d) error(%v)", cid, status, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpNewVideo  update video SrcType\Status\XcodeState by cid.
func (d *Dao) TxUpNewVideo(tx *sql.Tx, v *archive.Video) (rows int64, err error) {
	res, err := tx.Exec(_upNewVideoSQL, v.SrcType, v.Status, v.XcodeState, v.Cid)
	if err != nil {
		log.Error("d.upSimNewVideo.Exec(%s,%d,%d,%d) error(%v)", v.SrcType, v.Status, v.XcodeState, v.Cid, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// NewVideoFn get video by filename
func (d *Dao) NewVideoFn(c context.Context, filename string) (v *archive.Video, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	row := d.rddb.QueryRow(c, _newVideoFnSQL, hash64, filename)
	v = &archive.Video{}
	var dimStr string
	if err = row.Scan(&v.Cid, &v.Filename, &v.SrcType, &v.Resolutions, &v.Playurl, &v.Status, &v.XcodeState, &v.Duration, &v.Filesize, &v.Attribute, &v.FailCode, &v.CTime, &v.MTime, &dimStr); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			v = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	v.Dimension, _ = d.parseDimensions(dimStr)
	return
}

// NewVideoByFn get video by filename
func (d *Dao) NewVideoByFn(c context.Context, filename string) (v *archive.Video, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	row := d.rddb.QueryRow(c, _newVideoByFnSQL, hash64, filename)
	v = &archive.Video{}
	var dimStr string
	if err = row.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
		&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &v.Status, &v.CTime, &v.MTime, &dimStr); err != nil {
		if err == sql.ErrNoRows {
			v = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	v.Dimension, _ = d.parseDimensions(dimStr)
	return
}

// NewCidsByFns  get cids map in batches by filenames and hash64s.
func (d *Dao) NewCidsByFns(c context.Context, nvs []*archive.Video) (cids map[string]int64, err error) {
	var (
		buf     bytes.Buffer
		hash64s []int64
	)
	for _, v := range nvs {
		buf.WriteByte('\'')
		buf.WriteString(v.Filename)
		buf.WriteString("',")
		hash64s = append(hash64s, int64(farm.Hash64([]byte(v.Filename))))
	}
	buf.Truncate(buf.Len() - 1)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_newVideoDataCidsFnSQL, xstr.JoinInts(hash64s), buf.String()))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	cids = make(map[string]int64)
	for rows.Next() {
		var (
			cid      int64
			filename string
		)
		if err = rows.Scan(&cid, &filename); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		cids[filename] = cid
	}
	return
}

// SimpleArcVideos get simple videos from avr
func (d *Dao) SimpleArcVideos(c context.Context, aid int64) (vs []*archive.SimpleVideo, err error) {
	rows, err := d.rddb.Query(c, _newsimpleArcVideoSQL, aid)
	if err != nil {
		log.Error("d.videosStmt.Query(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &archive.SimpleVideo{}
		if err = rows.Scan(&v.Cid, &v.Title, &v.Index, &v.Status, &v.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		vs = append(vs, v)
	}
	return
}

// NewVideos get videos info by aid.
func (d *Dao) NewVideos(c context.Context, aid int64) (vs []*archive.Video, err error) {
	rows, err := d.rddb.Query(c, _newVideosSQL, aid)
	if err != nil {
		log.Error("d.videosStmt.Query(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &archive.Video{}
		var (
			avrState, vState int16
			dimStr           string
		)
		if err = rows.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
			&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &avrState, &vState, &v.CTime, &v.MTime, &dimStr); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		v.Dimension, _ = d.parseDimensions(dimStr)
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

// NewVideoMap get video map info by aid.
func (d *Dao) NewVideoMap(c context.Context, aid int64) (vm map[string]*archive.Video, cvm map[int64]*archive.Video, err error) {
	rows, err := d.rddb.Query(c, _newVideosSQL, aid)
	if err != nil {
		log.Error("d.videosStmt.Query(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	vm = make(map[string]*archive.Video)
	cvm = make(map[int64]*archive.Video)
	for rows.Next() {
		v := &archive.Video{}
		var (
			avrState, vState int16
			dimStr           string
		)
		if err = rows.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
			&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &avrState, &vState, &v.CTime, &v.MTime, &dimStr); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		v.Dimension, _ = d.parseDimensions(dimStr)
		// 2 state map to 1
		if avrState == archive.VideoStatusDelete {
			v.Status = archive.VideoStatusDelete
		} else {
			v.Status = vState
		}
		cvm[v.Cid] = v
		vm[v.Filename] = v
	}
	return
}

// NewVideoByCID get video by cid.
func (d *Dao) NewVideoByCID(c context.Context, cid int64) (v *archive.Video, err error) {
	row := d.rddb.QueryRow(c, _newvideoCidSQL, cid)
	v = &archive.Video{}
	var (
		avrState, vState int16
		dimStr           string
	)
	if err = row.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
		&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &avrState, &vState, &v.CTime, &v.MTime, &dimStr); err != nil {
		if err == sql.ErrNoRows {
			v = nil
			err = nil
		}
		log.Error("row.Scan error(%v)", err)
		return
	}
	v.Dimension, _ = d.parseDimensions(dimStr)
	// 2 state map to 1
	if avrState == archive.VideoStatusDelete {
		v.Status = archive.VideoStatusDelete
	} else {
		v.Status = vState
	}
	return
}

// NewVideosByCID multi get video by cids.
func (d *Dao) NewVideosByCID(c context.Context, cids []int64) (vm map[int64]map[int64]*archive.Video, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_newVideosCidSQL, xstr.JoinInts(cids)))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	vm = make(map[int64]map[int64]*archive.Video)
	for rows.Next() {
		var (
			avrState, vState int16
			dimStr           string
		)
		v := &archive.Video{}
		if err = rows.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
			&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &avrState, &vState, &v.CTime, &v.MTime, &dimStr); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		v.Dimension, _ = d.parseDimensions(dimStr)
		// 2 state map to 1
		if avrState == archive.VideoStatusDelete {
			v.Status = archive.VideoStatusDelete
		} else {
			v.Status = vState
		}
		if vv, ok := vm[v.Aid]; !ok {
			vm[v.Aid] = map[int64]*archive.Video{
				v.Cid: v,
			}
		} else {
			vv[v.Cid] = v
		}
	}
	return
}

// NewVideosByFn multi get video by filenames.
func (d *Dao) NewVideosByFn(c context.Context, fns []string) (vm map[int64]map[string]*archive.Video, err error) {
	var (
		buf     bytes.Buffer
		hash64s []int64
	)
	for _, fn := range fns {
		buf.WriteByte('\'')
		buf.WriteString(fn)
		buf.WriteString("',")
		hash64s = append(hash64s, int64(farm.Hash64([]byte(fn))))
	}
	buf.Truncate(buf.Len() - 1)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_newVideosFnSQL, xstr.JoinInts(hash64s), buf.String()))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	vm = make(map[int64]map[string]*archive.Video)
	for rows.Next() {
		var (
			avrState, vState int16
			dimStr           string
		)
		v := &archive.Video{}
		if err = rows.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
			&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &avrState, &vState, &v.CTime, &v.MTime, &dimStr); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		v.Dimension, _ = d.parseDimensions(dimStr)
		// 2 state map to 1
		if avrState == archive.VideoStatusDelete {
			v.Status = archive.VideoStatusDelete
		} else {
			v.Status = vState
		}
		if vv, ok := vm[v.Aid]; !ok {
			vm[v.Aid] = map[string]*archive.Video{
				v.Filename: v,
			}
		} else {
			vv[v.Filename] = v
		}
	}
	return
}

// CheckNewVideosTimeout check 48 timeout by add filenames.
func (d *Dao) CheckNewVideosTimeout(c context.Context, fns []string) (has bool, filename string, err error) {
	var (
		buf     bytes.Buffer
		hash64s []int64
	)
	for _, fn := range fns {
		buf.WriteByte('\'')
		buf.WriteString(fn)
		buf.WriteString("',")
		hash64s = append(hash64s, int64(farm.Hash64([]byte(fn))))
	}
	buf.Truncate(buf.Len() - 1)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_newVideosTimeoutSQL, xstr.JoinInts(hash64s), buf.String()))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	now := time.Now().Unix()
	for rows.Next() {
		v := &archive.VideoFn{}
		if err = rows.Scan(&v.Cid, &v.Filename, &v.CTime, &v.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if now-v.CTime.Time().Unix() > archive.VideoFilenameTimeout {
			log.Error("this video filename(%v) timeout (%+v)", v.Filename, v)
			has = true
			filename = v.Filename
			err = nil
			return
		}
	}
	return
}

// NewVideosReason get videos audit reason.
func (d *Dao) NewVideosReason(c context.Context, aid int64) (res map[int64]string, err error) {
	rows, err := d.rddb.Query(c, _newVidReasonSQL, aid)
	if err != nil {
		log.Error("d.vdoRsnStmt.Query(%d)|error(%v)", aid, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]string)
	for rows.Next() {
		var (
			vid    int64
			reason string
		)
		if err = rows.Scan(&vid, &reason); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res[vid] = reason
	}
	return
}

// parseDimensions 解析分辨率
func (d *Dao) parseDimensions(dim string) (dimensions *archive.Dimension, err error) {
	dimensions = &archive.Dimension{}
	if dim == "" || dim == "0,0,0" {
		return
	}
	dims, err := xstr.SplitInts(dim)
	if err != nil {
		log.Error("d.parseDimensions() xstr.SplitInts(%s) error(%v)", dim, err)
		return
	}
	if len(dims) != 3 {
		return
	}
	dimensions = &archive.Dimension{
		Width:  dims[0],
		Height: dims[1],
		Rotate: dims[2],
	}
	return
}
