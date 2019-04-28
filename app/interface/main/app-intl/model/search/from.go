package search

import (
	"encoding/json"

	"go-common/library/log"
	xtime "go-common/library/time"
)

// search const
const (
	TypeVideo        = "video"
	TypeLive         = "live_room"
	TypeMediaBangumi = "media_bangumi"
	TypeMediaFt      = "media_ft"
	TypeArticle      = "article"
	TypeSpecial      = "special_card"
	TypeBanner       = "banner"
	TypeUser         = "user"
	TypeBiliUser     = "bili_user"
	TypeGame         = "game"
	TypeSpecialS     = "special_card_small"
	TypeConverge     = "content_card"
	TypeQuery        = "query"
	TypeLiveMaster   = "live_master"
	TypeTwitter      = "twitter"
	// TypeLiveComic    = "comic"

	SuggestionJump = 99
	SuggestionAV   = "video"
	SuggestionLive = "live"
)

// Search all
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
		UpUser       *Page `json:"upuser,omitempty"`
		BiliUser     *Page `json:"bili_user,omitempty"`
		User         *Page `json:"user,omitempty"`
		Movie        *Page `json:"movie,omitempty"`
		Article      *Page `json:"article,omitempty"`
		LiveRoom     *Page `json:"live_room,omitempty"`
		LiveUser     *Page `json:"live_user,omitempty"`
		LiveAll      *Page `json:"live_all,omitempty"`
		MediaBangumi *Page `json:"media_bangumi,omitempty"`
		MediaFt      *Page `json:"media_ft,omitempty"`
	} `json:"pageinfo,omitempty"`
	Result struct {
		Bangumi      []*Bangumi `json:"bangumi,omitempty"`
		UpUser       []*User    `json:"upuser,omitempty"`
		BiliUser     []*User    `json:"bili_user,omitempty"`
		User         []*User    `json:"user,omitempty"`
		Movie        []*Movie   `json:"movie,omitempty"`
		LiveRoom     []*Live    `json:"live_room,omitempty"`
		LiveUser     []*Live    `json:"live_user,omitempty"`
		Video        []*Video   `json:"video,omitempty"`
		MediaBangumi []*Media   `json:"media_bangumi,omitempty"`
		MediaFt      []*Media   `json:"media_ft,omitempty"`
	} `json:"result,omitempty"`
	FlowResult      []*Flow `json:"flow_result,omitempty"`
	FlowPlaceholder int     `json:"flow_placeholder,omitempty"`
	EggInfo         *struct {
		Source    int64 `json:"source,omitempty"`
		ShowCount int   `json:"show_count,omitempty"`
	} `json:"egg_info,omitempty"`
}

// NoResultRcmd no result rcmd
type NoResultRcmd struct {
	Code           int      `json:"code,omitempty"`
	Msg            string   `json:"msg,omitempty"`
	ReqType        string   `json:"req_type,omitempty"`
	Result         []*Video `json:"result,omitempty"`
	NumResults     int      `json:"numResults,omitempty"`
	Page           int      `json:"page,omitempty"`
	Trackid        string   `json:"seid,omitempty"`
	SuggestKeyword string   `json:"suggest_keyword,omitempty"`
	RecommendTips  string   `json:"recommend_tips,omitempty"`
}

// RecommendPre search at pre-page
type RecommendPre struct {
	Code      int    `json:"code,omitempty"`
	Msg       string `json:"msg,omitempty"`
	NumResult int    `json:"numResult,omitempty"`
	Trackid   string `json:"seid,omitempty"`
	Result    []*struct {
		Type  string `json:"type,omitempty"`
		Query string `json:"query,omitempty"`
		List  []*struct {
			Type string `json:"source_type,omitempty"`
			ID   int64  `json:"source_id,omitempty"`
		} `json:"rec_list,omitempty"`
	} `json:"result,omitempty"`
}

// Page struct
type Page struct {
	NumResults int `json:"numResults"`
	Pages      int `json:"pages"`
}

// Bangumi struct
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

// Movie struct
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

