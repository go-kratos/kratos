package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/app/admin/main/workflow/model/search"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// SetChallResult will set a challenge result
func (s *Service) SetChallResult(c context.Context, crp *param.ChallResParam) (err error) {
	var (
		chall *model.Chall
		g     *model.Group
		rows  int64
	)
	// Check challenge and group is exist
	if chall, err = s.dao.Chall(c, crp.Cid); err != nil {
		log.Error("s.dao.Chall() error(%v)", err)
		return
	}
	if chall == nil {
		log.Error("Challenge(%d) not exist", crp.Cid)
		err = ecode.NothingFound
		return
	}

	chall.SetState(uint32(crp.State), 0)
	chall.AdminID = crp.AdminID

	tx := s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if rows, err = s.dao.TxUpChall(tx, chall); err != nil {
		tx.Rollback()
		log.Error("s.dao.TxUpChall(%v) err(%v)", chall, err)
		return
	}
	if g, err = s.dao.GroupByID(c, chall.Gid); err != nil {
		tx.Rollback()
		log.Error("s.dao.GroupByID() error(%v)", err)
		return
	}
	if g == nil {
		tx.Rollback()
		log.Error("Group(%d) not exist", chall.Gid)
		err = ecode.NothingFound
		return
	}

	// decrease group pending stat count
	g.Handling -= int32(rows)
	if g.Handling < 0 {
		g.Handling = 0
	}
	if err = s.dao.TxUpGroupHandling(tx, g.ID, g.Handling); err != nil {
		tx.Rollback()
		log.Error("Failed to update group stat count(%d) all(%d) by gid %d, error(%v)",
			g.Count, g.Handling, g.ID, err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("tx.Commit() error(%v)", err)
		return
	}

	s.task(func() {
		s.afterSetChallResult(crp, chall)
	})
	return
}

// BatchSetChallResult will set a set of challenges result
func (s *Service) BatchSetChallResult(c context.Context, bcrp *param.BatchChallResParam) (err error) {
	var (
		cids      []int64 // all requested cids
		pcids     []int64 // pending cids
		gids      []int64
		challs    map[int64]*model.Chall
		groups    map[int64]*model.Group
		gidToDecr map[int64]int32
	)
	// collect cids
	cids = append(cids, bcrp.Cids...)
	if challs, err = s.dao.StateChalls(c, cids, model.Pending); err != nil {
		log.Error("s.dao.StateChalls() error(%v)", err)
		return
	}
	// collect gids
	gidToDecr = make(map[int64]int32, len(challs))
	for _, c := range challs {
		pcids = append(pcids, int64(c.Cid))
		gids = append(gids, int64(c.Gid))
		if _, ok := gidToDecr[c.Gid]; !ok {
			gidToDecr[c.Gid] = 0
		}
		gidToDecr[c.Gid]++
	}
	if groups, err = s.dao.Groups(c, gids); err != nil {
		log.Error("s.dao.Groups() error(%v)", err)
		return
	}
	tx := s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()

	if err = s.dao.TxBatchUpChallByIDs(tx, pcids, bcrp.State); err != nil {
		tx.Rollback()
		log.Error("s.dao.TxBatchUpChallByIDs(%+v) error(%v)", bcrp, err)
		return
	}

	// decrease pending stat count by gidToDesc
	for gid, decr := range gidToDecr {
		g, ok := groups[gid]
		if !ok {
			log.Error("Failed to retrive group by gid %d", gid)
			continue
		}
		// decr group handling counts
		g.Handling -= decr
		if g.Handling < 0 {
			g.Handling = 0
		}
		if err = s.dao.TxUpGroupHandling(tx, gid, g.Handling); err != nil {
			tx.Rollback()
			log.Error("Failed to update group stat count(%d) all(%d) by gid %d, error(%v)",
				g.Count, g.Handling, g.ID, err)
			return
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("tx.Commit() error(%v)", err)
		return
	}

	//double write state to dispatch_state field
	changedChalls, err := s.dao.Challs(c, pcids)
	for cid := range changedChalls {
		changedChalls[cid].SetState(uint32(bcrp.State), 0)
		if err = s.dao.BatchUpChallByIDs([]int64{int64(cid)}, changedChalls[cid].DispatchState, bcrp.AdminID); err != nil {
			log.Error("s.dao.BatchUpChallByIDs(%v,%d,%d) error(%v)", cid, changedChalls[cid].DispatchState, bcrp.AdminID, err)
			return
		}
	}

	s.task(func() {
		s.afterBatchSetChallResult(bcrp, challs)
	})
	return
}

// UpChallExtraV2 will update extra data of a challenge by cid
func (s *Service) UpChallExtraV2(c context.Context, cep *param.ChallExtraParam) (err error) {
	var (
		obj   *model.Business
		chall *model.Chall
	)
	if chall, err = s.dao.Chall(c, cep.Cid); err != nil {
		log.Error("s.dao.Chall(%d) error(%v)", cep.Cid, err)
		return
	}
	if obj, err = s.dao.LastBusRec(c, chall.Business, chall.Oid); err != nil {
		log.Error("s.dao.LastBusRec(%d, %d) error(%v)", chall.Business, chall.Oid, err)
		return
	}
	if obj == nil {
		err = ecode.NothingFound
		return
	}
	parsed := make(map[string]interface{})
	if obj.Extra != "" {
		if err1 := json.Unmarshal([]byte(obj.Extra), &parsed); err1 != nil {
			log.Error("json.Unmarshal(%s) error(%v)", obj.Extra, err1)
		}
	}

	tx := s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("s.UpChallExtra() panic(%v)", r)
		}
	}()

	for k, v := range cep.Extra {
		parsed[k] = v
	}
	if _, err = s.dao.TxUpChallExtraV2(tx, chall.Business, chall.Oid, cep.AdminID, parsed); err != nil {
		tx.Rollback()
		log.Error("s.dao.UpChallExtra(%d, %d, %v) error(%v)", cep.Cid, cep.AdminID, parsed, err)
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	return
}

// BatchUpChallExtraV2 update business object extra field by business oid
func (s *Service) BatchUpChallExtraV2(c context.Context, bcep *param.BatchChallExtraParam) (err error) {
	var (
		challs map[int64]*model.Chall
		bs     *model.Business
	)
	if challs, err = s.dao.Challs(c, bcep.Cids); err != nil {
		log.Error("s.dao.Challs(%v) error(%v)", bcep.Cids, err)
		return
	}
	parsed := make(map[int64]map[string]interface{}, len(challs))
	for cid, chall := range challs {
		if bs, err = s.dao.LastBusRec(c, chall.Business, chall.Oid); err != nil {
			log.Error("s.dao.LastBusRec(%d, %d) error(%v)", chall.Business, chall.Oid, err)
			return
		}
		inner := make(map[string]interface{})
		if bs.Extra != "" {
			if err1 := json.Unmarshal([]byte(bs.Extra), &inner); err1 != nil {
				log.Error("json.Unmarshal(%s) error(%v)", bs.Extra, err1)
			}
		}
		parsed[cid] = inner
	}

	tx := s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	for _, ex := range parsed {
		for k, v := range bcep.Extra {
			ex[k] = v
		}
	}
	for cid, ex := range parsed {
		chall := challs[cid]
		if _, err = s.dao.TxUpChallExtraV2(tx, chall.Business, chall.Oid, bcep.AdminID, ex); err != nil {
			tx.Rollback()
			log.Error("s.dao.TxUpChallExtra(%d, %d, %v) error(%v)", cid, bcep.AdminID, ex, err)
			return
		}
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	return
}

// UpChallExtraV3 .
func (s *Service) UpChallExtraV3(c context.Context, cep3 *param.ChallExtraParamV3) (err error) {
	//todo
	return
}

// ChallList will list challenges by several conditions
// Deprecated
func (s *Service) ChallList(c context.Context, cond *search.ChallSearchCommonCond) (challPage *search.ChallListPageCommon, err error) {
	var (
		resp         *search.ChallSearchCommonResp
		cids         []int64
		challs       map[int64]*model.Chall
		challLastLog map[int64]string
		attPaths     map[int64][]string
		t            *model.TagMeta
	)

	if resp, err = s.dao.SearchChallenge(c, cond); err != nil {
		log.Error("s.dao.SearchChall() error(%v)", err)
		return
	}
	cids = make([]int64, 0, len(resp.Result))
	for _, r := range resp.Result {
		cids = append(cids, r.ID)
	}

	if challs, err = s.dao.Challs(c, cids); err != nil {
		log.Error("s.dao.Challs() error(%v)", err)
		return
	}
	// read state from new state field
	for cid := range challs {
		challs[cid].FromState()
	}

	if challLastLog, err = s.LastLog(c, cids, []int{model.WLogModuleChallenge}); err != nil {
		log.Error("s.dao.BatchLastLog() error(%v)", err)
		err = nil
	}

	if attPaths, err = s.dao.AttPathsByCids(c, cids); err != nil {
		log.Error("s.dao.AttPathsByCids() error(%v)", err)
		return
	}

	challPage = &search.ChallListPageCommon{}
	challList := make([]*model.Chall, 0, len(resp.Result))
	for _, r := range resp.Result {
		cid := r.ID
		c, ok := challs[cid]
		if !ok {
			log.Warn("Invalid challenge id %d", r.ID)
			continue
		}

		// fill tag
		if t, err = s.tag(c.Business, c.Tid); err != nil {
			log.Error("Failed to retrive tag by bid(%d) tag_id(%d)", c.Business, c.Tid)
			err = nil
		} else {
			c.Tag = t.Name
			c.Round = t.RID
		}

		// fill last log
		if l, ok := challLastLog[cid]; ok {
			c.LastLog = l
		}

		// fill attachments
		c.Attachments = make([]string, 0)
		if ps, ok := attPaths[cid]; ok {
			c.Attachments = ps
			c.FixAttachments()
		}

		c.FormatState()
		challList = append(challList, c)
	}

	challPage.Items = challList
	challPage.Page = &model.Page{
		Num:   resp.Page.Num,
		Size:  resp.Page.Size,
		Total: resp.Page.Total,
	}

	return
}

// ChallListCommon will list challenges by several conditions
func (s *Service) ChallListCommon(c context.Context, cond *search.ChallSearchCommonCond) (challPage *search.ChallListPageCommon, err error) {
	var (
		challSearchResp      *search.ChallSearchCommonResp
		business             int8
		cids                 []int64
		uids                 []int64
		oids                 []int64
		uNames               map[int64]string
		mids                 []int64
		challs               map[int64]*model.Chall
		challLastLog         map[int64]string
		challLastEvent       map[int64]*model.Event
		attPaths             map[int64][]string
		archives             map[int64]*model.Archive
		authors              map[int64]*model.Account
		lastAuditLogSchRes   *search.AuditLogSearchResult           // 最近的稿件操作日志
		lastAuditLogExtraMap map[int64]*search.ArchiveAuditLogExtra // 稿件操作日志 extra 信息
		chall                *model.Chall
		ps                   []string //attachments
		ok                   bool
		l                    string
	)
	business = cond.Business
	if challSearchResp, err = s.dao.SearchChallenge(c, cond); err != nil {
		log.Error("s.dao.SearchChallenge() error(%v)", err)
		return
	}

	// no result in es
	if challSearchResp.Page.Total == 0 {
		challPage = &search.ChallListPageCommon{}
		challPage.Items = make([]*model.Chall, 0)
		challPage.Page = &model.Page{
			Num:   challSearchResp.Page.Num,
			Size:  challSearchResp.Page.Size,
			Total: challSearchResp.Page.Total,
		}
		return
	}

	cids = make([]int64, 0, len(challSearchResp.Result))
	oids = make([]int64, 0, len(challSearchResp.Result))
	for _, r := range challSearchResp.Result {
		cids = append(cids, r.ID)
	}

	if challs, err = s.dao.Challs(c, cids); err != nil {
		log.Error("s.dao.Challs(%v) error(%v)", cids, err)
		return
	}
	// read state from new state field
	for cid, c := range challs {
		challs[cid].FromState()
		uids = append(uids, int64(c.AdminID))
		uids = append(uids, int64(c.AssigneeAdminID))
		oids = append(oids, c.Oid)
		mids = append(mids, c.Mid)
	}

	if challLastLog, err = s.LastLog(c, cids, []int{model.WLogModuleChallenge}); err != nil {
		log.Error("s.batchLastLog(%v,%v) error(%v)", cids, model.WLogModuleChallenge, err)
		err = nil
	}

	// admin unames
	if uNames, err = s.dao.BatchUNameByUID(c, uids); err != nil {
		log.Error("s.dao.SearchUNameByUid(%v) error(%v)", uids, err)
		err = nil
	}

	if attPaths, err = s.dao.AttPathsByCids(c, cids); err != nil {
		log.Error("s.dao.AttPathsByCids() error(%v)", err)
		return
	}

	if challLastEvent, err = s.batchLastEvent(c, cids); err != nil {
		log.Error("s.batchLastEvent(%v) error(%v)", cids, err)
		return
	}
	switch business {
	case model.ArchiveAppeal:
		if archives, err = s.dao.ArchiveRPC(c, oids); err != nil {
			log.Error("s.dao.ArchiveInfos(%v) error(%v)", oids, err)
			err = nil
		} else {
			cond := &search.AuditLogGroupSearchCond{
				Businesses: []int64{3},
				Order:      "ctime",
				PS:         100,
				PN:         1,
				Sort:       "desc",
				Oids:       oids,
				Group:      []string{"oid"},
			}
			if lastAuditLogSchRes, err = s.dao.SearchAuditLogGroup(c, cond); err != nil {
				log.Error("s.dao.SearchAuditLogGroup(%+v) error(%v)", cond, err)
				err = nil
			}
			lastAuditLogExtraMap = make(map[int64]*search.ArchiveAuditLogExtra)
			for _, l := range lastAuditLogSchRes.Data.Result {
				extra := new(search.ArchiveAuditLogExtra)
				if err = json.Unmarshal([]byte(l.ExtraData), extra); err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", l.ExtraData, err)
					continue
				}
				lastAuditLogExtraMap[l.Oid] = extra
			}
		}
	case model.CreditAppeal:
		authors = s.dao.AccountInfoRPC(c, mids)

	case model.ArchiveAudit:
		if archives, err = s.dao.ArchiveRPC(c, oids); err != nil {
			log.Error("s.dao.ArchiveInfos(%v) error(%v)", oids, err)
			err = nil
		} else {
			cond := &search.AuditLogGroupSearchCond{
				Businesses: []int64{3},
				Order:      "ctime",
				PS:         100,
				PN:         1,
				Sort:       "desc",
				Oids:       oids,
				Group:      []string{"oid"},
			}
			if lastAuditLogSchRes, err = s.dao.SearchAuditLogGroup(c, cond); err != nil {
				log.Error("s.dao.SearchAuditLogGroup(%+v) error(%v)", cond, err)
				err = nil
			}
			lastAuditLogExtraMap = make(map[int64]*search.ArchiveAuditLogExtra)
			for _, l := range lastAuditLogSchRes.Data.Result {
				extra := new(search.ArchiveAuditLogExtra)
				if err = json.Unmarshal([]byte(l.ExtraData), extra); err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", l.ExtraData, err)
					continue
				}
				lastAuditLogExtraMap[l.Oid] = extra
			}
		}
	}

	challPage = &search.ChallListPageCommon{}
	challList := make([]*model.Chall, 0, len(challSearchResp.Result))
	for _, r := range challSearchResp.Result {
		cid := r.ID
		chall, ok = challs[cid]
		if !ok {
			log.Warn("Invalid challenge id %d", r.ID)
			continue
		}

		// fill tag
		var t *model.TagMeta
		if t, err = s.tag(chall.Business, chall.Tid); err != nil {
			log.Error("s.tag(%d,%d) error(%v)", chall.Business, chall.Tid, err)
			err = nil
		} else {
			chall.Tag = t.Name
			chall.Round = t.RID
		}

		// fill last log
		if l, ok = challLastLog[cid]; ok {
			chall.LastLog = l
		}

		// fill last event
		chall.LastEvent = challLastEvent[cid]

		// fill attachments
		chall.Attachments = make([]string, 0)
		if ps, ok = attPaths[cid]; ok {
			chall.Attachments = ps
			chall.FixAttachments()
		}
		chall.FormatState()

		//fill business object
		if chall.BusinessObject, err = s.dao.LastBusRec(c, chall.Business, chall.Oid); err != nil {
			log.Error("s.dao.LastBusRec(%d, %d) error(%v)", chall.Business, chall.Oid, err)
			err = nil
		}

		//todo: add challenge meta
		switch business {
		case model.ArchiveAppeal: // 稿件申诉
			var (
				archive *model.Archive
				extra   *search.ArchiveAuditLogExtra
			)
			if archive, ok = archives[chall.Oid]; !ok {
				log.Warn("failed get archive info oid(%d) cid(%d)", chall.Oid, r.ID)
			}
			if extra, ok = lastAuditLogExtraMap[chall.Oid]; !ok {
				log.Warn("not exist archive operate last audit log extra oid(%d) cid(%d)", chall.Oid, r.ID)
			}
			if extra != nil && archive != nil {
				archive.OPName = extra.Content.UName
				archive.OPContent = extra.Diff
				archive.OPRemark = extra.Content.Note
			}
			chall.Meta = archive
		case model.CreditAppeal: //小黑屋
			chall.Meta = chall.BusinessObject
			if _, ok = authors[chall.Mid]; ok {
				chall.MName = authors[chall.Mid].Name
			}
		case model.ArchiveAudit: //稿件审核
			var archive *model.Archive
			var extra *search.ArchiveAuditLogExtra
			if archive, ok = archives[chall.Oid]; !ok {
				log.Warn("failed get archive info oid(%d) cid(%d)", chall.Oid, r.ID)
			}
			if extra, ok = lastAuditLogExtraMap[chall.Oid]; !ok {
				log.Warn("not exist archive operate last audit log extra oid(%d) cid(%d)", chall.Oid, r.ID)
			}
			if extra != nil && archive != nil {
				archive.OPName = extra.Content.UName
				archive.OPContent = extra.Diff
				archive.OPRemark = extra.Content.Note
			}
			chall.Meta = archive
		}

		chall.AssigneeAdminName = uNames[chall.AssigneeAdminID]
		chall.AdminName = uNames[chall.AdminID]
		challList = append(challList, chall)
	}

	challPage.Items = challList
	challPage.Page = &model.Page{
		Num:   challSearchResp.Page.Num,
		Size:  challSearchResp.Page.Size,
		Total: challSearchResp.Page.Total,
	}
	return
}

