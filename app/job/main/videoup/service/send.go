package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/message"
	"go-common/app/job/main/videoup/model/redis"
	"go-common/library/log"
)

func (s *Service) sendChangeMsg(c context.Context, a *archive.Archive, v *archive.Video, m *message.Videoup) (err error) {
	const (
		_changeTypeID    = "1_7_13"
		_changeCopyright = "1_7_14"
		_changeTitle     = "1_7_15"
		_changeCover     = "1_7_16"
	)
	var (
		code  string
		title string
		msg   string
	)
	title = "稿件信息变动通知"
	now := time.Now().Unix()
	if m.ChangeTitle {
		code = _changeTitle
		msg = fmt.Sprintf(`您的稿件《%s》 #{（av%d）}{"http://www.bilibili.com/video/av%d/"} 标题内容已被管理员修改，相关规则请查阅投稿页和帮助中心。`, a.Title, a.Aid, a.Aid)
		s.msg.Send(c, code, title, msg, a.Mid, now)
	}
	if m.ChangeTypeID {
		code = _changeTypeID
		msg = fmt.Sprintf(`您的稿件《%s》（#{（av%d）}{"http://www.bilibili.com/video/av%d/"} 不符合分区分类规则，已被管理员移动至%s区，相关规则请查阅投稿页和帮助中心`, a.Title, a.Aid, a.Aid, s.TypeMap[a.TypeID])
		s.msg.Send(c, code, title, msg, a.Mid, now)
	}
	if m.ChangeCopyright {
		code = _changeCopyright
		msg = fmt.Sprintf(`您的稿件《%s》#{（av%d）}{"http://www.bilibili.com/video/av%d/"} 投稿类型已被管理员修改，相关规则请查阅投稿页和帮助中心。`, a.Title, a.Aid, a.Aid)
		s.msg.Send(c, code, title, msg, a.Mid, now)
	}
	if m.ChangeCover {
		code = _changeCover
		msg = fmt.Sprintf(`您的稿件《%s》#{（av%d）}{"http://www.bilibili.com/video/av%d/"} 封面已被管理员修改，相关规则请查阅投稿页和帮助中心。`, a.Title, a.Aid, a.Aid)
		s.msg.Send(c, code, title, msg, a.Mid, now)
	}
	return
}

// sendMissionMsg
func (s *Service) sendMissionMsg(c context.Context, a *archive.Archive) (err error) {
	var (
		code  = "1_7_21"
		title = "【您的稿件已通过审核】"
		msg   = `您的《%s》（%d）已经通过审核，但由于不符合本次征稿活动的规则，故该稿件无法参与本次活动的评选。 #{点击查看>>}{"http://www.bilibili.com/video/av%d/"} 如果您有疑问，请联系help@bilibili.com。更多活动信息请关注哔哩哔哩活动。`
		now   = time.Now().Unix()
	)
	msg = fmt.Sprintf(msg, a.Title, a.Aid, a.Aid)
	if err = s.msg.Send(c, code, title, msg, a.Mid, now); err != nil {
		log.Error("s.msg.Send(%s,%s,%s,%d,%d) error(%v)", code, title, msg, a.Mid, now, err)
		return
	}
	return
}

// sendNewUpperMsg
func (s *Service) sendNewUpperMsg(c context.Context, mid, aid int64) (err error) {
	var (
		code  = "1_7_22"
		title = "【 %s，有位神秘人访问了你的作品！】"
		msg   = `终于看到你的投稿啦|ω・）！欢迎加入UP主大家庭，与我们分享你的热爱。我是你的贴身小秘创作君，请收下我悄悄为你准备的 #{入门福利}{"http://member.bilibili.com/studio/annyroal/newcomer-letter?aid=%d"}，一定要亲自打开噢(/ω＼)！想了解更多UP主资讯，欢迎关注 #{@哔哩哔哩创作中心}{"https://space.bilibili.com/37090048/?from=message"} ！`
		now   = time.Now().Unix()
	)
	upper, err := s.profile(c, mid)
	if err != nil {
		log.Error("s.profile(%d) error(%v)", mid, err)
		return
	}
	title = fmt.Sprintf(title, upper.Profile.Name)
	msg = fmt.Sprintf(msg, aid)
	if err = s.msg.Send(c, code, title, msg, mid, now); err != nil {
		log.Error("s.msg.Send(%s,%s,%s,%d,%d) error(%v)", code, title, msg, mid, now, err)
		return
	}
	return
}

