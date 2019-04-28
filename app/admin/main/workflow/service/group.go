package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"sort"
	"strconv"
	"time"

	"go-common/app/admin/main/credit/model/blocked"
	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/app/admin/main/workflow/model/search"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
	"golang.org/x/sync/errgroup"
)

// GroupListV3 .
func (s *Service) GroupListV3(c context.Context, cond *search.GroupSearchCommonCond) (grpPage *model.GroupListPage, err error) {
	var (
		ok                    bool
		groupSearchCommonResp *search.GroupSearchCommonResp
		wg                    = &errgroup.Group{}
	)

	if groupSearchCommonResp, err = s.dao.SearchGroup(c, cond); err != nil {
		log.Error("Failed to s.dao.SearchGroup(%v): %v", cond, err)
		err = ecode.WkfSearchGroupFailed
		return
	}
	if len(groupSearchCommonResp.Result) == 0 {
		grpPage = &model.GroupListPage{
			Items: []*model.Group{},
			Page:  &model.Page{},
		}
		return
	}
	eids, oids, gids, mids := []int64{}, []int64{}, []int64{}, []int64{}
	fut := make(map[int64]int64, len(groupSearchCommonResp.Result))
	for _, v := range groupSearchCommonResp.Result {
		gids = append(gids, v.ID)
		oids = append(oids, v.Oid)
		eids = append(eids, v.Eid)
		if v.Mid > 0 {
			mids = append(mids, v.Mid)
		}
		fut[v.ID] = v.FirstUserTid
	}
	// last chall mids
	cscc := &search.ChallSearchCommonCond{
		Fields:   []string{"id", "gid", "mid"},
		Business: cond.Business,
		Gids:     gids,
		Order:    "ctime",
		Sort:     "desc",
		Distinct: []string{"gid"},
		PS:       999,
		PN:       1,
	}
	var (
		cResp      *search.ChallSearchCommonResp
		lastChalls = make(map[int64]*model.Chall)
	)

	if cResp, err = s.dao.SearchChallenge(c, cscc); err != nil {
		log.Error("s.dao.SearchChallenge(%v) error(%v)", cscc, err)
	} else {
		for _, r := range cResp.Result {
			if r.Mid > 0 {
				mids = append(mids, r.Mid)
			}
			lastChalls[r.Gid] = &model.Chall{Cid: r.ID, Mid: r.Mid}
		}
	}

	// group object
	var groups map[int64]*model.Group
	wg.Go(func() error {
		if groups, err = s.dao.Groups(c, gids); err != nil {
			log.Error("Failed to s.dao.Groups(%v): %v", gids, err)
			return err
		}
		return nil
	})

	// todo load data from db, search, external api
	// tag count of challenge
	var gidToChallTagCount map[int64]map[int64]int64
	wg.Go(func() error {
		// todo: count tag in es
		if gidToChallTagCount, err = s.dao.ChallTagsCountV3(c, gids); err != nil {
			log.Error("s.dao.ChallTagsCountV3(%v) error(%v)", gids, err)
			err = nil
		} else {
			log.Info("gidToChallTagCount (%+v)", gidToChallTagCount)
		}
		return nil
	})

	var grpLastLog map[int64]string
	wg.Go(func() error {
		if grpLastLog, err = s.LastLog(c, gids, []int{model.WLogModuleGroup, model.WLogModuleRoleShift}); err != nil {
			log.Error("s.LastLog(%v,%d) error(%v)", gids, model.WLogModuleGroup, err)
			err = nil
		}
		return nil
	})

	// search account
	var users map[int64]*model.Account
	wg.Go(func() error {
		users = s.dao.AccountInfoRPC(c, mids)
		//get uper group
		var uperTagMap map[int64][]*model.SpecialTag
		if uperTagMap, err = s.dao.BatchUperSpecial(c, mids); err != nil {
			log.Error("s.dao.BatchUperSpecial(%v) error(%v)", mids, err)
			err = nil
		} else {
			for id, user := range users {
				var st []*model.SpecialTag
				if st, ok = uperTagMap[id]; !ok {
					log.Warn("not find special tag mid(%d)", id)
					continue
				}
				user.SpecialTag = st
			}
		}
		return nil
	})

	// object table
	// todo: judge if need read local object
	var bus map[int64]*model.Business
	var archives map[int64]*model.Archive
	wg.Go(func() error {
		if bus, err = s.dao.BusObjectByGids(c, gids); err != nil {
			log.Error("s.dao.BusObjectByGids(%v) error(%v)", gids, err)
			err = nil
		}

		// search archive
		// todo: judge if need search archive
		aids := []int64{}
		for _, b := range bus {
			aids = append(aids, b.Oid)
		}
		if archives, err = s.dao.ArchiveRPC(c, aids); err != nil {
			log.Error("s.dao.ArchiveRPC(%v) error(%v)", oids, err)
			err = nil
		}
		return nil
	})

	// external meta
	var metas map[int64]*model.GroupMeta
	wg.Go(func() error {
		if metas, err = s.externalMeta(c, cond.Business, gids, oids, eids); err != nil {
			log.Error("s.ExternalMeta(%d,%v,%v,%v) error(%v)", cond.Business, gids, oids, eids, err)
			err = nil
		} else {
			log.Info("external meta (%+v)", metas)
		}
		return nil
	})

	// wait all wg.go()
	if err = wg.Wait(); err != nil {
		return
	}

	//todo make response
	grpPage = new(model.GroupListPage)
	rgs := make([]*model.Group, 0, len(groupSearchCommonResp.Result))
	for _, v := range groupSearchCommonResp.Result {
		var (
			rg *model.Group
			ok bool
		)
		if rg, ok = groups[v.ID]; !ok {
			log.Warn("Failed to retrive group by group id %d", v.ID)
			continue
		}
		if !dbCheck(cond, rg) {
			continue
		}
		rg.TypeID = v.TypeID
		// fill last log
		var l string
		if l, ok = grpLastLog[v.ID]; ok {
			rg.LastLog = l
		}
		// fill tag name
		// using tid from group row can ensure the lastest tid fetched
		if tid := rg.Tid; tid != 0 {
			var t *model.TagMeta
			if t, err = s.tag(rg.Business, tid); err != nil {
				log.Error("Failed to retrive tag by bid(%d) tag_id(%d)", rg.Business, tid)
				err = nil
			} else {
				rg.Tag = t.Name
			}
		}
		rg.ChallengeTags = model.ChallTagSlice{}
		var tc map[int64]int64
		if tc, ok = gidToChallTagCount[v.ID]; ok {
			total := int64(0)
			for _, count := range tc {
				total += count
			}
			for tid, c := range tc {
				tname := "<Unknow>"
				tround := int8(0)
				var t *model.TagMeta
				if t, err = s.tag(rg.Business, tid); err != nil {
					log.Error("Failed to retrive tag by bid(%d) tag_id(%d)", rg.Business, tid)
					err = nil
				} else {
					tname = t.Name
					tround = t.RID
				}
				ct := &model.ChallTag{
					ID:      tid,
					Tag:     tname,
					Count:   c,
					Percent: 100,
					Round:   tround,
				}
				if total != 0 {
					ct.Percent = round((float64(c) / float64(total)) * 100)
				}
				rg.ChallengeTags = append(rg.ChallengeTags, ct)
			}
			log.Warn("challenge tags of gid(%d) %+v", rg.ID, rg.ChallengeTags)
		} else {
			log.Warn("not found chall tag count of gid(%d)", v.ID)
		}
		sort.Sort(rg.ChallengeTags)
		// fill last producer
		var lc *model.Chall
		if lc, ok = lastChalls[rg.ID]; ok {
			if rg.LastProducer, ok = users[lc.Mid]; !ok {
				log.Warn("gid(%d) has last producer mid(%d) but not found account", rg.ID, lc.Mid)
			}
			log.Info("gid(%d) load last producer mid(%d) success (%+v)", rg.ID, lc.Mid, rg.LastProducer)
		} else {
			log.Warn("not found account of last producer gid(%d)", rg.ID)
		}
		// fill meta
		s.warpMeta(cond.Business, rg, metas, archives, users, bus)
		// oid to string
		rg.OidStr = strconv.FormatInt(rg.Oid, 10)
		// eid to string
		rg.EidStr = strconv.FormatInt(rg.Eid, 10)
		rgs = append(rgs, rg)
		// fill first_user_tid
		rg.FirstUserTid = fut[rg.ID]
	}
	grpPage.Items = rgs
	grpPage.Page = &model.Page{
		Num:   groupSearchCommonResp.Page.Num,
		Size:  groupSearchCommonResp.Page.Size,
		Total: groupSearchCommonResp.Page.Total,
	}
	return
}

