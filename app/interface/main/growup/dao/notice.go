package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/library/log"

	"go-common/app/interface/main/growup/model"
)

const (
	_latestNoticeSQL = "SELECT title,link,type FROM notice WHERE status = 1 AND platform IN (?, 3) ORDER BY mtime DESC LIMIT 1"
	_noticeSQL       = "SELECT id,title,link,ctime,type FROM notice WHERE %s platform IN (?, 3) AND status=1 ORDER BY mtime DESC LIMIT ?,?"
	_noticeCountSQL  = "SELECT count(*) from notice WHERE %s platform IN (?, 3) AND status=1"
)

// LatestNotice latest notice
func (d *Dao) LatestNotice(c context.Context, platform int) (n *model.Notice, err error) {
	row := d.rddb.QueryRow(c, _latestNoticeSQL, platform)
	n = &model.Notice{}
	if err = row.Scan(&n.Title, &n.Link, &n.Type); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}

// Notices notices
func (d *Dao) Notices(c context.Context, typ string, platform int, offset, limit int) (notices []*model.Notice, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_noticeSQL, typ), platform, offset, limit)
	if err != nil {
		log.Error("d.db.Query Notices error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := &model.Notice{}
		err = rows.Scan(&n.ID, &n.Title, &n.Link, &n.CTime, &n.Type)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		notices = append(notices, n)
	}
	return
}

// NoticeCount notice count
func (d *Dao) NoticeCount(c context.Context, typ string, platform int) (count int64, err error) {
	row := d.rddb.QueryRow(c, fmt.Sprintf(_noticeCountSQL, typ), platform)
	if err = row.Scan(&count); err != nil {
		log.Error("row scan error(%v)", err)
	}
	return
}
