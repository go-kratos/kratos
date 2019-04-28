package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
)

//GetNetList .
func (s *Service) GetNetList(c context.Context, pm *net.ListNetParam) (result *net.ListNetRes, err error) {
	var (
		unames map[int64]string
	)
	if result, err = s.gorm.NetList(c, pm); err != nil {
		return
	}
	if len(result.Result) == 0 {
		return
	}

	//get username
	uid := []int64{}
	for _, item := range result.Result {
		if item.UID > 0 {
			uid = append(uid, item.UID)
		}
	}
	if len(uid) == 0 {
		return
	}

	if unames, err = s.http.GetUnames(c, uid); err != nil || len(unames) == 0 {
		log.Error("GetNetList s.http.GetUnames error(%v) or empty uid(%+v)", err, uid)
		err = nil
	}
	for _, item := range result.Result {
		item.Username = unames[item.UID]
	}

	return
}

//GetNetByBusiness .
func (s *Service) GetNetByBusiness(c context.Context, businessID int64) (res map[int64]string, err error) {
	var (
		list []*net.Net
	)

	res = map[int64]string{}
	if list, err = s.gorm.NetsByBusiness(c, []int64{businessID}, true); err != nil {
		log.Error("GetNetByBusiness s.gorm.NetsByBusiness(%d) error(%v)", businessID, err)
		return
	}

	for _, item := range list {
		res[item.ID] = item.ChName
	}
	return
}

//ShowNet .
func (s *Service) ShowNet(c context.Context, id int64) (r *net.Net, err error) {
	if r, err = s.gorm.NetByID(c, id); err != nil {
		log.Error("ShowNet s.gorm.NetByID(%d) error(%v)", id, err)
	}
	return
}

func (s *Service) netCheckUnique(c context.Context, chName string) (err error, msg string) {
	var exist *net.Net
	if exist, err = s.gorm.NetByUnique(c, chName); err != nil && err != gorm.ErrRecordNotFound {
		log.Error("netCheckUnique s.gorm.NetByID(%s) error(%v)", chName, err)
		return
	}
	if exist != nil {
		err = ecode.AegisUniqueAlreadyExist
		msg = fmt.Sprintf(ecode.AegisUniqueAlreadyExist.Message(), "网", chName)
	}
	return
}

//AddNet .
func (s *Service) AddNet(c context.Context, n *net.Net) (id int64, err error, msg string) {
	if err, msg = s.netCheckUnique(c, n.ChName); err != nil {
		return
	}
	//网初建为禁用状态
	n.DisableTime = time.Now()
	if err = s.gorm.AddItem(c, nil, n); err != nil {
		return
	}
	id = n.ID

	//日志
	oper := &model.NetConfOper{
		OID:    n.ID,
		Action: model.LogNetActionNew,
		UID:    n.UID,
		NetID:  n.ID,
		ChName: n.ChName,
		FlowID: n.StartFlowID,
		Diff:   []string{model.LogFieldTemp(model.LogFieldPID, n.PID, 0, false)},
	}
	s.sendNetConfLog(c, model.LogTypeNetConf, oper)
	return
}

//UpdateNet .
func (s *Service) UpdateNet(c context.Context, uid int64, n *net.NetEditParam) (err error, msg string) {
	var (
		old     *net.Net
		updates = map[string]interface{}{}
		diff    = []string{}
	)
	if old, err = s.gorm.NetByID(c, n.ID); err != nil {
		log.Error("UpdateNet s.gorm.NetByID(%d) error(%v)", n.ID, err)
		return
	}
	if n.ChName != old.ChName {
		if err, msg = s.netCheckUnique(c, n.ChName); err != nil {
			return
		}

		diff = append(diff, model.LogFieldTemp(model.LogFieldChName, n.ChName, old.ChName, true))
		old.ChName = n.ChName
		updates["ch_name"] = n.ChName
	}
	if n.Description != old.Description {
		old.Description = n.Description
		updates["description"] = n.Description
	}
	if err = s.gorm.UpdateFields(c, nil, net.TableNet, old.ID, updates); err != nil {
		return
	}
	s.delNetCache(c, old)

	//日志
	if len(diff) == 0 {
		return
	}
	oper := &model.NetConfOper{
		OID:    old.ID,
		Action: model.LogNetActionUpdate,
		UID:    uid,
		NetID:  old.ID,
		ChName: old.ChName,
		FlowID: old.StartFlowID,
		Diff:   diff,
	}
	s.sendNetConfLog(c, model.LogTypeNetConf, oper)
	return
}

