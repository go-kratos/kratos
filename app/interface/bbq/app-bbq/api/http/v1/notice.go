package v1

import (
	notice "go-common/app/service/bbq/notice-service/api/v1"
)

// NoticeNumResponse .
type NoticeNumResponse struct {
	RedDot int64 `json:"red_dot"`
}

// NoticeOverviewResponse .
type NoticeOverviewResponse struct {
	Notices []*NoticeOverview `json:"notices,omitempty"`
}

// NoticeOverview .
type NoticeOverview struct {
	UnreadNum  int64  `json:"unread_num"`
	Name       string `json:"name"`
	NoticeType int32  `json:"notice_type"`
	ShowType   int32  `json:"show_type"`
}

// NoticeListRequest .
type NoticeListRequest struct {
	Mid        int64
	NoticeType int32  `form:"notice_type" validated:"required"`
	CursorNext string `form:"cursor_next" validated:"required"`
}

// NoticeListResponse .
type NoticeListResponse struct {
	HasMore bool         `json:"has_more"`
	List    []*NoticeMsg `json:"list,omitempty"`
}

// NoticeMsg .
type NoticeMsg struct {
	*notice.NoticeBase
	ShowType    int32     `json:"show_type"`
	State       int32     `json:"state"`
	UserInfo    *UserInfo `json:"user_info,omitempty"`
	Pic         string    `json:"pic"`
	CursorValue string    `json:"cursor_value"`
	ErrMsg      string    `json:"err_msg"`
}