// UpGroup will update a group
func (s *Service) UpGroup(c context.Context, gp *param.GroupParam) (err error) {
	var (
		tx *gorm.DB
		l  *model.WLog
		g  *model.Group
	)

	// double write rid
	tMeta := &model.TagMeta{}
	if tMeta, err = s.tag(gp.Business, gp.Tid); err != nil {
		log.Error("TagListCache not found bid(%d) tag_id(%d)", gp.Business, gp.Tid)
		return
	}
	gp.Rid = tMeta.RID

	// Check group and tag is exist
	if g, err = s.dao.GroupByOid(c, gp.Oid, gp.Business); err != nil {
		log.Error("s.dao.GroupByOid(%d, %d) error(%v)", gp.Oid, gp.Business, err)
		return
	}
	if g == nil {
		log.Error("Group(%d, %d) not exist", gp.Oid, gp.Business)
		err = ecode.WkfGroupNotFound
		return
	}

	tx = s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("Service.UpGroup() panic(%v)", r)
		}
	}()

	if err = s.dao.TxUpGroup(tx, gp.Oid, gp.Business, gp.Tid, gp.Note, gp.Rid); err != nil {
		tx.Rollback()
		log.Error("s.TxUpGroup(%d, %d, %d, %s, %d) error(%v)", gp.Oid, gp.Business, gp.Tid, gp.Note, gp.Rid, err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("tx.Commit() error(%v)", err)
		return
	}

	s.task(func() {
		l = &model.WLog{
			AdminID:  gp.AdminID,
			Admin:    gp.AdminName,
			Oid:      g.Oid,
			Business: g.Business,
			Target:   g.ID,
			Module:   model.WLogModuleGroup,
			Remark:   fmt.Sprintf(`工单编号 %d “管理 Tag”更新为“%s”`, g.ID, tMeta.Name),
			Note:     gp.Note,
		}
		s.writeAuditLog(l)
	})
	return
}

