package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

//启动流程网之前的检查
func (s *Service) startNet(c context.Context, businessID int64, netID int64) (result *net.TriggerResult, err error) {
	var (
		n                  *net.Net
		resultTokenPackage *net.TokenPackage
	)

	//网可用性
	if n, err = s.netByID(c, netID); err != nil {
		log.Error("startNet s.netByID(%d) error(%v)", netID, err)
		return
	}
	if n == nil || !n.IsAvailable() || n.StartFlowID <= 0 {
		log.Error("startNet s.netByID(%d) not found/disabled/start_flow_id=0", netID)
		err = ecode.RequestErr
		return
	}
	if n.BusinessID != businessID {
		log.Error("startNet s.netByID(%d) business(%d) != param business(%d)", netID, n.BusinessID, businessID)
		err = ecode.RequestErr
		return
	}

	//初始节点的token集合
	if resultTokenPackage, err = s.computeNewFlowPackage(c, n.StartFlowID, nil, nil, true); err != nil {
		log.Error("startNet s.computeNewFlowPackage(%d) error(%+v) netid(%d)", n.StartFlowID, err, netID)
		return
	}

	result = &net.TriggerResult{
		NetID:       netID,
		ResultToken: resultTokenPackage,
		NewFlowID:   n.StartFlowID,
		From:        model.LogFromStart,
	}
	return
}

/*FetchJumpFlowInfo ..
* 跳流程详情页, 跳流程是单独db层面的更新
* 参数：flowid(可选)
* 获取网下所有flow和变迁的所有操作项
 */
func (s *Service) FetchJumpFlowInfo(c context.Context, flowID int64) (res *net.JumpInfo, err error) {
	var (
		flow       *net.Flow
		runningNet int64
		flowList   []*net.Flow
	)

	res = &net.JumpInfo{
		Flows:      []*net.SimpleInfo{},
		Operations: []*net.TranOperation{},
	}
	if flow, err = s.flowByID(c, flowID); err != nil {
		log.Error("FetchJumpFlowInfo s.flowByID(%d) error(%v)", flowID, err)
		return
	}
	if flow == nil {
		log.Error("FetchJumpFlowInfo flow(%d) not found", flowID)
		err = ecode.AegisFlowNotFound
		return
	}
	runningNet = flow.NetID
	if flowList, err = s.flowsByNet(c, runningNet); err != nil {
		log.Error("FetchJumpFlowInfo s.flowsByNet(%d) error(%v) flowid(%d)", runningNet, err, flowID)
		return
	}
	flowInfo := []*net.SimpleInfo{}
	for _, item := range flowList {
		flowInfo = append(flowInfo, &net.SimpleInfo{
			ID:     item.ID,
			ChName: item.ChName,
		})
	}

	//网下的所有变迁可操作项
	if res.Operations, err = s.tranOpersByNet(c, []int64{runningNet}, []int8{net.BindTypeTransition}); err != nil {
		log.Error("FetchJumpFlowInfo s.tranOpersByNet(%d) error(%v) flowid(%d)", runningNet, err, flowID)
		return
	}
	res.Flows = flowInfo
	return
}

/**
 * 跳流程提交：rid, oldflowid, newflowid, binds
 * 可以单独提供flowid,单独提供binds
 * 单独提供flowid时，不做流程流转
 */
