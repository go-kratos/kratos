package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

var (
	_resLogByTidSQL = "SELECT id,oid,type,tid,mid,role,action,remark,state,ctime,mtime FROM resource_tag_log_%s WHERE oid=? AND `type`=? AND tid=? AND action = 0 ;"
)

// ResourceLogByTid return resource logs from mysql.
func (d *Dao) ResourceLogByTid(c context.Context, oid, tid int64, typ int32) (res []*model.ResourceLog, err error) {
	res = make([]*model.ResourceLog, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_resLogByTidSQL, d.hit(oid)), oid, typ, tid)
	if err != nil {
		log.Error("d.db.Query(%d,%d,%d) error(%v)", oid, tid, typ, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ResourceLog{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Action, &r.Remark, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

var (
	_resLogsSQL = "SELECT id,oid,type,tid,mid,role,action,remark,state,ctime,mtime FROM resource_tag_log_%s WHERE oid=? AND type=? AND state=0 AND role!=2 ORDER BY id DESC LIMIT ?,? ;"
)

// ResourceLogs return resource logs from mysql.
func (d *Dao) ResourceLogs(c context.Context, oid int64, typ int32, ps, pn int) (res []*model.ResourceLog, err error) {
	res = make([]*model.ResourceLog, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_resLogsSQL, d.hit(oid)), oid, typ, ps, pn)
	if err != nil {
		log.Error("d.db.Query(%d,%d,%d,%d) error(%v)", oid, typ, ps, pn, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ResourceLog{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Action, &r.Remark, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

var (
	_resLogsAdminSQL = "SELECT id,oid,type,tid,mid,role,action,remark,state,ctime,mtime FROM resource_tag_log_%s WHERE oid=? AND type=? AND state=0 ORDER BY id DESC LIMIT ?,? ;"
)

// ResourceLogsAdmin return resource logs from mysql include admin op.
func (d *Dao) ResourceLogsAdmin(c context.Context, oid int64, typ int32, ps, pn int) (res []*model.ResourceLog, err error) {
	res = make([]*model.ResourceLog, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_resLogsAdminSQL, d.hit(oid)), oid, typ, ps, pn)
	if err != nil {
		log.Error("d.db.Query(%d,%d,%d,%d) error(%v)", oid, typ, ps, pn, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ResourceLog{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Action, &r.Remark, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

var (
	_resLogSQL = "SELECT id,oid,type,tid,mid,role,action,remark,state,ctime,mtime FROM resource_tag_log_%s WHERE id=? and type = ? AND state = 0 ;"
)

// ResourceLog return resource logs from mysql.
func (d *Dao) ResourceLog(c context.Context, oid, logID int64, typ int32) (r *model.ResourceLog, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_resLogSQL, d.hit(oid)), logID, typ)
	r = &model.ResourceLog{}
	if err = row.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Action, &r.Remark, &r.State, &r.CTime, &r.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("rows.Scan() error(%v)", err)
		}
		r = nil
	}
	return
}

var (
	_addResLogSQL = "INSERT IGNORE INTO resource_tag_log_%s (oid,type,tid,tname,mid,role,action,remark,state) VALUES (?,?,?,?,?,?,?,?,?) ;"
)

// AddResourceLog add a resource log into mysql.
func (d *Dao) AddResourceLog(c context.Context, tname string, m *model.ResourceLog) (id int64, err error) {
	row, err := d.db.Exec(c, fmt.Sprintf(_addResLogSQL, d.hit(m.Oid)), m.Oid, m.Type, m.Tid, tname, m.Mid, m.Role, m.Action, m.Remark, m.State)
	if err != nil {
		log.Error("d.db.Exec(%v) error(%v)", m, err)
		return
	}
	return row.LastInsertId()
}
