package model

// App represents App info.
type App struct {
	AppID     int64  `json:"appid" gorm:"column:appid"`
	AppName   string `json:"app_name" gorm:"column:app_name"`
	AppKey    string `json:"appkey" gorm:"column:appkey"`
	AppSecret string `json:"app_secret" gorm:"column:app_secret"`
	Enabled   int64  `json:"enabled" gorm:"column:enabled" default:"1"`
}

// AppListParams represents SearchAppParams info.
type AppListParams struct {
	AppKey  string `form:"sappkey"`
	AppName string `form:"app_name"`
	PS      int64  `form:"ps" default:"10"`
	PN      int64  `form:"pn" default:"1"`
}

// AppParams .
type AppParams struct {
	AppID   int64  `json:"appid" form:"appid" validate:"required"`
	AppName string `json:"app_name" form:"app_name" validate:"required"`
}
