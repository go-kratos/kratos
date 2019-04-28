package service

import (
	"context"
	"fmt"
	"go-common/app/admin/main/reply/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	accmdl "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"time"

	"net/url"
	"strconv"
)

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

// 专门给喷子，带狗发送的消息通知, 20181206 针对大忽悠事件
func (s *Service) NotifyTroll(c context.Context, mid int64) (err error) {
	var ok bool
	if ok, err = s.dao.ExsistsDelMid(c, mid); err != nil || ok {
		return
	}
	title := "评论处理通知"
	msg := fmt.Sprintf("您好，根据#{关于规范“主播吴织亚切大忽悠事件”相关言论、信息发布的公告}{\"%s\"}，您的相关评论已被清理。对于这一事件的讨论请移步公告中告知的区域进行讨论。", s.conf.Reply.Link)
	if err = s.dao.SendReplyDelMsg(c, mid, title, msg, time.Now()); err != nil {
		return
	}
	return s.dao.SetDelMid(c, mid)
}

// TitleLink TitleLink
func (s *Service) TitleLink(c context.Context, oid int64, typ int32) (title, link string, err error) {
	switch typ {
	case model.SubTypeArchive:
		arg := &arcmdl.ArgAid2{
			Aid: oid,
		}
		var m *api.Arc
		m, err = s.arcSrv.Archive3(c, arg)
		if err != nil || m == nil {
			log.Error("s.arcSrv.Archive3(%v) ret:%v error(%v)", arg, m, err)
			return
		}
		if m.RedirectURL != "" {
			// NOTE mobile native jump
			var uri *url.URL
			if uri, err = url.Parse(m.RedirectURL); err == nil {
				q := uri.Query()
				q.Set("aid", strconv.FormatInt(oid, 10))
				uri.RawQuery = q.Encode()
				link = uri.String()
			}
		} else {
			link = fmt.Sprintf("http://www.bilibili.com/video/av%d/", oid)
		}
		title = m.Title
	case model.SubTypeTopic:
		if title, link, err = s.dao.TopicTitle(c, oid); err != nil {
			log.Error("s.noticeDao.Topic(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeActivity:
		if title, link, err = s.dao.TopicTitle(c, oid); err != nil {
			log.Error("s.noticeDao.Activity(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeForbiden:
		title, link, err = s.dao.BanTitle(c, oid)
		if err != nil {
			return
		}
	case model.SubTypeNotice:
		title, link, err = s.dao.NoticeTitle(c, oid)
		if err != nil {
			return
		}
	case model.SubTypeActArc:
		if title, link, err = s.dao.TopicTitle(c, oid); err != nil {
			log.Error("s.noticeDao.ActivitySub(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeArticle:
		arg := &artmdl.ArgAids{
			Aids: []int64{oid},
		}
		var m map[int64]*artmdl.Meta
		m, err = s.articleSrv.ArticleMetas(c, arg)
		if err != nil || m == nil {
			log.Error("s.articleSrv.ArticleMetas(%v) ret:%v error(%v)", arg, m, err)
			return
		}
		if meta, ok := m[oid]; ok {
			title = meta.Title
			link = fmt.Sprintf("http://www.bilibili.com/read/cv%d", oid)
		}
	case model.SubTypeLiveVideo:
		if title, link, err = s.dao.LiveVideoTitle(c, oid); err != nil {
			log.Error("s.noticeDao.LiveSmallVideo(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeLiveAct:
		if title, link, err = s.dao.LiveActivityTitle(c, oid); err != nil {
			log.Error("s.noticeDao.LiveActivity(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeLiveNotice:
		// NOTE 忽略直播公告跳转链接
		return
	case model.SubTypeLivePicture:
		if title, link, err = s.dao.LivePictureTitle(c, oid); err != nil {
			log.Error("s.noticeDao.LivePiture(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeCredit:
		if title, link, err = s.dao.CreditTitle(c, oid); err != nil {
			log.Error("s.noticeDao.Credit(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeDynamic:
		if title, link, err = s.dao.DynamicTitle(c, oid); err != nil {
			log.Error("s.noticeDao.Dynamic(%d) error(%v)", oid, err)
			return
		}
	default:
		return
	}
	return
}

func (s *Service) moralAndNotify(c context.Context, rp *model.Reply, moral int32, notify bool, rptMid, adid int64, adname, remark string, reason, freason int32, ftime int64, isPunish bool) {
	var err error
	title, link, _, msg := s.messageInfo(c, rp)
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
			log.Error("s.reportNotify(%d,%d,%d) error(%v)", rp.Oid, rp.Type, rp.ID, err)
		}
	}
	if ftime != 0 {
		if err = s.dao.BlockAccount(c, rp.Mid, ftime, notify, freason, title, rp.Content.Message, link, adname, remark); err != nil {
			log.Error("s.reportNotify(%d,%d,%d) error(%v)", rp.Oid, rp.Type, rp.ID, err)
		}
	}
	if err = s.reportNotify(c, rp, title, link, msg, ftime, reason, freason, isPunish); err != nil {
		log.Error("s.reportNotify(%d,%d,%d) error(%v)", rp.Oid, rp.Type, rp.ID, err)
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
		if err = s.dao.SendReplyDelMsg(c, rp.Mid, mt, mc, rp.MTime.Time()); err != nil {
			log.Error("s.messageDao.DeleteReply failed, (%d) error(%v)", rp.Mid, err)
		}
		log.Info("notify oid:%d type:%d rpID:%d reason:%d content:%s", rp.Oid, rp.Type, rp.ID, reason, mc)
	} else {
		log.Warn("no notify oid:%d type:%d rpid:%d", rp.Oid, rp.Type, rp.ID)
	}
}

func (s *Service) messageInfo(c context.Context, rp *model.Reply) (title, link, jump, msg string) {
	tmpMsg := []rune(rp.Content.Message)
	if len(tmpMsg) > 80 {
		msg = string(tmpMsg[:80])
	} else {
		msg = rp.Content.Message
	}
	var err error
	switch rp.Type {
	case model.SubTypeArchive:
		arg := &arcmdl.ArgAid2{
			Aid: rp.Oid,
		}
		var m *api.Arc
		m, err = s.arcSrv.Archive3(c, arg)
		if err != nil || m == nil {
			log.Error("s.arcSrv.Archive3(%v) ret:%v error(%v)", arg, m, err)
			return
		}
		if m.RedirectURL != "" {
			// NOTE mobile native jump
			var uri *url.URL
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
	case model.SubTypeTopic:
		if title, link, err = s.dao.TopicTitle(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Topic(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeActivity:
		if title, link, err = s.dao.TopicTitle(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Activity(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeForbiden:
		title, link, err = s.dao.BanTitle(c, rp.Oid)
		if err != nil {
			return
		}
	case model.SubTypeNotice:
		title, link, err = s.dao.NoticeTitle(c, rp.Oid)
		if err != nil {
			return
		}
	case model.SubTypeActArc:
		if title, link, err = s.dao.TopicTitle(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.ActivitySub(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeArticle:
		arg := &artmdl.ArgAids{
			Aids: []int64{rp.Oid},
		}
		var m map[int64]*artmdl.Meta
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
		if title, link, err = s.dao.LiveVideoTitle(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.LiveSmallVideo(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeLiveAct:
		if title, link, err = s.dao.LiveActivityTitle(c, rp.Oid); err != nil {
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
		if title, link, err = s.dao.LivePictureTitle(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.LivePiture(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeCredit:
		if title, link, err = s.dao.CreditTitle(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Credit(%d) error(%v)", rp.Oid, err)
			return
		}
	case model.SubTypeDynamic:
		if title, link, err = s.dao.DynamicTitle(c, rp.Oid); err != nil {
			log.Error("s.noticeDao.Dynamic(%d) error(%v)", rp.Oid, err)
			return
		}
	default:
		return
	}
	tmp := []rune(title)
	if len(tmp) > 40 {
		title = string(tmp[:40])
	}
	jump = fmt.Sprintf("%s#reply%d", link, rp.ID)
	log.Info("messageInfo(%d,%d) title:%s link:%s jump:%s msg:%s", rp.Type, rp.Oid, title, link, jump, msg)
	return
}

func (s *Service) reportNotify(c context.Context, rp *model.Reply, title, link, msg string, ftime int64, reason, freason int32, isPunish bool) (err error) {
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
	if rptUsers, err = s.dao.ReportUsers(c, rp.Oid, rp.Type, rp.ID); err != nil {
		log.Error("s.dao.ReportUsers(%d,%d,%d) error(%v)", rp.Oid, rp.Type, rp.ID, err)
		return
	}
	for _, rptUser = range rptUsers {
		// send the message
		if err = s.dao.SendReportAcceptMsg(c, rptUser.Mid, mt, mc, rp.MTime.Time()); err != nil {
			log.Error("s.MessageAcceptReport failed, (%d) error(%v)", rp.Mid, err)
		}
	}
	if _, err = s.dao.SetUserReported(c, rp.Oid, rp.Type, rp.ID, rp.MTime.Time()); err != nil {
		log.Error("s.dao.SetUserReported(%d, %d, %d) error(%v)", rp.Oid, rp.Type, rp.ID)
	}
	return
}
