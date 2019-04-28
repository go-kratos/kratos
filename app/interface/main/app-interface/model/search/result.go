package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "hash/crc32"
	"math"
	"regexp"
	"strconv"
	"strings"

	"go-common/app/interface/main/app-interface/model"
	bangumimdl "go-common/app/interface/main/app-interface/model/bangumi"
	"go-common/app/interface/main/app-interface/model/bplus"
	"go-common/app/interface/main/app-interface/model/live"
	tagmdl "go-common/app/interface/main/app-interface/model/tag"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
	xtime "go-common/library/time"
)

var (
	getHightLight = regexp.MustCompile(`<em.*?em>`)
	payBadge      = &model.ReasonStyle{Text: "付费",
		TextColor:        "#FFFFFFFF",
		TextColorNight:   "#E5E5E5",
		BgColor:          "#FAAB4B",
		BgColorNight:     "#BA833F",
		BorderColor:      "#FAAB4B",
		BorderColorNight: "#BA833F",
		BgStyle:          BgStyleFill,
	}
	cooperationBadge = &model.ReasonStyle{Text: "合作",
		TextColor:        "#FFFFFFFF",
		TextColorNight:   "#E5E5E5",
		BgColor:          "#FB7299",
		BgColorNight:     "#BB5B76",
		BorderColor:      "#FB7299",
		BorderColorNight: "#BB5B76",
		BgStyle:          BgStyleFill,
	}
	videoStrongStyle = &model.ReasonStyle{
		TextColor:        "#FFFFFFFF",
		TextColorNight:   "#E5E5E5",
		BgColor:          "#FAAB4B",
		BgColorNight:     "#BA833F",
		BorderColor:      "#FAAB4B",
		BorderColorNight: "#BA833F",
		BgStyle:          BgStyleFill,
	}
	videoWeekStyle = &model.ReasonStyle{
		TextColor:        "#FAAB4B",
		TextColorNight:   "#BA833F",
		BgColor:          "",
		BgColorNight:     "",
		BorderColor:      "#FAAB4B",
		BorderColorNight: "#BA833F",
		BgStyle:          BgStyleStroke,
	}
)

// search const
const (
	_emptyLiveCover  = "https://static.hdslb.com/images/transparent.gif"
	_emptyLiveCover2 = "https://i0.hdslb.com/bfs/live/0477300d2adf65062a3d1fb7ef92122b82213b0f.png"

	StarSpace   = 1
	StarChannel = 2

	BgStyleFill              = int8(1)
	BgStyleStroke            = int8(2)
	BgStyleFillAndStroke     = int8(3)
	BgStyleNoFillAndNoStroke = int8(4)
)

// Result struct
type Result struct {
	Trackid   string      `json:"trackid,omitempty"`
	Page      int         `json:"page,omitempty"`
	NavInfo   []*NavInfo  `json:"nav,omitempty"`
	Items     ResultItems `json:"items,omitempty"`
	Item      []*Item     `json:"item,omitempty"`
	Array     int         `json:"array,omitempty"`
	Attribute int32       `json:"attribute"`
	EasterEgg *EasterEgg  `json:"easter_egg,omitempty"`
}

// ResultItems struct
type ResultItems struct {
	SuggestKeyWord *Item   `json:"suggest_keyword,omitempty"`
	Operation      []*Item `json:"operation,omitempty"`
	Season2        []*Item `json:"season2,omitempty"`
	Season         []*Item `json:"season,omitempty"`
	Upper          []*Item `json:"upper,omitempty"`
	Movie2         []*Item `json:"movie2,omitempty"`
	Movie          []*Item `json:"movie,omitempty"`
	Archive        []*Item `json:"archive,omitempty"`
	LiveRoom       []*Item `json:"live_room,omitempty"`
	LiveUser       []*Item `json:"live_user,omitempty"`
}

// NavInfo struct
type NavInfo struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
	Pages int    `json:"pages"`
	Type  int    `json:"type"`
	Show  int    `json:"show_more,omitempty"`
}

// TypeSearch struct
type TypeSearch struct {
	TrackID string  `json:"trackid"`
	Pages   int     `json:"pages"`
	Total   int     `json:"total"`
	Items   []*Item `json:"items,omitempty"`
}

// TypeSearchLiveAll struct
type TypeSearchLiveAll struct {
	TrackID string      `json:"trackid"`
	Pages   int         `json:"pages"`
	Total   int         `json:"total"`
	Master  *TypeSearch `json:"live_master,omitempty"`
	Room    *TypeSearch `json:"live_room,omitempty"`
}

// Suggestion struct
type Suggestion struct {
	TrackID string      `json:"trackid"`
	UpUser  interface{} `json:"upuser,omitempty"`
	Bangumi interface{} `json:"bangumi,omitempty"`
	Suggest []string    `json:"suggest,omitempty"`
}

// Suggestion2 struct
type Suggestion2 struct {
	TrackID string  `json:"trackid"`
	List    []*Item `json:"list,omitempty"`
}

// SuggestionResult3 struct
type SuggestionResult3 struct {
	TrackID string  `json:"trackid"`
	List    []*Item `json:"list,omitempty"`
}

// RecommendResult struct
type RecommendResult struct {
	TrackID string  `json:"trackid"`
	Title   string  `json:"title,omitempty"`
	Pages   int     `json:"pages"`
	Items   []*Item `json:"list,omitempty"`
}

// DefaultWordResult struct
type DefaultWordResult struct {
	TrackID string  `json:"trackid"`
	Title   string  `json:"title,omitempty"`
	Pages   int     `json:"pages"`
	Items   []*Item `json:"items,omitempty"`
}

// NoResultRcndResult struct
type NoResultRcndResult struct {
	TrackID string  `json:"trackid"`
	Title   string  `json:"title,omitempty"`
	Pages   int     `json:"pages"`
	Items   []*Item `json:"items,omitempty"`
}

// EasterEgg struct
type EasterEgg struct {
	ID        int64 `json:"id,omitempty"`
	ShowCount int   `json:"show_count,omitempty"`
}

