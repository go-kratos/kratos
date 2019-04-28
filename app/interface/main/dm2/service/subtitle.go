package service

import (
	"context"
	"time"

	"go-common/app/interface/main/dm2/model"
	account "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	archive "go-common/app/service/main/archive/model/archive"
	filterCli "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_filterArea = "subtitle"
)

// SubtitleFilter .
func (s *Service) SubtitleFilter(c context.Context, words string) (hit []string, err error) {
	var (
		reply *filterCli.HitReply
	)
	if len(words) > model.SubtitleContentSizeLimit {
		err = ecode.SubtitleSizeLimit
		return
	}
	if reply, err = s.filterRPC.Hit(c, &filterCli.HitReq{
		Area: _filterArea,
		Msg:  words,
	}); err != nil {
		log.Error("SubtitleFilter(params:%+v),error(%v)", words, err)
		return
	}
	hit = reply.Hits
	return
}

// SubtitlePermission .
func (s *Service) SubtitlePermission(c context.Context, aid int64, oid int64, tp int32, mid int64) (err error) {
	var (
		subject *model.Subject
	)
	if subject, err = s.subject(c, tp, oid); err != nil {
		log.Error("params(tp:%v,oid:%v).error(%v)", tp, oid, err)
		return
	}
	if s.checkAidOid(c, aid, oid); err != nil {
		return
	}
	if err = s.checkSubtitlePermission(c, aid, oid, tp, mid, subject); err != nil {
		return
	}
	return
}

// DelSubtitle .
func (s *Service) DelSubtitle(c context.Context, oid int64, subtitleID int64, mid int64) (err error) {
	var (
		subtitle *model.Subtitle
	)
	if subtitle, err = s.getSubtitle(c, oid, subtitleID); err != nil {
		log.Error("params(oid:%v, subtitleID:%v) error(%v)", oid, subtitleID, err)
		return
	}
	if subtitle == nil {
		err = ecode.NothingFound
		return
	}
	if subtitle.Mid != mid {
		err = ecode.SubtitlePermissionDenied
		return
	}
	if subtitle.Status != model.SubtitleStatusDraft &&
		subtitle.Status != model.SubtitleStatusToAudit &&
		subtitle.Status != model.SubtitleStatusAuditBack &&
		subtitle.Status != model.SubtitleStatusManagerBack {
		err = ecode.SubtitleDelUnExist
		return
	}
	subtitle.Status = model.SubtitleStatusRemove
	subtitle.PubTime = time.Now().Unix()
	if err = s.updateSubtitle(c, subtitle); err != nil {
		return
	}
	s.subtitleReportDelete(c, oid, subtitleID)
	return
}

// addSubtitle new a subtitle draft
func (s *Service) addSubtitle(c context.Context, draft *model.Subtitle) (insertID int64, err error) {
	if insertID, err = s.dao.AddSubtitle(c, draft); err != nil {
		log.Error("params(draft:%+v).error(%v)", draft, err)
		return
	}
	s.dao.DelSubtitleDraftCache(context.Background(), draft.Oid, draft.Type, draft.Mid, draft.Lan)
	s.dao.DelSubtitleCache(context.Background(), draft.Oid, draft.ID)
	return
}

// updateSubtitle  update an exist subtitle
func (s *Service) updateSubtitle(c context.Context, subtitle *model.Subtitle) (err error) {
	if err = s.dao.UpdateSubtitle(c, subtitle); err != nil {
		log.Error("params(draft:%+v).error(%v)", subtitle, err)
		return
	}
	s.dao.DelSubtitleDraftCache(context.Background(), subtitle.Oid, subtitle.Type, subtitle.Mid, subtitle.Lan)
	s.dao.DelSubtitleCache(context.Background(), subtitle.Oid, subtitle.ID)
	return
}

// SubtitleSign .
func (s *Service) SubtitleSign(c context.Context, oid int64, tp int32, mid int64, subtitleID int64, isSign bool) (err error) {
	var (
		subtitle *model.Subtitle
	)
	if subtitle, err = s.getSubtitle(c, oid, subtitleID); err != nil {
		log.Error("params(oid:%v,subtitleID:%v).error(%v)", oid, subtitleID, err)
		return
	}
	if subtitle == nil {
		err = ecode.NothingFound
		return
	}
	if mid != subtitle.Mid {
		err = ecode.SubtitlePermissionDenied
		return
	}
	if subtitle.Status != model.SubtitleStatusDraft &&
		subtitle.Status != model.SubtitleStatusToAudit &&
		subtitle.Status != model.SubtitleStatusAuditBack &&
		subtitle.Status != model.SubtitleStatusPublish &&
		subtitle.Status != model.SubtitleStatusCheckToAudit &&
		subtitle.Status != model.SubtitleStatusCheckPublish &&
		subtitle.Status != model.SubtitleStatusManagerBack {
		err = ecode.SubtitlePermissionDenied
		return
	}
	subtitle.IsSign = isSign
	if err = s.dao.UpdateSubtitle(c, subtitle); err != nil {
		log.Error("params(%+v).error(%v)", subtitle, err)
		return
	}
	if err = s.dao.DelSubtitleCache(c, oid, subtitleID); err != nil {
		log.Error("DelSubtitleCache.params(oid:%v,subtitleID:%v).error(%v)", oid, subtitleID, err)
		return
	}
	if err = s.dao.DelVideoSubtitleCache(c, oid, tp); err != nil {
		log.Error("DelVideoSubtitleCache.params(oid:%v,tp:%v).error(%v)", oid, tp, err)
		return
	}
	return
}

