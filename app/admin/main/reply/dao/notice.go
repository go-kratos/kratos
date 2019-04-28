package dao

import (
	"context"

	"go-common/app/admin/main/reply/model"
	"go-common/library/database/sql"
	xtime "go-common/library/time"
)

const (
	_selOneNotice          = "SELECT id,plat,version,condi,build,title,content,link,stime,etime,status,ctime,mtime,client_type FROM notice where id=?"
	_selCountNoticeSQL     = "SELECT count(*) as count FROM notice"
	_selAllNoticeSQL       = "SELECT id,plat,version,condi,build,title,content,link,stime,etime,status,ctime,mtime,client_type FROM notice ORDER BY stime DESC limit ?,?"
	_addNoticeSQL          = "INSERT INTO notice (plat,version,condi,build,title,content,link,status,stime,etime,client_type) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	_updateNotcieSQL       = "UPDATE notice set plat=?,version=?,condi=?,build=?,title=?,content=?,link=?,stime=?,etime=?,client_type=? where id=?"
	_updateNoticeStatusSQL = "UPDATE notice set status=? where id=?"
	_delNoticeSQL          = "DELETE FROM notice where id=?"
	_selOnlineNoticeSQL    = "SELECT id,plat,version,condi,build,title,content,link,stime,etime,status,ctime,mtime,client_type  FROM notice WHERE plat=? and (stime<=? and etime>?) and status=1"
)

// RangeNotice 获取发布时间与某一时间范围有交集的公告.
func (dao *Dao) RangeNotice(c context.Context, plat model.NoticePlat, stime xtime.Time, etime xtime.Time) (nts []*model.Notice, err error) {
	rows, err := dao.db.Query(c, _selOnlineNoticeSQL, plat, etime.Time(), stime.Time())
	if err != nil {
		return
	}
	defer rows.Close()
	nts = make([]*model.Notice, 0)
	for rows.Next() {
		nt := new(model.Notice)
		if err = rows.Scan(&nt.ID, &nt.Plat, &nt.Version, &nt.Condition, &nt.Build, &nt.Title, &nt.Content, &nt.Link, &nt.StartTime, &nt.EndTime, &nt.Status, &nt.CreateTime, &nt.ModifyTime, &nt.ClientType); err != nil {
			return
		}
		nts = append(nts, nt)
	}
	return
}

// Notice get one notice detail.
func (dao *Dao) Notice(c context.Context, id uint32) (nt *model.Notice, err error) {
	row := dao.db.QueryRow(c, _selOneNotice, id)
	nt = new(model.Notice)
	if err = row.Scan(&nt.ID, &nt.Plat, &nt.Version, &nt.Condition, &nt.Build, &nt.Title, &nt.Content, &nt.Link, &nt.StartTime, &nt.EndTime, &nt.Status, &nt.CreateTime, &nt.ModifyTime, &nt.ClientType); err != nil {
		if err == sql.ErrNoRows {
			nt = nil
			err = nil
		}
	}
	return
}

// CountNotice return notice count.
func (dao *Dao) CountNotice(c context.Context) (count int64, err error) {
	row := dao.db.QueryRow(c, _selCountNoticeSQL)
	err = row.Scan(&count)
	return
}

// ListNotice retrive reply's notice list from db by offset and count.
func (dao *Dao) ListNotice(c context.Context, offset int64, count int64) (nts []*model.Notice, err error) {
	rows, err := dao.db.Query(c, _selAllNoticeSQL, offset, count)
	if err != nil {
		return
	}
	defer rows.Close()
	nts = make([]*model.Notice, 0)
	for rows.Next() {
		nt := new(model.Notice)
		if err = rows.Scan(&nt.ID, &nt.Plat, &nt.Version, &nt.Condition, &nt.Build, &nt.Title, &nt.Content, &nt.Link, &nt.StartTime, &nt.EndTime, &nt.Status, &nt.CreateTime, &nt.ModifyTime, &nt.ClientType); err != nil {
			return
		}
		nts = append(nts, nt)
	}
	err = rows.Err()
	return
}

// CreateNotice insert a notice into db.
func (dao *Dao) CreateNotice(c context.Context, nt *model.Notice) (rows int64, err error) {
	res, err := dao.db.Exec(c, _addNoticeSQL, nt.Plat, nt.Version, nt.Condition, nt.Build, nt.Title, nt.Content, nt.Link, nt.Status, nt.StartTime, nt.EndTime, nt.ClientType)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

// UpdateNotice update main notice's main fileds.
func (dao *Dao) UpdateNotice(c context.Context, nt *model.Notice) (rows int64, err error) {
	res, err := dao.db.Exec(c, _updateNotcieSQL, nt.Plat, nt.Version, nt.Condition, nt.Build, nt.Title, nt.Content, nt.Link, nt.StartTime, nt.EndTime, nt.ClientType, nt.ID)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpdateNoticeStatus change notice's status to offline\online.
func (dao *Dao) UpdateNoticeStatus(c context.Context, status model.NoticeStatus, id uint32) (rows int64, err error) {
	res, err := dao.db.Exec(c, _updateNoticeStatusSQL, status, id)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// DeleteNotice delete a notice entry from db.
func (dao *Dao) DeleteNotice(c context.Context, id uint32) (rows int64, err error) {
	res, err := dao.db.Exec(c, _delNoticeSQL, id)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