// UpGroupRole will 流转工单
func (s *Service) UpGroupRole(c context.Context, grsp *param.GroupRoleSetParam) (err error) {
	var groups map[int64]*model.Group
	// Check group and tag is exist
	if groups, err = s.dao.Groups(c, grsp.GID); err != nil {
		log.Error("s.dao.Groups(%v) error(%v)", grsp.GID, err)
		return
	}
	if len(groups) == 0 {
		log.Error("Group(%v) not exist", grsp.GID)
		err = ecode.WkfGroupNotFound
		return
	}
	// check bid
	for _, g := range groups {
		if g.Business != grsp.BID {
			err = ecode.WkfBusinessNotConsistent
			return
		}
	}
	// check tid available
	var tMeta *model.TagMeta
	if tMeta, err = s.tag(grsp.BID, grsp.TID); err != nil {
		return
	}
	grsp.RID = tMeta.RID

	if err = s.dao.UpGroupRole(c, grsp); err != nil {
		log.Error("s.UpGroupRole(%+v) error(%v)", grsp, err)
		return
	}
	s.task(func() {
		s.afterSetGroupRole(grsp, groups)
	})
	return
}

// SetGroupResult will set a group result
func (s *Service) SetGroupResult(c context.Context, grp *param.GroupResParam) (err error) {
	var (
		tx         *gorm.DB
		g          *model.Group
		tinyChalls map[int64]*model.TinyChall
		cids       []int64
	)

	if g, err = s.dao.GroupByOid(c, grp.Oid, grp.Business); err != nil {
		log.Error("s.dao.GroupByOid() error(%v)", err)
		return
	}
	if g == nil {
		log.Error("Group(%d, %d) not exist", grp.Oid, grp.Business)
		err = ecode.NothingFound
		return
	}
	if g.State != model.Pending {
		log.Error("Group(%d, %d) not pending", grp.Oid, grp.Business)
		return
	}

	tx = s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("Service.SetGroupResult() panic(%v)", r)
		}
	}()

	// Update Group State
	if err = s.dao.TxUpGroupState(tx, g.ID, grp.State); err != nil {
		tx.Rollback()
		log.Error("s.txUpGroupState(%+v) error(%v)", g, err)
		return
	}

	// set grouo handling field to 0
	if err = s.dao.TxUpGroupHandling(tx, g.ID, 0); err != nil {
		tx.Rollback()
		log.Error("s.txUpGroupHandling(%+v, 0) error(%v)", g.ID, err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("tx.Commit() error(%v)", err)
		return
	}

	s.task(func() {
		var result []*search.ChallSearchCommonData
		cond := &search.ChallSearchCommonCond{
			Fields: []string{"id", "gid", "mid", "state", "ctime"},
			Gids:   []int64{g.ID},
			States: []int64{int64(model.Pending)},
		}
		if result, err = s.dao.SearchChallengeMultiPage(context.Background(), cond); err != nil {
			log.Error("s.dao.SearchChallengeMultiPage(%+v) error(%v)", cond, err)
			return
		}
		tinyChalls = make(map[int64]*model.TinyChall, len(cids))
		for _, c := range result {
			cids = append(cids, c.ID)
			tc := &model.TinyChall{
				Cid: c.ID,
				Gid: c.Gid,
				Mid: c.Mid,
			}
			var ctime time.Time
			ctime, err = time.ParseInLocation("2006-01-02 15:04:05", c.CTime, time.Local)
			if err != nil {
				log.Error("time.Parse(%v) error(%v)", c.CTime, err)
			}
			tc.CTime.Scan(ctime)
			log.Info("tc.CTime.Scan(%v) get (%v) cid(%v) gid(%v)", ctime.Unix(), tc.CTime.Time(), tc.Cid, tc.Gid)
			if str, ok := c.State.(string); ok {
				st, _ := strconv.Atoi(str)
				tc.State = int8(st)
			}
			if f, ok := c.State.(float64); ok {
				tc.State = int8(math.Floor(f))
			}
			tinyChalls[c.ID] = tc
		}

		if err = s.dao.BatchUpChallByIDs(cids, uint32(grp.State), grp.AdminID); err != nil {
			log.Error("s.dao.TxBatchUpChallByIDs(%v,%d) error(%v)", cids, grp.State, err)
			return
		}
		s.afterSetGrpResult(grp, g, tinyChalls)
	})
	return
}

