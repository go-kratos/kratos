package model

// ArgAuthCode get open_id args.
type ArgAuthCode struct {
	Code  string `form:"auth_code" validate:"required"`
	APPID int64
	IP    string
}

// OpenIDResp open_id resp.
type OpenIDResp struct {
	OpenID string `json:"open_id"`
}

// OAuth2InfoResp oauth2 resp.
type OAuth2InfoResp struct {
	Mid   int64  `json:"mid"`
	Uname string `json:"uname"`
}

// ArgBind bind args.
type ArgBind struct {
	OpenID    string `form:"open_id" validate:"required"`
	OutOpenID string `form:"out_open_id" validate:"required"`
	AppID     int64
}

// ArgUserInfoByOpenID args.
type ArgUserInfoByOpenID struct {
	OpenID string `form:"open_id" validate:"required"`
	AppID  int64
	IP     string
}

// ArgBindInfo bind info args.
type ArgBindInfo struct {
	Mid   int64
	AppID int64
}

// ArgBilibiliPrizeGrant args.
type ArgBilibiliPrizeGrant struct {
	PrizeKey string `form:"prize_key" validate:"required"`
	UniqueNo string `form:"unique_no" validate:"required"`
	OpenID   string `form:"open_id" validate:"required"`
	AppID    int64
}

// ArgBilibiliVipGrant bilibili vip grant args.
type ArgBilibiliVipGrant struct {
	OpenID     string `form:"open_id" validate:"required"`
	OutOpenID  string `form:"out_open_id" validate:"required"`
	OutOrderNO string `form:"out_order_no" validate:"required"`
	Duration   int32  `form:"duration" validate:"required"`
	AppID      int64
}

// ArgOpenAuthCallBack args.
type ArgOpenAuthCallBack struct {
	ThirdCode string `form:"auth_code" validate:"required"`
	State     string `form:"state" validate:"required"`
	Mid       int64
	AppID     int64
}
