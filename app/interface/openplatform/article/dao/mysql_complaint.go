package dao

import (
	"context"
	"database/sql"

	"go-common/library/log"
)

const (
	_addComplaintsSQL     = "INSERT INTO article_complaints(article_id,mid,type,reason,image_urls) VALUES (?,?,?,?,?)"
	_complaintExistSQL    = "SELECT id FROM article_complaints WHERE article_id=? AND mid=? AND state=0"
	_complaintProtectSQL  = "SELECT protect FROM article_complain_articles WHERE article_id=? AND deleted_time=0"
	_addComplaintCountSQL = "INSERT INTO article_complain_articles(article_id,count) VALUES (?,1) ON DUPLICATE KEY UPDATE count=count+1,state=0"

	_articleProtected = 1 // 0: no pretected  1: protected
)

// AddComplaint add complaint.
func (d *Dao) AddComplaint(c context.Context, aid, mid, ctype int64, reason, imageUrls string) (err error) {
	if _, err = d.addComplaintStmt.Exec(c, aid, mid, ctype, reason, imageUrls); err != nil {
		PromError("db:新增投诉")
		log.Error("dao.addComplaintStmt.exec(%s, %v, %v, %v, %v) error(%+v)", aid, mid, ctype, reason, imageUrls, err)
	}
	return
}

// ComplaintExist .
func (d *Dao) ComplaintExist(c context.Context, aid, mid int64) (exist bool, err error) {
	var id int
	if err = d.complaintExistStmt.QueryRow(c, aid, mid).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("d.complaintExistStmt.QueryRow(%d,%d) error(%+v)", aid, mid, err)
			PromError("db:判断之前是否投诉过")
		}
		return
	}
	exist = true
	return
}

// ComplaintProtected .
func (d *Dao) ComplaintProtected(c context.Context, aid int64) (protected bool, err error) {
	var p int
	if err = d.complaintProtectStmt.QueryRow(c, aid).Scan(&p); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("d.complaintProtectStmt.QueryRow(%d) error(%+v)", aid, err)
			PromError("db:判断文章是否被保护")
		}
		return
	}
	if p == _articleProtected {
		protected = true
	}
	return
}

// AddComplaintCount .
func (d *Dao) AddComplaintCount(c context.Context, aid int64) (err error) {
	if _, err = d.addComplaintCountStmt.Exec(c, aid); err != nil {
		log.Error("d.addComplaintCountStmt.Exec(%d) error(%+v)", aid, err)
		PromError("db:增加投诉计数")
	}
	return
}
