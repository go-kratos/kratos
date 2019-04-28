package model

import (
	libTime "go-common/library/time"
)

// const text
const (
	// TimeLayout time layout
	TimeLayout = "2006-01-02 15:04:05"
)

// NewbieLetterReq newbie request
type NewbieLetterReq struct {
	Aid int64 `form:"aid" validate:"required"`
	Mid int64
}

// Category category
type Category struct {
	ID   int64  `json:"id"`
	Pid  int64  `json:"pid"`
	Name string `json:"name"`
}

// CategoriesRes category result
type CategoriesRes struct {
	Code    int                 `json:"code"`
	Data    map[int64]*Category `json:"data"`
	Message string              `json:"message"`
}

// Activity activity
type Activity struct {
	ID         int64  `json:"-"`
	AndroidUrl string `json:"-"`
	H5Cover    string `json:"-"`
	ActUrl     string `json:"act_url"`
	IosUrl     string `json:"-"`
	Cover      string `json:"cover"`
	Type       int32  `json:"type"`
}

// ActivitiesRes activities result
type ActivitiesRes struct {
	Code    int         `json:"code"`
	Data    []*Activity `json:"data"`
	Message string      `json:"message"`
	TTL     int64       `json:"ttl"`
}

// VideoUpArchive video up archive
type VideoUpArchive struct {
	Aid   int64  `json:"aid"`
	Mid   int64  `json:"mid"`
	Tid   int64  `json:"tid"`
	Title string `json:"title"`
	PTime int64  `json:"ptime"`
}

// VideoUpVideo video up video
/*type VideoUpVideo struct {

}*/

// VideoUp video up
type VideoUp struct {
	Archive *VideoUpArchive `json:"archive"`
	//Videos		[]*VideoUpVideo	`json:"videos"`
}

// VideoUpRes video up result
type VideoUpRes struct {
	Code    int      `json:"code"`
	Data    *VideoUp `json:"data"`
	Message string   `json:"message"`
}

// Relation relation
type Relation struct {
	Mid       int64  `json:"mid"`
	Attribute int    `json:"attribute"`
	Face      string `json:"face"`
	Name      string `json:"name"`
}

// RelationsRes relation result
type RelationsRes struct {
	Code    int                 `json:"code"`
	Data    map[int64]*Relation `json:"data"`
	Message string              `json:"message"`
	TTL     int64               `json:"ttl"`
}

// RecommendUp table recommend up
type RecommendUp struct {
	ID       int64        `json:"id"`
	Mid      int64        `json:"mid"`
	Tid      int64        `json:"tid"`
	SubTid   int64        `json:"sub_tid"`
	Reason   string       `json:"reason"`
	Operator string       `json:"operator"`
	CTime    libTime.Time `json:"ctime"`
	MTime    libTime.Time `json:"mtime"`
}

// NewbieLetterArchive newbie letter archive
type NewbieLetterArchive struct {
	Mid   int64  `json:"-"`
	Tid   int64  `json:"-"`
	Title string `json:"title"`
	PTime string `json:"ptime"`
}

// NewbieLetterUpInfo newbie letter up info
type NewbieLetterUpInfo struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
}

// NewbieLetterRes newbie letter result
type NewbieLetterRes struct {
	UperInfo   *NewbieLetterUpInfo  `json:"uper_info"`
	Activities []*Activity          `json:"activities"`
	Relations  []*Relation          `json:"relations"`
	Archive    *NewbieLetterArchive `json:"archive"`
	Talent     string               `json:"talent"`
	Area       string               `json:"area"`
}
