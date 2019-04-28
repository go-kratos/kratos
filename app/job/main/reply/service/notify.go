package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	artmdl "go-common/app/interface/openplatform/article/model"
	model "go-common/app/job/main/reply/model/reply"
	accmdl "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	epmdl "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	"go-common/library/log"
)

const (
	_mcReply      = "1_1_1"
	_mcCntArticle = "1_1_4"
	_mcCntDynamic = "1_1_5"
	_mcCntClip    = "1_1_6"
	_mcCntAlbum   = "1_1_7"
	_mcCntArchive = "1_1_8"

	_msgTitleSize   = 40
	_msgContentSize = 80
)

func (s *Service) notifyReply(c context.Context, sub *model.Subject, rp *model.Reply) {
	s.notify.Do(c, func(c context.Context) {
		if len(rp.Content.Ats) > 0 {
			title, link, jump, nativeJump, msg := s.messageInfo(c, rp)
			if link != "" {
				atmt := fmt.Sprintf("#{%s}{\"%s\"}评论中@了你", title, link)
				cont := fmt.Sprintf("#{%s}{\"%s\"}", msg, jump)
				if err := s.messageDao.At(c, rp.Mid, rp.Content.Ats, atmt, cont, extraInfo(nativeJump), rp.CTime.Time()); err != nil {
					log.Error("s.messageDao.At failed , mid(%d) err(%v)", rp.Mid, err)
				}
			}
		}
		if err := s.notifyCnt(c, sub, rp); err != nil {
			log.Error("s.notifyCnt(%v,%v) error(%v)", sub, rp, err)
		}
	})
}

func (s *Service) notifyReplyReply(c context.Context, sub *model.Subject, rootRp, parentRp, rp *model.Reply) {
	s.notify.Do(c, func(c context.Context) {
		if err := s.notifyCnt(c, sub, rp); err != nil {
			log.Error("s.notifyCnt(%v,%v) error(%v)", sub, rp, err)

		}
		// notify parent reply
		if rp.Mid == rootRp.Mid && rp.Root == rp.Parent && len(rp.Content.Ats) == 0 {
			return
		}
		title, link, jump, nativeJump, msg := s.messageInfo(c, rp)
		if title == "" || link == "" {
			return
		}
		// 5.29 改为根评论内容推送
		//rpmt = fmt.Sprintf("#{%s}{\"%s\"}评论中回复了你", title, link)
		rpmt := []rune(parentRp.Content.Message)
		if len(rpmt) > _msgContentSize {
			rpmt = rpmt[:_msgContentSize]
		}
		atmt := fmt.Sprintf("#{%s}{\"%s\"}评论中@了你", title, link)
		cont := fmt.Sprintf("#{%s}{\"%s\"}", msg, jump)
		// notify
		if rp.Mid != rootRp.Mid && !s.getBlackListRelation(c, rootRp.Mid, rp.Mid) {
			if err := s.messageDao.Reply(c, _mcReply, "", rp.Mid, rootRp.Mid, string(rpmt), cont, extraInfo(nativeJump), rp.CTime.Time()); err != nil {
				log.Error("s.messageDao.Reply failed , mid(%d) Parent(%d),  err(%v)", rp.Mid, rootRp.Mid, err)
			}
		}
		if rp.Root != rp.Parent {
			if parentRp != nil && rootRp.Mid != parentRp.Mid && rp.Mid != parentRp.Mid && !s.getBlackListRelation(c, parentRp.Mid, rp.Mid) {
				if err := s.messageDao.Reply(c, _mcReply, "", rp.Mid, parentRp.Mid, string(rpmt), cont, extraInfo(nativeJump), rp.CTime.Time()); err != nil {
					log.Error("s.messageDao.Reply failed , mid(%d) Parent(%d),  err(%v)", rp.Mid, parentRp.Mid, err)
				}
			}
		}
		var ats []int64
		for _, mid := range rp.Content.Ats {
			if mid != parentRp.Mid {
				ats = append(ats, mid)
			}
		}
		if len(ats) > 0 {
			if err := s.messageDao.At(c, rp.Mid, ats, atmt, cont, extraInfo(nativeJump), rp.CTime.Time()); err != nil {
				log.Error("s.messageDao.At failed , mid(%d),  err(%v)", rp.Mid, err)
			}
		}
	})
}

