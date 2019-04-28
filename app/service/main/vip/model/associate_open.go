package model

// ArgRegisterOpenID add open_id args.
type ArgRegisterOpenID struct {
	Mid   int64
	AppID int64
}

// RegisterOpenIDResp  register openid resp.
type RegisterOpenIDResp struct {
	OpenID string
}

// ArgOpenAuthCallBack callback args.
type ArgOpenAuthCallBack struct {
	Mid       int64
	ThirdCode string
	AppID     int64
}

// ArgUserInfoByOpenID args.
type ArgUserInfoByOpenID struct {
	AppID  int64
	OpenID string
	IP     string
}

// UserInfoByOpenIDResp resp.
type UserInfoByOpenIDResp struct {
	Name      string `json:"name"`
	BindState int32  `json:"bind_state"`
	OutOpenID string `json:"out_open_id"`
}
