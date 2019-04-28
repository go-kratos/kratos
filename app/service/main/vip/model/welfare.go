package model

import (
	"go-common/library/time"
	xtime "go-common/library/time"
)

// ArgWelfareList args for welfare list.
type ArgWelfareList struct {
	Tid       int64      `form:"tid"`
	Recommend int64      `form:"recommend"`
	Ps        int64      `form:"ps"`
	Pn        int64      `form:"pn"`
	NowTime   xtime.Time `form:"-"`
}

// WelfareListResp response for welfare list.
type WelfareListResp struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	HomepageUri string `json:"homepage_uri"`
	BackdropUri string `json:"backdrop_uri"`
	Tid         int32  `json:"tid"`
	Rank        int32  `json:"rank"`
}

// WelfareTypeListResp response for welfare type list.
type WelfareTypeListResp struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// ArgWelfareInfo args for welfare info.
type ArgWelfareInfo struct {
	ID  int64 `form:"id"`
	MID int64 `form:"mid"`
}

// WelfareInfoResp response for welfare info.
type WelfareInfoResp struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Desc        string    `json:"desc"`
	ReceiveRate int       `json:"receive_rate"`
	HomepageUri string    `json:"homepage_uri"`
	BackdropUri string    `json:"backdrop_uri"`
	Finished    bool      `json:"finished"`
	Received    bool      `json:"received"`
	UsageForm   int32     `json:"usage_form"`
	VipType     int64     `json:"vip_type"`
	Stime       time.Time `json:"stime"`
	Etime       time.Time `json:"etime"`
}

// ArgWelfareReceive args for welfare receive.
type ArgWelfareReceive struct {
	Wid int64 `form:"wid"`
	Mid int64 `form:"mid"`
}

// WelfareReceiveResp response for welfare receive.
type WelfareReceiveResp struct {
}

// WelfareBatchResp response for welfare batch.
type WelfareBatchResp struct {
	Id            int       `json:"id"`
	ReceivedCount int       `json:"received_count"`
	Count         int       `json:"count"`
	Vtime         time.Time `json:"vtime"`
}

// ReceivedCodeResp response for welfare code.
type ReceivedCodeResp struct {
	ID    int       `json:"id"`
	Mtime time.Time `json:"mtime"`
}

// UnReceivedCodeResp response for welfare unreceive.
type UnReceivedCodeResp struct {
	Id   int    `json:"id"`
	Bid  int    `json:"bid"`
	Code string `json:"code"`
}

// ReceiveRecordResp response for welfare record.
type ReceiveRecordResp struct {
	Id        int `json:"id"`
	Mid       int `json:"mid"`
	Wid       int `json:"wid"`
	MonthYear int `json:"month_year"`
	Count     int `json:"count"`
}

// MyWelfareResp response for my welfare.
type MyWelfareResp struct {
	Wid        int32     `json:"wid"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	UsageForm  int32     `json:"usage_form"`
	ReceiveUri string    `json:"receive_uri"`
	Code       string    `json:"code"`
	Expired    bool      `json:"expired"`
	Stime      time.Time `json:"stime"`
	Etime      time.Time `json:"etime"`
}