// ChallDetail will retrive challenge by cid
func (s *Service) ChallDetail(c context.Context, cid int64) (chall *model.Chall, err error) {
	var (
		cl        map[int64]string
		attPaths  []string
		t         *model.TagMeta
		bs        *model.Business
		lastEvent *model.Event
		archives  map[int64]*model.Archive
		authors   map[int64]*model.Account
		ok        bool
		l         string
	)
	if chall, err = s.dao.Chall(c, cid); err != nil || chall == nil {
		log.Error("Failed to s.dao.Chall(%d) or chall not found: %v", cid, err)
		err = ecode.NothingFound
		return
	}
	// read state from new state field
	chall.FromState()

	if attPaths, err = s.dao.AttPathsByCid(c, cid); err != nil {
		log.Error("Failed to s.dao.AttPathsByCid(%d): %v", cid, err)
		err = nil
	}
	if cl, err = s.LastLog(c, []int64{cid}, []int{model.WLogModuleChallenge}); err != nil {
		log.Error("Failed to s.dao.LastLog(%d): %v", cid, err)
		err = nil
	}
	if lastEvent, err = s.dao.LastEventByCid(c, cid); err != nil {
		log.Error("Failed to s.dao.EventsByCid(%d): %v", cid, err)
		err = nil
	}
	if t, err = s.tag(chall.Business, chall.Tid); err != nil {
		log.Error("Failed to s.tag(%d,%d) error(%v)", chall.Business, chall.Tid, err)
		// fixme: to debug
		err = nil
	} else {
		chall.Tag = t.Name
		chall.Round = t.RID
	}
	if bs, err = s.dao.LastBusRec(c, chall.Business, chall.Oid); err != nil {
		log.Error("Failed to s.dao.BusRecByCid(%d): %v", cid, err)
		err = nil
	}

	business := chall.Business
	switch business {
	case model.ArchiveAppeal:
		if archives, err = s.dao.ArchiveRPC(c, []int64{chall.Oid}); err != nil {
			log.Error("s.dao.ArchiveInfos(%v) error(%v)", chall.Oid, err)
			err = nil
		}
	}

	switch business {
	case model.ArchiveAppeal: // 稿件申诉
		if chall.Meta, ok = archives[chall.Oid]; !ok {
			log.Warn("failed get archive info oid(%d) cid(%d)", chall.Oid, chall.Cid)
		}
	case model.CreditAppeal: //小黑屋
		chall.Meta = chall.BusinessObject
		if authors = s.dao.AccountInfoRPC(c, []int64{chall.Mid}); authors != nil {
			if _, ok = authors[chall.Mid]; ok {
				chall.MName = authors[chall.Mid].Name
			}
		}
	}

	if l, ok = cl[cid]; ok {
		chall.LastLog = l
	} else {
		log.Info("cid(%d) not found last log", cid)
	}
	chall.LastEvent = lastEvent
	chall.BusinessObject = bs

	chall.Attachments = attPaths
	chall.FixAttachments()

	return
}

