package api

import (
	"hash/crc32"
	"strconv"
	"strings"

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

// IsNormal is
func (a *Arc) IsNormal() bool {
	return a.State >= StateOpen || a.State == StateForbidFixed
}

// RegionArc RegionArc
type RegionArc struct {
	Aid       int64
	Attribute int32
	Copyright int8
	PubDate   time.Time
}

// AllowShow AllowShow
func (ra *RegionArc) AllowShow() bool {
	return ra.attrVal(AttrBitNoWeb) == AttrNo && ra.attrVal(AttrBitNoMobile) == AttrNo
}

func (ra *RegionArc) attrVal(bit uint) int32 {
	return (ra.Attribute >> bit) & int32(1)
}

// AttrVal get attr val by bit.
func (a *Arc) AttrVal(bit uint) int32 {
	return (a.Attribute >> bit) & int32(1)
}

// FillDimension is
func (a *Arc) FillDimension(d string) {
	if d == "" || d == "0,0,0" {
		return
	}
	ds := strings.Split(d, ",")
	if len(ds) != 3 {
		return
	}
	var (
		width, height, rotate int64
		err                   error
	)
	if width, err = strconv.ParseInt(ds[0], 10, 64); err != nil {
		return
	}
	if height, err = strconv.ParseInt(ds[1], 10, 64); err != nil {
		return
	}
	if rotate, err = strconv.ParseInt(ds[2], 10, 64); err != nil {
		return
	}
	a.Dimension = Dimension{
		Width:  width,
		Height: height,
		Rotate: rotate,
	}
}

// FillDimension is
func (v *Page) FillDimension(d string) {
	if d == "" || d == "0,0,0" {
		return
	}
	ds := strings.Split(d, ",")
	if len(ds) != 3 {
		return
	}
	var (
		width, height, rotate int64
		err                   error
	)
	if width, err = strconv.ParseInt(ds[0], 10, 64); err != nil {
		return
	}
	if height, err = strconv.ParseInt(ds[1], 10, 64); err != nil {
		return
	}
	if rotate, err = strconv.ParseInt(ds[2], 10, 64); err != nil {
		return
	}
	v.Dimension = Dimension{
		Width:  width,
		Height: height,
		Rotate: rotate,
	}
}

// Fill file archive some field.
func (a *Arc) Fill() {
	a.Tags = _emptyTags
	a.Pic = coverURL(a.Pic)
	a.Rights.Bp = a.AttrVal(AttrBitAllowBp)
	a.Rights.Movie = a.AttrVal(AttrBitIsMovie)
	a.Rights.Pay = a.AttrVal(AttrBitBadgepay)
	a.Rights.HD5 = a.AttrVal(AttrBitHasHD5)
	a.Rights.NoReprint = a.AttrVal(AttrBitNoReprint)
	a.Rights.UGCPay = a.AttrVal(AttrBitUGCPay)
	a.Rights.IsCooperation = a.AttrVal(AttrBitIsCooperation)
	if a.FirstCid == 0 ||
		a.Access == AccessMember ||
		a.AttrVal(AttrBitIsPGC) == AttrYes ||
		a.AttrVal(AttrBitAllowBp) == AttrYes ||
		a.AttrVal(AttrBitBadgepay) == AttrYes ||
		a.AttrVal(AttrBitOverseaLock) == AttrYes ||
		a.AttrVal(AttrBitUGCPay) == AttrYes ||
		a.AttrVal(AttrBitLimitArea) == AttrYes {
		return
	}
	a.Rights.Autoplay = 1
}

// coverURL convert cover url to full url.
func coverURL(uri string) (cover string) {
	if uri == "" {
		cover = "http://static.hdslb.com/images/transparent.gif"
		return
	}
	cover = uri
	if strings.Index(uri, "http://") == 0 {
		return
	}
	if len(uri) >= 10 && uri[:10] == "/templets/" {
		return
	}
	if strings.HasPrefix(uri, "group1") {
		cover = "http://i0.hdslb.com/" + uri
		return
	}
	if pos := strings.Index(uri, "/uploads/"); pos != -1 && (pos == 0 || pos == 3) {
		cover = uri[pos+8:]
	}
	cover = strings.Replace(cover, "{IMG}", "", -1)
	cover = "http://i" + strconv.FormatInt(int64(crc32.ChecksumIEEE([]byte(cover)))%3, 10) + ".hdslb.com" + cover
	return
}

// FillStat file stat, check access.
func (a *Arc) FillStat() {
	if a.Access > 0 {
		a.Stat.View = 0
	}
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
