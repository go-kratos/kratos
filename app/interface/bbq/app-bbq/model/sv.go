package model

import (
	"go-common/library/time"
)

//SVBvcKey ..
type SVBvcKey struct {
	SVID            int64  `json:"svid"`
	Path            string `json:"path"`
	ResolutionRetio string `json:"resolution_retio"`
	CodeRate        int16  `json:"code_rate"`
	VideoCode       string `json:"video_code"`
	FileSize        int64  `json:"file_size"`
	Duration        int64  `json:"duration"`
}

// ParamScore 打分参数
type ParamScore struct {
	SVID  int64 `form:"svid" validate:"gt=0,required"`
	Score int64 `form:"score" validate:"gt=0,required"`
}

// ParamStatistic 统计参数
type ParamStatistic struct {
	SVIDs string `form:"svid" validate:"required"`
}

// SvInfo svList response
type SvInfo struct {
	SVID        int64     `json:"svid"`
	TID         int64     `json:"tid"`
	SubTID      int64     `json:"sub_tid"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	MID         int64     `json:"mid"`
	Duration    int64     `json:"duration"`
	Pubtime     time.Time `json:"pubtime"`
	Ctime       time.Time `json:"ctime"`
	AVID        int64     `json:"avid"`
	CID         int64     `json:"cid"`
	State       int16     `json:"state"`
	Original    int16     `json:"original"`
	From        int16     `json:"from"`
	VerID       int64     `json:"ver_id"`
	Ver         int64     `json:"ver"`
	Tag         string    `json:"tag"`
	CoverURL    string    `json:"cover_url"`
	CoverWidth  int       `json:"cover_width"`
	CoverHeight int       `json:"cover_height"`
}

// SvStInfo static info
type SvStInfo struct {
	SVID      int64 `json:"svid"`
	Play      int64 `json:"view"` //和上层的play重复，因此改成view
	Subtitles int64 `json:"subtitles"`
	Like      int64 `json:"like"`
	Share     int64 `json:"share"`
	Reply     int64 `json:"reply"`
	Report    int64 `json:"report"`
}

// SvTag SvTag struct
type SvTag struct {
	SVID  int64
	TagID int64
}
