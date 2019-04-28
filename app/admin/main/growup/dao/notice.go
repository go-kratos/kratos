package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/growup/model"
	"go-common/library/log"
)

const (
	_inNoticeSQL     = "INSERT INTO notice(title,type,platform,link,status) VALUES (?,?,?,?,?)"
	_noticesSQL      = "SELECT id,title,type,platform,link,status FROM notice WHERE id > ? %s LIMIT ?"
	_noticeCountSQL  = "SELECT count(*) FROM notice WHERE id > 0 %s"
	_updateNoticeSQL = "UPDATE notice SET %s WHERE id=?"
)

// InsertNotice insert notice
func (d *Dao) InsertNotice(c context.Context, notice *model.Notice) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _inNoticeSQL, notice.Title, notice.Type, notice.Platform, notice.Link, notice.Status)
	if err != nil {
		log.Error("d.db.Exec insert notice error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// NoticeCount get notice count
func (d *Dao) NoticeCount(c context.Context, query string) (count int, err error) {
	row := d.rddb.QueryRow(c, fmt.Sprintf(_noticeCountSQL, query))
	if err = row.Scan(&count); err != nil {
		log.Error("d.db.notice count error(%v)", err)
	}
	return
}

// Notices get notices
func (d *Dao) Notices(c context.Context, query string, offset int, limit int) (notices []*model.Notice, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_noticesSQL, query), offset, limit)
	if err != nil {
		log.Error("d.db.Query notice error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := &model.Notice{}
		err = rows.Scan(&n.ID, &n.Title, &n.Type, &n.Platform, &n.Link, &n.Status)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		notices = append(notices, n)
	}
	return
}

// UpdateNotice update notice
func (d *Dao) UpdateNotice(c context.Context, kv string, id int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_updateNoticeSQL, kv), id)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