// UpChallBusState will update business_state field of a challenge
func (s *Service) UpChallBusState(c context.Context, cid int64, assigneeAdminid int64, assigneeAdminName string, busState int8) (err error) {
	// TODO(zhoujiahui): record in log?
	if cid <= 0 {
		err = ecode.WkfChallNotFound
		return
	}

	if err = s.dao.UpChallBusState(c, cid, busState, assigneeAdminid); err != nil {
		return
	}
	s.task(func() {
		var (
			result []*search.ChallSearchCommonData
			challs []*model.Chall
		)
		cond := &search.ChallSearchCommonCond{
			Fields: []string{"id", "oid", "business", "mid", "typeid"},
			IDs:    []int64{cid},
		}
		if result, err = s.dao.SearchChallengeMultiPage(context.Background(), cond); err != nil {
			log.Error("s.dao.SearchChallengeMultiPage(%+v) error(%v)", cond, err)
			return
		}
		for _, r := range result {
			c := &model.Chall{
				Cid:               r.ID,
				Oid:               r.Oid,
				Business:          r.Business,
				BusinessState:     busState,
				AssigneeAdminID:   assigneeAdminid,
				AssigneeAdminName: assigneeAdminName,
				Mid:               r.Mid,
				TypeID:            r.TypeID,
			}
			challs = append(challs, c)
		}
		s.afterSetBusinessState(challs)
	})
	return
}

