package model

import (
	"go-common/library/time"
)

// Welfare vip_welfare table
type Welfare struct {
	ID          int       `json:"id"`
	WelfareName string    `json:"welfare_name" form:"welfare_name"`
	WelfareDesc string    `json:"welfare_desc" form:"welfare_desc"`
	HomepageUri string    `json:"homepage_uri" form:"homepage_uri"`
	BackdropUri string    `json:"backdrop_uri" form:"backdrop_uri"`
	Recommend   int       `json:"recommend"`
	Rank        int       `json:"rank"`
	Tid         int       `json:"tid"`
	UsageForm   int       `json:"usage_form" form:"usage_form"`
	ReceiveRate int       `json:"receive_rate" form:"receive_rate"`
	ReceiveUri  string    `json:"receive_uri" form:"receive_uri"`
	VipType     int       `json:"vip_type" form:"vip_type"`
	State       int       `json:"state"`
	OperID      int       `json:"oper_id"`
	OperName    string    `json:"oper_name"`
	Stime       time.Time `json:"stime"`
	Etime       time.Time `json:"etime"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

// WelfareType vip_welfare_type table
type WelfareType struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	State    int    `json:"state"`
	OperID   int    `json:"oper_id"`
	OperName string `json:"oper_name"`
}

// WelfareCodeBatch vip_welfare_code_batch table
type WelfareCodeBatch struct {
	ID            int       `json:"id" gorm:"-;primary_key;AUTO_INCREMENT" form:"id"`
	BatchName     string    `json:"batch_name"`
	Wid           int       `json:"wid"`
	Count         int       `json:"count"`
	ReceivedCount int       `json:"received_count"`
	Ver           int       `json:"ver"`
	State         int       `json:"state"`
	OperID        int       `json:"oper_id"`
	OperName      string    `json:"oper_name"`
	Vtime         time.Time `json:"vtime"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

// WelfareCode vip_welfare_code table
type WelfareCode struct {
	ID    int       `json:"id"`
	Bid   int       `json:"bid"`
	Wid   int       `json:"wid"`
	Code  string    `json:"code"`
	Mid   int       `json:"mid"`
	State int       `json:"state"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

//WelfareTypeRes welfare type response
type WelfareTypeRes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// WelfareReq save or update welfare params
type WelfareReq struct {
	ID          int       `form:"id"`
	WelfareName string    `form:"name" validate:"required"`
	WelfareDesc string    `form:"desc" validate:"required"`
	HomepageUri string    `form:"homepage_uri"`
	BackdropUri string    `form:"backdrop_uri"`
	Recommend   int       `form:"recommend"`
	Rank        int       `form:"rank" validate:"required"`
	Tid         int       `form:"tid"`
	UsageForm   int       `form:"usage_form" validate:"required"`
	ReceiveRate int       `form:"receive_rate"`
	ReceiveUri  string    `form:"receive_uri"`
	VipType     int       `form:"vip_type" validate:"required"`
	Stime       time.Time `form:"stime" validate:"required"`
	Etime       time.Time `form:"etime" validate:"required"`
	OperID      int       `json:"-"`
	OperName    string    `json:"-"`
}

//WelfareRes welfare type response
type WelfareRes struct {
	ID            int       `json:"id" gorm:"column:id"`
	Name          string    `json:"name" gorm:"column:welfare_name"`
	Desc          string    `json:"desc" gorm:"column:welfare_desc"`
	TID           int       `json:"tid" gorm:"column:tid"`
	HomepageUri   string    `json:"homepage_uri"`
	BackdropUri   string    `json:"backdrop_uri"`
	Recommend     int       `json:"recommend"`
	Rank          int       `json:"rank"`
	UsageForm     int       `json:"usage_form"`
	Stime         time.Time `json:"stime"`
	Etime         time.Time `json:"etime"`
	ReceiveRate   int       `json:"receive_rate"`
	ReceiveUri    string    `json:"receive_uri"`
	VipType       int       `json:"vip_type"`
	ReceivedCount int       `json:"received_count"`
	Count         int       `json:"count"`
}

//WelfareBatchRes welfare batch response
type WelfareBatchRes struct {
	ID            int       `json:"id"`
	Name          string    `json:"batch_name" gorm:"column:batch_name"`
	WID           int       `json:"wid" gorm:"column:wid"`
	Ver           int       `json:"ver"`
	OperId        int       `json:"oper_id"`
	OperName      string    `json:"oper_name"`
	Vtime         time.Time `json:"vtime"`
	Ctime         time.Time `json:"ctime"`
	ReceivedCount int       `json:"received_count"`
	Count         int       `json:"count"`
}

// TableName vip_welfare_type
func (*WelfareType) TableName() string {
	return "vip_welfare_type"
}

// TableName vip_welfare
func (*Welfare) TableName() string {
	return "vip_welfare"
}

// TableName vip_welfare_code_batch
func (*WelfareCodeBatch) TableName() string {
	return "vip_welfare_code_batch"
}

// TableName vip_welfare_code
func (*WelfareCode) TableName() string {
	return "vip_welfare_code"
}
