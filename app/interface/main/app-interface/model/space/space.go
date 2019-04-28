package space

import (
	"strconv"

	"encoding/json"
	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/audio"
	"go-common/app/interface/main/app-interface/model/bangumi"
	"go-common/app/interface/main/app-interface/model/community"
	"go-common/app/interface/main/app-interface/model/elec"
	"go-common/app/interface/main/app-interface/model/favorite"
	tag "go-common/app/interface/main/tag/model"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

// Space struct
type Space struct {
	Relation  int             `json:"relation"`
	Medal     int             `json:"medal,omitempty"`
	Attention uint32          `json:"attention,omitempty"`
	Setting   *Setting        `json:"setting,omitempty"`
	Tab       *Tab            `json:"tab,omitempty"`
	Card      *Card           `json:"card,omitempty"`
	Space     *Mob            `json:"images,omitempty"`
	Shop      *Shop           `json:"shop,omitempty"`
	Live      json.RawMessage `json:"live,omitempty"`
	Elec      *elec.Info      `json:"elec,omitempty"`
	Archive   *ArcList        `json:"archive,omitempty"`
	Article   *ArticleList    `json:"article,omitempty"`
	Clip      *ClipList       `json:"clip,omitempty"`
	Album     *AlbumList      `json:"album,omitempty"`
	Favourite *FavList        `json:"favourite,omitempty"`
	Season    *BangumiList    `json:"season,omitempty"`
	CoinArc   *ArcList        `json:"coin_archive,omitempty"`
	LikeArc   *ArcList        `json:"like_archive,omitempty"`
	Audios    *AudioList      `json:"audios,omitempty"`
	Community *CommuList      `json:"community,omitempty"`
}

// Card struct
type Card struct {
	Mid            string        `json:"mid"`
	Name           string        `json:"name"`
	Approve        bool          `json:"approve"`
	Sex            string        `json:"sex"`
	Rank           string        `json:"rank"`
	Face           string        `json:"face"`
	DisplayRank    string        `json:"DisplayRank"`
	Regtime        int64         `json:"regtime"`
	Spacesta       int           `json:"spacesta"`
	Birthday       string        `json:"birthday"`
	Place          string        `json:"place"`
	Description    string        `json:"description"`
	Article        int           `json:"article"`
	Attentions     []int64       `json:"attentions"`
	Fans           int           `json:"fans"`
	Friend         int           `json:"friend"`
	Attention      int           `json:"attention"`
	Sign           string        `json:"sign"`
	LevelInfo      LevelInfo     `json:"level_info"`
	Pendant        PendantInfo   `json:"pendant"`
	Nameplate      NameplateInfo `json:"nameplate"`
	OfficialVerify OfficialInfo  `json:"official_verify"`
	Vip            struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
	FansGroup  int  `json:"fans_group,omitempty"`
	Audio      int  `json:"audio,omitempty"`
	FansUnread bool `json:"fans_unread,omitempty"`
}

// Mob struct
type Mob struct {
	ImgURL string `json:"imgUrl"`
}

// Shop struct
type Shop struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// LevelInfo struct
type LevelInfo struct {
	Cur     int32       `json:"current_level"`
	Min     int32       `json:"current_min"`
	NowExp  int32       `json:"current_exp"`
	NextExp interface{} `json:"next_exp"`
}

// PendantInfo struct
type PendantInfo struct {
	Pid    int    `json:"pid"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Expire int    `json:"expire"`
}

// NameplateInfo struct
type NameplateInfo struct {
	Nid        int    `json:"nid"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ImageSmall string `json:"image_small"`
	Level      string `json:"level"`
	Condition  string `json:"condition"`
}

// OfficialInfo struct
type OfficialInfo struct {
	Type  int8   `json:"type"`
	Desc  string `json:"desc"`
	Role  int8   `json:"role"`
	Title string `json:"title"`
}

// Setting struct
type Setting struct {
	Channel    int `json:"channel,omitempty"`
	FavVideo   int `json:"fav_video"`
	CoinsVideo int `json:"coins_video"`
	LikesVideo int `json:"likes_video"`
	Bangumi    int `json:"bangumi"`
	PlayedGame int `json:"played_game"`
	Groups     int `json:"groups"`
}

// TagList struct
type TagList struct {
	Count int        `json:"count"`
	Tags  []*tag.Tag `json:"item"`
}

// ArcList struct
type ArcList struct {
	Count int        `json:"count"`
	Item  []*ArcItem `json:"item"`
}

// ArticleList struct
type ArticleList struct {
	Count      int             `json:"count"`
	Item       []*ArticleItem  `json:"item"`
	ListsCount int             `json:"lists_count"`
	Lists      []*article.List `json:"lists"`
}

// CommuList struct
type CommuList struct {
	Count int         `json:"count"`
	Item  []*CommItem `json:"item"`
}

// FavList struct
type FavList struct {
	Count int                `json:"count"`
	Item  []*favorite.Folder `json:"item"`
}

// BangumiList struct
type BangumiList struct {
	Count int            `json:"count"`
	Item  []*BangumiItem `json:"item"`
}

// AudioList struct
type AudioList struct {
	Count int          `json:"count"`
	Item  []*AudioItem `json:"item"`
}

// ClipList struct
type ClipList struct {
	Count  int     `json:"count"`
	More   int     `json:"has_more"`
	Offset int     `json:"next_offset"`
	Item   []*Item `json:"item"`
}

// AlbumList struct
type AlbumList struct {
	Count  int     `json:"count"`
	More   int     `json:"has_more"`
	Offset int     `json:"next_offset"`
	Item   []*Item `json:"item"`
}

// ArcItem struct
type ArcItem struct {
	Title    string `json:"title"`
	TypeName string `json:"tname"`
	Cover    string `json:"cover"`
	URI      string `json:"uri"`
	Param    string `json:"param"`
	Goto     string `json:"goto"`
	Length   string `json:"length"`
	Duration int64  `json:"duration"`
	// av
	Play    int        `json:"play"`
	Danmaku int        `json:"danmaku"`
	CTime   xtime.Time `json:"ctime"`
	UGCPay  int32      `json:"ugc_pay"`
}

// ArticleItem struct
type ArticleItem struct {
	*article.Meta
	URI   string `json:"uri"`
	Param string `json:"param"`
	Goto  string `json:"goto"`
}

// BangumiItem struct
type BangumiItem struct {
	Title         string     `json:"title"`
	Cover         string     `json:"cover"`
	URI           string     `json:"uri"`
	Param         string     `json:"param"`
	Goto          string     `json:"goto"`
	Finish        int8       `json:"finish"`
	Index         string     `json:"index"`
	MTime         xtime.Time `json:"mtime"`
	NewestEpIndex string     `json:"newest_ep_index"`
	IsStarted     int        `json:"is_started"`
	IsFinish      string     `json:"is_finish"`
	NewestEpID    string     `json:"newest_ep_id"`
	TotalCount    string     `json:"total_count"`
	Attention     string     `json:"attention"`
}

// CommItem struct
type CommItem struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Desc           string `json:"desc"`
	Thumb          string `json:"thumb"`
	PostCount      int    `json:"post_count"`
	MemberCount    int    `json:"member_count"`
	PostNickname   string `json:"post_nickname"`
	MemberNickname string `json:"member_nickname"`
}

