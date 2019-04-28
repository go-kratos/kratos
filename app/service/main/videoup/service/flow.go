package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_emptyItem = []*archive.FlowData{}
)

// txAddFlowID add private flow id.
func (s *Service) txAddFlowID(c context.Context, tx *sql.Tx, aid, flowID int64, ap *archive.ArcParam) (err error) {
	//backup flow_id in archive_forbid mark as up does
	if _, err = s.arc.TxUpForbid(tx, aid, flowID); err != nil {
		log.Error("s.arc.TxUpForbid(%d,%d) error(%v)", aid, flowID, err)
		return
	}
	//archive_porder
	if _, err = s.arc.TxUpPorder(tx, aid, ap); err != nil {
		log.Error("s.arc.TxUpPorder(%d,%d) error(%v)", aid, flowID, err)
		return
	}
	//flow_design.pool=2
	if _, err = s.txAddFlow(tx, aid, archive.AutoOperUID, flowID, archive.FlowDesignPrivate, "用户稿件添加私单类型"); err != nil {
		log.Error("s.txAddFlow(%d,%d,%d,%d,%s) error(%v)", aid, archive.AutoOperUID, flowID, archive.FlowDesignPrivate, "用户稿件添加私单类型", err)
		return
	}
	log.Info("aid(%d) update autoFlowID(%d)", aid, flowID)
	return
}

//getMidForbid auto add forbid by flow_design  up list mid setting
func (s *Service) getMidForbid(c context.Context, aid int64, ap *archive.ArcParam) (midForbid *archive.ForbidAttr, err error) {
	var (
		hasForbid bool
	)
	log.Info("s.txAddForbidGroup param (%+V) ", ap)
	if hasForbid, midForbid, err = s.getForbidAttrByMid(c, ap.Mid); err != nil || !hasForbid {
		log.Error("s.arc.getForbidAttrByMid(%s) is empty and  error(%v)", midForbid, err)
		err = nil
		return
	}
	midForbid.Convert()
	log.Info("s.txAddForbidGroup param (%+V) and midForbid is (%+v)", ap, midForbid)
	return
}

//hasForbid check mid has forbid or not
func (s *Service) hasForbid(c context.Context, mid int64) (ok bool) {
	_, ok = s.forbidMidsCache[mid]
	return
}

func (s *Service) mergeForbid(c context.Context, a []*archive.ForbidAttr) (sumForbid *archive.ForbidAttr) {
	sumForbid = &archive.ForbidAttr{}
	for _, item := range a {
		//config forbid
		item.Convert()
		sumForbid.Convert()
		ok := int32(1)
		//rank
		if item.Rank.Main == ok {
			sumForbid.Rank.Main = 1
		}
		if item.Rank.RecentArc == ok {
			sumForbid.Rank.RecentArc = 1
		}
		if item.Rank.AllArc == ok {
			sumForbid.Rank.AllArc = 1
		}
		//Recommend
		if item.Recommend.Main == ok {
			sumForbid.Recommend.Main = 1
		}
		//Search
		if item.SearchV == ok {
			sumForbid.SearchV = 1
		}
		//PushBlog
		if item.PushBlogV == ok {
			sumForbid.PushBlogV = 1
		}
		//Dynamic
		if item.Dynamic.Main == ok {
			sumForbid.Dynamic.Main = 1
		}
		//show
		if item.Show.Main == ok {
			sumForbid.Show.Main = 1
		}
		if item.Show.Mobile == ok {
			sumForbid.Show.Mobile = 1
		}
		if item.Show.Web == ok {
			sumForbid.Show.Web = 1
		}
		if item.Show.Oversea == ok {
			sumForbid.Show.Oversea = 1
		}
		if item.Show.Online == ok {
			sumForbid.Show.Online = 1
		}
	}
	sumForbid.Reverse()
	return

}