// RecommendPreResult struct
type RecommendPreResult struct {
	TrackID string  `json:"trackid"`
	Total   int     `json:"total"`
	Items   []*Item `json:"items,omitempty"`
}

// Item struct
type Item struct {
	TrackID        string `json:"trackid,omitempty"`
	LinkType       string `json:"linktype,omitempty"`
	Position       int    `json:"position,omitempty"`
	SuggestKeyword string `json:"suggest_keyword,omitempty"`
	Title          string `json:"title,omitempty"`
	Name           string `json:"name,omitempty"`
	Cover          string `json:"cover,omitempty"`
	URI            string `json:"uri,omitempty"`
	Param          string `json:"param,omitempty"`
	Goto           string `json:"goto,omitempty"`
	// av
	Play       int                  `json:"play,omitempty"`
	Danmaku    int                  `json:"danmaku,omitempty"`
	Author     string               `json:"author,omitempty"`
	ViewType   string               `json:"view_type,omitempty"`
	PTime      xtime.Time           `json:"ptime,omitempty"`
	RecTags    []string             `json:"rec_tags,omitempty"`
	IsPay      int                  `json:"is_pay,omitempty"`
	NewRecTags []*model.ReasonStyle `json:"new_rec_tags,omitempty"`
	// bangumi season
	SeasonID       int64   `json:"season_id,omitempty"`
	SeasonType     int     `json:"season_type,omitempty"`
	SeasonTypeName string  `json:"season_type_name,omitempty"`
	Finish         int8    `json:"finish,omitempty"`
	Started        int8    `json:"started,omitempty"`
	Index          string  `json:"index,omitempty"`
	NewestCat      string  `json:"newest_cat,omitempty"`
	NewestSeason   string  `json:"newest_season,omitempty"`
	CatDesc        string  `json:"cat_desc,omitempty"`
	TotalCount     int     `json:"total_count,omitempty"`
	MediaType      int     `json:"media_type,omitempty"`
	PlayState      int     `json:"play_state,omitempty"`
	Style          string  `json:"style,omitempty"`
	Styles         string  `json:"styles,omitempty"`
	CV             string  `json:"cv,omitempty"`
	Rating         float64 `json:"rating,omitempty"`
	Vote           int     `json:"vote,omitempty"`
	RatingCount    int     `json:"rating_count,omitempty"`
	// BadgeType    int     `json:"badge_type,omitempty"`
	OutName string `json:"out_name,omitempty"`
	OutIcon string `json:"out_icon,omitempty"`
	OutURL  string `json:"out_url,omitempty"`
	// upper
	Sign           string          `json:"sign,omitempty"`
	Fans           int             `json:"fans,omitempty"`
	Level          int             `json:"level,omitempty"`
	Desc           string          `json:"desc,omitempty"`
	OfficialVerify *OfficialVerify `json:"official_verify,omitempty"`
	AvItems        []*Item         `json:"av_items,omitempty"`
	Item           []*Item         `json:"item,omitempty"`
	CTime          int64           `json:"ctime,omitempty"`
	IsUp           bool            `json:"is_up,omitempty"`
	LiveURI        string          `json:"live_uri,omitempty"`
	// movie
	ScreenDate string `json:"screen_date,omitempty"`
	Area       string `json:"area,omitempty"`
	CoverMark  string `json:"cover_mark,omitempty"`
	// arc and sp
	Arcs int `json:"archives,omitempty"`
	// arc and movie
	Duration    string `json:"duration,omitempty"`
	DurationInt int64  `json:"duration_int,omitempty"`
	Actors      string `json:"actors,omitempty"`
	Staff       string `json:"staff,omitempty"`
	Length      int    `json:"length,omitempty"`
	Status      int    `json:"status,omitempty"`
	// live
	RoomID      int64  `json:"roomid,omitempty"`
	Mid         int64  `json:"mid,omitempty"`
	Type        string `json:"type,omitempty"`
	Attentions  int    `json:"attentions,omitempty"`
	LiveStatus  int    `json:"live_status,omitempty"`
	Tags        string `json:"tags,omitempty"`
	Region      int    `json:"region,omitempty"`
	Online      int    `json:"online,omitempty"`
	ShortID     int    `json:"short_id,omitempty"`
	CateName    string `json:"area_v2_name,omitempty"`
	IsSelection int    `json:"is_selection,omitempty"`
	// article
	ID         int64    `json:"id,omitempty"`
	TemplateID int      `json:"template_id,omitempty"`
	ImageUrls  []string `json:"image_urls,omitempty"`
	View       int      `json:"view,omitempty"`
	Like       int      `json:"like,omitempty"`
	Reply      int      `json:"reply,omitempty"`
	// special
	Badge      string      `json:"badge,omitempty"`
	RcmdReason *RcmdReason `json:"rcmd_reason,omitempty"`
	// media bangumi and mdeia ft
	Prompt   string  `json:"prompt,omitempty"`
	Episodes []*Item `json:"episodes,omitempty"`
	Label    string  `json:"label,omitempty"`
	// game
	Reserve string `json:"reserve,omitempty"`
	// user
	Face string `json:"face,omitempty"`
	// suggest
	From      string  `json:"from,omitempty"`
	KeyWord   string  `json:"keyword,omitempty"`
	CoverSize float64 `json:"cover_size,omitempty"`
	SugType   string  `json:"sug_type,omitempty"`
	TermType  int     `json:"term_type,omitempty"`
	// rcmd query
	List       []*Item `json:"list,omitempty"`
	FromSource string  `json:"from_source,omitempty"`
	// live master
	UCover         string `json:"ucover,omitempty"`
	VerifyType     int    `json:"verify_type,omitempty"`
	VerifyDesc     string `json:"verify_desc,omitempty"`
	LevelColor     int64  `json:"level_color,omitempty"`
	IsAttention    int    `json:"is_atten,omitempty"`
	CateParentName string `json:"cate_parent_name,omitempty"`
	CateNameNew    string `json:"cate_name,omitempty"`
	Glory          *Glory `json:"glory_info,omitempty"`
	// twitter
	Covers     []string `json:"covers,omitempty"`
	CoverCount int      `json:"cover_count,omitempty"`
	Upper      *Item    `json:"upper,omitempty"`
	State      *Item    `json:"stat,omitempty"`
	// star
	TagItems []*Item `json:"tag_items,omitempty"`
	TagID    int64   `json:"tag_id,omitempty"`
	URIType  int     `json:"uri_type,omitempty"`
	// ticket
	ShowTime      string `json:"show_time,omitempty"`
	City          string `json:"city,omitempty"`
	Venue         string `json:"venue,omitempty"`
	Price         int    `json:"price,omitempty"`
	PriceComplete string `json:"price_complete,omitempty"`
	PriceType     int    `json:"price_type,omitempty"`
	ReqNum        int    `json:"required_number,omitempty"`
	// product
	ShopName string `json:"shop_name,omitempty"`
	// specialer_guide
	Phone    string               `json:"phone,omitempty"`
	Badges   []*model.ReasonStyle `json:"badges,omitempty"`
	ComicURL string               `json:"comic_url,omitempty"`
	// suggest_keyword
	SugKeyWordType int `json:"sugKeyWord_type,omitempty"`
}

