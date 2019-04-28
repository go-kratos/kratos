package model

import "time"

// Figure user figure model
type Figure struct {
	ID              int64     `json:"-"`
	Mid             int64     `json:"mid"`
	Score           int32     `json:"score"`
	LawfulScore     int32     `json:"lawful_score"`
	WideScore       int32     `json:"wide_score"`
	FriendlyScore   int32     `json:"friendly_score"`
	BountyScore     int32     `json:"bounty_score"`
	CreativityScore int32     `json:"creativity_score"`
	Ver             int64     `json:"-"`
	Ctime           time.Time `json:"-"`
	Mtime           time.Time `json:"-"`
}

// Rank rank fiigure
type Rank struct {
	ID         int64     `json:"-"`
	ScoreFrom  int32     `json:"score_from"`
	ScoreTo    int32     `json:"score_to"`
	Percentage int8      `json:"percentage"`
	Ver        int64     `json:"ver"`
	Ctime      time.Time `json:"-"`
	Mtime      time.Time `json:"-"`
}
