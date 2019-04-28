package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/business"
	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/app/admin/main/aegis/model/resource"
	taskmod "go-common/app/admin/main/aegis/model/task"
	uprpc "go-common/app/service/main/up/api/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

// ListBizFlow .
func (s *Service) ListBizFlow(c context.Context, tp int8, bizID []int64, flowID []int64) (list []*business.BizItem, err error) {
	// 1. 获取每个类型的业务
	list = []*business.BizItem{}
	res, err := s.gorm.BusinessList(c, tp, bizID, true)
	if err != nil {
		err = ecode.AegisBusinessCfgErr
		return
	}
	if len(res) == 0 {
		return
	}
	bizMap := map[int64]*business.Business{}
	accessBiz := []int64{}
	for _, item := range res {
		bizMap[item.ID] = item
		accessBiz = append(accessBiz, item.ID)
	}

	//获取指定业务下的可分发流程节点
	var bizFlow map[int64]map[int64]string
	if bizFlow, err = s.dispatchFlow(c, accessBiz, flowID); err != nil {
		return
	}
	for bizID, item := range bizFlow {
		if bizMap[bizID] == nil {
			continue
		}

		bizItem := &business.BizItem{
			BizID:   bizID,
			BizName: bizMap[bizID].Name,
			BizType: bizMap[bizID].TP,
			Flows:   item,
		}
		list = append(list, bizItem)
	}

	sort.Sort(business.BizItemArr(list))
	return
}

// Next reveive next auditing task
func (s *Service) Next(c context.Context, opt *taskmod.NextOptions) (infos []*model.AuditInfo, err error) {
	log.Info("Next opt(%+v)", opt)
	if opt.Debug == 1 {
		info := model.GetEmptyInfo()
		infos = append(infos, info)
		return
	}
	// 1. TODO: 根据业务 动态判定dispatch_count和seize_count
	/*
		获取任务后，要根据资源状态判断释放需要审核，不用审核的直接关任务。 再次循环获取
	*/
	var (
		tasks []*taskmod.Task
	)
	for i := 0; i < 3; i++ {
		tasks, _, err = s.NextTask(c, opt)
		if err != nil {
			log.Error("NextTask err(%v)", err)
			err = ecode.AegisTaskErr
			return
		}
		if len(tasks) == 0 {
			log.Info("NextTask empty task")
			return
		}

		for _, task := range tasks {
			log.Info("Next task(%+v)", task)
			info, err := s.auditInfoByTask(c, task, &opt.BaseOptions)
			if err != nil {
				err = nil
				continue
			}
			infos = append(infos, info)
		}
		if len(infos) > 0 {
			break
		}
	}

	return
}

// InfoTask .
func (s *Service) InfoTask(c context.Context, opt *common.BaseOptions, taskid int64) (info *model.AuditInfo, err error) {
	task, err := s.Task(c, taskid)
	if err != nil || task == nil {
		err = ecode.AegisTaskErr
		return
	}

	return s.auditInfoByTask(c, task, opt)
}

// ListByTask 任务列表
func (s *Service) ListByTask(c context.Context, opt *taskmod.ListOptions) (list []*model.ListTaskItem, err error) {
	log.Info("ListByTask opt(%+v)", opt)
	var (
		tasks []*taskmod.Task
		count int64
	)

	if opt.Debug == 1 {
		list = []*model.ListTaskItem{{}}
		opt.Total = 1
	} else {
		tasks, count, err = s.ListTasks(c, opt)
		if err != nil {
			log.Error("ListTasks err(%v)", err)
			err = ecode.AegisTaskErr
			return
		}
		if len(tasks) == 0 {
			return
		}
		opt.Total = int(count)
	}

	for _, task := range tasks {
		item := &model.ListTaskItem{}

		item.Task = task
		if int(time.Since(task.Gtime.Time()).Minutes()) < 10 {
			item.GTstr = common.WaitTime(task.Gtime.Time())
		}
		item.CTstr = task.Ctime.Time().Format("2006-01-02 15:04:05")
		item.MTstr = task.Mtime.Time().Format("2006-01-02 15:04:05")
		item.WaitTime = common.WaitTime(task.Ctime.Time())
		ids, _ := xstr.SplitInts(task.Group)
		item.UserGroup = s.getUserGroup(c, ids)
		list = append(list, item)
	}

	// 补充oid,content
	s.mulIDtoName(c, list, s.gorm.ListHelperForTask, "RID", "OID", "Content", "Metas")
	// 补充user_info,user_group
	s.mulIDtoName(c, list, s.listHelpUser, "MID", "UserGroup", "UserInfo")
	// 补充uname
	s.mulIDtoName(c, list, s.transUnames, "UID", "UserName")
	// 将mid替换为昵称 或者 cuser
	for _, item := range list {
		if item.MID != 0 && item.UserInfo != nil {
			item.MidStr = item.UserInfo.Name
		} else if item.Metas != nil {
			if val, ok := item.Metas["cuser"]; ok {
				item.MidStr = fmt.Sprint(val)
			}
		}
	}
	return
}

