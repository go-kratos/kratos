package model

import "encoding/json"

// Page es page
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// SearchVideo search video.
type SearchVideo struct {
	AID int64 `json:"aid"`
}

// SearchEsp big search esports.
type SearchEsp struct {
	Code       int             `json:"code,omitempty"`
	Seid       string          `json:"seid"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pagesize"`
	NumResults int             `json:"numResults"`
	NumPages   int             `json:"numPages"`
	Result     json.RawMessage `json:"result"`
}

// FilterES  filter ES video and match
type FilterES struct {
	GroupByGid []struct {
		DocCount int    `json:"doc_count"`
		Key      string `json:"key"`
	} `json:"group_by_gid"`
	GroupByMatch []struct {
		DocCount int    `json:"doc_count"`
		Key      string `json:"key"`
	} `json:"group_by_match"`
	GroupByTag []struct {
		DocCount int    `json:"doc_count"`
		Key      string `json:"key"`
	} `json:"group_by_tag"`
	GroupByTeam []struct {
		DocCount int    `json:"doc_count"`
		Key      string `json:"key"`
	} `json:"group_by_team"`
	GroupByYear []struct {
		DocCount int    `json:"doc_count"`
		Key      string `json:"key"`
	} `json:"group_by_year"`
}