func (s *Service) jumpFlow(c context.Context, tx *gorm.DB, rid int64, oldFlowID int64, newFlowID int64, binds []int64) (result *net.JumpFlowResult, err error) {
	var (
		frs              []*net.FlowResource
		oldFlow, newFlow *net.Flow
		triggerResult    *net.TriggerResult
		bindList         []*net.TokenBind
		submitPackage    *net.TokenPackage
	)

	log.Info("jumpFlow before rid(%d) oldflowid(%d) newflowid(%d) binds(%v)", rid, oldFlowID, newFlowID, binds)
	if rid <= 0 || oldFlowID <= 0 || (len(binds) == 0 && newFlowID <= 0) {
		err = ecode.RequestErr
		return
	}
	//检查资源是否运行在该节点
	if frs, err = s.gorm.FRByUniques(c, []int64{rid}, []int64{oldFlowID}, true); err != nil {
		log.Error("jumpFlow s.gorm.FRByUniques(%d,%d) error(%v)", rid, oldFlowID, err)
		return
	}
	if len(frs) == 0 {
		log.Error("jumpFlow rid(%d) not running at oldflowid(%d)", rid, oldFlowID)
		err = ecode.AegisNotRunInFlow
		return
	}
	if oldFlow, err = s.flowByID(c, oldFlowID); err != nil {
		log.Error("jumpFlow s.flowByID(%d) error(%v) rid(%d)", oldFlowID, err, rid)
		return
	}
	if oldFlow == nil {
		log.Error("jumpFlow oldflowid(%d) not found, rid(%d)", oldFlowID, rid)
		err = ecode.AegisFlowNotFound
		return
	}

	/**
	 * jumpFlow 指定一批资源，更新为相同现状
	 * 不需要前端指定现状flowid
	 * 资源现状从单线跳到多线程的某个节点时----由运营决定是到sync-split节点还是并发分支上的某个节点------fr数目从1变成1/还是1变成n
	 * 资源现状从多线程跳到某个节点时----节点在并发分支上，如何确定现状flow和新节点都在并发分支上----fr数目从1变成1
	 *                           ----节点在单线上，结束所有运行中都并发现状，加一个单线分支---fr数目从n变成1
	 */
	if len(binds) > 0 {
		if bindList, err = s.tokenBinds(c, binds, true); err != nil {
			log.Error("jumpFlow s.tokenBinds(%v) error(%v) rid(%d)", binds, err, rid)
			return
		}
		if len(bindList) == 0 {
			log.Error("jumpFlow binds(%v) not found, rid(%d)", binds, rid)
			err = ecode.RequestErr
			return
		}
		if submitPackage, err = s.newTokenPackage(c, bindList, true); err != nil {
			log.Error("jumpFlow s.newTokenPackage(%v) error(%v) rid(%d)", binds, err, rid)
			return
		}
	}

	if newFlowID > 0 && newFlowID != oldFlowID {
		if newFlow, err = s.flowByID(c, newFlowID); err != nil {
			log.Error("jumpFlow newflow(%d) error(%v) rid(%d)", newFlowID, err, rid)
			return
		}
		if newFlow == nil {
			log.Error("jumpFlow newflow(%d) not found, rid(%d)", newFlowID, rid)
			err = ecode.AegisFlowNotFound
			return
		}
		if oldFlow.NetID != newFlow.NetID {
			log.Error("jumpFlow rid(%d) run at oldflowid(%d)/net(%d), can't jump to newflowid(%d)/net(%d)", rid, oldFlowID, oldFlow.NetID, newFlowID, newFlow.NetID)
			err = ecode.RequestErr
			return
		}
		triggerResult = &net.TriggerResult{
			RID:         rid,
			NetID:       oldFlow.NetID,
			NewFlowID:   newFlowID,
			OldFlowID:   strconv.FormatInt(oldFlowID, 10),
			SubmitToken: submitPackage,
			ResultToken: submitPackage,
		}

		if err = s.reachNewFlowDB(c, tx, triggerResult); err != nil {
			log.Error("jumpFlow s.reachNewFlowDB error(%v) triggerresult(%+v)", err, triggerResult)
			return
		}
	} else {
		newFlowID = oldFlowID
	}

	result = &net.JumpFlowResult{
		RID:         rid,
		NetID:       oldFlow.NetID,
		SubmitToken: submitPackage,
		ResultToken: submitPackage,
		NewFlowID:   newFlowID,
		OldFlowID:   strconv.FormatInt(oldFlowID, 10),
	}
	return
}

func (s *Service) afterJumpFlow(c context.Context, res *net.JumpFlowResult, bizid int64) (err error) {
	trigger := &net.TriggerResult{
		RID:         res.RID,
		NetID:       res.NetID,
		SubmitToken: res.SubmitToken,
		ResultToken: res.ResultToken,
		NewFlowID:   res.NewFlowID,
		OldFlowID:   res.OldFlowID,
		From:        model.LogFromJump,
	}

	s.afterReachNewFlow(c, trigger, bizid)
	return
}

/**
 * 批量审核提交--单个遍历处理
 * 参数：business, rid, binds
 * flowid(由rid反查现状)
 * transition_id（由flow_id查到的下游transition & token_bind集合 与 token_id的交集决定），非叶子节点若找不到则报错
 * 找到了,则计算新flow & result_token
 * 计算result_token规则：
 *   变迁trigger=人工，必须提供binds;
 *   变迁trigger=其他，binds可以为空 & （新flow要么有绑定token，要么有向线上output不为空）
 */
