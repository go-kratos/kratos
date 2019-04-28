package model

import (
	"encoding/json"

	accmdl "go-common/app/service/main/account/model"
)

// WxSearchType .
const WxSearchType = "wechat"

// Search all search.
type Search struct {
	Code           int             `json:"code,omitempty"`
	SeID           string          `json:"seid,omitempty"`
	Page           int             `json:"page,omitempty"`
	PageSize       int             `json:"pagesize,omitempty"`
	Total          int             `json:"total,omitempty"`
	NumResults     int             `json:"numResults"`
	NumPages       int             `json:"numPages"`
	SuggestKeyword string          `json:"suggest_keyword"`
	RqtType        string          `json:"rqt_type,omitempty"`
	CostTime       json.RawMessage `json:"cost_time,omitempty"`
	ExpList        json.RawMessage `json:"exp_list,omitempty"`
	EggHit         int             `json:"egg_hit"`
	PageInfo       json.RawMessage `json:"pageinfo,omitempty"`
	Result         json.RawMessage `json:"result,omitempty"`
	TopTList       json.RawMessage `json:"top_tlist,omitempty"`
	EggInfo        *struct {
		ID     int64 `json:"id,omitempty"`
		Source int64 `json:"source,omitempty"`
	} `json:"egg_info,omitempty"`
	ShowColumn int `json:"show_column"`
}

// SearchTypeRes search type res.
type SearchTypeRes struct {
	Code           int             `json:"code,omitempty"`
	SeID           string          `json:"seid,omitempty"`
	Page           int             `json:"page,omitempty"`
	PageSize       int             `json:"pagesize,omitempty"`
	Total          int             `json:"total,omitempty"`
	NumResults     int             `json:"numResults"`
	NumPages       int             `json:"numPages"`
	SuggestKeyword string          `json:"suggest_keyword"`
	RqtType        string          `json:"rqt_type,omitempty"`
	CostTime       json.RawMessage `json:"cost_time,omitempty"`
	ExpList        json.RawMessage `json:"exp_list,omitempty"`
	EggHit         int             `json:"egg_hit"`
	PageInfo       json.RawMessage `json:"pageinfo,omitempty"`
	Result         json.RawMessage `json:"result,omitempty"`
	ShowColumn     int             `json:"show_column"`
}

// SearchRec search recommend.
type SearchRec struct {
	Code           int             `json:"code,omitempty"`
	SeID           string          `json:"seid,omitempty"`
	Page           int             `json:"page,omitempty"`
	PageSize       int             `json:"pagesize,omitempty"`
	Total          int             `json:"total,omitempty"`
	NumResults     int             `json:"numResults"`
	NumPages       int             `json:"numPages"`
	SuggestKeyword string          `json:"suggest_keyword"`
	RqtType        string          `json:"rqt_type,omitempty"`
	CostTime       json.RawMessage `json:"cost_time,omitempty"`
	ExpList        json.RawMessage `json:"exp_list,omitempty"`
	EggHit         int             `json:"egg_hit"`
	Result         json.RawMessage `json:"result,omitempty"`
}

// SearchAllArg search all api arguments.
type SearchAllArg struct {
	Pn           int    `form:"page"`
	Keyword      string `form:"keyword" validate:"required"`
	Rid          int    `form:"tids"`
	Duration     int    `form:"duration" validate:"gte=0,lte=4"`
	FromSource   string `form:"from_source"`
	Highlight    int    `form:"highlight"`
	SingleColumn int    `form:"-"`
}

// SearchTypeArg search type api arguments.
type SearchTypeArg struct {
	Pn         int    `form:"page" validate:"min=1" default:"1"`
	SearchType string `form:"search_type" validate:"required"`
	Keyword    string `form:"keyword" validate:"required"`
	Order      string `form:"order"`
	Rid        int64  `form:"tids"`
	FromSource string `form:"from_source"`
	Platform   string `form:"platform" default:"web"`
	Duration   int    `form:"duration" validate:"min=0,max=4"`
	// article
	CategoryID int64 `form:"category_id"`
	// special
	VpNum int `form:"vp_num"`
	// bili user
	BiliUserVl   int `form:"bili_user_vl" default:"3"`
	UserType     int `form:"user_type" validate:"min=0,max=3"`
	OrderSort    int `form:"order_sort"`
	Highlight    int `form:"highlight"`
	SingleColumn int `form:"-"`
}

// SearchDefault search default
type SearchDefault struct {
	Trackid  string `json:"seid"`
	ID       int64  `json:"id"`
	ShowName string `json:"show_name"`
	Name     string `json:"name"`
	Type     int    `json:"type"`
}

