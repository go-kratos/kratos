package archive

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	_inArcOperSQL     = "INSERT INTO archive_oper (aid,uid,typeid,state,content,round,attribute,last_id,remark) VALUES (?,?,?,?,?,?,?,?,?)"
	_inVideoOperSQL   = "INSERT INTO archive_video_oper (aid,uid,vid,status,content,attribute,last_id,remark) VALUES (?,?,?,?,?,?,?,?)"
	_upVideoOperSQL   = "UPDATE archive_video_oper SET last_id=? WHERE id=?"
	_arcOperSQL       = "SELECT id,aid,uid,typeid,state,content,round,attribute,last_id,remark FROM archive_oper WHERE aid = ? ORDER BY ctime DESC"
	_arcPassedOperSQL = "SELECT id FROM archive_oper WHERE aid=? AND state>=? LIMIT 1"
	_videoOperSQL     = "SELECT id,aid,uid,vid,status,content,attribute,last_id,remark,ctime FROM archive_video_oper WHERE vid = ? ORDER BY ctime DESC"
	_operAttrSQL      = "SELECT attribute, ctime  FROM archive_video_oper WHERE vid=? ORDER BY ctime DESC;"
)

// AddArcOper insert archive_oper.
func (d *Dao) AddArcOper(c context.Context, aid, adminID int64, attribute int32, typeID, state int16, round int8, lastID int64, content, remark string) (rows int64, err error) {
	res, err := d.db.Exec(c, _inArcOperSQL, aid, adminID, typeID, state, content, round, attribute, lastID, remark)
	if err != nil {
		log.Error("d.inArcOper.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// AddVideoOper insert archive_video_oper.
func (d *Dao) AddVideoOper(c context.Context, aid, adminID, vid int64, attribute int32, status int16, lastID int64, content, remark string) (id int64, err error) {
	res, err := d.db.Exec(c, _inVideoOperSQL, aid, adminID, vid, status, content, attribute, lastID, remark)
	if err != nil {
		log.Error("d.inVideoOper.Exec error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}

// UpVideoOper update archive_video_oper last_id by id.
func (d *Dao) UpVideoOper(c context.Context, lastID, id int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _upVideoOperSQL, lastID, id)
	if err != nil {
		log.Error("d.upVideoOper.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// ArchiveOper select archive_oper.
func (d *Dao) ArchiveOper(c context.Context, aid int64) (oper *archive.ArcOper, err error) {
	row := d.rddb.QueryRow(c, _arcOperSQL, aid)
	oper = &archive.ArcOper{}
	if err = row.Scan(&oper.ID, &oper.Aid, &oper.UID, &oper.TypeID, &oper.State, &oper.Content, &oper.Round, &oper.Attribute, &oper.LastID, &oper.Remark); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// VideoOper select archive_video_oper.
func (d *Dao) VideoOper(c context.Context, vid int64) (oper *archive.VideoOper, err error) {
	row := d.rddb.QueryRow(c, _videoOperSQL, vid)
	oper = &archive.VideoOper{}
	if err = row.Scan(&oper.ID, &oper.AID, &oper.UID, &oper.VID, &oper.Status, &oper.Content, &oper.Attribute, &oper.LastID, &oper.Remark, &oper.CTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// PassedOper check archive passed
func (d *Dao) PassedOper(c context.Context, aid int64) (id int64, err error) {
	row := d.rddb.QueryRow(c, _arcPassedOperSQL, aid, archive.StateOpen)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

//VideoOperAttrsCtimes 获取vid的审核属性记录，按照ctime排序
func (d *Dao) VideoOperAttrsCtimes(c context.Context, vid int64) (attrs []int32, ctimes []int64, err error) {
	rows, err := d.rddb.Query(c, _operAttrSQL, vid)
	if err != nil {
		log.Error("VideoOperAttrsCtimes d.rddb.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			ctime time.Time
			attr  int32
		)

		if err = rows.Scan(&attr, &ctime); err != nil {
			log.Error("VideoOperAttrsCtimes rows.Scan error(%v)", err)
			return
		}
		attrs = append(attrs, attr)
		ctimes = append(ctimes, int64(ctime))
	}
	return
}
