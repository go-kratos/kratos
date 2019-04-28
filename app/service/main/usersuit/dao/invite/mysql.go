package dao

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/service/main/usersuit/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addInviteSQL       = "INSERT INTO invite_code(mid,code,buy_ip,buy_ip_ng,expires,ctime) VALUES(?,?,?,?,?,?)"
	_updateInviteSQL    = "UPDATE invite_code SET imid=?,used_at=? WHERE code=?"
	_getInviteSQL       = "SELECT mid,imid,code,expires,used_at,ctime FROM invite_code WHERE code=?"
	_getInvitesSQL      = "SELECT mid,imid,code,expires,used_at,ctime FROM invite_code WHERE mid=?"
	_getCurrentCountSQL = "SELECT count(1) FROM invite_code WHERE mid=? AND ctime>=? AND ctime<=?"
)

// Begin begin transaction
func (d *Dao) Begin(c context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(c)
}

// TxAddInvite transaction invite
func (d *Dao) TxAddInvite(c context.Context, tx *xsql.Tx, inv *model.Invite) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_addInviteSQL, inv.Mid, inv.Code, inv.IP, inv.IPng, inv.Expires, inv.Ctime); err != nil {
		log.Error("add invite, dao.db.Exec(%d, %s, %d, %v, %d, %v) error(%v)", inv.Mid, inv.Code, inv.IP, inv.IPng, inv.Expires, inv.Ctime, err)
		return
	}
	return res.RowsAffected()
}

// UpdateInvite update invite
func (d *Dao) UpdateInvite(c context.Context, imid int64, usedAt int64, code string) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateInviteSQL, imid, usedAt, code); err != nil {
		log.Error("update invite, dao.db.Exec(%d, %d, %s) error(%v)", imid, usedAt, code, err)
		return
	}
	return res.RowsAffected()
}

// Invite invite info
func (d *Dao) Invite(c context.Context, code string) (res *model.Invite, err error) {
	row := d.db.QueryRow(c, _getInviteSQL, code)
	res = new(model.Invite)
	if err = row.Scan(&res.Mid, &res.Imid, &res.Code, &res.Expires, &res.UsedAt, &res.Ctime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("get invite, row.Scan() error(%v)", err)
		}
	}
	return
}

// Invites invite list
func (d *Dao) Invites(c context.Context, mid int64) (res []*model.Invite, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _getInvitesSQL, mid); err != nil {
		log.Error("get invites, dao.db.Query(%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		inv := new(model.Invite)
		if err = rows.Scan(&inv.Mid, &inv.Imid, &inv.Code, &inv.Expires, &inv.UsedAt, &inv.Ctime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res = append(res, inv)
	}
	err = rows.Err()
	return
}

// CurrentCount current count
func (d *Dao) CurrentCount(c context.Context, mid int64, start, end time.Time) (res int64, err error) {
	row := d.db.QueryRow(c, _getCurrentCountSQL, mid, start, end)
	if err = row.Scan(&res); err != nil {
		log.Error("get current count, mid: %d, start: %d, end: %d, row.Scan() error(%v)", mid, start, end, err)
	}
	return
}