// BatchSetGroupResult will set a set of groups result
func (s *Service) BatchSetGroupResult(c context.Context, bgrp *param.BatchGroupResParam) (err error) {
	var (
		tx         *gorm.DB
		groups     map[int64]*model.Group
		tinyChalls map[int64]*model.TinyChall
		gids       []int64
		cids       []int64
	)

	tx = s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("Service.BatchSetGroupResult() panic(%v)", r)
		}
	}()

	if groups, err = s.dao.TxGroupsByOidsStates(tx, bgrp.Oids, bgrp.Business, model.Pending); err != nil {
		log.Error("s.dao.TxGroupsByOidsStates() error(%v)", err)
		return
	}
	if len(groups) <= 0 {
		log.Warn("No pending groups found with conditon(%+v, %d, %d)", bgrp.Oids, bgrp.Business, model.Pending)
		return
	}
	// collect all gids
	for gid := range groups {
		gids = append(gids, int64(gid))
	}

	if err = s.dao.TxBatchUpGroupState(tx, gids, bgrp.State); err != nil {
		tx.Rollback()
		log.Error("s.TxBatchUpGroupState(%+v) error(%v)", bgrp, err)
		return
	}

	// Set group handling count to 0
	if err = s.dao.TxBatchUpGroupHandling(tx, gids, 0); err != nil {
		tx.Rollback()
		log.Error("s.TxBatchUpGroupHandling(%+v, 0) error(%v)", gids, err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("tx.Commit() error(%v)", err)
		return
	}

	s.task(func() {
		var result []*search.ChallSearchCommonData
		cond := &search.ChallSearchCommonCond{
			Fields: []string{"id", "gid", "mid", "state", "ctime"},
			Gids:   gids,
			States: []int64{int64(model.Pending)},
		}
		if result, err = s.dao.SearchChallengeMultiPage(context.Background(), cond); err != nil {
			log.Error("s.dao.SearchChallengeMultiPage(%+v) error(%v)", cond, err)
			return
		}
		tinyChalls = make(map[int64]*model.TinyChall, len(cids))
		for _, c := range result {
			cids = append(cids, c.ID)
			tc := &model.TinyChall{
				Cid: c.ID,
				Gid: c.Gid,
				Mid: c.Mid,
			}
			var ctime time.Time
			ctime, err = time.ParseInLocation("2006-01-02 15:04:05", c.CTime, time.Local)
			if err != nil {
				log.Error("time.Parse(%v) error(%v)", c.CTime, err)
			}
			tc.CTime.Scan(ctime)
			log.Info("tc.CTime.Scan(%v) get (%v) cid(%v) gid(%v)", ctime.Unix(), tc.CTime.Time(), tc.Cid, tc.Gid)
			if str, ok := c.State.(string); ok {
				st, _ := strconv.Atoi(str)
				tc.State = int8(st)
			}
			if f, ok := c.State.(float64); ok {
				tc.State = int8(math.Floor(f))
			}
			tinyChalls[c.ID] = tc
		}

		if err = s.dao.BatchUpChallByIDs(cids, uint32(bgrp.State), bgrp.AdminID); err != nil {
			log.Error("s.dao.BatchUpChallByIDs(%v,%d) error(%v)", cids, bgrp.State, err)
			return
		}
		s.afterBatchSetGrpResult(bgrp, groups, tinyChalls)
	})
	return
}

