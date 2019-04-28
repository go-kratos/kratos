package model

import "go-common/library/ecode"

// Pagination Pagination.
type Pagination struct {
	PageSize int `form:"page_size" json:"page_size"`
	PageNum  int `form:"page_num" json:"page_num"`
}

// Verify verify the value of pageNum and pageSize.
func (p *Pagination) Verify() error {
	if p.PageNum < 0 {
		return ecode.MerlinIllegalPageNumErr
	} else if p.PageNum == 0 {
		p.PageNum = DefaultPageNum
	}
	if p.PageSize < 0 {
		return ecode.MerlinIllegalPageSizeErr
	} else if p.PageSize == 0 {
		p.PageSize = DefaultPageSize
	}
	return nil
}

// HookURLUpdateReq Hook URL Update Req.
type HookURLUpdateReq struct {
	ID          int64    `json:"id"`
	URL         string   `json:"url"`
	WorkspaceID int      `json:"workspace_id"`
	Status      int      `json:"status"`
	Events      []string `json:"events"`
}

// QueryHookURLReq Query Hook URL Req
type QueryHookURLReq struct {
	Pagination
	HookURLUpdateReq
	UpdateBy string `json:"update_by"`
}

// QueryHookURLRep Query Hook URL Rep.
type QueryHookURLRep struct {
	Pagination
	Total    int64      `json:"total"`
	HookUrls []*HookUrl `json:"hook_urls"`
}

// EventRequest Event Request.
type EventRequest struct {
	Event       Event  `json:"event"`
	WorkspaceID string `json:"workspace_id"`
	EventID     string `json:"id"`
	Created     string `json:"created"`
	Secret      string `json:"secret"`
}

// EventCallBackRequest Event CallBack Request.
type EventCallBackRequest struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// QueryEventLogReq Query Event Log Req
type QueryEventLogReq struct {
	Pagination
	Event       Event `json:"event"`
	WorkspaceID int   `json:"workspace_id"`
	EventID     int   `json:"id"`
}

// QueryEventLogRep Query Event Log Rep.
type QueryEventLogRep struct {
	Pagination
	Total     int64       `json:"total"`
	EventLogs []*EventLog `json:"event_logs"`
}
