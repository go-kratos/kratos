package datamodel

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

/* ------------- */

//ViewerTypeData type's play count
type ViewerTypeData struct {
	Tid  int    `json:"tid"`
	Name string `json:"name"`
	Play int64  `json:"play"`
}

//ViewerTagData viewer tag struct
type ViewerTagData struct {
	Idx   int    `json:"idx"`
	TagID int    `json:"tag_id"`
	Name  string `json:"name"`
}

//ViewerTypeTagInfo struct for viewer type and tag
type ViewerTypeTagInfo struct {
	Type []*ViewerTypeData `json:"type"`
	Tag  []*ViewerTagData  `json:"tag"`
}

// ViewerTrendInfo struct for viewer trend
type ViewerTrendInfo struct {
	Fans  *ViewerTypeTagInfo `json:"fans"`
	Guest *ViewerTypeTagInfo `json:"guest"`
}

/* ------------- */

//ViewerAreaData viewer area data
type ViewerAreaData struct {
	Area    string `json:"area"`
	Viewers int64  `json:"viewers"`
}

//ViewerAreaInfo viewer area info
type ViewerAreaInfo struct {
	Fans  []ViewerAreaData `json:"fans"`
	Guest []ViewerAreaData `json:"guest"`
}

/* ------------- */

//ViewerBaseData base data
//f:plat0	web-pc播放
//f:plat1	web-h5播放
//f:plat2	站外播放
//f:plat3	ios播放
//f:plat4	android播放
type ViewerBaseData struct {
	Male         int64 `json:"male"`
	Female       int64 `json:"female"`
	Age1         int64 `json:"age_1"`          //0~16岁
	Age2         int64 `json:"age_2"`          // 16~25
	Age3         int64 `json:"age_3"`          //25~40
	Age4         int64 `json:"age_4"`          // 40+
	PlatPC       int64 `json:"plat_pc"`        // pc
	PlatH5       int64 `json:"plat_h5"`        // h5
	PlatOut      int64 `json:"plat_out"`       // 站外播放
	PlatIOS      int64 `json:"plat_ios"`       // ios播放
	PlatAndroid  int64 `json:"plat_android"`   // android播放
	PlatOtherApp int64 `json:"plat_other_app"` // 其他播放
}

//ViewerBaseInfo base info
type ViewerBaseInfo struct {
	Fans  *ViewerBaseData `json:"fans"`
	Guest *ViewerBaseData `json:"guest"`
}

/* ------------- */

//FanSummaryData fans summary data
type FanSummaryData struct {
	Total        int32 `json:"total" family:"f" qualifier:"all"`            // 粉丝总数
	Active       int32 `json:"active" family:"f" qualifier:"act"`           // 活跃粉丝数
	Inc          int32 `json:"inc" family:"f" qualifier:"inc"`              // 新增粉丝
	Medal        int32 `json:"medal" family:"f" qualifier:"mdl"`            // 领取勋章粉丝
	Elec         int32 `json:"elec" family:"f" qualifier:"elec"`            //充电粉丝
	MedalDiff    int32 `json:"medal_diff" family:"f" qualifier:"mdl_diff"`  //活跃粉丝（增量）
	ActiveDiff   int32 `json:"active_diff" family:"f" qualifier:"act_diff"` //领取勋章粉丝（增量）
	ElecDiff     int32 `json:"elec_diff" family:"f" qualifier:"elec_diff"`  //充电粉丝（增量）
	Inter        int32 `json:"inter" family:"f" qualifier:"inter"`          //观看活跃度*10000
	ViewPercent  int32 `json:"view_percent" family:"f" qualifier:"v"`       //播放粉丝占比*10000
	DmPercent    int32 `json:"dm_percent" family:"f" qualifier:"da"`        //弹幕粉丝占比*10000
	ReplyPercent int32 `json:"reply_percent" family:"f" qualifier:"re"`     //评论粉丝占比*10000
	CoinPercent  int32 `json:"coin_percent" family:"f" qualifier:"co"`      //投币粉丝占比*10000
	FavorPercent int32 `json:"favor_percent" family:"f" qualifier:"fv"`     //收藏粉丝占比*10000
	SharePercent int32 `json:"share_percent" family:"f" qualifier:"sh"`     //分享粉丝占比*10000
	LikePercent  int32 `json:"like_percent" family:"f" qualifier:"lk"`      //点赞粉丝占比*10000
}

//FanInfo fans info
type FanInfo struct {
	Summary FanSummaryData `json:"summary"`
}

//RelationFanHistoryData data for relation fans follow and unfollow
type RelationFanHistoryData struct {
	FollowData   map[string]int `family:"a"`
	UnfollowData map[string]int `family:"u"`
}