// SetGroupState 修改未处理的工单状态
func (s *Service) SetGroupState(c context.Context, gssp *param.GroupStateSetParam) (err error) {
	var (
		groups     map[int64]*model.Group
		tinyChalls map[int64]*model.TinyChall
		gids, cids []int64
		newRid     int8
	)
	//check tid 有效处理验证tid
	if gssp.State == model.Effective {
		var tmeta *model.TagMeta
		if tmeta, err = s.tag(gssp.Business, gssp.Tid); err != nil {
			return
		}
		newRid = tmeta.RID
	}

	tx := s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("s.SetGroupState() panic(%v)", r)
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	if groups, err = s.dao.Groups(c, gssp.ID); err != nil {
		log.Error("s.dao.TxGroups(%v) error(%v)", gssp.ID, err)
		return
	}

	// ignore group if state not pending
	for id, g := range groups {
		if g.State != model.Pending || g.Business != gssp.Business || g.Rid != gssp.Rid {
			delete(groups, id)
			continue
		}
		gids = append(gids, id)
	}

	if len(gids) <= 0 {
		log.Warn("No settable groups found with conditon(%+v)", *gssp)
		err = ecode.WkfGroupNotFound
		return
	}

	// 有效处理同步修改tid & rid
	if gssp.State == model.Effective {
		if err = s.dao.TxSetGroupStateTid(tx, gids, gssp.State, newRid, gssp.Tid); err != nil {
			log.Error("s.TxSetGroupStateTid(%v,%d,%d) error(%v)", gids, gssp.State, gssp.Tid, err)
			return
		}
	} else {
		if err = s.dao.TxSimpleSetGroupState(tx, gids, gssp.State); err != nil {
			log.Error("s.TxSimpleSetGroupState(%v,%d,%d) error(%v)", gids, gssp.State, gssp.Tid, err)
			return
		}
	}

	// Set group handling count to 0
	if err = s.dao.TxBatchUpGroupHandling(tx, gids, 0); err != nil {
		log.Error("s.TxBatchUpGroupHandling(%+v, 0) error(%v)", gids, err)
		return
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}

	s.task(func() {
		// group bus object
		var gidToBus map[int64]*model.Business
		if gidToBus, err = s.dao.BusObjectByGids(context.Background(), gids); err != nil {
			log.Error("s.dao.BusObjectByGids(%v) error(%v)", gids, err)
			return
		}
		for _, g := range groups {
			g.BusinessObject = gidToBus[g.ID]
		}
		var result []*search.ChallSearchCommonData
		cond := &search.ChallSearchCommonCond{
			Fields: []string{"id", "gid", "mid", "state", "title", "oid", "tid", "business"},
			Gids:   gids,
			States: []int64{int64(model.Pending)},
		}
		if result, err = s.dao.SearchChallengeMultiPage(context.Background(), cond); err != nil {
			log.Error("s.dao.SearchChallengeMultiPage(%+v) error(%v)", cond, err)
			return
		}
		tinyChalls = make(map[int64]*model.TinyChall, len(cids)) // map[id]*tc
		for _, c := range result {
			cids = append(cids, c.ID)
			tc := &model.TinyChall{
				Cid:   c.ID,
				Gid:   c.Gid,
				Mid:   c.Mid,
				Title: c.Title,
			}
			if str, ok := c.State.(string); ok {
				st, _ := strconv.Atoi(str)
				tc.State = int8(st)
			}
			if f, ok := c.State.(float64); ok {
				tc.State = int8(math.Floor(f))
			}
			tinyChalls[c.ID] = tc
		}
		// async set challenge state
		if err = s.dao.BatchUpChallByIDs(cids, uint32(gssp.State), gssp.AdminID); err != nil {
			log.Error("s.dao.BatchUpChallByIDs(%v,%d) error(%v)", cids, gssp.State, err)
			return
		}
		s.afterSetGroupState(gssp, groups, tinyChalls)
	})
	return
}

