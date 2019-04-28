package model

import "go-common/library/time"

const (
	//WatermarkWhite 水印白名单
	WatermarkWhite = 1
	//WatermarkDefault 水印默认值
	WatermarkDefault = 0
	//OrderDesc 降序
	OrderDesc = 1
)

// WaterMarkList def.
type WaterMarkList struct {
	ID           string    `form:"id" json:"id"`
	Epid         string    `form:"epid" json:"epid"`
	SeasonID     string    `form:"season_id" json:"season_id"`
	Category     string    `form:"category" json:"category"`
	SeasonTitle  string    `form:"season_title" json:"season_title"`
	ContentTitle string    `form:"content_title" json:"content_title"`
	MarkTime     time.Time `form:"mark_time" json:"mark_time"`
}

// WaterMarkOne is used for only selecting some field from gorm query
type WaterMarkOne struct {
	ID       string `form:"id" json:"id"`
	Mark     uint8  `form:"mark" json:"mark"`
	MarkTime string `form:"mark_time" json:"mark_time"`
}

// WaterMarkListPager is used for return items and pager info
type WaterMarkListPager struct {
	Items []*WaterMarkList `json:"items"`
	Page  *Page            `json:"page"`
}

// WaterMarkListParam is use for watermarklist function query param valid
type WaterMarkListParam struct {
	Order    uint8  `form:"id" json:"order" default:"1"`
	EpID     string `form:"epid" json:"epid"`
	SeasonID string `form:"season_id" json:"season_id"`
	Category string `form:"category" json:"category"`
	Pn       int    `form:"pn" json:"pn;Min(1)" default:"1"`
	Ps       int    `form:"ps" json:"ps;Min(1)" default:"20"`
}

// TransReq is the request for transcode consulting
type TransReq struct {
	Order    int    `form:"order" default:"1"` // 1=desc,2=asc
	EpID     int64  `form:"epid"`
	SeasonID int64  `form:"season_id"`
	Title    string `form:"title"`
	Category int    `form:"category" validate:"min=0,max=5"`
	Status   string `form:"status"`
	Pn       int    `form:"pn" default:"1"`
}

// TransReply is the response for transList
type TransReply struct {
	EpID       int64  `json:"epid"`
	SeasonID   int64  `json:"season_id"`
	Category   string `json:"category"`
	Etitle     string `json:"etitle"`
	Stitle     string `json:"stitle"`
	Transcoded int    `json:"transcoded"`
	ApplyTime  string `json:"apply_time"`
	MarkTime   string `json:"mark_time"`
}

// TransPager is used for return items and pager info
type TransPager struct {
	Items   []*TransReply `json:"items"`
	Page    *Page         `json:"page"`
	CountSn int           `json:"count_sn"`
}

// AddEpIDResp is for function addEpID to return success and not exist and invalid values
type AddEpIDResp struct {
	Succ     []int64
	NotExist []int64
	Invalids []int64
}

// TableName select watermark list
func (a WaterMarkList) TableName() string {
	return "tv_content"
}

//TableName only select watermark one
func (a WaterMarkOne) TableName() string {
	return "tv_content"
}
