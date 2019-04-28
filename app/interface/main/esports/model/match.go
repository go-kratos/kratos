package model

import (
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/time"
)

// Filter filter struct
type Filter struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	SubTitle  string `json:"sub_title"`
	Logo      string `json:"logo"`
	Rank      int    `json:"rank"`
	URL       string `json:"url"`
	DataFocus string `json:"data_focus"`
	FocusURL  string `json:"focus_url"`
}

// Year year struct
type Year struct {
	ID   int64 `json:"id"`
	Year int64 `json:"year"`
	Aid  int64 `json:"aid"`
}

// Calendar calendar struct
type Calendar struct {
	Stime string `json:"stime"`
	Count int64  `json:"count"`
}

// Season season struct
type Season struct {
	ID        int64     `json:"id"`
	Mid       int64     `json:"mid"`
	Title     string    `json:"title"`
	SubTitle  string    `json:"sub_title"`
	Stime     int64     `json:"stime"`
	Etime     int64     `json:"etime"`
	Sponsor   string    `json:"sponsor"`
	Logo      string    `json:"logo"`
	Dic       string    `json:"dic"`
	Status    int64     `json:"status"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
	Rank      int64     `json:"rank"`
	IsApp     int64     `json:"is_app"`
	URL       string    `json:"url"`
	DataFocus string    `json:"data_focus"`
	FocusURL  string    `json:"focus_url"`
}

// Contest contest struct
type Contest struct {
	ID              int64       `json:"id"`
	GameStage       string      `json:"game_stage"`
	Stime           int64       `json:"stime"`
	Etime           int64       `json:"etime"`
	HomeID          int64       `json:"home_id"`
	AwayID          int64       `json:"away_id"`
	HomeScore       int64       `json:"home_score"`
	AwayScore       int64       `json:"away_score"`
	LiveRoom        int64       `json:"live_room"`
	Aid             int64       `json:"aid"`
	Collection      int64       `json:"collection"`
	GameState       int64       `json:"game_state"`
	Dic             string      `json:"dic"`
	Ctime           string      `json:"ctime"`
	Mtime           string      `json:"mtime"`
	Status          int64       `json:"status"`
	Sid             int64       `json:"sid"`
	Mid             int64       `json:"mid"`
	Season          interface{} `json:"season"`
	HomeTeam        interface{} `json:"home_team"`
	AwayTeam        interface{} `json:"away_team"`
	Special         int         `json:"special"`
	SuccessTeam     int64       `json:"success_team"`
	SuccessTeaminfo interface{} `json:"success_teaminfo"`
	SpecialName     string      `json:"special_name"`
	SpecialTips     string      `json:"special_tips"`
	SpecialImage    string      `json:"special_image"`
	Playback        string      `json:"playback"`
	CollectionURL   string      `json:"collection_url"`
	LiveURL         string      `json:"live_url"`
	DataType        int64       `json:"data_type"`
	MatchID         int64       `json:"match_id"`
}

// ContestsData contest data struct
type ContestsData struct {
	ID         int64  `json:"id"`
	Cid        int64  `json:"cid"`
	URL        string `json:"url"`
	PointData  int64  `json:"point_data"`
	GameStatus int64  `json:"game_status"`
	DataType   int64  `json:"-"`
}

//ContestDataPage contest data pager
type ContestDataPage struct {
	Contest *Contest        `json:"contest"`
	Detail  []*ContestsData `json:"detail"`
}

// ElaSub elasticsearch sub contest.
type ElaSub struct {
	SeasonStime int64 `json:"season_stime"`
	Mid         int64 `json:"mid"`
	Stime       int64 `json:"stime"`
	Oid         int64 `json:"oid"`
	State       int64 `json:"state"`
	Sid         int64 `json:"sid"`
}

// Tree match Active
type Tree struct {
	ID        int64 `json:"id" form:"id"`
	MaID      int64 `json:"ma_id,omitempty" form:"ma_id" validate:"required"`
	MadID     int64 `json:"mad_id,omitempty" form:"mad_id" validate:"required"`
	Pid       int64 `json:"pid" form:"pid"`
	RootID    int64 `json:"root_id" form:"root_id"`
	GameRank  int64 `json:"game_rank,omitempty" form:"game_rank" validate:"required"`
	Mid       int64 `json:"mid" form:"mid"`
	IsDeleted int   `json:"is_deleted,omitempty" form:"is_deleted"`
}

// Team .
type Team struct {
	ID         int64  `json:"id" form:"id"`
	Title      string `json:"title" form:"title" validate:"required"`
	SubTitle   string `json:"sub_title" form:"sub_title"`
	ETitle     string `json:"e_title" form:"e_title"`
	CreateTime int64  `json:"create_time" form:"create_time"`
	Area       string `json:"area" form:"area"`
	Logo       string `json:"logo" form:"logo" validate:"required"`
	UID        int64  `json:"uid" form:"uid" gorm:"column:uid"`
	Members    string `json:"members" form:"members"`
	Dic        string `json:"dic" form:"dic"`
	IsDeleted  int    `json:"is_deleted" form:"is_deleted"`
}

// ContestInfo .
type ContestInfo struct {
	*Contest
	HomeName    string `json:"home_name"`
	AwayName    string `json:"away_name"`
	SuccessName string `json:"success_name" form:"success_name"`
}

// TreeList .
type TreeList struct {
	*Tree
	*ContestInfo
}

// Active match Active
type Active struct {
	ID           int64  `json:"id"`
	Mid          int64  `json:"mid"`
	Sid          int64  `json:"sid"`
	Background   string `json:"background"`
	Liveid       int64  `json:"live_id"`
	Intr         string `json:"intr"`
	Focus        string `json:"focus"`
	URL          string `json:"url"`
	BackColor    string `json:"back_color"`
	ColorStep    string `json:"color_step"`
	H5Background string `json:"h5_background"`
	H5BackColor  string `json:"h5_back_color"`
	IntrLogo     string `json:"intr_logo"`
	IntrTitle    string `json:"intr_title"`
	IntrText     string `json:"intr_text"`
	H5Focus      string `json:"h5_focus"`
	H5Url        string `json:"h5_url"`
}

// Module match module
type Module struct {
	ID   int64  `json:"id"`
	MAid int64  `json:"ma_id"`
	Name string `json:"name"`
	Oids string `json:"oids"`
}

//ActiveDetail 活动页数据模块
type ActiveDetail struct {
	ID           int64  `json:"id"`
	Maid         int64  `json:"ma_id"`
	GameType     int    `json:"game_type"`
	STime        int64  `json:"stime"`
	ETime        int64  `json:"etime"`
	ScoreID      int64  `json:"score_id"`
	GameStage    string `json:"game_stage"`
	KnockoutType int    `json:"knockout_type"`
	WinnerType   int    `json:"winner_type"`
	Online       int    `json:"online"`
}

//ActivePage 活动页
type ActivePage struct {
	Active       *Active         `json:"active,omitempty"`
	Videos       []*arcmdl.Arc   `json:"video_first,omitempty"`
	Modules      []*Module       `json:"video_module,omitempty"`
	ActiveDetail []*ActiveDetail `json:"active_detail,omitempty"`
	Season       *Season         `json:"season,omitempty"`
}
