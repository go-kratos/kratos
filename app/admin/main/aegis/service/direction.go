package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
)

//ShowDirection .
func (s *Service) ShowDirection(c context.Context, id int64) (r *net.ShowDirectionResult, err error) {
	var (
		d *net.Direction
		f *net.Flow
		t *net.Transition
	)

	if d, err = s.gorm.DirectionByID(c, id); err != nil {
		return
	}
	if f, err = s.gorm.FlowByID(c, d.FlowID); err != nil {
		return
	}
	if t, err = s.gorm.TransitionByID(c, d.TransitionID); err != nil {
		return
	}
	r = &net.ShowDirectionResult{
		Direction:      d,
		FlowName:       f.ChName,
		TransitionName: t.ChName,
	}
	return
}

//GetDirectionList .
func (s *Service) GetDirectionList(c context.Context, pm *net.ListDirectionParam) (result *net.ListDirectionRes, err error) {
	var (
		n      *net.Net
		flows  map[int64]string
		trans  map[int64]string
		unames map[int64]string
		fids   = []int64{}
		tids   = []int64{}
		uid    = []int64{}
	)
	if n, err = s.gorm.NetByID(c, pm.NetID); err != nil {
		return
	}
	if result, err = s.gorm.DirectionList(c, pm); err != nil {
		return
	}
	if len(result.Result) == 0 {
		return
	}
	for _, item := range result.Result {
		fids = append(fids, item.FlowID)
		tids = append(tids, item.TransitionID)
		uid = append(uid, item.UID)
	}
	if flows, err = s.gorm.ColumnMapString(c, net.TableFlow, "ch_name", fids, ""); err != nil {
		return
	}
	if trans, err = s.gorm.ColumnMapString(c, net.TableTransition, "ch_name", tids, ""); err != nil {
		return
	}
	if unames, err = s.http.GetUnames(c, uid); err != nil {
		log.Error("GetDirectionList s.http.GetUnames error(%v)", err)
		err = nil
	}

	for _, item := range result.Result {
		item.NetName = n.ChName
		item.FlowName = flows[item.FlowID]
		item.TransitionName = trans[item.TransitionID]
		item.UserName = unames[item.UID]
	}
	return
}

//SwitchDirection .
func (s *Service) SwitchDirection(c context.Context, id int64, needDisable bool) (err error, msg string) {
	var (
		old       *net.Direction
		action    string
		canUpdate bool
	)
	if old, err = s.gorm.DirectionByID(c, id); err != nil {
		log.Error("SwitchDirection s.gorm.DirectionByID(%d) error(%v) needDisable(%v)", id, err, needDisable)
		return
	}
	available := old.IsAvailable()
	if available == !needDisable {
		return
	}

	if canUpdate, err = s.beforeUpdate(c, old.NetID); err != nil {
		log.Error("SwitchDirection s.beforeUpdate error(%v)", err)
		return
	}
	if !canUpdate {
		log.Error("SwitchDirection can't update id(%d) needdisable(%v)", id, needDisable)
		return
	}

	if needDisable {
		old.DisableTime = time.Now()
		action = model.LogNetActionDisable
	} else {
		if err = s.checkDirectionBindAvailable(c, old.FlowID, old.TransitionID); err != nil {
			log.Error("SwitchDirection s.checkDirectionBindAvailable(%+v) error(%v) needDisable(%v)", old, err, needDisable)
			return
		}
		if err, msg = s.checkDirConflict(c, old); err != nil {
			log.Error("SwitchDirection s.checkDirConflict(%+v) error(%v) needDisable(%v)", old, err, needDisable)
			return
		}
		old.DisableTime = net.Recovered
		action = model.LogNetActionAvailable
	}

	tx, err := s.gorm.BeginTx(c)
	if err != nil {
		log.Error("SwitchDirection s.gorm.BeginTx error(%v)", err)
		return
	}
	if err = s.gorm.DisableNet(c, tx, old.NetID); err != nil {
		tx.Rollback()
		return
	}
	if err = s.gorm.UpdateFields(c, tx, net.TableDirection, id, map[string]interface{}{"disable_time": old.DisableTime}); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("SwitchDirection tx.Commit error(%v)", err)
		return
	}
	s.delDirCache(c, old)

	//日志
	oper := &model.NetConfOper{
		OID:    old.ID,
		Action: action,
		UID:    old.UID,
		NetID:  old.NetID,
		FlowID: old.FlowID,
		TranID: old.TransitionID,
	}
	s.sendNetConfLog(c, model.LogTypeDirConf, oper)
	return
}

