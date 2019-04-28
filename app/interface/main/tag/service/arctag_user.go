package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-common/app/interface/main/tag/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	taGrpcModel "go-common/app/service/main/tag/api"
	rpcModel "go-common/app/service/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
)

const (
	_adminFace = "http://static.hdslb.com/images/member/noface.gif"
	_adminName = "管理员"
)

func keySingle(oid int64, tp int32) string {
	return fmt.Sprintf("%d_%d", oid, tp)
}

// ArcTags .
func (s *Service) ArcTags(c context.Context, aid, mid int64) (reTs []*model.Tag, err error) {
	reTs = make([]*model.Tag, 0)
	v, err, _ := s.singleGroup.Do(keySingle(aid, rpcModel.ResTypeArchive), func() (res interface{}, err error) {
		return s.dao.ResTagMap(c, aid, rpcModel.ResTypeArchive)
	})
	if err != nil {
		return
	}
	arcTags := v.(map[int64]*taGrpcModel.Resource)
	var (
		tids    []int64
		allTids []int64
	)
	for _, v := range arcTags {
		if v.State != model.ResTagStateNormal && v.State != model.ResTagStateRegion {
			continue
		}
		allTids = append(allTids, v.Tid)
		if v.State == model.ResTagStateNormal {
			tids = append(tids, v.Tid)
		}
	}
	if len(tids) == 0 {
		return
	}
	_, channelIDs := s.resourceChannel(allTids, model.ManagerNo)
	channelMap := make(map[int64]struct{}, len(channelIDs))
	for _, tid := range channelIDs {
		channelMap[tid] = struct{}{}
		allTids = append(allTids, tid)
	}
	tagMap, err := s.dao.TagMap(c, allTids, mid)
	if err != nil {
		return
	}
	actionMap, _ := s.resAction(c, mid, aid, rpcModel.ResTypeArchive, tids)
	resTags := make([]*model.Tag, 0, len(channelIDs)+len(tids))
	for _, tid := range channelIDs {
		if tag, ok := tagMap[tid]; ok && tag.State == model.TagStateNormal {
			t := &model.Tag{
				ID:           tag.Id,
				Name:         tag.Name,
				Cover:        tag.Cover,
				Content:      tag.Content,
				ShortContent: tag.ShortContent,
				HeadCover:    tag.HeadCover,
				Type:         int8(tag.Type),
				State:        int8(tag.State),
				CTime:        tag.Ctime,
				MTime:        tag.Mtime,
				IsAtten:      int8(tag.Attention),
				Attribute:    int8(tag.Attr),
			}
			t.Count.Use = int(tag.Bind)
			t.Count.Atten = int(tag.Sub)
			resTags = append(resTags, t)
		}
	}
	for _, tid := range tids {
		if _, ok := channelMap[tid]; ok {
			continue
		}
		resource, ok := arcTags[tid]
		if !ok {
			continue
		}
		tag, ok := tagMap[tid]
		if !ok || tag.State != model.TagStateNormal {
			continue
		}
		t := &model.Tag{
			ID:           tag.Id,
			Name:         tag.Name,
			Cover:        tag.Cover,
			Content:      tag.Content,
			ShortContent: tag.ShortContent,
			HeadCover:    tag.HeadCover,
			Type:         int8(tag.Type),
			State:        int8(tag.State),
			CTime:        tag.Ctime,
			MTime:        tag.Mtime,
			IsAtten:      int8(tag.Attention),
			Role:         int8(resource.Role),
			Likes:        int64(resource.Like),
			Hates:        int64(resource.Hate),
			Attribute:    int8(resource.Attr),
		}
		t.Count.Use = int(tag.Bind)
		t.Count.Atten = int(tag.Sub)
		if action, ok := actionMap[tid]; ok {
			switch action {
			case rpcModel.UserActionLike:
				t.Liked = 1
			case rpcModel.UserActionHate:
				t.Hated = 1
			}
		}
		reTs = append(reTs, t)
	}
	sort.Sort(model.Tags(reTs))
	resTags = append(resTags, reTs...)
	return resTags, err
}

