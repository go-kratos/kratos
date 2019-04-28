package service

import (
	"context"
	"fmt"
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/dao"
	"go-common/app/interface/bbq/app-bbq/model"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"net/url"
	"strconv"
)

//CommentSubCursor 游标评论列表
func (s *Service) CommentSubCursor(c context.Context, mid int64, arg *v1.CommentSubCursorReq, device *bm.Device) (res *model.SubCursorRes, err error) {
	res = new(model.SubCursorRes)
	if _, err = s.dao.VideoBase(c, mid, arg.SvID); err != nil {
		log.Warnw(c, "log", "get video base fail", "svid", arg.SvID)
		return
	}

	req := map[string]interface{}{
		"oid":  arg.SvID,
		"type": arg.Type,
		"sort": arg.Sort,
		"root": arg.Root,
	}
	if len(arg.Access) != 0 {
		req["access_key"] = arg.Access
	}
	if arg.RpID != 0 {
		req["rpid"] = arg.RpID
	}
	if arg.Size != 0 {
		req["size"] = arg.Size
	}
	if arg.MinID > 0 && arg.MaxID > 0 {
		err = ecode.ParamInvalid
		return
	}
	if arg.MinID != 0 {
		req["min_id"] = arg.MinID
	}
	if arg.MaxID != 0 {
		req["max_id"] = arg.MaxID
	}

	res, err = s.dao.ReplySubCursor(c, req)
	return
}

//CommentList 游标评论列表
func (s *Service) CommentList(c context.Context, arg *v1.CommentListReq, device *bm.Device) (res *model.ReplyList, err error) {
	res = new(model.ReplyList)
	if _, err = s.dao.VideoBase(c, arg.MID, arg.SvID); err != nil {
		log.Warnw(c, "log", "get video base fail", "svid", arg.SvID)
		return
	}

	req := map[string]interface{}{
		"oid":   arg.SvID,
		"type":  arg.Type,
		"sort":  arg.Sort,
		"nohot": arg.NoHot,
	}

	if len(arg.Access) != 0 {
		req["access_key"] = arg.Access
	}
	if arg.Pn != 0 {
		req["pn"] = arg.Pn
	}
	if arg.Ps != 0 {
		req["ps"] = arg.Ps
	}
	if device.Build != 0 {
		req["build"] = arg.Build
	}
	if device.RawPlatform != "" {
		req["plat"] = arg.Plat
	}
	if device.RawMobiApp != "" {
		req["mobi_app"] = device.RawMobiApp
	}
	res, err = s.dao.ReplyList(c, req)
	return
}

//CommentCursor 游标评论列表
func (s *Service) CommentCursor(c context.Context, arg *v1.CommentCursorReq, device *bm.Device) (res *model.CursorRes, err error) {
	res = new(model.CursorRes)
	if _, err = s.dao.VideoBase(c, arg.MID, arg.SvID); err != nil {
		log.Warnw(c, "log", "get video base fail", "svid", arg.SvID)
		return
	}

	req := map[string]interface{}{
		"oid":    arg.SvID,
		"type":   arg.Type,
		"sort":   arg.Sort,
		"max_id": arg.MaxID,
		"min_id": arg.MinID,
		"size":   arg.Size,
	}
	if arg.RpID != 0 {
		req["rpid"] = arg.RpID
	}
	if len(arg.Access) != 0 {
		req["access_key"] = arg.Access
	}
	res, err = s.dao.ReplyCursor(c, req)
	return
}

