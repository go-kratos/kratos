package dao

import (
	"context"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_limitUsersSQL      = "SELECT id,mid,name,creator,ctime,mtime FROM limit_user ORDER BY ctime DESC LIMIT ?,?"
	_LimitUserCountSQL  = "SELECT count(*) FROM limit_user"
	_limitUserSQL       = "SELECT id,mid,name,creator,ctime,mtime FROM limit_user WHERE mid=?"
	_addLimitUserSQL    = "INSERT INTO limit_user(mid,name,creator) VALUES (?,?,?)"
	_delLimitUserSQL    = "DELETE FROM limit_user WHERE mid=?"
	_limitResByOidSQL   = "SELECT id,oid,type,author,operation,ctime,mtime FROM limit_resource WHERE oid=? AND type=?"
	_resLimitCountSQL   = "SELECT count(*) FROM limit_resource WHERE operation=?"
	_updateResLimitSQL  = "UPDATE limit_resource res SET res.operation=? WHERE res.oid=? AND res.type=?"
	_addLimitResSQL     = "INSERT INTO limit_resource(oid,type,operation) VALUES (?,?,?)"
	_resLimitByStateSQL = "SELECT id,oid,type,author,operation,ctime,mtime FROM limit_resource WHERE operation=? ORDER BY ctime DESC LIMIT ?,?;"
)

// LimitUsers limit users.
func (d *Dao) LimitUsers(c context.Context, start, end int32) (res []*model.LimitUser, err error) {
	rows, err := d.db.Query(c, _limitUsersSQL, start, end)
	if err != nil {
		log.Error("query limit user(%d,%d) error(%v)", start, end, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &model.LimitUser{}
		if err = rows.Scan(&u.ID, &u.Mid, &u.Name, &u.Creator, &u.CTime, &u.MTime); err != nil {
			log.Error("rows scan limit user error(%v)", err)
			return
		}
		res = append(res, u)
	}
	return
}

// LimitUserCount limit user count.
func (d *Dao) LimitUserCount(c context.Context) (count int64, err error) {
	row := d.db.QueryRow(c, _LimitUserCountSQL)
	if err = row.Scan(&count); err != nil {
		log.Error("cacu limit user count error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// LimitUser query limit user info by user mid.
func (d *Dao) LimitUser(c context.Context, mid int64) (res *model.LimitUser, err error) {
	res = new(model.LimitUser)
	row := d.db.QueryRow(c, _limitUserSQL, mid)
	if err = row.Scan(&res.ID, &res.Mid, &res.Name, &res.Creator, &res.CTime, &res.MTime); err != nil {
		log.Error("row.Scan limitUser (%d) error(%v)", mid, err)
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		}
	}
	return
}

// InsertLimitUser insert limit user.
func (d *Dao) InsertLimitUser(c context.Context, mid int64, name, cname string) (id int64, err error) {
	res, err := d.db.Exec(c, _addLimitUserSQL, mid, name, cname)
	if err != nil {
		log.Error("insert limit user(%d,%s,%s) error(%v)", mid, name, cname, err)
		return
	}
	return res.LastInsertId()
}

// DelLimitUser delete limit user by mid.
func (d *Dao) DelLimitUser(c context.Context, mid int64) (affect int64, err error) {
	res, err := d.db.Exec(c, _delLimitUserSQL, mid)
	if err != nil {
		log.Error("del limit user(%d) error(%v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// ResLimitByOid get resource limit by oid.
func (d *Dao) ResLimitByOid(c context.Context, oid int64, typ int32) (res *model.LimitRes, err error) {
	res = new(model.LimitRes)
	rows := d.db.QueryRow(c, _limitResByOidSQL, oid, typ)
	if err = rows.Scan(&res.ID, &res.Oid, &res.Type, &res.Author, &res.Operation, &res.CTime, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("rows.Scan() error(%v)", err)
		}
	}
	return
}

// ResLimitCount resource limit count.
func (d *Dao) ResLimitCount(c context.Context, state int32) (count int64, err error) {
	row := d.db.QueryRow(c, _resLimitCountSQL, state)
	if err = row.Scan(&count); err != nil {
		log.Error("CountNotice row.Scan err (%v)", err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// UpResLimitState update resource limit state by oid.
func (d *Dao) UpResLimitState(c context.Context, oid int64, tp int32, opera int32) (affect int64, err error) {
	res, err := d.db.Exec(c, _updateResLimitSQL, opera, oid, tp)
	if err != nil {
		log.Error("update limit res(%d,%d,%d) error(%v)", oid, tp, opera, err)
		return
	}
	return res.RowsAffected()
}

// ResLimitAdd add resource limit.
func (d *Dao) ResLimitAdd(c context.Context, oid int64, tp, operation int32) (id int64, err error) {
	res, err := d.db.Exec(c, _addLimitResSQL, oid, tp, operation)
	if err != nil {
		log.Error("insert reslimit(%d,%d,%d) error(%v)", oid, tp, operation, err)
		return
	}
	return res.LastInsertId()
}

// ResLimitByOpState get resource limit by Operation state.
func (d *Dao) ResLimitByOpState(c context.Context, state, start, end int32) (res []*model.LimitRes, oids []int64, err error) {
	rows, err := d.db.Query(c, _resLimitByStateSQL, state, start, end)
	if err != nil {
		log.Error("query limit resource by limitstate(%d,%d,%d) error(%v)", state, start, end, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.LimitRes{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Author, &r.Operation, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		oids = append(oids, r.Oid)
		res = append(res, r)
	}
	return
}
