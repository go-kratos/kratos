package model

import (
	"time"
)

// ActMatchs def.
type ActMatchs struct {
	ID       int64     `json:"id" form:"id"`
	Name     string    `json:"name" form:"name"`
	URL      string    `json:"url" form:"url"`
	Cover    string    `json:"cover" form:"cover"`
	SID      int64     `json:"sid" form:"sid" gorm:"column:sid"`
	MaxStake int64     `json:"max_stake" form:"max_stake"`
	Stake    int8      `json:"stake" form:"stake"`
	Status   int8      `json:"status" form:"status"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// ActMatchsObject def.
type ActMatchsObject struct {
	ID        int64     `json:"id" form:"id"`
	HomeName  string    `json:"home_name" form:"home_name"`
	HomeLogo  string    `json:"home_logo" form:"home_logo"`
	HomeScore int64     `json:"home_score" form:"home_score"`
	AwayName  string    `json:"away_name" form:"away_name"`
	AwayLogo  string    `json:"away_logo" form:"away_logo"`
	AwayScore int64     `json:"away_score" form:"away_score"`
	SID       int64     `json:"sid" gorm:"column:sid"  form:"sid"`
	MatchID   int64     `json:"match_id" form:"match_id"`
	GameStime time.Time `json:"game_stime" form:"game_stime" time_format:"2006-01-02 15:04:05"`
	Stime     time.Time `json:"stime" form:"stime" time_format:"2006-01-02 15:04:05"`
	Etime     time.Time `json:"etime" form:"etime" time_format:"2006-01-02 15:04:05"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
	Result    int8      `json:"result"  form:"result"`
	Status    int8      `json:"status"  form:"status"`
}

// TableName ActMatchsObject def.
func (ActMatchsObject) TableName() string {
	return "act_matchs_object"
}