// InfoResource .
func (s *Service) InfoResource(c context.Context, opt *common.BaseOptions) (info *model.AuditInfo, err error) {
	if opt.Debug == 1 {
		return model.GetEmptyInfo(), nil
	}
	// 根据oid查资源
	rsc, err := s.gorm.ResByOID(c, opt.BusinessID, opt.OID)
	if err != nil || rsc == nil {
		err = ecode.AegisResourceErr
		return
	}

	return s.auditInfoByRsc(c, rsc, opt.NetID)
}

// ListByResource 资源列表
func (s *Service) ListByResource(c context.Context, arg *model.SearchParams) (columns []*model.Column, list []*model.ListRscItem, op []*net.TranOperation, err error) {
	if arg.Debug == 1 {
		arg.Total = 1
		return []*model.Column{}, []*model.ListRscItem{model.EmptyListItem()}, []*net.TranOperation{{}}, nil
	}
	// 搜索返回资源信息
	var sres *model.SearchRes
	if sres, err = s.http.ResourceES(c, arg); err != nil || sres == nil {
		err = ecode.AegisSearchErr
		return
	}
	arg.Total = sres.Page.Total
	if arg.FilterOff {
		list = sres.Resources
	} else {
		list = s.listParseState(c, arg.State, sres.Resources)
	}

	// 补充粉丝数和分组
	if len(list) > 0 {
		g, _ := errgroup.WithContext(c)
		g.Go(func() error {
			s.mulIDtoName(c, list, s.listHelpUser, "MID", "UserGroup", "UserInfo")
			return nil
		})
		g.Go(func() error {
			s.listHightLight(c, arg.BusinessID, list)
			return nil
		})
		g.Go(func() error {
			s.mulIDtoName(c, list, s.listMetas, "ID", "MetaData", "Metas")
			return nil
		})
		g.Wait()
	}
	columns = s.getColumns(c, arg.BusinessID)

	// 批量操作项
	op, err = s.fetchBatchOperations(c, arg.BusinessID, 0)
	return
}

// listMetas 补充列表里面的 metadata metas
func (s *Service) listMetas(c context.Context, ids []int64) (res map[int64][]interface{}, err error) {
	res = make(map[int64][]interface{})
	var metas map[int64]string

	if metas, err = s.gorm.MetaByRID(c, ids); err != nil {
		return
	}
	for id, meta := range metas {
		mmeta := make(map[string]interface{})
		if len(meta) > 0 {
			if err = json.Unmarshal([]byte(meta), &mmeta); err != nil {
				log.Error("listMetas json.Unmarshal error(%v)", err)
				err = nil
			}
		}
		res[id] = []interface{}{meta, mmeta}
	}
	return
}

//状态筛选，防止搜索列表更新不及时
func (s *Service) listParseState(c context.Context, state int64, list []*model.ListRscItem) (hitlist []*model.ListRscItem) {
	var arrids []int64
	for _, item := range list {
		arrids = append(arrids, item.ID)
	}
	hitids, err := s.gorm.ResourceHit(c, arrids)
	if err != nil {
		return list
	}

	//搜索待审列表时,过滤掉已被领取的任务，避免提交冲突
	taskhitids, _ := s.gorm.TaskHitAuditing(c, arrids)
	for _, item := range list {
		if _, ok := taskhitids[item.ID]; ok {
			continue
		}
		if st, ok := hitids[item.ID]; ok {
			if state == -12345 || (state == st) {
				item.State = st
				hitlist = append(hitlist, item)
			}
		}
	}
	return
}

// listHelpUser 补充列表里面的 user_info user_group
func (s *Service) listHelpUser(c context.Context, mids []int64) (res map[int64][]interface{}, err error) {
	res = make(map[int64][]interface{})
	//mids去零
	for i, v := range mids {
		if v > 0 {
			continue
		}

		if i == len(mids)-1 {
			mids = mids[:i]
		} else {
			mids = append(mids[:i], mids[i+1:]...)
		}
	}
	if len(mids) == 0 {
		return
	}

	infos, err := s.rpc.UserInfos(c, mids)
	if err != nil {
		infos = make(map[int64]*model.UserInfo)
	}

	upspecials, err := s.rpc.UpsSpecial(c, mids)
	if err != nil {
		upspecials = make(map[int64]*uprpc.UpSpecial)
	}

	for _, mid := range mids {
		gids := []int64{}
		if gs, ok := upspecials[mid]; ok {
			gids = gs.GroupIDs
		}
		res[mid] = []interface{}{s.getUserGroup(c, gids), infos[mid]}
	}
	log.Info("listRscHelper res(%+v)", res)
	return
}