// BatchUpChallBusState will update business_state field of a set of challenges
func (s *Service) BatchUpChallBusState(c context.Context, cids []int64, assigneeAdminid int64, assigneeAdminName string, busState int8) (err error) {
	// TODO(zhoujiahui): record in log?
	if len(cids) <= 0 {
		return
	}

	if err = s.dao.BatchUpChallBusState(c, cids, busState, assigneeAdminid); err != nil {
		return
	}

	s.task(func() {
		var (
			result []*search.ChallSearchCommonData
			challs []*model.Chall
		)
		cond := &search.ChallSearchCommonCond{
			Fields: []string{"id", "oid", "business", "mid", "typeid"},
			IDs:    cids,
		}
		if result, err = s.dao.SearchChallengeMultiPage(context.Background(), cond); err != nil {
			log.Error("s.dao.SearchChallengeMultiPage(%+v) error(%v)", cond, err)
			return
		}
		for _, r := range result {
			c := &model.Chall{
				Cid:               r.ID,
				Oid:               r.Oid,
				Business:          r.Business,
				BusinessState:     busState,
				AssigneeAdminID:   assigneeAdminid,
				AssigneeAdminName: assigneeAdminName,
				Mid:               r.Mid,
				TypeID:            r.TypeID,
			}
			challs = append(challs, c)
		}
		s.afterSetBusinessState(challs)
	})
	return
}

