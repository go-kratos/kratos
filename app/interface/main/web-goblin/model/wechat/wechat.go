package wechat

// AccessToken wechat access token.
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}
