package model

import (
	"errors"
	"go-common/library/time"
)

const (
	TopicCacheMiss = -1
	// http mode
	HttpMode4Http  = 1 // eg "http://a.bilibili.com"
	HttpMode4Https = 2 // eg "https://a.bilibili.com"
	HttpMode4Both  = 3 // eg "//a.bilibili.com"
)

var ErrTopicRequest = errors.New("Get topic info request error")

// TopicsResult topics.
type TopicsResult struct {
	Code int `json:"code"`
	Data struct {
		List []*Topic `json:"list"`
	} `json:"data"`
}

type TopicFav struct {
	ID    int64     `json:"id"`
	Mid   int64     `json:"mid"`
	TpID  int64     `json:"tpid"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

type Topic struct {
	ID       int64     `json:"id"`
	TpID     int64     `json:"tp_id"`
	MID      int64     `json:"mid"`
	FavAt    time.Time `json:"fav_at"`
	State    int64     `json:"state"`
	Stime    string    `json:"stime"`
	Etime    string    `json:"etime"`
	Ctime    string    `json:"ctime"`
	Mtime    string    `json:"mtime"`
	Name     string    `json:"name"`
	Author   string    `json:"author"`
	PCUrl    string    `json:"pc_url"`
	H5Url    string    `json:"h5_url"`
	PCCover  string    `json:"pc_cover"`
	H5Cover  string    `json:"h5_cover"`
	Rank     int64     `json:"rank"`
	PageName string    `json:"page_name"`
	Plat     int64     `json:"plat"`
	Desc     string    `json:"desc"`
	Click    int64     `json:"click"`
	TPType   int64     `json:"type"`
	Mold     int64     `json:"mold"`
	Series   int64     `json:"series"`
	Dept     int64     `json:"dept"`
	ReplyID  int64     `json:"reply_id"`
}

type TopicList struct {
	PageNum  int      `json:"page"`
	PageSize int      `json:"pagesize"`
	Total    int64    `json:"total"`
	List     []*Topic `json:"list"`
}
