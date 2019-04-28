package model

// AuthToken for auth
type AuthToken struct {
	ID      int64  `json:"id"`
	Mid     int64  `json:"mid"`
	AppID   int64  `json:"appid"`
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
	Type    int64  `json:"type"`
}