// Glory live struct
type Glory struct {
	Title string  `json:"title,omitempty"`
	Total int     `json:"total"`
	Items []*Item `json:"items,omitempty"`
}

// RcmdReason struct
type RcmdReason struct {
	Content string `json:"content,omitempty"`
}

// UserResult struct
type UserResult struct {
	Items []*Item `json:"items,omitempty"`
}

// DefaultWords struct
type DefaultWords struct {
	Trackid   string `json:"trackid,omitempty"`
	Param     string `json:"param,omitempty"`
	Show      string `json:"show,omitempty"`
	Word      string `json:"word,omitempty"`
	ShowFront int    `json:"show_front,omitempty"`
}

// FromSeason form func
func (i *Item) FromSeason(b *Bangumi, bangumi string) {
	i.Title = b.Title
	i.Cover = b.Cover
	i.Goto = model.GotoBangumi
	i.Param = strconv.Itoa(int(b.SeasonID))
	i.URI = model.FillURI(bangumi, i.Param, nil)
	i.Finish = int8(b.IsFinish)
	i.Started = int8(b.IsStarted)
	i.Index = b.NewestEpIndex
	i.NewestCat = b.NewestCat
	i.NewestSeason = b.NewestSeason
	i.TotalCount = b.TotalCount
	var buf bytes.Buffer
	if b.CatList.TV != 0 {
		buf.WriteString(`TV(`)
		buf.WriteString(strconv.Itoa(b.CatList.TV))
		buf.WriteString(`) `)
	}
	if b.CatList.Movie != 0 {
		buf.WriteString(`剧场版(`)
		buf.WriteString(strconv.Itoa(b.CatList.Movie))
		buf.WriteString(`) `)
	}
	if b.CatList.Ova != 0 {
		buf.WriteString(`OVA/OAD/SP(`)
		buf.WriteString(strconv.Itoa(b.CatList.Ova))
		buf.WriteString(`)`)
	}
	i.CatDesc = buf.String()
}

// FromUpUser form func
func (i *Item) FromUpUser(u *User, as map[int64]*api.Arc, lv *live.RoomInfo) {
	i.Title = u.Name
	i.Cover = u.Pic
	i.Goto = model.GotoAuthor
	i.OfficialVerify = u.OfficialVerify
	i.Param = strconv.Itoa(int(u.Mid))
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Mid = u.Mid
	i.Sign = u.Usign
	i.Fans = u.Fans
	i.Level = u.Level
	i.Arcs = u.Videos
	i.AvItems = make([]*Item, 0, len(u.Res))
	for _, v := range u.Res {
		vi := &Item{}
		vi.Title = v.Title
		vi.Cover = v.Pic
		vi.Goto = model.GotoAv
		vi.Param = strconv.Itoa(int(v.Aid))
		vi.URI = model.FillURI(vi.Goto, vi.Param, model.AvHandler(archive.BuildArchive3(as[v.Aid])))
		a, ok := as[v.Aid]
		if ok {
			vi.Play = int(a.Stat.View)
			vi.Danmaku = int(a.Stat.Danmaku)
			if a.Rights.UGCPay == 1 {
				vi.Badges = append(vi.Badges, payBadge)
			}
			if a.Rights.IsCooperation == 1 {
				vi.Badges = append(vi.Badges, cooperationBadge)
			}
		} else {
			switch play := v.Play.(type) {
			case float64:
				vi.Play = int(play)
			case string:
				vi.Play, _ = strconv.Atoi(play)
			}
			vi.Danmaku = v.Danmaku
		}
		vi.IsPay = v.IsPay
		vi.CTime = v.Pubdate
		vi.Duration = v.Duration
		i.AvItems = append(i.AvItems, vi)
	}
	i.LiveStatus = u.IsLive
	i.RoomID = u.RoomID
	i.IsUp = u.IsUpuser == 1
	if i.RoomID != 0 {
		i.LiveURI = model.FillURI(model.GotoLive, strconv.Itoa(int(u.RoomID)), model.LiveHandler(lv))
	}
}

// FromUser form func
func (i *Item) FromUser(u *User, as map[int64]*api.Arc, lv *live.RoomInfo) {
	i.Title = u.Name
	i.Cover = u.Pic
	i.Goto = model.GotoAuthor
	i.OfficialVerify = u.OfficialVerify
	i.Param = strconv.Itoa(int(u.Mid))
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Mid = u.Mid
	i.Sign = u.Usign
	i.Fans = u.Fans
	i.Level = u.Level
	i.Arcs = u.Videos
	i.AvItems = make([]*Item, 0, len(u.Res))
	i.LiveStatus = u.IsLive
	i.RoomID = u.RoomID
	if i.RoomID != 0 {
		i.LiveURI = model.FillURI(model.GotoLive, strconv.Itoa(int(u.RoomID)), model.LiveHandler(lv))
	}
	if u.IsUpuser == 1 {
		for _, v := range u.Res {
			vi := &Item{}
			vi.Title = v.Title
			vi.Cover = v.Pic
			vi.Goto = model.GotoAv
			vi.Param = strconv.Itoa(int(v.Aid))
			vi.URI = model.FillURI(vi.Goto, vi.Param, model.AvHandler(archive.BuildArchive3(as[v.Aid])))
			a, ok := as[v.Aid]
			if ok {
				vi.Play = int(a.Stat.View)
				vi.Danmaku = int(a.Stat.Danmaku)
				if a.Rights.UGCPay == 1 {
					vi.Badges = append(vi.Badges, payBadge)
				}
				if a.Rights.IsCooperation == 1 {
					vi.Badges = append(vi.Badges, cooperationBadge)
				}
			} else {
				switch play := v.Play.(type) {
				case float64:
					vi.Play = int(play)
				case string:
					vi.Play, _ = strconv.Atoi(play)
				}
				vi.Danmaku = v.Danmaku
			}
			vi.IsPay = v.IsPay
			vi.CTime = v.Pubdate
			vi.Duration = v.Duration
			i.AvItems = append(i.AvItems, vi)
		}
		i.IsUp = true
	}
}

