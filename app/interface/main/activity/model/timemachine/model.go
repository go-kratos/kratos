package timemachine

import "go-common/app/service/main/archive/api"

// Timemachine .
type Timemachine struct {
	Mid                 int64      `json:"mid"`
	Face                string     `json:"face"`
	Uname               string     `json:"uname"`
	RegTime             string     `json:"reg_time"`
	RegDay              int64      `json:"reg_day"`
	IsUp                int64      `json:"is_up"`
	DurationHour        int64      `json:"duration_hour"`
	ArchiveVv           int64      `json:"archive_vv"`
	LikeTagID           int64      `json:"like_tag_id"`
	LikeTagName         string     `json:"like_tag_name"`
	LikeTagDescFirst    string     `json:"like_tag_desc_first"`
	LikeTagDescSecond   string     `json:"like_tag_desc_second"`
	LikeTagDescSecond2  string     `json:"like_tag_desc_second2"`
	LikeSubtidVv        int64      `json:"like_subtid_vv"`
	Likes3Arcs          []*TmArc   `json:"likes_3arcs"`
	LikeBestUpID        int64      `json:"like_best_upid"`
	LikeBestUpName      string     `json:"like_best_up_name"`
	LikeBestUpFace      string     `json:"like_best_up_face"`
	LikeUpAvDuration    int64      `json:"like_up_av_duration"`
	LikeUpLiveDuration  int64      `json:"like_up_live_duration"`
	LikeUpDuration      int64      `json:"like_up_duration"`
	LikeUp3Arcs         []*TmArc   `json:"like_up_3arcs"`
	LikeLiveUpSubTname  string     `json:"like_live_up_sub_tname"`
	BrainwashCirTime    string     `json:"brainwash_cir_time"`
	BrainwashCirArc     *TmArc     `json:"brainwash_cir_arc"`
	BrainwashCirVv      int64      `json:"brainwash_cir_vv"`
	FirstSubmitArc      *TmArc     `json:"first_submit_arc"`
	FirstSubmitTime     string     `json:"first_submit_time"`
	FirstSubmitType     int64      `json:"first_submit_type"`
	SubmitAvsRds        string     `json:"submit_avs_rds"`
	BestArc             *TmArc     `json:"best_arc"`
	BestAvidType        int64      `json:"best_avid_type"`
	BestArcOld          *TmArc     `json:"best_arc_old"`
	BestAvidOldType     int64      `json:"best_avid_old_type"`
	OldAvVv             int64      `json:"old_av_vv"`
	AllVv               int64      `json:"all_vv"`
	UpLiveDuration      int64      `json:"up_live_duration"`
	IsLiveUp            int64      `json:"is_live_up"`
	ValidLiveDays       int64      `json:"valid_live_days"`
	MaxCdnNumDate       string     `json:"max_cdn_num_date"`
	MaxCdnNum           int64      `json:"max_cdn_num"`
	AddAttentions       int64      `json:"add_attentions"`
	Fans                int64      `json:"fans"`
	UpBestFanVv         *FavVv     `json:"up_best_fan_vv"`
	UpBestFanLiveMinute *FanMinute `json:"up_best_fan_live_minute"`
	WinRatio            string     `json:"win_ratio"`
	Like2Tnames         string     `json:"like_2tnames"`
	Like2SubTnames      string     `json:"like_2sub_tnames"`
	LikeSubDesc1        string     `json:"like_sub_desc1"`
	LikeSubDesc2        string     `json:"like_sub_desc2"`
	LikeSubDesc3        string     `json:"like_sub_desc3"`
}

// AidView aid view.
type AidView struct {
	Aid  int64 `json:"aid"`
	View int64 `json:"view"`
}

// FavVv .
type FavVv struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
	Face string `json:"face"`
	Vv   int64  `json:"vv"`
}

// FanMinute .
type FanMinute struct {
	Mid    int64  `json:"mid"`
	Name   string `json:"name"`
	Face   string `json:"face"`
	Minute int64  `json:"minute"`
}

// TmArc time machine arc.
type TmArc struct {
	Aid    int64      `json:"aid"`
	Title  string     `json:"title"`
	Cover  string     `json:"cover"`
	Author api.Author `json:"author"`
}

// TagDesc tag desc.
type TagDesc struct {
	TagID      int64  `json:"tag_id"`
	TagName    string `json:"tag_name"`
	Desc1      string `json:"desc1"`
	Desc2Line1 string `json:"desc2_line1"`
	Desc2Line2 string `json:"desc2_line2"`
}

// TagRegionDesc tag region desc.
type TagRegionDesc struct {
	RID        int64  `json:"rid"`
	Desc1      string `json:"desc1"`
	Desc2Line1 string `json:"desc2_line1"`
	Desc2Line2 string `json:"desc2_line2"`
}

// RegionDesc region desc.
type RegionDesc struct {
	RID   int64  `json:"rid"`
	Desc1 string `json:"desc1"`
	Desc2 string `json:"desc2"`
	Desc3 string `json:"desc3"`
}
