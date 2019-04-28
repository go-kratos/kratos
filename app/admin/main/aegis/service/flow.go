package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

/**
 * 业务下所有需要分发任务的flow,用net+flow中文名拼接
 */
func (s *Service) dispatchFlow(c context.Context, businessID []int64, limitFlow []int64) (res map[int64]map[int64]string, err error) {
	var (
		nets   []*net.Net
		tranID []int64
		dirs   []*net.Direction
		flows  []*net.Flow
	)

	//biz:flow:flow_name
	res = map[int64]map[int64]string{}
	//业务下所有可用网
	if nets, err = s.gorm.NetsByBusiness(c, businessID, false); err != nil {
		log.Error("dispatchFlow s.gorm.NetsByBusiness(%v) error(%v)", businessID, err)
		return
	}
	netID := []int64{}
	netMap := map[int64]*net.Net{}
	for _, item := range nets {
		netID = append(netID, item.ID)
		netMap[item.ID] = item
	}
	//网下所有可用变迁
	if tranID, err = s.tranIDByNet(c, netID, true, false); err != nil {
		log.Error("dispatchFlow s.gorm.TransitionIDByNet(%v) error(%v) businessid(%d)", netID, err, businessID)
		return
	}
	if len(tranID) == 0 {
		return
	}
	//变迁所有可用的被指向的有向线
	if dirs, err = s.gorm.DirectionByTransitionID(c, tranID, net.DirInput, false); err != nil {
		log.Error("dispatchFlow s.gorm.DirectionByTransitionID error(%v) businessid(%d)", err, businessID)
		return
	}
	limitFlowMap := map[int64]int{}
	for _, item := range limitFlow {
		limitFlowMap[item] = 1
	}
	accessFlow := []int64{}
	for _, item := range dirs {
		if len(limitFlowMap) > 0 && limitFlowMap[item.FlowID] <= 0 {
			continue
		}
		accessFlow = append(accessFlow, item.FlowID)
	}
	if len(accessFlow) == 0 {
		return
	}
	//拼接每个节点的中文名
	if flows, err = s.flows(c, accessFlow, false); err != nil {
		log.Error("dispatchFlow s.flows error(%v) businessid(%d)", err, businessID)
		return
	}
	sort.Sort(net.FlowArr(flows))
	for _, item := range flows {
		if item == nil || netMap[item.NetID] == nil {
			continue
		}

		nt := netMap[item.NetID]
		if _, exist := res[nt.BusinessID]; !exist {
			res[nt.BusinessID] = map[int64]string{}
		}
		res[nt.BusinessID][item.ID] = nt.ChName + item.ChName
	}
	return
}

//ShowFlow .
func (s *Service) ShowFlow(c context.Context, id int64) (r *net.ShowFlowResult, err error) {
	var (
		f       *net.Flow
		details map[int64][]*net.TokenBind
		n       *net.Net
	)

	if f, err = s.gorm.FlowByID(c, id); err != nil {
		return
	}
	if details, err = s.gorm.TokenBindByElement(c, []int64{id}, []int8{net.BindTypeFlow}, true); err != nil {
		return
	}
	if n, err = s.gorm.NetByID(c, f.NetID); err != nil {
		return
	}
	r = &net.ShowFlowResult{
		Flow:    f,
		Tokens:  details[id],
		IsStart: n.StartFlowID == id,
	}
	return
}

//GetFlowList .
func (s *Service) GetFlowList(c context.Context, pm *net.ListNetElementParam) (result *net.ListFlowRes, err error) {
	var (
		flowID []int64
		tks    map[int64][]*net.TokenBind
		n      *net.Net
		uid    = []int64{}
		unames map[int64]string
	)

	if result, err = s.gorm.FlowList(c, pm); err != nil {
		return
	}
	if len(result.Result) == 0 {
		return
	}
	for _, item := range result.Result {
		flowID = append(flowID, item.ID)
		uid = append(uid, item.UID)
	}
	if tks, err = s.gorm.TokenBindByElement(c, flowID, []int8{net.BindTypeFlow}, true); err != nil {
		return
	}
	if n, err = s.gorm.NetByID(c, pm.NetID); err != nil {
		return
	}
	if unames, err = s.http.GetUnames(c, uid); err != nil {
		log.Error("GetFlowList s.http.GetUnames error(%v)", err)
		err = nil
	}
	for _, item := range result.Result {
		item.IsStart = item.ID == n.StartFlowID
		item.Username = unames[item.UID]
		for _, bd := range tks[item.ID] {
			item.Tokens = append(item.Tokens, bd.ChName)
		}
	}

	return
}