// SetChallBusState will update business_state field of a set of challenges
func (s *Service) SetChallBusState(c context.Context, bcbsp *param.BatchChallBusStateParam) (err error) {
	// TODO(zhoujiahui): record in log?
	if len(bcbsp.Cids) <= 0 {
		return
	}

	if err = s.dao.BatchUpChallBusState(c, bcbsp.Cids, bcbsp.BusState, bcbsp.AssigneeAdminID); err != nil {
		return
	}

	s.task(func() {
		var (
			result []*search.ChallSearchCommonData
			challs []*model.Chall
		)
		cond := &search.ChallSearchCommonCond{
			Fields: []string{"id", "oid", "business", "mid", "typeid"},
			IDs:    bcbsp.Cids,
		}
		if result, err = s.dao.SearchChallengeMultiPage(context.Background(), cond); err != nil {
			log.Error("s.dao.SearchChallengeMultiPage(%+v) error(%v)", cond, err)
			return
		}
		for _, r := range result {
			c := &model.Chall{
				Cid:               r.ID,
				Oid:               r.Oid,
				Business:          r.Business,
				BusinessState:     bcbsp.BusState,
				AssigneeAdminID:   bcbsp.AssigneeAdminID,
				AssigneeAdminName: bcbsp.AssigneeAdminName,
				Mid:               r.Mid,
				TypeID:            r.TypeID,
			}
			challs = append(challs, c)
		}
		s.afterSetBusinessState(challs)
	})
	return
}

