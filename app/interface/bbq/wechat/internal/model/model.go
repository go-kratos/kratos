package model

// const 常量
const ()

// WXToken 返回token
type WXToken struct {
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
}

// WXTicket 返回ticket
type WXTicket struct {
	Ticket string `json:"ticket"`
}

// TokenReq .
type TokenReq struct {
	Noncestr  string `json:"nonce" form:"nonce" validate:"required"`
	Timestamp string `json:"timestamp" form:"timestamp" validate:"required"`
}