// SetPublicReferee 移交众裁
func (s *Service) SetPublicReferee(c context.Context, gspr *param.GroupStatePublicReferee) (err error) {
	tx := s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("s.SetPublicReferee() panic(%v)", r)
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	var groups map[int64]*model.Group
	if groups, err = s.dao.TxGroups(tx, gspr.ID); err != nil {
		log.Error("s.dao.TxGroups(%v) error(%v)", gspr.ID, err)
		return
	}

	// ignore group if state not pending
	var gids, oids, eids []int64
	for id, g := range groups {
		if g.State != model.Pending || g.Business != gspr.Business {
			delete(groups, id)
			continue
		}
		gids = append(gids, id)
		oids = append(oids, g.Oid)
		eids = append(eids, g.Eid)
	}

	if len(gids) <= 0 {
		log.Warn("No pending groups found with conditon(%+v)", *gspr)
		err = ecode.WkfGroupNotFound
		return
	}
	// set state to public referee
	gspr.State = model.PublicReferee

	// start block add case
	// object
	var gidBus map[int64]*model.Business
	if gidBus, err = s.dao.BusObjectByGids(c, gids); err != nil {
		log.Error("s.dao.BusObjectByGids(%v) error(%v)", gids, err)
		err = ecode.WkfSetPublicRefereeFailed
		return
	}
	// external
	var metas map[int64]*model.GroupMeta
	if metas, err = s.externalMeta(c, gspr.Business, gids, oids, eids); err != nil {
		log.Error("s.SearchMeta(%d,%v,%v,%v) error(%v)", gspr.Business, gids, oids, eids, err)
		err = ecode.WkfSetPublicRefereeFailed
		return
	}

	var data []model.BlockCaseAdd
	for _, g := range groups {
		log.Info("start add case gid(%d)", g.ID)
		if _, ok := gidBus[g.ID]; !ok {
			log.Warn("gid(%d) not found bus object", g.ID)
			err = ecode.WkfSetPublicRefereeFailed
			return
		}
		bus := gidBus[g.ID]
		var (
			extra       map[string]interface{}
			external    map[string]interface{}
			link, title string
			ok          bool
			ctime       float64
		)
		if err = json.Unmarshal([]byte(bus.Extra), &extra); err != nil {
			log.Error("json.Unmarshal(%s) failed error(%v)", bus.Extra, err)
			err = ecode.WkfSetPublicRefereeFailed
			return
		}
		// redirect url
		if link, ok = extra["link"].(string); !ok {
			log.Error("gid(%d) assert business extra link failed", g.ID)
			err = ecode.WkfSetPublicRefereeFailed
			return
		}
		// object title
		if title, ok = extra["title"].(string); !ok {
			log.Error("gid(%d) assert business extra title failed", g.ID)
			err = ecode.WkfSetPublicRefereeFailed
			return
		}

		if _, ok = metas[g.ID]; !ok {
			log.Error("gid(%d) not found meta data", g.ID)
			err = ecode.WkfSetPublicRefereeFailed
			return
		}
		if external, ok = metas[g.ID].External.(map[string]interface{}); !ok {
			log.Error("gid(%d) external meta data assert failed", g.ID)
			err = ecode.WkfSetPublicRefereeFailed
			return
		}
		// business object ctime
		if ctime, ok = external["ctime"].(float64); !ok {
			log.Error("gid(%d) external ctime assert failed", g.ID)
			err = ecode.WkfSetPublicRefereeFailed
			return
		}

		d := model.BlockCaseAdd{
			RpID:          g.Eid,
			Oid:           g.Oid,
			Type:          g.Fid,
			Mid:           bus.Mid,
			Operator:      gspr.AdminName,
			OperID:        gspr.AdminID,
			OriginContent: bus.Title,
			ReasonType:    g.Tid,
			BusinessTime:  int64(ctime),
			OriginType:    int64(blocked.OriginReply), // fixme: support multi business
			OriginTitle:   title,
			OriginURL:     link,
		}
		data = append(data, d)
	}

	// request credit
	var content []byte
	if content, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%v) error(%v)", data, err)
		err = ecode.WkfSetPublicRefereeFailed
		return
	}
	uv := url.Values{}
	uv.Set("data", string(content))
	if err = s.dao.AddCreditCase(c, uv); err != nil {
		return
	}

	// set group set only
	if err = s.dao.TxSimpleSetGroupState(tx, gids, gspr.State); err != nil {
		log.Error("s.TxSimpleSetGroupState(%v,%d) error(%v)", gids, gspr.State, err)
		return
	}
	// set group handling count to 0
	if err = s.dao.TxBatchUpGroupHandling(tx, gids, 0); err != nil {
		log.Error("s.TxBatchUpGroupHandling(%+v, 0) error(%v)", gids, err)
		return
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}

	s.task(func() {
		// async set challenge state
		var (
			result []*search.ChallSearchCommonData
			cids   []int64
		)
		cond := &search.ChallSearchCommonCond{
			Fields: []string{"id"},
			Gids:   gids,
			States: []int64{int64(model.Pending)},
		}
		if result, err = s.dao.SearchChallengeMultiPage(context.Background(), cond); err != nil {
			log.Error("s.dao.SearchChallengeMultiPage(%+v) error(%v)", cond, err)
			return
		}
		for _, c := range result {
			cids = append(cids, c.ID)
		}
		if err = s.dao.BatchUpChallByIDs(cids, uint32(gspr.State), gspr.AdminID); err != nil {
			log.Error("s.dao.BatchUpChallByIDs(%v,%d) error(%v)", cids, gspr.State, err)
			return
		}
		s.afterSimpleSetState(gspr, groups)
	})
	return
}