// FromMovie form func
func (i *Item) FromMovie(m *Movie, as map[int64]*api.Arc) {
	i.Title = m.Title
	i.Desc = m.Desc
	if m.Type == "movie" {
		i.Cover = m.Cover
		i.Param = strconv.Itoa(int(m.Aid))
		i.Goto = model.GotoAv
		i.URI = model.FillURI(i.Goto, i.Param, model.AvHandler(archive.BuildArchive3(as[m.Aid])))
		i.CoverMark = model.StatusMark(m.Status)
	} else if m.Type == "special" {
		i.Param = m.SpID
		i.Goto = model.GotoSp
		i.URI = model.FillURI(i.Goto, i.Param, nil)
		i.Cover = m.Pic
	}
	i.Staff = m.Staff
	i.Actors = m.Actors
	i.Area = m.Area
	i.Length = m.Length
	i.Status = m.Status
	i.ScreenDate = m.ScreenDate
}

// FromVideo form func
func (i *Item) FromVideo(v *Video, a *api.Arc, cooperation bool) {
	i.Title = v.Title
	i.Cover = v.Pic
	i.Author = v.Author
	i.Param = strconv.Itoa(int(v.ID))
	i.Goto = model.GotoAv
	if a != nil {
		i.Face = a.Author.Face
		i.URI = model.FillURI(i.Goto, i.Param, model.AvHandler(archive.BuildArchive3(a)))
		i.Play = int(a.Stat.View)
		i.Danmaku = int(a.Stat.Danmaku)
		i.Mid = a.Author.Mid
		if a.Rights.UGCPay == 1 {
			i.Badges = append(i.Badges, payBadge)
		}
		if a.Rights.IsCooperation == 1 {
			i.Badges = append(i.Badges, cooperationBadge)
			if i.Author != "" && cooperation {
				i.Author += " 等联合创作"
			}
		}
	} else {
		i.URI = model.FillURI(i.Goto, i.Param, nil)
		switch play := v.Play.(type) {
		case float64:
			i.Play = int(play)
		case string:
			i.Play, _ = strconv.Atoi(play)
		}
		i.Danmaku = v.Danmaku
	}
	i.IsPay = v.IsPay
	i.Desc = v.Desc
	i.Duration = v.Duration
	i.ViewType = v.ViewType
	i.RecTags = v.RecTags
	for _, r := range v.NewRecTags {
		if r.Name != "" {
			switch r.Style {
			case BgStyleFill:
				videoStrongStyle.Text = r.Name
				i.NewRecTags = append(i.NewRecTags, videoStrongStyle)
			case BgStyleStroke:
				videoWeekStyle.Text = r.Name
				i.NewRecTags = append(i.NewRecTags, videoWeekStyle)
			}
		}
	}
}

// FromLive form func
func (i *Item) FromLive(l *Live, lv *live.RoomInfo) {
	i.RoomID = l.RoomID
	i.Mid = l.UID
	i.Title = l.Title
	i.Type = l.Type
	if l.Cover == "" {
		i.Cover = l.Uface
	} else {
		i.Cover = l.Cover
	}
	i.Name = l.Uname
	i.Online = l.Online
	i.Attentions = l.Attentions
	i.Goto = model.GotoLive
	if i.Type == "live_user" {
		i.Param = strconv.Itoa(int(i.Mid))
	} else {
		i.Param = strconv.Itoa(int(i.RoomID))
	}
	i.URI = model.FillURI(i.Goto, i.Param, model.LiveHandler(lv))
	i.Tags = l.Tags
	i.Region = l.Area
	i.Badge = "直播"
}

// FromLive2 form func
func (i *Item) FromLive2(l *Live, lv *live.RoomInfo) {
	i.RoomID = l.RoomID
	i.Mid = l.UID
	i.Title = l.Title
	i.Type = l.Type
	if l.UserCover != "" && l.UserCover != _emptyLiveCover {
		i.Cover = l.UserCover
	} else if l.Cover != "" && l.Cover != _emptyLiveCover {
		i.Cover = l.Cover
	} else {
		i.Cover = _emptyLiveCover2
	}
	i.Name = l.Uname
	i.Online = l.Online
	i.Attentions = l.Attentions
	i.Goto = model.GotoLive
	if i.Type == "live_user" {
		i.Param = strconv.Itoa(int(i.Mid))
	} else {
		i.Param = strconv.Itoa(int(i.RoomID))
	}
	i.URI = model.FillURI(i.Goto, i.Param, model.LiveHandler(lv))
	i.Tags = l.Tags
	i.Region = l.Area
	i.Badge = "直播"
	i.ShortID = l.ShortID
	i.CateName = l.CateName
	i.LiveStatus = l.LiveStatus
}

// FromArticle form func
func (i *Item) FromArticle(a *Article, acc *account.Info) {
	i.ID = a.ID
	i.Mid = a.Mid
	if acc != nil {
		i.Author = acc.Name
	}
	i.TemplateID = a.TemplateID
	i.Title = a.Title
	i.Desc = a.Desc
	i.ImageUrls = a.ImageUrls
	i.View = a.View
	i.Play = a.View
	i.Like = a.Like
	i.Reply = a.Reply
	i.Badge = "专栏"
	i.Goto = model.GotoArticle
	i.Param = strconv.Itoa(int(a.ID))
	i.URI = model.FillURI(i.Goto, i.Param, nil)
}

