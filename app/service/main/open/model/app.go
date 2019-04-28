package model

// App represents App info.
type App struct {
	AppID     int64  `json:"appid" gorm:"primary_key"`
	AppName   string `json:"app_name" gorm:"column:app_name"`
	AppKey    string `json:"appkey" gorm:"column:appkey"`
	AppSecret string `json:"app_secret" gorm:"column:app_secret"`
}