// SubtitleLock .
func (s *Service) SubtitleLock(c context.Context, oid int64, tp int32, mid int64, subtitleID int64, isLock bool) (err error) {
	var (
		subject  *model.Subject
		subtitle *model.Subtitle
	)
	if subtitle, err = s.getSubtitle(c, oid, subtitleID); err != nil {
		log.Error("params(oid:%v,subtitleID:%v).error(%v)", oid, subtitleID, err)
		return
	}
	if subtitle == nil {
		err = ecode.NothingFound
		return
	}
	if subject, err = s.subject(c, tp, oid); err != nil {
		log.Error("params(oid:%v,tp:%v).error(%v)", oid, tp, err)
		return
	}
	if mid != subject.Mid {
		err = ecode.SubtitlePermissionDenied
		return
	}
	if subtitle.Status != model.SubtitleStatusPublish &&
		subtitle.Status != model.SubtitleStatusCheckPublish {
		err = ecode.SubtitleNotPublish
		return
	}
	subtitle.IsLock = isLock
	if err = s.dao.UpdateSubtitle(c, subtitle); err != nil {
		log.Error("params(%+v).error(%v)", subtitle, err)
		return
	}
	if err = s.dao.DelSubtitleCache(c, oid, subtitleID); err != nil {
		log.Error("DelSubtitleCache.params(oid:%v,subtitleID:%v).error(%v)", oid, subtitleID, err)
		return
	}
	if err = s.dao.DelVideoSubtitleCache(c, oid, tp); err != nil {
		log.Error("DelVideoSubtitleCache.params(oid:%v,tp:%v).error(%v)", oid, tp, err)
		return
	}
	return
}

// ArchiveName .
func (s *Service) ArchiveName(c context.Context, aid int64) (arcvhiveName string, err error) {
	var (
		res *api.Arc
	)
	if res, err = s.arcRPC.Archive3(c, &archive.ArgAid2{
		Aid: aid,
	}); err != nil {
		log.Error("params(aid:%v).error(%v)", aid, err)
		return
	}
	arcvhiveName = res.Title
	return
}

// SubtitleShow .
func (s *Service) SubtitleShow(c context.Context, oid int64, subtitleID int64, mid int64) (subtitleShow *model.SubtitleShow, err error) {
	var (
		subtitle   *model.Subtitle
		canShow    bool
		res        *api.Arc
		infoReply  *account.InfoReply
		showStatus model.SubtitleStatus
	)
	if subtitle, err = s.getSubtitle(c, oid, subtitleID); err != nil {
		log.Error("params(oid:%v,subtitleID:%v).error(%v)", oid, subtitleID, err)
		return
	}
	if subtitle == nil {
		err = ecode.NothingFound
		return
	}
	showStatus = subtitle.Status
	// 发布的状态都可见
	// 非发布的状态本人可见
	// 审核状态 up  主可见
	switch subtitle.Status {
	case model.SubtitleStatusPublish:
		canShow = true
	case model.SubtitleStatusToAudit,
		model.SubtitleStatusCheckPublish:
		if subtitle.UpMid == mid || subtitle.Mid == mid {
			canShow = true
		}
	case model.SubtitleStatusDraft,
		model.SubtitleStatusAuditBack,
		model.SubtitleStatusCheckToAudit:
		if subtitle.Mid == mid {
			canShow = true
		}
	case model.SubtitleStatusManagerBack:
		if subtitle.Mid == mid {
			canShow = true
		}
		showStatus = model.SubtitleStatusAuditBack
	default:
		err = ecode.SubtitlePermissionDenied
		return
	}
	if !canShow {
		err = ecode.SubtitlePermissionDenied
		return
	}
	lan, lanDoc := s.subtitleLans.GetByID(int64(subtitle.Lan))
	subtitleShow = &model.SubtitleShow{
		ID:            subtitle.ID,
		Oid:           subtitle.Oid,
		Type:          subtitle.Type,
		Aid:           subtitle.Aid,
		Lan:           lan,
		LanDoc:        lanDoc,
		Mid:           subtitle.Mid,
		IsSign:        subtitle.IsSign,
		IsLock:        subtitle.IsLock,
		Status:        showStatus,
		SubtitleURL:   subtitle.SubtitleURL,
		RejectComment: subtitle.RejectComment,
	}
	if subtitle.UpMid == mid {
		subtitleShow.UpperStatus = model.UpperStatusUpper
	}
	if subtitle.Mid == mid {
		subtitleShow.AuthorStatus = model.AuthorStatusAuthor
	}
	if res, err = s.arcRPC.Archive3(c, &archive.ArgAid2{
		Aid: subtitle.Aid,
	}); err != nil {
		log.Error("params(aid:%v).error(%v)", subtitle.Aid, err)
		err = nil
	} else {
		subtitleShow.ArchiveName = res.Title
	}
	if infoReply, err = s.accountRPC.Info3(c, &account.MidReq{
		Mid: subtitle.AuthorID,
	}); err != nil {
		log.Error("params(mid:%v).error(%v)", subtitle.Mid, err)
		err = nil
	} else {
		subtitleShow.Author = infoReply.GetInfo().GetName()
	}
	return
}
