package model

import xtime "go-common/library/time"

// Notice struct
type Notice struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

// Reason struct
type Reason struct {
	ID      int64  `json:"id"`
	Reason  string `json:"reason"`
	Content string `json:"content"`
}

// KPI is jury kpi info.
type KPI struct {
	ID           int64      `json:"id"`
	Mid          int64      `json:"mid"`
	Number       int64      `json:"number"`
	Day          xtime.Time `json:"day"`
	Rate         int64      `json:"rate"`
	Rank         int64      `json:"rank"`
	RankPer      int64      `json:"rankper"`
	RankTotal    int64      `json:"rankTotal"`
	Point        int64      `json:"point"`
	ActiveDays   int64      `json:"activeDays"`
	VoteTotal    int64      `json:"voteTotal"`
	VoteRadio    int64      `json:"voteRadio"`
	BlockedTotal int64      `json:"blockedTotal"`
	TermStart    xtime.Time `json:"termStart"`
	TermEnd      xtime.Time `json:"termEnd"`
	OpinionLikes int64      `json:"opinion_likes"`
}

//KPIData is jury kpi data info.
type KPIData struct {
	KPI
	VoteRealTotal int64 `json:"vote_real_total"`
}

// Opinion jury vote opinion.
type Opinion struct {
	// user info.
	Mid     int64  `json:"mid,omitempty"`
	Face    string `json:"face,omitempty"`
	Name    string `json:"name,omitempty"`
	OpID    int64  `json:"opid"`
	Vote    int8   `json:"vote"`
	Content string `json:"content"`
	Attr    int8   `json:"attr"`
	Hate    int64  `json:"hate"`
	Like    int64  `json:"like"`
}

// OpinionRes get opinion response.
type OpinionRes struct {
	Count   int        `json:"count"`
	Opinion []*Opinion `json:"opinion"`
}

// SimCase struct
type SimCase struct {
	ID         int64      `json:"id"`
	Mid        int64      `json:"mid"`
	VoteRule   int64      `json:"vote_rule"`
	VoteBreak  int64      `json:"vote_break"`
	VoteDelete int64      `json:"vote_delete"`
	CaseType   int8       `json:"case_type"`
	Stime      xtime.Time `json:"start_time"`
	Etime      xtime.Time `json:"end_time"`
}