// SearchUpRecArg search up rec arg.
type SearchUpRecArg struct {
	ServiceArea string  `form:"service_area" validate:"required"`
	Platform    string  `form:"platform" validate:"required"`
	ContextID   int64   `form:"context_id"`
	MainTids    []int64 `form:"main_tids,split"`
	SubTids     []int64 `form:"sub_tids,split"`
	MobiApp     string  `form:"mobi_app"`
	Device      string  `form:"device"`
	Build       int64   `form:"build"`
	Ps          int     `form:"ps" default:"5" validate:"min=1,max=15"`
}

// SearchUpRecRes .
type SearchUpRecRes struct {
	UpID      int64  `json:"up_id"`
	RecReason string `json:"rec_reason"`
	Tid       int16  `json:"tid"`
	SecondTid int16  `json:"second_tid"`
}

// UpRecInfo .
type UpRecInfo struct {
	Mid      int64               `json:"mid"`
	Name     string              `json:"name"`
	Face     string              `json:"face"`
	Official accmdl.OfficialInfo `json:"official"`
	Follower int64               `json:"follower"`
	Vip      struct {
		Type   int32 `json:"type"`
		Status int32 `json:"status"`
	} `json:"vip"`
	RecReason   string `json:"rec_reason"`
	Tid         int16  `json:"tid"`
	Tname       string `json:"tname"`
	SecondTid   int16  `json:"second_tid"`
	SecondTname string `json:"second_tname"`
	Sign        string `json:"sign"`
}

// UpRecData .
type UpRecData struct {
	TrackID string       `json:"track_id"`
	List    []*UpRecInfo `json:"list"`
}

// SearchEgg .
type SearchEgg struct {
	Plat map[int64][]*struct {
		EggID int64  `json:"egg_id"`
		Plat  int    `json:"plat"`
		URL   string `json:"url"`
		MD5   string `json:"md5"`
		Size  int64  `json:"size"`
	} `json:"plat"`
	ShowCount int `json:"show_count"`
}

// SearchEggRes .
type SearchEggRes struct {
	EggID     int64              `json:"egg_id"`
	ShowCount int                `json:"show_count"`
	Source    []*SearchEggSource `json:"source"`
}

// SearchEggSource .
type SearchEggSource struct {
	URL  string `json:"url"`
	MD5  string `json:"md5"`
	Size int64  `json:"size"`
}

// SearchType search types
const (
	SearchTypeAll      = "all"
	SearchTypeVideo    = "video"
	SearchTypeBangumi  = "media_bangumi"
	SearchTypePGC      = "media_ft"
	SearchTypeLive     = "live"
	SearchTypeLiveRoom = "live_room"
	SearchTypeLiveUser = "live_user"
	SearchTypeArticle  = "article"
	SearchTypeSpecial  = "special"
	SearchTypeTopic    = "topic"
	SearchTypeUser     = "bili_user"
	SearchTypePhoto    = "photo"
	WxSearchTypeAll    = "wx_all"
)

// SearchDefaultArg search default params.
var SearchDefaultArg = map[string]map[string]int{
	SearchTypeAll: {
		"highlight":         1,
		"video_num":         20,
		"media_bangumi_num": 3,
		"media_ft_num":      3,
		"is_new_pgc":        1,
		"live_room_num":     1,
		"card_num":          1,
		"activity":          1,
		"bili_user_num":     1,
		"bili_user_vl":      3,
		"user_num":          1,
		"user_video_limit":  3,
		"is_star":           1,
	},
	SearchTypeVideo: {
		"highlight":  1,
		"pagesize":   20,
		"is_new_pgc": 1,
	},
	SearchTypeBangumi: {
		"highlight": 1,
		"pagesize":  20,
	},
	SearchTypePGC: {
		"highlight": 1,
		"pagesize":  20,
	},
	SearchTypeLive: {
		"highlight":     1,
		"live_user_num": 6,
		"live_room_num": 40,
	},
	SearchTypeLiveRoom: {
		"highlight": 1,
		"pagesize":  40,
	},
	SearchTypeLiveUser: {
		"highlight": 1,
		"pagesize":  30,
	},
	SearchTypeArticle: {
		"highlight": 1,
		"pagesize":  20,
	},
	SearchTypeSpecial: {
		"pagesize": 20,
	},
	SearchTypeTopic: {
		"pagesize": 20,
	},
	SearchTypeUser: {
		"highlight": 1,
		"pagesize":  20,
	},
	SearchTypePhoto: {
		"pagesize": 20,
	},
	WxSearchTypeAll: {
		"video_num":         20,
		"media_bangumi_num": 3,
		"media_ft_num":      3,
		"is_new_pgc":        1,
	},
}
