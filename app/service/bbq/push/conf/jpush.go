package conf

import (
	xtime "go-common/library/time"
)

// JPushConfig 极光推送配置
type JPushConfig struct {
	AppKey         string
	SecretKey      string
	Timeout        xtime.Duration
	ApnsProduction bool
}