func (s *Service) computeBatchTriggerResult(c context.Context, businessID int64, rid int64, bindID []int64) (result *net.TriggerResult, err error) {
	var (
		oldFlow *net.Flow
	)

	log.Info("computeBatchTriggerResult start business(%d) rid(%d) bindID(%+v)", businessID, rid, bindID)

	//反查指定资源所运行的节点
	if oldFlow, _, err = s.flowsByRunning(c, rid, businessID, 0, true); err != nil {
		return
	}

	if result, err = s.computeResult(c, rid, oldFlow, bindID, true); err != nil {
		log.Error("computeBatchTriggerResult s.computeResult error(%v) rid(%d) businessid(%d) bindID(%v) oldflow(%+v)", err, rid, businessID, bindID, oldFlow)
		return
	}
	result.From = model.LogFromBatch
	log.Info("computeBatchTriggerResult end business(%d) rid(%d) bindID(%+v) result(%+v)", businessID, rid, bindID, result)
	return
}

/**
 * 详情页审核提交
 * 参数：rid, flowid, binds
 * 当flowid=最后一个，不流转--binds可能为多个; 其余情况，均正常流转到下一个flow--binds只有一个
 * 计算result_token规则：
 *   变迁trigger=人工，必须提供binds;
 *   变迁trigger=其他，binds可以为空 & （新flow要么有绑定token，要么有向线上output不为空）
 */
func (s *Service) computeTriggerResult(c context.Context, rid int64, oldFlowID int64, binds []int64) (result *net.TriggerResult, err error) {
	var (
		oldFlow *net.Flow
		frs     []*net.FlowResource
	)

	log.Info("computeTriggerResult start rid(%d) oldflow(%d) binds(%+v)", rid, oldFlowID, binds)
	if rid <= 0 || oldFlowID <= 0 || len(binds) == 0 {
		err = ecode.RequestErr
		return
	}

	if frs, err = s.gorm.FRByUniques(c, []int64{rid}, []int64{oldFlowID}, true); err != nil {
		log.Error("computeTriggerResult s.gorm.FRByUniques(%d) error(%v) rid(%d) binds(%v)", oldFlowID, err, rid, binds)
		return
	}
	if len(frs) == 0 {
		log.Error("computeTriggerResult rid(%d) not running at oldflow(%d) binds(%v)", rid, oldFlowID, binds)
		err = ecode.AegisNotRunInFlow
		return
	}
	if oldFlow, err = s.flowByID(c, oldFlowID); err != nil {
		log.Error("computeTriggerResult s.flowByID(%d) error(%v) rid(%d) binds(%v)", oldFlowID, err, rid, binds)
		return
	}
	if oldFlow == nil {
		err = ecode.AegisFlowNotFound
		return
	}

	if result, err = s.computeResult(c, rid, oldFlow, binds, false); err != nil {
		log.Error("computeTriggerResult s.computeResult error(%v) binds(%v) flow(%+v) rid(%d)", err, binds, oldFlow, rid)
		return
	}
	result.From = model.LogFromSingle
	log.Info("computeTriggerResult end rid(%d) oldflow(%d) binds(%+v) result(%+v)", rid, oldFlowID, binds, result)
	return
}

