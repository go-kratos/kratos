package service

import (
	"context"
	"strings"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

/**
 * txAddFlow 新增流量套餐的记录
 * return int64, error/nil
 */
func (s *Service) txAddFlow(tx *sql.Tx, pool int8, oid, groupID, uid int64, remark string) (id int64, err error) {
	if id, err = s.arc.TxAddFlow(tx, pool, oid, uid, groupID, remark); err != nil {
		log.Error("txAddFlow s.arc.TxAddFlow(%d,%d,%d,%d,%s) error(%v)", pool, oid, uid, groupID, remark, err)
		return
	}
	if id <= 0 {
		return
	}
	if _, err = s.arc.TxAddFlowLog(tx, pool, archive.FlowLogAdd, oid, uid, groupID, remark); err != nil {
		log.Error("txAddFlow s.arc.TxAddFlowLog(%d,%d,%d,%d,%s) error(%v)", pool, oid, uid, groupID, remark, err)
		return
	}
	return
}

/**
 * txUpFlowState 更新流量套餐的状态
 * return error/nil
 */
func (s *Service) txUpFlowState(tx *sql.Tx, state int8, uid int64, f *archive.FlowData) (err error) {
	if f == nil {
		return
	}
	var rows int64
	if rows, err = s.arc.TxUpFlowState(tx, f.ID, state); err != nil {
		log.Error("updateFlowState s.arc.TxUpFlowState error(%v) id(%d) state(%d)", err, f.ID, state)
		return
	}
	if rows <= 0 {
		return
	}

	action := archive.FlowLogUpdate
	if state == archive.FlowDelete {
		action = archive.FlowLogDel
	}
	if _, err = s.arc.TxAddFlowLog(tx, f.Pool, action, f.OID, uid, f.GroupID, "审核后台修改状态"); err != nil {
		log.Error("updateFlowState s.arc.TxAddFlowLog error(%v) pool(%d) oid(%d) uid(%d) state(%d)", err, f.Pool, f.OID, uid, state)
		return
	}
	return
}

/**
 * txAddOrUpdateFlowState 新增或更新流量套餐的状态
 * return *archive.FlowData/nil, bool, error/nil
 */
func (s *Service) txAddOrUpdateFlowState(c context.Context, tx *sql.Tx, oid, groupID, uid int64, pool, state int8, remark string) (flow *archive.FlowData, diff string, err error) {
	var (
		old, nw int8
	)

	defer func() {
		if err == nil && old != nw {
			tagID := archive.FlowOperType[groupID]
			if tagID > 0 {
				stateMap := map[int8]string{archive.FlowOpen: "是", archive.FlowDelete: "否"}
				diff = strings.TrimSpace(archive.Operformat(tagID, stateMap[old], stateMap[nw], archive.OperStyleOne))
			}
		}
	}()

	if flow, err = s.arc.FlowUnique(c, oid, groupID, pool); err != nil {
		log.Error("txAddOrUpdateFlowState s.arc.FlowUnique(%d,%d,%d) error(%v) state(%d)", oid, groupID, pool, err, state)
		return
	}
	//无数据前提下，新状态=state就没必要添加数据啦
	if flow == nil && state == archive.FlowDelete {
		return
	}
	if flow == nil {
		flow = &archive.FlowData{Pool: pool, OID: oid, GroupID: groupID, State: archive.FlowOpen}
		if flow.ID, err = s.txAddFlow(tx, flow.Pool, flow.OID, flow.GroupID, uid, remark); err != nil {
			log.Error("txAddOrUpdateFlowState s.txAddFlow error(%v) flow(%+v) state(%d)", err, flow, state)
			return
		}
		old = archive.FlowDelete
		nw = archive.FlowOpen
	} else {
		old = flow.State
		nw = state
	}
	if flow.State == state {
		return
	}

	if err = s.txUpFlowState(tx, state, uid, flow); err != nil {
		log.Error("txAddOrUpdateFlowState s.txUpdateFlowState error(%v) flow(%+v) state(%d) ", err, flow, state)
		return
	}
	flow.State = state
	nw = state
	return
}