func (s *Service) transferAttrByScope(c context.Context, a *archive.Archive, ap *archive.ArcParam) (forbid *archive.ForbidAttr, err error) {
	var forbidJSON json.RawMessage
	a.AttrSet(archive.AttrYes, archive.AttrBitIsPorder)
	forbid = &archive.ForbidAttr{}
	autoFlowID, _ := s.getGroupIDByScope(c, ap)
	if autoFlowID <= 0 {
		return
	}
	log.Info("aid (%d) autoFlowID is (%+v) ", a.Aid, autoFlowID)
	for _, flow := range s.flowsCache {
		if flow.ID == autoFlowID {
			forbidJSON = flow.Value
			break
		}
	}
	log.Info("aid (%d) flowForbidJSON  is (%+v) ", a.Aid, forbidJSON)
	if err = json.Unmarshal(forbidJSON, forbid); err != nil {
		log.Error("transferAttrByScope json.Unmarshal(%+v) error(%v)", forbidJSON, err)
		return
	}
	forbid.Convert()
	if forbid.RankV >= 1 {
		a.AttrSet(archive.AttrYes, archive.AttrBitNoRank)
	}
	if forbid.DynamicV >= 1 {
		a.AttrSet(archive.AttrYes, archive.AttrBitNoDynamic)
	}
	if forbid.RecommendV >= 1 {
		a.AttrSet(archive.AttrYes, archive.AttrBitNoRecommend)
	}
	if forbid.SearchV >= 1 {
		a.AttrSet(archive.AttrYes, archive.AttrBitNoSearch)
	}
	if (forbid.ShowV>>3)&1 == 1 {
		a.AttrSet(archive.AttrYes, archive.AttrBitOverseaLock)
	}
	if forbid.PushBlogV == 1 {
		a.AttrSet(archive.AttrYes, archive.AttrBitNoPushBplus)
	}
	log.Info("aid (%d) transferAttrByScope result is (%+v) ", a.Aid, forbid)
	return
}

//transferAttr forbid group affect archive attr
func (s *Service) transferAttr(c context.Context, a *archive.Archive, ap *archive.ArcParam) {
	var (
		hasForbid bool
		forbid    *archive.ForbidAttr
		err       error
	)
	//mid flow forbid  attr
	if hasForbid, forbid, err = s.getForbidAttrByMid(c, ap.Mid); err != nil || !hasForbid {
		log.Error("transferAttr (%+v) s.getForbidAttrByMid(%+v) error(%v)", a, forbid, err)
		return
	}
	if hasForbid {
		log.Info("(%+v) transferAttr s.hasForbid(%+v)", a, forbid)
		if forbid.RankV >= 1 {
			a.AttrSet(archive.AttrYes, archive.AttrBitNoRank)
		}
		if forbid.DynamicV == 1 {
			a.AttrSet(archive.AttrYes, archive.AttrBitNoDynamic)
		}
		if forbid.RecommendV == 1 {
			a.AttrSet(archive.AttrYes, archive.AttrBitNoRecommend)
		}
		if forbid.SearchV == 1 {
			a.AttrSet(archive.AttrYes, archive.AttrBitNoSearch)
		}
		if forbid.PushBlogV == 1 {
			a.AttrSet(archive.AttrYes, archive.AttrBitNoPushBplus)
		}
		if (forbid.ShowV>>3)&1 == 1 {
			a.AttrSet(archive.AttrYes, archive.AttrBitOverseaLock)
		}
		log.Info("transferAttr (%+v) ", a)
	}
}

//getGroupIDByScope flow scope
func (s *Service) getGroupIDByScope(c context.Context, ap *archive.ArcParam) (groupID int64, err error) {
	if groupID, err = s.arc.FindGroupIDByScope(c, archive.FlowDesignPrivate, ap.Porder.IndustryID, ap.Porder.BrandID, ap.Porder.Official); err != nil {
		log.Error("aid(%d) group_id(%d) is not active,err(%+v)", ap.Aid, groupID, err)
		return
	}
	return
}

func (s *Service) getForbidAttrByMid(c context.Context, mid int64) (ok bool, forbid *archive.ForbidAttr, err error) {
	var forbidJSON string
	if forbidJSON, ok = s.forbidMidsCache[mid]; !ok {
		log.Warn("mid(%d) forbidMidsCache(%s) is empty", mid, forbidJSON)
		return
	}
	log.Warn("mid(%d) forbidJson(%+v)", mid, forbidJSON)
	forbid = &archive.ForbidAttr{}
	if err = json.Unmarshal([]byte(forbidJSON), forbid); err != nil {
		log.Error("forbidJson json.Unmarshal(%+v) error(%v)", forbidJSON, err)
		return
	}
	log.Warn("mid(%d) forbid data(%+v)", mid, forbid)
	return
}