func (s *Service) checkDirectionUnique(c context.Context, netID int64, FlowID int64, transitionID int64, direction int8) (err error, msg string) {
	var exist *net.Direction
	if exist, err = s.gorm.DirectionByUnique(c, netID, FlowID, transitionID, direction); err != nil {
		return
	}
	if exist != nil {
		err = ecode.AegisUniqueAlreadyExist
		msg = fmt.Sprintf(ecode.AegisUniqueAlreadyExist.Message(), "有向线", "")
	}
	return
}

func (s *Service) checkDirectionBindAvailable(c context.Context, flowID, transitionID int64) (err error) {
	var (
		f *net.Flow
		t *net.Transition
	)

	if flowID > 0 {
		if f, err = s.gorm.FlowByID(c, flowID); err != nil {
			return
		}
		if !f.IsAvailable() {
			err = ecode.AegisFlowDisabled
			return

		}
	}

	if transitionID > 0 {
		if t, err = s.gorm.TransitionByID(c, transitionID); err != nil {
			return
		}
		if !t.IsAvailable() {
			err = ecode.AegisTranDisabled
			return
		}
	}

	return
}

/**
 * 不同order兼容性不同：
 * order=顺序，flow<->tran一对一
 * order=todo
 */
func (s *Service) checkDirConflict(c context.Context, d *net.Direction) (err error, msg string) {
	var (
		flowDir, tranDir      []*net.Direction
		orderDes, conflictDes string
	)
	if d == nil {
		return
	}
	if flowDir, err = s.gorm.DirectionByFlowID(c, []int64{d.FlowID}, d.Direction); err != nil {
		log.Error("checkDirConflict s.gorm.DirectionByFlowID error(%v) direction(%+v)", err, d)
		return
	}
	flowDirLen := len(flowDir)
	for k, item := range flowDir {
		if item.ID != d.ID {
			continue
		}
		if k < flowDirLen-1 {
			flowDir = append(flowDir[:k], flowDir[k+1:]...)
		} else {
			flowDir = flowDir[:k]
		}
		flowDirLen--
	}

	if tranDir, err = s.gorm.DirectionByTransitionID(c, []int64{d.TransitionID}, d.Direction, true); err != nil {
		log.Error("checkDirConflict s.gorm.DirectionByTransitionID error(%v) direction(%+v)", err, d)
		return
	}
	tranDirLen := len(tranDir)
	for k, item := range tranDir {
		if item.ID != d.ID {
			continue
		}
		if k < tranDirLen-1 {
			tranDir = append(tranDir[:k], tranDir[k+1:]...)
		} else {
			tranDir = tranDir[:k]
		}
		tranDirLen--
	}
	//无任何已有线
	if flowDirLen == 0 && tranDirLen == 0 {
		return
	}

	if d.Order == net.DirOrderSequence {
		conflictDes = "节点或变迁已被绑定"
	} else {
		//不支持的顺序报错
		log.Error("checkDirConflict order(%d) is not supported! direction(%+v)", d.Order, d)
		err = ecode.RequestErr
		return
	}

	//顺序兼容性报错
	log.Error("checkDirConflict direction(%+v) not allowed!", d)
	err = ecode.AegisDirOrderConflict
	orderDes = net.DirOrderDesc[d.Order]
	msg = fmt.Sprintf(ecode.AegisDirOrderConflict.Message(), orderDes, conflictDes)
	return
}