//SwitchNet .
func (s *Service) SwitchNet(c context.Context, id int64, needDisable bool) (err error) {
	var (
		n      *net.Net
		action string
		valid  bool
	)
	if n, err = s.gorm.NetByID(c, id); err != nil {
		log.Error("SwitchNet s.gorm.NetByID(%d) error(%v) needDisable(%v)", id, err, needDisable)
		return
	}
	available := n.IsAvailable()
	if available == !needDisable {
		return
	}

	if needDisable {
		n.DisableTime = time.Now()
		action = model.LogNetActionDisable
	} else {
		if valid, err = s.checkNetValid(c, id); err != nil {
			log.Error("SwitchNet s.checkNetValid error(%v) id(%d)", err, id)
			return
		}
		if !valid {
			log.Error("SwitchNet id(%d) isn't valid, can't be available", id)
			err = ecode.AegisNetErr
			return
		}
		n.DisableTime = net.Recovered
		action = model.LogNetActionAvailable
	}
	if err = s.gorm.UpdateFields(c, nil, net.TableNet, id, map[string]interface{}{"disable_time": n.DisableTime}); err != nil {
		return
	}
	s.delNetCache(c, n)

	//日志
	oper := &model.NetConfOper{
		OID:    n.ID,
		Action: action,
		UID:    n.UID,
		NetID:  n.ID,
		ChName: n.ChName,
		FlowID: n.StartFlowID,
	}
	s.sendNetConfLog(c, model.LogTypeNetConf, oper)
	return
}

//初步检查流程网的可用性：流转完整性
func (s *Service) checkNetValid(c context.Context, netID int64) (valid bool, err error) {
	var (
		n    *net.Net
		flow *net.Flow
		dirs []*net.Direction
	)
	if n, err = s.netByID(c, netID); err != nil || n == nil || n.StartFlowID <= 0 {
		log.Error("checkNetValid s.netByID(%d) error(%v)/not found/start_flow_id=0, net(%+v)", netID, err, n)
		return
	}
	if flow, err = s.flowByID(c, n.StartFlowID); err != nil || flow == nil || !flow.IsAvailable() {
		log.Error("checkNetValid s.flowByID(%d) error(%v)/flow not found/disabled, netid(%d), flow(%+v)", n.StartFlowID, err, netID, flow)
		return
	}
	if dirs, err = s.gorm.DirectionByNet(c, netID); err != nil {
		log.Error("checkNetValid s.gorm.DirectionByNet(%d) error(%v)", netID, err)
		return
	}
	tranPrevMap := map[int64][]int64{}
	tranNextMap := map[int64][]int64{}
	flowPrevMap := map[int64][]int64{}
	for _, item := range dirs {
		if item.Direction == net.DirInput {
			tranPrevMap[item.TransitionID] = append(tranPrevMap[item.TransitionID], item.FlowID)
			continue
		}

		flowPrevMap[item.FlowID] = append(flowPrevMap[item.FlowID], item.TransitionID)
		tranNextMap[item.TransitionID] = append(tranNextMap[item.TransitionID], item.FlowID)
	}

	/**
	flow next empty/trans---dir=input
	flow prev empty(start)/trans---dir=output
	tran next flows---dir=output
	tran prev flows---dir=input
	*/
	for flowID, trans := range flowPrevMap {
		if len(trans) == 0 && flowID != n.StartFlowID {
			log.Error("checkNetValid flow(%d) no previous transition", flowID)
			return
		}
	}
	for tranID, flows := range tranPrevMap {
		prv := len(flows)
		nxt := len(tranNextMap[tranID])
		if prv == 0 || nxt == 0 {
			log.Error("checkNetValid transition(%d) no prev(%d)/next(%d) flow", tranID, prv, nxt)
			return
		}
	}
	for tranID, flows := range tranNextMap {
		prv := len(tranPrevMap[tranID])
		nxt := len(flows)
		if prv == 0 || nxt == 0 {
			log.Error("checkNetValid transition(%d) no prev(%d)/next(%d) flow", tranID, prv, nxt)
			return
		}
	}

	valid = true
	return
}
