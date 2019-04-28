package feed

import (
	cdm "go-common/app/interface/main/app-card/model"
)

// IndexParam struct
type IndexParam struct {
	Build    int    `form:"build"`
	Platform string `form:"platform"`
	MobiApp  string `form:"mobi_app"`
	Device   string `form:"device"`
	Network  string `form:"network"`
	// idx, err := strconv.ParseInt(idxStr, 10, 64)
	// if err != nil || idx < 0 {
	// 	idx = 0
	// }
	Idx int64 `form:"idx" default:"0"`
	// pull, err := strconv.ParseBool(pullStr)
	// if err != nil {
	// 	pull = true
	// }
	Pull   bool             `form:"pull" default:"true"`
	Column cdm.ColumnStatus `form:"column"`
	// loginEvent, err := strconv.Atoi(loginEventStr)
	// if err != nil {
	// 	loginEvent = 0
	// }
	LoginEvent   int    `form:"login_event" default:"0"`
	OpenEvent    string `form:"open_event"`
	BannerHash   string `form:"banner_hash"`
	AdExtra      string `form:"ad_extra"`
	Qn           int    `form:"qn" default:"0"`
	Interest     string `form:"interest"`
	Flush        int    `form:"flush"`
	AutoPlayCard int    `form:"autoplay_card"`
	Fnver        int    `form:"fnver" default:"0"`
	Fnval        int    `form:"fnval" default:"0"`
	DeviceType   int    `form:"device_type"`
	Locale       string `form:"locale"`
}
