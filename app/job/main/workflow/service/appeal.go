package service

import (
	"context"
	"time"

	"go-common/app/job/main/workflow/model"
	"go-common/library/log"

	"go-common/library/sync/errgroup.v2"
)

// 工作台单条过期
func (s *Service) singleExpireproc() {
	errGroup := errgroup.Group{}
	for {
		for _, attr := range s.businessAttr {
			if attr.DealType != model.PDealType {
				continue
			}
			time.Sleep(1 * time.Second)
			errGroup.Go(func(ctx context.Context) error {
				return s.singleExpire(ctx, attr.Bid)
			})
		}
		errGroup.Wait()
	}
}

func (s *Service) singleExpire(c context.Context, bid int) (err error) {
	var delIDs []int64
	if delIDs, err = s.dao.SingleExpire(c, bid); err != nil {
		log.Error("s.dao.SingleExpire() bid:%d error: %v", bid, err)
		return
	}
	if len(delIDs) == 0 {
		return
	}
	log.Info("bid(%d) apid in (%v) are single expire", bid, delIDs)
	// set db state
	if err = s.dao.SetAppealAssignState(c, delIDs, model.AssignStateNotDispatch); err != nil {
		log.Error("s.dao.SetAppealAssignState() error: %v", err)
		return
	}
	if err = s.dao.DelSingleExpire(c, bid, delIDs); err != nil {
		log.Error("s.dao.SetSingleExpire() bid:%d error: %v", bid, err)
		return
	}
	if err = s.delRelatedMissions(c, bid, delIDs); err != nil {
		log.Error("s.delRelatedMissions() error: %v", err)
	}
	return
}

// 反馈整体过期
func (s *Service) overallExpireproc() {
	errGroup := errgroup.Group{}
	for {
		for _, attr := range s.businessAttr {
			if attr.DealType != model.PDealType {
				continue
			}
			time.Sleep(1 * time.Second)
			errGroup.Go(func(ctx context.Context) error {
				return s.overallExpire(ctx, attr.Bid)
			})
		}
		errGroup.Wait()
	}
}

func (s *Service) overallExpire(c context.Context, bid int) (err error) {
	time.Sleep(10 * time.Second)
	cond := model.AppealSearchCond{
		Fields:        []string{"id"},
		Bid:           []int{bid},
		AssignState:   []int8{model.AssignStatePoped},
		TransferState: []int8{model.TransferStatePendingSystemReply},
		DTimeTo:       time.Now().Add(-time.Minute * 5).Format("2006-01-02 15:04:05"),
		PS:            101,
		PN:            1,
		Order:         "id",
		Sort:          "desc",
	}
	var res *model.AppealSearchRes
	if res, err = s.dao.SearchAppeal(c, cond); err != nil {
		log.Error("s.dao.SearchAppeal() error(%+v)", err)
		return
	}
	if len(res.Result) == 0 {
		log.Warn("no appeal is expire overall")
		return
	}
	ids := make([]int64, 0, len(res.Result))
	for _, r := range res.Result {
		ids = append(ids, r.ID)
	}
	log.Info("apids(%v) are expire overall", ids)
	if err = s.delRelatedMissions(c, bid, ids); err != nil {
		log.Error("s.delRelatedMissions() error: %v", err)
	}
	if err = s.dao.SetAppealAssignState(c, ids, model.AssignStateNotDispatch); err != nil {
		log.Error("s.dao.SetAppealAssignState(%v,%d) error(%v)", ids, model.AssignStateNotDispatch, err)
	}
	return
}

// 释放用户未评价反馈
func (s *Service) releaseExpireproc() {
	for {
		for _, attr := range s.businessAttr {
			var (
				c   = context.TODO()
				err error
			)

			if attr.DealType != model.PDealType {
				continue
			}
			time.Sleep(1 * time.Minute)
			cond := model.AppealSearchCond{
				Fields:        []string{"id", "mid"},
				Bid:           []int{attr.Bid},
				TTimeTo:       time.Now().AddDate(0, 0, -3).Format("2006-01-02 15:04:05"),
				TransferState: []int8{model.TransferStatePendingSystemReply, model.TransferStateAdminReplyReaded, model.TransferStateAdminReplyNotReaded},
				PS:            50,
				PN:            1,
				Order:         "id",
				Sort:          "desc",
			}
			var res *model.AppealSearchRes
			if res, err = s.dao.SearchAppeal(c, cond); err != nil {
				log.Error("d.SearchAppeal error(%+v)", cond)
				continue
			}
			ids := make([]int64, 0, len(res.Result))
			for _, r := range res.Result {
				var e *model.Event
				if e, err = s.dao.LastEvent(r.ID); err != nil {
					log.Error("s.dao.LastEvent() id(%d) error:%v", r.ID, err)
					continue
				}
				// 关闭最后一个对话为管理员回复的申诉
				if e.Event == model.EventAdminReply {
					ids = append(ids, r.ID)
				}
			}
			if len(ids) == 0 {
				continue
			}

			// 删除权重参数
			if err = s.dao.DelUperInfo(c, ids); err != nil {
				log.Error("s.dao.DelUperInfo(%v), error(%v)", ids, err)
				continue
			}

			// todo delete single expire zset member
			if err = s.dao.DelSingleExpire(c, attr.Bid, ids); err != nil {
				log.Error("s.dao.DelSingleExpire(%d, %v), error(%v)", attr.Bid, ids, err)
				continue
			}

			if err = s.dao.SetAppealTransferState(c, ids, model.TransferStateAutoClosedExpire); err != nil {
				log.Error("s.dao.SetAppealTransferState() ids(%v) transferstate(%d) error(%v)", ids, model.TransferStateAutoClosedExpire, err)
				continue
			}
			log.Info("apids (%v) are expire user not set degree", ids)
		}
	}
}