//CommentAdd 发表评论评论服务
func (s *Service) CommentAdd(c context.Context, mid int64, arg *v1.CommentAddReq, device *bm.Device) (res *model.AddRes, err error) {
	res = new(model.AddRes)
	// 屏蔽词
	level, filterErr := s.dao.Filter(c, arg.Message, dao.FilterAreaReply)
	if filterErr != nil {
		log.Errorv(c, log.KV("log", "filter fail"))
	} else if level >= dao.FilterLevel {
		err = ecode.CommentFilterErr
		log.Warnv(c, log.KV("log", fmt.Sprintf("content filter fail: content=%s, level=%d", arg.Message, level)))
		return
	}

	var upMid int64
	videoBase, err := s.dao.VideoBase(c, mid, arg.SvID)
	if err != nil {
		log.Warnw(c, "log", "get video base fail", "svid", arg.SvID)
		return
	}
	upMid = videoBase.Mid
	parentMid := upMid

	req := map[string]interface{}{
		"oid":        arg.SvID,
		"type":       arg.Type,
		"message":    arg.Message,
		"access_key": arg.AccessKey,
	}
	if arg.At != "" {
		req["at"] = arg.At
	}
	if arg.Parent != 0 {
		req["parent"] = arg.Parent
		req["root"] = arg.Root
	} else if arg.Root != 0 {
		req["root"] = arg.Root
	}
	if arg.Plat != 0 {
		req["plat"] = arg.Plat
	}
	if arg.Device != "" {
		req["device"] = arg.Plat
	}
	if arg.Code != "" {
		req["code"] = arg.Code
	}
	res, err = s.dao.ReplyAdd(c, req)
	//wrap error
	switch ecode.Cause(err).Code() {
	case ecode.ReplyDeniedAsCaptcha.Code():
		err = ecode.CommentForbidden
	case ecode.ReplyContentOver.Code():
		err = ecode.CommentLengthIllegal
	}
	// 推送评论给通知中心
	if err == nil {
		title := "评论了你的作品"
		bizType := int32(notice.NoticeBizTypeSv)
		rootID := res.RpID
		if arg.Parent != 0 {
			title = "回复了你的评论"
			bizType = int32(notice.NoticeBizTypeComment)
			rootID = arg.Parent
			// get root comment's owner
			list, tmpErr := s.dao.ReplyMinfo(c, arg.SvID, []int64{arg.Parent})
			if tmpErr != nil || len(list) == 0 {
				log.Warnv(c, log.KV("log", "get root reply rpid info fail"), log.KV("rsp_size", len(list)))
				return
			}
			reply, exists := list[arg.Parent]
			if !exists {
				log.Errorv(c, log.KV("log", "not found reply rpid's info"))
				return
			} else if reply.Mid == 0 {
				log.Errorv(c, log.KV("log", "reply rpid's owner mid=0"))
				return
			}
			parentMid = reply.Mid
		}
		if parentMid == mid {
			log.V(1).Infov(c, log.KV("log", "action_mid=mid"), log.KV("mid", mid))
			return
		}
		urlVal := make(url.Values)
		urlVal.Add("svid", strconv.FormatInt(arg.SvID, 10))
		urlVal.Add("rootid", strconv.FormatInt(rootID, 10))
		urlVal.Add("rpid", strconv.FormatInt(res.RpID, 10))
		jumpURL := fmt.Sprintf("qing://commentdetail?%s", urlVal.Encode())
		notice := &notice.NoticeBase{
			Mid: parentMid, ActionMid: mid, SvId: arg.SvID, NoticeType: notice.NoticeTypeComment, Title: title, Text: arg.Message,
			JumpUrl: jumpURL, BizType: bizType, BizId: res.RpID}
		tmpErr := s.dao.CreateNotice(c, notice)
		if tmpErr != nil {
			log.Error("create comment notice fail: notice_msg=%s", notice.String())
		}
	}

	return
}

//CommentLike 评论点赞服务
func (s *Service) CommentLike(c context.Context, mid int64, arg *v1.CommentLikeReq, device *bm.Device) (err error) {
	if _, err = s.dao.VideoBase(c, mid, arg.SvID); err != nil {
		log.Warnw(c, "log", "get video base fail", "svid", arg.SvID)
		return
	}

	req := map[string]interface{}{
		"oid":        arg.SvID,
		"type":       arg.Type,
		"rpid":       arg.RpID,
		"action":     arg.Action,
		"access_key": arg.AccessKey,
	}
	err = s.dao.ReplyLike(c, req)
	if ecode.Cause(err).Code() == ecode.ReplyForbidAction.Code() {
		err = ecode.CommentForbidLike
		return
	}
	// TODO: 推送评论给通知中心
	if arg.Action == 1 && err == nil {
		// get root comment's owner
		list, tmpErr := s.dao.ReplyMinfo(c, arg.SvID, []int64{arg.RpID})
		if tmpErr != nil || len(list) == 0 {
			log.Warnv(c, log.KV("log", "get root rpid info fail"))
			return
		}
		reply, exists := list[arg.RpID]
		if !exists {
			log.Errorv(c, log.KV("log", "not found reply rpid's info"))
			return
		} else if reply.Mid == 0 {
			log.Errorv(c, log.KV("log", "reply rpid's owner mid=0"))
			return
		}
		parentMid := reply.Mid
		if parentMid == mid {
			log.V(1).Infov(c, log.KV("log", "action_mid=mid"), log.KV("mid", mid))
			return
		}

		text := ""
		if reply.Content != nil {
			text = reply.Content.Message
		}
		title := "点赞了你的评论"
		bizType := int32(notice.NoticeBizTypeComment)
		notice := &notice.NoticeBase{
			Mid: parentMid, ActionMid: mid, SvId: arg.SvID, NoticeType: notice.NoticeTypeLike, Title: title, Text: text,
			BizType: bizType, BizId: arg.RpID}
		tmpErr = s.dao.CreateNotice(c, notice)
		if tmpErr != nil {
			log.Errorv(c, log.KV("log", "create like notice fail: notice_msg="+notice.String()+", err="+err.Error()))
			return
		}
	}

	return
}

//CommentReport 评论举报服务
func (s *Service) CommentReport(c context.Context, arg *v1.CommentReportReq) (err error) {
	req := map[string]interface{}{
		"oid":        arg.SvID,
		"type":       arg.Type,
		"rpid":       arg.RpID,
		"reason":     arg.Reason,
		"access_key": arg.AccessKey,
	}
	if arg.Content != "" {
		req["content"] = arg.Content
	}
	err = s.dao.ReplyReport(c, req)
	return
}
