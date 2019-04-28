package model

import (
	"fmt"
	"time"
)

const (
	// TotalScore of up rating
	TotalScore = 600
	// LowerBoundScore ...
	LowerBoundScore = 240
)

// Rating of up
type Rating struct {
	Score      *Score       `json:"score"`
	Rank       *Rank        `json:"rank"`
	Prize      *Prize       `json:"prize"`
	Privileges []*Privilege `json:"privileges"`
}

// Score of rating
type Score struct {
	MID       int64     `json:"mid"`
	Magnetic  int       `json:"magnetic"`   // 电磁力
	Creative  int       `json:"creative"`   // 创作力
	Influence int       `json:"influence"`  // 影响力
	Credit    int       `json:"credit"`     // 信用分
	CDate     time.Time `json:"c_date"`     // 统计月
	StatStart time.Time `json:"stat_start"` // 统计周期开始日
	StatEnd   time.Time `json:"stat_end"`   // 统计周期结束日
}

// Rank of rating
type Rank struct {
	Level RankLevel `json:"level"`
	Desc  string    `json:"desc"`
}

// RankLevel of rank
type RankLevel int8

// RankLevel list
const (
	RankLevelSuper RankLevel = 10 * (1 + iota)
	RankLevelStrong
	RankLevelStandout
	RankLevelNormal
	RankLevelNone
)

// Ranks list all levels of rank
var Ranks = []RankLevel{
	RankLevelSuper,
	RankLevelStrong,
	RankLevelStandout,
	RankLevelNormal,
	RankLevelNone,
}

// rank meta info
var rankMeta = map[RankLevel]struct {
	score int
	desc  string
}{
	RankLevelSuper:    {int(0.9 * TotalScore), "超能力"},
	RankLevelStrong:   {int(0.75 * TotalScore), "强能力"},
	RankLevelStandout: {int(0.6 * TotalScore), "异能力"},
	RankLevelNormal:   {int(0.3 * TotalScore), "常能力"},
	RankLevelNone:     {0, "新能力"},
}

// Score of rankLevel
func (r RankLevel) Score() int {
	if m, ok := rankMeta[r]; ok {
		return m.score
	}
	return RankLevelNone.Score()
}

// Rank content of rankLevel
func (r RankLevel) Rank() *Rank {
	if m, ok := rankMeta[r]; ok {
		return &Rank{
			Level: r,
			Desc:  m.desc,
		}
	}
	return RankLevelNone.Rank()
}

// Prize of rating
type Prize struct {
	Level   PrizeLevel `json:"level"`
	Desc    string     `json:"desc"`
	Content string     `json:"content"`
}

// PrizeLevel of prize
type PrizeLevel int8

// Prize Level List
const (
	PrizeLevelOne PrizeLevel = 10 * (1 + iota)
	PrizeLevelTwo
	PrizeLevelThree
	PrizeLevelFour
	PrizeLevelFive
)

// Prizes list prize levels by priority
var Prizes = []PrizeLevel{
	PrizeLevelOne,
	PrizeLevelTwo,
	PrizeLevelThree,
	PrizeLevelFour,
	PrizeLevelFive,
}

var prizeMeta = map[PrizeLevel]struct {
	desc    string
	content func(arg ...interface{}) string
}{
	PrizeLevelOne: {desc: "睥睨众生奖", content: func(...interface{}) string {
		return "恭喜你获得超高的电磁力，那可真是会当临绝顶，一览众山小吖"
	}},
	PrizeLevelTwo: {desc: "稳如泰山奖", content: func(...interface{}) string {
		return "稳如泰山是你的优点，也可能是你的天花板，试着努力突破一下吧"
	}},
	PrizeLevelThree: {desc: "飞速进步奖", content: func(arg ...interface{}) string {
		return fmt.Sprintf("本月电磁力上升%d分，真是付出了超级多努力呢，请继续加油吧", arg[0])
	}},
	PrizeLevelFour: {desc: "特别有趣奖", content: func(...interface{}) string {
		return "看来你是被2233娘选中的孩子，希望这样的幸运能够继续支撑你努力"
	}},
	PrizeLevelFive: {desc: "全村希望奖", content: func(...interface{}) string {
		return "作为全村的希望，未来的你一定会感谢现在持续努力的自己"
	}},
}

// Prize constructor
func (p PrizeLevel) Prize(arg ...interface{}) *Prize {
	if meta, ok := prizeMeta[p]; ok {
		return &Prize{
			Level:   p,
			Desc:    meta.desc,
			Content: meta.content(arg...),
		}
	}
	return PrizeLevelFive.Prize(arg...)
}

// Privilege of rating
type Privilege struct{}