// UpBusChallsBusState will update business_state field of a set of challenges with same business and oid
func (s *Service) UpBusChallsBusState(c context.Context, business, busState int8, preBusStates []int8, oid int64, assigneeAdminid int64, extra map[string]interface{}) (cids []int64, err error) {
	// TODO(zhoujiahui): record in log?
	tx := s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()

	if cids, err = s.dao.TxChallsByBusStates(tx, business, oid, preBusStates); err != nil {
		tx.Rollback()
		log.Error("s.dao.TxChallsByBusStates(%d, %d, %s) error(%v)", business, oid, preBusStates, err)
		return
	}

	if err = s.dao.TxUpChallsBusStateByIDs(tx, cids, busState, assigneeAdminid); err != nil {
		tx.Rollback()
		log.Error("s.dao.TxUpChallsBusStateByIDs(%s, %d) error(%v)", cids, busState, err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("tx.Commit() error(%v)", err)
		return
	}

	bcep := &param.BatchChallExtraParam{
		Cids:    cids,
		AdminID: assigneeAdminid,
		Extra:   extra,
	}

	if err = s.BatchUpChallExtraV2(c, bcep); err != nil {
		log.Error("s.BatchUpChallExtra(%v) error(%v)", bcep, err)
		return
	}

	// double write to new field
	challs, err := s.dao.Challs(c, cids)
	if err != nil {
		return
	}
	for cid := range challs {
		challs[cid].SetState(uint32(busState), uint8(1))
		if err = s.dao.ORM.Table("workflow_chall").Where("id=?", cid).Update("dispatch_state", challs[cid].DispatchState).Error; err != nil {
			err = errors.Wrapf(err, "cid(%d), dispatch_state(%d)", cid, challs[cid].DispatchState)
			return
		}
	}

	return
}

// RstChallResult will reset challenge and its linked group state as Pending
func (s *Service) RstChallResult(c context.Context, crp *param.ChallRstParam) (err error) {
	var (
		chall *model.Chall
		group *model.Group
	)
	if chall, err = s.dao.Chall(c, crp.Cid); err != nil || chall == nil {
		log.Error("Failed to query challenge(%d) or it not exist: %v", crp.Cid, err)
		return
	}
	if group, err = s.dao.GroupByID(c, chall.Gid); err != nil || group == nil {
		log.Error("Failed to query group(%d) or it not exist: %v", chall.Gid, err)
		return
	}
	// update new field
	chall.SetState(uint32(crp.State), 0)

	tx := s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("s.RstChallResult() panic(%v)", r)
		}
	}()
	if err = tx.Model(chall).UpdateColumn(map[string]interface{}{
		"state":          crp.State,
		"dispatch_state": chall.DispatchState,
	}).Error; err != nil {
		tx.Rollback()
		log.Error("Failed to set chall(%v) as pending: %v", chall, err)
		return
	}
	if err = tx.Model(group).UpdateColumn(map[string]interface{}{
		"state": crp.State,
	}).Error; err != nil {
		tx.Rollback()
		log.Error("Failed to set group(%v) as pending: %v", group, err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("Failed to tx.Commit() in RstChallResult: %v", err)
		return
	}

	s.task(func() {
		cl := &model.WLog{
			AdminID:  crp.AdminID,
			Admin:    crp.AdminName,
			Oid:      chall.Oid,
			Business: chall.Business,
			Target:   chall.Cid,
			Module:   model.WLogModuleChallenge,
			Remark:   fmt.Sprintf(`“工单详情编号 %d”设置为 %s 移交复审`, chall.Cid, s.StateDescr(chall.Business, 0, crp.State)),
			Note:     crp.Reason,
		}
		gl := &model.WLog{
			AdminID:  crp.AdminID,
			Admin:    crp.AdminName,
			Oid:      group.Oid,
			Business: group.Business,
			Target:   group.ID,
			Module:   model.WLogModuleGroup,
			Remark:   fmt.Sprintf(`“工单编号 %d”设置为 %s 移交复审`, group.ID, s.StateDescr(chall.Business, 0, crp.State)),
			Note:     crp.Reason,
		}
		s.writeAuditLog(cl)
		s.writeAuditLog(gl)
	})
	return
}