func (s *Service) notifyLike(c context.Context, mid int64, rp *model.Reply) {
	s.notify.Do(c, func(c context.Context) {
		if ok, num := s.notifyLikeNum(c, rp, mid); ok {
			_, _, jump, nativeJump, msg := s.messageInfo(c, rp)
			if jump == "" {
				return
			}
			// NOTE content and title is opposite
			cont := fmt.Sprintf("等%d人赞了你的回复", num)
			rpmt := fmt.Sprintf("#{%s}{\"%s\"}", msg, jump)
			if err := s.messageDao.Like(c, mid, rp.Mid, rpmt, cont, extraInfo(nativeJump), rp.CTime.Time()); err != nil {
				log.Error("s.messageDao.Reply failed , mid(%d) Parent(%d),  err(%v)", mid, rp.Mid, err)
			}
		} else {
			log.Warn("Didn't satify notify condition, omit notify!")
		}
	})
}

// notifyLike check if need notify user when receive like
func (s *Service) notifyLikeNum(c context.Context, rp *model.Reply, mid int64) (ok bool, num int64) {
	if rp.Mid == mid || rp.Like <= 0 {
		ok = false
		return
	}
	num = int64(rp.Like)
	// NOTE  if  num >1000 send when num%1000==0
	if num < 10 || (num < 100 && num%10 == 0) || (num < 1000 && num%100 == 0) || num%1000 == 0 {
		ok = true
	}
	return
}

func (s *Service) notifyCnt(c context.Context, sub *model.Subject, rp *model.Reply) (err error) {
	max, err := s.dao.Redis.NotifyCnt(c, sub.Oid, sub.Type)
	if err != nil {
		log.Error("redis.NotifyCnt(%d,%d) error(%v)", sub.Oid, sub.Type, err)
		return
	}
	if sub.ACount <= max {
		log.Warn("notifyCnt ignore oid:%d type:%d current:%d max:%d", sub.Oid, sub.Type, sub.ACount, max)
		return
	}
	if err = s.dao.Redis.SetNotifyCnt(c, sub.Oid, sub.Type, sub.ACount); err != nil {
		log.Error("redis.SetNotifyCnt(%d,%d,%d) error(%v)", sub.Oid, sub.Type, sub.ACount, err)
		return
	}
	switch sub.Type {
	case model.SubTypeVideo:
		return s.notifyArchiveCnt(c, sub, rp)
	case model.SubTypeArticle:
		return s.notifyArticleCnt(c, sub, rp)
	case model.SubTypeDynamic:
		return s.notifyDynamicCnt(c, sub, rp, _mcCntDynamic)
	case model.SubTypeLiveVideo:
		return s.notifyDynamicCnt(c, sub, rp, _mcCntClip)
	case model.SubTypeLivePicture:
		return s.notifyDynamicCnt(c, sub, rp, _mcCntAlbum)
	default:
		return
	}
}
func (s *Service) notifyDynamicCnt(c context.Context, sub *model.Subject, rp *model.Reply, mc string) (err error) {
	if !shouldNotifyLow(sub.ACount) {
		return
	}
	title, link, _, nativeJump, msg := s.messageInfo(c, rp)
	if title == "" || link == "" {
		return
	}
	resID := fmt.Sprintf("%d_%d", rp.Oid, rp.Type)
	notifyTitle := fmt.Sprintf("#{%s}{\"%s\"}收到了第%d条评论", title, link, sub.ACount)
	notifyContent := fmt.Sprintf("#{%s}{\"%s\"}", msg, link)
	if err = s.messageDao.Reply(c, mc, resID, rp.Mid, sub.Mid, notifyTitle, notifyContent, extraInfo(nativeJump), rp.CTime.Time()); err != nil {
		log.Error("s.messageDao.Reply(mid:%d,oid:%d,type:%d,acount:%d) error(%v)", rp.Mid, rp.Oid, rp.Type, sub.ACount, err)
	}
	return
}

func (s *Service) notifyArchiveCnt(c context.Context, sub *model.Subject, rp *model.Reply) (err error) {
	if !shouldNotifyLow(sub.ACount) {
		return
	}
	title, link, _, nativeJump, msg := s.messageInfo(c, rp)
	if title == "" || link == "" {
		return
	}
	resID := fmt.Sprintf("%d_%d", rp.Oid, rp.Type)
	notifyTitle := fmt.Sprintf("你的投稿收到了第%d条评论", sub.ACount)
	notifyContent := fmt.Sprintf("你投稿的视频“#{%s}{\"%s\"}”收到了第%d条评论：『%s』", title, link, sub.ACount, msg)
	if err = s.messageDao.System(c, _mcCntArchive, resID, sub.Mid, notifyTitle, notifyContent, extraInfo(nativeJump), rp.CTime.Time()); err != nil {
		log.Error("s.messageDao.System(mid:%d,oid:%d,type:%d,acount:%d) error(%v)", rp.Mid, rp.Oid, rp.Type, sub.ACount, err)
	}
	return
}