func (s *Service) computeResult(c context.Context, rid int64, flow *net.Flow, bindID []int64, fromBatch bool) (result *net.TriggerResult, err error) {
	var (
		hitAudit, hitHelper                    bool //审核流转和非流转操作是互斥的，不能同时提交
		dirs                                   []*net.Direction
		isLast                                 bool
		triggerBind                            = []*net.TokenBind{}
		binds                                  []*net.TokenBind
		trans                                  []*net.Transition
		triggerTranID, newFlowID               int64
		submitTokenPackage, resultTokenPackage *net.TokenPackage
		resultDir                              *net.Direction
	)

	//非叶子节点只能有一个bind
	if rid == 0 || flow == nil {
		log.Error("computeResult rid(%d)=0/flow(%+v)=ni", rid, flow)
		err = ecode.RequestErr
		return
	}

	if binds, err = s.tokenBinds(c, bindID, true); err != nil {
		log.Error("computeBatchTriggerResult s.tokenBinds(%+v) error(%v)", bindID, err)
		return
	}
	tranID := []int64{}
	for _, item := range binds {
		if item.IsBatch() != fromBatch {
			continue
		}

		tranID = append(tranID, item.ElementID)
		if (fromBatch && item.Type == net.BindTypeTranBatch) || (!fromBatch && item.Type == net.BindTypeTransition) {
			hitAudit = true
		}
		if (fromBatch && item.Type == net.BindTypeTranHelpBatch) || (!fromBatch && item.Type == net.BindTypeTranHelp) {
			hitHelper = true
		}
	}
	if len(binds) == 0 || (hitAudit && hitHelper) {
		log.Error("computeBatchTriggerResult binds(%v) not found/both hit audit && helper are not allowed", binds)
		err = ecode.RequestErr
		return
	}

	if hitAudit {
		//同网下、同flow下游的变迁过滤, 非叶子节点必须命中一个变迁
		if dirs, err = s.dirByFlow(c, []int64{flow.ID}, net.DirInput); err != nil {
			log.Error("computeResult s.dirByFlow(%d) error(%v)", flow.ID, err)
			return
		}
		isLast = len(dirs) == 0 //是否为叶子节点
		if !isLast {
			tranID = []int64{}
			for _, item := range dirs {
				tranID = append(tranID, item.TransitionID)
			}
		}
	}

	if trans, err = s.transitions(c, tranID, true); err != nil {
		log.Error("computeResult s.transitions(%v) error(%v) rid(%d) flow(%v)", tranID, err, rid, flow)
		return
	}
	tranMap := map[int64]*net.Transition{}
	for _, item := range trans {
		tranMap[item.ID] = item
	}
	triggerBindID := []int64{}
	triggerTrans := []int64{}
	for _, item := range binds {
		if tranMap[item.ElementID] == nil || tranMap[item.ElementID].Trigger != net.TriggerManual {
			continue
		}

		//叶子节点=同网过滤|非叶子节点=同flow下游变迁过滤
		if (isLast && tranMap[item.ElementID].NetID == flow.NetID) || !isLast {
			triggerBind = append(triggerBind, item)
			triggerTrans = append(triggerTrans, item.ElementID)
			triggerBindID = append(triggerBindID, item.ID)
		}
	}
	if len(triggerBind) == 0 || (!isLast && len(triggerBind) > 1) {
		log.Error("computeResult no triggered(%+v)/non-leaf-flow triggered >2, rid(%d) flow(%+v)", triggerBindID, rid, flow)
		err = ecode.AegisNotTriggerFlow
		return
	}

	//计算提交的变迁 & 提交令牌集合
	if submitTokenPackage, err = s.newTokenPackage(c, triggerBind, true); err != nil {
		log.Error("computeResult s.newTokenPackage error(%v)", err)
		return
	}

	if hitHelper || isLast {
		//不流转的情况
		newFlowID = flow.ID
		resultTokenPackage = submitTokenPackage
	} else {
		//中间节点找到下游新节点
		triggerTranID = triggerTrans[0]
		if resultDir, err = s.fetchTranNextEnableDirs(c, triggerTranID); err != nil {
			log.Error("computeResult s.fetchTranNextEnableDirs(%d) error(%v)", triggerTranID, err)
			return
		}

		newFlowID = resultDir.FlowID
		if resultTokenPackage, err = s.computeNewFlowPackage(c, newFlowID, resultDir, submitTokenPackage, false); err != nil {
			log.Error("computeResult s.computeNewFlowPackage error(%v)", err)
			return
		}
	}

	log.Info("submitpk(%+v)  result(%+v)", submitTokenPackage, resultTokenPackage)
	result = &net.TriggerResult{
		RID:          rid,
		NetID:        flow.NetID,
		SubmitToken:  submitTokenPackage,
		ResultToken:  resultTokenPackage,
		NewFlowID:    newFlowID,
		OldFlowID:    strconv.FormatInt(flow.ID, 10),
		TransitionID: triggerTrans,
	}
	log.Info("computeResult end islast(%v) frombatch(%v) result(%+v)", isLast, fromBatch, result)
	return
}

/**
 * 资源审核详情页，获取运行资源所在网的节点下的所有可允许下游可操作项，去掉批量可操作项:
 * 若为网中的最后一个flow, 获取网中所有变迁的可操作项(bind级别去重)--提交时只改token，不改flow
 * 若为中间节点：只有一个（非并发），直接获取变迁的可操作项; 并发，需指定哪个变迁(先抛错)
 * 前提：同一rid只在同一业务下的一个net下运行
 * 任务列表进入详情页参数：rid,businessid,netid(可选)
 */
