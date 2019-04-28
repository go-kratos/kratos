package model

import (
	"encoding/json"
	"time"
)

type extraData struct {
	IPPort string `json:"ip_port"`
}

// LoginLog user login active log
type LoginLog struct {
	Mid       int64  `json:"mid"`
	IP        string `json:"ip"`
	Buvid     string `json:"buvid"`
	ExtraData string `json:"extra_data,omitempty"`
	Business  int    `json:"business"`
	CTime     string `json:"ctime"`
}

// NewLoginLog new loginLog
func NewLoginLog(mid int64, ip string, ipport string, buvid string) *LoginLog {
	ed, _ := json.Marshal(&extraData{IPPort: ipport})
	return &LoginLog{
		Mid:       mid,
		IP:        ip,
		Buvid:     buvid,
		ExtraData: string(ed),
		Business:  53,
		CTime:     time.Now().Format("2006-01-02 15:04:05"),
	}
}
