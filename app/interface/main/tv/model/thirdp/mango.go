package thirdp

import (
	"go-common/app/interface/main/tv/model"
	"go-common/library/time"
)

// MangoSn is the plus version of Dangbei Season
type MangoSn struct {
	SID int64 `json:"sid"`
	DBeiSeason
	OriginName string    `json:"origin_name"`
	Alias      string    `json:"alias"`
	Autorised  bool      `json:"autorised"`
	Mtime      time.Time `json:"mtime"`
	EpCover    string    `json:"ep_cover"`
}

// ToMangoSn transforms an seasonCMS to mangoSn
func ToMangoSn(s *model.SeasonCMS, mtime time.Time, autorised bool) *MangoSn {
	mSn := &MangoSn{
		DBeiSeason: *DBeiSn(s),
		SID:        s.SeasonID,
		OriginName: s.OriginName,
		Alias:      s.Alias,
		Mtime:      mtime,
		Autorised:  autorised,
	}
	mSn.SeasonID = nil
	return mSn
}

// MangoEP is mango ep structure
type MangoEP struct {
	model.EpCMS
	SeasonID  int64     `json:"sid"`
	Autorised bool      `json:"autorised"`
	Mtime     time.Time `json:"mtime"`
}

// MangoSnPage is mango sn page structure
type MangoSnPage struct {
	List  []*MangoSn      `json:"list"`
	Pager *model.IdxPager `json:"pager"`
}

// MangoEpPage is mango ep page structure
type MangoEpPage struct {
	SeasonID int64           `json:"sid"`
	List     []*MangoEP      `json:"list"`
	Pager    *model.IdxPager `json:"pager"`
}

// MangoArc is mango archive structure
type MangoArc struct {
	AVID      int64     `json:"avid"`
	Cover     string    `json:"cover"`
	Desc      string    `json:"desc"`
	Title     string    `json:"title"`
	PlayTime  string    `json:"play_time"`
	Category1 string    `json:"category_1"`
	Category2 string    `json:"category_2"`
	Autorised bool      `json:"autorised"`
	Mtime     time.Time `json:"mtime"`
}

// MangoVideo is mango video structure
type MangoVideo struct {
	CID       int64     `json:"cid"`
	Page      int       `json:"page"`
	Desc      string    `json:"desc"`
	Title     string    `json:"title"`
	Duration  int64     `json:"duration"`
	Autorised bool      `json:"autorised"`
	Mtime     time.Time `json:"mtime"`
}

// MangoArcPage is mango arc page structure
type MangoArcPage struct {
	List  []*MangoArc     `json:"list"`
	Pager *model.IdxPager `json:"pager"`
}

// MangoVideoPage is mango video page structure
type MangoVideoPage struct {
	AVID  int64           `json:"avid"`
	List  []*MangoVideo   `json:"list"`
	Pager *model.IdxPager `json:"pager"`
}

// RespSid is response structure for sid
type RespSid struct {
	Sid   int64
	Mtime time.Time
}

// PickSids picks sids from resp slices
func PickSids(resps []*RespSid) (sids []int64) {
	for _, v := range resps {
		sids = append(sids, v.Sid)
	}
	return
}

// ToMangoArc transforms an ArcCMS to MangoArc
func ToMangoArc(a *model.ArcCMS, mtime time.Time, cat1, cat2 string) *MangoArc {
	return &MangoArc{
		AVID:      a.AID,
		Cover:     a.Cover,
		Desc:      a.Content,
		Title:     a.Title,
		PlayTime:  a.Pubtime.Time().Format("2006-01-02"),
		Autorised: a.CanPlay(),
		Mtime:     mtime,
		Category1: cat1,
		Category2: cat2,
	}
}
