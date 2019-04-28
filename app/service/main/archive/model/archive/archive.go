package archive

import (
	"go-common/app/service/main/archive/api"
	"go-common/library/time"
)

// 各属性地址见 http://syncsvn.bilibili.co/platform/doc/blob/master/archive/field/state.md

// all const
const (
	// open state
	StateOpen    = 0
	StateOrange  = 1
	AccessMember = int32(10000)
	// forbid state
	StateForbidWait       = -1
	StateForbidRecicle    = -2
	StateForbidPolice     = -3
	StateForbidLock       = -4
	StateForbidFixed      = -6
	StateForbidLater      = -7
	StateForbidAdminDelay = -10
	StateForbidXcodeFail  = -16
	StateForbidSubmit     = -30
	StateForbidUserDelay  = -40
	StateForbidUpDelete   = -100
	// copyright
	CopyrightUnknow   = int8(0)
	CopyrightOriginal = int8(1)
	CopyrightCopy     = int8(2)

	// attribute yes and no
	AttrYes = int32(1)
	AttrNo  = int32(0)
	// attribute bit
	AttrBitNoRank    = uint(0)
	AttrBitNoDynamic = uint(1)
	AttrBitNoWeb     = uint(2)
	AttrBitNoMobile  = uint(3)
	// AttrBitNoSearch    = uint(4)
	AttrBitOverseaLock = uint(5)
	AttrBitNoRecommend = uint(6)
	AttrBitNoReprint   = uint(7)
	AttrBitHasHD5      = uint(8)
	AttrBitIsPGC       = uint(9)
	AttrBitAllowBp     = uint(10)
	AttrBitIsBangumi   = uint(11)
	AttrBitIsPorder    = uint(12)
	AttrBitLimitArea   = uint(13)
	AttrBitAllowTag    = uint(14)
	// AttrBitIsFromArcApi  = uint(15)
	AttrBitJumpUrl       = uint(16)
	AttrBitIsMovie       = uint(17)
	AttrBitBadgepay      = uint(18)
	AttrBitUGCPay        = uint(22)
	AttrBitHasBGM        = uint(23)
	AttrBitIsCooperation = uint(24)
	AttrBitHasViewpoint  = uint(25)
	AttrBitHasArgument   = uint(26)
)

var (
	_emptyTags = []string{}
)

// AidPubTime aid's pubdate and copyright
type AidPubTime struct {
	Aid       int64     `json:"aid"`
	PubDate   time.Time `json:"pubdate"`
	Copyright int8      `json:"copyright"`
}

// ArchiveWithPlayer with first player info
type ArchiveWithPlayer struct {
	*Archive3
	PlayerInfo *PlayerInfo `json:"player_info,omitempty"`
}

// ArchiveWithBvc with first player info
type ArchiveWithBvc struct {
	*Archive3
	PlayerInfo *BvcVideoItem `json:"player_info,omitempty"`
}

// PlayerInfo player info
type PlayerInfo struct {
	Cid                int64                     `json:"cid"`
	ExpireTime         int64                     `json:"expire_time,omitempty"`
	FileInfo           map[int][]*PlayerFileInfo `json:"file_info"`
	SupportQuality     []int                     `json:"support_quality"`
	SupportFormats     []string                  `json:"support_formats"`
	SupportDescription []string                  `json:"support_description"`
	Quality            int                       `json:"quality"`
	URL                string                    `json:"url,omitempty"`
	VideoCodecid       uint32                    `json:"video_codecid"`
	VideoProject       bool                      `json:"video_project"`
	Fnver              int                       `json:"fnver"`
	Fnval              int                       `json:"fnval"`
	Dash               *ResponseDash             `json:"dash,omitempty"`
}

// PlayerFileInfo is
type PlayerFileInfo struct {
	TimeLength int64  `json:"timelength"`
	FileSize   int64  `json:"filesize"`
	Ahead      string `json:"ahead,omitempty"`
	Vhead      string `json:"vhead,omitempty"`
	URL        string `json:"url,omitempty"`
	Order      int64  `json:"order,omitempty"`
}

// ArcType arctype
type ArcType struct {
	ID   int16  `json:"id"`
	Pid  int16  `json:"pid"`
	Name string `json:"name"`
}

// Videoshot videoshot
type Videoshot struct {
	PvData string   `json:"pvdata"`
	XLen   int      `json:"img_x_len"`
	YLen   int      `json:"img_y_len"`
	XSize  int      `json:"img_x_size"`
	YSize  int      `json:"img_y_size"`
	Image  []string `json:"image"`
	Attr   int32    `json:"-"`
}

// RankArchives3 RankArchives3
type RankArchives3 struct {
	Archives []*api.Arc `json:"archives"`
	Count    int        `json:"count"`
}

// LikedArchives3 LikedArchives3
type LikedArchives3 struct {
	Archives []*Archive3 `json:"archives"`
	Count    int64
}

// RegionArc RegionArc
type RegionArc struct {
	Aid       int64
	Attribute int32
	Copyright int8
	PubDate   time.Time
}

// IsNormal is
func (a *Archive3) IsNormal() bool {
	return a.State >= StateOpen || a.State == StateForbidFixed
}

// AttrVal get attr val by bit.
func (a *Archive3) AttrVal(bit uint) int32 {
	return (a.Attribute >> bit) & int32(1)
}

func (ra *RegionArc) attrVal(bit uint) int32 {
	return (ra.Attribute >> bit) & int32(1)
}