func (s *Service) fetchResourceTranInfo(c context.Context, rid int64, businessID int64, netID int64) (result *net.TransitionInfo, err error) {
	var (
		runningFlow *net.Flow
		runningNet  int64
	)
	if rid <= 0 || businessID <= 0 {
		err = ecode.RequestErr
		return
	}

	if runningFlow, runningNet, err = s.flowsByRunning(c, rid, businessID, netID, false); err != nil {
		return
	}
	flowID := runningFlow.ID

	if result, err = s.fetchTaskTranInfo(c, rid, flowID, runningNet); err != nil {
		log.Error("fetchResourceTranInfo s.fetchTaskTranInfo error(%v) rid(%d) businessid(%d) netid(%d) runningnet(%d)", err, rid, businessID, netID, runningNet)
	}
	return
}

/**
* 任务详情页，获取flowid下的所有可允许变迁的可操作项，去掉批量可操作项:
* 若为网中的最后一个flow, 获取网中所有变迁的可操作项(bind级别去重)--提交时只改token，不改flow
* 若为中间节点：只有一个（非并发），直接获取变迁的可操作项;并发，需指定哪个变迁(先抛错)
* 前提：同一rid只在同一业务下的一个net下运行
* 资源列表进入详情页参数：rid, flowid, netid(可选)
 */
func (s *Service) fetchTaskTranInfo(c context.Context, rid int64, flowID int64, netID int64) (result *net.TransitionInfo, err error) {
	var (
		transitionID []int64
		flow         *net.Flow
		enableDir    []*net.Direction
		trans        []*net.Transition
	)

	if rid <= 0 || flowID <= 0 {
		err = ecode.RequestErr
		return
	}

	if enableDir, err = s.fetchFlowNextEnableDirs(c, flowID); err != nil {
		log.Error("fetchTaskTranInfo s.fetchFlowNextEnableDirs(%d) error(%v)", flowID, err)
		return
	}
	//作为分支的叶子节点,获取整个网的所有可操作项
	if len(enableDir) == 0 {
		if netID <= 0 {
			if flow, err = s.flowByID(c, flowID); err != nil {
				log.Error("fetchTaskTranInfo s.flowByID error(%v) rid(%d) flowid(%d)", err, rid, flowID)
				return
			}
			if flow == nil {
				log.Error("fetchTaskTranInfo flow(%d) not found rid(%d) ", flowID, rid)
				err = ecode.AegisFlowNotFound
				return
			}
			netID = flow.NetID
		}

		if transitionID, err = s.tranIDByNet(c, []int64{netID}, true, true); err != nil {
			log.Error("fetchTaskTranInfo s.tranIDByNet(%d) error(%v) rid(%d) flowid(%d)", netID, err, rid, flowID)
			return
		}
	} else {
		transitionID = []int64{}
		for _, item := range enableDir {
			transitionID = append(transitionID, item.TransitionID)
		}
		if trans, err = s.transitions(c, transitionID, false); err != nil {
			log.Error("fetchTaskTranInfo s.transitions(%v) error(%v) rid(%d) flowid(%d)", transitionID, err, rid, flowID)
			return
		}
		transitionID = []int64{}
		for _, item := range trans {
			if item.Trigger == net.TriggerManual {
				transitionID = append(transitionID, item.ID)
			}
		}
	}

	result = &net.TransitionInfo{
		RID:    rid,
		FlowID: flowID,
	}
	if len(transitionID) == 0 {
		return
	}

	if result.Operations, err = s.tranOpers(c, transitionID, []int8{net.BindTypeTransition, net.BindTypeTranHelp}); err != nil {
		log.Error("fetchTaskTranInfo s.tranOpers(%v) error(%v) rid(%d) flowid(%d)", transitionID, err, rid, flowID)
		return
	}
	return
}

/**
 * 全部资源列表页，获取所有变迁的所有批量操作项
 * 参数：businessid,netid(选填)
 */
func (s *Service) fetchBatchOperations(c context.Context, businessID int64, netID int64) (operations []*net.TranOperation, err error) {
	var (
		netIDList []int64
	)

	operations = []*net.TranOperation{}
	if netID > 0 {
		netIDList = []int64{netID}
	} else {
		//查询business_id下的所有net
		if netIDList, err = s.netIDByBusiness(c, businessID); err != nil {
			log.Error("fetchBatchOperations s.netIDByBusiness error(%v) businessid(%d) netid(%d)", err, businessID, netID)
			return
		}
		if len(netIDList) == 0 {
			log.Error("fetchBatchOperations business(%d) no net", businessID)
			return
		}
	}

	//查询所有net下的所有变迁可操作项
	if operations, err = s.tranOpersByNet(c, netIDList, []int8{net.BindTypeTranBatch, net.BindTypeTranHelpBatch}); err != nil {
		log.Error("fetchBatchOperations s.tranOpersByNet(%v) error(%v) businessid(%d) netid(%d)", netIDList, err, businessID, netID)
	}
	return
}

