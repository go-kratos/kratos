package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_relationsByTidSQL    = "SELECT r.id,r.oid,r.type,r.tid,r.mid,r.role,r.enjoy,r.hate,r.attr,r.state,r.ctime,r.mtime FROM tag_resource_%s as r WHERE r.tid=? AND r.type=3 AND r.state=0 ORDER BY r.ctime DESC LIMIT ?,?"
	_tagResCountSQL       = "SELECT count(*) FROM tag_resource_%s WHERE tid=? AND state=0"
	_resTagCountSQL       = "SELECT count(*) FROM resource_tag_%s WHERE oid=? AND type=? AND state=0"
	_relationsByOidSQL    = "SELECT r.id,r.oid,r.type,r.tid,r.mid,r.role,r.enjoy,r.hate,r.attr,r.state,r.ctime,r.mtime FROM resource_tag_%s as r WHERE r.oid=? AND r.type=? AND r.state=0 ORDER BY r.ctime DESC LIMIT ?,?"
	_relationSQL          = "SELECT id,oid,type,tid,state FROM resource_tag_%s WHERE oid=? AND tid=? AND type=? AND state=?"
	_insertByOidSQL       = "INSERT IGNORE INTO resource_tag_%s(oid,type,tid,mid,role,enjoy,hate,attr,state) VALUE(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE role=?,attr=?,state=?"
	_insertByTidSQL       = "INSERT IGNORE INTO tag_resource_%s(oid,type,tid,mid,role,enjoy,hate,attr,state) VALUE(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE role=?,attr=?,state=?"
	_updateAttrResTagSQL  = "UPDATE resource_tag_%s r SET r.attr=? WHERE oid=? AND tid=? AND type=?"
	_updateAttrTagResSQL  = "UPDATE tag_resource_%s r SET r.attr=? WHERE oid=? AND tid=? AND type=?"
	_updateStateResTagSQL = "UPDATE resource_tag_%s r SET r.state=? WHERE oid=? AND tid=? AND type=?"
	_updateStateTagResSQL = "UPDATE tag_resource_%s r SET r.state=? WHERE oid=? AND tid=? AND type=?"
)

// RelationsByTid RelationsByTid.
func (d *Dao) RelationsByTid(c context.Context, tid int64, start, end int32) (res []*model.Resource, oids []int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_relationsByTidSQL, d.hit(tid)), tid, start, end)
	if err != nil {
		log.Error("query relations by tid(%d,%d,%d) error(%v)", tid, start, end, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Resource{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Enjoy, &r.Hate, &r.Attr, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows scan(Relation{}) error(%v)", err)
			return
		}
		res = append(res, r)
		oids = append(oids, r.Oid)
	}
	return
}

// RelationsByOid RelationsByOid.
func (d *Dao) RelationsByOid(c context.Context, oid int64, tp, start, end int32) (res []*model.Resource, tids []int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_relationsByOidSQL, d.hit(oid)), oid, tp, start, end)
	if err != nil {
		log.Error("query relations by oid(%d,%d,%d,%d) error(%v)", oid, tp, start, end, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Resource{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Role, &r.Enjoy, &r.Hate, &r.Attr, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan(model.Relation{}) error(%v)", err)
			return
		}
		res = append(res, r)
		tids = append(tids, r.Tid)
	}
	return
}

// TagResCount count tag-resource.
func (d *Dao) TagResCount(c context.Context, tid int64) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_tagResCountSQL, d.hit(tid)), tid)
	if err = row.Scan(&count); err != nil {
		log.Error("query tag_res_count(%d) error (%v)", tid, err)
	}
	return
}

// ResTagCount count resource-tag.
func (d *Dao) ResTagCount(c context.Context, oid int64, tp int32) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_resTagCountSQL, d.hit(oid)), oid, tp)
	if err = row.Scan(&count); err != nil {
		log.Error("query res_tag_count(%d,%d) error (%v)", oid, tp, err)
	}
	return
}

// Relation Relation.
func (d *Dao) Relation(c context.Context, oid, tid int64, tp, state int32) (res *model.Resource, err error) {
	res = new(model.Resource)
	row := d.db.QueryRow(c, fmt.Sprintf(_relationSQL, d.hit(oid)), oid, tid, tp, state)
	if err = row.Scan(&res.ID, &res.Oid, &res.Type, &res.Tid, &res.State); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("query relation(%d,%d,%d,%d) error(%v)", oid, tid, tp, state, err)
		}
	}
	return
}

// TxInsertResTag tran insert  relation by oid.
func (d *Dao) TxInsertResTag(tx *sql.Tx, relation *model.Relation) (id int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_insertByOidSQL, d.hit(relation.Oid)),
		relation.Oid, relation.Type, relation.Tid, relation.Mid, relation.Role, relation.Enjoy, relation.Hate, relation.Attr, relation.State, relation.Role, relation.Attr, relation.State)
	if err != nil {
		log.Error("insert relation by oid(%d,%d,%d) error(%v)", relation.Oid, relation.Tid, relation.Type, err)
		return
	}
	return res.LastInsertId()
}

// TxInsertTagRes tran insert by tid.
func (d *Dao) TxInsertTagRes(tx *sql.Tx, relation *model.Relation) (id int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_insertByTidSQL, d.hit(relation.Tid)),
		relation.Oid, relation.Type, relation.Tid, relation.Mid, relation.Role, relation.Enjoy, relation.Hate, relation.Attr, relation.State, relation.Role, relation.Attr, relation.State)
	if err != nil {
		log.Error("insert relation by tid(%d,%d,%d) error(%v)", relation.Oid, relation.Tid, relation.Type, err)
		return
	}
	return res.LastInsertId()
}

// TxUpdateAttrResTag tran update attr.
func (d *Dao) TxUpdateAttrResTag(tx *sql.Tx, tid, oid int64, tp, attr int32) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateAttrResTagSQL, d.hit(oid)), attr, oid, tid, tp)
	if err != nil {
		log.Error("update res-tag attr(%d,%d,%d,%d) error(%v)", oid, tid, tp, attr, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateAttrTagRes tran update attr by tid.
func (d *Dao) TxUpdateAttrTagRes(tx *sql.Tx, tid, oid int64, tp, attr int32) (rowsCount int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateAttrTagResSQL, d.hit(tid)), attr, oid, tid, tp)
	if err != nil {
		log.Error("update tag-res attr(%d,%d,%d,%d) error(%v)", oid, tid, tp, attr, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateStateResTag tran delete resource-tag by oid.
func (d *Dao) TxUpdateStateResTag(tx *sql.Tx, tid, oid int64, tp, state int32) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateStateResTagSQL, d.hit(oid)), state, oid, tid, tp)
	if err != nil {
		log.Error("update res-tag state(%d,%d,%d,%d) error(%v)", oid, tid, tp, state, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateStateTagRes tran delete resource-tag by oid.
func (d *Dao) TxUpdateStateTagRes(tx *sql.Tx, tid, oid int64, tp, state int32) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateStateTagResSQL, d.hit(tid)), state, oid, tid, tp)
	if err != nil {
		log.Error("update tag-res state(%d,%d,%d,%d) error(%v)", oid, tid, tp, state, err)
		return
	}
	return res.RowsAffected()
}