// FromOperate form func
func (i *Item) FromOperate(o *Operate, gt string) {
	i.Title = o.Title
	i.Cover = o.Cover
	i.URI = o.RedirectURL
	i.Param = strconv.FormatInt(o.ID, 10)
	i.Desc = o.Desc
	i.Badge = o.Corner
	i.Goto = gt
	if o.RecReason != "" {
		i.RcmdReason = &RcmdReason{Content: o.RecReason}
	}
}

// FromConverge form func
func (i *Item) FromConverge(o *Operate, am map[int64]*api.Arc, rm map[int64]*live.Room, artm map[int64]*article.Meta) {
	const _convergeMinCount = 2
	cis := make([]*Item, 0, len(o.ContentList))
	for _, c := range o.ContentList {
		ci := &Item{}
		switch c.Type {
		case 0:
			if a, ok := am[c.ID]; ok && a.IsNormal() {
				ci.Title = a.Title
				ci.Cover = a.Pic
				ci.Goto = model.GotoAv
				ci.Param = strconv.FormatInt(a.Aid, 10)
				ci.URI = model.FillURI(ci.Goto, ci.Param, model.AvHandler(archive.BuildArchive3(a)))
				ci.fillArcStat(a)
				cis = append(cis, ci)
			}
		case 1:
			if r, ok := rm[c.ID]; ok {
				if r.LiveStatus == 0 {
					continue
				}
				ci.Title = r.Title
				ci.Cover = r.Cover
				ci.Goto = model.GotoLive
				ci.Param = strconv.FormatInt(r.RoomID, 10)
				ci.Online = r.Online
				ci.URI = model.FillURI(ci.Goto, ci.Param, nil) + "?broadcast_type=" + strconv.Itoa(r.BroadcastType)
				ci.Badge = "直播"
				cis = append(cis, ci)
			}
		case 2:
			if art, ok := artm[c.ID]; ok {
				ci.Title = art.Title
				ci.Desc = art.Summary
				if len(art.ImageURLs) != 0 {
					ci.Cover = art.ImageURLs[0]
				}
				ci.Goto = model.GotoArticle
				ci.Param = strconv.FormatInt(art.ID, 10)
				ci.URI = model.FillURI(ci.Goto, ci.Param, nil)
				if art.Stats != nil {
					ci.fillArtStat(art)
				}
				ci.Badge = "文章"
				cis = append(cis, ci)
			}
		}
	}
	if len(cis) < _convergeMinCount {
		return
	}
	i.Item = cis
	i.Title = o.Title
	i.Cover = o.Cover
	i.URI = o.RedirectURL
	i.Param = strconv.FormatInt(o.ID, 10)
	i.Desc = o.Desc
	i.Badge = o.Corner
	i.Goto = model.GotoConverge
	if o.RecReason != "" {
		i.RcmdReason = &RcmdReason{Content: o.RecReason}
	}
}

// FromMedia form func
func (i *Item) FromMedia(m *Media, prompt string, gt string, bangumis map[string]*bangumimdl.Card) {
	i.Title = m.Title
	if i.Title == "" {
		i.Title = m.OrgTitle
	}
	i.Cover = m.Cover
	i.Goto = gt
	i.Param = strconv.Itoa(int(m.MediaID))
	i.URI = m.GotoURL
	i.MediaType = m.MediaType
	i.PlayState = m.PlayState
	i.Style = m.Styles
	i.CV = m.CV
	i.Staff = m.Staff
	if m.MediaScore != nil {
		i.Rating = m.MediaScore.Score
		i.Vote = m.MediaScore.UserCount
	}
	i.PTime = m.Pubtime
	areas := strings.Split(m.Areas, "、")
	if len(areas) != 0 {
		i.Area = areas[0]
	}
	i.Prompt = prompt
	i.OutName = m.AllNetName
	i.OutIcon = m.AllNetIcon
	i.OutURL = m.AllNetURL
	var hit string
	for _, v := range m.HitColumns {
		if v == "cv" {
			hit = v
			break
		} else if v == "staff" {
			hit = v
		}
	}
	if hit == "cv" {
		for _, v := range getHightLight.FindAllStringSubmatch(m.CV, -1) {
			if m.MediaType == 7 {
				i.Label = fmt.Sprintf("嘉宾: %v...", v[0])
				break
			}
			if gt == model.GotoBangumi {
				i.Label = fmt.Sprintf("声优: %v...", v[0])
				break
			} else if gt == model.GotoMovie {
				i.Label = fmt.Sprintf("演员: %v...", v[0])
				break
			}
		}
	} else if hit == "staff" {
		for _, v := range getHightLight.FindAllStringSubmatch(m.Staff, -1) {
			i.Label = fmt.Sprintf("制作人员: %v...", v[0])
			break
		}
	} else if hit == "" {
		i.Label = FormPGCLabel(m.MediaType, m.Styles, m.Staff, m.CV)
	}
	// get from PGC API.
	i.SeasonID = m.SeasonID
	ssID := strconv.Itoa(int(m.SeasonID))
	if bgm, ok := bangumis[ssID]; ok {
		i.Badge = model.FormMediaType(bgm.SeasonType)
		i.SeasonTypeName = bgm.SeasonTypeName
		i.IsAttention = bgm.IsFollow
		i.IsSelection = bgm.IsSelection
		i.SeasonType = bgm.SeasonType
		for _, v := range bgm.Episodes {
			tmp := &Item{
				Param:  strconv.Itoa(int(v.ID)),
				Index:  v.Index,
				Badges: v.Badges,
			}
			tmp.URI = model.FillURI(model.GotoEP, tmp.Param, nil)
			i.Episodes = append(i.Episodes, tmp)
		}
	}
	var (
		badges []*model.ReasonStyle
		err    error
	)
	err = json.Unmarshal(m.DisplayInfo, &badges)
	if err != nil {
		log.Error("%v", err)
		return
	}
	i.Badges = badges
}

