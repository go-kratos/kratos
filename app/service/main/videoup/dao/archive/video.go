package archive

import (
	"fmt"
	"strings"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// insert
	_inVideoSQL = `INSERT INTO archive_video (id,aid,eptitle,description,filename,src_type,cid,index_order,attribute,duration,filesize,resolutions,playurl,failinfo,xcode_state,status)
					VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_inAuditsSQL = "INSERT INTO archive_video_audit (vid,aid,tid,oname,reason) VALUES %s"
	// update
	_upVideoSQL     = `UPDATE archive_video SET eptitle=?,description=?,index_order=?,status=? WHERE id=?`
	_upVdoStatusSQL = "UPDATE archive_video SET status=? WHERE aid=? AND filename=?"
	_upVdoXcodeSQL  = "UPDATE archive_video SET xcode_state=? WHERE aid=? AND filename=?"
	_upVdoAttrSQL   = "UPDATE archive_video SET attribute=? WHERE aid=? AND filename=?"
	_upVdoCidSQL    = "UPDATE archive_video SET cid=? WHERE aid=? AND filename=?"
)

// TxAddVideo insert archive video.
func (d *Dao) TxAddVideo(tx *sql.Tx, v *archive.Video) (id int64, err error) {
	res, err := tx.Exec(_inVideoSQL, v.ID, v.Aid, v.Title, v.Desc, v.Filename, v.SrcType, v.Cid, v.Index, v.Attribute, v.Duration, v.Filesize, v.Resolutions, v.Playurl, v.FailCode, v.XcodeState, v.Status)
	if err != nil {
		log.Error("d.inVideo.Exec error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}

// TxUpVideo update video.
func (d *Dao) TxUpVideo(tx *sql.Tx, v *archive.Video) (rows int64, err error) {
	res, err := tx.Exec(_upVideoSQL, v.Title, v.Desc, v.Index, v.Status, v.ID)
	if err != nil {
		log.Error("d.upVideo.Exec(%v) error(%v)", v, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVideoStatus update video status.
func (d *Dao) TxUpVideoStatus(tx *sql.Tx, aid int64, filename string, status int16) (rows int64, err error) {
	res, err := tx.Exec(_upVdoStatusSQL, status, aid, filename)
	if err != nil {
		log.Error("d.upVideoStatus.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVideoXcode update video fail_code.
func (d *Dao) TxUpVideoXcode(tx *sql.Tx, aid int64, filename string, xCodeState int8) (rows int64, err error) {
	res, err := tx.Exec(_upVdoXcodeSQL, xCodeState, aid, filename)
	if err != nil {
		log.Error("d.upVdoXcode.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVideoAttr update video attribute.
func (d *Dao) TxUpVideoAttr(tx *sql.Tx, aid int64, filename string, attribute int32) (rows int64, err error) {
	res, err := tx.Exec(_upVdoAttrSQL, attribute, aid, filename)
	if err != nil {
		log.Error("d.upVideoAttr.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVideoCid update video attribute.
func (d *Dao) TxUpVideoCid(tx *sql.Tx, aid int64, filename string, cid int64) (rows int64, err error) {
	res, err := tx.Exec(_upVdoCidSQL, cid, aid, filename)
	if err != nil {
		log.Error("d.upVideoCid.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxAddAudit insert video audit.
func (d *Dao) TxAddAudit(tx *sql.Tx, vs []*archive.Video) (rows int64, err error) {
	var args = make([]string, 0, len(vs))
	for _, v := range vs {
		args = append(args, fmt.Sprintf(`(%d,%d,%d,'%s','%s')`, v.ID, v.Aid, 0, "videoup-service", ""))
	}
	res, err := tx.Exec(fmt.Sprintf(_inAuditsSQL, strings.Join(args, ",")))
	if err != nil {
		log.Error("d.inAudit.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