// AudioItem struct
type AudioItem struct {
	ID       int64      `json:"id"`
	Aid      int64      `json:"aid"`
	UID      int64      `json:"uid"`
	Title    string     `json:"title"`
	Cover    string     `json:"cover"`
	Author   string     `json:"author"`
	Schema   string     `json:"schema"`
	Duration int64      `json:"duration"`
	Play     int        `json:"play"`
	Reply    int        `json:"reply"`
	IsOff    int        `json:"isOff"`
	AuthType int        `json:"authType"`
	CTime    xtime.Time `json:"ctime"`
}

// FromSeason func
func (i *BangumiItem) FromSeason(b *bangumi.Season) {
	i.Title = b.Title
	i.Cover = b.Cover
	i.Goto = model.GotoBangumi
	i.Param = b.SeasonID
	i.URI = model.FillURI(model.GotoBangumiWeb, b.SeasonID, nil)
	i.IsStarted = b.IsStarted
	if b.IsFinish == "1" {
		i.Finish = 1
	}
	i.NewestEpIndex = b.NewestEpIndex
	i.TotalCount = b.TotalCount
	if b.UserSeason != nil {
		i.Attention = b.UserSeason.Attention
	}
}

// FromCoinArc func
func (i *ArcItem) FromCoinArc(a *api.Arc) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Param = strconv.FormatInt(int64(a.Aid), 10)
	i.URI = model.FillURI(model.GotoAv, i.Param, nil)
	i.Goto = model.GotoAv
	i.Danmaku = int(a.Stat.Danmaku)
	i.Duration = a.Duration
	i.CTime = a.PubDate
	i.Play = int(a.Stat.View)
}

// FromLikeArc fun
func (i *ArcItem) FromLikeArc(a *api.Arc) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Param = strconv.FormatInt(int64(a.Aid), 10)
	i.URI = model.FillURI(model.GotoAv, i.Param, nil)
	i.Goto = model.GotoAv
	i.Danmaku = int(a.Stat.Danmaku)
	i.Duration = a.Duration
	i.CTime = a.PubDate
	i.Play = int(a.Stat.View)
}

// FromArticle func
func (i *ArticleItem) FromArticle(a *article.Meta) {
	i.Meta = a
	i.Param = strconv.FormatInt(int64(a.ID), 10)
	i.URI = model.FillURI(model.GotoArticle, i.Param, nil)
	i.Goto = model.GotoArticle

}

// FromArc func
func (i *ArcItem) FromArc(c *api.Arc) {
	i.Title = c.Title
	i.Cover = c.Pic
	i.TypeName = c.TypeName
	i.Param = strconv.FormatInt(int64(c.Aid), 10)
	i.URI = model.FillURI(model.GotoAv, i.Param, nil)
	i.Goto = model.GotoAv
	i.Danmaku = int(c.Stat.Danmaku)
	i.CTime = c.PubDate
	i.Duration = c.Duration
	i.Play = int(c.Stat.View)
	i.UGCPay = c.Rights.UGCPay
}

// FromCommunity func
func (i *CommItem) FromCommunity(c *community.Community) {
	i.ID = c.ID
	i.Name = c.Name
	i.Desc = c.Desc
	i.Thumb = c.Thumb
	i.PostCount = c.PostCount
	i.MemberCount = c.MemberCount
	i.PostNickname = c.PostNickname
	i.MemberNickname = c.MemberNickname
}

// FromAudio func
func (i *AudioItem) FromAudio(a *audio.Audio) {
	i.ID = a.ID
	i.Aid = a.Aid
	i.UID = a.UID
	i.Title = a.Title
	i.Cover = a.Cover
	i.Author = a.Author
	i.Schema = a.Schema
	i.Duration = a.Duration
	i.Play = a.Play
	i.Reply = a.Reply
	i.IsOff = a.IsOff
	i.AuthType = a.AuthType
	i.CTime = a.CTime
}
