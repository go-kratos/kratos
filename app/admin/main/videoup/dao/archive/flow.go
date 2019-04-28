package archive

import (
	"context"
	"fmt"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_upFlowSQL             = "UPDATE flow_design SET group_id=?,uid=? WHERE id=?"
	_inFlowLogSQL          = "INSERT into flow_design_log(pool,oid,group_id,uid,action,remark)  VALUES (?,?,?,?,?,?)"
	_inFlowSQL             = "INSERT into flow_design(pool,oid,group_id,uid,remark)  VALUES (?,?,?,?,?)"
	_flowsSQL              = "SELECT id,name FROM flow_group WHERE  state=0"
	_flowPoolSQL           = "SELECT id FROM flow_design WHERE pool=? AND oid=? AND state=0 order by id desc limit 1"
	_findGroupIDByScopeSQL = "SELECT group_id FROM flow_scope WHERE  pool= ? AND industry_id=? AND brand_id=? AND official=? AND state=0  order by id desc limit 1;"
	_upFlowStateSQL        = "UPDATE flow_design SET state=? WHERE id=?"
	_flowsByOIDSQL         = "SELECT fd.id,fd.pool,fd.oid,fd.group_id,fd.parent,fd.state,fg.value FROM flow_design fd LEFT JOIN flow_group fg ON fd.group_id=fg.id WHERE fd.oid=? AND fd.state=0 AND fg.state=0"
	_flowUniqueSQL         = "SELECT id,pool,oid,group_id,parent,state FROM flow_design WHERE oid=? AND pool=? AND group_id=? LIMIT 1"
	_flowGroupPool         = "SELECT id, pool FROM flow_group WHERE id IN (%s)"
)

// TxUpFlow tx up flow_design.
func (d *Dao) TxUpFlow(tx *sql.Tx, flowID, groupID, UID int64) (rows int64, err error) {
	res, err := tx.Exec(_upFlowSQL, groupID, UID, flowID)
	if err != nil {
		log.Error("d.TxUpFlow.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxAddFlow tx add flow_design.
func (d *Dao) TxAddFlow(tx *sql.Tx, pool int8, oid, uid, groupID int64, remark string) (id int64, err error) {
	res, err := tx.Exec(_inFlowSQL, pool, oid, groupID, uid, remark)
	if err != nil {
		log.Error("d.TxAddFlow.Exec() error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}

//FindGroupIDByScope .
func (d *Dao) FindGroupIDByScope(c context.Context, pool int8, IndustryID, brandID int64, official int8) (groupID int64, err error) {
	row := d.rddb.QueryRow(c, _findGroupIDByScopeSQL, pool, IndustryID, brandID, official)
	if err = row.Scan(&groupID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		groupID = 1
		log.Info("FindGroupIDByScope match no scope AND hit default scope (%v)", groupID)
	}
	return
}

// TxAddFlowLog tx add flow_design log.
func (d *Dao) TxAddFlowLog(tx *sql.Tx, pool, action int8, oid, uid, groupID int64, remark string) (id int64, err error) {
	res, err := tx.Exec(_inFlowLogSQL, pool, oid, groupID, uid, action, remark)
	if err != nil {
		log.Error("d._inFlowLog.Exec() error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}

// Flows get flow_control id and remark.
func (d *Dao) Flows(c context.Context) (fs map[int64]string, err error) {
	rows, err := d.rddb.Query(c, _flowsSQL)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", err)
		return
	}
	defer rows.Close()
	fs = make(map[int64]string)
	for rows.Next() {
		var (
			id   int64
			name string
		)
		if err = rows.Scan(&id, &name); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		fs[id] = name
	}
	return
}

//FlowByPool .
func (d *Dao) FlowByPool(pool int8, oid int64) (id int64, err error) {
	row := d.rddb.QueryRow(context.TODO(), _flowPoolSQL, pool, oid)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

//TxUpFlowState 更新pool!=1的流量套餐资源的状态
func (d *Dao) TxUpFlowState(tx *sql.Tx, id int64, state int8) (rows int64, err error) {
	res, err := tx.Exec(_upFlowStateSQL, state, id)
	if err != nil {
		log.Error("TxUpFlowState.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

//FlowsByOID 获取所有命中的流量套餐记录
func (d *Dao) FlowsByOID(c context.Context, oid int64) (res []*archive.FlowData, err error) {
	var (
		rows *sql.Rows
	)
	res = []*archive.FlowData{}
	if rows, err = d.rddb.Query(context.TODO(), _flowsByOIDSQL, oid); err != nil {
		log.Error("FlowsByOID d.rddb.Query error(%v) oid(%d)", err, oid)
		return
	}
	defer rows.Close()

	for rows.Next() {
		f := &archive.FlowData{}
		if err = rows.Scan(&f.ID, &f.Pool, &f.OID, &f.GroupID, &f.Parent, &f.State, &f.GroupValue); err != nil {
			log.Error("FlowsByOID rows.Scan error(%v) oid(%d)", err, oid)
			return
		}
		res = append(res, f)
	}
	return
}

//FlowUnique 获取命中 指定流量套餐的记录
func (d *Dao) FlowUnique(c context.Context, oid, groupID int64, pool int8) (f *archive.FlowData, err error) {
	f = &archive.FlowData{}
	if err = d.rddb.QueryRow(context.TODO(), _flowUniqueSQL, oid, pool, groupID).Scan(&f.ID, &f.Pool, &f.OID, &f.GroupID, &f.Parent, &f.State); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			f = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

//FlowGroupPools 获取指定流量套餐的pool
func (d *Dao) FlowGroupPools(c context.Context, ids []int64) (res map[int64]int8, err error) {
	var (
		rows *sql.Rows
		id   int64
		pool int8
	)
	res = map[int64]int8{}
	idstr := xstr.JoinInts(ids)
	if rows, err = d.rddb.Query(context.TODO(), fmt.Sprintf(_flowGroupPool, idstr)); err != nil {
		log.Error("FlowGroupPools d.rddb.Query error(%v) ids(%s)", err, idstr)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&id, &pool); err != nil {
			log.Error("FlowGroupPools rows.Scan error(%v) ids(%d)", err, idstr)
			return
		}
		res[id] = pool
	}
	return
}
