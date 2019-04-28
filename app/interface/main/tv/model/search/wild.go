package search

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	v1 "go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

// default .
const (
	GotoBangumi    = "bangumi"
	GotoAv         = "av"
	GotoWeb        = "web"
	GotoMovie      = "movie"
	GotoBangumiWeb = "bangumi_web"
	GotoSp         = "sp"
	GotoLive       = "live"
	GotoGame       = "game"
	GotoAuthor     = "author"
	GotoClip       = "clip"
	GotoAlbum      = "album"
	GotoArticle    = "article"
	GotoAudio      = "audio"
	GotoSpecial    = "special"
	GotoBanner     = "banner"
	GotoSpecialS   = "special_s"
	GotoConverge   = "converge"
	GotoPGC        = "pgc"
	GotoChannel    = "channel"
	GotoEP         = "ep"
	GotoTwitter    = "twitter"

	CoverIng      = "即将上映"
	CoverPay      = "付费观看"
	CoverFree     = "免费观看"
	CoverVipFree  = "付费观看"
	CoverVipOnly  = "专享"
	CoverVipFirst = "抢先"
)

// UserSearch user search request .
type UserSearch struct {
	SearchType string `form:"search_type" validate:"required"`
	Order      string `form:"order"`
	Category   int    `form:"category"`
	Platform   string `form:"platform"`
	Build      string `form:"build"`
	MobiAPP    string `form:"mobi_app"`
	Device     string `form:"device"`
	Keyword    string `form:"keyword"  validate:"required"`
	Page       int    `form:"page"  validate:"required,min=1"`
	Pagesize   int    `form:"pagesize"`
	UserType   int    `form:"user_type"`
	Highlight  int    `form:"highlight"`
	OrderSort  int    `form:"order_sort"`
	FromSource string `form:"from_source"`
	Buvid      string `form:"buvid"`
	Duration   int    `form:"duration"` // 视频时长
	// 传递参数，给到dao用
	SeasonNum int   `form:"season_num"`
	MovieNum  int   `form:"movie_num"`
	RID       int   `form:"rid"`
	MID       int64 `form:"mid"`
}

// User struct res .
type User struct {
	Type       string `json:"type"`
	Mid        int64  `json:"mid,omitempty"`
	Name       string `json:"uname,omitempty"`
	Usign      string `json:"usign,omitempty"`
	Fans       int    `json:"fans,omitempty"`
	Videos     int    `json:"videos,omitempty"`
	Pic        string `json:"upic,omitempty"`
	VerifyInfo string `json:"verify_info"`
	Level      int    `json:"level,omitempty"`
	Gender     int    `json:"gender"`
	IsUpuser   int    `json:"is_upuser,omitempty"`
	IsLive     int    `json:"is_live,omitempty"`
	RoomID     int64  `json:"room_id,omitempty"`
	Res        []*struct {
		Aid      int64       `json:"aid,omitempty"`
		Title    string      `json:"title,omitempty"`
		Pubdate  int64       `json:"pubdate,omitempty"`
		ArcURL   string      `json:"arcurl,omitempty"`
		Pic      string      `json:"pic,omitempty"`
		Play     interface{} `json:"play,omitempty"`
		Danmaku  int         `json:"dm,omitempty"`
		Coin     int         `json:"coin"`
		Fav      int         `json:"fav"`
		Desc     string      `json:"desc"`
		Duration string      `json:"duration,omitempty"`
	} `json:"res,omitempty"`
	OfficialVerify *OfficialVerify `json:"official_verify,omitempty"`
	*ResultResponse
}

// OfficialVerify struct .
type OfficialVerify struct {
	Type int    `json:"type"`
	Desc string `json:"desc,omitempty"`
}