//AddDirection .
func (s *Service) AddDirection(c context.Context, uid int64, d *net.DirEditParam) (id int64, err error, msg string) {
	var (
		canUpdate bool
	)
	if canUpdate, err = s.beforeUpdate(c, d.NetID); err != nil {
		log.Error("AddDirection s.beforeUpdate error(%v)", err)
		return
	}
	if !canUpdate {
		log.Error("AddDirection can't update param(%+v)", d)
		return
	}
	if err, msg = s.checkDirectionUnique(c, d.NetID, d.FlowID, d.TransitionID, d.Direction); err != nil {
		return
	}
	if err = s.checkDirectionBindAvailable(c, d.FlowID, d.TransitionID); err != nil {
		return
	}
	dir := &net.Direction{
		NetID:        d.NetID,
		FlowID:       d.FlowID,
		TransitionID: d.TransitionID,
		Direction:    d.Direction,
		Order:        d.Order,
		Guard:        d.Guard,
		Output:       d.Output,
		UID:          uid,
	}
	if err, msg = s.checkDirConflict(c, dir); err != nil {
		return
	}

	tx, err := s.gorm.BeginTx(c)
	if err != nil {
		log.Error("AddDirection s.gorm.BeginTx error(%v)", err)
		return
	}
	if err = s.gorm.DisableNet(c, tx, dir.NetID); err != nil {
		tx.Rollback()
		return
	}
	if err = s.gorm.AddItem(c, tx, dir); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("AddDirection tx.Commit error(%v)", err)
		return
	}
	id = d.ID

	//日志
	oper := &model.NetConfOper{
		OID:    dir.ID,
		Action: model.LogNetActionNew,
		UID:    dir.UID,
		NetID:  dir.NetID,
		FlowID: dir.FlowID,
		TranID: dir.TransitionID,
		Diff: []string{
			model.LogFieldTemp(model.LogFieldDirection, net.DirDirectionDesc[dir.Direction], "", false),
			model.LogFieldTemp(model.LogFieldOrder, net.DirOrderDesc[dir.Order], "", false),
			model.LogFieldTemp(model.LogFieldGuard, dir.Guard, "", false),
			model.LogFieldTemp(model.LogFieldOutput, dir.Output, "", false),
		},
	}
	s.sendNetConfLog(c, model.LogTypeDirConf, oper)
	return
}

//UpdateDirection .
func (s *Service) UpdateDirection(c context.Context, uid int64, d *net.DirEditParam) (err error, msg string) {
	var (
		old                                  *net.Direction
		canUpdate, checkUnique, orderChanged bool
		updates                              = map[string]interface{}{}
		diff                                 = []string{}
	)
	if old, err = s.gorm.DirectionByID(c, d.ID); err != nil {
		log.Error("UpdateDirection s.gorm.DirectionByID(%d) error(%v)", d.ID, err)
		return
	}
	if canUpdate, err = s.beforeUpdate(c, old.NetID); err != nil {
		log.Error("UpdateDirection s.beforeUpdate error(%v)", err)
		return
	}
	if !canUpdate {
		log.Error("UpdateDirection can't update param(%+v)", d)
		return
	}

	cp := *old
	nw := &cp
	if d.FlowID != old.FlowID {
		if err = s.checkDirectionBindAvailable(c, d.FlowID, 0); err != nil {
			return
		}
		checkUnique = true
		nw.FlowID = d.FlowID
		updates["flow_id"] = d.FlowID
	}
	if d.TransitionID != old.TransitionID {
		if err = s.checkDirectionBindAvailable(c, 0, d.TransitionID); err != nil {
			return
		}
		checkUnique = true
		nw.TransitionID = d.TransitionID
		updates["transition_id"] = d.TransitionID
	}
	if d.Direction != old.Direction {
		checkUnique = true
		diff = append(diff, model.LogFieldTemp(model.LogFieldDirection, net.DirDirectionDesc[nw.Direction], net.DirDirectionDesc[old.Direction], true))
		nw.Direction = d.Direction
		updates["direction"] = d.Direction
	}
	if checkUnique {
		if err, msg = s.checkDirectionUnique(c, nw.NetID, nw.FlowID, nw.TransitionID, nw.Direction); err != nil {
			return
		}
	}
	if d.Order != old.Order {
		diff = append(diff, model.LogFieldTemp(model.LogFieldOrder, net.DirOrderDesc[nw.Order], net.DirOrderDesc[old.Order], true))
		nw.Order = d.Order
		updates["order"] = d.Order
		orderChanged = true
	}
	//todo-- guard,output
	if checkUnique || orderChanged {
		if err, msg = s.checkDirConflict(c, nw); err != nil {
			return
		}
	}
	if len(updates) <= 0 {
		return
	}

	tx, err := s.gorm.BeginTx(c)
	if err != nil {
		log.Error("UpdateDirection s.gorm.BeginTx error(%v)", err)
		return
	}
	if err = s.gorm.DisableNet(c, tx, old.NetID); err != nil {
		tx.Rollback()
		return
	}
	if err = s.gorm.UpdateFields(c, tx, net.TableDirection, old.ID, updates); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("UpdateDirection tx.Commit error(%v)", err)
		return
	}
	s.delDirCache(c, old)

	//日志
	oper := &model.NetConfOper{
		OID:    nw.ID,
		Action: model.LogNetActionUpdate,
		UID:    nw.UID,
		NetID:  nw.NetID,
		FlowID: nw.FlowID,
		TranID: nw.TransitionID,
		Diff:   diff,
	}
	s.sendNetConfLog(c, model.LogTypeDirConf, oper)
	return
}

