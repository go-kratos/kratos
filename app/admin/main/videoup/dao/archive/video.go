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
	_inVdoSQL = `INSERT INTO archive_video(filename,cid,aid,eptitle,description,src_type,duration,filesize,resolutions,playurl,failinfo,index_order,
	               attribute,xcode_state,status,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_upVdoSQL       = "UPDATE archive_video SET eptitle=?,description=? WHERE id=?"
	_upVdoIndexSQL  = "UPDATE archive_video SET index_order=? WHERE id=?"
	_upVdoLinkSQL   = "UPDATE archive_video SET weblink=? WHERE id=?"
	_upVdoStatusSQL = "UPDATE archive_video SET status=? WHERE id=?"
	_upVdoAttrSQL   = "UPDATE archive_video SET attribute=attribute&(~(1<<?))|(?<<?) WHERE id=?"
	_videoIDSQL     = `SELECT id,filename,cid,aid,eptitle,description,src_type,duration,filesize,resolutions,playurl,failinfo,
						index_order,attribute,xcode_state,status,ctime,mtime FROM archive_video WHERE id=? LIMIT 1`
	_videoIDsSQL = `SELECT id,filename,cid,aid,eptitle,description,src_type,duration,filesize,resolutions,playurl,failinfo,
						index_order,attribute,xcode_state,status,ctime,mtime FROM archive_video WHERE id in (%s)`
	_videoAidSQL = `SELECT id,filename,cid,aid,eptitle,description,src_type,duration,filesize,resolutions,playurl,failinfo,
						index_order,attribute,xcode_state,status,ctime,mtime FROM archive_video WHERE aid=? and status != -100 ORDER BY index_order ASC`
	_videoStatesSQL = "SELECT vr.id,vr.state AS vr_state,v.status AS v_status FROM archive_video_relation AS vr LEFT JOIN video AS v on vr.cid = v.id WHERE vr.id IN (%s)"
	_aidByVidsSQL   = "SELECT id,aid FROM archive_video_relation WHERE id IN (%s)"
)

// TxAddVideo insert video.
func (d *Dao) TxAddVideo(tx *sql.Tx, v *archive.Video) (vid int64, err error) {
	res, err := tx.Exec(_inVdoSQL, v.Filename, v.Cid, v.Aid, v.Title, v.Desc, v.SrcType, v.Duration, v.Filesize, v.Resolutions,
		v.Playurl, v.FailCode, v.Index, v.Attribute, v.XcodeState, v.Status, v.CTime, v.MTime)
	if err != nil {
		log.Error("d.inVideo.Exec error(%v)", err)
		return
	}
	vid, err = res.LastInsertId()
	return
}

// TxUpVideo update video by id.
func (d *Dao) TxUpVideo(tx *sql.Tx, vid int64, title, desc string) (rows int64, err error) {
	res, err := tx.Exec(_upVdoSQL, title, desc, vid)
	if err != nil {
		log.Error("d.upVideo.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVideoIndex update video index by id.
func (d *Dao) TxUpVideoIndex(tx *sql.Tx, vid int64, index int) (rows int64, err error) {
	res, err := tx.Exec(_upVdoIndexSQL, index, vid)
	if err != nil {
		log.Error("d.upVideoIndex.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVideoLink update weblink.
func (d *Dao) TxUpVideoLink(tx *sql.Tx, id int64, weblink string) (rows int64, err error) {
	res, err := tx.Exec(_upVdoLinkSQL, weblink, id)
	if err != nil {
		log.Error("d.upVideoLink.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVideoStatus update video status by id.
func (d *Dao) TxUpVideoStatus(tx *sql.Tx, id int64, status int16) (rows int64, err error) {
	res, err := tx.Exec(_upVdoStatusSQL, status, id)
	if err != nil {
		log.Error("d.upVideoStatus.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVideoAttr update video attribute by id.
func (d *Dao) TxUpVideoAttr(tx *sql.Tx, id int64, bit uint, val int32) (rows int64, err error) {
	res, err := tx.Exec(_upVdoAttrSQL, bit, val, bit, id)
	if err != nil {
		log.Error("d.upVideoAttr.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// VideoByID Video get video info by id. TODO Depreciated
func (d *Dao) VideoByID(c context.Context, id int64) (v *archive.Video, err error) {
	row := d.rddb.QueryRow(c, _videoIDSQL, id)
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

// VideoByIDs Video get video info by ids. TODO Depreciated
func (d *Dao) VideoByIDs(c context.Context, id []int64) (vs []*archive.Video, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_videoIDsSQL, xstr.JoinInts(id)))
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

// VideosByAid Video get video info by aid. TODO Depreciated
func (d *Dao) VideosByAid(c context.Context, aid int64) (vs []*archive.Video, err error) {
	rows, err := d.rddb.Query(c, _videoAidSQL, aid)
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

// VideoStateMap get archive id and state map
func (d *Dao) VideoStateMap(c context.Context, vids []int64) (sMap map[int64]int, err error) {
	sMap = make(map[int64]int)
	if len(vids) == 0 {
		return
	}
	rows, err := d.rddb.Query(c, fmt.Sprintf(_videoStatesSQL, xstr.JoinInts(vids)))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := struct {
			ID     int64
			State  int
			Status int
		}{}
		if err = rows.Scan(&a.ID, &a.State, &a.Status); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if a.State == -100 {
			sMap[a.ID] = -100
		} else {
			sMap[a.ID] = a.Status
		}
	}
	return
}

// VideoAidMap 批量通过视频id获取稿件id
func (d *Dao) VideoAidMap(c context.Context, vids []int64) (vMap map[int64]int64, err error) {
	var (
		aid, vid int64
	)
	vMap = make(map[int64]int64)
	if len(vids) == 0 {
		return
	}
	rows, err := d.rddb.Query(c, fmt.Sprintf(_aidByVidsSQL, xstr.JoinInts(vids)))
	defer rows.Close()
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	for rows.Next() {
		if err = rows.Scan(&vid, &aid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		vMap[vid] = aid
	}
	return
}