// BusinessList will retrive business object by cids
// Deprecated
func (s *Service) BusinessList(c context.Context, cids []int64) (cidToBus map[int64]*model.Business, err error) {
	if cidToBus, err = s.dao.BatchBusRecByCids(c, cids); err != nil {
		log.Error("s.dao.BatchBusRecByCids(%v) error(%v)", cids, err)
		return
	}
	return
}

// UpChall will update challenge tid
func (s *Service) UpChall(c context.Context, cup *param.ChallUpParam) (err error) {
	var (
		chall *model.Chall
		t     *model.TagMeta
	)
	// Check group and tag is exist
	if chall, err = s.dao.Chall(c, cup.Cid); err != nil {
		log.Error("s.dao.Chall(%d) error(%v)", cup.Cid, err)
		return
	}
	if chall == nil {
		log.Error("Challenge(%d) not exist", cup.Cid)
		err = ecode.WkfChallNotFound
		return
	}
	if t, err = s.tag(chall.Business, cup.Tid); err != nil {
		log.Error("bid(%d) tag_id(%d) not found in cache", chall.Business, cup.Tid)
		return
	}
	tx := s.dao.ORM.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.ORM.Begin() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if err = s.dao.TxUpChallTag(tx, cup.Cid, cup.Tid); err != nil {
		tx.Rollback()
		log.Error("s.TxUpChallTag(%d, %d) error(%v)", cup.Cid, cup.Tid, err)
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.task(func() {
		l := &model.WLog{
			Oid:      chall.Oid,
			Business: chall.Business,
			Target:   chall.Cid,
			Module:   model.WLogModuleChallenge,
			AdminID:  cup.AdminID,
			Admin:    cup.AdminName,
			Remark:   fmt.Sprintf(`工单编号 %d “管理 Tag”更新为“%s”`, cup.Cid, t.Name),
		}
		s.writeAuditLog(l)
	})
	return
}

// ChallListV3 .
func (s *Service) ChallListV3(c context.Context, cond *search.ChallSearchCommonCond) (challPage *search.ChallListPageCommon, err error) {
	var (
		challSearchResp *search.ChallSearchCommonResp
		uNames          map[int64]string
		challs          map[int64]*model.Chall
		challLastLog    map[int64]string
		challLastEvent  map[int64]*model.Event
		attPaths        map[int64][]string
		gidToBus        map[int64]*model.Business
		users           map[int64]*model.Account
		b               *model.Business
		chall           *model.Chall
		ps              []string //attachments
		ok              bool
		meta            interface{}
	)
	if challSearchResp, err = s.dao.SearchChallenge(c, cond); err != nil {
		log.Error("s.dao.SearchChallenge() error(%v)", err)
		err = ecode.WkfSearchChallFailed
		return
	}
	// no result in es
	if challSearchResp.Page.Total == 0 {
		challPage = &search.ChallListPageCommon{}
		challPage.Items = make([]*model.Chall, 0)
		challPage.Page = &model.Page{
			Num:   challSearchResp.Page.Num,
			Size:  challSearchResp.Page.Size,
			Total: challSearchResp.Page.Total,
		}
		return
	}
	cids := make([]int64, 0, len(challSearchResp.Result))
	oids := make([]int64, 0, len(challSearchResp.Result))
	mids := make([]int64, 0, len(challSearchResp.Result))
	gids := make([]int64, 0, len(challSearchResp.Result))
	uids := make([]int64, 0, len(challSearchResp.Result)*2)
	for _, r := range challSearchResp.Result {
		cids = append(cids, r.ID)
		oids = append(oids, r.Oid)
		mids = append(mids, r.Mid)
		gids = append(gids, r.Gid)
	}
	if challs, err = s.dao.Challs(c, cids); err != nil {
		log.Error("s.dao.Challs(%v) error(%v)", cids, err)
		return
	}
	// read state from new state field
	for cid, c := range challs {
		challs[cid].FromState()
		uids = append(uids, int64(c.AdminID))
		uids = append(uids, int64(c.AssigneeAdminID))
	}
	if challLastLog, err = s.LastLog(c, cids, []int{model.WLogModuleChallenge}); err != nil {
		log.Error("s.LastLog(%v,%v) error(%v)", cids, model.WLogModuleChallenge, err)
		err = nil
	}

	// admin unames
	if uNames, err = s.dao.BatchUNameByUID(c, uids); err != nil {
		log.Error("s.dao.SearchUNameByUid(%v) error(%v)", uids, err)
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
	// user account
	users = s.dao.AccountInfoRPC(c, mids)
	// load meta
	meta = s.searchChallMeta(c, cond.Business, oids, users)

	challPage = &search.ChallListPageCommon{}
	challList := make([]*model.Chall, 0, len(challSearchResp.Result))
	for _, r := range challSearchResp.Result {
		var (
			t *model.TagMeta
			l string
		)
		cid := r.ID
		chall, ok = challs[cid]
		if !ok {
			log.Warn("Invalid challenge id %d", r.ID)
			continue
		}
		// fill tag
		if t, err = s.tag(chall.Business, chall.Tid); err != nil {
			log.Error("s.tag(%d,%d) error(%v)", chall.Business, chall.Tid, err)
			err = nil
		} else {
			chall.Tag = t.Name
			chall.Round = t.RID
		}
		// fill last log
		if l, ok = challLastLog[cid]; ok {
			chall.LastLog = l
		}
		// fill last event
		chall.LastEvent = challLastEvent[cid]
		// fill attachments
		chall.Attachments = make([]string, 0)
		if ps, ok = attPaths[cid]; ok {
			chall.Attachments = ps
			chall.FixAttachments()
		}
		chall.FormatState()

		//fill business object
		if b, ok = gidToBus[chall.Gid]; ok {
			chall.BusinessObject = b
		}
		//fill meta
		s.wrapChallMeta(cond.Business, chall, meta)

		chall.AssigneeAdminName = uNames[chall.AssigneeAdminID]
		chall.AdminName = uNames[chall.AdminID]
		if chall.Producer, ok = users[chall.Mid]; !ok {
			log.Warn("failed get producer info mid(%d)", chall.Mid)
		}
		chall.OidStr = strconv.FormatInt(chall.Oid, 10)
		challList = append(challList, chall)
	}
	challPage.Items = challList
	challPage.Page = &model.Page{
		Num:   challSearchResp.Page.Num,
		Size:  challSearchResp.Page.Size,
		Total: challSearchResp.Page.Total,
	}
	return
}

func (s *Service) searchChallMeta(c context.Context, business int8, oids []int64, users map[int64]*model.Account) (meta interface{}) {
	var (
		ok       bool
		err      error
		archives map[int64]*model.Archive
	)

	switch business {
	case model.ArchiveAppeal, model.ArchiveAudit:
		var (
			resp                 *search.AuditLogSearchCommonResult
			lastAuditLogExtraMap = make(map[int64]*search.ArchiveAuditLogExtra)
		)
		if archives, err = s.dao.ArchiveRPC(c, oids); err != nil {
			log.Error("s.dao.ArchiveInfosV3(%v) error(%v)", oids, err)
			err = nil
			return archives
		}

		cond := &search.AuditReportSearchCond{
			Business:      3,
			Fields:        []string{"oid", "extra_data"},
			Order:         "ctime",
			Sort:          "desc",
			Oid:           oids,
			Distinct:      "oid",
			IndexTimeType: "month",
			IndexTimeFrom: time.Now().AddDate(0, -6, 0),
			IndexTimeEnd:  time.Now(),
		}

		if resp, err = s.dao.SearchAuditReportLog(c, cond); err != nil {
			log.Error("s.dao.SearchAuditReportLog(%+v) error(%v)", cond, err)
			err = nil
		}
		for _, l := range resp.Result {
			extra := new(search.ArchiveAuditLogExtra)
			if err = json.Unmarshal([]byte(l.ExtraData), extra); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", l.ExtraData, err)
				continue
			}
			lastAuditLogExtraMap[l.Oid] = extra
		}
		for _, l := range resp.Result {
			extra := new(search.ArchiveAuditLogExtra)
			if err = json.Unmarshal([]byte(l.ExtraData), extra); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", l.ExtraData, err)
				continue
			}
		}
		for oid, a := range archives {
			if a.Composer, ok = users[a.Mid]; !ok {
				log.Warn("failed get account info mid(%d)", a.Mid)
				a.Composer = &model.Account{}
				continue
			}
			var extra *search.ArchiveAuditLogExtra
			if extra, ok = lastAuditLogExtraMap[oid]; !ok {
				log.Warn("failed get audit log of archive(%d)", oid)
				continue
			} else {
				a.OPName = extra.Content.UName
				a.OPContent = extra.Diff
				a.OPRemark = extra.Content.Note
			}
		}
	}
	return
}

