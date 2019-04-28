package model

import (
	"go-common/library/time"
)

// SpecialAward special award info
type SpecialAward struct {
	AwardID      int64     `json:"award_id"`
	AwardName    string    `json:"award_name"`
	Divisions    []string  `json:"divisions"`
	CycleStart   time.Time `json:"cycle_start"`
	CycleEnd     time.Time `json:"cycle_end"`
	AnnounceDate time.Time `json:"announce_date"`
	Duration     int64     `json:"duration"`
	OpenStatus   int       `json:"open_status"`
}

// Resource award resource
type Resource struct {
	ResourceType  int
	ResourceIndex int
	Content       string
}

// WinningRecord winning record
type WinningRecord struct {
	AwardID   int64  `json:"award_id"`
	AwardName string `json:"award_name"`
	PrizeID   int64  `json:"prize_id"`
	State     int    `json:"state"`
}

// Poster poster
type Poster struct {
	AwardName string `json:"award_name"`
	Nickname  string `json:"nickname"`
	Face      string `json:"face"`
	PrizeName string `json:"prize_name"`
	Date      string `json:"date"`
	Bonus     int64  `json:"bonus"`
}

// SimpleSpecialAward simplify
type SimpleSpecialAward struct {
	AwardName  string    `json:"award_name"`
	AwardID    int64     `json:"award_id"`
	CycleStart time.Time `json:"cycle_start"`
}

// QA question & answer
type QA struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// UpAwardState up award state
type UpAwardState struct {
	AwardID   int64  `json:"-"`
	AwardName string `json:"award_name"`
	State     int    `json:"state"`
}

// AwardUpStatus up status
type AwardUpStatus struct {
	Joined    bool `json:"joined"`
	Qualified bool `json:"qualified"`
}
