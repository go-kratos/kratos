package model

import (
	"go-common/library/ecode"
	"sort"
)

// const var
const (
	WilsonLHRRAlgorithm      = "wilsonLHRR"
	WilsonLHRRFluidAlgorithm = "wilsonLHRRFluid"
	OriginAlgorithm          = "origin"
	LikeDescAlgorithm        = "likeDesc"

	StateInactive = int(0)
	StateActive   = int(1)

	SlotsNum = 100

	DefaultSlotName  = "default"
	DefaultAlgorithm = "default"
	DefaultWeight    = ""
)

// EventMsg EventMsg
type EventMsg struct {
	Action string `json:"action"`
	Oid    int64  `json:"oid"`
	Tp     int    `json:"tp"`
}

// SlotsMapping slot name mapping
type SlotsMapping struct {
	Name  string
	Slots []int
	State int
}

// SlotsStat slots stat
type SlotsStat struct {
	Name      string
	Slots     []int
	Algorithm string
	Weight    string
	State     int
}

// StatisticsStats StatisticsStats
type StatisticsStats []*StatisticsStat

// GroupByName group statistics by name
func (s StatisticsStats) GroupByName() (res map[string]StatisticsStats) {
	res = make(map[string]StatisticsStats)
	for _, stat := range s {
		if _, ok := res[stat.Name]; ok {
			res[stat.Name] = append(res[stat.Name], stat)
		} else {
			var tmp []*StatisticsStat
			tmp = append(tmp, stat)
			res[stat.Name] = tmp
		}
	}
	return
}

// StatisticsStat 实验组或者对照组的各项统计
type StatisticsStat struct {
	// 流量所属槽位 0~99
	Slot int
	// 所属实验组名
	Name          string
	State         int
	Date          int
	Hour          int
	HotLike       int64
	HotHate       int64
	HotReport     int64
	HotChildReply int64
	// 整个评论区
	TotalLike       int64
	TotalHate       int64
	TotalReport     int64
	TotalChildReply int64
	TotalRootReply  int64
	// 用户点开评论区次数
	View uint32
	// 评论列表接口调用次数
	TotalView uint32
	// 热门评论接口调用次数
	HotView uint32
	// 更多热门评论点击次数
	HotClick uint32
	// 用户在评论首页看到的热门评论被点赞点踩评论以及举报的次数

	// UV的统计数据
	HotLikeUV   int64
	HotHateUV   int64
	HotReportUV int64
	HotChildUV  int64

	TotalLikeUV   int64
	TotalHateUV   int64
	TotalReportUV int64
	TotalChildUV  int64
	TotalRootUV   int64
}

// Merge Merge
func (stat1 *StatisticsStat) Merge(stat2 *StatisticsStat) (stat3 *StatisticsStat) {
	stat3 = new(StatisticsStat)
	stat3.View = stat1.View + stat2.View
	stat3.HotView = stat1.HotView + stat2.HotView
	stat3.HotClick = stat1.HotClick + stat2.HotClick
	stat3.TotalView = stat1.TotalView + stat2.TotalView
	return
}

// DivideByPercent ...
func (stat1 *StatisticsStat) DivideByPercent(percent int64) (stat2 *StatisticsStat) {
	stat2 = new(StatisticsStat)
	if percent <= 0 {
		return
	}
	stat2.Name = stat1.Name
	stat2.Date = stat1.Date
	stat2.Hour = stat1.Hour
	stat2.View = stat1.View / uint32(percent)
	stat2.HotView = stat1.HotView / uint32(percent)
	stat2.HotClick = stat1.HotClick / uint32(percent)
	stat2.TotalView = stat1.TotalView / uint32(percent)
	stat2.HotLike = stat1.HotLike / percent
	stat2.HotHate = stat1.HotHate / percent
	stat2.HotChildReply = stat1.HotChildReply / percent
	stat2.HotReport = stat1.HotReport / percent
	stat2.TotalLike = stat1.TotalLike / percent
	stat2.TotalHate = stat1.TotalHate / percent
	stat2.TotalReport = stat1.TotalReport / percent
	stat2.TotalRootReply = stat1.TotalRootReply / percent
	stat2.TotalChildReply = stat1.TotalChildReply / percent
	return
}

