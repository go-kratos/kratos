package model

// OriginPerm origin token.
type OriginPerm struct {
	Mid          int64  `json:"mid"`
	AppID        int32  `json:"appid"`
	AppSubID     int32  `json:"app_subid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreateAt     int64  `json:"create_at"`
	Expires      int64  `json:"expires"`
	Ctime        string `json:"ctime"`
	Mtime        string `json:"mtime"`
}

// OriginAsoAccount origin aso account.
type OriginAsoAccount struct {
	Mid   int64  `json:"mid"`
	Uname string `json:"uname"`
	Mtime string `json:"modify_time"`
}

// OriginMember origin member.
type OriginMember struct {
	Mid   int64  `json:"mid"`
	Face  string `json:"face"`
	Mtime string `json:"modify_time"`
}

// Info account info.
type Info struct {
	Mid   int64  `json:"mid"`
	Uname string `json:"uname"`
	Face  string `json:"face"`
	Email string `json:"email"`
	Tel   string `json:"tel"`
}

// Equals check info equals.
func (m *Info) Equals(other *Info) bool {
	return m.Mid == other.Mid && m.Uname == other.Uname && m.Face == other.Face
}

// AsoAccount aso account.
type AsoAccount struct {
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
	Mtime          string `json:"mtime"`
}

// Equals check perm equals, check all non primary key fields exclude create_at.
func (m *Perm) Equals(other *Perm) bool {
	return m.Mid == other.Mid && m.AppID == other.AppID && m.AccessToken == other.AccessToken && m.RefreshToken == other.RefreshToken && m.AppSubID == other.AppSubID && m.Expires == other.Expires
}
