package service

import (
	"context"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/app/admin/main/workflow/model/search"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// PlatformChallCount will return count of challenges which are backlog of an admin
func (s *Service) PlatformChallCount(c context.Context, assigneeAdminID int64, permissionMap map[int8]int64) (challCount *search.ChallCount, err error) {
	var challSearchCommonResp *search.ChallSearchCommonResp

	if challCount, err = s.dao.ChallCountCache(c, assigneeAdminID); err != nil {
		log.Warn("s.dao.ChallCountCache(%d) error(%v)", assigneeAdminID, err)
		err = nil
	}
	if challCount != nil {
		return
	}

	// not fit cache, need to search es
	challCount = new(search.ChallCount)
	challCount.BusinessCount = make(map[int8]int64)
	for business, round := range permissionMap {
		cond := new(search.ChallSearchCommonCond)
		cond.Fields = []string{"id"}
		cond.Business = business
		cond.AssigneeAdminIDs = []int64{assigneeAdminID}
		cond.PN = 1
		cond.PS = 1000
		cond.Order = "id"
		cond.Sort = "desc"
		if round == model.FeedbackRound {
			cond.BusinessStates = []int64{0, 1}
		} else {
			cond.States = []int64{0}
		}
		if challSearchCommonResp, err = s.dao.SearchChallenge(c, cond); err != nil {
			log.Error("s.dao.SearchChallenge(%v) error(%v)", cond, err)
			return
		}
		challCount.BusinessCount[business] = int64(challSearchCommonResp.Page.Total)
		challCount.TotalCount += int64(challSearchCommonResp.Page.Total)

	}
	if err = s.dao.UpChallCountCache(c, challCount, assigneeAdminID); err != nil {
		log.Error("s.dao.UpChallCountCache(%d) error(%v)", assigneeAdminID, err)
		err = nil
	}
	return
}

// PlatformChallListPending will return challenges which are backlog of an admin
func (s *Service) PlatformChallListPending(c context.Context, assigneeAdminID int64, permissionMap map[int8]int64, pclp *param.ChallListParam) (challPage *search.ChallListPageCommon, err error) {
	var (
		challSearchCommonResp *search.ChallSearchCommonResp
		cids                  []int64
		uids                  []int64
		challs                map[int64]*model.Chall
		challLastLog          map[int64]string
		challLastEvent        map[int64]*model.Event
		attPaths              map[int64][]string
		uNames                map[int64]string
		attr                  *model.BusinessAttr
		rcids                 []int64
		exist                 bool
		ok                    bool
		pMeta                 map[int8]map[int64][]int64
		t                     *model.TagMeta
		l                     string
		gidToBus              map[int64]*model.Business
	)
	rand := pclp.R
	log.Info("assignee_adminid(%d) call pending rand(%d)", assigneeAdminID, rand)
	pMetas := model.PlatformMetas()
	// todo: if assign type = 0 judge if admin is online
	if exist, err = s.dao.IsOnline(c, assigneeAdminID); err != nil {
		log.Info("s.dao.IsOnline(%d) error(%v)", assigneeAdminID, err)
		return
	}
	for i, business := range pclp.Businesses {
		pMeta, ok = pMetas[business]
		if !ok {
			log.Error("not read platform meta of business(%d)", business)
		}
		if attr, ok = s.busAttrCache[business]; !ok {
			log.Error("can not find business(%d) attr", business)
			continue
		}
		assignNum := pclp.AssignNum[i]
		// assignNum not allow over assignMax
		if assignNum > attr.AssignMax {
			assignNum = attr.AssignMax
		}
		round, ok := permissionMap[business]
		if !ok {
			log.Warn("uid(%d) not has permission of business(%d) rand(%d)", assigneeAdminID, business, rand)
			continue
		}
		// need get mission from redis list (not assigneed)
		if attr.AssignType == 0 {
			// assignType == 0 need judge checkin
			if !exist {
				log.Info("uid(%d) not checkin platform rand(%d)", assigneeAdminID, rand)
				continue
			}
			// get mission from es first
			cond := &search.ChallSearchCommonCond{
				Fields:           []string{"id"},
				AssigneeAdminIDs: []int64{assigneeAdminID},
				PS:               int(assignNum),
				PN:               pclp.PN,
			}
			if round >= model.AuditRoundMin && round <= model.AuditRoundMax {
				cond.Business = business
				cond.States, ok = pMeta[0][0]
				if !ok {
					continue
				}
			}
			if round == model.FeedbackRound { //feedback flow
				cond.Business = business
				cond.BusinessStates, ok = pMeta[0][1]
				if !ok {
					continue
				}
			}
			if challSearchCommonResp, err = s.dao.SearchChallenge(c, cond); err != nil {
				log.Error("s.dao.SearchChallenge(%+v) error(%v)", cond, err)
				return
			}

			for _, r := range challSearchCommonResp.Result {
				cids = append(cids, r.ID)
			}

			log.Warn("uid(%d) has mission in db, cids:(%v)  rand(%d)", assigneeAdminID, cids, rand)
			assignNum = assignNum - int8(len(challSearchCommonResp.Result))
			if assignNum <= 0 {
				log.Warn("uid(%d) cids:(%v) business(%d) round(%d) not need consume, continue  rand(%d)", assigneeAdminID, cids, business, round, rand)
				continue
			}
			log.Warn("uid(%d) wanna consume redis business(%d) round(%d) num(%d) already has cids(%v)  rand(%d)", assigneeAdminID, business, round, assignNum, cids, rand)
			// mission from redis
			rcids, err = s.dao.RedisRPOPCids(c, business, round, assignNum)
			if err != nil {
				log.Error("s.dao.RedisRPOPCids(%d,%d,%d) error(%v)", business, round, assignNum, err)
				err = errors.WithStack(err)
				return nil, err
			}
			tx := s.dao.ORM.Begin()
			if tx.Error != nil {
				return
			}
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
					log.Error("s.PlatformChallListPending() panic(%v)", r)
				}
			}()
			// set dispatch_time if consume mission from redis
			if err = s.dao.TxUpChallAssignee(tx, rcids); err != nil {
				log.Error("s.dao.TxUpChallAssignee(%v,%d)", rcids, assigneeAdminID)
				return nil, err
			}
			log.Warn("uid(%d) consume cids(%v) business(%d) round(%d)  rand(%d)", assigneeAdminID, rcids, business, round, rand)
			cids = append(cids, rcids...)
			// set challenge business_state to pending
			if err = s.dao.TxUpChallsBusStateByIDs(tx, rcids, 1, assigneeAdminID); err != nil {
				log.Error("s.dao.TxUpChallsBusStateByIDs(%v,%d,%d)", cids, 1, assigneeAdminID)
				return nil, err
			}
			if err = tx.Commit().Error; err != nil {
				tx.Rollback()
				log.Error("Failed to tx.Commit(): %v", err)
				return
			}
		} else if attr.AssignType == 1 { // get mission only from es search (already assigneed)
			cond := &search.ChallSearchCommonCond{
				Fields:           []string{"id"},
				AssigneeAdminIDs: []int64{assigneeAdminID},
				PS:               int(assignNum),
				PN:               1,
			}
			if round >= model.AuditRoundMin && round <= model.AuditRoundMax {
				cond.Business = business
				cond.States = pMetas[business][0][0]
			} else {
				continue //feedback flow not support assign type 0
			}

			if challSearchCommonResp, err = s.dao.SearchChallenge(c, cond); err != nil {
				log.Error("s.dao.SearchChallenge(%+v) error(%v)", cond, err)
				return
			}
			for _, r := range challSearchCommonResp.Result {
				cids = append(cids, r.ID)
			}
		}
	}
	log.Info("after a pending uid(%d) rand(%d)", assigneeAdminID, rand)
	challPage = &search.ChallListPageCommon{}
	if len(cids) == 0 {
		challPage.Items = make([]*model.Chall, 0)
		challPage.Page = &model.Page{
			Num:   pclp.PN,
			Size:  pclp.PS,
			Total: 0,
		}
		return
	}

	if challs, err = s.dao.Challs(c, cids); err != nil {
		log.Error("s.dao.Challs(%v) error(%v)", cids, err)
		return
	}

	if challLastLog, err = s.LastLog(c, cids, []int{model.WLogModuleChallenge, model.WLogModuleReply}); err != nil {
		log.Error("s.batchLastLog(%v,%v) error(%v)", cids, model.WLogModuleChallenge, err)
		err = nil
	}

	if attPaths, err = s.dao.AttPathsByCids(c, cids); err != nil {
		log.Error("s.dao.AttPathsByCids() error(%v)", err)
		return
	}

	cond := &search.ChallSearchCommonCond{
		Fields: []string{"id", "gid"},
		IDs:    cids,
		PS:     1000,
		PN:     1,
	}
	if challSearchCommonResp, err = s.dao.SearchChallenge(c, cond); err != nil {
		log.Error("s.dao.SearchChallenge(%+v) error(%v)", cond, err)
		return
	}
	var gids []int64
	for _, r := range challSearchCommonResp.Result {
		gids = append(gids, r.Gid)
	}

	if gidToBus, err = s.dao.BusObjectByGids(c, gids); err != nil {
		log.Error("s.dao.BusObjectByGids(%v) error(%v)", gids, err)
		return
	}

	if challLastEvent, err = s.batchLastEvent(c, cids); err != nil {
		log.Error("s.batchLastEvent(%v) error(%v)", cids, err)
		return
	}

	for _, c := range challs {
		uids = append(uids, int64(c.AdminID))
		uids = append(uids, int64(c.AssigneeAdminID))
	}

	if uNames, err = s.dao.BatchUNameByUID(c, uids); err != nil {
		log.Error("s.dao.SearchUNameByUid(%v) error(%v)", uids, err)
		err = nil
	}

	challList := make([]*model.Chall, 0, len(challSearchCommonResp.Result))
	for _, cid := range cids {
		c, ok := challs[cid]
		if !ok {
			log.Warn("Invalid challenge id %d", cid)
			continue
		}

		// fill tag
		if t, err = s.tag(c.Business, c.Tid); err != nil {
			log.Error("s.tag(%d,%d) error(%v)", c.Business, c.Tid, err)
			err = nil
		} else {
			c.Tag = t.Name
			c.Round = t.RID
		}

		// fill last log
		if l, ok = challLastLog[cid]; ok {
			c.LastLog = l
		}

		// fill last event
		c.LastEvent = challLastEvent[cid]

		// fill attachments
		c.Attachments = make([]string, 0)
		if ps, ok := attPaths[cid]; ok {
			c.Attachments = ps
			c.FixAttachments()
		}

		//fill business object
		if b, ok := gidToBus[c.Gid]; ok {
			c.BusinessObject = b
		} else {
			log.Warn("failed to find bus object gid(%d) cid(%d)", c.Gid, c.Cid)
		}

		c.AssigneeAdminName = uNames[c.AssigneeAdminID]
		c.AdminName = uNames[c.AdminID]
		c.FromState()
		challList = append(challList, c)
	}

	challPage.Items = challList
	challPage.Page = &model.Page{
		Num:   challSearchCommonResp.Page.Num,
		Size:  challSearchCommonResp.Page.Size,
		Total: len(cids),
	}

	return
}