// FromGame form func
func (i *Item) FromGame(g *Game) {
	i.Title = g.Title
	i.Cover = g.Cover
	i.Desc = g.Desc
	i.Rating = g.View
	var reserve string
	if g.Status == 1 || g.Status == 2 {
		if g.Like < 10000 {
			reserve = strconv.FormatInt(g.Like, 10) + "人预约"
		} else {
			reserve = strconv.FormatFloat(float64(g.Like)/10000, 'f', 1, 64) + "万人预约"
		}
	}
	i.Reserve = reserve
	i.Goto = model.GotoGame
	i.Param = strconv.FormatInt(g.ID, 10)
	i.URI = g.RedirectURL
}

// fillArcStat fill func
func (i *Item) fillArcStat(a *api.Arc) {
	if a.Access == 0 {
		i.Play = int(a.Stat.View)
	}
	i.Danmaku = int(a.Stat.Danmaku)
	i.Reply = int(a.Stat.Reply)
	i.Like = int(a.Stat.Like)
}

// fillArtStat fill func
func (i *Item) fillArtStat(m *article.Meta) {
	i.Play = int(m.Stats.View)
	i.Reply = int(m.Stats.Reply)
}

// FromSuggest form func
func (i *Item) FromSuggest(st *SuggestTag) {
	i.From = "search"
	if st.SpID == SuggestionJump {
		switch st.Type {
		case SuggestionAV:
			i.Title = st.Value
			i.Goto = model.GotoAv
			i.URI = model.FillURI(i.Goto, strconv.Itoa(int(st.Ref)), nil)
		case SuggestionLive:
			i.Title = st.Value
			i.Goto = model.GotoLive
			i.URI = model.FillURI(i.Goto, strconv.Itoa(int(st.Ref)), nil)
		}
	} else {
		i.Title = st.Value
	}
}

// FromSuggest2 form func
func (i *Item) FromSuggest2(st *SuggestTag, as map[int64]*api.Arc, ls map[int64]*live.RoomInfo) {
	i.From = "search"
	if st.SpID == SuggestionJump {
		switch st.Type {
		case SuggestionAV:
			i.Title = st.Value
			i.Goto = model.GotoAv
			i.URI = model.FillURI(i.Goto, strconv.Itoa(int(st.Ref)), model.AvHandler(archive.BuildArchive3(as[st.Ref])))
		case SuggestionLive:
			var (
				l  *live.RoomInfo
				ok bool
			)
			i.Title = st.Value
			i.Goto = model.GotoLive
			if l, ok = ls[st.Ref]; !ok {
				for _, v := range ls {
					if v.ShortID == st.Ref {
						l = v
						break
					}
				}
			}
			i.URI = model.FillURI(i.Goto, strconv.Itoa(int(st.Ref)), model.LiveHandler(l))
			if strings.Contains(i.URI, "broadcast_type") {
				i.URI += "&extra_jump_from=23004"
			} else {
				i.URI += "?extra_jump_from=23004"
			}
		}
	} else {
		i.Title = st.Value
	}
}

// FromSuggest3 form func
func (i *Item) FromSuggest3(st *Sug, as map[int64]*api.Arc, ls map[int64]*live.RoomInfo) {
	i.From = "search"
	i.Title = st.ShowName
	i.KeyWord = st.Term
	i.Position = st.Pos
	i.Cover = st.Cover
	i.CoverSize = st.CoverSize
	i.SugType = st.SubType
	i.TermType = st.TermType
	if st.TermType == SuggestionJump {
		switch st.SubType {
		case SuggestionAV:
			i.Goto = model.GotoAv
			i.URI = model.FillURI(i.Goto, strconv.Itoa(int(st.Ref)), model.AvHandler(archive.BuildArchive3(as[st.Ref])))
			i.SugType = "视频"
		case SuggestionLive:
			var (
				l  *live.RoomInfo
				ok bool
			)
			i.Goto = model.GotoLive
			if l, ok = ls[st.Ref]; !ok {
				for _, v := range ls {
					if v.ShortID == st.Ref {
						l = v
						break
					}
				}
			}
			i.URI = model.FillURI(i.Goto, strconv.Itoa(int(st.Ref)), model.LiveHandler(l))
			if strings.Contains(i.URI, "broadcast_type") {
				i.URI += "&extra_jump_from=23004"
			} else {
				i.URI += "?extra_jump_from=23004"
			}
			i.SugType = "直播"
		case SuggestionArticle:
			i.Goto = model.GotoArticle
			i.URI = model.FillURI(i.Goto, strconv.Itoa(int(st.Ref)), nil)
			if !strings.Contains(i.URI, "column_from") {
				i.URI += "?column_from=search"
			}
			i.SugType = "专栏"
		}
	} else if st.TermType == SuggestionJumpUser && st.User != nil {
		i.Title = st.User.Name
		i.Cover = st.User.Face
		i.Goto = model.GotoAuthor
		i.OfficialVerify = &OfficialVerify{Type: st.User.OfficialVerifyType}
		i.Param = strconv.Itoa(int(st.User.Mid))
		i.URI = model.FillURI(i.Goto, i.Param, nil)
		i.Mid = st.User.Mid
		i.Fans = st.User.Fans
		i.Level = st.User.Level
		i.Arcs = st.User.Videos
	} else if st.TermType == SuggestionJumpPGC && st.PGC != nil {
		var styles []string
		i.Title = st.PGC.Title
		i.Cover = st.PGC.Cover
		i.PTime = st.PGC.Pubtime
		i.URI = st.PGC.GotoURL
		if pt := i.PTime.Time().Format("2006"); pt != "" {
			styles = append(styles, pt)
		}
		i.SeasonTypeName = model.FormMediaType(st.PGC.MediaType)
		if i.SeasonTypeName != "" {
			styles = append(styles, i.SeasonTypeName)
		}
		i.Goto = model.GotoPGC
		i.Param = strconv.Itoa(int(st.PGC.MediaID))
		i.Area = st.PGC.Areas
		if i.Area != "" {
			styles = append(styles, i.Area)
		}
		i.Style = st.PGC.Styles
		if len(styles) > 0 {
			i.Styles = strings.Join(styles, "|")
		}
		i.Label = FormPGCLabel(st.PGC.MediaType, st.PGC.Styles, st.PGC.Staff, st.PGC.CV)
		i.Rating = st.PGC.MediaScore
		i.Vote = st.PGC.MediaUserCount
		i.Badges = st.PGC.Badges
	}
}