/**
 * isDirEnable 检查到下一步的可能性，由order&guard决定
 */
func (s *Service) isDirEnable(dir *net.Direction) (enable bool) {
	if dir == nil {
		return
	}
	if (dir.Order != net.DirOrderOrSplit && dir.Order != net.DirOrderOrResultSplit) ||
		(dir.Order == net.DirOrderOrResultSplit && dir.Guard == "") {
		enable = true
		return
	}

	//todo--compute guard expression-v2

	return
}

/**
 * dispatchDirs 方向线是否可达下一步
 * 加入资源维度：若资源没有任何可达性如何办，这是配置错误--后期加入动态检查
 */
func (s *Service) dispatchDirs(dirs []*net.Direction) (enableDir []*net.Direction) {
	var (
		enable bool
	)

	//对每个有向线, 过滤所有资源, 区分哪些资源可达
	enableDir = []*net.Direction{}
	for _, dir := range dirs {
		enable = s.isDirEnable(dir)
		if enable {
			enableDir = append(enableDir, dir)
			continue
		}
	}

	return
}

/**
 * fetchFlowNextEnableDirs 从flow出发，找到下一步变迁
 * 后期加入资源维度：若资源没有任何可达性如何办，这是配置错误--后期加入动态检查
 * 被应用在：
 * 1. 新节点查看下一步变迁是否需要创建分发任务;
 * 2. 审核提交后，通过flowid找到被触发的是哪个变迁(必须是同一个---后期加入动态检查);
 * 3. 通过flowid获取可用变迁的操作项;
 * todo -- 任务方存储了transitionid后，可以移除2+3逻辑
 */
func (s *Service) fetchFlowNextEnableDirs(c context.Context, flowID int64) (enableDir []*net.Direction, err error) {
	var (
		list []*net.Direction
	)

	//找到以flow为起点的所有方向线
	if list, err = s.dirByFlow(c, []int64{flowID}, net.DirInput); err != nil {
		log.Error("fetchFlowNextEnableDirs s.dirByFlow(%d) error(%v)", flowID, err)
		return
	}
	if len(list) == 0 { //没配置下一步，正常情况
		return
	}

	if enableDir = s.dispatchDirs(list); len(enableDir) == 0 {
		err = ecode.AegisFlowNoEnableTran
		log.Error("fetchFlowNextEnableDirs s.dispatchDirs flowid(%v) no enable transition", flowID)
	}
	return
}

/**
 * fetchTranNextEnableDirs 从变迁出发，找到下一步flow
 * 后期加入资源维度
 */
func (s *Service) fetchTranNextEnableDirs(c context.Context, tranID int64) (resultDir *net.Direction, err error) {
	var (
		list      []*net.Direction
		enableDir []*net.Direction
	)
	if list, err = s.dirByTran(c, []int64{tranID}, net.DirOutput, true); err != nil {
		log.Error("fetchTranNextEnableDirs s.dirByTran(%d) error(%v)", tranID, err)
		return
	}
	if len(list) == 0 { //变迁后面没配置flow，为配置错误--后期都需要添加动态检查
		err = ecode.AegisTranNoFlow
		log.Error("fetchTranNextEnableDirs s.gorm.DirectionByTransitionID(%d) no flow", tranID)
		return
	}

	enableDir = s.dispatchDirs(list)
	if len(enableDir) == 0 { //变迁后面，没有任何可用flow，为配置错误---后期需添加动态检查
		err = ecode.AegisTranNoFlow
		log.Error("fetchTranNextEnableDirs transition(%d) has no enable flow", tranID)
	}

	//todo--允许变迁和节点的一对多对应吗？
	resultDir = enableDir[0]
	return
}

func (s *Service) beforeUpdate(c context.Context, netID int64) (ok bool, err error) {
	var (
		fr     *net.FlowResource
		flowID []int64
	)

	if flowID, err = s.flowIDByNet(c, netID); err != nil {
		return
	}

	if len(flowID) > 0 {
		if fr, err = s.gorm.FRByFlow(c, flowID); err != nil {
			log.Error("beforeUpdate s.gorm.FRByFlow error(%v) netid(%d)", err, netID)
			return
		}
	}

	ok = fr == nil
	return
}
