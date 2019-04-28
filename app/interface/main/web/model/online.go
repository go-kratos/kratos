package model

import v1 "go-common/app/service/main/archive/api"

// Online struct of online api response
type Online struct {
	RegionCount map[int16]int `json:"region_count"`
	AllCount    int64         `json:"all_count"`
	WebOnline   int64         `json:"web_online"`
	PlayOnline  int64         `json:"play_online"`
}

// OnlineCount struct of online count api data
type OnlineCount struct {
	ConnCount int64 `json:"conn_count"`
	IPCount   int64 `json:"ip_count"`
}

// LiveOnlineCount struct of live online count api data
type LiveOnlineCount struct {
	IPConnect   int64 `json:"ip_connect"`
	TotalOnline int64 `json:"total_online"`
}

// OnlineAid online aids and count
type OnlineAid struct {
	Aid   int64 `json:"aid"`
	Count int64 `json:"count"`
}

// OnlineArc archive whit online count
type OnlineArc struct {
	*v1.Arc
	OnlineCount int64 `json:"online_count"`
}
