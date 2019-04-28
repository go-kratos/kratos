package view

import (
	"go-common/app/interface/main/app-view/model/elec"
	"go-common/app/service/main/archive/model/archive"
	xtime "go-common/library/time"
)

// BnjMain is
type BnjMain struct {
	*archive.Archive3
	// now user
	ReqUser *ReqUser   `json:"req_user,omitempty"`
	Pages   []*Page    `json:"pages,omitempty"`
	Elec    *elec.Info `json:"elec,omitempty"`
	Relates []*BnjItem `json:"relates"`
	// player_icon
	PlayerIcon    *PlayerIcon `json:"player_icon,omitempty"`
	ElecBigText   string      `json:"elec_big_text"`
	ElecSmallText string      `json:"elec_small_text"`
}

// BnjList is
type BnjList struct {
	Item []*BnjItem `json:"list"`
}

// BnjItem is
type BnjItem struct {
	Aid       int64             `json:"aid"`
	Cid       int64             `json:"cid"`
	Tid       int32             `json:"tid"`
	Pic       string            `json:"pic"`
	Copyright int32             `json:"copyright"`
	PubDate   xtime.Time        `json:"pubdate"`
	IsAd      int               `json:"is_ad"`
	Title     string            `json:"title"`
	Desc      string            `json:"desc,omitempty"`
	Stat      archive.Stat3     `json:"stat,omitempty"`
	Duration  int64             `json:"duration,omitempty"`
	Author    archive.Author3   `json:"owner,omitempty"`
	Dimension archive.Dimension `json:"dimension,omitempty"`
	ReqUser   *ReqUser          `json:"req_user,omitempty"`
	Pages     []*Page           `json:"pages,omitempty"`
	Rights    archive.Rights3   `json:"rights"`
}