// FromQuery form func
func (i *Item) FromQuery(qs []*Query) {
	i.Goto = model.GOtoRecommendWord
	for _, q := range qs {
		i.List = append(i.List, &Item{Param: strconv.FormatInt(q.ID, 10), Title: q.Name, Type: q.Type, FromSource: q.FromSource})
	}
}

func (i *Item) FromComic(c *Comic) {
	i.ID = c.ID
	i.Title = c.Title
	if len(c.Author) > 0 {
		i.Name = fmt.Sprintf("作者: %v", strings.Join(c.Author, "、"))
	}
	i.Style = c.Styles
	i.Cover = c.Cover
	i.URI = c.URL
	i.ComicURL = c.ComicURL
	i.Param = strconv.FormatInt(c.ID, 10)
	i.Goto = model.GotoComic
	i.Badge = "漫画"
}

// FromLiveMaster form func
func (i *Item) FromLiveMaster(l *Live, lv *live.RoomInfo) {
	i.Type = l.Type
	i.Name = l.Uname
	i.UCover = l.Uface
	i.Attentions = l.Fans
	i.VerifyType = l.VerifyType
	i.VerifyDesc = l.VerifyDesc
	i.Title = l.Title
	if l.Cover != "" && l.Cover != _emptyLiveCover {
		i.Cover = l.Cover
	} else {
		i.Cover = _emptyLiveCover2
	}
	i.Goto = model.GotoLive
	i.Mid = l.UID
	i.RoomID = l.RoomID
	i.Param = strconv.Itoa(int(i.RoomID))
	i.URI = model.FillURI(i.Goto, i.Param, model.LiveHandler(lv))
	i.Online = l.Online
	i.LiveStatus = l.LiveStatus
	i.CateParentName = l.CateParentName
	i.CateNameNew = l.CateName
}

// FromChannel form func
func (i *Item) FromChannel(c *Channel, avm map[int64]*api.Arc, bangumis map[int32]*seasongrpc.CardInfoProto, lm map[int64]*live.RoomInfo, tagMyInfos []*tagmdl.Tag) {
	i.ID = c.TagID
	i.Title = c.TagName
	i.Cover = c.Cover
	i.Param = strconv.FormatInt(c.TagID, 10)
	i.Goto = model.GotoChannel
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Type = c.Type
	i.Attentions = c.AttenCount
	i.Desc = c.Desc
	for _, myInfo := range tagMyInfos {
		if myInfo != nil && myInfo.TagID == c.TagID {
			i.IsAttention = myInfo.IsAtten
			break
		}
	}
	var (
		item        []*Item
		cooperation bool
	)
	for _, v := range c.Values {
		ii := &Item{TrackID: v.TrackID, LinkType: v.LinkType, Position: v.Position}
		switch v.Type {
		case TypeVideo:
			ii.FromVideo(v.Video, avm[v.Video.ID], cooperation)
			//case TypeLive:
			//	ii.FromLive(v.Live, lm[v.Live.RoomID])
			//case TypeMediaBangumi, TypeMediaFt:
			//	if bangumi, ok := bangumis[int32(v.Media.SeasonID)]; ok {
			//		ii.FromTagPGC(v.Media, bangumi)
			//	}
			//case TypeTicket:
			//	ii.FromTicket(v.Ticket)
			//case TypeProduct:
			//	ii.FromProduct(v.Product)
			//case TypeArticle:
			//	ii.FromArticle(v.Article, nil)
		}
		if ii.Goto != "" {
			item = append(item, ii)
		}
	}
	i.Item = item
}

// FromTwitter form twitter
func (i *Item) FromTwitter(t *Twitter, details map[int64]*bplus.Detail, isUP, isCount, isNew bool) {
	var (
		gt, id string
	)
	i.Title = t.Content
	i.Covers = t.Cover
	i.CoverCount = t.CoverCount
	i.Param = strconv.FormatInt(t.ID, 10)
	i.Goto = model.GotoTwitter
	if isNew {
		gt = model.GotoDynamic
		id = i.Param
	} else {
		gt = model.GotoTwitter
		id = strconv.FormatInt(t.PicID, 10)
	}
	i.URI = model.FillURI(gt, id, nil)
	if detail, ok := details[t.ID]; ok {
		if isUP {
			ii := &Item{
				Mid:   detail.Mid,
				Title: detail.NickName,
				Cover: detail.FaceImg,
			}
			i.Upper = ii
		}
		if isCount {
			ii := &Item{
				Play:  detail.ViewCount,
				Like:  detail.LikeCount,
				Reply: detail.CommentCount,
			}
			i.State = ii
		}
	}
}

// FromRcmdPre from rcmd pre.
func (i *Item) FromRcmdPre(id int64, a *api.Arc, bangumi *seasongrpc.CardInfoProto) {
	if a != nil {
		i.Title = a.Title
		i.Cover = a.Pic
		i.Author = a.Author.Name
		i.Param = strconv.Itoa(int(id))
		i.Goto = model.GotoAv
		i.URI = model.FillURI(i.Goto, i.Param, model.AvHandler(archive.BuildArchive3(a)))
		i.fillArcStat(a)
		i.Desc = a.Desc
		i.DurationInt = a.Duration
	} else if bangumi != nil {
		i.Title = bangumi.Title
		i.Cover = bangumi.Cover
		i.Param = strconv.Itoa(int(id))
		i.Goto = model.GotoPGC
		i.URI = model.FillURI(i.Goto, i.Param, nil)
		i.Badge = bangumi.SeasonTypeName
		i.Started = int8(bangumi.IsStarted)
		i.Play = int(bangumi.Stat.View)
		if bangumi.Rating != nil {
			i.Rating = float64(bangumi.Rating.Score)
			i.RatingCount = int(bangumi.Rating.Count)
		}
		i.MediaType = int(bangumi.SeasonType) // 1：番剧，2：电影，3：纪录片，4：国漫，5：电视剧
		if bangumi.Stat != nil {
			i.Attentions = int(bangumi.Stat.Follow)
		}
		if bangumi.NewEp != nil {
			i.Label = bangumi.NewEp.IndexShow
		}
	}
}