/*
避免超时，控制文本长度小于3000
过滤相同的文本，减少不必要请求
*/
func (s *Service) listHightLight(c context.Context, bizid int64, list []*model.ListRscItem) {

	var (
		area string
		err  error
	)

	if area = s.getConfig(c, bizid, business.TypeFiler); len(area) == 0 {
		log.Warn("sigleHightLight(%d) 没有文本高亮配置(%v)", bizid, err)
		return
	}

	arrcontent := []string{}
	for _, item := range list {
		arrcontent = append(arrcontent, item.Content)
	}
	arrset := stringset(arrcontent)
	hits, err := s.concurrentHightList(c, area, arrset)
	if err != nil {
		log.Error("listHightLight error(%v)", err)
		return
	}

	hitset := stringset(hits)
	for _, item := range list {
		item.Hit = hitset
	}
}

func (s *Service) concurrentHightList(c context.Context, area string, mapset []string) (hits []string, err error) {
	var eg errgroup.Group
	msgs := joinstr(mapset, "msg=", 3000)
	for _, msg := range msgs {
		var m = msg
		eg.Go(func() error {
			var (
				e   error
				hit []string
			)
			hit, e = s.http.FilterMulti(context.Background(), area, m)
			hits = append(hits, hit...)
			return e
		})
	}
	err = eg.Wait()
	return
}

func (s *Service) sigleHightLight(c context.Context, bizid int64, content string) (hit []string) {
	var (
		area string
		err  error
	)

	if area = s.getConfig(c, bizid, business.TypeFiler); len(area) == 0 {
		log.Warn("sigleHightLight(%d) 没有文本高亮配置(%v)", bizid, err)
		return
	}

	hits, err := s.http.FilterMulti(c, area, "msg="+content)
	if err != nil {
		log.Error("sigleHightLight error(%v)", err)
		return
	}
	return hits
}

func (s *Service) getColumns(c context.Context, bizid int64) (columns []*model.Column) {
	cfg := s.getConfig(c, bizid, business.TypeRscListAdapter)
	if len(cfg) == 0 {
		log.Warn("getColumns empty config")
		return
	}
	if err := json.Unmarshal([]byte(cfg), &columns); err != nil {
		log.Error("getColumns err(%v)", err)
	}
	return
}

//Submit .
func (s *Service) Submit(c context.Context, opt *model.SubmitOptions) (err error) {
	log.Info("Pre submit(%+v)", opt)

	// 1. 获取flow流转结果
	result, err := s.computeTriggerResult(c, opt.RID, opt.FlowID, opt.Binds)
	if err != nil {
		return err
	}

	esupsert, err := s.submit(c, "submit", opt, result)
	if err != nil {
		return err
	}

	//更新es
	s.http.UpsertES(c, []*model.UpsertItem{esupsert})
	return
}

func (s *Service) submit(c context.Context, action string, opt *model.SubmitOptions, result *net.TriggerResult) (esupsert *model.UpsertItem, err error) {
	ormTx, err := s.gorm.BeginTx(c)
	if err != nil {
		log.Error("tx error(%v)", err)
		return nil, ecode.ServerErr
	}

	var (
		rsc                *resource.Resource
		mid, ouid, otaskid int64
		ostate             int8
		taskstate          int8
	)

	if rsc, err = s.gorm.ResourceByOID(c, opt.OID, opt.BusinessID); err != nil || rsc == nil {
		ormTx.Rollback()
		err = ecode.AegisResourceErr
		return
	}

	if rsc.ID != opt.RID {
		ormTx.Rollback()
		err = ecode.AegisResourceErr
		return
	}
	mid = rsc.MID

	// 2. 更新resource 和 resource_result
	if err = s.TxUpdateResource(ormTx, opt.RID, result.ResultToken.Values, opt.Result); err != nil {
		log.Error("submit TxUpdateResource(%v)", err)
		ormTx.Rollback()
		err = ecode.AegisResourceErr
		return
	}
	log.Info("submit TxUpdateResource success")

	// 3. 更新task
	if result.ResultToken.HitAudit && opt.TaskID > 0 {
		taskstate = taskmod.TaskStateSubmit
		if action == "batch" { // 目前的接口我们区分不了 资源提交还是任务提交
			taskstate = taskmod.TaskStateRscSb
		}

		// TODO 看来还是要把资源提交和任务提交分开。 资源提交不能确定要关的是哪个任务
		if ostate, otaskid, ouid, err = s.TxSubmitTask(c, ormTx, &opt.BaseOptions, taskstate); err != nil {
			log.Error("submit TxSubmitTask(%v)", err)
			ormTx.Rollback()
			err = ecode.AegisTaskErr
			return
		}
		log.Info("submit TxSubmitTask success")
		if otaskid != opt.TaskID {
			log.Warn("submit different taskid(%d-->%d)", opt.TaskID, otaskid)
			opt.TaskID = otaskid
		}
	}

	if result.ResultToken.HitAudit {
		// 4. 更新flow
		if err = s.reachNewFlowDB(c, ormTx, result); err != nil {
			log.Error("submit reachNewFlowDB(%v)", err)
			ormTx.Rollback()
			return
		}
		log.Info("submit reachNewFlowDB success")

		opt.Result.State, _ = strconv.Atoi(fmt.Sprint(result.ResultToken.Values["state"]))
	}

	// 5. 更新业务
	if err = s.syncResource(c, opt, mid, result.ResultToken); err != nil {
		log.Error("submit syncResource(%v)", err)
		ormTx.Rollback()
		return
	}
	log.Info("submit syncResource success")

	if err = ormTx.Commit().Error; err != nil {
		ormTx.Rollback()
		return
	}

	// 6. 任务缓存更新
	if result.ResultToken.HitAudit && opt.TaskID > 0 {
		s.submitTaskCache(c, &opt.BaseOptions, ostate, otaskid, ouid)
	}

	// 7. 任务流转
	if err = s.afterReachNewFlow(c, result, opt.BusinessID); err != nil {
		log.Error("submit afterReachNewFlow(%v)", err)
		return
	}
	log.Info("submit afterReachNewFlow success")

	// 8. 查询下最终数据库的结果,记录
	res, err := s.gorm.ResourceRes(c, opt.RID)
	if err != nil || res == nil {
		return
	}

	//8. 记录日志
	if opt.Result != nil {
		opt.Result.State = int(res.State)
	}
	s.logSubmit(c, action, opt, result)

	//9.更新es
	esupsert = &model.UpsertItem{
		ID:     res.ID,
		State:  int(res.State),
		Extra1: res.Extra1,
		Extra2: res.Extra2,
		Extra3: res.Extra3,
		Extra4: res.Extra4,
	}
	return
}

