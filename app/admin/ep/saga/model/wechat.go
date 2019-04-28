package model

import "time"

// MaxWechatLen ...
const MaxWechatLen = 254 //企业微信内容最大长度

// AppConfig def
type AppConfig struct {
	AppID     int    // 企业微信：SAGA应用的appId
	AppSecret string // 企业微信：SAGA应用的secret
}

// Notification def
type Notification struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
}

// Text def
type Text struct {
	Content string `json:"content"`
}

// TxtNotification 文本消息
type TxtNotification struct {
	Notification
	Body Text `json:"text"`
	Safe int  `json:"safe"`
}

// AllowUserInfo 应用可见名单列表
type AllowUserInfo struct {
	Users []*UserInfo `json:"user"`
}

// UserInfo only contain userid now
type UserInfo struct {
	UserID string `json:"userid"`
}

// QueryContactLogRequest Query Contact Log Request.
type QueryContactLogRequest struct {
	Pagination
	UserID      int64  `form:"user_id"`
	UserName    string `form:"user_name"`
	OperateUser string `form:"operate_user"`
	OperateType string `form:"operate_type"`
}

// QueryContactRequest Query Contact Log Request.
type QueryContactRequest struct {
	Pagination
}

// AboundContactLog Abound Contact Log.
type AboundContactLog struct {
	ContactLog
	Name string `json:"machine_name"`
}

// ContactLog Contact Log.
type ContactLog struct {
	ID            int64     `json:"-" gorm:"column:id"`
	Username      string    `json:"username" gorm:"column:username"`
	MachineID     int64     `json:"machine_id" gorm:"column:machine_id"`
	OperateType   string    `json:"operate_type" gorm:"column:operation_type"`
	OperateResult string    `json:"operate_result" gorm:"column:operation_result"`
	OperateTime   time.Time `json:"operate_time" gorm:"column:ctime;default:current_timestamp"`
	UTime         time.Time `json:"-" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
	Type          int       `json:"type" gorm:"column:type"`
}

// Contact Contact info.
type Contact struct {
	ID       int64  `json:"-" gorm:"column:id"`
	Username string `json:"user_name" gorm:"column:user_name"`
	UserID   string `json:"user_id" gorm:"column:user_id"`
}

// PaginateContactLog Paginate Contact Log.
type PaginateContactLog struct {
	Total       int64               `json:"total"`
	PageNum     int                 `json:"page_num"`
	PageSize    int                 `json:"page_size"`
	MachineLogs []*AboundContactLog `json:"machine_logs"`
}

// PaginateContact Paginate Contact.
type PaginateContact struct {
	Total    int64          `json:"total"`
	PageNum  int            `json:"page_num"`
	PageSize int            `json:"page_size"`
	Contacts []*ContactInfo `json:"contacts"`
}

// CreateChatReq ...
type CreateChatReq struct {
	Name     string   `json:"name" validate:"required"`
	Owner    string   `json:"owner" validate:"required"`
	UserList []string `json:"userlist" validate:"required"`
	ChatID   string   `json:"chatid" validate:"required"`
}

// WechatCreateLog ...
type WechatCreateLog struct {
	ID     int       `json:"id" gorm:"column:id"`
	Name   string    `json:"name" gorm:"column:name"`
	Owner  string    `json:"owner" gorm:"column:owner"`
	ChatID string    `json:"chatid" gorm:"column:chatid"`
	Cuser  string    `json:"cuser" gorm:"column:cuser"`
	Ctime  time.Time `form:"ctime" json:"ctime" gorm:"column:ctime"`
	Status int       `form:"status" json:"status" gorm:"column:status"` //1创建 2修改 3同步中 4同步完成 5同步失败
}

// WechatChatLog ...
type WechatChatLog struct {
	ID      int    `json:"id" gorm:"column:id"`
	ChatID  string `json:"chatid" gorm:"column:chatid"`
	MsgType string `json:"msgtype" gorm:"column:msgtype"`
	Content string `json:"content" gorm:"column:content"`
	Safe    int    `json:"safe" gorm:"column:safe"`
	Status  int    `form:"status" json:"status" gorm:"column:status"` //1成功 0失败
}

// WechatMessageLog ...
type WechatMessageLog struct {
	ID      int    `json:"id" gorm:"column:id"`
	Touser  string `json:"touser" gorm:"column:touser"`
	Content string `json:"content" gorm:"column:content"`
	Status  int    `form:"status" json:"status" gorm:"column:status"` //1成功 0失败
}

// ChatResp ...
type ChatResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// CreateChatResp ...
type CreateChatResp struct {
	*ChatResp
	ChatID string `json:"chatid"`
}

// CreateChatLog ...
type CreateChatLog struct {
	*WechatCreateLog
	Buttons []string `json:"buttons"`
}

// CreateChatLogResp ...
type CreateChatLogResp struct {
	Total int `json:"total"`
	*Pagination
	Logs []*CreateChatLog `json:"logs,omitempty"`
}

// GetChatResp ...
type GetChatResp struct {
	*ChatResp
	ChatInfo *CreateChatReq `json:"chat_info"`
}

// SendChatReq ...
type SendChatReq struct {
	ChatID  string `json:"chatid"`
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Safe int `json:"safe"`
}

// SendMessageReq ...
type SendMessageReq struct {
	Touser  []string `json:"touser"`
	Content string   `json:"content"`
}

// UpdateChatReq ...
type UpdateChatReq struct {
	ChatID      string   `json:"chatid"`
	Name        string   `json:"name"`
	Owner       string   `json:"owner"`
	AddUserList []string `json:"add_user_list"`
	DelUserList []string `json:"del_user_list"`
}
