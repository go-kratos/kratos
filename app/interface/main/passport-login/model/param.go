package model

// ParamLogin .
type ParamLogin struct {
	UserName string `form:"username"`
	Pwd      string `form:"pwd"`
}

// ParamModifyAuth .
type ParamModifyAuth struct {
	Mid     int64  `form:"mid"`
	Token   string `form:"token"`
	Refresh string `form:"refresh"`
	Session string `form:"session"`
	AppID   int64  `form:"appid"`
}
