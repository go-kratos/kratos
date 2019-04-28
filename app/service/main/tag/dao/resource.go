package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_addResTag = "INSERT IGNORE INTO resource_tag_%s (oid,type,tid,mid,role,enjoy,hate,attr,state) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?,role=?,state=0"
	_addTagRes = "INSERT IGNORE INTO tag_resource_%s (oid,type,tid,mid,role,enjoy,hate,attr,state) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?,role=?,state=0"
)

// TxAddResource add a resource into mysql.
func (d *Dao) TxAddResource(tx *sql.Tx, r *model.Resource) (int64, error) {
	res, err := tx.Exec(fmt.Sprintf(_addResTag, d.hit(r.Oid)), r.Oid, r.Type, r.Tid, r.Mid, r.Role, r.Like, r.Hate, r.Attr, r.State, r.Mid, r.Role)
	if err != nil {
		log.Error("tx.Exec(%v) error(%v)", r, err)
		return 0, err
	}
	_, err = tx.Exec(fmt.Sprintf(_addTagRes, d.hit(r.Tid)), r.Oid, r.Type, r.Tid, r.Mid, r.Role, r.Like, r.Hate, r.Attr, r.State, r.Mid, r.Role)
	if err != nil {
		log.Error("tx.Exec(%v) error(%v)", r, err)
		return 0, err
	}
	return res.RowsAffected()
}

var (
	_delResTag = "UPDATE resource_tag_%s SET state=1,mid=?,role=? WHERE oid=? AND type=? AND tid=?"
	_delTagRes = "UPDATE tag_resource_%s SET state=1,mid=?,role=? WHERE tid=? AND type=? AND oid=?"
)