// Search all .
type Search struct {
	Code           int    `json:"code,omitempty"`
	Trackid        string `json:"seid,omitempty"`
	Page           int    `json:"page,omitempty"`
	PageSize       int    `json:"pagesize,omitempty"`
	Total          int    `json:"total,omitempty"`
	NumResults     int    `json:"numResults,omitempty"`
	NumPages       int    `json:"numPages,omitempty"`
	SuggestKeyword string `json:"suggest_keyword,omitempty"`
	Attribute      int32  `json:"exp_bits,omitempty"`
	PageInfo       struct {
		Bangumi      *Page `json:"bangumi,omitempty"`
		UpUser       *Page `json:"upuser,omitempty"`
		BiliUser     *Page `json:"bili_user,omitempty"`
		User         *Page `json:"user,omitempty"`
		Movie        *Page `json:"movie,omitempty"`
		Film         *Page `json:"pgc,omitempty"`
		MediaBangumi *Page `json:"media_bangumi,omitempty"`
		MediaFt      *Page `json:"media_ft,omitempty"`
	} `json:"pageinfo,omitempty"`
	Result struct {
		Bangumi      []*Bangumi `json:"bangumi,omitempty"`
		UpUser       []*User    `json:"upuser,omitempty"`
		BiliUser     []*User    `json:"bili_user,omitempty"`
		User         []*User    `json:"user,omitempty"`
		Movie        []*Movie   `json:"movie,omitempty"`
		Video        []*Video   `json:"video,omitempty"`
		MediaBangumi []*Media   `json:"media_bangumi,omitempty"`
		MediaFt      []*Media   `json:"media_ft,omitempty"`
	} `json:"result,omitempty"`
}

// Media struct .
type Media struct {
	MediaID    int64  `json:"media_id,omitempty"`
	SeasonID   int64  `json:"season_id,omitempty"`
	Title      string `json:"title,omitempty"`
	OrgTitle   string `json:"org_title,omitempty"`
	Styles     string `json:"styles,omitempty"`
	Cover      string `json:"cover,omitempty"`
	PlayState  int    `json:"play_state,omitempty"`
	MediaScore *struct {
		Score     float64 `json:"score,omitempty"`
		UserCount int     `json:"user_count,omitempty"`
	} `json:"media_score,omitempty"`
	MediaType  int        `json:"media_type,omitempty"`
	CV         string     `json:"cv,omitempty"`
	Staff      string     `json:"staff,omitempty"`
	Areas      string     `json:"areas,omitempty"`
	GotoURL    string     `json:"goto_url,omitempty"`
	Pubtime    xtime.Time `json:"pubtime,omitempty"`
	HitColumns []string   `json:"hit_columns,omitempty"`
}

// Movie struct .
type Movie struct {
	Title      string `json:"title"`
	SpID       string `json:"spid"`
	Type       string `json:"type"`
	Aid        int64  `json:"aid"`
	Desc       string `json:"description"`
	Actors     string `json:"actors"`
	Staff      string `json:"staff"`
	Cover      string `json:"cover"`
	Pic        string `json:"pic"`
	ScreenDate string `json:"screenDate"`
	Area       string `json:"area"`
	Status     int    `json:"status"`
	Length     int    `json:"length"`
	Pages      int    `json:"numPages"`
}

// Video struct .
type Video struct {
	ID       int64       `json:"id"`
	Author   string      `json:"author"`
	Title    string      `json:"title"`
	Pic      string      `json:"pic"`
	Desc     string      `json:"description"`
	Play     interface{} `json:"play"`
	Danmaku  int         `json:"video_review"`
	Duration string      `json:"duration"`
	Pages    int         `json:"numPages"`
	ViewType string      `json:"view_type"`
	RecTags  []string    `json:"rec_tags"`
}

// ResultAll struct .
type ResultAll struct {
	Trackid   string      `json:"trackid,omitempty"`
	Page      int         `json:"page,omitempty"`
	NavInfo   []*NavInfo  `json:"nav,omitempty"`
	Items     ResultItems `json:"items,omitempty"`
	Item      []*Item     `json:"item,omitempty"` // 混排的数据(未用到)
	Attribute int32       `json:"attribute"`      // 实验中开关
}

// ResultItems struct .
type ResultItems struct {
	Season2 []*Item `json:"season2,omitempty"`
	Season  []*Item `json:"season,omitempty"` // 老数据字段（未用到）
	Upper   []*Item `json:"upper,omitempty"`
	Movie2  []*Item `json:"movie2,omitempty"`
	Movie   []*Item `json:"movie,omitempty"` // 老数据字段（未用到）
	Archive []*Item `json:"archive,omitempty"`
}