func (s *Service) logSubmit(c context.Context, action string, opt *model.SubmitOptions, res interface{}) {
	s.async.Do(c, func(ctx context.Context) {
		// 1. 操作日志
		s.sendAuditLog(ctx, action, opt, res, model.LogTypeAuditSubmit)
		// 2. resource日志
		s.sendRscSubmitLog(ctx, action, opt, res)
	})
}

// JumpFlow 跳流程提交
func (s *Service) JumpFlow(c context.Context, opt *model.SubmitOptions) (err error) {
	var (
		ostate        int8
		ouid, otaskid int64
		esupsert      *model.UpsertItem
	)

	ormTx, err := s.gorm.BeginTx(c)
	if err != nil {
		log.Error("tx error(%v)", err)
		return ecode.ServerErr
	}

	// 1. 提交结果到flow
	result, err := s.jumpFlow(c, ormTx, opt.RID, opt.FlowID, opt.NewFlowID, opt.Binds)
	if err != nil {
		ormTx.Rollback()
		return
	}

	// 2. 更新task
	if result.ResultToken.HitAudit && opt.TaskID > 0 { //资源列表直接提交可能没有任务信息
		if ostate, otaskid, ouid, err = s.TxSubmitTask(c, ormTx, &opt.BaseOptions, taskmod.TaskStateClosed); err != nil {
			ormTx.Rollback()
			err = ecode.AegisTaskErr
			return
		}
		if otaskid != opt.TaskID {
			log.Warn("submit different taskid(%d-->%d)", opt.TaskID, otaskid)
			opt.TaskID = otaskid
		}
	}

	if result.SubmitToken != nil {
		var (
			rsc *resource.Resource
			mid int64
		)
		if rsc, err = s.gorm.ResourceByOID(c, opt.OID, opt.BusinessID); err != nil || rsc == nil {
			ormTx.Rollback()
			err = ecode.AegisResourceErr
			return
		}

		if rsc.ID != opt.RID {
			ormTx.Rollback()
			err = ecode.AegisResourceErr
			return
		}
		mid = rsc.MID

		// 3. 更新resource 和 resource_result
		if err = s.TxUpdateResource(ormTx, opt.RID, result.ResultToken.Values, opt.Result); err != nil {
			ormTx.Rollback()
			err = ecode.AegisResourceErr
			return
		}

		// 4. 更新业务
		if err = s.syncResource(c, opt, mid, result.ResultToken); err != nil {
			ormTx.Rollback()
			err = ecode.AegisBusinessSyncErr
			return
		}

	}

	if err = ormTx.Commit().Error; err != nil {
		return
	}

	// 5. 任务缓存更新
	if result.ResultToken.HitAudit && opt.TaskID > 0 {
		s.submitTaskCache(c, &opt.BaseOptions, ostate, otaskid, ouid)
	}

	// 6. 任务流转
	if result.ResultToken.HitAudit {
		s.afterJumpFlow(c, result, opt.BusinessID)
	}

	rscr, err := s.gorm.ResourceRes(c, opt.RID)
	if err != nil || rscr != nil {
		return
	}
	//7. 记录日志
	if opt.Result != nil {
		opt.Result.State = int(rscr.State)
	}
	s.logSubmit(c, "jump", opt, result)

	//8.更新es
	esupsert = &model.UpsertItem{
		ID:     rscr.ID,
		State:  int(rscr.State),
		Extra1: rscr.Extra1,
		Extra2: rscr.Extra2,
		Extra3: rscr.Extra3,
		Extra4: rscr.Extra4,
	}
	s.http.UpsertES(c, []*model.UpsertItem{esupsert})

	return
}

