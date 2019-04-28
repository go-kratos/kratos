package pgc

import (
	"go-common/library/time"
)

// TVEpContent reprensents the content table
type TVEpContent struct {
	ID        int64
	SeasonID  int64
	CID       int64
	Title     string
	LongTitle string
	Cover     string
	Length    int32
	IsDeleted int8
	Order     int
	Ctime     time.Time
	Mtime     time.Time
	PayStatus int
}

// TVEpSeason represents the season table
type TVEpSeason struct {
	ID         int64
	OriginName string
	Title      string
	Alias      string
	Category   int8
	Desc       string
	Style      string
	Area       string
	PlayTime   time.Time
	Info       int8
	State      int8
	TotalNum   int32
	Upinfo     string
	Staff      string
	Role       string
	Copyright  string
	IsDeleted  int8
	Ctime      time.Time
	Mtime      time.Time
	Check      int8
	AuditTime  int
	Cover      string
	Valid      int    `json:"valid"`
	Producer   string `json:"producer"`
	Version    string `json:"version"`
	Status     int
}

// Offset used for mysql offset
type Offset struct {
	Begin int
	End   int
}

// TableName gives the table name of content
func (*TVEpContent) TableName() string {
	return "tv_ep_content"
}

// TableName gives the table name of season
func (*TVEpSeason) TableName() string {
	return "tv_ep_season"
}