func (s *Service) notifyArticleCnt(c context.Context, sub *model.Subject, rp *model.Reply) (err error) {
	if !shouldNotifyMiddle(sub.ACount) {
		return
	}
	title, link, _, nativeJump, _ := s.messageInfo(c, rp)
	if title == "" || link == "" {
		return
	}
	resID := fmt.Sprintf("%d_%d", rp.Oid, rp.Type)
	notifyTitle := fmt.Sprintf("你的专栏文章评论数达到了%d", sub.ACount)
	notifyContent := fmt.Sprintf("你投稿的专栏文章“#{%s}{\"%s\"}”评论数达到了%d！去回应一下大家的评论吧～  #{点击前往}{\"%s\"}", title, link, sub.ACount, link)
	if err = s.messageDao.System(c, _mcCntArticle, resID, sub.Mid, notifyTitle, notifyContent, extraInfo(nativeJump), rp.CTime.Time()); err != nil {
		log.Error("s.messageDao.System(mid:%d,oid:%d,type:%d,acount:%d) error(%v)", rp.Mid, rp.Oid, rp.Type, sub.ACount, err)
	}
	return
}

func shouldNotifyLow(n int) (ok bool) {
	switch {
	case n <= 0:
		ok = false
	case n == 1 || n == 10 || n == 30 || n == 50:
		ok = true
	case n <= 1000:
		ok = (n%100 == 0)
	default:
		ok = (n%10000 == 0)
	}
	return
}

func shouldNotifyMiddle(n int) (ok bool) {
	switch {
	case n <= 0:
		ok = false
	case n <= 10:
		ok = true
	case n <= 100:
		ok = (n%10 == 0)
	case n <= 1000:
		ok = (n%100 == 0)
	default:
		ok = (n%10000 == 0)
	}
	return
}

// filterViolationMsg every two characters, the third character processing for *.
func filterViolationMsg(msg string) string {
	s := []rune(msg)
	for i := 0; i < len(s); i++ {
		if i%3 != 0 {
			s[i] = '*'
		}
	}
	return string(s)
}

// moralAndNotify del moral and notify user.
func (s *Service) moralAndNotify(c context.Context, rp *model.Reply, moral int, notify bool, rptMid, adid int64, adname, remark string, reason, freason int8, ftime int64, isPunish bool) (err error) {
	title, link, _, _, msg := s.messageInfo(c, rp)
	smsg := []rune(msg)
	if len(smsg) > 50 {
		smsg = smsg[:50]
	}
	if moral > 0 {
		reason := "发布的评论违规并被管理员删除 - " + string(smsg)
		if rptMid > 0 {
			reason = "发布的评论被举报并被管理员删除 - " + string(smsg)
		}
		arg := &accmdl.MoralReq{
			Mid:    rp.Mid,
			Moral:  -float64(moral),
			Oper:   adname,
			Reason: reason,
			Remark: remark,
		}
		if _, err = s.accSrv.AddMoral3(c, arg); err != nil {
			log.Error("s.accSrv.AddMoral3(%d) error(%v)", rp.Mid, err)
		}
	}
	msg = filterViolationMsg(msg)
	if title != "" && link != "" && rptMid > 0 {
		if err = s.reportNotify(c, rp, title, link, msg, ftime, reason, freason, isPunish); err != nil {
			log.Error("s.reportNotify(%d,%d,%d) error(%v)", rp.Oid, rp.Type, rp.RpID, err)
		}
	}
	if !notify {
		return
	}
	if title != "" && link != "" {
		// notify message
		mt := "评论违规处理通知"
		mc := fmt.Sprintf("您好，您在#{%s}{\"%s\"}下的评论 『%s』 ", title, link, msg)
		if rptMid > 0 {
			mc = fmt.Sprintf("您好，根据用户举报，您在#{%s}{\"%s\"}下的评论 『%s』 ", title, link, msg)
		}
		if isPunish {
			mc += "，已被处罚"
		} else {
			mc += "，已被移除"
		}
		// forbidden
		if ftime > 0 {
			mc += fmt.Sprintf("，并被封禁%d天。", ftime)
		} else if ftime == -1 {
			mc += "，并被永久封禁。"
		} else {
			mc += "。"
		}
		// forbid reason
		if ar, ok := model.ForbidReason[freason]; ok {
			mc += "理由：" + ar + "。"
			// community rules
			switch {
			case freason == model.ForbidReasonSpoiler || freason == model.ForbidReasonAd || freason == model.ForbidReasonUnlimitedSign || freason == model.ForbidReasonMeaningless:
				mc += model.NotifyComRules
			case freason == model.ForbidReasonProvoke || freason == model.ForbidReasonAttack:
				mc += model.NotifyComProvoke
			default:
				mc += model.NofityComProhibited
			}
		} else { // report reason
			if ar, ok := model.ReportReason[reason]; ok {
				mc += "理由：" + ar + "。"
			}
			// community rules
			switch {
			case reason == model.ReportReasonSpoiler || reason == model.ReportReasonAd || reason == model.ReportReasonUnlimitedSign || reason == model.ReportReasonMeaningless:
				mc += model.NotifyComRules
			case reason == model.ReportReasonUnrelated:
				mc += model.NotifyComUnrelated
			case reason == model.ReportReasonProvoke || reason == model.ReportReasonAttack:
				mc += model.NotifyComProvoke
			default:
				mc += model.NofityComProhibited
			}
		}
		// send the message
		if err = s.messageDao.DeleteReply(c, rp.Mid, mt, mc, rp.MTime.Time()); err != nil {
			log.Error("s.messageDao.DeleteReply failed, (%d) error(%v)", rp.Mid, err)
		}
		log.Info("notify oid:%d type:%d rpID:%d reason:%d content:%s", rp.Oid, rp.Type, rp.RpID, reason, mc)
	} else {
		log.Warn("no notify oid:%d type:%d rpid:%d", rp.Oid, rp.Type, rp.RpID)
	}
	return
}

