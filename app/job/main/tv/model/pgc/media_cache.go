package pgc

import "go-common/library/time"

// SeasonCMS defines the elements could be changed from TV CMS side
type SeasonCMS struct {
	SeasonID    int
	Cover       string
	Desc        string
	Title       string
	UpInfo      string // season update information
	Category    int    // - cn, jp, movie, tv, documentary
	Area        string // - cn, jp, others
	Playtime    time.Time
	Role        string
	Staff       string
	NewestOrder int // the newest passed ep's order
	NewestEPID  int // the newest passed ep's ID
	NewestNb    int // the newest ep's number ( after keyword filter )
	TotalNum    int
	Style       string
	OriginName  string // new fields
	Alias       string // new fields
	PayStatus   int    // season's pay status, 0||2 = free, others = pay, pass by conf
}

// EpCMS defines the elements could be changed from TV CMS side
type EpCMS struct {
	EPID     int    `json:"epid"`
	Cover    string `json:"cover"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// new fields
	PayStatus int `json:"pay_status"`
}