// NavInfo struct .
type NavInfo struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
	Pages int    `json:"pages"`
	Type  int    `json:"type"`
	Show  int    `json:"show_more,omitempty"`
}

// Item struct .
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
	Play     int        `json:"play,omitempty"`
	Danmaku  int        `json:"danmaku,omitempty"`
	Author   string     `json:"author,omitempty"`
	ViewType string     `json:"view_type,omitempty"`
	PTime    xtime.Time `json:"ptime,omitempty"`
	RecTags  []string   `json:"rec_tags,omitempty"`
	// media bangumi and mdeia ft
	Prompt   string  `json:"prompt,omitempty"`
	Episodes []*Item `json:"episodes,omitempty"`
	Label    string  `json:"label,omitempty"`
	// bangumi season
	Finish       int8    `json:"finish,omitempty"`
	Started      int8    `json:"started,omitempty"`
	Index        string  `json:"index,omitempty"`
	NewestCat    string  `json:"newest_cat,omitempty"`
	NewestSeason string  `json:"newest_season,omitempty"`
	CatDesc      string  `json:"cat_desc,omitempty"`
	TotalCount   int     `json:"total_count,omitempty"`
	MediaType    int     `json:"media_type,omitempty"`
	PlayState    int     `json:"play_state,omitempty"`
	Style        string  `json:"style,omitempty"`
	CV           string  `json:"cv,omitempty"`
	Rating       float64 `json:"rating,omitempty"`
	Vote         int     `json:"vote,omitempty"`
	RatingCount  int     `json:"rating_count,omitempty"`
	BadgeType    int     `json:"badge_type,omitempty"`
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
	// user
	Face string `json:"face,omitempty"`
	// arc and sp
	Arcs int `json:"archives,omitempty"`
	// arc and movie
	Duration    string `json:"duration,omitempty"`
	DurationInt int64  `json:"duration_int,omitempty"`
	Actors      string `json:"actors,omitempty"`
	Staff       string `json:"staff,omitempty"`
	Length      int    `json:"length,omitempty"`
	Status      int    `json:"status,omitempty"`
}

// Bangumi struct .
type Bangumi struct {
	Name          string `json:"name,omitempty"`
	SeasonID      int    `json:"season_id,omitempty"`
	Title         string `json:"title,omitempty"`
	Cover         string `json:"cover,omitempty"`
	Evaluate      string `json:"evaluate,omitempty"`
	NewestEpID    int    `json:"newest_ep_id,omitempty"`
	NewestEpIndex string `json:"newest_ep_index,omitempty"`
	IsFinish      int    `json:"is_finish,omitempty"`
	IsStarted     int    `json:"is_started,omitempty"`
	NewestCat     string `json:"newest_cat,omitempty"`
	NewestSeason  string `json:"newest_season,omitempty"`
	TotalCount    int    `json:"total_count,omitempty"`
	Pages         int    `json:"numPages,omitempty"`
	CatList       *struct {
		TV    int `json:"tv"`
		Movie int `json:"movie"`
		Ova   int `json:"ova"`
	} `json:"catlist,omitempty"`
}

// TypeSearch struct .
type TypeSearch struct {
	TrackID string  `json:"trackid"`
	Pages   int     `json:"pages"`
	Total   int     `json:"total"`
	Items   []*Item `json:"items,omitempty"`
}

// Card for bangumi .
type Card struct {
	SeasonID    int64      `json:"season_id"`
	IsFollow    int        `json:"is_follow"`
	IsSelection int        `json:"is_selection"`
	Badge       string     `json:"badge"`
	BadgeType   int        `json:"badge_type"`
	Episodes    []*Episode `json:"episodes"`
}

// Episode for bangumi card .
type Episode struct {
	ID         int64  `json:"id"`
	Badge      string `json:"badge"`
	BadgeType  int    `json:"badge_type"`
	Status     int    `json:"status"`
	Cover      string `json:"cover"`
	Index      string `json:"index"`
	IndexTitle string `json:"index_title"`
}