func (s *Service) reportNotify(c context.Context, rp *model.Reply, title, link, msg string, ftime int64, reason, freason int8, isPunish bool) (err error) {
	var (
		rptUser  *model.ReportUser
		rptUsers map[int64]*model.ReportUser
	)
	mt := "举报处理结果通知"
	mc := fmt.Sprintf("您好，您在#{%s}{\"%s\"}下举报的评论 『%s』 ", title, link, msg)
	if isPunish {
		mc += "已被处罚"
	} else {
		mc += "已被移除"
	}
	// forbidden
	if ftime > 0 {
		mc += fmt.Sprintf("，并被封禁%d天。", ftime)
	} else if ftime == -1 {
		mc += "，该用户已被永久封禁。"
	} else {
		mc += "。"
	}
	// forbid reason
	if ar, ok := model.ForbidReason[freason]; ok {
		mc += "理由：" + ar + "。"
	} else if ar, ok := model.ReportReason[reason]; ok {
		mc += "理由：" + ar + "。"
	}
	// community rules
	mc += model.NotifyComRulesReport
	if rptUsers, err = s.dao.Report.GetUsers(c, rp.Oid, rp.Type, rp.RpID); err != nil {
		log.Error("reportUser.GetUsers(%d,%d,%d) error(%v)", rp.Oid, rp.Type, rp.RpID, err)
		return
	}
	for _, rptUser = range rptUsers {
		// send the message
		if err = s.messageDao.AcceptReport(c, rptUser.Mid, mt, mc, rp.MTime.Time()); err != nil {
			log.Error("s.messageDao.DeleteReply failed, (%d) error(%v)", rp.Mid, err)
		}
	}
	if _, err = s.dao.Report.SetUserReported(c, rp.Oid, rp.Type, rp.RpID, rp.MTime.Time()); err != nil {
		log.Error("s.dao.Report.SetUserReported(%d, %d, %d) error(%v)", rp.Oid, rp.Type, rp.RpID)
	}
	return
}

