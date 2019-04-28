package model

import (
	"go-common/app/interface/main/dm2/model"
	tagmdl "go-common/app/interface/main/tag/model"
	accmdl "go-common/app/service/main/account/model"
	arcmdl "go-common/app/service/main/archive/api"
	ugcmdl "go-common/app/service/main/ugcpay/api/grpc/v1"
)

// View view data
type View struct {
	// archive data
	*arcmdl.Arc
	NoCache bool `json:"no_cache"`
	// video data pages
	Pages    []*arcmdl.Page         `json:"pages,omitempty"`
	Subtitle *Subtitle              `json:"subtitle"`
	Asset    *ugcmdl.AssetQueryResp `json:"asset,omitempty"`
}

// AssetRelation .
type AssetRelation struct {
	State int `json:"state"`
}

// Stat archive stat web struct
type Stat struct {
	Aid       int64       `json:"aid"`
	View      interface{} `json:"view"`
	Danmaku   int32       `json:"danmaku"`
	Reply     int32       `json:"reply"`
	Fav       int32       `json:"favorite"`
	Coin      int32       `json:"coin"`
	Share     int32       `json:"share"`
	Like      int32       `json:"like"`
	NowRank   int32       `json:"now_rank"`
	HisRank   int32       `json:"his_rank"`
	NoReprint int32       `json:"no_reprint"`
	Copyright int32       `json:"copyright"`
}

// Detail detail data
type Detail struct {
	View    *View
	Card    *Card
	Tags    []*tagmdl.Tag
	Reply   *ReplyHot
	Related []*arcmdl.Arc
}

// ArchiveUserCoins .
type ArchiveUserCoins struct {
	Multiply int64 `json:"multiply"`
}

// Subtitle dm subTitle.
type Subtitle struct {
	AllowSubmit bool            `json:"allow_submit"`
	List        []*SubtitleItem `json:"list"`
}

// SubtitleItem dm subTitle.
type SubtitleItem struct {
	*model.VideoSubtitle
	Author *accmdl.Info `json:"author"`
}

// TripleRes struct
type TripleRes struct {
	Like      bool  `json:"like"`
	Coin      bool  `json:"coin"`
	Fav       bool  `json:"fav"`
	Multiply  int64 `json:"multiply"`
	UpID      int64 `json:"-"`
	Anticheat bool  `json:"-"`
}

var (
	// StatAllowStates archive stat allow states
	statAllowStates = []int32{-9, -15, -30}
)

// CheckAllowState check archive stat allow state
func CheckAllowState(arc *arcmdl.Arc) bool {
	if arc.IsNormal() {
		return true
	}
	for _, allow := range statAllowStates {
		if arc.State == allow {
			return true
		}
	}
	return false
}

// FmtArc fmt grpc arc to archive3
func FmtArc(arc *arcmdl.Arc) (data *arcmdl.Arc) {
	data = &arcmdl.Arc{
		Aid:         arc.Aid,
		Videos:      arc.Videos,
		TypeID:      arc.TypeID,
		TypeName:    arc.TypeName,
		Copyright:   arc.Copyright,
		Pic:         arc.Pic,
		Title:       arc.Title,
		PubDate:     arc.PubDate,
		Ctime:       arc.Ctime,
		Desc:        arc.Desc,
		State:       arc.State,
		Access:      arc.Access,
		Attribute:   arc.Attribute,
		Tag:         arc.Tag,
		Tags:        arc.Tags,
		Duration:    arc.Duration,
		MissionID:   arc.MissionID,
		OrderID:     arc.OrderID,
		RedirectURL: arc.RedirectURL,
		Forward:     arc.Forward,
		Rights: arcmdl.Rights{
			Bp:        arc.Rights.Bp,
			Elec:      arc.Rights.Elec,
			Download:  arc.Rights.Download,
			Movie:     arc.Rights.Movie,
			Pay:       arc.Rights.Pay,
			HD5:       arc.Rights.HD5,
			NoReprint: arc.Rights.NoReprint,
			Autoplay:  arc.Rights.Autoplay,
			UGCPay:    arc.Rights.UGCPay,
		},
		Author: arcmdl.Author{
			Mid:  arc.Author.Mid,
			Name: arc.Author.Name,
			Face: arc.Author.Face,
		},
		Stat: arcmdl.Stat{
			Aid:     arc.Stat.Aid,
			View:    arc.Stat.View,
			Danmaku: arc.Stat.Danmaku,
			Reply:   arc.Stat.Reply,
			Fav:     arc.Stat.Fav,
			Coin:    arc.Stat.Coin,
			Share:   arc.Stat.Share,
			NowRank: arc.Stat.NowRank,
			HisRank: arc.Stat.HisRank,
			Like:    arc.Stat.Like,
			DisLike: arc.Stat.DisLike,
		},
		ReportResult: arc.ReportResult,
		Dynamic:      arc.Dynamic,
		FirstCid:     arc.FirstCid,
		Dimension: arcmdl.Dimension{
			Width:  arc.Dimension.Width,
			Height: arc.Dimension.Height,
			Rotate: arc.Dimension.Rotate,
		},
	}
	return
}