func (s *Service) wrapChallMeta(business int8, chall *model.Chall, meta interface{}) {
	switch business {
	case model.ArchiveAppeal: // 稿件申诉
		var (
			a        *model.Archive
			archives map[int64]*model.Archive
			ok       bool
		)
		if archives, ok = meta.(map[int64]*model.Archive); !ok {
			return
		}
		if a, ok = archives[chall.Oid]; !ok {
			log.Warn("failed get archive info oid(%d) cid(%d)", chall.Oid, chall.Cid)
			chall.Meta = struct{}{}
			break
		}
		chall.Meta = a

	case model.CreditAppeal: //小黑屋
		if chall.BusinessObject == nil {
			log.Warn("can not get credit appeal info oid(%d) cid(%d) business(%d)", chall.Oid, chall.Cid, chall.Business)
			break
		}
		cMeta := &model.CreditMeta{
			Business: chall.BusinessObject,
		}
		chall.Meta = cMeta

	case model.ArchiveAudit: //稿件审核
		var (
			a        *model.Archive
			archives map[int64]*model.Archive
			ok       bool
		)
		if archives, ok = meta.(map[int64]*model.Archive); !ok {
			return
		}
		if a, ok = archives[chall.Oid]; !ok {
			log.Warn("failed get archive info oid(%d) cid(%d)", chall.Oid, chall.Cid)
		}
		chall.Meta = a
	}
}
