package reply

import (
	"context"
	"time"

	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_selAllNoticeSQL = `SELECT id,plat,condi,build,title,content,link,client_type  FROM notice WHERE stime <? and etime> ? and status=1 `
)

// NoticeDao notice dao.
type NoticeDao struct {
	db *sql.DB
}

// NewNoticeDao new a notice dao and return.
func NewNoticeDao(db *sql.DB) (dao *NoticeDao) {
	dao = &NoticeDao{
		db: db,
	}
	return
}

// ReplyNotice get reply notice infos from db
func (dao *NoticeDao) ReplyNotice(c context.Context) (nts []*reply.Notice, err error) {
	now := time.Now()
	rows, err := dao.db.Query(c, _selAllNoticeSQL, now, now)
	if err != nil {
		log.Error("dao.selAllResStmt query error (%v)", err)
		return
	}
	defer rows.Close()
	nts = make([]*reply.Notice, 0)
	for rows.Next() {
		nt := &reply.Notice{}
		if err = rows.Scan(&nt.ID, &nt.Plat, &nt.Condition, &nt.Build, &nt.Title, &nt.Content, &nt.Link, &nt.ClientType); err != nil {
			log.Error("rows.Scan err (%v)", err)
			return
		}
		nts = append(nts, nt)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}
