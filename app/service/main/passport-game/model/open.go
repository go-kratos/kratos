package model

// App app info.
type App struct {
	AppID     int32  `json:"appid"`
	AppKey    string `json:"appkey"`
	AppSecret string `json:"app_secret"`
}
