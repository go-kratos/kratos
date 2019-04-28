package archive

import (
	"context"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_oidFlowCount   = "SELECT COUNT(*) as count FROM flow_design WHERE state = 0 AND pool = ? AND group_id = ? AND oid = ?"
	_inFlowSQL      = "INSERT into flow_design(pool,oid,group_id,uid,remark)  VALUES (?,?,?,?,?)"
	_inFlowLogSQL   = "INSERT into flow_design_log(pool,oid,group_id,uid,action,remark)  VALUES (?,?,?,?,?,?)"
	_upFlowStateSQL = "UPDATE flow_design SET state=? WHERE id=?"
	_flowUniqueSQL  = "SELECT id,pool,oid,group_id,parent,state FROM flow_design WHERE oid=? AND pool=? AND group_id=? LIMIT 1"
)

// HasFlowGroup check if has flow group record
func (d *Dao) HasFlowGroup(c context.Context, pool int, gid, oid int64) (has bool, err error) {
	var (
		count int
	)
	row := d.db.QueryRow(c, _oidFlowCount, pool, gid, oid)
	if err = row.Scan(&count); err != nil {
		log.Error("d.hasFlowGroup err(%v)", err)
		return
	}
	has = count > 0
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

// TxUpFlowState 更新pool!=1的流量套餐资源的状态
// return int64, error/nil
func (d *Dao) TxUpFlowState(tx *sql.Tx, id int64, state int8) (rows int64, err error) {
	res, err := tx.Exec(_upFlowStateSQL, state, id)
	if err != nil {
		log.Error("TxUpFlowState.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// FlowUnique 获取命中 指定流量套餐的记录
// return *archive.FlowData/nil, error/nil
func (d *Dao) FlowUnique(c context.Context, oid, groupID int64, pool int8) (f *archive.FlowData, err error) {
	f = &archive.FlowData{}
	if err = d.db.QueryRow(context.TODO(), _flowUniqueSQL, oid, pool, groupID).Scan(&f.ID, &f.Pool, &f.OID, &f.GroupID, &f.Parent, &f.State); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			f = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