//GetFlowByNet .
func (s *Service) GetFlowByNet(c context.Context, netID int64) (result map[int64]string, err error) {
	var (
		flows []*net.Flow
	)
	result = map[int64]string{}
	if flows, err = s.gorm.FlowsByNet(c, []int64{netID}); err != nil {
		log.Error("GetFlowByNet s.gorm.FlowsByNet(%d) error(%v)", netID, err)
		return
	}

	for _, item := range flows {
		result[item.ID] = item.ChName
	}
	return
}

func (s *Service) checkFlowUnique(c context.Context, netID int64, name string) (err error, msg string) {
	var exist *net.Flow
	if exist, err = s.gorm.FlowByUnique(c, netID, name); err != nil {
		log.Error("checkFlowUnique s.gorm.FlowByUnique(%d,%s) error(%v)", netID, name, err)
		return
	}
	if exist != nil {
		err = ecode.AegisUniqueAlreadyExist
		msg = fmt.Sprintf(ecode.AegisUniqueAlreadyExist.Message(), "节点", name)
	}
	return
}

func (s *Service) checkStartFlowBind(oldFlow *net.Flow, tokenIDList []int64) (err error, msg string) {
	if oldFlow != nil && !oldFlow.IsAvailable() {
		err = ecode.AegisFlowDisabled
		msg = fmt.Sprintf("%s,不能作为初始节点", ecode.AegisFlowDisabled.Message())
		return
	}
	//第一版动态审核初始接入状态：敏感待审、非敏感待审、高频转发待审，非一个确定性值，而系统不提供条件判断和guard解析，由
	//配置初始节点没有令牌，而业务start时自动传入state支持(后续接入统一初始状态进而条件分状态的逻辑后，去掉state字段，且加上该判断)
	//if len(tokenIDList) == 0 {
	//	err = ecode.AegisFlowNoToken
	//	msg = fmt.Sprintf("%s,不能作为初始节点", ecode.AegisFlowNoToken.Message())
	//}

	return
}

//AddFlow .
func (s *Service) AddFlow(c context.Context, uid int64, f *net.FlowEditParam) (id int64, err error, msg string) {
	var (
		tx       *gorm.DB
		diff     = []string{}
		diffBind string
	)
	if err, msg = s.checkFlowUnique(c, f.NetID, f.Name); err != nil {
		return
	}
	if f.IsStart {
		if err, msg = s.checkStartFlowBind(nil, f.TokenIDList); err != nil {
			return
		}
		diff = append(diff, model.LogFieldTemp(model.LogFieldStartFlow, f.IsStart, false, false))
	}
	flow := &net.Flow{
		NetID:       f.NetID,
		Name:        f.Name,
		ChName:      f.ChName,
		Description: f.Description,
		UID:         uid,
	}

	//db update
	tx, err = s.gorm.BeginTx(c)
	if err != nil {
		log.Error("AddFlow s.gorm.BeginTx error(%v)", err)
		return
	}
	if err = s.gorm.AddItem(c, tx, flow); err != nil {
		tx.Rollback()
		return
	}
	if diffBind, _, err, msg = s.compareFlowBind(c, tx, flow.ID, f.TokenIDList, false); err != nil {
		log.Error("AddFlow s.compareFlowBind error(%v) params(%+v)", err, f)
		tx.Rollback()
		return
	}
	if diffBind != "" {
		diff = append(diff, diffBind)
	}
	if f.IsStart {
		if err = s.gorm.NetBindStartFlow(c, tx, flow.NetID, flow.ID); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("AddFlow tx.Commit error(%v)", err)
		return
	}

	id = flow.ID

	//日志
	diff = append(diff, model.LogFieldTemp(model.LogFieldChName, f.ChName, "", false))
	diff = append(diff, model.LogFieldTemp(model.LogFieldName, f.Name, "", false))
	oper := &model.NetConfOper{
		OID:    flow.ID,
		Action: model.LogNetActionNew,
		UID:    flow.UID,
		NetID:  flow.NetID,
		ChName: flow.ChName,
		FlowID: flow.ID,
		Diff:   diff,
	}
	s.sendNetConfLog(c, model.LogTypeFlowConf, oper)
	return
}

