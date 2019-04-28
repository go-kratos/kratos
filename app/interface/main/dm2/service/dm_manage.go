package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/model/oplog"
	account "go-common/app/service/main/account/api"
	assmdl "go-common/app/service/main/assist/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_assistUpperLimit = 100
)

// isAssist check if the user is assist of upper or not.
func (s *Service) isAssist(c context.Context, mid, uid int64) (err error) {
	arg := assmdl.ArgAssist{
		Mid:       mid,
		AssistMid: uid,
		Type:      assmdl.TypeDm,
		RealIP:    "",
	}
	res, err := s.assRPC.Assist(c, &arg)
	if err != nil {
		log.Error("s.assRPC.Assist(%v) error(%v)", arg, err)
		return ecode.AccessDenied
	}
	if res.Assist == 1 && res.Allow == 1 {
		return nil
	}
	if res.Assist == 1 && res.Count > _assistUpperLimit {
		return ecode.DMAssistOpToMuch
	}
	return ecode.AccessDenied
}

// isUpper check if the user is upper.
func (s *Service) isUpper(mid, uid int64) bool {
	return mid == uid
}

// EditDMState change dm state
// 0：正常、1：删除、10：用户删除、11：举报脚本删除
func (s *Service) EditDMState(c context.Context, tp int32, mid, oid int64, state int32, dmids []int64, source oplog.Source, operatorType oplog.OperatorType) (err error) {
	var (
		affect, action int64
	)
	if source <= 0 {
		source = oplog.SourceUp
	}
	if operatorType <= 0 {
		operatorType = oplog.OperatorUp
	}
	sub, err := s.subject(c, tp, oid)
	if err != nil {
		return
	}
	switch state {
	case model.StateNormal, model.StateDelete:
		var isAssist bool
		if !s.isUpper(sub.Mid, mid) {
			if err = s.isAssist(c, sub.Mid, mid); err != nil {
				return
			}
		}
		affect, err = s.dao.UpdateDMStat(c, tp, oid, state, dmids)
		if err != nil {
			log.Error("s.dao.UpdateDMStat(oid:%d state:%d dmids:%v) error(%v)", oid, state, dmids, err)
			return
		}
		if affect > 0 {
			if state == model.StateDelete && sub.IsMonitoring() {
				s.oidLock.Lock()
				s.moniOidMap[sub.Oid] = struct{}{}
				s.oidLock.Unlock()
			}
			s.OpLog(c, oid, mid, time.Now().Unix(), int(tp), dmids, "status", "", strconv.FormatInt(int64(state), 10), "更新弹幕状态", source, operatorType)
			if state == model.StateDelete {
				action = assmdl.ActDelete
				affect = -affect
			}
			if sub.Count+affect < 0 {
				affect = -sub.Count
			}
			if _, err = s.dao.IncrSubjectCount(c, tp, oid, affect); err != nil {
				return
			}
			if isAssist {
				for _, dmid := range dmids {
					s.addAssistLog(sub.Mid, mid, oid, action, dmid)
				}
			}
		}
	case model.StateUserDelete:
		affect, err = s.dao.UpdateUserDMStat(c, tp, oid, mid, state, dmids)
		if err != nil {
			log.Error("s.dao.UpdateUserDMStat(mid:%d oid:%d state:%d dmids:%v) error(%v)", mid, oid, state, dmids, err)
			return
		}
		if affect > 0 {
			if sub.IsMonitoring() {
				s.oidLock.Lock()
				s.moniOidMap[sub.Oid] = struct{}{}
				s.oidLock.Unlock()
			}
			s.OpLog(c, oid, mid, time.Now().Unix(), int(tp), dmids, "status", "", fmt.Sprint(state), "更新弹幕状态", source, operatorType)
			affect = -affect
			if sub.Count+affect < 0 {
				affect = -sub.Count
			}
			if _, err = s.dao.IncrSubjectCount(c, tp, oid, affect); err != nil {
				return
			}
		}
	case model.StateScriptDelete:
		affect, err = s.dao.UpdateDMStat(c, tp, oid, state, dmids)
		if err != nil {
			log.Error("s.dao.UpdateDMStat(oid:%d state:%d dmids:%v) error(%v)", oid, state, dmids, err)
			return
		}
		if affect > 0 {
			if sub.IsMonitoring() {
				s.oidLock.Lock()
				s.moniOidMap[sub.Oid] = struct{}{}
				s.oidLock.Unlock()
			}
			s.OpLog(c, oid, mid, time.Now().Unix(), int(tp), dmids, "status", "", fmt.Sprint(state), "更新弹幕状态", source, operatorType)
			affect = -affect
			if sub.Count+affect < 0 {
				affect = -sub.Count
			}
			if _, err = s.dao.IncrSubjectCount(c, tp, oid, affect); err != nil {
				return
			}
		}
	default:
		err = ecode.RequestErr
	}
	return
}