// StatusMark cover status mark .
func StatusMark(status int) string {
	if status == 0 {
		return CoverIng
	} else if status == 1 {
		return CoverPay
	} else if status == 2 {
		return CoverFree
	} else if status == 3 {
		return CoverVipFree
	} else if status == 4 {
		return CoverVipOnly
	} else if status == 5 {
		return CoverVipFirst
	}
	return ""
}

// FillURI deal app schema .
func FillURI(gt, param string, f func(uri string) string) (uri string) {
	switch gt {
	case GotoAv, "":
		uri = "bilibili://video/" + param
	case GotoLive:
		uri = "bilibili://live/" + param
	case GotoBangumi:
		uri = "bilibili://bangumi/season/" + param
	case GotoBangumiWeb:
		uri = "http://bangumi.bilibili.com/anime/" + param
	case GotoGame:
		uri = "bilibili://game_center/detail?id=" + param + "&sourceType=adPut"
	case GotoSp:
		uri = "bilibili://splist/" + param
	case GotoAuthor:
		uri = "bilibili://author/" + param
	case GotoClip:
		uri = "bilibili://clip/" + param
	case GotoAlbum:
		uri = "bilibili://album/" + param
	case GotoArticle:
		uri = "bilibili://article/" + param
	case GotoWeb:
		uri = param
	case GotoPGC:
		uri = "https://www.bilibili.com/bangumi/play/ss" + param
	case GotoChannel:
		uri = "bilibili://pegasus/channel/" + param + "/"
	case GotoEP:
		uri = "https://www.bilibili.com/bangumi/play/ep" + param
	case GotoTwitter:
		uri = "bilibili://pictureshow/detail/" + param
	}
	if f != nil {
		uri = f(uri)
	}
	return
}

// search const
var getHightLight = regexp.MustCompile(`<em.*?em>`)

