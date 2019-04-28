package search

import (
	"encoding/json"

	"go-common/app/interface/main/app-interface/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// search const
const (
	TypeVideo          = "video"
	TypeLive           = "live_room"
	TypeMediaBangumi   = "media_bangumi"
	TypeMediaFt        = "media_ft"
	TypeArticle        = "article"
	TypeSpecial        = "special_card"
	TypeBanner         = "banner"
	TypeUser           = "user"
	TypeBiliUser       = "bili_user"
	TypeGame           = "game"
	TypeSpecialS       = "special_card_small"
	TypeConverge       = "content_card"
	TypeQuery          = "query"
	TypeLiveMaster     = "live_master"
	TypeTwitter        = "twitter"
	TypeComic          = "comic"
	TypeStar           = "star"
	TypeTicket         = "ticket"
	TypeProduct        = "product"
	TypeSpecialerGuide = "special_guide_card"
	TypeChannel        = "tag"

	SuggestionJump     = 99
	SuggestionJumpUser = 81
	SuggestionJumpPGC  = 82
	SuggestionAV       = "video"
	SuggestionLive     = "live"
	SuggestionArticle  = "article"

	SearchLiveAllAndroid = 5275000
	SearchLiveAllIOS     = 6800

	SearchEggInfoAndroid = 5270000

	LiveBroadcastTypeAndroid = 5305000

	SearchTwitterAndroid = 5315000
	SearchTwitterIOS     = 8111

	SearchNewIPad   = 8231
	SearchNewIPadHD = 12041

	SearchConvergeIOS     = 8140
	SearchConvergeAndroid = 5320000

	SearchStarIOS     = 8220
	SearchStarAndroid = 5335000

	SearchTicketIOS     = 8220
	SearchTicketAndroid = 5335000

	SearchProductIOS     = 8220
	SearchProductAndroid = 5335000
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
	CrrQuery       string `json:"crr_query,omitempty"`
	Attribute      int32  `json:"exp_bits,omitempty"`
	PageInfo       struct {
		Bangumi      *Page `json:"bangumi,omitempty"`
		UpUser       *Page `json:"upuser,omitempty"`
		BiliUser     *Page `json:"bili_user,omitempty"`
		User         *Page `json:"user,omitempty"`
		Movie        *Page `json:"movie,omitempty"`
		Film         *Page `json:"pgc,omitempty"`
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
		IsPay    int         `json:"is_pay,omitempty"`
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
	ID         int64       `json:"id"`
	Author     string      `json:"author"`
	Title      string      `json:"title"`
	Pic        string      `json:"pic"`
	Desc       string      `json:"description"`
	Play       interface{} `json:"play"`
	Danmaku    int         `json:"video_review"`
	Duration   string      `json:"duration"`
	Pages      int         `json:"numPages"`
	ViewType   string      `json:"view_type"`
	RecTags    []string    `json:"rec_tags"`
	IsPay      int         `json:"is_pay"`
	NewRecTags []*RecTag   `json:"new_rec_tags"`
}

// RecTag from video
type RecTag struct {
	Name  string `json:"tag_name"`
	Style int8   `json:"tag_style"`
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
	MediaType   int             `json:"media_type,omitempty"`
	CV          string          `json:"cv,omitempty"`
	Staff       string          `json:"staff,omitempty"`
	Areas       string          `json:"areas,omitempty"`
	GotoURL     string          `json:"goto_url,omitempty"`
	Pubtime     xtime.Time      `json:"pubtime,omitempty"`
	HitColumns  []string        `json:"hit_columns,omitempty"`
	AllNetName  string          `json:"all_net_name,omitempty"`
	AllNetIcon  string          `json:"all_net_icon,omitempty"`
	AllNetURL   string          `json:"all_net_url,omitempty"`
	DisplayInfo json.RawMessage `json:"display_info,omitempty"`
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
	ShowName  string          `json:"show_name,omitempty"`
	Term      string          `json:"term,omitempty"`
	Ref       int64           `json:"ref,omitempty"`
	TermType  int             `json:"term_type,omitempty"`
	SubType   string          `json:"sub_type,omitempty"`
	Pos       int             `json:"pos,omitempty"`
	Cover     string          `json:"cover,omitempty"`
	CoverSize float64         `json:"cover_size,omitempty"`
	Value     json.RawMessage `json:"value,omitempty"`
	PGC       *SugPGC         `json:"-"`
	User      *SugUser        `json:"-"`
}

// SugPGC fro sug
type SugPGC struct {
	MediaID        int64                `json:"media_id,omitempty"`
	SeasonID       int64                `json:"season_id,omitempty"`
	Title          string               `json:"title,omitempty"`
	MediaType      int                  `json:"media_type,omitempty"`
	GotoURL        string               `json:"goto_url,omitempty"`
	Areas          string               `json:"areas,omitempty"`
	Pubtime        xtime.Time           `json:"pubtime,omitempty"`
	FixPubTime     string               `json:"fix_pubtime_str,omitempty"`
	Styles         string               `json:"styles,omitempty"`
	CV             string               `json:"cv,omitempty"`
	Staff          string               `json:"staff,omitempty"`
	MediaScore     float64              `json:"media_score,omitempty"`
	MediaUserCount int                  `json:"media_user_cnt,omitempty"`
	Cover          string               `json:"cover,omitempty"`
	Badges         []*model.ReasonStyle `json:"badges,omitempty"`
}

// SugUser fro sug
type SugUser struct {
	Mid                int64  `json:"uid,omitempty"`
	Face               string `json:"face,omitempty"`
	Name               string `json:"uname,omitempty"`
	Fans               int    `json:"fans,omitempty"`
	Videos             int    `json:"videos,omitempty"`
	Level              int    `json:"level,omitempty"`
	OfficialVerifyType int    `json:"verify_type,omitempty"`
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

// Comic struct
type Comic struct {
	ID       int64    `json:"id,omitempty"`
	Title    string   `json:"title,omitempty"`
	Author   []string `json:"author,omitempty"`
	Cover    string   `json:"cover,omitempty"`
	Styles   string   `json:"styles,omitempty"`
	URL      string   `json:"url,omitempty"`
	ComicURL string   `json:"sq_url,omitempty"`
}

// Channel struct
type Channel struct {
	Type       string  `json:"type,omitempty"`
	TagID      int64   `json:"tag_id,omitempty"`
	TagName    string  `json:"tag_name,omitempty"`
	AttenCount int     `json:"atten_count,omitempty"`
	Cover      string  `json:"cover,omitempty"`
	Desc       string  `json:"desc,omitempty"`
	Values     []*Flow `json:"value_list,omitempty"`
}

// Twitter twitter.
type Twitter struct {
	ID         int64    `json:"id,omitempty"`
	PicID      int64    `json:"pic_id"`
	Cover      []string `json:"cover,omitempty"`
	CoverCount int      `json:"cover_count,omitempty"`
	Content    string   `json:"content,omitempty"`
}

// Star struct
type Star struct {
	ID      int64  `json:"id,omitempty"`
	Cover   string `json:"cover,omitempty"`
	Desc    string `json:"desc,omitempty"`
	Title   string `json:"title,omitempty"`
	MID     int64  `json:"mid,omitempty"`
	TagID   int64  `json:"tag_id,omitempty"`
	TagList []*struct {
		TagName   string `json:"tagname,omitempty"`
		KeyWord   string `json:"searchtagname,omitempty"`
		ValueList []*struct {
			Type  string `json:"type,omitempty"`
			Video *Video `json:"values,omitempty"`
		} `json:"value_list,omitempty"`
	} `json:"tag_list,omitempty"`
}

// Ticket for search.
type Ticket struct {
	ID        int64  `json:"id,omitempty"`
	Title     string `json:"project_name,omitempty"`
	Cover     string `json:"cover,omitempty"`
	ShowTime  string `json:"show_time,omitempty"`
	CityName  string `json:"city_name,omitempty"`
	VenueName string `json:"venue_name,omitempty"`
	PriceLow  int    `json:"price_low,omitempty"`
	PriceType int    `json:"need_up,omitempty"`
	ReqNum    int    `json:"required_number,omitempty"`
	URL       string `json:"url,omitempty"`
}

// Product for search.
type Product struct {
	ID        int64  `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Cover     string `json:"cover,omitempty"`
	ShopName  string `json:"shop_name,omitempty"`
	Price     int    `json:"price,omitempty"`
	PriceType int    `json:"need_up,omitempty"`
	ReqNum    int    `json:"required_number,omitempty"`
	URL       string `json:"url,omitempty"`
}

// SpecialerGuide fro search
type SpecialerGuide struct {
	ID    int64  `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"`
	Cover string `json:"cover,omitempty"`
	Tel   string `json:"tel,omitempty"`
}

// Flow struct
type Flow struct {
	LinkType       string          `json:"linktype,omitempty"`
	Position       int             `json:"position,omitempty"`
	Type           string          `json:"type,omitempty"`
	TypeName       string          `json:"type_name,omitempty"`
	Value          json.RawMessage `json:"value,omitempty"`
	Video          *Video
	Live           *Live
	Operate        *Operate
	Article        *Article
	Media          *Media
	User           *User
	Game           *Game
	Query          []*Query
	Twitter        *Twitter
	Comic          *Comic
	Star           *Star
	Ticket         *Ticket
	Product        *Product
	SpecialerGuide *SpecialerGuide
	Channel        *Channel
	TrackID        string `json:"trackid,omitempty"`
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
	case TypeComic:
		err = json.Unmarshal(f.Value, &f.Comic)
	case TypeTwitter:
		err = json.Unmarshal(f.Value, &f.Twitter)
	case TypeStar:
		err = json.Unmarshal(f.Value, &f.Star)
	case TypeTicket:
		err = json.Unmarshal(f.Value, &f.Ticket)
	case TypeProduct:
		err = json.Unmarshal(f.Value, &f.Product)
	case TypeSpecialerGuide:
		err = json.Unmarshal(f.Value, &f.SpecialerGuide)
	case TypeChannel:
		if err = json.Unmarshal(f.Value, &f.Channel); err == nil {
			if f.Channel != nil && len(f.Channel.Values) > 0 {
				for _, value := range f.Channel.Values {
					value.Change()
				}
			}
		}
	}
	if err != nil {
		log.Error("Change json.Unmarshal(%s) error(%+v)", f.Value, err)
	}
}

// SugChange chagne sug value
func (s *Sug) SugChange() {
	var err error
	switch s.TermType {
	case SuggestionJumpUser:
		err = json.Unmarshal(s.Value, &s.PGC)
	case SuggestionJumpPGC:
		err = json.Unmarshal(s.Value, &s.User)
	}
	if err != nil {
		log.Error("SugChange json.Unmarshal(%s) error(%+v)", s.Value, err)
	}
}