// GroupPendingCount 当前 bid/rid 待办工单数
func (s *Service) GroupPendingCount(c context.Context, cond *search.GroupSearchCommonCond) (gpc *model.GroupPendingCount, err error) {
	var groupSearchCommonResp *search.GroupSearchCommonResp
	if groupSearchCommonResp, err = s.dao.SearchGroup(c, cond); err != nil {
		log.Error("Failed to s.dao.SearchGroup(%v): %v", cond, err)
		err = ecode.WkfSearchGroupFailed
		return
	}
	gpc = &model.GroupPendingCount{
		Total: groupSearchCommonResp.Page.Total,
	}
	return
}

// ExternalMeta external dependency
func (s *Service) externalMeta(c context.Context, business int8, gids, oids, eids []int64) (metas map[int64]*model.GroupMeta, err error) {
	metas = make(map[int64]*model.GroupMeta)
	// check if has external uri
	if _, ok := s.callbackCache[business]; !ok {
		return
	}
	uri := ""
	if uri = s.callbackCache[business].ExternalAPI; uri == "" {
		log.Warn("bid %d not found external api", business)
		return
	}

	// search meta
	//todo: common extra info
	var data map[string]interface{}
	if data, err = s.dao.CommonExtraInfo(c, business, uri, gids, oids, eids); err != nil {
		log.Error("s.dao.CommonExtraInfo() error(%v)", err)
		return
	}
	log.Info("bid(%d) external data(%v)", business, data)
	for gidStr, ext := range data {
		gid, _ := strconv.ParseInt(gidStr, 10, 64)
		metas[gid] = &model.GroupMeta{
			External: ext,
		}
	}
	return
}