// PlatformChallListHandlingDone list handling challenges of admin
func (s *Service) PlatformChallListHandlingDone(c *bm.Context, pchlp *param.ChallHandlingDoneListParam, permissionMap map[int8]int64, assigneeAdminID int64, feature int8) (challPage interface{}, err error) {
	pMetas := model.PlatformMetas()
	business := pchlp.Businesses

	round := permissionMap[business]
	if _, ok := pMetas[business]; !ok { // business not in platform
		err = errors.Wrap(ecode.MethodNotAllowed, "business not in platform")
		return
	}
	if _, ok := pMetas[business][feature]; !ok { // business not has platform state
		err = errors.Wrap(ecode.MethodNotAllowed, "business not has platform state")
		return
	}
	cond := &search.ChallSearchCommonCond{
		Fields:           []string{"id", "gid"},
		Business:         business,
		AssigneeAdminIDs: []int64{assigneeAdminID},
		PS:               pchlp.PS,
		PN:               pchlp.PN,
		Sort:             pchlp.Sort,
		Order:            pchlp.Order,
	}
	if round >= model.AuditRoundMin && round <= model.AuditRoundMax { //audit flow
		cond.Business = business
		cond.States = pMetas[business][feature][0]
	}
	if round == model.FeedbackRound { //feedback flow
		cond.Business = business
		cond.BusinessStates = pMetas[business][feature][1]
	}

	return s.ChallsWrap(c, cond)
}

