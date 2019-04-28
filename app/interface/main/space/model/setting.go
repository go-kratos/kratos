package model

// DefaultPrivacy default privacy.
var (
	PcyBangumi     = "bangumi"
	PcyTag         = "tags"
	PcyFavVideo    = "fav_video"
	PcyCoinVideo   = "coins_video"
	PcyGroup       = "groups"
	PcyGame        = "played_game"
	PcyChannel     = "channel"
	PcyUserInfo    = "user_info"
	PcyLikeVideo   = "likes_video"
	DefaultPrivacy = map[string]int{
		PcyBangumi:   1,
		PcyTag:       1,
		PcyFavVideo:  1,
		PcyCoinVideo: 1,
		PcyGroup:     1,
		PcyGame:      1,
		PcyChannel:   1,
		PcyUserInfo:  1,
		PcyLikeVideo: 1,
	}
	DefaultIndexOrder = []*IndexOrder{
		{ID: 1, Name: "我的稿件"},
		{ID: 8, Name: "我的专栏"},
		{ID: 7, Name: "我的频道"},
		{ID: 2, Name: "我的收藏夹"},
		{ID: 3, Name: "订阅番剧"},
		{ID: 4, Name: "订阅标签"},
		{ID: 5, Name: "最近投币的视频"},
		{ID: 6, Name: "我的圈子"},
		{ID: 9, Name: "我的相簿"},
		{ID: 21, Name: "公告"},
		{ID: 22, Name: "直播间"},
		{ID: 23, Name: "个人资料"},
		{ID: 24, Name: "官方活动"},
		{ID: 25, Name: "最近玩过的游戏"},
	}
	IndexOrderMap = indexOrderMap()
)

// Setting setting struct.
type Setting struct {
	Privacy    map[string]int `json:"privacy"`
	IndexOrder []*IndexOrder  `json:"index_order"`
}

// Privacy privacy struct.
type Privacy struct {
	Privacy string `json:"privacy"`
	Status  int    `json:"status"`
}

// IndexOrder index order struct.
type IndexOrder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func indexOrderMap() map[int]string {
	data := make(map[int]string, len(DefaultIndexOrder))
	for _, v := range DefaultIndexOrder {
		data[v.ID] = v.Name
	}
	return data
}
