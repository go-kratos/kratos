package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inVideoOperSQL = "INSERT INTO archive_video_oper (aid,uid,vid,status,content,attribute,last_id,remark) VALUES (?,?,?,?,?,?,?,?)"
	_videoOperSQL   = "SELECT id,aid,uid,vid,status,content,attribute,last_id,remark,ctime FROM archive_video_oper WHERE vid = ? ORDER BY ctime DESC"
	_upVideoOperSQL = "UPDATE archive_video_oper SET last_id=? WHERE id=?"
)

// AddVideoOper insert archive_video_oper.
func (d *Dao) AddVideoOper(c context.Context, aid, adminID, vid int64, attribute int32, status int16, lastID int64, content, remark string) (id int64, err error) {
	res, err := d.arcDB.Exec(c, _inVideoOperSQL, aid, adminID, vid, status, content, attribute, lastID, remark)
	if err != nil {
		log.Error("d.inVideoOper.Exec error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}

//VideoOpers get video oper history list
func (d *Dao) VideoOpers(c context.Context, vid int64) (op []*model.VOper, uids []int64, err error) {
	var (
		rows  *sql.Rows
		ctime time.Time
	)
	op = []*model.VOper{}
	uids = []int64{}
	if rows, err = d.arcReadDB.Query(c, _videoOperSQL, vid); err != nil {
		log.Error("d.arcReadDB.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := &model.VOper{}
		if err = rows.Scan(&p.ID, &p.AID, &p.UID, &p.VID, &p.Status, &p.Content, &p.Attribute, &p.LastID, &p.Remark, &ctime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		p.CTime = ctime.Format("2006-01-02 15:04:05")
		op = append(op, p)
		uids = append(uids, p.UID)
	}
	return
}

// UpVideoOper update archive_video_oper last_id by id.
func (d *Dao) UpVideoOper(c context.Context, lastID, id int64) (rows int64, err error) {
	res, err := d.arcDB.Exec(c, _upVideoOperSQL, lastID, id)
	if err != nil {
		log.Error("d.upVideoOper.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
