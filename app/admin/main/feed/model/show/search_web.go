package show

import (
	"go-common/app/admin/main/feed/model/common"
	xtime "go-common/library/time"
)

//SearchWeb search web
type SearchWeb struct {
	ID          int64             `json:"id" form:"id"`
	CardType    int               `json:"card_type" form:"card_type"`
	CardValue   string            `json:"card_value" form:"card_value"`
	Stime       xtime.Time        `json:"stime" form:"stime"`
	Etime       xtime.Time        `json:"etime" form:"etime"`
	Check       int               `json:"check" form:"check"`
	Status      int               `json:"status" form:"status"`
	Priority    int               `json:"priority" form:"priority"`
	Person      string            `json:"person" form:"person"`
	ApplyReason string            `json:"apply_reason" form:"apply_reason"`
	Deleted     int               `json:"deleted" form:"deleted"`
	Query       []*SearchWebQuery `json:"query" form:"query" gorm:"-"`
	Card        interface{}       `json:"card" gorm:"-"`
}

//SearchWebPager .
type SearchWebPager struct {
	Item []*SearchWeb `json:"item"`
	Page common.Page  `json:"page"`
}

// TableName .
func (a SearchWeb) TableName() string {
	return "search_web"
}

/*
---------------------------
 struct param
---------------------------
*/

//SearchWebAP add param
type SearchWebAP struct {
	ID          int64      `json:"id" form:"id"`
	CardType    int        `json:"card_type" form:"card_type" validate:"required"`
	CardValue   string     `json:"card_value" form:"card_value" validate:"required"`
	Stime       xtime.Time `json:"stime" form:"stime" validate:"required"`
	Etime       xtime.Time `json:"etime" form:"etime" validate:"required"`
	Priority    int        `json:"priority" form:"priority" validate:"required"`
	Check       int        `form:"check" default:"1"`
	Person      string     `json:"person" form:"person"`
	ApplyReason string     `json:"apply_reason" form:"apply_reason"`
	Query       string     `json:"query" form:"query" gorm:"-" validate:"required"`
}

//SearchWebUP update param
type SearchWebUP struct {
	ID          int64      `form:"id" validate:"required"`
	CardType    int        `json:"card_type" form:"card_type"`
	CardValue   string     `json:"card_value" form:"card_value"`
	Stime       xtime.Time `json:"stime" form:"stime"`
	Etime       xtime.Time `json:"etime" form:"etime"`
	Check       int        `json:"check" form:"check"`
	Status      int        `json:"status" form:"status"`
	Priority    int        `json:"priority" form:"priority"`
	Person      string     `json:"person" form:"person"`
	ApplyReason string     `json:"apply_reason" form:"apply_reason"`
	Query       string     `json:"query" form:"query" gorm:"-" validate:"required"`
}

//SearchWebLP list param
type SearchWebLP struct {
	ID     int    `form:"id"`
	Check  int    `form:"check"`
	Person string `form:"person"`
	STime  string `form:"stime"`
	ETime  string `form:"etime"`
	Ps     int    `form:"ps" default:"20"`
	Pn     int    `form:"pn" default:"1"`
}

//SearchWebOption option web card (online,hidden,pass,reject)
type SearchWebOption struct {
	ID     int64 `form:"id" validate:"required"`
	Check  int   `json:"check" form:"check"`
	Status int   `json:"status" form:"status"`
}

//SWTimeValid option web card (online,hidden,pass,reject)
type SWTimeValid struct {
	ID        int64
	Query     string
	Priority  int
	STime     xtime.Time
	ETime     xtime.Time
	CardValue string
}

// TableName .
func (a SearchWebOption) TableName() string {
	return "search_web"
}

// TableName .
func (a SearchWebAP) TableName() string {
	return "search_web"
}

// TableName .
func (a SearchWebUP) TableName() string {
	return "search_web"
}