// User struct
type User struct {
	Mid            int64           `json:"mid,omitempty"`
	Name           string          `json:"uname,omitempty"`
	SName          string          `json:"name,omitempty"`
	OfficialVerify *OfficialVerify `json:"official_verify,omitempty"`
	Usign          string          `json:"usign,omitempty"`
	Fans           int             `json:"fans,omitempty"`
	Videos         int             `json:"videos,omitempty"`
	Level          int             `json:"level,omitempty"`
	Pic            string          `json:"upic,omitempty"`
	Pages          int             `json:"numPages,omitempty"`
	Res            []*struct {
		Play     interface{} `json:"play,omitempty"`
		Danmaku  int         `json:"dm,omitempty"`
		Pubdate  int64       `json:"pubdate,omitempty"`
		Title    string      `json:"title,omitempty"`
		Aid      int64       `json:"aid,omitempty"`
		Pic      string      `json:"pic,omitempty"`
		ArcURL   string      `json:"arcurl,omitempty"`
		Duration string      `json:"duration,omitempty"`
	} `json:"res,omitempty"`
	IsLive   int   `json:"is_live,omitempty"`
	RoomID   int64 `json:"room_id,omitempty"`
	IsUpuser int   `json:"is_upuser,omitempty"`
}

// OfficialVerify struct
type OfficialVerify struct {
	Type int    `json:"type"`
	Desc string `json:"desc,omitempty"`
}

// Video struct
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

// Live struct
type Live struct {
	Total          int    `json:"total,omitempty"`
	Pages          int    `json:"pages"`
	UID            int64  `json:"uid,omitempty"`
	RoomID         int64  `json:"roomid,omitempty"`
	Type           string `json:"type,omitempty"`
	Title          string `json:"title,omitempty"`
	LiveStatus     int    `json:"live_status,omitempty"`
	ShortID        int    `json:"short_id,omitempty"`
	Uname          string `json:"uname,omitempty"`
	Uface          string `json:"uface,omitempty"`
	Cover          string `json:"cover,omitempty"`
	Online         int    `json:"online,omitempty"`
	Attentions     int    `json:"attentions,omitempty"`
	Tags           string `json:"tags,omitempty"`
	Area           int    `json:"area,omitempty"`
	CateName       string `json:"cate_name,omitempty"`
	CateParentName string `json:"cate_parent_name,omitempty"`
	UserCover      string `json:"user_cover,omitempty"`
	VerifyType     int    `json:"verify_type,omitempty"`
	VerifyDesc     string `json:"verify_desc,omitempty"`
	Fans           int    `json:"fans,omitempty"`
}

// Article struct
type Article struct {
	ID         int64    `json:"id"`
	Mid        int64    `json:"mid"`
	Uname      string   `json:"uname"`
	TemplateID int      `json:"template_id"`
	Title      string   `json:"title"`
	Desc       string   `json:"desc"`
	ImageUrls  []string `json:"image_urls"`
	View       int      `json:"view"`
	Like       int      `json:"like"`
	Reply      int      `json:"reply"`
}

// Media struct
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

// Query struct
type Query struct {
	Type       string `json:"type,omitempty"`
	Name       string `json:"name,omitempty"`
	ID         int64  `json:"id,omitempty"`
	FromSource string `json:"from_source,omitempty"`
}

// Hot struct
type Hot struct {
	Code    int    `json:"code,omitempty"`
	SeID    string `json:"seid,omitempty"`
	TrackID string `json:"trackid"`
	List    []struct {
		Keyword  string `json:"keyword"`
		Status   string `json:"status"`
		NameType string `json:"name_type"`
	} `json:"list"`
}

// Suggest struct
type Suggest struct {
	Code     int         `json:"code"`
	Stoken   string      `json:"stoken"`
	ResultBs interface{} `json:"result"`
	Result   struct {
		Accurate struct {
			UpUser  interface{} `json:"upuser,omitempty"`
			Bangumi interface{} `json:"bangumi,omitempty"`
		} `json:"accurate,omitempty"`
		Tag []*struct {
			Value string `json:"value,omitempty"`
		} `json:"tag,omitempty"`
	} `json:"-"`
}

// Suggest2 struct
type Suggest2 struct {
	Code   int    `json:"code"`
	Stoken string `json:"stoken"`
	Result *struct {
		Tag []*SuggestTag `json:"tag"`
	} `json:"result"`
}

// SuggestTag struct
type SuggestTag struct {
	Value string `json:"value,omitempty"`
	Ref   int64  `json:"ref,omitempty"`
	Name  string `json:"name,omitempty"`
	SpID  int    `json:"spid,omitempty"`
	Type  string `json:"type,omitempty"`
}

