package model

// ArgRefreshToken rpc refreshToken arg
type ArgRefreshToken struct {
	AppID        int32
	RefreshToken string
}

// ArgToken rpc token arg.
type ArgToken struct {
	AppID, SubID int32
	Mid          int64
}

// ArgCookie rpc cookie arg.
type ArgCookie struct {
	Mid, Expires int64
	PWD          string
}