// 进任务池
func (s *Service) enterPoolproc() {
	for {
		errGroup := errgroup.Group{}
		for _, attr := range s.businessAttr {
			if attr.DealType != model.PDealType {
				continue
			}
			time.Sleep(1 * time.Second)
			cond := model.AppealSearchCond{
				Fields:        []string{"id"},
				Bid:           []int{attr.Bid},
				AssignState:   []int8{model.AssignStateNotDispatch},
				AuditState:    []int8{model.AuditStateInvalid},
				TransferState: []int8{model.TransferStatePendingSystemReply, model.TransferStatePendingSystemNotReply},
				PS:            99,
				PN:            1,
				Order:         "id",
				Sort:          "desc",
			}
			errGroup.Go(func(ctx context.Context) error {
				return s.enterPool(ctx, cond)
			})
		}
		errGroup.Wait()
	}
}

func (s *Service) enterPool(c context.Context, cond model.AppealSearchCond) (err error) {
	var res *model.AppealSearchRes
	if res, err = s.dao.SearchAppeal(c, cond); err != nil {
		log.Error("s.dao.SearchAppeal error(%+v)", cond)
		return
	}
	ids := make([]int64, len(res.Result))
	for _, r := range res.Result {
		ids = append(ids, r.ID)
	}
	var appeals []*model.Appeal
	if appeals, err = s.dao.Appeals(c, ids); err != nil {
		log.Error("s.dao.Appeals(%v) error(%v)", ids, err)
		return
	}

	ApIDMap := make(map[int64]int64) //map[ap_id]weight
	ApIDs := make([]int64, 0)
	// check state
	for _, ap := range appeals {
		if ap.AssignState == model.AssignStateNotDispatch && ap.AuditState == model.AuditStateInvalid && ap.TransferState == model.TransferStatePendingSystemReply {
			ApIDMap[ap.ApID] = ap.Weight
			ApIDs = append(ApIDs, ap.ApID)
		}
	}
	if len(ApIDs) == 0 {
		log.Warn("bid(%v) not found apids after check db should enter pool!", cond.Bid)
		return
	}

	log.Info("bid(%v) apids(%v) ApIDMap(%v) should set into mission pool", cond.Bid, ApIDs, ApIDMap)
	if err = s.dao.SetAppealAssignState(c, ApIDs, model.AssignStatePushed); err != nil {
		log.Error("s.dao.SetAppealAssignState(%v) err(%v)", ApIDs, err)
		return
	}
	//set sorted set
	if err = s.dao.SetWeightSortedSet(c, cond.Bid[0], ApIDMap); err != nil {
		log.Error("s.dao.SetWeightSortedSet() error(%v)", err)
	}
	return
}

// 刷新权重值 每3分钟刷新一次
func (s *Service) refreshWeightproc() {
	errGroup := errgroup.Group{}
	for {
		time.Sleep(3 * time.Second) // todo 3 min
		for _, attr := range s.businessAttr {
			if attr.DealType != model.PDealType {
				continue
			}
			time.Sleep(1 * time.Second)
			cond := model.AppealSearchCond{
				Bid:           []int{attr.Bid},
				Fields:        []string{"id", "mid", "weight"},
				AuditState:    []int8{model.AuditStateInvalid},
				TransferState: []int8{model.TransferStatePendingSystemReply},
				PS:            100,
				PN:            1,
				Order:         "id",
				Sort:          "desc",
			}
			errGroup.Go(func(ctx context.Context) error {
				time.Sleep(3 * time.Second)
				return s.refreshWeight(ctx, cond)
			})
		}
		errGroup.Wait()
	}
}

