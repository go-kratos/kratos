package archive

//稿件禁止项目
const (
	ForbidRank      = "rank"
	ForbidDynamic   = "dynamic"
	ForbidRecommend = "recommend"
	ForbidShow      = "show"

	ForbidRankMain      = 0
	ForbidRankRecentArc = 1
	ForbidRankAllArc    = 2

	ForbidDynamicMain = 0

	ForbidRecommendMain = 0

	ForbidShowMain    = 0
	ForbidShowMobile  = 1
	ForbidShowWeb     = 2
	ForbidShowOversea = 3
	ForbidShowOnline  = 4
)

// ForbidAttr forbid attribute
type ForbidAttr struct {
	Aid        int64 `json:"aid"`
	OnFlowID   int64 `json:"-"`
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
	}
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
