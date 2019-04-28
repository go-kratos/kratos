package model

// Token user access token.
type Token struct {
	Mid     int64  `json:"mid"`
	Appid   int64  `json:"appid"`
	Subid   int64  `json:"appSubid"`
	Token   string `json:"accessToken"`
	RToken  string `json:"refreshToken"`
	CTime   int64  `json:"createAt"`
	Expires int64  `json:"expires"`
	Type    int64  `json:"type"`
}
