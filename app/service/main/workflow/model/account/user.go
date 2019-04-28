package account

// User for sobot userinfo
type User struct {
	Mid    int64                  `json:"mid"`
	UName  string                 `json:"uname"`
	Tel    string                 `json:"tel"`
	EMail  string                 `json:"email"`
	Status int32                  `json:"status"`
	Formal int32                  `json:"formal"`
	Moral  int32                  `json:"moral"`
	Level  int32                  `json:"level"`
	Exp    string                 `json:"exp"`
	Coin   float64                `json:"coin"`
	BCoin  float64                `json:"bcoin"`
	Medal  string                 `json:"medal"`
	Up     map[string]interface{} `json:"up"`

	// extra field for further buseinss
	Extra map[string]interface{} `json:"extra"`
}