type PGCPlayer struct {
	PlayerInfo *PlayerInfo `json:"player_info"`
	Aid        int64       `json:"aid"`
}

// BuildArchive3 build Archive3 with new proto
func BuildArchive3(a *api.Arc) (arc *Archive3) {
	if a == nil {
		return nil
	}
	arc = &Archive3{
		Aid:         a.Aid,
		Videos:      a.Videos,
		TypeID:      a.TypeID,
		TypeName:    a.TypeName,
		Copyright:   a.Copyright,
		Pic:         a.Pic,
		Title:       a.Title,
		PubDate:     a.PubDate,
		Ctime:       a.Ctime,
		Desc:        a.Desc,
		State:       a.State,
		Access:      a.Access,
		Attribute:   a.Attribute,
		Duration:    a.Duration,
		MissionID:   a.MissionID,
		OrderID:     a.OrderID,
		RedirectURL: a.RedirectURL,
		Forward:     a.Forward,
		Rights: Rights3{
			Bp:            a.Rights.Bp,
			Elec:          a.Rights.Elec,
			Download:      a.Rights.Download,
			Movie:         a.Rights.Movie,
			Pay:           a.Rights.Pay,
			HD5:           a.Rights.HD5,
			NoReprint:     a.Rights.NoReprint,
			Autoplay:      a.Rights.Autoplay,
			UGCPay:        a.Rights.UGCPay,
			IsCooperation: a.Rights.IsCooperation,
		},
		Author: Author3{
			Mid:  a.Author.Mid,
			Name: a.Author.Name,
			Face: a.Author.Face,
		},
		Stat: Stat3{
			Aid:     a.Stat.Aid,
			View:    a.Stat.View,
			Danmaku: a.Stat.Danmaku,
			Reply:   a.Stat.Reply,
			Fav:     a.Stat.Fav,
			Coin:    a.Stat.Coin,
			Share:   a.Stat.Share,
			NowRank: a.Stat.NowRank,
			HisRank: a.Stat.HisRank,
			Like:    a.Stat.Like,
			DisLike: a.Stat.DisLike,
		},
		ReportResult: a.ReportResult,
		Dynamic:      a.Dynamic,
		FirstCid:     a.FirstCid,
		Dimension: Dimension{
			Width:  a.Dimension.Width,
			Height: a.Dimension.Height,
			Rotate: a.Dimension.Rotate,
		},
	}
	return
}

// BuildView3 build View3 with new proto.
func BuildView3(a *api.Arc, pages []*api.Page) (v *View3) {
	if a == nil {
		return nil
	}
	var arc = &Archive3{
		Aid:         a.Aid,
		Videos:      a.Videos,
		TypeID:      a.TypeID,
		TypeName:    a.TypeName,
		Copyright:   a.Copyright,
		Pic:         a.Pic,
		Title:       a.Title,
		PubDate:     a.PubDate,
		Ctime:       a.Ctime,
		Desc:        a.Desc,
		State:       a.State,
		Access:      a.Access,
		Attribute:   a.Attribute,
		Duration:    a.Duration,
		MissionID:   a.MissionID,
		OrderID:     a.OrderID,
		RedirectURL: a.RedirectURL,
		Forward:     a.Forward,
		Rights: Rights3{
			Bp:            a.Rights.Bp,
			Elec:          a.Rights.Elec,
			Download:      a.Rights.Download,
			Movie:         a.Rights.Movie,
			Pay:           a.Rights.Pay,
			HD5:           a.Rights.HD5,
			NoReprint:     a.Rights.NoReprint,
			Autoplay:      a.Rights.Autoplay,
			UGCPay:        a.Rights.UGCPay,
			IsCooperation: a.Rights.IsCooperation,
		},
		Author: Author3{
			Mid:  a.Author.Mid,
			Name: a.Author.Name,
			Face: a.Author.Face,
		},
		Stat: Stat3{
			Aid:     a.Stat.Aid,
			View:    a.Stat.View,
			Danmaku: a.Stat.Danmaku,
			Reply:   a.Stat.Reply,
			Fav:     a.Stat.Fav,
			Coin:    a.Stat.Coin,
			Share:   a.Stat.Share,
			NowRank: a.Stat.NowRank,
			HisRank: a.Stat.HisRank,
			Like:    a.Stat.Like,
			DisLike: a.Stat.DisLike,
		},
		ReportResult: a.ReportResult,
		Dynamic:      a.Dynamic,
		FirstCid:     a.FirstCid,
		Dimension: Dimension{
			Width:  a.Dimension.Width,
			Height: a.Dimension.Height,
			Rotate: a.Dimension.Rotate,
		},
	}
	for _, s := range a.StaffInfo {
		arc.StaffInfo = append(arc.StaffInfo, &StaffInfo{
			Mid:   s.Mid,
			Title: s.Title,
		})
	}
	var viewPages []*Page3
	for _, v := range pages {
		viewPages = append(viewPages, &Page3{
			Cid:      v.Cid,
			Page:     v.Page,
			From:     v.From,
			Part:     v.Part,
			Duration: v.Duration,
			Vid:      v.Vid,
			Desc:     v.Desc,
			WebLink:  v.WebLink,
			Dimension: Dimension{
				Width:  v.Dimension.Width,
				Height: v.Dimension.Height,
				Rotate: v.Dimension.Rotate,
			},
		})
	}
	v = &View3{
		Archive3: arc,
		Pages:    viewPages,
	}
	return
}