// PlatformChallListCreated list created challenges of admin
func (s *Service) PlatformChallListCreated(c context.Context, cond *search.ChallSearchCommonCond) (challPage *search.ChallListPageCommon, err error) {
	return s.ChallsWrap(c, cond)
}

// PlatformRelease admin offline
func (s *Service) PlatformRelease(c context.Context, permissionMap map[int8]int64, assigneeAdminID int64) (err error) {
	var (
		challSearchCommonResp *search.ChallSearchCommonResp
		cids                  []int64
		attr                  *model.BusinessAttr
		ok                    bool
	)
	cids = make([]int64, 0)

	for business, round := range permissionMap {
		cond := &search.ChallSearchCommonCond{
			Fields:           []string{"id"},
			AssigneeAdminIDs: []int64{assigneeAdminID},
			PN:               1,
			PS:               1000,
		}
		if attr, ok = s.busAttrCache[business]; !ok {
			log.Error("can not find business(%d) attr", business)
			continue
		}
		if attr.AssignType == 1 {
			continue
		} else { //任务消费 退出需要释放待处理状态的工单
			if round == model.FeedbackRound { //客服
				cond.BusinessStates = []int64{0, 1}
			} else {
				cond.States = []int64{0}
			}
			cond.Business = business
			if challSearchCommonResp, err = s.dao.SearchChallenge(c, cond); err != nil {
				log.Error("s.dao.SearchChallenge(%v) error(%v)", cond, err)
				return
			}

			for _, r := range challSearchCommonResp.Result {
				cids = append(cids, r.ID)
			}
		}
	}
	if err = s.dao.BatchResetAssigneeAdminID(cids); err != nil {
		return
	}
	err = s.dao.DelOnline(c, assigneeAdminID)
	// add report
	log.Info("uid(%d) offline success err(%v)", assigneeAdminID, err)
	return
}

