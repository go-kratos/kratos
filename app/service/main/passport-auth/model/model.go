package model

// RefreshTokenResp refreshToken response
type RefreshTokenResp struct {
	Mid          int64  `json:"mid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expires      int64  `json:"expires"`
}

// CookieResp cookie response
type CookieResp struct {
	Mid     int64  `json:"mid"`
	Session string `json:"session"`
	CSRF    string `json:"csrf"`
	Expires int64  `json:"expires"`
}
