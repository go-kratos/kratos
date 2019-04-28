package model

// Setting .
type Setting struct {
	// 是否开启申请
	ApplyOpen bool `json:"apply_info"`
	// 申请限制
	ApplyLimit int64 `json:"apply_limit"`
	// 申请被拒绝后的冷冻期
	ApplyFrozenDuration int64 `json:"frozen_duration"`
	// 是否在推荐页面展示最新投稿
	ShowRecommendNewArticles bool `json:"show_rec_new_arts"`
	// 是否展示web端排行榜的说明
	ShowRankNote bool `json:"show_rank_note"`
	// 是否展示app专栏主要的排行榜
	ShowAppHomeRank bool `json:"show_app_home_rank"`
	// 详情页展示稍后在看
	ShowLaterWatch bool `json:"show_later_watch"`
	// 详情页展示小窗播放
	ShowSmallWindow bool `json:"show_small_window"`
	// 热点标签
	ShowHotspot bool `json:"show_hotspot"`
}
