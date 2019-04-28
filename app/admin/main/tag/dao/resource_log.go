package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/tag/model"
	"go-common/library/log"
)

const (
	_resTagLogCountSQL = "SELECT count(*) FROM resource_tag_log_%s  WHERE oid=? AND type=? %s"
	_resTagLogsSQL     = "SELECT id,oid,mid,tid,tname,type,role,action,ctime,state FROM resource_tag_log_%s  WHERE oid=? AND type=? %s ORDER BY id DESC LIMIT ?,?"
	_upResLogStateSQL  = "UPDATE resource_tag_log_%s SET state=? WHERE oid=? AND type=? AND id=?;"
)

func resourceLogSQL(role, action int32) string {
	sqlStr := ""
	if role != model.ResRoleALL {
		sqlStr = sqlStr + fmt.Sprintf(" AND role=%d ", role)
	}
	if action != model.ResTagALL {
		sqlStr = sqlStr + fmt.Sprintf(" AND action=%d ", action)
	}
	return sqlStr
}

// ResTagLogCount ResTagLogCount.
func (d *Dao) ResTagLogCount(c context.Context, oid int64, tp, role, action int32) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_resTagLogCountSQL, d.hit(oid), resourceLogSQL(role, action)), oid, tp)
	if err = row.Scan(&count); err != nil {
		log.Error("query restaglog count(%d,%d,%d,%d) error(%v)", oid, tp, role, action, err)
	}
	return
}

// ResourceLogs ResourceLogs.
func (d *Dao) ResourceLogs(c context.Context, oid int64, tp, role, action, start, end int32) (res []*model.ResTagLog, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_resTagLogsSQL, d.hit(oid), resourceLogSQL(role, action)), oid, tp, start, end)
	if err != nil {
		log.Error("query resource log(%d,%d,%d,%d) error(%v)", oid, tp, role, action, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ResTagLog{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Mid, &r.Tid, &r.Tname, &r.Typ, &r.Role, &r.Action, &r.CTime, &r.State); err != nil {
			log.Error("scan resource log error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// UpdateResLogState update resource tag log state.
func (d *Dao) UpdateResLogState(c context.Context, id, oid int64, tp, state int32) (affect int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upResLogStateSQL, d.hit(oid)), state, oid, tp, id)
	if err != nil {
		log.Error("update res log state(%d,%d,%d,%d) error(%v)", state, oid, tp, id, err)
		return
	}
	return res.RowsAffected()
}