// PlatformCheckIn admin online
func (s *Service) PlatformCheckIn(c context.Context, assigneeAdminID int64) (err error) {
	err = s.dao.AddOnline(c, assigneeAdminID)
	// add report
	log.Info("uid(%d) online success err(%v)", assigneeAdminID, err)
	return
}

// PlatformOnlineList .
func (s *Service) PlatformOnlineList(c context.Context) (err error) {
	var onlineAdminIDs []int64
	if onlineAdminIDs, err = s.dao.ListOnline(c); err != nil {
		return
	}

	// search login/out time, last 24h operate
	s.dao.LogInOutTime(c, onlineAdminIDs)
	return
}

// ChallsWrap warp challenges list result
func (s *Service) ChallsWrap(c context.Context, cond *search.ChallSearchCommonCond) (challPageCommon *search.ChallListPageCommon, err error) {
	var (
		challSearchCommonResp *search.ChallSearchCommonResp
		challLastLog          map[int64]string
		challLastEvent        map[int64]*model.Event
		attPaths              map[int64][]string
		gidToBus              map[int64]*model.Business
		uNames                map[int64]string
		challs                map[int64]*model.Chall
		cids                  []int64
		uids                  []int64
		gids                  []int64
		t                     *model.TagMeta
		l                     string
	)
	challSearchCommonResp, err = s.dao.SearchChallenge(c, cond)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	cids = make([]int64, 0)
	uids = make([]int64, 0, len(challSearchCommonResp.Result)*2)
	gids = make([]int64, 0)
	for _, r := range challSearchCommonResp.Result {
		cids = append(cids, r.ID)
		gids = append(gids, r.Gid)
	}
	challPageCommon = new(search.ChallListPageCommon)
	if len(cids) == 0 {
		challPageCommon.Items = make([]*model.Chall, 0)
		challPageCommon.Page = &model.Page{
			Num:   cond.PN,
			Size:  cond.PS,
			Total: 0,
		}
		return
	}
	if challs, err = s.dao.Challs(c, cids); err != nil {
		log.Error("s.dao.Challs(%v) error(%v)", cids, err)
		return
	}

	if challLastLog, err = s.LastLog(c, cids, []int{model.WLogModuleChallenge, model.WLogModuleReply}); err != nil {
		log.Error("s.batchLastLog(%v,%v) error(%v)", cids, model.WLogModuleChallenge, err)
		err = nil
	}

	if attPaths, err = s.dao.AttPathsByCids(c, cids); err != nil {
		log.Error("s.dao.AttPathsByCids() error(%v)", err)
		return
	}

	if gidToBus, err = s.dao.BusObjectByGids(c, gids); err != nil {
		log.Error("s.dao.BusObjectByGids(%v) error(%v)", gids, err)
		return
	}

	if challLastEvent, err = s.batchLastEvent(c, cids); err != nil {
		log.Error("s.batchLastEvent(%v) error(%v)", cids, err)
		return
	}

	for _, c := range challs {
		uids = append(uids, int64(c.AdminID))
		uids = append(uids, int64(c.AssigneeAdminID))
	}

	if uNames, err = s.dao.BatchUNameByUID(c, uids); err != nil {
		log.Error("s.dao.SearchUNameByUid(%v) error(%v)", uids, err)
		err = nil
	}

	challList := make([]*model.Chall, 0, len(cids))
	for _, cid := range cids {
		c, ok := challs[cid]
		if !ok {
			log.Warn("Invalid challenge id %d", cid)
			continue
		}

		// fill tag
		if t, err = s.tag(c.Business, c.Tid); err != nil {
			log.Error("s.tag(%d,%d) error(%v)", c.Business, c.Tid, err)
			err = nil
		} else {
			c.Tag = t.Name
			c.Round = t.RID
		}

		// fill last log
		if l, ok = challLastLog[cid]; ok {
			c.LastLog = l
		}

		// fill last event
		c.LastEvent = challLastEvent[cid]

		// fill attachments
		c.Attachments = make([]string, 0)
		if ps, ok := attPaths[cid]; ok {
			c.Attachments = ps
			c.FixAttachments()
		}

		//fill business object
		if b, ok := gidToBus[c.Gid]; ok {
			c.BusinessObject = b
		} else {
			log.Warn("failed to find bus object gid(%d) cid(%d)", c.Gid, c.Cid)
		}

		c.AssigneeAdminName = uNames[c.AssigneeAdminID]
		c.AdminName = uNames[c.AdminID]
		c.FromState()

		challList = append(challList, c)
	}

	challPageCommon.Items = challList
	challPageCommon.Page = &model.Page{
		Num:   challSearchCommonResp.Page.Num,
		Size:  challSearchCommonResp.Page.Size,
		Total: challSearchCommonResp.Page.Total,
	}

	return
}
