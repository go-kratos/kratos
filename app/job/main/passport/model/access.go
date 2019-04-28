package model

// AccessInfo aso_app_perm table
type AccessInfo struct {
	AppID   int32  `json:"appid"`
	Mid     int64  `json:"mid"`
	Token   string `json:"access_token"`
	Expires int64  `json:"expires"`
}
