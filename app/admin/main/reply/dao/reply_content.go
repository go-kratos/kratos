package dao

import (
	"context"
	"fmt"
	"go-common/app/admin/main/reply/model"
	xsql "go-common/library/database/sql"
	"go-common/library/xstr"
	"time"
)

const _contSharding int64 = 200

const (
	_selReplyContentSQL  = "SELECT rpid,message,ats,ip,plat,device,ctime,mtime FROM reply_content_%d WHERE rpid=?"
	_selReplyContentsSQL = "SELECT rpid,message,ats,ip,plat,device,ctime,mtime FROM reply_content_%d WHERE rpid IN (%s)"
	_selContSQL          = "SELECT rpid,message,ats,ip,plat,device FROM reply_content_%d WHERE rpid=?"
	_selContsSQL         = "SELECT rpid,message,ats,ip,plat,device FROM reply_content_%d WHERE rpid IN (%s)"
	_upContMsgSQL        = "UPDATE reply_content_%d SET message=?,mtime=? WHERE rpid=?"
)

// UpReplyContent update reply content's message.
func (d *Dao) UpReplyContent(c context.Context, oid int64, rpID int64, msg string, now time.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upContMsgSQL, hit(oid)), msg, now, rpID)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// ReplyContent get a ReplyContent from database.
func (d *Dao) ReplyContent(c context.Context, oid, rpID int64) (rc *model.ReplyContent, err error) {
	rc = new(model.ReplyContent)
	row := d.db.QueryRow(c, fmt.Sprintf(_selReplyContentSQL, hit(oid)), rpID)
	if err = row.Scan(&rc.ID, &rc.Message, &rc.Ats, &rc.IP, &rc.Plat, &rc.Device, &rc.CTime, &rc.MTime); err != nil {
		if err == xsql.ErrNoRows {
			rc = nil
			err = nil
		}
	}
	return
}

// ReplyContents get reply contents by ids.
func (d *Dao) ReplyContents(c context.Context, oids []int64, rpIds []int64) (rcMap map[int64]*model.ReplyContent, err error) {
	hitMap := make(map[int64][]int64)
	for i, oid := range oids {
		hitMap[hit(oid)] = append(hitMap[hit(oid)], rpIds[i])
	}
	rcMap = make(map[int64]*model.ReplyContent, len(rpIds))

	for hit, ids := range hitMap {
		var rows *xsql.Rows
		rows, err = d.db.Query(c, fmt.Sprintf(_selReplyContentsSQL, hit, xstr.JoinInts(ids)))
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			rc := &model.ReplyContent{}
			if err = rows.Scan(&rc.ID, &rc.Message, &rc.Ats, &rc.IP, &rc.Plat, &rc.Device, &rc.CTime, &rc.MTime); err != nil {
				return
			}
			rcMap[rc.ID] = rc
		}
		if err = rows.Err(); err != nil {
			return
		}
	}
	return
}
