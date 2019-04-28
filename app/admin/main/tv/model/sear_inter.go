package model

import (
	"go-common/library/time"
)

//SearInter reprensents the search intervene
type SearInter struct {
	ID         int64     `json:"id" params:"id"`
	Searchword string    `json:"searchword" params:"searchword"`
	Rank       int64     `json:"rank" params:"rank"`
	Tag        string    `json:"tag" params:"tag"`
	Deleted    int8      `json:"deleted"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}

// TableName gives the table name of search intervene
func (*SearInter) TableName() string {
	return "search_intervene"
}

// SearInterPager search intervene pager
type SearInterPager struct {
	TotalCount int          `json:"total_count"`
	Pn         int          `json:"pn"`
	Ps         int          `json:"ps"`
	Items      []*SearInter `json:"items"`
	PubState   int8
	PubTime    string
}

//OutSearchInter output search intervene
type OutSearchInter struct {
	Keyword string `json:"keyword"`
	Status  string `json:"status"`
}

//PublishStatus search intervene publish status state 0-unPublish 1-publish
type PublishStatus struct {
	Time  string
	State int8
}
