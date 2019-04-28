package dao

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/admin/main/usersuit/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addIgnoreInviteSQL = "INSERT IGNORE INTO invite_code(mid,code,buy_ip,buy_ip_ng,expires,ctime) VALUES(?,?,?,?,?,?)"
	_getRangeInvitesSQL = "SELECT mid,imid,code,buy_ip,buy_ip_ng,expires,used_at,ctime FROM invite_code WHERE mid=? AND ctime>=? AND ctime<=?"
)

// AddIgnoreInvite add ignore invite code.
func (d *Dao) AddIgnoreInvite(c context.Context, inv *model.Invite) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _addIgnoreInviteSQL, inv.Mid, inv.Code, inv.IP, inv.IPng, inv.Expires, inv.Ctime); err != nil {
		log.Error("add invite, dao.db.Exec(%d, %s, %d, %v, %d, %v) error(%v)", inv.Mid, inv.Code, inv.IP, inv.IPng, inv.Expires, inv.Ctime, err)
		return
	}
	return res.RowsAffected()
}

// RangeInvites range invites.
func (d *Dao) RangeInvites(c context.Context, mid int64, start, end time.Time) (res []*model.Invite, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _getRangeInvitesSQL, mid, start, end); err != nil {
		log.Error("get range invites, dao.db.Query(%v, %v, %v) error(%v)", mid, start, end, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		inv := new(model.Invite)
		if err = rows.Scan(&inv.Mid, &inv.Imid, &inv.Code, &inv.IP, &inv.IPng, &inv.Expires, &inv.UsedAt, &inv.Ctime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res = append(res, inv)
	}
	err = rows.Err()
	return
}