func (s *Service) refreshWeight(c context.Context, cond model.AppealSearchCond) (err error) {
	var res *model.AppealSearchRes
	if res, err = s.dao.SearchAppeal(c, cond); err != nil {
		log.Error("d.SearchAppeal error(%+v)", cond)
		return
	}
	appeals := make([]*model.Appeal, 0, len(res.Result))
	apIDs := make([]int64, 0, len(res.Result))
	newWeight := make(map[int64]int64, len(res.Result))
	newWeightInSortedSet := make(map[int64]int64, len(res.Result))
	for _, ap := range res.Result {
		appeals = append(appeals, &model.Appeal{
			ApID:   ap.ID,
			Mid:    ap.Mid,
			Weight: ap.Weight,
		})
		apIDs = append(apIDs, ap.ID)
	}
	if len(appeals) == 0 {
		log.Warn("no appeal is in feedback")
		return
	}

	var params []int64
	if params, err = s.dao.UperInfoCache(c, apIDs); err != nil { // 读用户维度的权重参数
		log.Error("s.dao.UperInfoCache(%v) error(%v)", apIDs, err)
		return
	}
	// fixme cache miss
	if len(params) != len(apIDs) {
		log.Warn("len params not equre len mids")
		return
	}
	for i, ap := range appeals {
		p := params[i]
		incr := s.calcWeight(ap, p)
		newWeight[ap.ApID] = ap.Weight + incr
		if ap.AssignState == model.AssignStatePushed { // should rewrite weight in db
			newWeightInSortedSet[ap.ApID] = newWeight[ap.ApID]
		}
	}
	log.Info("appeals in feedback weight(%+v) mids(%v)", newWeight, apIDs)
	tx := s.dao.WriteORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.WriteORM.Begin() error(%v)", err)
		return
	}
	if err = s.dao.TxSetWeight(tx, newWeight); err != nil {
		log.Error("s.dao.SetWeight(%v) error(%v)", newWeight, err)
		tx.Rollback()
		return
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("tx.Commit() error:%v", err)
		tx.Rollback()
		return
	}
	// async write sorted set
	if len(newWeightInSortedSet) == 0 {
		return
	}
	_ = s.cache.Do(c, func(c context.Context) {
		if err = s.dao.SetWeightSortedSet(c, cond.Bid[0], newWeightInSortedSet); err != nil {
			log.Error("s.dao.SetWeightSortedSet() error(%v)", err)
		}
		log.Info("appeals in feedback sorted set bid(%d) weight(%+v)", cond.Bid[0], newWeightInSortedSet)
	})
	return
}

func (s *Service) delRelatedMissions(c context.Context, bid int, delIDs []int64) (err error) {
	cond := model.AppealSearchCond{
		Fields: []string{"id", "transfer_admin"},
		IDs:    delIDs,
		Order:  "id",
		Sort:   "asc",
		PS:     100,
		PN:     1,
	}

	var res *model.AppealSearchRes
	if res, err = s.dao.SearchAppeal(c, cond); err != nil {
		log.Error("s.dao.SearchAppeal() error:%v")
		return
	}
	tadmin := make(map[int][]int64)
	for _, r := range res.Result {
		if _, ok := tadmin[r.TransferAdmin]; !ok {
			tadmin[r.TransferAdmin] = make([]int64, 0)
		}
		tadmin[r.TransferAdmin] = append(tadmin[r.TransferAdmin], r.ID)
	}
	for uid, ids := range tadmin {
		if err = s.dao.DelRelatedMissions(c, bid, uid, ids); err != nil {
			log.Error("s.dao.DelRelatedMissions() bid(%d) uid(%d) ap_ids(%v) error:%v", bid, uid, ids, err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

// calcWeight 计算权重值
func (s *Service) calcWeight(ap *model.Appeal, p int64) (incr int64) {
	switch p { // 粉丝/特殊用户组维度
	case 1:
		incr += 8
	case 2:
		incr += 10
	case 3:
		incr += 12
	case 4:
		incr += 15
	case 5:
		incr += 18
	default:
		incr += 8
	}
	switch { // 时间维度
	case time.Since(ap.CTime.Time()) < 3*time.Minute:
		incr += 0
	case time.Since(ap.CTime.Time()) >= 3*time.Minute && time.Since(ap.CTime.Time()) < 6*time.Minute:
		incr += 3
	case time.Since(ap.CTime.Time()) >= 6*time.Minute && time.Since(ap.CTime.Time()) < 9*time.Minute:
		incr += 6
	case time.Since(ap.CTime.Time()) >= 9*time.Minute && time.Since(ap.CTime.Time()) < 15*time.Minute:
		incr += 9
	case time.Since(ap.CTime.Time()) >= 15*time.Minute:
		incr += 12
	default:
		incr += 0
	}
	return incr
}
