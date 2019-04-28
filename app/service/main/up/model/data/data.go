package data

import (
	xtime "go-common/library/time"

	"go-common/app/interface/main/creative/model/medal"
)

// Stat info
type Stat struct {
	FanLast     int64             `json:"fan_last"`
	Fan         int64             `json:"fan"`
	DmLast      int64             `json:"dm_last"`
	Dm          int64             `json:"dm"`
	CommentLast int64             `json:"comment_last"`
	Comment     int64             `json:"comment"`
	Play        int64             `json:"play"`
	PlayLast    int64             `json:"play_last"`
	Fav         int64             `json:"fav"`
	FavLast     int64             `json:"fav_last"`
	Like        int64             `json:"like"`
	LikeLast    int64             `json:"like_last"`
	Day30       map[string]int    `json:"30,omitempty"`
	Arcs        map[string][]*Arc `json:"arcs"`
}

// Arc Arcs info.
type Arc struct {
	Aid   int64  `json:"aid"`
	Title string `json:"title"`
	Click int64  `json:"click"`
}

// Tags Arcs info
type Tags struct {
	Tags []string `json:"tags"`
}

// CheckedTag tags with checked
type CheckedTag struct {
	Tag     string `json:"tag"`
	Checked int    `json:"checked"`
}

// AppStat arc stat.
type AppStat struct {
	Date string `json:"date"`
	Num  int64  `json:"num"`
}

// AppStatList for arc stat list.
type AppStatList struct {
	Danmu   []*AppStat `json:"danmu"`
	View    []*AppStat `json:"view"`
	Fans    []*AppStat `json:"fans"`
	Comment []*AppStat `json:"comment"`
	Show    int8       `json:"show"`
}

// ViewerTrend for up trend data.
type ViewerTrend struct {
	Tag map[int]string   `json:"tag"`
	Ty  map[string]int64 `json:"ty"`
}

// ViewerIncr for up increment data.
type ViewerIncr struct {
	Arcs      []*ArcInc      `json:"arc_inc"`
	TotalIncr int            `json:"total_inc"`
	TyRank    map[string]int `json:"type_rank"`
}

// ArcInc for archive increment data.
type ArcInc struct {
	AID   int64      `json:"aid"`
	Incr  int        `json:"incr"`
	Title string     `json:"title"`
	PTime xtime.Time `json:"ptime"`
}

// PeriodTip  period tip for data.
type PeriodTip struct {
	ModuleOne   string `json:"module_one"`
	ModuleTwo   string `json:"module_two"`
	ModuleThree string `json:"module_three"`
	ModuleFour  string `json:"module_four"`
}

// AppViewerIncr for up increment data.
type AppViewerIncr struct {
	DateKey   int64     `json:"date_key"`
	Arcs      []*ArcInc `json:"arc_inc"`
	TotalIncr int       `json:"total_inc"`
	TyRank    []*Rank   `json:"type_rank"`
}

// ThirtyDay for 30 days data.
type ThirtyDay struct {
	DateKey   int64 `json:"date_key"`
	TotalIncr int64 `json:"total_inc"`
}

// Rank type rank for up data.
type Rank struct {
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

// CreatorDataShow for display archive/article data module.
type CreatorDataShow struct {
	Archive int `json:"archive"`
	Article int `json:"article"`
}

// AppFan for stat.
type AppFan struct {
	Summary map[string]int64 `json:"summary"`
}

//for fan manager top mids.
const (
	//Total 粉丝管理-累计数据
	Total = iota
	//Seven 粉丝管理-7日数据
	Seven
	//Thirty 粉丝管理-30日数据
	Thirty
	//Ninety 粉丝管理-90日数据
	Ninety
	//PlayDuration 播放时长
	PlayDuration = "video_play"
	//VideoAct 视频互动
	VideoAct = "video_act"
	//DynamicAct 动态互动
	DynamicAct = "dynamic_act"
)

// WebFan for stat.
type WebFan struct {
	RankMap   map[string]map[string]int32  `json:"-"`
	Summary   map[string]int32             `json:"summary"`
	RankList  map[string][]*RankInfo       `json:"rank_list"`
	RankMedal map[string][]*medal.FansRank `json:"rank_medal"`
	Source    map[string]int32             `json:"source"`
}

// RankInfo str
type RankInfo struct {
	MID   int64  `json:"mid"`
	Uname string `json:"uname"`
	Photo string `json:"photo"`
}

// PlaySource for play soucre.
type PlaySource struct {
	PlayProportion map[string]int32 `json:"play_proportion"`
	PageSource     map[string]int32 `json:"page_source"`
}

// ArchivePlay for archive play.
type ArchivePlay struct {
	AID         int64  `json:"aid"`
	View        int32  `json:"view"`
	Rate        int32  `json:"rate"`
	CTime       int32  `json:"ctime"`
	Duration    int64  `json:"duration"`
	AvgDuration int64  `json:"avg_duration"`
	Title       string `json:"title"`
}

//ArchivePlayList for arc play list.
type ArchivePlayList struct {
	ArcPlayList []*ArchivePlay `json:"arc_play_list"`
}