// MutiArcTags muti archive tags.
func (s *Service) MutiArcTags(c context.Context, mid int64, oids []int64) (res map[int64][]*model.Tag, err error) {
	res = make(map[int64][]*model.Tag, len(oids))
	resouceMap, err := s.dao.ResTags(c, oids, rpcModel.ResTypeArchive)
	if err != nil {
		return
	}
	var (
		tids          = make([]int64, 0)
		arcTidsMap    = make(map[int64][]int64, len(oids))
		allArcTidsMap = make(map[int64][]int64, len(oids))
	)
	for oid, arcTags := range resouceMap {
		var (
			arcTids    = make([]int64, 0, len(arcTags))
			allArcTids = make([]int64, 0, len(arcTags))
		)
		for _, v := range arcTags {
			if v.State != model.ResTagStateNormal && v.State != model.ResTagStateRegion {
				continue
			}
			allArcTids = append(allArcTids, v.Tid)
			tids = append(tids, v.Tid)
			if v.State == model.ResTagStateNormal {
				arcTids = append(arcTids, v.Tid)
			}
		}
		arcTidsMap[oid] = arcTids
		allArcTidsMap[oid] = allArcTids
	}
	channelMap, channelIDs := s.resourceChannels(c, model.ManagerNo, allArcTidsMap)
	tids = append(tids, channelIDs...)
	tagMap, err := s.dao.TagMap(c, tids, mid)
	if err != nil {
		return
	}
	for oid, v := range arcTidsMap {
		k, ok := channelMap[oid]
		var arcTags = make([]*model.Tag, 0, len(v)+len(k))
		for channelID := range k {
			if tag, exist := tagMap[channelID]; exist && tag.State == model.TagStateNormal {
				t := &model.Tag{
					ID:           tag.Id,
					Name:         tag.Name,
					Cover:        tag.Cover,
					Content:      tag.Content,
					ShortContent: tag.ShortContent,
					HeadCover:    tag.HeadCover,
					Type:         int8(tag.Type),
					State:        int8(tag.State),
					CTime:        tag.Ctime,
					MTime:        tag.Mtime,
					IsAtten:      int8(tag.Attention),
					Attribute:    int8(tag.Attr),
				}
				t.Count.Use = int(tag.Bind)
				t.Count.Atten = int(tag.Sub)
				arcTags = append(arcTags, t)
			}
		}
		for _, tid := range v {
			if ok {
				if _, b := k[tid]; b {
					continue
				}
			}
			if tag, exist := tagMap[tid]; exist && tag.State == model.TagStateNormal {
				t := &model.Tag{
					ID:           tag.Id,
					Name:         tag.Name,
					Cover:        tag.Cover,
					Content:      tag.Content,
					ShortContent: tag.ShortContent,
					HeadCover:    tag.HeadCover,
					Type:         int8(tag.Type),
					State:        int8(tag.State),
					CTime:        tag.Ctime,
					MTime:        tag.Mtime,
					IsAtten:      int8(tag.Attention),
					Attribute:    int8(tag.Attr),
				}
				t.Count.Use = int(tag.Bind)
				t.Count.Atten = int(tag.Sub)
				arcTags = append(arcTags, t)
			}
		}
		res[oid] = arcTags
	}
	return
}

// Logs .
func (s *Service) Logs(c context.Context, aid, mid int64, pn, ps int) (atls []*model.ArcTagLog, err error) {
	var (
		mids []int64
		mm   map[int64]*account.Card
		tagm = make(map[int64]string)
	)
	if _, err = s.normalArchive(c, aid); err != nil {
		atls = []*model.ArcTagLog{}
		return
	}
	var rls []*rpcModel.ResourceLog
	rls, err = s.resTagLog(c, mid, aid, rpcModel.ResTypeArchive, pn, ps)
	if err != nil {
		return
	}
	if len(rls) == 0 {
		atls = []*model.ArcTagLog{}
		return
	}
	var (
		tids []int64
		ts   []*rpcModel.Tag
	)
	for _, v := range rls {
		atl := &model.ArcTagLog{
			Lid:    v.ID,
			Aid:    v.Oid,
			Tid:    v.Tid,
			Mid:    v.Mid,
			Role:   int8(v.Role),
			Action: int8(v.Action),
			CTime:  v.CTime,
		}
		if v.State == 1 {
			atl.IsDeal = 1
		}
		if v.State == 2 {
			atl.IsReport = 1
		}
		atls = append(atls, atl)
		tids = append(tids, v.Tid)
	}
	ts, err = s.tags(c, tids, mid)
	if err != nil {
		return
	}
	for _, v := range ts {
		if v != nil {
			tagm[v.ID] = v.Name
		}
	}
	for _, atl := range atls {
		if len(tagm) > 0 {
			if name, ok := tagm[atl.Tid]; ok {
				atl.Tname = name
			}
		}
		if atl.Role != 2 {
			mids = append(mids, atl.Mid)
		}
	}
	if mm, err = s.dao.UserCards(c, mids); err != nil {
		return
	}
	for _, atl := range atls {
		if atl.Role == 2 {
			atl.Face = _adminFace
			atl.UName = _adminName
			continue
		}
		if _, ok := mm[atl.Mid]; ok {
			atl.Face = mm[atl.Mid].Face
			atl.UName = mm[atl.Mid].Name
		}
	}
	return
}