// BatchSubmit 批量提交, 超过10个的做异步
func (s *Service) BatchSubmit(c context.Context, opt *model.BatchOption) (tip *model.Tip, err error) {
	tip = &model.Tip{Fail: make(map[int64]string)}

	if len(opt.RIDs) > 10 {
		go s.processBatch(context.Background(), opt.RIDs[10:], opt, nil)
		rids := opt.RIDs[:10]
		s.processBatch(c, rids, opt, tip)
		tip.Async = append(tip.Async, opt.RIDs[10:]...)
	} else {
		s.processBatch(c, opt.RIDs, opt, tip)
	}

	return
}

func (s *Service) processBatch(c context.Context, rids []int64, opt *model.BatchOption, tip *model.Tip) {
	var (
		err       error
		esupdate  *model.UpsertItem
		esupdates = []*model.UpsertItem{}
	)
	for _, rid := range rids {
		var (
			taskid int64
			oid    string
		)
		if oid, err = s.gorm.OidByRID(c, rid); err != nil {
			if tip != nil {
				tip.Fail[rid] = err.Error()
			}
			continue
		}

		res, err := s.computeBatchTriggerResult(c, opt.BusinessID, rid, opt.Binds)
		if err != nil {
			if tip != nil {
				tip.Fail[rid] = ecode.Cause(err).Message()
			}
			continue
		}
		if task, _ := s.gorm.TaskByRID(c, rid, 0); task != nil {
			taskid = task.ID
		}

		smtOpt := &model.SubmitOptions{
			EngineOption: model.EngineOption{
				BaseOptions: common.BaseOptions{
					BusinessID: opt.BusinessID,
					OID:        oid,
					UID:        opt.UID,
					RID:        rid,
					Uname:      opt.Uname,
				},
				Result: &resource.Result{
					Attribute:    -1, // -1表示不更新
					ReasonID:     opt.ReasonID,
					RejectReason: opt.RejectReason,
				},
				TaskID: taskid,
				ExtraData: map[string]interface{}{
					"notify": opt.Notify,
				},
			},
			Binds: opt.Binds,
		}
		if esupdate, err = s.submit(c, "batch", smtOpt, res); err != nil {
			if tip != nil {
				tip.Fail[rid] = ecode.Cause(err).Message()
			}
		} else {
			if tip != nil {
				tip.Success = append(tip.Success, rid)
			}
			esupdates = append(esupdates, esupdate)
		}
	}

	//更新es
	s.http.UpsertES(c, esupdates)
}

// Add 业务方添加资源
func (s *Service) Add(c context.Context, opt *model.AddOption) (err error) {
	var res *net.TriggerResult
	defer func() {
		// 5. 记录资源添加日志
		s.sendRscLog(c, "add", opt, res, nil, err)
	}()

	business.AdaptAddOpt(opt, s.getAdapter(c, opt.BusinessID))
	if b, _ := s.gorm.Business(c, opt.BusinessID); b == nil || b.State != 0 {
		err = ecode.AegisBusinessSyncErr
		return
	}

	var rid int64
	ormTx, err := s.gorm.BeginTx(c)
	if err != nil {
		log.Error("tx error(%v)", err)
		return ecode.ServerErr
	}

	// 2. 根据net_id,计算是否可添加
	res, err = s.startNet(c, opt.BusinessID, opt.NetID)
	if err != nil {
		ormTx.Rollback()
		return
	}

	// 3. 事务新增资源
	rsc := &resource.Result{State: opt.State}
	if res.ResultToken != nil {
		if state, ok := res.ResultToken.Values["state"]; ok {
			rsc.State, _ = strconv.Atoi(fmt.Sprint(state))
		}
	}

	rid, err = s.TxAddResource(ormTx, &opt.Resource, rsc)
	if err != nil {
		ormTx.Rollback()
		err = ecode.AegisResourceErr
		return
	}
	res.RID = rid

	// 4. 事务新增flow,根据创建rid
	if err = s.reachNewFlowDB(c, ormTx, res); err != nil {
		ormTx.Rollback()
		return
	}

	if err = ormTx.Commit().Error; err != nil {
		log.Error("Commit opt(%+v) error(%v)", opt, err)
		return
	}

	go func() {
		s.afterReachNewFlow(context.TODO(), res, opt.BusinessID)

		// 1. 操作日志
		s.sendAuditLog(context.TODO(), "add", &model.SubmitOptions{
			EngineOption: model.EngineOption{
				BaseOptions: common.BaseOptions{
					BusinessID: opt.BusinessID,
					FlowID:     res.NewFlowID,
					RID:        rid,
					UID:        399,
					Uname:      "业务方",
				},
				Result: rsc,
			},
		}, res, model.LogTypeAuditAdd)
	}()
	return
}

