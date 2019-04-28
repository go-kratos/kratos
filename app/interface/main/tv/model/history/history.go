package history

import (
	hismodel "go-common/app/interface/main/history/model"
	"go-common/app/interface/main/tv/model"
)

// HisRes is the history resource model
type HisRes struct {
	Mid          int64              `json:"mid,omitempty"`
	Oid          int64              `json:"oid"` //
	Sid          int64              `json:"season_id,omitempty"`
	Epid         int64              `json:"epid,omitempty"`
	Cid          int64              `json:"cid,omitempty"`
	Business     string             `json:"-"`
	DT           int8               `json:"dt,omitempty"`
	Pro          int64              `json:"pro"`
	PageDuration int64              `json:"duration"`
	Unix         int64              `json:"view_at"`
	Type         int                `json:"type"`              // 1=pgc, 2=ugc
	Title        string             `json:"title"`             // common
	Cover        string             `json:"cover"`             // common
	Page         *HisPage           `json:"page,omitempty"`    // ugc page
	EPMeta       *HisEP             `json:"bangumi,omitempty"` // pgc ep
	CornerMark   *model.SnVipCorner `json:"cornermark"`
}

// HisEP is history EP struct
type HisEP struct {
	EPID      int64  `json:"ep_id"`
	Cover     string `json:"cover"`
	Title     string `json:"title"`
	LongTitle string `json:"long_title"`
}

// HisPage is history page struct
type HisPage struct {
	CID  int64  `json:"cid"`
	Part string `json:"part"`
	Page int    `json:"page"`
}

// HisMC is history structure in MC
type HisMC struct {
	MID        int64
	Res        []*HisRes
	LastViewAt int64 // timestamp the first view_at of cursor
}

// Dur is duration struct
type Dur struct {
	Oid      int64
	Duration int64
}

// RespCacheHis is the the response of cacheHis function
type RespCacheHis struct {
	Filtered []*HisRes
	Res      []*hismodel.Resource
	UseCache bool
}

// ReqCombineHis is the request for the combineHis function
type ReqCombineHis struct {
	Mid    int64
	OriRes []*hismodel.Resource
	OkSids map[int64]int
	OkAids map[int64]int
}
