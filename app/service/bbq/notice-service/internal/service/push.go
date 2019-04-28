package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/json-iterator/go"

	"go-common/app/service/bbq/notice-service/api/v1"
	"go-common/app/service/bbq/notice-service/internal/conf"
	"go-common/app/service/bbq/notice-service/internal/model"
	push "go-common/app/service/bbq/push/api/grpc/v1"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func (s *Service) needPush(c context.Context, notice *v1.NoticeBase) bool {
	// TODO: 获取粉丝数量
	F, err := s.dao.FetchUserFansNum(c, notice.Mid)
	if err != nil {
		log.Errorv(c, log.KV("log", "FetchUserFansNum error"), log.KV("error", err))
	}

	var strategyMap *map[string]*conf.PushStrategy
	if F <= 1000 {
		strategyMap = &s.c.L1PushStrategy
	} else {
		strategyMap = &s.c.L2PushStrategy
	}

	var strategy *conf.PushStrategy
	switch notice.NoticeType {
	case model.NoticeTypeLike:
		strategy = (*strategyMap)["like"]
	case model.NoticeTypeComment:
		strategy = (*strategyMap)["comment"]
	case model.NoticeTypeFan:
		strategy = (*strategyMap)["follow"]
	case model.NoticeTypeSysMsg:
		strategy = (*strategyMap)["sysmsg"]
	}

	// A = -1 and B = -1 无限制
	if strategy.A == -1 && strategy.B == -1 {
		return true
	}
	// A = -1 触发无限制
	if strategy.A != -1 {
		count, err := s.dao.IncrHourPushAction(c, notice.Mid, notice.NoticeType, strategy.T)
		if count < strategy.A || err != nil {
			return false
		}
	}
	// B = -1 全天无限制
	if strategy.B != -1 {
		count, err := s.dao.IncrDailyPushCount(c, notice.Mid)
		if count > strategy.B || err != nil {
			return false
		}
	}
	// 推送触发条件：T小时内触发A次且当日推送次数小于等于B

	return true
}

func (s *Service) pushNotification(c context.Context, nid int64, notice *v1.NoticeBase) (err error) {
	// redis push action_mid
	err = s.dao.SetPushActionMid(c, notice.Mid, notice.ActionMid, notice.NoticeType)
	if err != nil {
		log.Errorv(c, log.KV("log", "SetPushActionMid error"), log.KV("error", err))
		return
	}

	if !s.needPush(c, notice) {
		return
	}

	var title, content string
	switch notice.NoticeType {
	case model.NoticeTypeLike:
		if notice.BizType == model.NoticeBizTypeSv {
			// 视频点赞
			content = model.PushMsgVideoLike
		} else if notice.BizType == model.NoticeBizTypeComment {
			// 评论点赞
			content = model.PushMsgCommentLike
		}
	case model.NoticeTypeComment:
		if notice.BizType == model.NoticeBizTypeSv {
			// 视频评论
			content = model.PushMsgVideoComment
		} else if notice.BizType == model.NoticeBizTypeComment {
			// 评论回复
			content = model.PushMsgCommentReply
		}
	case model.NoticeTypeFan:
		// 关注
		content = model.PushMsgFollow
	case model.NoticeTypeSysMsg:
		// 系统消息
		// return s.pushMessage(c, nid, notice)
		if notice.BizType == model.NoticeBizTypeCmsReview {
			// 审核类通知不推送
			return
		}
		content = notice.Text
	}

	// 填写内容详情
	midList, err := s.dao.GetPushActionMid(c, notice.Mid, notice.NoticeType)
	if err != nil || len(midList) == 0 {
		log.Errorv(c, log.KV("log", "GetPushActionMid error"), log.KV("error", err))
		return
	}
	nameList, err := s.dao.GetUserName(c, midList, 2)
	if err != nil {
		log.Errorv(c, log.KV("log", "GetUserName error"), log.KV("error", err))
		return
	}
	unreadInfo, err := s.dao.GetUnreadInfo(c, notice.Mid)
	if err != nil {
		log.Errorv(c, log.KV("log", "GetUnreadInfo error"), log.KV("error", err))
		return
	}
	if len(nameList) > 1 {
		tmp := fmt.Sprintf("等%d人", unreadInfo[int(notice.NoticeType)-1].UnreadNum)
		content = fmt.Sprintf(content, strings.Join(nameList, ","), tmp)
	} else {
		content = fmt.Sprintf(content, strings.Join(nameList, ","), "")
	}

	schema := fmt.Sprintf(model.PushSchemaNotice, notice.NoticeType)
	ext := make(map[string]string)
	ext["scheme"] = schema
	extStr, _ := jsoniter.Marshal(ext)

	dev, err := s.dao.FetchPushDev(c, notice.Mid)
	if err != nil {
		log.Errorv(c, log.KV("log", "FetchPushDev error"), log.KV("error", err))
		return
	}
	dev.SendNo = nid
	devs := []*push.Device{dev}

	body := &push.NotificationBody{
		Title:   title,
		Content: content,
		Extra:   string(extStr),
	}
	req := &push.NotificationRequest{
		Dev:  devs,
		Body: body,
	}
	result, err := s.dao.PushNotice(c, req)
	if err != nil {
		log.Errorv(c, log.KV("log", "PushNotice error"), log.KV("error", err), log.KV("result", result))
		return
	}
	err = s.dao.ClearHourPushAction(c, notice.Mid, notice.NoticeType)
	if err != nil {
		log.Errorv(c, log.KV("log", "hour push action clear error"), log.KV("error", err), log.KV("notice_type", notice.NoticeType))
	}

	// 埋点
	tracer, _ := trace.FromContext(c)
	s.Infoc.Info(tracer, notice.Mid, notice.Buvid, nid, notice.NoticeType, notice.BizId, notice.BizType, time.Now().Unix(), result)

	return
}

// func (s *Service) pushMessage(c context.Context, nid int64, notice *v1.NoticeBase) (err error) {
// 	dev, err := s.dao.FetchPushDev(c, notice.Mid)
// 	if err != nil {
// 		return
// 	}
// 	dev.SendNo = nid
// 	devs := []*push.Device{dev}

// 	schema := fmt.Sprintf(model.PushSchemaNotice, notice.NoticeType)
// 	ext := make(map[string]string)
// 	ext["shcema"] = schema
// 	extStr, _ := jsoniter.Marshal(ext)
// 	body := &push.MessageBody{
// 		Title:       notice.Title,
// 		Content:     notice.Text,
// 		ContentType: "text",
// 		Extra:       string(extStr),
// 	}
// 	req := &push.MessageRequest{
// 		Dev:  devs,
// 		Body: body,
// 	}
// 	result, err := s.dao.PushMessage(c, req)
// 	if err != nil {
// 		log.Errorv(c, log.KV("log", "PushMessage error"), log.KV("error", err), log.KV("result", result))
// 		return
// 	}
// 	err = s.dao.ClearHourPushAction(c, notice.Mid, notice.NoticeType)
// 	if err != nil {
// 		log.Errorv(c, log.KV("log", "hour push action clear error"), log.KV("error", err), log.KV("notice_type", notice.NoticeType))
// 	}

// 	// 埋点
// 	tracer, _ := trace.FromContext(c)
// 	s.Infoc.Info(tracer, notice.Mid, notice.Buvid, nid, notice.NoticeType, notice.BizId, notice.BizType, time.Now().Unix(), result)

// 	return
// }