// txAddAppFeedID  add appfeed white aids.
func (s *Service) txAddAppFeedID(c context.Context, tx *sql.Tx, ap *archive.ArcParam) (err error) {
	var (
		state   int8
		operUID int64
		remark  string
		//group_id=11天马白名单 方案
		groupID = int64(11)
		//pool=0稿件池子
		pool = archive.FlowDesignAppFeed
	)
	if state, _ = s.arc.CheckActGroupID(c, groupID); state == 1 {
		log.Warn("aid(%d) group_id(%d) is not active", ap.Aid, groupID)
		return
	}
	//私单跟天马互斥
	if ap.Porder.IndustryID > 0 {
		log.Warn("aid(%d) flow_id(%d) is private order_id", ap.Aid, ap.Porder.IndustryID)
		return
	}
	//商单自动进天马白名单
	if ap.OrderID > 0 {
		operUID = archive.CMOperUID
		remark = "商单后台添加天马稿件"
	} else {
		uid, ok := s.isWhiteMid(ap.Mid)
		if !ok {
			return
		}
		operUID = uid
		remark = "白名单用户添加天马稿件"
	}
	if _, err = s.txAddFlow(tx, ap.Aid, operUID, groupID, pool, remark); err != nil {
		log.Error("s.txAddFlow(%d,%d,%d,%d,%s) error(%v)", ap.Aid, operUID, groupID, pool, remark, err)
		return
	}
	log.Info("aid(%d) update appFeed(%d) operUID(%d)", ap.Aid, groupID, operUID)
	return
}

// txAddFlow a public addflow func.
func (s *Service) txAddFlow(tx *sql.Tx, aid, operUID, groupID int64, pool int8, remark string) (id int64, err error) {
	if id, err = s.arc.TxAddFlow(tx, aid, operUID, groupID, pool, remark); err != nil {
		log.Error("s.arc.TxAddFlow(%d,%d,%d,%d,%s) error(%v)", aid, operUID, groupID, pool, remark, err)
		return
	}
	if id <= 0 {
		return
	}
	if _, err = s.arc.TxAddFlowLog(tx, aid, operUID, groupID, pool, remark); err != nil {
		log.Error("s.arc.TxAddFlowLog(%d,%d,%d,%d,%s) error(%v)", aid, operUID, groupID, pool, remark, err)
		return
	}
	return
}

// txUpFlowState .
func (s *Service) txUpFlowState(tx *sql.Tx, id, aid, operUID, groupID int64, pool, state int8, remark string) (err error) {
	var rows int64
	if rows, err = s.arc.TxUpFlowState(tx, state, id); err != nil {
		log.Error("s.arc.txUpFlowState(%d,%d,%d,%d,%s) error(%v)", aid, operUID, groupID, pool, remark, err)
		return
	}
	if rows <= 0 {
		return
	}
	if _, err = s.arc.TxAddFlowLog(tx, aid, operUID, groupID, pool, remark); err != nil {
		log.Error("s.arc.TxAddFlowLog(%d,%d,%d,%d,%s) error(%v)", aid, operUID, groupID, pool, remark, err)
		return
	}
	return
}

//AddByOid  根据mid 配置自动创建对应 flow_design.pool 数据集
func (s *Service) AddByOid(c context.Context, pool int8, oid, uid int64, pm map[string]int32) (err error) {
	//up pool is forbid
	if pool != archive.PoolArcPGC {
		log.Info("flowDesign pool is wrong(%d)", pool)
		err = ecode.RequestErr
		return
	}
	// begin tran
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	_, err = s.txBatchUpFlowsState(c, tx, oid, uid, pm)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	return
}