// MergeByDate MergeByDate
func (stat1 *StatisticsStat) MergeByDate(stat2 *StatisticsStat) (stat3 *StatisticsStat) {
	stat3 = new(StatisticsStat)
	stat3.Name = stat1.Name
	stat3.Date = stat1.Date
	stat3.View = stat1.View + stat2.View
	stat3.HotView = stat1.HotView + stat2.HotView
	stat3.HotClick = stat1.HotClick + stat2.HotClick
	stat3.TotalView = stat1.TotalView + stat2.TotalView
	stat3.HotLike = stat1.HotLike + stat2.HotLike
	stat3.HotHate = stat1.HotHate + stat2.HotHate
	stat3.HotChildReply = stat1.HotChildReply + stat2.HotChildReply
	stat3.HotReport = stat1.HotReport + stat2.HotReport
	stat3.TotalLike = stat1.TotalLike + stat2.TotalLike
	stat3.TotalHate = stat1.TotalHate + stat2.TotalHate
	stat3.TotalReport = stat1.TotalReport + stat2.TotalReport
	stat3.TotalRootReply = stat1.TotalRootReply + stat2.TotalRootReply
	stat3.TotalChildReply = stat1.TotalChildReply + stat2.TotalChildReply
	return
}

// WilsonLHRRWeight wilson score interval weight
type WilsonLHRRWeight struct {
	Like   float64 `json:"like"`
	Hate   float64 `json:"hate"`
	Reply  float64 `json:"reply"`
	Report float64 `json:"report"`
}

// Validate Validate
func (weight WilsonLHRRWeight) Validate() (err error) {
	if weight.Report*weight.Reply*weight.Hate*weight.Like <= 0 {
		err = ecode.RequestErr
		return
	}
	return
}

// WilsonLHRRFluidWeight WilsonLHRRFluidWeight
type WilsonLHRRFluidWeight struct {
	Like   float64 `json:"like"`
	Hate   float64 `json:"hate"`
	Reply  float64 `json:"reply"`
	Report float64 `json:"report"`
	Slope  float64 `json:"slope"`
}

// Validate Validate
func (weight WilsonLHRRFluidWeight) Validate() (err error) {
	if weight.Report*weight.Reply*weight.Hate*weight.Like*weight.Slope <= 0 {
		err = ecode.RequestErr
		return
	}
	return
}

// SSReq ss req
type SSReq struct {
	DateFrom int64 `form:"date_from" validate:"required"`
	DateEnd  int64 `form:"date_end" validate:"required"`
	Hour     bool  `form:"hour"`
}

// SSHourRes ss res
type SSHourRes struct {
	Legend []string                     `json:"legend"`
	XAxis  []string                     `json:"x_axis"`
	Series map[string][]*StatisticsStat `json:"series"`
}

// Sort ...
func (s *SSHourRes) Sort() {
	sort.Strings(s.Legend)
	sort.Strings(s.XAxis)
	for _, v := range s.Series {
		sort.Slice(v, func(i, j int) bool { return v[i].Date*100+v[i].Hour < v[j].Date*100+v[j].Hour })
	}
}

// SSDateRes ss res
type SSDateRes struct {
	Legend []string                     `json:"legend"`
	XAxis  []int                        `json:"x_axis"`
	Series map[string][]*StatisticsStat `json:"series"`
}

// Sort ...
func (s *SSDateRes) Sort() {
	sort.Strings(s.Legend)
	sort.Ints(s.XAxis)
	for _, v := range s.Series {
		sort.Slice(v, func(i, j int) bool { return v[i].Date < v[j].Date })
	}
}
