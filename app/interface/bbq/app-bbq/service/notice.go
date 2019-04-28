package service

import (
	"context"
	"encoding/json"
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/conf"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/app/service/bbq/common"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	"go-common/library/ecode"
	"go-common/library/log"
)

var noticeTypes map[int32]bool

const (
	_noticeTypeLike    = 1
	_noticeTypeComment = 2
	_noticeTypeFan     = 3
	_noticeTypeSysMsg  = 4
)

func init() {
	noticeTypes = map[int32]bool{
		_noticeTypeLike:    true,
		_noticeTypeComment: true,
		_noticeTypeFan:     true,
		_noticeTypeSysMsg:  true,
	}
}

// GetNoticeNum 通知红点
func (s *Service) GetNoticeNum(ctx context.Context, mid int64) (res *v1.NoticeNumResponse, err error) {
	res = new(v1.NoticeNumResponse)
	unreadList, err := s.dao.GetNoticeUnread(ctx, mid)
	if err != nil {
		log.Warnv(ctx, log.KV("log", "get notice unread fail: err="+err.Error()))
		return
	}
	var redDot int64
	for _, v := range unreadList {
		redDot += v.UnreadNum
	}
	res.RedDot = redDot
	return
}

// NoticeOverview 通知中心概述
func (s *Service) NoticeOverview(ctx context.Context, mid int64) (res *v1.NoticeOverviewResponse, err error) {
	res = new(v1.NoticeOverviewResponse)
	unreadList, err := s.dao.GetNoticeUnread(ctx, mid)
	if err != nil {
		log.Warnv(ctx, log.KV("log", "get notice unread fail: err="+err.Error()))
		return
	}
	unreadMap := make(map[int32]*notice.UnreadItem)
	for _, item := range unreadList {
		unreadMap[item.NoticeType] = item
	}

	for _, v := range conf.Conf.Notices {
		var unreadNum int64
		if unreadInfo, exists := unreadMap[v.NoticeType]; exists {
			unreadNum = unreadInfo.UnreadNum
		}
		res.Notices = append(res.Notices, &v1.NoticeOverview{NoticeType: v.NoticeType, ShowType: v.ShowType, Name: v.Name, UnreadNum: unreadNum})
	}

	return
}

// NoticeList 请求通知列表，组装成通知消息列表
func (s *Service) NoticeList(ctx context.Context, req *v1.NoticeListRequest) (res *v1.NoticeListResponse, err error) {
	res = new(v1.NoticeListResponse)
	// 0. 校验请求合法性
	if _, exists := noticeTypes[req.NoticeType]; !exists {
		return nil, ecode.NoticeTypeErr
	}

	cursor, _, err := parseCursor("", req.CursorNext)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "parse cursor fail: cursor_next="+req.CursorNext))
		return
	}

	// 1. 请求notice列表
	list, err := s.dao.NoticeList(ctx, req.NoticeType, req.Mid, cursor.CursorID)
	if err != nil {
		return
	}

	// 2. 根据notice请求相应业务方
	// 2.0 整理notice列表
	var actionMIDs []int64 // 不是map，有可能重复
	var svIDs []int64
	for _, item := range list {
		if item.ActionMid != 0 {
			actionMIDs = append(actionMIDs, item.ActionMid)
		}
		if item.SvId != 0 {
			svIDs = append(svIDs, item.SvId)
		}
	}
	needUserInfo := true
	needFollowState := false
	needSvInfo := false
	switch req.NoticeType {
	case _noticeTypeLike:
		needSvInfo = true
	case _noticeTypeComment:
		needSvInfo = true
	case _noticeTypeFan:
		needFollowState = true
	case _noticeTypeSysMsg:
		needUserInfo = false
	}
	// 2. 请求用户信息
	var userInfos map[int64]*v1.UserInfo
	if needUserInfo && len(actionMIDs) > 0 {
		userInfos, err = s.dao.BatchUserInfo(ctx, req.Mid, actionMIDs, false, false, needFollowState)
		if err != nil {
			log.Errorv(ctx, log.KV("log", "batch fetch user infos fail"))
			err = nil
		}
	}
	// 2. 请求视频信息
	svInfos := make(map[int64]*model.SvInfo)
	if needSvInfo && len(svIDs) > 0 {
		svList, _, err := s.dao.GetVideoDetail(ctx, svIDs)
		if err != nil {
			log.Warnv(ctx, log.KV("log", "batch fetch sv detail fail"))
			err = nil
		} else {
			for _, sv := range svList {
				svInfos[sv.SVID] = sv
			}
		}
	}

	// 3. 组成回包
	for _, item := range list {
		var noticeMsg v1.NoticeMsg
		noticeMsg.NoticeBase = item
		// showtype没有用，这里就按照noticetype返回
		// md，想的好好的扩展灵活性，完全没被践行
		noticeMsg.ShowType = req.NoticeType
		if needUserInfo && (noticeMsg.ActionMid != 0) {
			if val, exists := userInfos[noticeMsg.ActionMid]; exists {
				noticeMsg.UserInfo = val
			}
		}
		if needSvInfo && (noticeMsg.SvId != 0) {
			if val, exists := svInfos[noticeMsg.SvId]; exists && noticeMsg.State >= common.VideoStPendingPassReview {
				noticeMsg.Pic = val.CoverURL
			} else {
				noticeMsg.State = 1
				noticeMsg.ErrMsg = "视频不见了"
			}
		}
		// cursor_value
		var cursor model.CursorValue
		cursor.CursorID = item.Id
		jsonBytes, _ := json.Marshal(cursor)
		noticeMsg.CursorValue = string(jsonBytes)

		res.List = append(res.List, &noticeMsg)
	}
	return
}