// Update 业务方修改资源参数
func (s *Service) Update(c context.Context, opt *model.UpdateOption) (err error) {
	var rows int64
	business.AdaptUpdateOpt(opt, s.getAdapter(c, opt.BusinessID))
	for key, val := range opt.Update {
		if _, ok := model.UpdateKeys[key]; !ok {
			delete(opt.Update, key)
			continue
		}
		switch key {
		case "extra1", "extra2", "extra3", "extra4", "extra5", "extra6":
			if reflect.TypeOf(val).Kind() != reflect.Float64 && reflect.TypeOf(val).Kind() != reflect.Int {
				return ecode.RequestErr
			}
		case "extra1s", "extra2s", "extra3s", "extra4s":
			if reflect.TypeOf(val).Kind() != reflect.String {
				return ecode.RequestErr
			}
		case "extratime1", "octime", "ptime": //时间类型不校验了
		case "metadata": //如果是map，json转化为字符串
			if reflect.TypeOf(val).Kind() == reflect.Map {
				var bs []byte
				if bs, err = json.Marshal(val); err != nil {
					return ecode.RequestErr
				}
				opt.Update["metadata"] = string(bs)
			}
		}
	}
	if len(opt.Update) == 0 {
		err = ecode.RequestErr
	} else {
		rows, err = s.gorm.UpdateResource(c, opt.BusinessID, opt.OID, opt.Update)
	}
	if rows > 0 || err != nil {
		s.sendRscLog(c, "update", &model.AddOption{
			Resource: resource.Resource{
				BusinessID: opt.BusinessID,
				OID:        opt.OID,
			}, NetID: opt.NetID,
		}, nil, opt.Update, err)
	}
	return
}

//CancelByOper 手动删除任务
func (s *Service) CancelByOper(c context.Context, businessID int64, oids []string, uid int64, username string) (err error) {
	if err = s.Cancel(c, businessID, oids, uid, username); err != nil {
		return
	}

	//upsert es
	go func(businessID int64, oids []string) {
		var (
			esupserts []*model.UpsertItem
			ctx       = context.TODO()
		)
		if esupserts, err = s.gorm.UpsertByOIDs(ctx, businessID, oids); err != nil || len(esupserts) == 0 {
			return
		}
		s.http.UpsertES(ctx, esupserts)
	}(businessID, oids)
	return
}

// Cancel 取消相关资源的所有流程
func (s *Service) Cancel(c context.Context, businessID int64, oids []string, uid int64, username string) (err error) {
	var (
		dstate  int
		ridFlow map[int64]string
	)
	defer func() {
		// 5. 记录资源注销日志
		s.sendRscCancleLog(c, businessID, oids, uid, username, err)
	}()

	ridstr, err := s.gorm.RidsByOids(c, businessID, oids)
	if err != nil {
		err = nil
		return
	}
	rids, err := xstr.SplitInts(ridstr)
	if err != nil || len(rids) == 0 {
		err = nil
		return
	}

	ormTx, err := s.gorm.BeginTx(c)
	if err != nil {
		log.Error("tx error(%v)", err)
		return ecode.ServerErr
	}
	defer ormTx.Commit()

	// 1. 资源修改为已删除状态
	if dstate, err = s.TxDelResource(c, ormTx, businessID, rids); err != nil {
		ormTx.Rollback()
		err = ecode.AegisResourceErr
		return
	}

	// 2. 取消相关flow流程
	if ridFlow, err = s.cancelNet(c, ormTx, rids); err != nil {
		ormTx.Rollback()
		return
	}

	// 3. 取消相关task
	if err = s.gorm.TxCloseTasks(ormTx, rids, uid); err != nil {
		ormTx.Rollback()
		err = ecode.AegisTaskErr
		return
	}

	// 操作日志
	for _, rid := range rids {
		s.sendAuditLog(c, "cancel", &model.SubmitOptions{
			EngineOption: model.EngineOption{
				BaseOptions: common.BaseOptions{
					BusinessID: businessID,
					RID:        rid,
					UID:        uid,
					Uname:      username,
				},
				Result: &resource.Result{
					State: dstate,
				},
			},
		}, &net.TriggerResult{
			OldFlowID: ridFlow[rid],
		}, model.LogTypeAuditCancel)
	}
	return
}

