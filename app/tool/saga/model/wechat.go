package model

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