// EditDMPool edit dm pool.
func (s *Service) EditDMPool(c context.Context, tp int32, mid, oid int64, pool int32, ids []int64, source oplog.Source, operatorType oplog.OperatorType) (err error) {
	var (
		isAssist bool
		affect   int64
	)
	if pool != model.PoolNormal && pool != model.PoolSubtitle {
		err = ecode.RequestErr
		return
	}
	// pool 2 dm can't move to other pool
	dmids := make([]int64, 0, len(ids))
	indexs, _, err := s.dao.IndexsByid(c, tp, oid, ids)
	for dmid, index := range indexs {
		if index.Pool != model.PoolSpecial {
			dmids = append(dmids, dmid)
		}
	}
	if len(dmids) <= 0 {
		return
	}
	if source <= 0 {
		source = oplog.SourceUp
	}
	if operatorType <= 0 {
		operatorType = oplog.OperatorUp
	}
	sub, err := s.subject(c, tp, oid)
	if err != nil {
		return
	}
	// maximum batch move count to subtitle pool is 300 when the rank of
	// user is equal or less than 15000
	if pool == model.PoolSubtitle {
		var (
			reply *account.ProfileReply
		)
		if reply, err = s.accountRPC.Profile3(c, &account.MidReq{Mid: mid}); err != nil {
			log.Error("accRPC.Profile3(%v) error(%v)", mid, err)
			return
		}
		if reply.Profile.Rank <= 15000 && int(sub.MoveCnt)+len(dmids) > 300 {
			err = ecode.DMPoolLimit
			return
		}
	}
	if !s.isUpper(sub.Mid, mid) {
		if err = s.isAssist(c, sub.Mid, mid); err != nil {
			return
		}
	}
	if sub.Childpool < pool {
		if _, err = s.dao.UpSubjectPool(c, tp, oid, pool); err != nil {
			return
		}
	}
	if affect, err = s.dao.UpdateDMPool(c, tp, oid, pool, dmids); err != nil {
		log.Error("s.dao.UpdateDMPool(oid:%d pool:%d dmids:%v) error(%v)", oid, pool, dmids, err)
		return
	}
	if affect > 0 {
		if pool == model.PoolNormal {
			s.dao.IncrSubMoveCount(c, sub.Type, sub.Oid, -affect) // NOTE update move_count,ignore error
		} else {
			s.dao.IncrSubMoveCount(c, sub.Type, sub.Oid, affect) // NOTE update move_count,ignore error
		}
		s.OpLog(c, oid, mid, time.Now().Unix(), int(tp), dmids, "pool", "", strconv.FormatInt(int64(pool), 10), "弹幕池变更", source, operatorType)
	}
	if isAssist {
		for _, dmid := range dmids {
			s.addAssistLog(sub.Mid, mid, oid, assmdl.ActProtect, dmid)
		}
	}
	return
}

// EditDMAttr update dm attribute.
func (s *Service) EditDMAttr(c context.Context, tp int32, mid, oid int64, bit uint, value int32, dmids []int64, source oplog.Source, operatorType oplog.OperatorType) (affectIds []int64, err error) {
	var isAssist bool
	affectIds = make([]int64, 0, len(dmids))
	if value != model.AttrNo && value != model.AttrYes {
		err = ecode.RequestErr
		return
	}
	if source <= 0 {
		source = oplog.SourceUp
	}
	if operatorType <= 0 {
		operatorType = oplog.OperatorUp
	}
	sub, err := s.subject(c, tp, oid)
	if err != nil {
		return
	}
	if !s.isUpper(sub.Mid, mid) {
		if err = s.isAssist(c, sub.Mid, mid); err != nil {
			return
		}
	}
	idxMap, _, err := s.dao.IndexsByid(c, tp, oid, dmids)
	if err != nil {
		return
	}
	for dmid, idx := range idxMap {
		if !model.IsDMEditAble(idx.State) {
			continue
		}
		idx.AttrSet(value, bit)
		if _, err = s.dao.UpdateDMAttr(c, tp, oid, dmid, idx.Attr); err != nil {
			continue
		}
		s.OpLog(c, oid, mid, time.Now().Unix(), int(tp), []int64{dmid}, "attribute", "", fmt.Sprintf("bit:%d,value:%d", bit, value), "弹幕保护状态变更", source, operatorType)
		if isAssist {
			s.addAssistLog(sub.Mid, mid, oid, assmdl.ActProtect, dmid)
		}
		affectIds = append(affectIds, dmid)
	}
	return
}

func (s *Service) addAssistLog(mid, assistMid, oid, action, dmid int64) {
	ct, err := s.dao.Content(context.TODO(), oid, dmid)
	if err != nil || ct == nil {
		return
	}
	detail := ct.Msg
	if len([]rune(ct.Msg)) > 50 {
		detail = string([]rune(ct.Msg)[:50])
	}
	arg := &assmdl.ArgAssistLogAdd{
		Mid:       mid,
		AssistMid: assistMid,
		Type:      assmdl.TypeDm,
		Action:    action,
		SubjectID: oid,
		ObjectID:  fmt.Sprint(dmid),
		Detail:    detail,
	}
	select {
	case s.assistLogChan <- arg:
	default:
		log.Error("assistLogChan is full")
	}
}

func (s *Service) assistLogproc() {
	for arg := range s.assistLogChan {
		if err := s.assRPC.AssistLogAdd(context.TODO(), arg); err != nil {
			log.Error("assRPC.AssistLogAdd(%v) error(%v)", arg, err)
		} else {
			log.Info("assRPC.AssistLogAdd(%v) success", arg)
		}
	}
}

// updateMonitorCnt update mcount of subject.
func (s *Service) updateMonitorCnt(c context.Context, sub *model.Subject) (err error) {
	var state, mcount int64
	if sub.AttrVal(model.AttrSubMonitorBefore) == model.AttrYes {
		state = int64(model.StateMonitorBefore)
	} else if sub.AttrVal(model.AttrSubMonitorAfter) == model.AttrYes {
		state = int64(model.StateMonitorAfter)
	} else {
		return
	}
	if mcount, err = s.dao.DMCount(c, sub.Type, sub.Oid, []int64{state}); err != nil {
		return
	}
	_, err = s.dao.UpSubjectMCount(c, sub.Type, sub.Oid, mcount)
	return
}
