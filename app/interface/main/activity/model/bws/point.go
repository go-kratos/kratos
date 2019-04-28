package bws

import "go-common/library/time"

// UserPoint .
type UserPoint struct {
	ID     int64     `json:"id"`
	Pid    int64     `json:"pid"`
	Points int64     `json:"points"`
	Ctime  time.Time `json:"ctime"`
}

// UserPointDetail .
type UserPointDetail struct {
	*UserPoint
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Fid          int64  `json:"fid"`
	Image        string `json:"image"`
	Unlocked     int64  `json:"unlocked"`
	LoseUnlocked int64  `json:"lose_unlocked"`
	LockType     int64  `json:"lockType"`
	Dic          string `json:"dic"`
	Rule         string `json:"rule"`
	Bid          int64  `json:"bid"`
}