func (s *Service) tranOpersByNet(c context.Context, netIDs []int64, bindTp []int8) (opers []*net.TranOperation, err error) {
	var (
		transitionID = []int64{}
	)

	opers = []*net.TranOperation{}
	if len(netIDs) == 0 {
		return
	}
	if transitionID, err = s.tranIDByNet(c, netIDs, true, true); err != nil {
		log.Error("tranOpersByNet s.tranIDByNet(%v) error(%v) isBatch(%v)", netIDs, err, bindTp)
		return
	}
	if len(transitionID) == 0 {
		return
	}
	if opers, err = s.tranOpers(c, transitionID, bindTp); err != nil {
		log.Error("tranOpersByNet s.tranOpers(%v) error(%v) netid(%v) bindtp(%v)", transitionID, err, netIDs, bindTp)
	}
	return
}

func (s *Service) tranOpers(c context.Context, tranID []int64, bindTp []int8) (opers []*net.TranOperation, err error) {
	var (
		binds = []*net.TokenBind{}
	)

	opers = []*net.TranOperation{}
	if binds, err = s.tokenBindByElement(c, tranID, bindTp); err != nil {
		log.Error("tranOpers s.tokenBindByElement error(%v) tranid(%v) bindtp(%v)", err, tranID, bindTp)
		return
	}
	//bind.token_id维度去重 + bind.id维度排序
	tokenIDBind := map[string][]int64{}
	unique := map[string]*net.TranOperation{}
	for _, item := range binds {
		if _, exist := tokenIDBind[item.TokenID]; !exist {
			tokenIDBind[item.TokenID] = []int64{item.ID}
			unique[item.TokenID] = &net.TranOperation{
				ChName: item.ChName,
			}
			continue
		}
		tokenIDBind[item.TokenID] = append(tokenIDBind[item.TokenID], item.ID)
	}

	for uniqueTokenID, item := range unique {
		item.BindIDList = xstr.JoinInts(tokenIDBind[uniqueTokenID])
		opers = append(opers, item)
	}
	sort.Sort(net.TranOperationArr(opers))
	return
}

/**
* 计算新节点的token集合
* 启动流程节点，只提供newflowID, 找绑定的token
* 流转过程中，根据bindIDs + 触发变迁的可允许有向线(变迁id + newflowid)计算新flow的结果：
* 1. 若有静态绑定token,返回绑定的
* 2. 若没静态绑定，且output="",返回bindID
* 3. 若没静态绑定，且output!="",解析output计算---计算可能涉及到触发变迁/上一个flowid
* 4. 其他情况，为配置错误，应该在配置时避免---todo
* 5. start临时支持无令牌绑定情况
 */
func (s *Service) computeNewFlowPackage(c context.Context, newFlowID int64, fromDir *net.Direction, submitPackage *net.TokenPackage, start bool) (resultTokens *net.TokenPackage, err error) {
	var (
		binds []*net.TokenBind
	)

	//flow绑定了tokens,直接返回
	if binds, err = s.tokenBindByElement(c, []int64{newFlowID}, []int8{net.BindTypeFlow}); err != nil {
		log.Error("computeNewFlowPackage s.tokenBindByElement error(%v) newflowid(%d) fromdir(%v) submit(%+v)", err, newFlowID, fromDir, submitPackage)
		return
	}
	if len(binds) > 0 {
		if resultTokens, err = s.newTokenPackage(c, binds, false); err != nil {
			log.Error("computeNewFlowPackage s.newTokenPackage error(%v) newflowid(%d) fromdir(%v) submit(%+v)", err, newFlowID, fromDir, submitPackage)
		}
		return
	}
	if start {
		return
	}

	//没绑定的，通过有向线的output计算
	if fromDir == nil {
		err = ecode.AegisFlowNoFromDir
		log.Error("computeNewFlowPackage newflowid(%d) has no tokens & fromdir, submit(%+v)", newFlowID, submitPackage)
		return
	}
	//output为空，根据提交内容
	if fromDir.Output == "" {
		resultTokens = submitPackage
		return
	}

	//todo--compute output，解析与prevFlowID + rid + submitTokens相关,return {{"rids":[],"tokens":[]}} --version2

	return
}

/**
 * 计算指定令牌关联的打包形式，包括：各类型的值+关联中文名+关联对应的令牌
 * sametokenid=true，需检查binds对应的tokenid是否一致
 */
