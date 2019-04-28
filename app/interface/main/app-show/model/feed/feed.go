package feed

import (
	"encoding/json"
	"strconv"

	clive "go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/card"
	"go-common/app/interface/main/app-show/model/dislike"
	"go-common/app/interface/main/app-show/model/tag"
	"go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

type Item struct {
	Title       string      `json:"title,omitempty"`
	Subtitle    string      `json:"subtitle,omitempty"`
	Cover       string      `json:"cover,omitempty"`
	URI         string      `json:"uri,omitempty"`
	Redirect    string      `json:"redirect,omitempty"`
	RedirectURI string      `json:"redirect_uri,omitempty"`
	Param       string      `json:"param,omitempty"`
	Goto        string      `json:"goto,omitempty"`
	ViewType    string      `json:"view_type,omitempty"`
	Kind        string      `json:"kind,omitempty"`
	Desc        string      `json:"desc,omitempty"`
	Play        int         `json:"play,omitempty"`
	Danmaku     int         `json:"danmaku,omitempty"`
	Reply       int         `json:"reply,omitempty"`
	Fav         int         `json:"favorite,omitempty"`
	Coin        int         `json:"coin,omitempty"`
	Share       int         `json:"share,omitempty"`
	Like        int         `json:"like,omitempty"`
	Count       int         `json:"count,omitempty"`
	Status      int8        `json:"status,omitempty"`
	Type        int8        `json:"type,omitempty"`
	Badge       string      `json:"badge,omitempty"`
	StatType    int8        `json:"stat_type,omitempty"`
	RcmdReason  *RcmdReason `json:"rcmd_reason,omitempty"`
	Item        []*Item     `json:"item,omitempty"`
	CornerMark  int8        `json:"-"`
	CardStyle   int8        `json:"-"`
	RcmdContent string      `json:"-"`
	// sortedset index
	Idx int64 `json:"idx,omitempty"`
	// av info
	Cid             int64                     `json:"cid,omitempty"`
	Rid             int16                     `json:"tid,omitempty"`
	TName           string                    `json:"tname,omitempty"`
	Tag             *Tag                      `json:"tag,omitempty"`
	DisklikeReasons []*dislike.DisklikeReason `json:"dislike_reasons,omitempty"`
	PTime           xtime.Time                `json:"ctime,omitempty"`
	Autoplay        int32                     `json:"autoplay,omitempty"`
	// av stat
	Duration int64 `json:"duration,omitempty"`
	// upper
	Mid      int64         `json:"mid,omitempty"`
	Name     string        `json:"name,omitempty"`
	Face     string        `json:"face,omitempty"`
	IsAtten  int8          `json:"is_atten,omitempty"`
	Fans     int64         `json:"fans,omitempty"`
	RecCnt   int           `json:"recent_count,omitempty"`
	Recent   []*Item       `json:"recent,omitempty"`
	Official *OfficialInfo `json:"official,omitempty"`
	// live
	Online int    `json:"online,omitempty"`
	Area   string `json:"area,omitempty"`
	AreaID int    `json:"area_id,omitempty"`
	Area2  *Area2 `json:"area2,omitempty"`
	// bangumi
	Index       string `json:"index,omitempty"`
	IndexTitle  string `json:"index_title,omitempty"`
	CoverMark   string `json:"cover_mark,omitempty"`
	Finish      bool   `json:"finish,omitempty"`
	LatestIndex string `json:"last_index,omitempty"`
	// bangumi recommend
	Updates int `json:"updates,omitempty"`
	// live or bangumi
	From int8 `json:"from,omitempty"`
	// adviertisement
	RequestID  string          `json:"request_id,omitempty"`
	CreativeID int64           `json:"creative_id,omitempty"`
	SrcID      int             `json:"src_id,omitempty"`
	IsAd       bool            `json:"is_ad,omitempty"`
	IsAdLoc    bool            `json:"is_ad_loc,omitempty"`
	AdCb       string          `json:"ad_cb,omitempty"`
	ShowURL    string          `json:"show_url,omitempty"`
	ClickURL   string          `json:"click_url,omitempty"`
	ClientIP   string          `json:"client_ip,omitempty"`
	CmMark     int64           `json:"cm_mark,omitempty"`
	AdIndex    int             `json:"ad_index,omitempty"`
	Extra      json.RawMessage `json:"extra,omitempty"`
	CardIndex  int             `json:"card_index,omitempty"`
	// activity
	STime string `json:"stime,omitempty"`
	ETime string `json:"etime,omitempty"`
	// tag
	Tags []*tag.Tag `json:"tags,omitempty"`
	// rank
	Cover1 string `json:"cover1,omitempty"`
	Cover2 string `json:"cover2,omitempty"`
	Cover3 string `json:"cover3,omitempty"`
	// banner`
	Hash string `json:"hash,omitempty"`
	// upper article
	Covers    []string  `json:"covers,omitempty"`
	Temple    int       `json:"temple,omitempty"`
	Category  *Category `json:"category,omitempty"`
	BannerURL string    `json:"banner_url,omitempty"`
	// game download
	GameDownloadButton *GameDownloadButton `json:"button,omitempty"`
	Download           int                 `json:"download,omitempty"`
	BigCover           string              `json:"big_cover,omitempty"`
	// special
	HideBadge bool    `json:"hide_badge,omitempty"`
	Ratio     float64 `json:"ratio,omitempty"`
	// shopping
	City  string `json:"city,omitempty"`
	PType string `json:"ptype,omitempty"`
	Price string `json:"price,omitempty"`
	// news
	Content string `json:"content,omitempty"`
	// bigdata source
	Source    string          `json:"-"`
	AvFeature json.RawMessage `json:"-"`
	// common
	GotoOrg  string `json:"-"`
	FromType string `json:"from_type,omitempty"`
	Pos      int    `json:"-"`
	Score    string `json:"score,omitempty"`
}

type Tag struct {
	TagID   int64     `json:"tag_id,omitempty"`
	TagName string    `json:"tag_name,omitempty"`
	IsAtten int8      `json:"is_atten,omitempty"`
	Count   *TagCount `json:"count,omitempty"`
	Name    string    `json:"name,omitempty"`
	URI     string    `json:"uri,omitempty"`
	//channel
	ID   int64  `json:"id,omitempty"`
	Face string `json:"face,omitempty"`
	Fans int64  `json:"fans,omitempty"`
}

type RcmdReason struct {
	ID           int    `json:"id,omitempty"`
	Content      string `json:"content,omitempty"`
	BgColor      string `json:"bg_color,omitempty"`
	IconLocation string `json:"icon_location,omitempty"`
	Message      string `json:"message,omitempty"`
}

type TagCount struct {
	Atten int `json:"atten,omitempty"`
}

type Category struct {
	ID       int64     `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Children *Category `json:"children,omitempty"`
}

type GameDownloadButton struct {
	Name        string `json:"name,omitempty"`
	URI         string `json:"uri,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

type Area2 struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Children *Area2 `json:"children,omitempty"`
}

type OfficialInfo struct {
	Role  int8   `json:"role,omitempty"`
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

func (i *Item) fillArcStat(a *api.Arc) {
	if a.Access == 0 {
		i.Play = int(a.Stat.View)
	}
	i.Danmaku = int(a.Stat.Danmaku)
	i.Reply = int(a.Stat.Reply)
	i.Fav = int(a.Stat.Fav)
	i.Coin = int(a.Stat.Coin)
	i.Share = int(a.Stat.Share)
	i.Like = int(a.Stat.Like)
}

func (i *Item) FromPlayerAv(a *api.Arc, uri string) {
	if i.Title == "" {
		i.Title = a.Title
	}
	if i.Cover == "" {
		i.Cover = model.CoverURLHTTPS(a.Pic)
	} else {
		i.Cover = model.CoverURLHTTPS(i.Cover)
	}
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.Goto = model.GotoAv
	i.URI = uri
	if i.URI == "" {
		i.URI = model.FillURI(i.Goto, i.Param, model.AvHandler(a))
	}
	i.Rid = int16(a.TypeID)
	i.TName = a.TypeName
	// i.Desc = a.Desc
	i.fillArcStat(a)
	i.Duration = a.Duration
	i.Mid = a.Author.Mid
	i.Name = a.Author.Name
	i.Face = a.Author.Face
	i.PTime = a.PubDate
	i.Autoplay = a.Rights.Autoplay
	i.Cid = a.FirstCid
	// TODO
	// if a.Stat.Like > 0 && a.Stat.DisLike > 0 {
	// 	percent := int(a.Stat.Like / (a.Stat.Like + a.Stat.DisLike) * 100)
	// 	if percent != 0 {
	// 		i.Desc = strconv.Itoa(percent) + "%的人推荐"
	// 	}
	// }
}

func (i *Item) FromRcmdReason(c *card.PopularCard) {
	var content string
	switch c.ReasonType {
	case 0:
		content = ""
	case 1:
		content = "编辑精选"
	case 2:
		content = "热门推荐"
	case 3:
		content = c.Reason
	}
	if content != "" {
		i.RcmdReason = &RcmdReason{ID: 1, Content: content, BgColor: "yellow", IconLocation: "left_top"}
		i.RcmdContent = content
	}
}

func (i *Item) FromRank(aids []int64, score map[int64]int64, am map[int64]*api.Arc) {
	var _rankCount = 3
	if len(aids) < _rankCount {
		return
	}
	items := []*Item{}
	for _, aid := range aids {
		if a, ok := am[aid]; ok {
			it := &Item{
				Title: a.Title,
				Cover: a.Pic,
				Goto:  model.GotoAv,
				Param: strconv.FormatInt(a.Aid, 10),
			}
			it.fillArcStat(a)
			it.Duration = a.Duration
			it.URI = model.FillURI(it.Goto, it.Param, model.AvHandler(a))
			if s, ok := score[aid]; ok {
				if s < 10000 {
					it.Score = model.Rounding(s, 0)
				} else if s >= 10000 && s < 100000000 {
					it.Score = model.Rounding(s, 10000) + "万"
				} else if s >= 100000000 {
					it.Score = model.Rounding(s, 100000000) + "亿"
				}
			}
			if it.Score != "" {
				it.Score = "综合评分:" + it.Score
			} else {
				it.Score = "综合评分:-"
			}
			items = append(items, it)
			if len(items) >= _rankCount {
				break
			}
		}
	}
	i.Title = "全站排行榜"
	i.Goto = model.GotoRank
	i.Item = items
	i.Param = "0"
	i.URI = "bilibili://rank?order_type=1&tid=0"
}

func (i *Item) FromHotTopic(hotTopics []*clive.TopicHot) {
	is := []*Item{}
	for _, t := range hotTopics {
		it := &Item{}
		it.Name = t.TName
		it.Param = strconv.Itoa(t.TID)
		it.Cover = t.ImageURL
		it.URI = model.FillURIHotTopic(it.Param, it.Name)
		is = append(is, it)
	}
	i.Item = is
	i.Title = "热门话题"
	i.Param = "0"
	i.Goto = model.GotoHotTopic
	i.URI = "activity://following/hot_topic_list"
	i.Desc = "更多热门话题"
}
