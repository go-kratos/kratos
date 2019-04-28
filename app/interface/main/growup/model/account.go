package model

// ActUpInfo account up info
type ActUpInfo struct {
	Nickname string `json:"nickname"`
	Face     string `json:"face"`
}

// Account Account
type Account struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
	Sex  string `json:"sex"`
	Face string `json:"face"`
	Sign string `json:"sign"`
	Rank int64  `json:"rank"`
}

// AccountInfosResult Account info result
type AccountInfosResult struct {
	Code    int                `json:"code"`
	Data    map[int64]*Account `json:"data"`
	Message string             `json:"message"`
	TTL     int64              `json:"ttl"`
}

// UpIdentify up identify
type UpIdentify struct {
	Article int `json:"article"`
	Pic     int `json:"pic"`
	Archive int `json:"archive"`
	Blink   int `json:"blink"`
}

// UperInfosResult info result
type UperInfosResult struct {
	Code    int                    `json:"code"`
	Data    map[string]*UpIdentify `json:"data"`
	Message string                 `json:"message"`
	TTL     int64                  `json:"ttl"`
}