// Add .
func (s *Service) Add(c context.Context, aid, mid int64, name string, now time.Time) (tid int64, err error) {
	var arc *api.Arc
	if arc, err = s.normalArchive(c, aid); err != nil {
		return
	}
	if err = s.policy(c, arc.Author.Mid, mid, now); err != nil {
		return
	}
	if mid != arc.Author.Mid {
		err = ecode.TagOnlyUpAdd
		return
	}
	if err = s.realName(c, mid); err != nil {
		return
	}
	if err = s.filter(c, name, now); err != nil {
		return
	}
	var (
		atm                  map[int64]*model.ArcTag
		tag                  *model.Tag
		tagType, userRole, _ = s.checkTypeRole(mid, arc.Author.Mid) // get tag type and user role
	)
	if tag, err = s.resCheckTag(c, mid, name, tagType, now); err != nil {
		return
	}
	// 活动tag不可外部加
	if tag.Type == model.OfficailActiveTag {
		err = ecode.TagIsOfficailTag
		return
	}
	tid = tag.ID
	if atm, err = s.checkArcTag(c, aid, tid, model.ArcTagAdd); err != nil {
		return
	}
	count, err := s.dao.SpamCache(c, mid, model.SpamAdd)
	if err != nil {
		return
	}
	if count >= s.c.Tag.ArcTagAddMaxNum {
		err = ecode.TagArcTagAddMaxFre
		return
	}
	if err = s.checkUser(c, aid, arc.Author.Mid, mid, tid, now, model.ArcTagAdd); err != nil {
		return
	}
	var rem map[string]*rpcModel.ResourceLog
	rem, err = s.resTagLogMap(c, mid, aid, rpcModel.ResTypeArchive, 1, 20)
	if err != nil {
		return
	}
	if len(rem) > 0 {
		k := fmt.Sprintf("%d_%d_%d_%d_%d", aid, rpcModel.ResTypeArchive, tid, mid, rpcModel.ResTagLogAdd)
		if rl, ok := rem[k]; ok && rl.State == 1 {
			return 0, ecode.TagAddNotRptPassed
		}
	}
	err = s.platformUserBind(c, aid, mid, tag.ID, rpcModel.ResTypeArchive, int32(userRole), rpcModel.ResTagLogAdd)
	if err != nil {
		return
	}
	s.dao.IncrSpamCache(c, mid, model.SpamAdd)
	var tids []int64
	for _, at := range atm {
		tids = append(tids, at.Tid)
	}
	tids = append(tids, tid)
	s.setArchiveTag(c, aid, tids)
	return
}

// Del .
func (s *Service) Del(c context.Context, aid, mid, tid int64, now time.Time) (err error) {
	if _, ok := s.channelMap[tid]; ok {
		return ecode.ChannelNotAllowDel
	}
	var arc *api.Arc
	if arc, err = s.normalArchive(c, aid); err != nil {
		return
	}
	if err = s.policy(c, arc.Author.Mid, mid, now); err != nil {
		return
	}
	if mid != arc.Author.Mid {
		return ecode.TagOnlyUpDel
	}
	if err = s.realName(c, mid); err != nil {
		return
	}
	count, err := s.dao.SpamCache(c, mid, model.SpamDel)
	if err != nil {
		return
	}
	if count >= s.c.Tag.ArcTagDelMaxNum {
		return ecode.TagArcTagDelMaxFre
	}
	if err = s.checkUser(c, aid, arc.Author.Mid, mid, tid, now, model.ArcTagDel); err != nil {
		return
	}
	var (
		_, _, opRole = s.checkTypeRole(mid, arc.Author.Mid)
		atm          map[int64]*model.ArcTag
	)
	if atm, err = s.checkArcTag(c, aid, tid, model.ArcTagDel); err != nil {
		return
	}
	var rem map[string]*rpcModel.ResourceLog
	rem, err = s.resTagLogMap(c, mid, aid, rpcModel.ResTypeArchive, 1, 20)
	if err != nil {
		return
	}
	if len(rem) > 0 {
		k := fmt.Sprintf("%d_%d_%d_%d_%d", aid, rpcModel.ResTypeArchive, tid, mid, rpcModel.ResTagLogDel)
		if rl, ok := rem[k]; ok && rl.State == 1 {
			return ecode.TagDelNotRptPassed
		}
	}
	err = s.platformUserBind(c, aid, mid, tid, rpcModel.ResTypeArchive, int32(opRole), rpcModel.ResTagLogDel)
	if err != nil {
		return
	}
	s.dao.IncrSpamCache(c, mid, model.SpamDel)
	var tids []int64
	delete(atm, tid)
	for _, at := range atm {
		tids = append(tids, at.Tid)
	}
	s.setArchiveTag(c, aid, tids)
	return
}