// Suggest3 struct
type Suggest3 struct {
	Code    int    `json:"code"`
	TrackID string `json:"trackid"`
	Result  []*Sug `json:"result"`
}

// Sug struct
type Sug struct {
	ShowName  string  `json:"show_name,omitempty"`
	Term      string  `json:"term,omitempty"`
	Ref       int64   `json:"ref,omitempty"`
	TermType  int     `json:"term_type,omitempty"`
	SubType   string  `json:"sub_type,omitempty"`
	Pos       int     `json:"pos,omitempty"`
	Cover     string  `json:"cover,omitempty"`
	CoverSize float64 `json:"cover_size,omitempty"`
}

// Operate struct
type Operate struct {
	ID          int64  `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Cover       string `json:"cover,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
	Desc        string `json:"desc,omitempty"`
	Corner      string `json:"corner,omitempty"`
	RecReason   string `json:"rec_reason,omitempty"`
	ContentList []*struct {
		Type int   `json:"type,omitempty"`
		ID   int64 `json:"id,omitempty"`
	} `json:"content_list,omitempty"`
}

// Game struct
type Game struct {
	ID          int64   `json:"id,omitempty"`
	Title       string  `json:"title,omitempty"`
	Cover       string  `json:"cover,omitempty"`
	Desc        string  `json:"description,omitempty"`
	View        float64 `json:"view,omitempty"`
	Like        int64   `json:"like,omitempty"`
	Status      int     `json:"status,omitempty"`
	RedirectURL string  `json:"redirect_url,omitempty"`
}

// type Comic struct {
// 	ID     int64  `json:"id,omitempty"`
// 	Title  string `json:"title,omitempty"`
// 	Cover  string `json:"cover,omitempty"`
// 	Uname  string `json:"uname,omitempty"`
// 	Areas  string `json:"areas,omitempty"`
// 	Styles string `json:"styles,omitempty"`
// 	URL    string `json:"url,omitempty"`
// 	Type   string `json:"type,omitempty"`
// }

// Channel struct
type Channel struct {
	Type       string `json:"type,omitempty"`
	TagID      int64  `json:"tag_id,omitempty"`
	TagName    string `json:"tag_name,omitempty"`
	AttenCount int    `json:"atten_count,omitempty"`
	Cover      string `json:"cover,omitempty"`
}

// Twitter twitter.
type Twitter struct {
	ID         int64    `json:"id,omitempty"`
	PicID      int64    `json:"pic_id"`
	Cover      []string `json:"cover,omitempty"`
	CoverCount int      `json:"cover_count,omitempty"`
	Content    string   `json:"content,omitempty"`
}

// Flow struct
type Flow struct {
	LinkType string          `json:"linktype,omitempty"`
	Position int             `json:"position,omitempty"`
	Type     string          `json:"type,omitempty"`
	TypeName string          `json:"type_name,omitempty"`
	Value    json.RawMessage `json:"value,omitempty"`
	Video    *Video
	Live     *Live
	Operate  *Operate
	Article  *Article
	Media    *Media
	User     *User
	Game     *Game
	Query    []*Query
	Twitter  *Twitter
	// Comic    *Comic
	TrackID string `json:"trackid,omitempty"`
}

// Change chagne flow
func (f *Flow) Change() {
	var err error
	switch f.Type {
	case TypeVideo:
		err = json.Unmarshal(f.Value, &f.Video)
	case TypeLive:
		err = json.Unmarshal(f.Value, &f.Live)
	case TypeMediaBangumi, TypeMediaFt:
		err = json.Unmarshal(f.Value, &f.Media)
	case TypeArticle:
		err = json.Unmarshal(f.Value, &f.Article)
	case TypeSpecial, TypeBanner, TypeSpecialS, TypeConverge:
		err = json.Unmarshal(f.Value, &f.Operate)
	case TypeUser, TypeBiliUser:
		err = json.Unmarshal(f.Value, &f.User)
	case TypeGame:
		err = json.Unmarshal(f.Value, &f.Game)
	case TypeQuery:
		err = json.Unmarshal(f.Value, &f.Query)
	// case TypeLiveComic:
	// 	err = json.Unmarshal(f.Value, &f.Comic)
	case TypeTwitter:
		err = json.Unmarshal(f.Value, &f.Twitter)
	}
	if err != nil {
		log.Error("json.Unmarshal(%s) error(%+v)", f.Value, err)
	}
}