func (s *Service) newTokenPackage(c context.Context, binds []*net.TokenBind, sameTokenID bool) (res *net.TokenPackage, err error) {
	var (
		tokenIDList   []int64
		tokens        []*net.Token
		value         interface{}
		sameTokenName string
		hitAudit      bool
	)

	tokenIDStr := ""
	tokenID := ""
	if sameTokenID {
		tokenID = binds[0].TokenID
		sameTokenName = binds[0].ChName
	}
	for _, item := range binds {
		if sameTokenID && tokenID != item.TokenID {
			err = ecode.RequestErr
			log.Error("newTokenPackage binds diff token(%s!=%s) sameTokenID(%v)", tokenID, item.TokenID, sameTokenID)
			return
		}

		if item.Type == net.BindTypeTransition || item.Type == net.BindTypeTranBatch {
			hitAudit = true
		}
		tokenIDStr = tokenIDStr + "," + item.TokenID
	}
	tokenIDStr = strings.TrimLeft(tokenIDStr, ",")

	if tokenIDList, err = xstr.SplitInts(tokenIDStr); err != nil {
		log.Error("newTokenPackage xstr.SplitInts(%s) error(%v) sameTokenID(%v)", tokenIDStr, err, sameTokenID)
		return
	}

	if tokens, err = s.tokens(c, tokenIDList); err != nil {
		log.Error("newTokenPackage s.tokens(%v) error(%v) sameTokenID(%v)", tokenIDStr, err, sameTokenID)
		return
	}
	if len(tokens) == 0 {
		log.Error("newTokenPackage tokens(%s) not found sameTokenID(%v)", tokenIDStr, sameTokenID)
		err = ecode.AegisTokenNotFound
		return
	}

	values := map[string]interface{}{}
	chName := sameTokenName
	for _, item := range tokens {
		if value, err = item.FormatValue(); err != nil {
			log.Error("NewTokenPackage item.FormatValue(%+v) error(%v) sameTokenID(%v)", item, err, sameTokenID)
			return
		}
		values[item.Name] = value
		if sameTokenName == "" {
			chName = chName + item.ChName
		}
	}

	res = &net.TokenPackage{
		Values:      values,
		TokenIDList: tokenIDList,
		ChName:      chName,
		HitAudit:    hitAudit,
	}
	return
}

/**
 * flowsByRunning 指定业务或netid，获取资源的运行节点, 只在一个节点上运行
 * rid反查现状得到flows,netid或businessid做过滤
 * businessID可选，netID可选，2者必须提供一个
 * onlyRunning 只过滤正常运行的资源现状
 */
func (s *Service) flowsByRunning(c context.Context, rid int64, businessID int64, netID int64, onlyRunning bool) (runningFlow *net.Flow, runningNetID int64, err error) {
	var (
		n    *net.Net
		nets []int64
		frs  []*net.FlowResource
	)
	if rid <= 0 || (businessID <= 0 && netID == 0) {
		err = ecode.RequestErr
		log.Error("flowsByRunning rid(%d)/businessid(%d)+netid(%d) are empty, onlyrunning(%v)", rid, businessID, netID, onlyRunning)
		return
	}

	if netID > 0 {
		if n, err = s.netByID(c, netID); err != nil {
			log.Error("flowsByRunning s.netByID(%d) error(%v) rid(%d) businessid(%d) onlyrunning(%v)", netID, err, rid, businessID, onlyRunning)
			return
		}
		if n == nil || n.BusinessID != businessID {
			log.Error("flowsByRunning net(%d) in business(%d) not found, rid(%d) onlyrunning(%v)", netID, businessID, rid, onlyRunning)
			err = ecode.RequestErr
			return
		}
		nets = []int64{netID}
	} else {
		if nets, err = s.netIDByBusiness(c, businessID); err != nil {
			log.Error("flowsByRunning s.netIDByBusiness(%d) error(%v) rid(%d) onlyrunning(%v)", businessID, err, rid, onlyRunning)
			return
		}
	}
	if frs, err = s.gorm.FRByNetRID(c, nets, []int64{rid}, onlyRunning); err != nil {
		log.Error("flowsByRunning s.gorm.FRByNetRID(%d) error(%v) businessid(%d) netid(%d)", rid, err, businessID, netID)
		return
	}
	if len(frs) == 0 {
		log.Error("flowsByRunning rid(%d) not running in business(%d)/netid(%d) onlyrunning(%v)", rid, businessID, netID, onlyRunning)
		err = ecode.AegisNotRunInRange
		return
	}

	flowID := []int64{}
	for _, item := range frs {
		if runningNetID == 0 {
			runningNetID = item.NetID
		} else if item.NetID != runningNetID {
			log.Error("flowsByRunning rid(%d) running is both net(%d) & net(%d), business(%d), onlyrunning(%v)", rid, runningNetID, item.NetID, businessID, onlyRunning)
			err = ecode.AegisRunInDiffNet
			return
		}

		flowID = append(flowID, item.FlowID)
	}
	if len(flowID) > 1 {
		log.Error("flowsByRunning rid(%d) running in flows(%+v)>=2, businessid(%d)/netid(%d) runningnet(%d), onlyrunning(%v)", rid, flowID, businessID, netID, runningNetID, onlyRunning)
		err = ecode.AegisNotRunInFlow
		return
	}

	if runningFlow, err = s.flowByID(c, flowID[0]); err != nil {
		log.Error("flowsByRunning s.flowByID(%d) error(%v), rid(%d) businessid(%d) netid(%d), onlyrunning(%v)", flowID[0], err, rid, businessID, netID, onlyRunning)
		return
	}
	if runningFlow == nil {
		log.Error("flowsByRunning rid(%d) running in flows(%d) not found, businessid(%d) netid(%d) onlyrunnning(%v)", rid, flowID[0], businessID, netID, onlyRunning)
		err = ecode.AegisFlowNotFound
		return
	}
	return
}

