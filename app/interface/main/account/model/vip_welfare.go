package model

import "go-common/library/time"

// WelfareRes welfare list info
type WelfareRes struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	HomepageUri string `json:"homepage_uri"`
	BackdropUri string `json:"backdrop_uri"`
	Tid         int    `json:"tid"`
	Rank        int    `json:"rank"`
}

// WelfareInfo welfare info
type WelfareInfo struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Desc         string    `json:"desc"`
	HomepageUri  string    `json:"homepage_uri"`
	BackdropUri  string    `json:"backdrop_uri"`
	SurplusCount int       `json:"surplus_count"`
	Received     bool      `json:"received"`
	Stime        time.Time `json:"stime"`
	Etime        time.Time `json:"etime"`
}
