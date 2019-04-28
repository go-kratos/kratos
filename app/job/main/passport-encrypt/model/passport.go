package model

// EncryptAccount aso account.
type EncryptAccount struct {
	Mid            int64  `json:"mid"`
	UserID         string `json:"userid"`
	Uname          string `json:"uname"`
	Pwd            string `json:"pwd"`
	Salt           string `json:"salt"`
	Email          string `json:"email"`
	Tel            []byte `json:"tel"`
	CountryID      int64  `json:"country_id"`
	MobileVerified int8   `json:"mobile_verified"`
	Isleak         int8   `json:"isleak"`
	Mtime          string `json:"mtime"`
}

// OriginAccount origin aso account.
type OriginAccount struct {
	Mid            int64  `json:"mid"`
	UserID         string `json:"userid"`
	Uname          string `json:"uname"`
	Pwd            string `json:"pwd"`
	Salt           string `json:"salt"`
	Email          string `json:"email"`
	Tel            string `json:"tel"`
	CountryID      int64  `json:"country_id"`
	MobileVerified int8   `json:"mobile_verified"`
	Isleak         int8   `json:"isleak"`
	Mtime          string `json:"modify_time"`
}
