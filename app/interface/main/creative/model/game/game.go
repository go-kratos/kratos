package game

import (
	xtime "go-common/library/time"
)

// ListItem str
type ListItem struct {
	GameBaseID int64  `json:"game_base_id"`
	GameName   string `json:"game_name"`
	Source     int8   `json:"source"`
	Letter     string `json:"letter"`
}

// ListWithPager fn
type ListWithPager struct {
	List  []*ListItem `json:"list"`
	Pn    int         `json:"pn"`
	Ps    int         `json:"ps"`
	Total int         `json:"total"`
}

// Info str
type Info struct {
	IsOnline  bool       `json:"is_online"`
	BaseID    int64      `json:"game_base_id"`
	Name      string     `json:"game_name"`
	Icon      string     `json:"game_icon"`
	Link      string     `json:"game_link"`
	Status    int        `json:"game_status"`
	BeginDate xtime.Time `json:"begin_date"`
}
