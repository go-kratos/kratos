package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

func (s *Service) prepareBeforeTrigger(c context.Context, rids []int64, flowID, bizid int64) (err error) {
	var (
		enableDir  []*net.Direction
		tids       []int64
		transition []*net.Transition
	)
	//以flow为起点的有向线
	if enableDir, err = s.fetchFlowNextEnableDirs(c, flowID); err != nil {
		log.Error("prepareBeforeTrigger s.fetchFlowNextEnableDirs error(%v) flowid(%d) rids(%v)", err, flowID, rids)
		return
	}
	//flow下游没有变迁，属于正常情况
	if len(enableDir) == 0 {
		return
	}

	//所有可用有向线的变迁限制
	tids = []int64{}
	for _, dir := range enableDir {
		tids = append(tids, dir.TransitionID)
	}
	if transition, err = s.transitions(c, tids, true); err != nil {
		log.Error("prepareBeforeTrigger s.transitions(%+v) error(%v) flowid(%d) rids(%v)", tids, err, flowID, rids)
		return
	}
	for _, item := range transition {
		if item.Trigger != net.TriggerManual || item.Limit <= 0 {
			continue
		}

		for _, rid := range rids {
			s.sendCreateTaskMsg(c, rid, flowID, item.Limit, bizid)
		}
	}
	return
}

//ShowTransition .
func (s *Service) ShowTransition(c context.Context, id int64) (r *net.ShowTransitionResult, err error) {
	var (
		t   *net.Transition
		tks map[int64][]*net.TokenBind
	)

	if t, err = s.gorm.TransitionByID(c, id); err != nil {
		return
	}
	if tks, err = s.gorm.TokenBindByElement(c, []int64{id}, net.BindTranType, true); err != nil {
		return
	}
	r = &net.ShowTransitionResult{
		Transition: t,
		Tokens:     tks[id],
	}
	return
}

//GetTranByNet .
func (s *Service) GetTranByNet(c context.Context, netID int64) (result map[int64]string, err error) {
	var (
		trans []*net.Transition
	)

	result = map[int64]string{}
	if trans, err = s.gorm.TranByNet(c, netID, true); err != nil {
		log.Error("GetTranByNet s.gorm.TranByNet(%d) error(%v)", netID, err)
		return
	}

	for _, item := range trans {
		result[item.ID] = item.ChName
	}
	return
}

//GetTransitionList .
func (s *Service) GetTransitionList(c context.Context, pm *net.ListNetElementParam) (result *net.ListTransitionRes, err error) {
	var (
		transitionID []int64
		tks          map[int64][]*net.TokenBind
		uid          = []int64{}
		unames       map[int64]string
	)

	if result, err = s.gorm.TransitionList(c, pm); err != nil {
		return
	}
	if len(result.Result) == 0 {
		return
	}
	for _, item := range result.Result {
		transitionID = append(transitionID, item.ID)
		uid = append(uid, item.UID)
	}
	if tks, err = s.gorm.TokenBindByElement(c, transitionID, net.BindTranType, true); err != nil {
		return
	}
	if unames, err = s.http.GetUnames(c, uid); err != nil {
		log.Error("GetTransitionList s.http.GetUnames error(%v)", err)
		err = nil
	}
	for _, item := range result.Result {
		item.Username = unames[item.UID]
		for _, bd := range tks[item.ID] {
			item.Tokens = append(item.Tokens, bd.ChName)
		}
	}
	return
}

func (s *Service) checkTransitionUnique(c context.Context, netID int64, name string) (err error, msg string) {
	var exist *net.Transition
	if exist, err = s.gorm.TransitionByUnique(c, netID, name); err != nil {
		log.Error("checkTransitionUnique s.gorm.TransitionByUnique(%d,%s) error(%v)", netID, name, err)
		return
	}
	if exist != nil {
		err = ecode.AegisUniqueAlreadyExist
		msg = fmt.Sprintf(ecode.AegisUniqueAlreadyExist.Message(), "变化", name)
	}
	return
}