// TxDelResource delete a resource.
func (d *Dao) TxDelResource(tx *sql.Tx, r *model.Resource) (int64, error) {
	res, err := tx.Exec(fmt.Sprintf(_delResTag, d.hit(r.Oid)), r.Mid, r.Role, r.Oid, r.Type, r.Tid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return 0, err
	}
	_, err = tx.Exec(fmt.Sprintf(_delTagRes, d.hit(r.Tid)), r.Mid, r.Role, r.Tid, r.Type, r.Oid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

var (
	_addDefaultResTag = "INSERT IGNORE INTO resource_tag_%s (oid,type,tid,mid,role,enjoy,hate,attr,state) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?,role=?,state=3"
	_addDefaultTagRes = "INSERT IGNORE INTO tag_resource_%s (oid,type,tid,mid,role,enjoy,hate,attr,state) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?,role=?,state=3"
)

// TxAddDefaultResource add default tag-res.
func (d *Dao) TxAddDefaultResource(tx *sql.Tx, r *model.Resource) (int64, error) {
	res, err := tx.Exec(fmt.Sprintf(_addDefaultResTag, d.hit(r.Oid)), r.Oid, r.Type, r.Tid, r.Mid, r.Role, r.Like, r.Hate, r.Attr, r.State, r.Mid, r.Role)
	if err != nil {
		log.Error("tx.Exec(%v) error(%v)", r, err)
		return 0, err
	}
	_, err = tx.Exec(fmt.Sprintf(_addDefaultTagRes, d.hit(r.Tid)), r.Oid, r.Type, r.Tid, r.Mid, r.Role, r.Like, r.Hate, r.Attr, r.State, r.Mid, r.Role)
	if err != nil {
		log.Error("tx.Exec(%v) error(%v)", r, err)
		return 0, err
	}
	return res.RowsAffected()
}

var (
	_resourceDefault = "SELECT id,oid,type,tid,mid,role,enjoy,hate,attr,state,ctime,mtime FROM resource_tag_%s WHERE oid=? AND type=? AND state=3"
)

// ResourceDefault return a resources by oid,type from msyql.
func (d *Dao) ResourceDefault(c context.Context, oid int64, typ int32) (res map[int64]*model.Resource, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_resourceDefault, d.hit(oid)), oid, typ)
	if err != nil {
		log.Error("d.ResourceDefault(%d,%d) error(%v)", oid, typ, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Resource)
	for rows.Next() {
		r := &model.Resource{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Like, &r.Hate, &r.Attr, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res[r.Tid] = r
	}
	return
}

var (
	_resource = "SELECT id,oid,type,tid,mid,role,enjoy,hate,attr,state,ctime,mtime FROM resource_tag_%s WHERE oid=? AND type=? AND tid=? AND state=0"
)

// Resource return a resources by tid from msyql.
func (d *Dao) Resource(c context.Context, oid, tid int64, typ int32) (r *model.Resource, err error) {
	r = new(model.Resource)
	row := d.db.QueryRow(c, fmt.Sprintf(_resource, d.hit(oid)), oid, typ, tid)
	if err = row.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Like, &r.Hate, &r.Attr, &r.State, &r.CTime, &r.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

var (
	_resources = "SELECT id,oid,type,tid,mid,role,enjoy,hate,attr,state,ctime,mtime FROM resource_tag_%s WHERE oid=? AND type=? AND state=0"
)

// Resources return resources by oid from mysql.
func (d *Dao) Resources(c context.Context, oid int64, typ int32) (res []*model.Resource, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_resources, d.hit(oid)), oid, typ)
	if err != nil {
		log.Error("d.oriTagsStmt(%d,%d) error(%v)", oid, typ, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Resource)
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Like, &r.Hate, &r.Attr, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// ResourceMap return resource map by oid from mysql.
func (d *Dao) ResourceMap(c context.Context, oid int64, typ int32) (res map[int64]*model.Resource, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_resources, d.hit(oid)), oid, typ)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Resource)
	for rows.Next() {
		r := &model.Resource{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Like, &r.Hate, &r.Attr, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res[r.Tid] = r
	}
	return
}

var (
	_resTagMapSQL = "SELECT id,oid,type,tid,mid,role,enjoy,hate,attr,state,ctime,mtime FROM resource_tag_%s WHERE oid=? AND type=?"
)

// ResTagMap .
func (d *Dao) ResTagMap(c context.Context, oid int64, tp int32) (res map[int64]*model.Resource, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_resTagMapSQL, d.hit(oid)), oid, tp)
	if err != nil {
		log.Error("d.dao.ResTagMap(%d,%d) error(%v)", oid, tp, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Resource)
	for rows.Next() {
		r := &model.Resource{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Like, &r.Hate, &r.Attr, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("d.dao.ResTagMap(%d,%d) rows.Scan() error(%v)", oid, tp, err)
			return
		}
		res[r.Tid] = r
	}
	err = rows.Err()
	return
}

var (
	_resTagsMapSQL = "SELECT id,oid,type,tid,mid,role,enjoy,hate,attr,state,ctime,mtime FROM resource_tag_%s WHERE oid in (%s) AND type=?"
)

// ResTagsMap .
func (d *Dao) ResTagsMap(c context.Context, oids []int64, tp int32) (res map[int64][]*model.Resource, err error) {
	res = make(map[int64][]*model.Resource, len(oids))
	for key, ids := range d.batchKey(oids) {
		var rows *sql.Rows
		if rows, err = d.db.Query(c, fmt.Sprintf(_resTagsMapSQL, key, xstr.JoinInts(ids)), tp); err != nil {
			log.Error("d.dao.ResTagsMap(%v,%d) error(%v)", oids, tp, err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			r := &model.Resource{}
			if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Like, &r.Hate, &r.Attr, &r.State, &r.CTime, &r.MTime); err != nil {
				log.Error("d.dao.ResTagsMap(%v,%d) rows.Scan() error(%v)", oids, tp, err)
				return
			}
			res[r.Oid] = append(res[r.Oid], r)
		}
		if err = rows.Err(); err != nil {
			return
		}
	}
	return
}

var (
	_resByTid = "SELECT oid,ctime FROM tag_resource_%s WHERE tid=? AND type=? AND state=0 ORDER BY id DESC LIMIT ?"
)

