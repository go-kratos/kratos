package like

import (
	garcmdl "go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

// Like struct
type Like struct {
	*Item
	Archive *garcmdl.Arc `json:"archive,omitempty"`
}

// Item like item struct.
type Item struct {
	ID       int64      `json:"id"`
	Wid      int64      `json:"wid"`
	Ctime    xtime.Time `json:"act_ctime"`
	Sid      int64      `json:"sid"`
	Type     int        `json:"type"`
	Mid      int64      `json:"mid"`
	State    int        `json:"state"`
	StickTop int        `json:"stick_top"`
	Mtime    xtime.Time `json:"mtime"`
}

// GroupItem .
type GroupItem struct {
	ID      int64  `json:"id"`
	Sid     int64  `json:"sid"`
	State   int    `json:"state"`
	Type    int    `json:"type"`
	Mid     int64  `json:"mid"`
	Wid     int64  `json:"wid"`
	Ctime   string `json:"ctime"`
	Likes   int    `json:"likes"`
	Liked   int    `json:"liked"`
	Message string `json:"message"`
	Device  string `json:"device"`
	Image   string `json:"image"`
	Plat    string `json:"plat"`
	Reply   string `json:"reply"`
	Link    string `json:"link"`
}

// List .
type List struct {
	*Item
	Object   interface{} `json:"object"`
	Like     int64       `json:"like"`
	Likes    int64       `json:"likes"`
	HasLikes int8        `json:"has_likes"`
	Click    int64       `json:"click"`
	Coin     int64       `json:"coin"`
	Share    int64       `json:"share"`
	Reply    int64       `json:"reply"`
	Dm       int64       `json:"dm"`
	Fav      int64       `json:"fav"`
}

// ListInfo .
type ListInfo struct {
	List []*List `json:"list"`
	*Page
}

// LidLikeRes .
type LidLikeRes struct {
	Score int64
	Lid   int64
}

// Extend like_extend .
type Extend struct {
	ID    int64      `json:"id"`
	Lid   int64      `json:"lid"`
	Like  int64      `json:"like"`
	Ctime xtime.Time `json:"ctime"`
	Mtime xtime.Time `json:"mtime"`
}

// Tag .
type Tag struct {
	ID   int64  `json:"tag_id,omitempty"`
	Name string `json:"tag_name,omitempty"`
}

// ArgTag .
type ArgTag struct {
	Archive *garcmdl.Arc `json:"archive,omitempty"`
	Tags    []string     `json:"tags,omitempty"`
}