func (s *Service) messageInfo(c context.Context, rp *model.Reply) (title, link, jump, nativeJump, msg string) {
	var (
		err           error
		native        bool
		subType       int
		extraIntentID int64
	)
	switch rp.Type {
	case model.SubTypeVideo:
		var (
			m   *api.Arc
			uri *url.URL
		)
		arg := &arcmdl.ArgAid2{
			Aid: rp.Oid,
		}
		m, err = s.arcSrv.Archive3(c, arg)
		if err != nil || m == nil {
			log.Error("s.arcSrv.Archive3(%v) ret:%v error(%v)", arg, m, err)
			return
		}
		if m.AttrVal(arcmdl.AttrBitIsBangumi) == 1 {
			req := &epmdl.EpAidReq{
				Aids: []int32{int32(rp.Oid)},
			}
			resp, err1 := s.bangumiSrv.ListByAids(c, req)
			if err1 != nil {
				log.Error("s.bangumiSrv.ListByAids(%v, %v) error(%v)", c, req, err1)
				return
			}
			if resp.Infos[int32(rp.Oid)] != nil {
				extraIntentID = int64(resp.Infos[int32(rp.Oid)].EpisodeId)
			}
			subType = 1
		}
		if m.RedirectURL != "" {
			// NOTE mobile jump
			if uri, err = url.Parse(m.RedirectURL); err == nil {
				q := uri.Query()
				q.Set("aid", strconv.FormatInt(rp.Oid, 10))
				uri.RawQuery = q.Encode()
				link = uri.String()
			}
		} else {
			link = fmt.Sprintf("http://www.bilibili.com/video/av%d/", rp.Oid)
		}
		title = m.Title
		native = true
	case model.SubTypeTopic:
		if title, link, err = s.noticeDao.Topic(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Topic(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeDrawyoo:
		if title, link, err = s.noticeDao.Drawyoo(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Drawyoo(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeActivity:
		if title, link, err = s.noticeDao.Activity(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Activity(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeForbiden:
		title, link, err = s.noticeDao.Ban(c, rp.Oid)
		if err != nil {
			return
		}
	case model.SubTypeNotice:
		title, link, err = s.noticeDao.Notice(c, rp.Oid)
		if err != nil {
			return
		}
	case model.SubTypeActArc:
		if title, link, err = s.noticeDao.ActivitySub(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.ActivitySub(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeArticle:
		var m map[int64]*artmdl.Meta
		arg := &artmdl.ArgAids{
			Aids: []int64{rp.Oid},
		}
		m, err = s.articleSrv.ArticleMetas(c, arg)
		if err != nil || m == nil {
			log.Error("s.articleSrv.ArticleMetas(%v) ret:%v error(%v)", arg, m, err)
			return
		}
		if meta, ok := m[rp.Oid]; ok {
			title = meta.Title
			link = fmt.Sprintf("http://www.bilibili.com/read/cv%d", rp.Oid)
		}
	case model.SubTypeLiveVideo:
		if title, link, err = s.noticeDao.LiveSmallVideo(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.LiveSmallVideo(%d) error(%v)", rp.Oid, err)
			return
		}
		native = true
	case model.SubTypeLiveAct:
		if title, link, err = s.noticeDao.LiveActivity(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.LiveActivity(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeLiveNotice:
		//if title, link, err = s.noticeDao.LiveNotice(c, rp.Oid); err != nil {
		//	log.Error("s.noticeDao.LiveNotice(%d) error(%v)", rp.Oid, err)
		//	return
		//}
		// NOTE 忽略直播公告跳转链接
		return
	case model.SubTypeLivePicture:
		if title, link, err = s.noticeDao.LivePicture(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.LivePiture(%d) error(%v)", rp.Oid, err)
			return
		}
		native = true
	case model.SubTypeCredit:
		if title, link, err = s.noticeDao.Credit(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Credit(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeDynamic:
		if title, link, err = s.noticeDao.Dynamic(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Dynamic(%d) error(%v)", rp.Oid, err)
			return
		}
		native = true
	case model.SubTypeAudio:
		if title, link, err = s.noticeDao.Audio(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Audio(%d) error(%v)", rp.Oid, err)
			return
		}
		native = true
	case model.SubTypeAudioPlaylist:
		if title, link, err = s.noticeDao.AudioPlayList(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.AudioPlayList(%d) error(%v)", rp.Oid, err)
			return
		}
		native = true

	default:
		return
	}
	tmp := []rune(title)
	if len(tmp) > _msgTitleSize {
		title = string(tmp[:_msgTitleSize])
	}
	jump = fmt.Sprintf("%s#reply%d", link, rp.RpID)
	tmp = []rune(rp.Content.Message)
	if len(tmp) > _msgContentSize {
		msg = string(tmp[:_msgContentSize])
	} else {
		msg = rp.Content.Message
	}
	if native {
		rootID := rp.Root
		if rootID == 0 {
			rootID = rp.RpID
		}
		nativeJump = fmt.Sprintf("bilibili://comment/detail/%d/%d/%d/?subType=%d&anchor=%d&showEnter=1&extraIntentId=%d", rp.Type, rp.Oid, rootID, subType, rp.RpID, extraIntentID)
	}
	return
}

func extraInfo(newJump string) string {
	var a = struct {
		CmNewURL struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"cm_new_url"`
	}{
		CmNewURL: struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}{
			Title:   newJump,
			Content: newJump,
		},
	}
	b, _ := json.Marshal(a)
	return string(b)
}