/**
 * 流转到新节点且已更新现状表后，所需做的处理
 * 1. 现状流转日志
 * 2. 下游变迁处理，比如：人工变迁是否需要分发
 */
func (s *Service) afterReachNewFlow(c context.Context, pm *net.TriggerResult, bizid int64) (err error) {
	s.sendNetTriggerLog(c, pm)
	err = s.prepareBeforeTrigger(c, []int64{pm.RID}, pm.NewFlowID, bizid)
	return
}

/**
 * 流转到达新节点，更新db操作，需在新节点计算token之后执行
 * 1. 现状表更新
 */
func (s *Service) reachNewFlowDB(c context.Context, tx *gorm.DB, pm *net.TriggerResult) (err error) {
	var (
		flow *net.Flow
	)
	if pm.NetID <= 0 || pm.RID <= 0 || pm.NewFlowID <= 0 {
		log.Error("reachNewFlowDB params error(%+v)", pm)
		err = ecode.RequestErr
		return
	}

	//新节点与netid的同网检查
	if flow, err = s.flowByID(c, pm.NewFlowID); err != nil {
		log.Error("reachNewFlowDB s.flowByID error(%v) params(%+v)", err, pm)
		return
	}
	if flow == nil || flow.NetID != pm.NetID {
		log.Error("reachNewFlowDB s.flowByID not found/not same net, params(%+v), flow(%+v)", pm, flow)
		err = ecode.RequestErr
		return
	}

	if err = s.updateFlowResources(c, tx, pm.NetID, pm.RID, pm.NewFlowID); err != nil {
		log.Error("reachNewFlowDB s.changeFlowResource error(%v) params(%+v)", err, pm)
		return
	}
	//todo--更新现状token情况

	return
}

//取消指定资源的所有运行中流程
func (s *Service) cancelNet(c context.Context, tx *gorm.DB, rids []int64) (res map[int64]string, err error) {
	var (
		frs []*net.FlowResource
	)

	if len(rids) == 0 {
		return
	}
	res = map[int64]string{}
	if frs, err = s.gorm.FRByUniques(c, rids, nil, true); err != nil {
		log.Error("cancelNet s.gorm.FRByUniques error(%v) rids(%+v)", err, rids)
		return
	}
	if len(frs) == 0 {
		return
	}
	if err = s.gorm.CancelFlowResource(c, tx, rids); err != nil {
		log.Error("cancelNet s.gorm.CancelFlowResource error(%+v) rids(%+v)", err, rids)
		return
	}

	//添加日志
	trigger := &net.TriggerResult{
		From: model.LogFromCancle,
	}
	for _, item := range frs {
		pref := res[item.RID]
		if pref != "" {
			pref = "," + pref
		}
		res[item.RID] = fmt.Sprintf("%s%d", pref, item.FlowID)
		trigger.RID = item.RID
		trigger.NetID = item.NetID
		trigger.OldFlowID = strconv.FormatInt(item.FlowID, 10)
		trigger.NewFlowID = item.FlowID
		s.sendNetTriggerLog(c, trigger)
	}
	return
}
