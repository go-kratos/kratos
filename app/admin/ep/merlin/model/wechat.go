package model

// MsgSendReq message send request of wechat
type MsgSendReq struct {
	Chatid  string         `json:"chatid" form:"chatid"`
	Msgtype string         `json:"msgtype" form:"msgtype"`
	Text    MsgSendReqText `json:"text" form:"test"`
	Safe    int            `json:"safe" form:"safe"`
}

// MsgSendReqText MegSendReq test
type MsgSendReqText struct {
	Content string `json:"content"`
}

// MsgSendRes message send response
type MsgSendRes struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	TTL     int                  `json:"ttl"`
	Data    WechatMegSendResData `json:"data"`
}

// WechatMegSendResData message send response data of wechat
type WechatMegSendResData struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
