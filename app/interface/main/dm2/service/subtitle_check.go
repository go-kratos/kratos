package service

import (
	"context"
	"encoding/json"

	"go-common/app/interface/main/dm2/model"
	account "go-common/app/service/main/account/api"
	archive "go-common/app/service/main/archive/model/archive"
	memberMdl "go-common/app/service/main/member/model/block"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_subtitleMaxLength = 50 * 10000
)

func (s *Service) checkAidOid(c context.Context, aid, oid int64) (err error) {
	if _, err = s.arcRPC.Video3(c, &archive.ArgVideo2{
		Aid: aid,
		Cid: oid,
	}); err != nil {
		log.Error("Video3(aid:%v,oid:%v),error(%v)", aid, oid, err)
		err = ecode.RequestErr
		return
	}
	return
}

func (s *Service) checkSubtitleLan(c context.Context, lan string) (err error) {
	if s.subtitleLans.GetByLan(lan) <= 0 {
		err = ecode.SubtitleIllegalLanguage
		return
	}
	return
}

func (s *Service) checkSubtitlePermission(c context.Context, aid, oid int64, tp int32, mid int64, subject *model.Subject) (err error) {
	if err = s.checkAudienceDraftAdd(c, aid, oid, tp, mid, subject); err != nil {
		log.Error("params(oid:%v, tp:%v, mid:%v,subject:%+v).err(%v)", oid, tp, mid, subject, err)
		return
	}
	return
}

func (s *Service) checkSubtitleData(c context.Context, aid, oid int64, data []byte) (detectErrs []*model.SubtitleDetectError, err error) {
	var (
		subtitleBody *model.SubtitleBody
		duration     int64
	)
	if len(data) > _subtitleMaxLength {
		err = ecode.SubtitleSizeLimit
		return
	}
	if err = json.Unmarshal(data, &subtitleBody); err != nil {
		err = ecode.SubtitleUnValid
		return
	}
	if duration, err = s.videoDuration(c, aid, oid); err != nil {
		return
	}
	if duration <= 0 {
		err = ecode.SubtitleTimeUnValid
		return
	}
	if detectErrs, err = subtitleBody.CheckItem(duration); err != nil {
		return
	}
	return
}

func (s *Service) checkSubtitleLocked(c context.Context, submit bool, oid int64, tp int32, lan string, mid int64) (err error) {
	var (
		lockSubtitle *model.Subtitle
	)
	if !submit {
		return
	}
	if lockSubtitle, err = s.isSubtitleLanLock(c, oid, tp, lan); err != nil {
		log.Error("params(oid:%v, tp:%v) error(%v)", oid, tp, err)
		return
	}
	if lockSubtitle != nil && lockSubtitle.IsLock && lockSubtitle.UpMid != mid && lockSubtitle.Mid != mid {
		err = ecode.SubtileLanLocked
		return
	}
	return
}

func (s *Service) checkSubtitleAuthor(c context.Context, oid, subtitleID int64, lan string, mid int64) (authorID int64, err error) {
	var (
		originSubtitle *model.Subtitle
		originLan      int64
	)
	if subtitleID <= 0 {
		authorID = mid
		return
	}
	if originSubtitle, err = s.getSubtitle(c, oid, subtitleID); err != nil {
		return
	}
	if originSubtitle == nil {
		err = ecode.NothingFound
		return
	}
	if originSubtitle.Status != model.SubtitleStatusAuditBack &&
		originSubtitle.Status != model.SubtitleStatusPublish &&
		originSubtitle.Status != model.SubtitleStatusCheckPublish &&
		originSubtitle.Status != model.SubtitleStatusManagerBack {
		err = ecode.SubtitleOriginUnValid
		return
	}
	if originLan = s.subtitleLans.GetByLan(lan); originLan <= 0 || originLan != int64(originSubtitle.Lan) {
		err = ecode.SubtitleIllegalLanguage
		return
	}
	authorID = originSubtitle.AuthorID
	return
}

func (s *Service) checkAudienceDraftAdd(c context.Context, aid, oid int64, tp int32, mid int64, subject *model.Subject) (err error) {
	var (
		profileReply    *account.ProfileReply
		blackReply      *account.BlacksReply
		blockInfo       *memberMdl.RPCResInfo
		resDm           []*model.UpFilter
		allow           bool
		closed          bool
		subtitleSubject *model.SubtitleSubject
	)
	if subtitleSubject, err = s.subtitleSubject(c, aid); err != nil {
		log.Error("subtitleSubject(aid:%v) error(%v)", aid, err)
		err = nil
	}
	if subtitleSubject != nil {
		allow = subtitleSubject.Allow
		closed = subtitleSubject.AttrVal(model.AttrSubtitleClose) == model.AttrYes
	}
	if closed {
		err = ecode.SubtitleDenied
		return
	}
	if subject.Mid == mid {
		return
	}
	if !allow {
		err = ecode.SubtitleDenied
		return
	}
	// 视频观众可以投稿
	// 账号绑定手机号
	if profileReply, err = s.accountRPC.Profile3(c, &account.MidReq{
		Mid: mid,
	}); err != nil {
		log.Error("accRPC.UserInfo(%v) error(%v)", mid, err)
		return
	}
	if profileReply.GetProfile().GetIdentification() == 0 && profileReply.GetProfile().GetTelStatus() == 0 {
		err = ecode.UserCheckNoPhone
		return
	}
	if profileReply.GetProfile().GetIdentification() == 0 && profileReply.GetProfile().GetTelStatus() == 2 {
		err = ecode.UserCheckInvalidPhone
		return
	}
	if profileReply.GetProfile().GetTelStatus() == 0 {
		err = ecode.UserCheckInvalidPhone
		return
	}
	// 用户等级大于2
	if profileReply.GetProfile().GetLevel() < 2 {
		err = ecode.UserLevelLow
		return
	}
	// 账号被拉黑
	if blackReply, err = s.accountRPC.Blacks3(c, &account.MidReq{
		Mid: subject.Mid,
	}); err != nil {
		log.Error("params(arg:%+v).err(%v)", subject.Mid, err)
		return
	}
	if _, ok := blackReply.GetBlackList()[mid]; ok {
		err = ecode.SubtitleUserBalcked
		return
	}
	if resDm, err = s.UpFilters(c, subject.Mid); err != nil {
		log.Error("params(mid:%+v).err(%v)", subject.Mid, err)
		return
	}
	hash := model.Hash(mid, 0)
	for _, uf := range resDm {
		if uf.Filter == hash {
			err = ecode.SubtitleUserBalcked
			return
		}
	}
	// 账号被封禁
	if blockInfo, err = s.memberRPC.BlockInfo(c, &memberMdl.RPCArgInfo{
		MID: mid,
	}); err != nil {
		log.Error("params(arg:%+v).err(%v)", mid, err)
		return
	}
	if blockInfo.BlockStatus != memberMdl.BlockStatusFalse {
		err = ecode.UserDisabled
		return
	}
	return
}