// FromStar form func
func (i *Item) FromStar(s *Star) {
	var cooperation bool
	i.Title = s.Title
	i.Cover = s.Cover
	i.Desc = s.Desc
	if i.URIType == StarSpace {
		i.URI = model.FillURI(model.GotoSpace, strconv.Itoa(int(s.MID)), nil)
	} else if i.URIType == StarChannel {
		i.URI = model.FillURI(model.GotoChannel, strconv.Itoa(int(s.TagID)), nil)
	}
	i.Param = strconv.Itoa(int(s.ID))
	i.Goto = model.GotoStar
	i.Mid = s.MID
	i.TagID = s.TagID
	i.TagItems = make([]*Item, 0, len(s.TagList))
	for _, v := range s.TagList {
		if v == nil {
			continue
		}
		vi := &Item{}
		vi.Title = v.TagName
		vi.KeyWord = v.KeyWord
		vi.Item = make([]*Item, 0, len(v.ValueList))
		for _, vv := range v.ValueList {
			if vv == nil || vv.Video == nil {
				continue
			}
			vvi := &Item{}
			vvi.FromVideo(vv.Video, nil, cooperation)
			vi.Item = append(vi.Item, vvi)
		}
		i.TagItems = append(i.TagItems, vi)
	}
}

// FromTicket from ticket
func (i *Item) FromTicket(t *Ticket) {
	i.ID = t.ID
	i.Param = strconv.Itoa(int(t.ID))
	i.Goto = model.GotoTicket
	i.Badge = "展演"
	i.Title = t.Title
	i.Cover = t.Cover
	i.ShowTime = t.ShowTime
	i.City = t.CityName
	i.Venue = t.VenueName
	i.Price = int(math.Ceil(float64(t.PriceLow) / 100))
	i.PriceComplete = strconv.FormatFloat(float64(t.PriceLow)/100, 'f', -1, 64)
	i.PriceType = t.PriceType
	i.ReqNum = t.ReqNum
	i.URI = t.URL
}

// FromProduct from ticket
func (i *Item) FromProduct(p *Product) {
	i.ID = p.ID
	i.Param = strconv.Itoa(int(p.ID))
	i.Goto = model.GotoProduct
	i.Badge = "商品"
	i.Title = p.Title
	i.Cover = p.Cover
	i.ShopName = p.ShopName
	i.Price = int(math.Ceil(float64(p.Price) / 100))
	i.PriceComplete = strconv.FormatFloat(float64(p.Price)/100, 'f', -1, 64)
	i.PriceType = p.PriceType
	i.ReqNum = p.ReqNum
	i.URI = p.URL
}

// FromSpecialerGuide from ticket
func (i *Item) FromSpecialerGuide(sg *SpecialerGuide) {
	i.ID = sg.ID
	i.Param = strconv.Itoa(int(sg.ID))
	i.Goto = model.GotoSpecialerGuide
	i.Title = sg.Title
	i.Cover = sg.Cover
	i.Desc = sg.Desc
	i.Phone = sg.Tel
}

func (i *Item) FromTagPGC(m *Media, bangumi *seasongrpc.CardInfoProto) {
	if m.SeasonID == 0 {
		return
	}
	ssid := strconv.Itoa(int(m.SeasonID))
	i.Title = bangumi.Title
	i.Cover = bangumi.Cover
	i.Param = strconv.Itoa(int(m.MediaID))
	i.Goto = model.GotoPGC
	i.URI = model.FillURI(i.Goto, ssid, nil)
	i.Badge = bangumi.SeasonTypeName
	i.Started = int8(bangumi.IsStarted)
	i.Play = int(bangumi.Stat.View)
	if bangumi.Rating != nil {
		i.Rating = float64(bangumi.Rating.Score)
		i.RatingCount = int(bangumi.Rating.Count)
	}
	i.MediaType = int(bangumi.SeasonType) // 1：番剧，2：电影，3：纪录片，4：国漫，5：电视剧
	if bangumi.Stat != nil {
		i.Attentions = int(bangumi.Stat.Follow)
	}
	if bangumi.NewEp != nil {
		i.Label = bangumi.NewEp.IndexShow
	}
}

// flowTest form func
// func flowTest(buvid string) (ok bool) {
// 	id := crc32.ChecksumIEEE([]byte(reverseString(buvid))) % 2
// 	if id%2 > 0 {
// 		ok = true
// 	}
// 	return
// }

// reverseString form func
// func reverseString(s string) string {
// 	runes := []rune(s)
// 	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
// 		runes[from], runes[to] = runes[to], runes[from]
// 	}
// 	return string(runes)
// }

func FormPGCLabel(mediaType int, styles, staff, cv string) (label string) {
	switch mediaType {
	case 1: // 演员
		label = strings.Replace(styles, "\n", "、", -1)
	case 2: // 电影
		label = "演员：" + strings.Replace(cv, "\n", "、", -1)
	case 3: // 纪录片
		label = strings.Replace(staff, "\n", "、", -1)
	case 4: // 国创
		label = strings.Replace(styles, "\n", "、", -1)
	case 5: // 电视剧
		label = "演员：" + strings.Replace(cv, "\n", "、", -1)
	case 7: // 综艺
		label = strings.Replace(cv, "\n", "、", -1)
	case 123: // 电视剧
		label = "演员：" + strings.Replace(cv, "\n", "、", -1)
	case 124: // 综艺
		label = strings.Replace(cv, "\n", "、", -1)
	case 125: // 纪录片
		label = strings.Replace(staff, "\n", "、", -1)
	case 126: // 电影
		label = "演员：" + strings.Replace(cv, "\n", "、", -1)
	case 127: // 动漫
		label = strings.Replace(styles, "\n", "、", -1)
	default:
		label = strings.Replace(cv, "\n", "、", -1)
	}
	return
}
