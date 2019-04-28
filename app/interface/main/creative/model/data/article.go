package data

import xtime "go-common/library/time"

const (
	_ byte = iota
	//ArtView 阅读
	ArtView
	//ArtReply 评论
	ArtReply
	//ArtShare 分享
	ArtShare
	//ArtCoin 硬币
	ArtCoin
	//ArtFavTBL 收藏
	ArtFavTBL
	//ArtLikeTBL 喜欢
	ArtLikeTBL
)

var (
	artTypeMap = map[byte]struct{}{
		ArtView:    {},
		ArtReply:   {},
		ArtShare:   {},
		ArtCoin:    {},
		ArtFavTBL:  {},
		ArtLikeTBL: {},
	}
)

//CheckType check article data type.
func CheckType(ty byte) bool {
	_, ok := artTypeMap[ty]
	return ok
}

// ArtTrend for article trend.
type ArtTrend struct {
	DateKey   int64 `json:"date_key"`
	TotalIncr int64 `json:"total_inc"`
}

// ArtRankMap for article rank source.
type ArtRankMap struct {
	AIDs  map[int]int64
	Incrs map[int]int
}

// ArtRankList for article top 10 list.
type ArtRankList struct {
	Arts []*ArtMeta `json:"art_rank"`
}

// ArtMeta for article rank meta data.
type ArtMeta struct {
	AID   int64      `json:"aid"`
	Incr  int        `json:"incr"`
	Title string     `json:"title"`
	PTime xtime.Time `json:"ptime"`
}

// ArtRead for article read source.
type ArtRead struct {
	Source map[string]int `family:"f" json:"source"`
}