func (s *Service) closeTask(c context.Context, task *taskmod.Task) (err error) {
	if err = s.gorm.CloseTask(c, task.ID); err != nil {
		return
	}
	if task.UID > 0 {
		s.submitTaskCache(c, &common.BaseOptions{
			BusinessID: task.BusinessID,
			FlowID:     task.FlowID,
		}, task.State, task.ID, task.UID)
	}
	return
}

// Upload to bfs
func (s *Service) Upload(c context.Context, fileName string, fileType string, timing int64, body []byte) (location string, err error) {
	if len(body) == 0 {
		err = ecode.FileNotExists
		return
	}
	if location, err = s.http.Upload(c, fileName, fileType, timing, body); err != nil {
		log.Error("s.upload.Upload() error(%v)", err)
	}
	return
}

func (s *Service) auditInfoByTask(c context.Context, task *taskmod.Task, opt *common.BaseOptions) (info *model.AuditInfo, err error) {
	var (
		userinfo          *model.UserInfo
		operhistorys, hit []string
		iframeurl         string
		actions           []*model.Action
	)
	// 1. 未完成的任务信息
	undoStat, _ := s.UnDoStat(c, opt)

	// 2. 获取业务信息 resource
	resource, err := s.ResourceRes(c, &resource.Args{RID: task.RID})
	if err != nil || resource == nil || !s.checkaudit(task.FlowID, resource.State) {
		log.Error("资源查找失败,删除任务 ResourceRes err(%v)", err)
		s.closeTask(c, task)
		err = ecode.AegisResourceErr
		return
	}

	// 3. 获取操作项 flow
	flow, err := s.fetchTaskTranInfo(c, task.RID, task.FlowID, opt.NetID)
	if err != nil {
		log.Error("fetchTransitionInfo(%d,%d) err(%v)", task.RID, task.FlowID, err)
		return
	}

	g, _ := errgroup.WithContext(c)

	// 4. 获取用户信息 account
	g.Go(func() error {
		userinfo, err = s.rpc.Profile(c, task.MID)
		if err != nil {
			log.Error("Profile(%d) err(%v)", task.MID, err)
		}
		return nil
	})

	// 5. 审核历史
	g.Go(func() error {
		operhistorys, err = s.auditLogByRID(c, task.RID)
		if err != nil {
			log.Error("AuditLog err(%v)", err)
		}
		return nil
	})

	// 6. iframe url
	g.Go(func() error {
		iframeurl = s.getConfig(c, task.BusinessID, business.TypeIframe)
		actions = s.getActions(c, task.BusinessID)

		var attrcfg map[string]uint
		if attrcfg, err = s.AttributeCFG(c, task.BusinessID); len(attrcfg) > 0 {
			resource.AttrParse(attrcfg)
		}
		resource.MetaParse()
		return nil
	})

	// 7. filter
	g.Go(func() error {
		hit = s.sigleHightLight(c, task.BusinessID, resource.Content)
		return nil
	})
	g.Wait()

	info = &model.AuditInfo{
		UnDoStat:     undoStat,
		Task:         task,
		Resource:     resource,
		UserInfo:     userinfo,
		Flow:         flow,
		OperHistorys: operhistorys,
		IFrame:       iframeurl,
		Actions:      actions,
		Hit:          hit,
	}
	if task.MID > 0 {
		s.getUserGroup(c, []int64{task.MID})
	}

	return
}

func (s *Service) auditInfoByRsc(c context.Context, rsc *resource.Res, netid int64) (info *model.AuditInfo, err error) {
	var (
		task              *taskmod.Task
		userinfo          *model.UserInfo
		flow              *net.TransitionInfo
		iframeurl         string
		attrcfg           map[string]uint
		operhistorys, hit []string
		actions           []*model.Action
	)

	g, _ := errgroup.WithContext(c)

	// 搜索未指定flow, 则检索task
	g.Go(func() error {
		task, err = s.gorm.TaskByRID(c, rsc.ID, 0)
		if err != nil {
			err = nil
			task = nil
		}
		return nil
	})

	// 1. 获取用户信息 account
	g.Go(func() error {
		if userinfo, err = s.rpc.Profile(c, rsc.MID); err != nil {
			log.Error("Profile(%d) err(%v)", rsc.MID, err)
		}
		return nil
	})

	// 2. 获取操作项 flow
	g.Go(func() error {
		flow, err = s.fetchResourceTranInfo(c, rsc.ID, rsc.BusinessID, netid)
		if err != nil {
			log.Error("fetchTransitionInfo(%d,%d,%d) err(%v)", rsc.ID, rsc.BusinessID, netid, err)
		}
		return nil
	})

	// 3. 审核历史
	g.Go(func() error {
		operhistorys, err = s.auditLogByRID(c, rsc.ID)
		if err != nil {
			log.Error("AuditLog(%d) err(%v)", rsc.ID, err)
		}
		return nil
	})

	// 4. iframe url
	g.Go(func() error {
		iframeurl = s.getConfig(c, rsc.BusinessID, business.TypeIframe)
		actions = s.getActions(c, rsc.BusinessID)

		if attrcfg, err = s.AttributeCFG(c, rsc.BusinessID); len(attrcfg) > 0 {
			rsc.AttrParse(attrcfg)
		}
		rsc.MetaParse()
		return nil
	})

	// 5. filter
	g.Go(func() error {
		hit = s.sigleHightLight(c, rsc.BusinessID, rsc.Content)
		return nil
	})
	g.Wait()

	info = &model.AuditInfo{
		Task:         task,
		Resource:     rsc,
		UserInfo:     userinfo,
		Flow:         flow,
		OperHistorys: operhistorys,
		IFrame:       iframeurl,
		Actions:      actions,
		Hit:          hit,
	}

	if rsc.MID > 0 {
		if upspecial, _ := s.rpc.UpSpecial(c, rsc.MID); upspecial != nil {
			info.UserGroup = s.getUserGroup(c, upspecial.GroupIDs)
		}
	}
	return
}

