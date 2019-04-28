package dao

import (
	"context"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	"go-common/library/log"
)

// NoticeList 获取通知列表
func (d *Dao) NoticeList(ctx context.Context, noticeType int32, mid, cursorID int64) (list []*notice.NoticeBase, err error) {
	req := &notice.ListNoticesReq{
		Mid:        mid,
		NoticeType: noticeType,
		CursorId:   cursorID,
	}

	res, err := d.noticeClient.ListNotices(ctx, req)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "notice-service:ListNotices fail"), log.KV("err", err))
		return
	}
	list = res.List
	return
}

// GetNoticeUnread 获取未读情况
func (d *Dao) GetNoticeUnread(ctx context.Context, mid int64) (list []*notice.UnreadItem, err error) {
	req := &notice.GetUnreadInfoRequest{Mid: mid}
	res, err := d.noticeClient.GetUnreadInfo(ctx, req)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "call notice service get unread info fail: err="+err.Error()))
		return
	}
	list = res.List
	log.V(1).Infov(ctx, log.KV("log", "call notice service get unread info: res="+res.String()))
	return
}

// CreateNotice 创建通知
func (d *Dao) CreateNotice(ctx context.Context, notice *notice.NoticeBase) (err error) {
	_, err = d.noticeClient.CreateNotice(ctx, notice)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "create notice fail: notice="+notice.String()))
		return
	}

	log.V(1).Infov(ctx, log.KV("log", "create notice: notice="+notice.String()))
	return
}