func (s *Service) sendMsg(c context.Context, a *archive.Archive, v *archive.Video) (err error) {
	const (
		_codePass      = "1_7_1"
		_codeRecycle   = "1_7_3"
		_codeLock      = "1_7_5"
		_codeXcodeFail = "1_7_7"
	)
	var (
		code   string
		title  string
		reason string
		msg    string
		title2 string
		msg2   string
	)
	switch a.State {
	case archive.StateOpen, archive.StateForbidUserDelay, archive.StateOrange:
		code = _codePass
		title = "您的稿件已通过审核"
		title2 = "您的视频已通过审核"
		msg = fmt.Sprintf(`您的稿件《%s》（av%d）已经通过审核，#{点击查看>>}{"http://www.bilibili.com/video/av%d/"}`, a.Title, a.Aid, a.Aid)
		msg2 = fmt.Sprintf(`您的视频《%s》（av%d）已经通过审核，#{点击查看>>}{"http://www.bilibili.com/video/av%d/"}`, a.Title, a.Aid, a.Aid)
	case archive.StateForbidRecicle:
		code = _codeRecycle
		title = "您的稿件被退回"
		title2 = "您的视频被退回"
		if v != nil {
			reason, _ = s.arc.Reason(c, v.ID)
		} else {
			reason = a.Reason
		}
		msg = fmt.Sprintf(`您的稿件《%s》（av%d）未能通过审核。原因：%s 您可以编辑稿件重新投稿，或者对审核结果进行申诉。`, a.Title, a.Aid, reason)
		msg2 = fmt.Sprintf(`您的视频《%s》（av%d）未能通过审核。原因：%s。您可以编辑稿件重新投稿，或者对审核结果进行申诉。#{点击进行编辑>>}{"https://member.bilibili.com/v/video/submit.html?type=edit&aid=%d"}`, a.Title, a.Aid, reason, a.Aid)
	case archive.StateForbidLock:
		code = _codeLock
		title = "您的视频被退回且锁定"
		title2 = title
		if v != nil {
			reason, _ = s.arc.Reason(c, v.ID)
		} else {
			reason = a.Reason
		}
		msg = fmt.Sprintf("您的稿件《%s》（av%d）未能通过审核且被锁定（锁定稿件无法被编辑）。原因：%s。", a.Title, a.Aid, reason)
		msg2 = fmt.Sprintf("您的视频《%s》（av%d）未能通过审核且被锁定（锁定稿件无法被编辑）。原因：%s。", a.Title, a.Aid, reason)
	case archive.StateForbidXcodeFail:
		if v == nil {
			log.Warn("(%d:%s)二转失败(-16)", a.Aid, title)
			return
		}
		code = _codeXcodeFail
		title = "您的视频未能成功转码"
		title2 = title
		msg = fmt.Sprintf(`您的稿件《%s》（av%d）未能成功转码。原因：%s 请检查视频文件是否可以正常播放后再重新上传视频后再进行投稿，#{点击进入编辑>>}{"http://member.bilibili.com/v/video/submit.html?type=edit&aid=%d"}`,
			a.Title, a.Aid, archive.XcodeFailMsgs[v.FailCode], a.Aid)
		msg2 = fmt.Sprintf(`您的视频《%s》（av%d）未能成功转码。原因：%s。请检查视频文件是否可以正常播放后再重新上传视频进行投稿。#{点击进行编辑>>}{"https://member.bilibili.com/v/video/submit.html?type=edit&aid=%d"}`,
			a.Title, a.Aid, archive.XcodeFailMsgs[v.FailCode], a.Aid)
	default:
		return
	}
	now := time.Now().Unix()
	s.msg.Send(c, code, title, msg, a.Mid, now)
	s.msg.Send(c, "113_1_1", title2, msg2, a.Mid, now)
	return
}

// sendAuditMsg send message when delay archive open publish or archive auto open or first round forbid
func (s *Service) sendAuditMsg(c context.Context, route string, aid int64) {
	var (
		msg = &message.Videoup{
			Route:     route,
			Aid:       aid,
			Timestamp: time.Now().Unix(),
		}
	)
	k := strconv.FormatInt(aid, 10)
	log.Info("s.sendAuditMsg() key(%s) msg(%v)", k, msg)
	if err := s.videoupPub.Send(c, k, msg); err != nil {
		log.Error("s.sendAuditMsg() key(%s) msg(%v) error (%v)", k, msg, err)
		s.syncRetry(c, aid, 0, redis.ActionForSendOpenMsg, msg.Route, "")
	}
}

// sendPostFirstRound send message when first round after async status
func (s *Service) sendPostFirstRound(c context.Context, route string, aid int64, filename string, adminChange bool) {
	var (
		msg = &message.Videoup{
			Route:       route,
			Aid:         aid,
			Filename:    filename,
			AdminChange: adminChange,
			Timestamp:   time.Now().Unix(),
		}
		bs []byte
	)
	k := strconv.FormatInt(aid, 10)
	log.Info("sendPostFirstRound key(%s) msg(%+v)", k, msg)
	if err := s.videoupPub.Send(c, k, msg); err != nil {
		log.Error("sendPostFirstRound s.videoupPub.Send key(%s) msg(%+v) error (%v)", k, msg, err)
		if bs, err = json.Marshal(msg); err != nil {
			log.Error("sendPostFirstRound json.Marshal error(%v)", err)
			return
		}
		s.syncRetry(c, aid, 0, redis.ActionForPostFirstRound, msg.Route, string(bs))
	}
}
