package dao

import (
	"context"

	"go-common/app/interface/main/feedback/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inReply       = "INSERT INTO reply (session_id,reply_id,content,img_url,log_url,ctime,mtime) VALUES (?,?,?,?,?,?,?)"
	_selReply      = "SELECT reply_id,type,content,img_url,log_url,ctime FROM reply WHERE session_id=? ORDER BY id DESC LIMIT ?,?"
	_selReplyByMid = "SELECT r.reply_id,r.type,r.content,r.img_url,r.log_url,r.ctime FROM reply r INNER JOIN session s ON r.session_id=s.id WHERE s.mid=? ORDER BY r.id DESC LIMIT ?,?"
	_selReplyBySid = "SELECT r.reply_id,r.type,r.content,r.img_url,r.log_url,r.ctime FROM reply r INNER JOIN session s ON r.session_id=s.id WHERE r.session_id=? AND s.mid=? ORDER BY r.id"
)

// TxAddReply implements add a new reply record
func (d *Dao) TxAddReply(tx *sql.Tx, r *model.Reply) (id int64, err error) {
	res, err := tx.Exec(_inReply, r.SessionID, r.ReplyID, r.Content, r.ImgURL, r.LogURL, r.CTime, r.MTime)
	if err != nil {
		log.Error("AddReply tx.Exec() error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// AddReply insert reply.
func (d *Dao) AddReply(c context.Context, r *model.Reply) (id int64, err error) {
	res, err := d.inReply.Exec(c, r.SessionID, r.ReplyID, r.Content, r.ImgURL, r.LogURL, r.CTime, r.MTime)
	if err != nil {
		log.Error("AddReply error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// Replys returns corresponding user feedback reply records
func (d *Dao) Replys(c context.Context, ssnID int64, offset, limit int) (rs []model.Reply, err error) {
	rows, err := d.selReply.Query(c, ssnID, offset, limit)
	if err != nil {
		log.Error("d.selReply.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = model.Reply{}
		if err = rows.Scan(&r.ReplyID, &r.Type, &r.Content, &r.ImgURL, &r.LogURL, &r.CTime); err != nil {
			log.Error("row.Scan() error(%s)", err)
			rs = nil
			return
		}
		rs = append(rs, r)
	}
	return
}

// WebReplys get by ssnID.
func (d *Dao) WebReplys(c context.Context, ssnID, mid int64) (rs []*model.Reply, err error) {
	rows, err := d.selReplyBySid.Query(c, ssnID, mid)
	if err != nil {
		log.Error("d.selReply.Query(%d, %d) error(%v)", ssnID, mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.Reply{}
		if err = rows.Scan(&r.ReplyID, &r.Type, &r.Content, &r.ImgURL, &r.LogURL, &r.CTime); err != nil {
			log.Error("row.Scan() error(%s)", err)
			rs = nil
			return
		}
		rs = append(rs, r)
	}
	return
}

// ReplysByMid returns corresponding user feedback reply records by mid
func (d *Dao) ReplysByMid(c context.Context, mid int64, offset, limit int) (rs []model.Reply, err error) {
	rows, err := d.selReplyByMid.Query(c, mid, offset, limit)
	if err != nil {
		log.Error("d.selReplyByMid.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = model.Reply{}
		if err = rows.Scan(&r.ReplyID, &r.Type, &r.Content, &r.ImgURL, &r.LogURL, &r.CTime); err != nil {
			log.Error("row.Scan() error(%s)", err)
			rs = nil
			return
		}
		rs = append(rs, r)
	}
	return
}