// WarpMeta .
func (s *Service) warpMeta(business int8, group *model.Group, metas map[int64]*model.GroupMeta, archives map[int64]*model.Archive, users map[int64]*model.Account, bus map[int64]*model.Business) {
	// not has external data
	if _, ok := metas[group.ID]; !ok {
		metas[group.ID] = &model.GroupMeta{}
	}
	var (
		b  *model.Business
		ok bool
	)
	if b, ok = bus[group.ID]; ok {
		metas[group.ID].Object = b
		group.Defendant = users[b.Mid]
		metas[group.ID].Archive = archives[b.Oid]
	}
	group.MetaData = metas[group.ID]
	log.Info("WarpMeta gid(%d) oid(%d) eid(%d) meta.External(%+v)", group.ID, group.Oid, group.Eid, metas[group.ID].External)
}

// UpGroupExtra  update business extra of gid, only cover business extra field
func (s *Service) UpGroupExtra(c context.Context, uep *param.UpExtraParam) (err error) {
	if err = s.dao.UpExtraV3(uep.Gids, uep.AdminID, uep.Extra); err != nil {
		log.Error("s.dao.UpExtraV3(%v, %d, %v) error(%v)", uep.Gids, uep.AdminID, uep.Extra, err)
	}
	// todo: after up extra
	return
}

// dbCheck check group state & rid in db
func dbCheck(cond *search.GroupSearchCommonCond, rg *model.Group) (isCheck bool) {
	var (
		isStateCheck bool
		isRidCheck   bool
	)
	// check db state
	if len(cond.States) == 0 {
		isStateCheck = true
	} else {
		for _, state := range cond.States {
			if rg.State == state {
				isStateCheck = true
				break
			}
		}
	}

	// check rid
	if len(cond.RID) == 0 {
		isRidCheck = true
	} else {
		for _, rid := range cond.RID {
			if rg.Rid == rid {
				isRidCheck = true
				break
			}
		}
	}

	if isStateCheck && isRidCheck {
		isCheck = true
	}
	return
}

func round(num float64) int32 {
	return int32(num + math.Copysign(0.5, num))
}