//AddByMid  根据mid 配置自动创建对应 flow_design.pool 数据集
func (s *Service) AddByMid(c context.Context, pool int8, mid, oid int64, state int8) (err error) {
	var (
		flows  []*archive.FlowData
		flowID int64
	)
	//up pool is forbid
	if pool == archive.PoolUp {
		log.Info("flowDesign pool is wrong(%d)", pool)
		err = ecode.RequestErr
		return
	}
	//check mid hit
	if flows, err = s.arc.CheckFlowMid(c, pool, mid); err != nil {
		log.Error("s.arc.CheckFlowMid(%d,%d) error(%v)", pool, mid, err)
		return
	}
	log.Info("flowDesign flows(%+v)", flows)
	if len(flows) <= 0 {
		log.Info("s.arc.CheckFlowMid(%d,%d) flows(%v) empty", pool, oid, flows)
		err = nil
		return
	}
	// begin tran
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	//create flow_design  flow data by mid
	for _, midFlow := range flows {
		//skip archive and porder flow
		// 因为稿件及私单业务mid 自动添加 已经具备 这里不要重复添加
		if midFlow.Pool == archive.PoolArc || midFlow.Pool == archive.PoolPorder {
			continue
		}
		//check flow_design already exist
		if flowID, err = s.arc.CheckFlowGroupID(c, pool, oid, midFlow.GroupID); err != nil {
			log.Error("s.arc.CheckFlowGroupID(%d,%d,%d) error(%v)", midFlow.Pool, oid, midFlow.GroupID, err)
			continue
		}
		log.Info("flowDesign flowID (%v) mid(%d)", flowID, mid)
		if flowID > 0 {
			//支持下线流量配置
			if state == archive.FlowDelete {
				if err = s.txUpFlowState(tx, flowID, oid, midFlow.UID, midFlow.GroupID, pool, state, "下线流量配置"); err != nil {
					log.Error("s.txUpFlowState(%d,%d,%d,%d,%d,%d,%s) error(%v)", flowID, oid, midFlow.UID, midFlow.GroupID, midFlow.Pool, state, midFlow.Remark, err)
					tx.Rollback()
					return
				}
			}
			continue
		}
		//check insert state
		if state == archive.FlowDelete {
			continue
		}
		//insert only new
		if _, err = s.txAddFlow(tx, oid, midFlow.UID, midFlow.GroupID, pool, fmt.Sprintf("[%s]-%s", "auto", midFlow.Remark)); err != nil {
			log.Error("s.txAddFlow(%d,%d,%d,%d,%s) error(%v)", oid, midFlow.UID, midFlow.GroupID, midFlow.Pool, midFlow.Remark, err)
			tx.Rollback()
			return
		}
		log.Info("flowDesign auto add flow(%+v) success", midFlow)
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
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
		if state < 0 {
			continue
		}
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
	//groupId 匹配业务
	if pools, err = s.arc.FlowGroupPools(c, groups); err != nil {
		log.Error("txBatchUpFlowsState s.arc.FlowGroupPools(%v) error(%v) params(%+v)", groups, err, pm)
		return
	}
	//写入业务流量
	for groupID, pool := range pools {
		if _, diff, err = s.txAddOrUpdateFlowState(c, tx, aid, groupID, uid, pool, groupStates[groupID], "流量接口修改"); err != nil {
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

/**
 * txAddOrUpdateFlowState 新增或更新流量套餐的状态
 * return *archive.FlowData/nil, bool, error/nil
 */
func (s *Service) txAddOrUpdateFlowState(c context.Context, tx *sql.Tx, oid, groupID, uid int64, pool, state int8, remark string) (flow *archive.FlowData, diff string, err error) {
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
		if flow.ID, err = s.txAddFlow(tx, flow.OID, uid, flow.GroupID, flow.Pool, remark); err != nil {
			log.Error("txAddOrUpdateFlowState s.txAddFlow error(%v) flow(%+v) state(%d)", err, flow, state)
			return
		}

	}
	if flow.State == state {
		return
	}

	if err = s.txUpFlowStateByFlow(tx, state, uid, flow); err != nil {
		log.Error("txAddOrUpdateFlowState s.txUpdateFlowState error(%v) flow(%+v) state(%d) ", err, flow, state)
		return
	}
	flow.State = state
	return
}

/**
 * txUpFlowState 更新流量套餐的状态
 * return error/nil
 */
func (s *Service) txUpFlowStateByFlow(tx *sql.Tx, state int8, uid int64, f *archive.FlowData) (err error) {
	if f == nil {
		return
	}
	var rows int64
	if rows, err = s.arc.TxUpFlowState(tx, state, f.ID); err != nil {
		log.Error("updateFlowState s.arc.TxUpFlowState error(%v) id(%d) state(%d)", err, f.ID, state)
		return
	}
	if rows <= 0 {
		return
	}

	//action := archive.FlowLogUpdate
	//if state == archive.FlowDelete {
	//	action = archive.FlowLogDel
	//}
	if _, err = s.arc.TxAddFlowLog(tx, f.OID, uid, f.GroupID, f.Pool, "审核后台修改状态"); err != nil {
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

//FlowPage .
func (s *Service) FlowPage(c context.Context, pool, gid, pn, ps int64) (data *archive.FlowPagerData, err error) {
	var (
		count int64
		flows []*archive.FlowData
	)
	data = &archive.FlowPagerData{
		Items: _emptyItem,
	}
	if count, err = s.arc.CountByGID(c, pool, gid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if count == 0 {
		return
	}

	if flows, err = s.arc.FlowPage(c, pool, gid, pn, ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(flows) == 0 {
		return
	}
	data = &archive.FlowPagerData{
		Items: flows,
		Pager: &archive.Pager{
			Num:   pn,
			Size:  ps,
			Total: count,
		},
	}
	return

}

//TagrgetFlows .
func (s *Service) TagrgetFlows(c context.Context, pool, gid int64, oids []int64) (res []int64, err error) {
	var (
		flows []*archive.FlowData
	)
	if flows, err = s.arc.OidsFlowByGID(c, pool, gid, xstr.JoinInts(oids)); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(flows) == 0 {
		res = []int64{}
		return
	}
	for _, v := range flows {
		res = append(res, v.OID)
	}
	return

}
