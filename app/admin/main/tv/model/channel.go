package model

import "go-common/library/time"

// Channel represents the table TV_RANK
type Channel struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Desc    string    `json:"desc"`
	Splash  string    `json:"splash"`
	Deleted int8      `json:"deleted"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime_nb"`
}

// ChannelFmt , mtimeFormat transforms the mtime timestamp
type ChannelFmt struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Desc        string    `json:"desc"`
	Splash      string    `json:"splash"`
	Deleted     int8      `json:"deleted"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime_nb,omitempty"`
	MtimeFormat string    `json:"mtime"`
}

//ChannelPager def.
type ChannelPager struct {
	TotalCount int64         `json:"total_count"`
	Pn         int           `json:"pn"`
	Ps         int           `json:"ps"`
	Items      []*ChannelFmt `json:"items"`
}

// ReqChannel def.
type ReqChannel struct {
	Page  int    `form:"page" default:"1"`
	Order int    `form:"order" default:"1"` // 1=desc,2=asc
	Title string `form:"title"`             // english name, precise search
	Desc  string `form:"desc"`              // chinese name, fuzzy search
}

// TableName tv_rank
func (c Channel) TableName() string {
	return "tv_channel"
}

// TableName tv_rank
func (c ChannelFmt) TableName() string {
	return "tv_channel"
}
