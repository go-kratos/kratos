package service

import (
	"context"
	"encoding/json"
	"strings"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// txUpFlowID update flow_id by videoParam.
func (s *Service) txUpFlowID(tx *sql.Tx, ap *archive.ArcParam) (err error) {
	if _, err = s.arc.TxUpFlowID(tx, ap.Aid, ap.OnFlowID); err != nil {
		log.Error("archive_forbid.on_flow_id s.TxUpFlowID(%d,%d) error(%v)", ap.Aid, ap.OnFlowID, err) // NOTE:  update  es index  question ,  if after update archive table
		return
	}
	var flowID int64
	if flowID, err = s.arc.FlowByPool(archive.PoolPrivateOrder, ap.Aid); err != nil {
		log.Error("flow_design s.arc.FlowByPool(%d,%d,%d) error(%v)", ap.Aid, ap.UID, ap.OnFlowID, err)
		return
	}
	if flowID > 0 {
		if _, err = s.arc.TxUpFlow(tx, flowID, ap.OnFlowID, ap.UID); err != nil {
			log.Error("flow_design s.arc.TxUpFlow(%d,%d,%d) error(%v)", ap.Aid, ap.UID, ap.OnFlowID, err)
			return
		}
		if _, err = s.arc.TxAddFlowLog(tx, archive.PoolPrivateOrder, archive.FlowLogUpdate, ap.Aid, ap.UID, ap.OnFlowID, "审核后台修改稿件私单类型"); err != nil {
			log.Error("s.arc.TxAddFlowLog(%d,%d,%d) error(%v)", ap.Aid, ap.UID, ap.OnFlowID, err)
			return
		}
	} else {
		if _, err = s.txAddFlow(tx, archive.PoolPrivateOrder, ap.Aid, ap.OnFlowID, ap.UID, "审核后台添加稿件私单类型"); err != nil {
			log.Error("flow_design s.arc.TxAddFlow(%d,%d,%d) error(%v)", ap.Aid, ap.UID, ap.OnFlowID, err)
			return
		}
	}

	log.Info("aid(%d)  flowid(%d)", ap.Aid, ap.OnFlowID)
	return
}

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
 * getFlowsByOID 命中哪些流量套餐
 * return []*archive.FlowData, error/nil
 */
func (s *Service) getFlowsByOID(c context.Context, oid int64) (flows []*archive.FlowData, err error) {
	if flows, err = s.arc.FlowsByOID(c, oid); err != nil {
		log.Error("getFlowsByOID s.arc.FlowsByOID error(%v) oid(%d)", err, oid)
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

//HitFlowGroups 命中哪些指定的流量套餐
func (s *Service) HitFlowGroups(c context.Context, oid int64, includePools []int8) (res map[string]int, err error) {
	var (
		flows []*archive.FlowData
	)
	res = map[string]int{}
	includes := map[int8]int8{}
	if flows, err = s.getFlowsByOID(c, oid); err != nil {
		return
	}
	for _, p := range includePools {
		includes[p] = 1
	}

	for _, f := range flows {
		if includes[f.Pool] != 1 {
			continue
		}

		//merge their values
		value := map[string]int{}
		if err = json.Unmarshal(f.GroupValue, &value); err != nil {
			log.Error("HitFlowGroups json.Unmarshal error(%v) value(%s) oid(%d) flow.id(%d)", err, string(f.GroupValue), oid, f.ID)
			return
		}
		for attr, val := range value {
			if val == 1 {
				res[attr] = val
			}
		}
	}
	log.Info("HitFlowGroups oid(%d) includepools(%v) res(%+v)", oid, includePools, res)
	return
}

func (s *Service) txBatchUpFlowsState(c context.Context, tx *sql.Tx, aid, uid int64, pm map[string]int32) (conts []string, err error) {
	var (
		diff   string
		groups []int64
		pools  map[int64]int8
	)

	groupStates := map[int64]int8{}
	for attr, state := range pm {
		groupID := archive.FlowAttrMap[attr]
		if groupID <= 0 {
			continue
		}
		groups = append(groups, groupID)
		if state == 0 {
			groupStates[groupID] = archive.FlowDelete
		} else {
			groupStates[groupID] = archive.FlowOpen
		}
	}
	if len(groups) <= 0 {
		return
	}

	if pools, err = s.arc.FlowGroupPools(c, groups); err != nil {
		log.Error("txBatchUpFlowsState s.arc.FlowGroupPools(%v) error(%v) params(%+v)", groups, err, pm)
		return
	}

	for groupID, pool := range pools {
		if _, diff, err = s.txAddOrUpdateFlowState(c, tx, aid, groupID, uid, pool, groupStates[groupID], "审核后台修改"); err != nil {
			log.Error("txBatchUpFlowsState s.txAddOrUpdateFlowState(%d,%d,%d,%d,%d) error(%v) params(%+v)", aid, groupID, uid, pool, groupStates[groupID], err, pm)
			return
		}
		if diff != "" {
			conts = append(conts, diff)
		}
	}
	log.Info("txBatchUpFlowsState aid(%d) params(%+v) conts(%v)", aid, pm, conts)
	return
}