// ResOidsByTid return resource oids by tid from mysql.
func (d *Dao) ResOidsByTid(c context.Context, tid, limit int64, typ int32) (res []*model.Resource, oids []int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_resByTid, d.hit(tid)), tid, typ, limit)
	oids = make([]int64, 0)
	if err != nil {
		log.Error("d.db.Query(%d,%d) error(%v)", tid, typ, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Resource{}
		if err = rows.Scan(&r.Oid, &r.CTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
		oids = append(oids, r.Oid)
	}
	return
}

//ResourcesByTid return resource by tid from mysql
func (d *Dao) ResourcesByTid(c context.Context, tid, limit int64, typeID int32) (res []*model.Resource, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_resByTid, d.hit(tid)), tid, typeID, limit)
	if err != nil {
		log.Error("d.db.Query(%d,%d) error(%v)", tid, typeID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Resource{}
		if err = rows.Scan(&r.Oid, &r.CTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		if r != nil {
			res = append(res, r)
		}
	}
	return
}

var (
	_oidsByTid = "SELECT oid FROM tag_resource_%s WHERE tid=? AND type=? AND state=0"
)

//ResAllOidByTid get resource all oids by tid,type
func (d *Dao) ResAllOidByTid(c context.Context, tid int64, typeID int32) (oids []int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_oidsByTid, d.hit(tid)), tid, typeID)
	if err != nil {
		log.Error("d.db.Query(%d,%d) error(%v)", tid, typeID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var oid int64
		if err = rows.Scan(&oid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		oids = append(oids, oid)
	}
	return
}

var (
	_upResTagLike = "UPDATE resource_tag_%s SET enjoy=enjoy+? WHERE oid=? AND type=? AND tid=?"
	_upTagResLike = "UPDATE tag_resource_%s SET enjoy=enjoy+? WHERE tid=? AND type=? AND oid=?"
)

// TxUpResLike increase a resource like to mysql.
func (d *Dao) TxUpResLike(tx *sql.Tx, oid, tid int64, typ, incr int32) (int64, error) {
	res, err := tx.Exec(fmt.Sprintf(_upResTagLike, d.hit(oid)), incr, oid, typ, tid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return 0, err
	}
	_, err = tx.Exec(fmt.Sprintf(_upTagResLike, d.hit(tid)), incr, tid, typ, oid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

var (
	_upResTagHate = "UPDATE resource_tag_%s SET hate=hate+? WHERE oid=? AND type=? AND tid=?"
	_upTagResHate = "UPDATE tag_resource_%s SET hate=hate+? WHERE tid=? AND type=? AND oid=?"
)

// TxUpResHate increase a resource hate to mysql.
func (d *Dao) TxUpResHate(tx *sql.Tx, oid, tid int64, typ, incr int32) (int64, error) {
	res, err := tx.Exec(fmt.Sprintf(_upResTagHate, d.hit(oid)), incr, oid, typ, tid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return 0, err
	}
	_, err = tx.Exec(fmt.Sprintf(_upTagResHate, d.hit(tid)), incr, tid, typ, oid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

var (
	_upResTagAttr = "UPDATE resource_tag_%s SET attr=attr&(~(1<<?))|(?<<?) WHERE oid=? AND type=? AND tid=?"
	_upTagResAttr = "UPDATE tag_resource_%s SET attr=attr&(~(1<<?))|(?<<?) WHERE tid=? AND type=? AND oid=? "
)

// TxUpResAttr update resource attr to mysql.
func (d *Dao) TxUpResAttr(tx *sql.Tx, oid, tid int64, bit uint, val, typ int32) (int64, error) {
	res, err := tx.Exec(fmt.Sprintf(_upResTagAttr, d.hit(oid)), bit, val, bit, oid, typ, tid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return 0, err
	}
	_, err = tx.Exec(fmt.Sprintf(_upTagResAttr, d.hit(tid)), bit, val, bit, tid, typ, oid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}