//AddTransition .
func (s *Service) AddTransition(c context.Context, uid int64, f *net.TransitionEditParam) (id int64, err error, msg string) {
	var (
		tx       *gorm.DB
		diff     = []string{}
		diffBind string
	)
	if err, msg = s.checkTransitionUnique(c, f.NetID, f.Name); err != nil {
		return
	}

	tran := &net.Transition{
		NetID:       f.NetID,
		Trigger:     f.Trigger,
		Limit:       f.Limit,
		Name:        f.Name,
		ChName:      f.ChName,
		Description: f.Description,
		UID:         uid,
	}
	//db update
	tx, err = s.gorm.BeginTx(c)
	if err != nil {
		log.Error("AddTransition s.gorm.BeginTx error(%v)", err)
		return
	}
	if err = s.gorm.AddItem(c, tx, tran); err != nil {
		tx.Rollback()
		return
	}
	if diffBind, _, err, msg = s.compareTranBind(c, tx, tran.ID, f.TokenList, false); err != nil {
		log.Error("AddTransition s.compareTranBind error(%v) params(%+v)", err, f)
		tx.Rollback()
		return
	}
	if diffBind != "" {
		diff = append(diff, diffBind)
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("AddTransition tx.Commit error(%v)", err)
		return
	}
	id = tran.ID

	//日志
	diff = append(diff, model.LogFieldTemp(model.LogFieldChName, f.ChName, "", false))
	diff = append(diff, model.LogFieldTemp(model.LogFieldName, f.Name, "", false))
	diff = append(diff, model.LogFieldTemp(model.LogFieldLimit, f.Limit, "", false))
	diff = append(diff, model.LogFieldTemp(model.LogFieldTrigger, net.TriggerDesc[f.Trigger], "", false))
	oper := &model.NetConfOper{
		OID:    tran.ID,
		Action: model.LogNetActionNew,
		UID:    tran.UID,
		NetID:  tran.NetID,
		ChName: tran.ChName,
		TranID: tran.ID,
		Diff:   diff,
	}
	s.sendNetConfLog(c, model.LogTypeTranConf, oper)
	return
}

//UpdateTransition .
func (s *Service) UpdateTransition(c context.Context, uid int64, f *net.TransitionEditParam) (err error, msg string) {
	var (
		old         *net.Transition
		updates     = map[string]interface{}{}
		tx          *gorm.DB
		diff        = []string{}
		diffBind    string
		changedBind []int64
	)
	if old, err = s.gorm.TransitionByID(c, f.ID); err != nil {
		log.Error("UpdateTransition s.gorm.TransitionByID(%d) error(%v)", f.ID, err)
		return
	}
	if f.Name != old.Name {
		if err, msg = s.checkTransitionUnique(c, f.NetID, f.Name); err != nil {
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
	if f.Limit != old.Limit {
		diff = append(diff, model.LogFieldTemp(model.LogFieldLimit, f.Limit, old.Limit, true))
		old.Limit = f.Limit
		updates["limit"] = f.Limit
	}
	if f.Trigger != old.Trigger {
		diff = append(diff, model.LogFieldTemp(model.LogFieldTrigger, net.TriggerDesc[f.Trigger], net.TriggerDesc[old.Trigger], true))
		old.Trigger = f.Trigger
		updates["trigger"] = f.Trigger
	}

	//db update
	tx, err = s.gorm.BeginTx(c)
	if err != nil {
		log.Error("UpdateTransition s.gorm.BeginTx error(%v)", err)
		return
	}
	if len(updates) > 0 {
		if err = s.gorm.UpdateFields(c, tx, net.TableTransition, old.ID, updates); err != nil {
			tx.Rollback()
			return
		}
	}
	if diffBind, changedBind, err, msg = s.compareTranBind(c, tx, old.ID, f.TokenList, true); err != nil {
		log.Error("UpdateTransition s.compareTranBind error(%v) params(%+v)", err, f)
		tx.Rollback()
		return
	}
	if diffBind != "" {
		diff = append(diff, diffBind)
	}

	if err = tx.Commit().Error; err != nil {
		log.Error("UpdateTransition tx.Commit error(%v)", err)
		return
	}
	s.delTranCache(c, old, changedBind)

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
		TranID: old.ID,
		Diff:   diff,
	}
	s.sendNetConfLog(c, model.LogTypeTranConf, oper)
	return
}

//SwitchTransition .
func (s *Service) SwitchTransition(c context.Context, id int64, needDisable bool) (err error) {
	var (
		old    *net.Transition
		dirs   []*net.Direction
		action string
	)
	if old, err = s.gorm.TransitionByID(c, id); err != nil {
		log.Error("SwitchTransition s.gorm.TransitionByID(%d) error(%v) needDisable(%v)", id, err, needDisable)
		return
	}
	available := old.IsAvailable()
	if available == !needDisable {
		return
	}

	if needDisable {
		if dirs, err = s.gorm.DirectionByTransitionID(c, []int64{id}, 0, true); err != nil {
			log.Error("SwitchTransition s.gorm.DirectionByTransitionID(%d) error(%v)", id, err)
			return
		}
		if len(dirs) > 0 {
			log.Error("SwitchTransition dir by transition(%d) founded", id)
			err = ecode.AegisTranBinded
			return
		}
		old.DisableTime = time.Now()
		action = model.LogNetActionDisable
	} else {
		old.DisableTime = net.Recovered
		action = model.LogNetActionAvailable
	}

	if err = s.gorm.UpdateFields(c, nil, net.TableTransition, id, map[string]interface{}{"disable_time": old.DisableTime}); err != nil {
		return
	}
	s.delTranCache(c, old, nil)

	//日志
	oper := &model.NetConfOper{
		OID:    old.ID,
		Action: action,
		UID:    old.UID,
		NetID:  old.NetID,
		ChName: old.ChName,
		TranID: old.ID,
	}
	s.sendNetConfLog(c, model.LogTypeTranConf, oper)
	return
}
