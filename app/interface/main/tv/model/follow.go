package model

// ResFollow is the result structure from PGC API
type ResFollow struct {
	Code    int       `json:"code"`
	Count   string    `json:"count"`
	Pages   string    `json:"pages"`
	Message string    `json:"message"`
	Result  []*Follow `json:"result"`
}

// Up is the uploader info
type Up struct {
	Mid string `json:"mid"`
	Up  string `json:"up"`
}

// EP is the newest EP info
type EP struct {
	AVID       string `json:"av_id"`
	Coins      string `json:"coins"`
	Cover      string `json:"cover"`
	Danmaku    string `json:"danmaku"`
	EpisodeID  string `json:"episode_id"`
	Index      string `json:"index"`
	IndexTitle string `json:"index_title"`
	IsWebplay  string `json:"is_webplay"`
	Page       string `json:"page"`
	Up         *Up    `json:"up"`
	UpdateTime string `json:"update_time"`
	WebplayURL string `json:"webplay_url"`
}

// Tag is the tag info
type Tag struct {
	Bottoms   string   `json:"bottoms"`
	Cover     string   `json:"cover"`
	Index     string   `json:"index"`
	OrderType string   `json:"orderType"`
	Seasons   []string `json:"seasons"`
	StyleID   string   `json:"style_id"`
	TagID     string   `json:"tag_id"`
	TagName   string   `json:"tag_name"`
	Tops      string   `json:"tops"`
	Type      string   `json:"type"`
}

// UserSeason is the user's season info
type UserSeason struct {
	Attention   string `json:"attention"`
	LastEPID    string `json:"last_ep_id"`
	LastEPIndex string `json:"last_ep_index"`
	LastTime    string `json:"last_time"`
}

// Follow is the main structure of the followed season
type Follow struct {
	Actor          []string     `json:"actor"`
	Alias          string       `json:"alias"`
	AllowBP        string       `json:"allow_bp"`
	AllowDownload  string       `json:"allow_download"`
	Area           string       `json:"area"`
	AreaLimit      int          `json:"arealimit"`
	BangumiID      string       `json:"bangumi_id"`
	BangumiTitle   string       `json:"bangumi_title"`
	Brief          string       `json:"brief"`
	Coins          string       `json:"coins"`
	Copyright      string       `json:"copyright"`
	Cover          string       `json:"cover"`
	DanmakuCount   string       `json:"danmaku_count"`
	Episodes       []string     `json:"episodes"`
	EDJump         int          `json:"ed_jump"`
	Evaluate       string       `json:"evaluate"`
	Favorites      string       `json:"favorites"`
	IsFinish       string       `json:"is_finish"`
	Progress       string       `json:"progress"`
	NewEP          *EP          `json:"new_ep"`
	NewestEPID     string       `json:"newest_ep_id"`
	NewestEPIndex  string       `json:"newest_ep_index"`
	PlayCount      string       `json:"play_count"`
	PubTime        string       `json:"pub_time"`
	RelatedSeasons []string     `json:"related_seasons"`
	SeasonID       string       `json:"season_id"`
	SeasonTitle    string       `json:"season_title"`
	Seasons        []string     `json:"seasons"`
	ShareURL       string       `json:"share_url"`
	SPID           string       `json:"spid"`
	SquareCover    string       `json:"squareCover"`
	Staff          string       `json:"staff"`
	Tag2s          []string     `json:"tag2s"`
	Tags           []*Tag       `json:"tags"`
	Title          string       `json:"title"`
	TotalCount     string       `json:"total_count"`
	UserSeason     *UserSeason  `json:"user_season"`
	Weekday        string       `json:"weekday"`
	CornerMark     *SnVipCorner `json:"cornermark"`
}
