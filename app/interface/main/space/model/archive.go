package model

import (
	v1 "go-common/app/service/main/archive/api"
	"go-common/library/time"
)

// UpArcStat up archives stat struct.
type UpArcStat struct {
	View  int64 `json:"view"`
	Reply int64 `json:"reply"`
	Dm    int64 `json:"dm"`
	Fans  int64 `json:"fans"`
}

// ArchiveReason archive with reason struct.
type ArchiveReason struct {
	*v1.Arc
	Reason string `json:"reason"`
}

// SearchRes search res data.
type SearchRes struct {
	TList map[string]*SearchTList `json:"tlist"`
	VList []*SearchVList          `json:"vlist"`
}

// SearchTList search cate list.
type SearchTList struct {
	Tid   int64  `json:"tid"`
	Count int64  `json:"count"`
	Name  string `json:"name"`
}

// SearchVList video list.
type SearchVList struct {
	Comment      int64       `json:"comment"`
	TypeID       int64       `json:"typeid"`
	Play         interface{} `json:"play"`
	Pic          string      `json:"pic"`
	SubTitle     string      `json:"subtitle"`
	Description  string      `json:"description"`
	Copyright    string      `json:"copyright"`
	Title        string      `json:"title"`
	Review       int64       `json:"review"`
	Author       string      `json:"author"`
	Mid          int64       `json:"mid"`
	Created      string      `json:"created"`
	Length       string      `json:"length"`
	VideoReview  int64       `json:"video_review"`
	Aid          int64       `json:"aid"`
	HideClick    bool        `json:"hide_click"`
	IsPay        int         `json:"is_pay"`
	IsUnionVideo int         `json:"is_union_video"`
}

// UpArc up archive struct
type UpArc struct {
	Count int64      `json:"count"`
	List  []*ArcItem `json:"list"`
}

// ArcItem space archive item.
type ArcItem struct {
	Aid      int64  `json:"aid"`
	Pic      string `json:"pic"`
	Title    string `json:"title"`
	Duration int64  `json:"duration"`
	Author   struct {
		Mid  int64  `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"author"`
	Stat struct {
		View    interface{} `json:"view"`
		Danmaku int32       `json:"danmaku"`
		Reply   int32       `json:"reply"`
		Fav     int32       `json:"favorite"`
		Coin    int32       `json:"coin"`
		Share   int32       `json:"share"`
		Like    int32       `json:"like"`
	} `json:"stat"`
	Rights  v1.Rights `json:"rights"`
	Pubdate time.Time `json:"pubdate"`
}

// FromArc from archive to space act item.
func (ac *ArcItem) FromArc(arc *v1.Arc) {
	ac.Aid = arc.Aid
	ac.Pic = arc.Pic
	ac.Title = arc.Title
	ac.Duration = arc.Duration
	ac.Author.Mid = arc.Author.Mid
	ac.Author.Name = arc.Author.Name
	ac.Author.Face = arc.Author.Face
	ac.Stat.View = arc.Stat.View
	if arc.Access >= 10000 {
		ac.Stat.View = "--"
	}
	ac.Stat.Danmaku = arc.Stat.Danmaku
	ac.Stat.Reply = arc.Stat.Reply
	ac.Stat.Fav = arc.Stat.Fav
	ac.Stat.Share = arc.Stat.Share
	ac.Stat.Like = arc.Stat.Like
	ac.Pubdate = arc.PubDate
	ac.Rights = arc.Rights
}