// UpdateFlow .
func (s *Service) UpdateFlow(c context.Context, uid int64, f *net.FlowEditParam) (err error, msg string) {
	var (
		old         *net.Flow
		n           *net.Net
		startFlowID int64 = -1
		updates           = map[string]interface{}{}
		tx          *gorm.DB
		diff        = []string{}
		diffBind    string
		changedBind []int64
	)
	if old, err = s.gorm.FlowByID(c, f.ID); err != nil {
		log.Error("UpdateFlow s.gorm.FlowByID(%d) error(%v)", f.ID, err)
		return
	}
	if n, err = s.gorm.NetByID(c, old.NetID); err != nil {
		log.Error("UpdateFlow s.gorm.NetByID(%d) error(%v) flowid(%d)", old.NetID, err, f.ID)
		return
	}
	if f.IsStart && n.StartFlowID != f.ID {
		startFlowID = f.ID
	} else if !f.IsStart && n.StartFlowID == f.ID {
		startFlowID = 0
	}

	if f.IsStart {
		if err, msg = s.checkStartFlowBind(old, f.TokenIDList); err != nil {
			return
		}
		diff = append(diff, model.LogFieldTemp(model.LogFieldStartFlow, true, false, true))
	}

	if f.Name != old.Name {
		if err, msg = s.checkFlowUnique(c, old.NetID, f.Name); err != nil {
			return
		}
		diff = append(diff, model.LogFieldTemp(model.LogFieldName, f.Name, old.Name, true))
		old.Name = f.Name
		updates["name"] = f.Name
	}
	if f.ChName != old.ChName {
		diff = append(diff, model.LogFieldTemp(model.LogFieldChName, f.ChName, old.ChName, true))
		old.ChName = f.ChName
		updates["ch_name"] = f.ChName
	}
	if f.Description != old.Description {
		old.Description = f.Description
		updates["description"] = f.Description
	}

	//db update
	tx, err = s.gorm.BeginTx(c)
	if err != nil {
		log.Error("UpdateFlow s.gorm.BeginTx error(%v)", err)
		return
	}
	if len(updates) > 0 {
		if err = s.gorm.UpdateFields(c, tx, net.TableFlow, old.ID, updates); err != nil {
			tx.Rollback()
			return
		}
	}
	if startFlowID >= 0 {
		if err = s.gorm.NetBindStartFlow(c, tx, n.ID, startFlowID); err != nil {
			tx.Rollback()
			return
		}
	}
	if diffBind, changedBind, err, msg = s.compareFlowBind(c, tx, f.ID, f.TokenIDList, true); err != nil {
		log.Error("updateFlow s.compareFlowBind error(%v) params(%+v)", err, f)
		tx.Rollback()
		return
	}
	if diffBind != "" {
		diff = append(diff, diffBind)
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("UpdateFlow tx.Commit error(%v)", err)
		return
	}
	s.delFlowCache(c, old, changedBind)

	//日志
	if len(diff) == 0 {
		return
	}
	oper := &model.NetConfOper{
		OID:    old.ID,
		Action: model.LogNetActionUpdate,
		UID:    uid,
		NetID:  old.NetID,
		ChName: old.ChName,
		FlowID: old.ID,
		Diff:   diff,
	}
	s.sendNetConfLog(c, model.LogTypeFlowConf, oper)
	return
}

//SwitchFlow .
func (s *Service) SwitchFlow(c context.Context, id int64, needDisable bool) (err error) {
	var (
		old    *net.Flow
		n      *net.Net
		dirs   []*net.Direction
		action string
	)
	if old, err = s.gorm.FlowByID(c, id); err != nil {
		log.Error("SwitchFlow s.gorm.FlowByID(%d) error(%v) needDisable(%v)", id, err, needDisable)
		return
	}
	log.Info("SwitchFlow id(%d) needdisable(%v) old-flow(%+v)", id, needDisable, old)
	available := old.IsAvailable()
	if available == !needDisable {
		return
	}

	if needDisable {
		if dirs, err = s.gorm.DirectionByFlowID(c, []int64{id}, 0); err != nil {
			log.Error("SwitchFlow s.gorm.DirectionByFlowID(%d) error(%v)", id, err)
			return
		}
		if len(dirs) > 0 {
			log.Error("SwitchFlow dir by flow(%d) founded", id)
			err = ecode.AegisFlowBinded
			return
		}
		if n, err = s.gorm.NetByID(c, old.NetID); err != nil {
			log.Error("SwitchFlow s.gorm.NetByID(%d) error(%v) flow(%d)", old.NetID, err, id)
			return
		}
		if n.StartFlowID == id {
			log.Error("SwitchFlow net(%d).startflow=flow(%d) founded", n.ID, id)
			err = ecode.AegisFlowBinded
			return
		}
		old.DisableTime = time.Now()
		action = model.LogNetActionDisable
	} else {
		old.DisableTime = net.Recovered
		action = model.LogNetActionAvailable
	}

	if err = s.gorm.UpdateFields(c, nil, net.TableFlow, id, map[string]interface{}{"disable_time": old.DisableTime}); err != nil {
		return
	}
	s.delFlowCache(c, old, nil)

	//日志
	oper := &model.NetConfOper{
		OID:    old.ID,
		Action: action,
		UID:    old.UID,
		NetID:  old.NetID,
		ChName: old.ChName,
		FlowID: old.ID,
	}
	s.sendNetConfLog(c, model.LogTypeFlowConf, oper)
	return
}