var (
	// AvHandler .
	AvHandler = func(a *v1.Arc) func(uri string) string {
		return func(uri string) string {
			if a == nil {
				return uri
			}
			if a.Dimension.Height != 0 || a.Dimension.Width != 0 {
				return fmt.Sprintf("%s?player_width=%d&player_height=%d&player_rotate=%d", uri, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
			}
			return uri
		}
	}
)

// FromSeason .
func (i *Item) FromSeason(b *Bangumi, bangumi string) {
	i.Title = b.Title
	i.Cover = b.Cover
	i.Goto = GotoBangumi
	i.Param = strconv.Itoa(int(b.SeasonID))
	i.URI = FillURI(bangumi, i.Param, nil)
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

// FromUpUser form func .
func (i *Item) FromUpUser(u *User, as map[int64]*v1.Arc) {
	i.Title = u.Name
	i.Cover = u.Pic
	i.Goto = GotoAuthor
	i.OfficialVerify = u.OfficialVerify
	i.Param = strconv.Itoa(int(u.Mid))
	i.URI = FillURI(i.Goto, i.Param, nil)
	i.Sign = u.Usign
	i.Fans = u.Fans
	i.Level = u.Level
	i.Arcs = u.Videos
	i.AvItems = make([]*Item, 0, len(u.Res))
	for _, v := range u.Res {
		vi := &Item{}
		vi.Title = v.Title
		vi.Cover = v.Pic
		vi.Goto = GotoAv
		vi.Param = strconv.Itoa(int(v.Aid))
		a, ok := as[v.Aid]
		if ok {
			vi.Play = int(a.Stat.View)
			vi.Danmaku = int(a.Stat.Danmaku)
		} else {
			switch play := v.Play.(type) {
			case float64:
				vi.Play = int(play)
			case string:
				vi.Play, _ = strconv.Atoi(play)
			}
			vi.Danmaku = v.Danmaku
		}
		vi.CTime = v.Pubdate
		vi.Duration = v.Duration
		i.AvItems = append(i.AvItems, vi)
	}
}

// FromUser form func .
func (i *Item) FromUser(u *User, as map[int64]*v1.Arc) {
	i.Title = u.Name
	i.Cover = u.Pic
	i.Goto = GotoAuthor
	i.OfficialVerify = u.OfficialVerify
	i.Param = strconv.Itoa(int(u.Mid))
	i.URI = FillURI(i.Goto, i.Param, nil)
	i.Sign = u.Usign
	i.Fans = u.Fans
	i.Level = u.Level
	i.Arcs = u.Videos
	i.AvItems = make([]*Item, 0, len(u.Res))
	if u.IsUpuser == 1 {
		for _, v := range u.Res {
			vi := &Item{}
			vi.Title = v.Title
			vi.Cover = v.Pic
			vi.Goto = GotoAv
			vi.Param = strconv.Itoa(int(v.Aid))
			a, ok := as[v.Aid]
			if ok {
				vi.Play = int(a.Stat.View)
				vi.Danmaku = int(a.Stat.Danmaku)
			} else {
				switch play := v.Play.(type) {
				case float64:
					vi.Play = int(play)
				case string:
					vi.Play, _ = strconv.Atoi(play)
				}
				vi.Danmaku = v.Danmaku
			}
			vi.CTime = v.Pubdate
			vi.Duration = v.Duration
			i.AvItems = append(i.AvItems, vi)
		}
		i.IsUp = true
	}
}

// FromMovie form func .
func (i *Item) FromMovie(m *Movie, as map[int64]*v1.Arc) {
	i.Title = m.Title
	i.Desc = m.Desc
	if m.Type == "movie" {
		i.Cover = m.Cover
		i.Param = strconv.Itoa(int(m.Aid))
		i.Goto = GotoAv
		i.URI = FillURI(i.Goto, i.Param, AvHandler(as[m.Aid]))
		i.CoverMark = StatusMark(m.Status)
	} else if m.Type == "special" {
		i.Param = m.SpID
		i.Goto = GotoSp
		i.URI = FillURI(i.Goto, i.Param, nil)
		i.Cover = m.Pic
	}
	i.Staff = m.Staff
	i.Actors = m.Actors
	i.Area = m.Area
	i.Length = m.Length
	i.Status = m.Status
	i.ScreenDate = m.ScreenDate
}

// FromVideo form func .
func (i *Item) FromVideo(v *Video, a *v1.Arc) {
	i.Title = v.Title
	i.Cover = v.Pic
	i.Author = v.Author
	i.Param = strconv.Itoa(int(v.ID))
	i.Goto = GotoAv
	if a != nil {
		i.Face = a.Author.Face
		i.URI = FillURI(i.Goto, i.Param, AvHandler(a))
		i.Play = int(a.Stat.View)
		i.Danmaku = int(a.Stat.Danmaku)
	} else {
		i.URI = FillURI(i.Goto, i.Param, nil)
		switch play := v.Play.(type) {
		case float64:
			i.Play = int(play)
		case string:
			i.Play, _ = strconv.Atoi(play)
		}
		i.Danmaku = v.Danmaku
	}
	i.Desc = v.Desc
	i.Duration = v.Duration
	i.ViewType = v.ViewType
	i.RecTags = v.RecTags
}

// FromMedia form func .
func (i *Item) FromMedia(m *Media, prompt string, gt string, bangumis map[string]*Card) {
	i.Title = m.Title
	if i.Title == "" {
		i.Title = m.OrgTitle
	}
	i.Cover = m.Cover
	i.Goto = gt
	i.Param = strconv.Itoa(int(m.SeasonID))
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
			if gt == GotoBangumi {
				i.Label = fmt.Sprintf("声优: %v...", v[0])
				break
			} else if gt == GotoMovie {
				i.Label = fmt.Sprintf("演员: %v...", v[0])
				break
			}
		}
	} else if hit == "staff" {
		for _, v := range getHightLight.FindAllStringSubmatch(m.Staff, -1) {
			i.Label = fmt.Sprintf("制作人员: %v...", v[0])
			break
		}
	}

	// get from PGC API .
	ssID := strconv.Itoa(int(m.SeasonID))
	if bgm, ok := bangumis[ssID]; ok {
		for _, v := range bgm.Episodes {
			tmp := &Item{
				Param:     strconv.Itoa(int(v.ID)),
				Index:     v.Index,
				BadgeType: v.BadgeType,
			}
			tmp.URI = FillURI(GotoEP, tmp.Param, nil)
			i.Episodes = append(i.Episodes, tmp)
		}
	}
}