func (s *Service) getActions(c context.Context, bizid int64) (actions []*model.Action) {
	cfg := s.getConfig(c, bizid, business.TypeAction)
	if len(cfg) == 0 {
		log.Error("getActions(%d) empty", bizid)
		return
	}
	if err := json.Unmarshal([]byte(cfg), &actions); err != nil {
		log.Error("getActions(%d) error(%v)", bizid, err)
	}
	return
}

func (s *Service) getAdapter(c context.Context, bizid int64) (adps []*business.Adapter) {
	cfg := s.getConfig(c, bizid, business.TypeAdapter)
	if len(cfg) == 0 {
		return
	}

	if err := json.Unmarshal([]byte(cfg), &adps); err != nil {
		log.Error("getAdapter cfg(%s) err(%v)", cfg, err)
	}
	return
}

// Gray 灰度
func (s *Service) Gray(opt *model.AddOption) (next bool) {
	if opt == nil {
		return
	}
	//未配置，则默认全量
	if len(s.gray[opt.BusinessID]) == 0 {
		next = true
		return
	}

	optval := reflect.ValueOf(opt).Elem()
	opttp := optval.Type()
	for _, opts := range s.gray[opt.BusinessID] {
		//策略与策略之间or, 策略的fields之间and
		okcnt := 0

		for _, item := range opts {
			f, exist := opttp.FieldByName(item.Name)
			if !exist {
				break
			}

			v := optval.FieldByIndex(f.Index)
			if strings.Contains(item.Value, fmt.Sprintf(",%v,", v.Interface())) {
				okcnt++
				continue
			}
		}
		if okcnt > 0 && okcnt == len(opts) {
			next = true
			return
		}
	}

	return
}

//checkaudit 检查任务对应的资源是否需要审核，不需要的不下发
func (s *Service) checkaudit(flowid int64, state int64) bool {
	if len(s.c.Auditstate) > 0 {
		if states, ok := s.c.Auditstate[fmt.Sprint(flowid)]; ok {
			if !strings.Contains(","+states+",", fmt.Sprint(state)) {
				return false
			}
		}
	}
	return true
}

//Auth 用户权限查询
func (s *Service) Auth(c context.Context, uid int64) (a *model.Auth, err error) {
	var (
		roles []*taskmod.Role
	)
	a = &model.Auth{}
	if s.Debug() == "local" {
		a.OK = true
		return
	}
	if s.IsAdmin(uid) {
		a.OK = true
		a.Admin = true
		return
	}

	//查看业务级别&任务级别的绑定权限
	bidBiz := map[int64]int64{}
	//任务级别权限
	for biz, bidFlows := range s.taskRoleCache {
		for bid := range bidFlows {
			bidBiz[bid] = biz
		}
	}
	//业务级别权限
	for biz, cfgs := range s.bizRoleCache {
		if bid, exist := cfgs[business.BizBIDMngID]; exist && bid > 0 {
			bidBiz[bid] = biz
		}
	}
	log.Info("Auth bidbiz(%+v) \r\n   taskrole(%+v) \r\n bizrole(%+v)", bidBiz, s.taskRoleCache, s.bizRoleCache)

	//用户角色
	if roles, err = s.http.GetUserRoles(c, uid); err != nil {
		log.Error("Auth s.http.GetUserRoles(%d) error(%v)", uid, err)
		return
	}

	a.Business = map[int64]int64{}
	for _, item := range roles {
		if bidBiz[item.BID] > 0 {
			a.Business[item.BID] = bidBiz[item.BID]
		}
	}
	a.OK = len(a.Business) > 0
	return
}
