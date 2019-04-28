package archive

const (
	// ForbidRank 禁止排行
	ForbidRank = "rank"
	// ForbidDynamic 动态禁止
	ForbidDynamic = "dynamic"
	// ForbidRecommend 禁止推荐
	ForbidRecommend = "recommend"
	// ForbidShow 禁止展示
	ForbidShow = "show"
	// ForbidRankMain forbid rank main
	ForbidRankMain = 0
	// ForbidRankRecentArc forbid rank recent archive
	ForbidRankRecentArc = 1
	// ForbidRankAllArc forbid rank all archive
	ForbidRankAllArc = 2
	// ForbidDynamicMain forbid dynamic main
	ForbidDynamicMain = 0
	// ForbidRecommendMain forbid recomment main
	ForbidRecommendMain = 0
	// ForbidShowMain forbid show main
	ForbidShowMain = 0
	// ForbidShowMobile forbid show mobile
	ForbidShowMobile = 1
	// ForbidShowWeb forbid show web
	ForbidShowWeb = 2
	// ForbidShowOversea forbid show oversea
	ForbidShowOversea = 3
	// ForbidShowOnline forbid show online
	ForbidShowOnline = 4
	//ForbidAttrChannel   forbid channel
	ForbidAttrChannel = 101
	//ForbidAttrHot  forbid hot
	ForbidAttrHot = 102
)

var (
	_forbidBits = map[string]map[uint]string{
		ForbidRank: map[uint]string{
			ForbidRankMain:      "所有排行禁止",
			ForbidRankRecentArc: "近期稿件排行禁止",
			ForbidRankAllArc:    "全部稿件排行禁止",
		},
		ForbidDynamic: map[uint]string{
			ForbidDynamicMain: "所有动态禁止",
		},
		ForbidRecommend: map[uint]string{
			ForbidRecommendMain: "所有推荐禁止",
		},
		ForbidShow: map[uint]string{
			ForbidShowMain:    "移动端最新/网页端最新/热度/在线等禁止",
			ForbidShowMobile:  "移动端最新禁止",
			ForbidShowWeb:     "网页端最新禁止",
			ForbidShowOversea: "海外禁止",
			ForbidShowOnline:  "在线列表禁止",
		},
	}
)

// ForbidAttr forbid attribute
type ForbidAttr struct {
	Aid        int64 `json:"aid"`
	OnFlowID   int64 `json:"on_flow_id"`
	RankV      int32 `json:"-"`
	DynamicV   int32 `json:"-"`
	RecommendV int32 `json:"-"`
	ShowV      int32 `json:"-"`
	// specific
	Rank struct {
		Main      int32 `json:"main"`
		RecentArc int32 `json:"recent_arc"`
		AllArc    int32 `json:"all_arc"`
	} `json:"rank_attr"`
	Dynamic struct {
		Main int32 `json:"main"`
	} `json:"dynamic_attr"`
	Recommend struct {
		Main int32 `json:"main"`
	} `json:"recommend_attr"`
	Show struct {
		Main    int32 `json:"main"`
		Mobile  int32 `json:"mobile"`
		Web     int32 `json:"web"`
		Oversea int32 `json:"oversea"`
		Online  int32 `json:"online"`
	} `json:"show_attr"`
}

// Convert convert db value into attr.
func (f *ForbidAttr) Convert() {
	// rank
	f.Rank.Main = f.RankV & 1
	f.Rank.RecentArc = (f.RankV >> 1) & 1
	f.Rank.AllArc = (f.RankV >> 2) & 1
	// dynamic
	f.Dynamic.Main = f.DynamicV & 1
	// recommend
	f.Recommend.Main = f.RecommendV & 1
	// show
	f.Show.Main = f.ShowV & 1
	f.Show.Mobile = (f.ShowV >> 1) & 1
	f.Show.Web = (f.ShowV >> 2) & 1
	f.Show.Oversea = (f.ShowV >> 3) & 1
	f.Show.Online = (f.ShowV >> 4) & 1
}

// Reverse reverse attr into db value.
// func (f *ForbidAttr) Reverse() {
// 	// rank
// 	f.RankV = (f.Rank.AllArc << 2) | (f.Rank.RecentArc << 1) | f.Rank.Main
// 	// dynamic
// 	f.DynamicV = f.Dynamic.Main
// 	// recommend
// 	f.RecommendV = f.Recommend.Main
// 	// show
// 	f.ShowV = (f.Show.Online << 4) | (f.Show.Oversea << 3) | (f.Show.Web << 2) | (f.Show.Mobile << 1) | f.Show.Main
// }

// SetAttr set forbid attr.
func (f *ForbidAttr) SetAttr(name string, v int32, bit uint) (change bool) {
	if name == ForbidRank {
		old := f.RankV
		f.RankV = f.RankV&(^(1 << bit)) | (v << bit)
		change = old == f.RankV
	} else if name == ForbidDynamic {
		old := f.DynamicV
		f.DynamicV = f.DynamicV&(^(1 << bit)) | (v << bit)
		change = old == f.DynamicV
	} else if name == ForbidRecommend {
		old := f.RecommendV
		f.RecommendV = f.RecommendV&(^(1 << bit)) | (v << bit)
		change = old == f.RecommendV
	} else if name == ForbidShow {
		old := f.ShowV
		f.ShowV = f.ShowV&(^(1 << bit)) | (v << bit)
		change = old == f.ShowV
	}
	return
}

// ForbidBitDesc return bit desc.
func ForbidBitDesc(name string, bit uint) (desc string) {
	return _forbidBits[name][bit]
}