// AddReport .
func (s *Service) AddReport(c context.Context, aid, tid, rptMid int64, reason int8, now time.Time) (err error) {
	if err = s.checkReportUser(c, rptMid); err != nil {
		return
	}
	var arc *api.Arc
	if arc, err = s.normalArchive(c, aid); err != nil {
		return
	}
	score, _ := s.userFigure(c, rptMid)
	var rtm map[int64]*rpcModel.Resource
	if rtm, err = s.resTagMap(c, rptMid, aid, rpcModel.ResTypeArchive); err != nil {
		return
	}
	if _, ok := rtm[tid]; !ok {
		if k, ok := s.channelMap[tid]; ok {
			appealInfo := &model.WorkflowAppealInfo{
				Business: model.WorkflowBusinessChannel,
				FID:      model.WorkflowFIDChannel,
				RID:      model.WorkflowRIDChannel,
				Score:    score,
				Oid:      aid,
				EID:      tid,
				RptMid:   rptMid,
				Mid:      arc.Author.Mid,
				RegionID: arc.TypeID,
				ReasonID: reason,
				TName:    k.Name,
				RealIP:   metadata.String(c, metadata.RemoteIP),
			}
			return s.workflowAppeal(c, appealInfo)
		}
		return ecode.TagResTagNotExist
	}
	err = s.dao.AddReport(c, aid, tid, rptMid, rpcModel.ResTypeArchive, arc.TypeID, int32(reason), int32(score))
	if err == nil {
		s.cacheCh.Do(c, func(ctx context.Context) {
			s.dao.IncrSpamCache(ctx, rptMid, model.SpamReport)
		})
	}
	return
}

// LogReport .
func (s *Service) LogReport(c context.Context, aid, logID, rptMid int64, reason int8, now time.Time) (err error) {
	if err = s.checkReportUser(c, rptMid); err != nil {
		return
	}
	var arc *api.Arc
	if arc, err = s.normalArchive(c, aid); err != nil {
		return
	}
	score, _ := s.userFigure(c, rptMid)
	err = s.reportAction(c, aid, logID, rptMid, rpcModel.ResTypeArchive, arc.TypeID, int32(reason), int32(score))
	if err == nil {
		s.cacheCh.Do(c, func(ctx context.Context) {
			s.dao.IncrSpamCache(ctx, rptMid, model.SpamReport)
		})
	}
	return
}

func (s *Service) setArchiveTag(c context.Context, aid int64, tids []int64) (err error) {
	var (
		names   []string
		nameStr string
		tags    []*model.Tag
	)
	if tags, err = s.infos(c, 0, tids); err != nil {
		return
	}
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	nameStr = strings.Join(names, ",")
	// up archive tag
	s.setArcTag(c, aid, nameStr)
	return
}

func (s *Service) workflowAppeal(c context.Context, appeal *model.WorkflowAppealInfo) (err error) {
	var state int32
	if state, err = s.workflowAppealState(c, appeal); err != nil {
		return
	}
	if state == model.AppealStatEffective || state == model.AppealStateInvalid {
		return ecode.TagRptNotRptPassed
	}
	appeals, err := s.workflowAppeaList(c, appeal)
	if err != nil {
		return
	}
	for _, rpt := range appeals {
		if rpt.Business == model.WorkflowBusinessChannel && rpt.Oid == appeal.Oid && rpt.Eid == appeal.EID && rpt.Mid == appeal.RptMid {
			return ecode.TagArcTagRpted
		}
	}
	return s.addWorkflowAppeal(c, appeal)
}

func (s *Service) checkReportUser(c context.Context, mid int64) (err error) {
	card, err := s.dao.UserCard(c, mid)
	if err != nil {
		return
	}
	if card.Level < int32(s.c.Tag.ArcTagRptLevel) { // 等级限制
		return ecode.TagArcRptLevelLower
	}
	if card.Silence != model.UserBannedNone {
		return ecode.TagArcAccountBlocked
	}
	count, err := s.dao.SpamCache(c, mid, model.SpamReport)
	if err != nil {
		return
	}
	if count >= s.c.Tag.ArcTagRptMaxNum {
		err = ecode.TagArcTagRptMaxFre
	}
	return
}
