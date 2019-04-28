package model

import (
	"go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

type Archive struct {
	Id    int64      `json:"id"`
	Mid   int64      `json:"mid"`
	Fid   int64      `json:"fid"`
	Aid   int64      `json:"aid"`
	CTime xtime.Time `json:"-"`
	MTime xtime.Time `json:"-"`
}

type SearchArchive struct {
	Code           int    `json:"code,omitempty"`
	Seid           string `json:"seid"`
	Page           int    `json:"page"`
	PageSize       int    `json:"pagesize"`
	NumPages       int    `json:"numPages,omitempty"`
	PageCount      int    `json:"pagecount"`
	NumResults     int    `json:"numResults,omitempty"`
	Total          int    `json:"total"`
	SuggestKeyword string `json:"suggest_keyword"`
	Mid            int64  `json:"mid"`
	Fid            int64  `json:"fid"`
	Tid            int    `json:"tid"`
	Order          string `json:"order"`
	Keyword        string `json:"keyword"`
	TList          []struct {
		Tid   int    `json:"tid"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	} `json:"tlist,omitempty"`
	Result   []*SearchArchiveResult `json:"result,omitempty"`
	Archives []*FavArchive          `json:"archives"`
}

type SearchArchiveResult struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Play    string `json:"play"`
	FavTime int64  `json:"fav_time"`
}

type FavArchive struct {
	*api.Arc
	FavAt          int64  `json:"fav_at"`
	PlayNum        string `json:"play_num"`
	HighlightTitle string `json:"highlight_title"`
}

type AppInfo struct {
	Platform string
	Build    string
	MobiApp  string
	Device   string
}
