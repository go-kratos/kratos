package model

import (
	xtime "go-common/library/time"
)

// const varialble
const (
	DefaultAlgorithm         = "default"
	WilsonLHRRAlgorithm      = "wilsonLHRR"
	WilsonLHRRFluidAlgorithm = "wilsonLHRRFluid"
	OriginAlgorithm          = "origin"
	LikeDescAlgorithm        = "likeDesc"
	StateInactive            = int(0)
	StateActive              = int(1)
	StateDelete              = int(2)

	// 用于redis中统计uv
	StatisticActionRootReply  = "rr"
	StatisticActionChildReply = "cr"
	StatisticActionLike       = "l"
	StatisticActionHate       = "h"
	StatisticActionReport     = "r"

	StatisticKindTotal = "t"
	StatisticKindHot   = "h"

	DatabusActionReply      = "reply"
	DatabusActionReport     = "report_add"
	DatabusActionLike       = "like"
	DatabusActionCancelLike = "like_cancel"
	DatabusActionHate       = "hate"
	DatabusActionCancelHate = "hate_cancel"
	// user upper or admin delete
	DatabusActionDel = "reply_del"
	// admin delete by report
	DatabusActionRptDel = "report_del"
	// admin recover
	DatabusActionRecover = "reply_recover"
	// admin or upper top reply
	DatabusActionTop = "top"
	// admin or upper untop reply
	DatabusActionUnTop = "untop"

	DatabusActionReIdx = "re_idx"

	// 只有大于等于3个赞且评论区根评论数目多余20才会被加入热门评论列表
	MinLikeCount      = 3
	MinRootReplyCount = 20
	SlotsNum          = 100
)

// Statistics const
var (
	StatisticActions    = []string{StatisticActionRootReply, StatisticActionChildReply, StatisticActionLike, StatisticActionHate, StatisticActionReport}
	StatisticKinds      = []string{StatisticKindTotal, StatisticKindHot}
	StatisticsDatabaseI = []string{"`name`", "`date`", "`hour`"}
	StatisticsDatabaseU = []string{"hot_like", "hot_hate", "hot_report", "hot_child", "total_like", "total_hate", "total_report", "total_root", "total_child"}
	StatisticsDatabaseS = []string{"hot_like_uv", "hot_hate_uv", "hot_report_uv", "hot_child_uv", "total_like_uv", "total_hate_uv", "total_report_uv", "total_child_uv", "total_root_uv"}
)

// ReplyScore reply score
type ReplyScore struct {
	RpID  int64
	Score float64
}

// ReplyStat 放在MC里的衡量一条根评论质量的各个参数
type ReplyStat struct {
	RpID        int64      `json:"rpid"`
	Like        int        `json:"like"`
	Hate        int        `json:"hate"`
	Reply       int        `json:"reply"`
	Report      int        `json:"report"`
	SubjectTime xtime.Time `json:"subject_time"`
	ReplyTime   xtime.Time `json:"reply_time"`
}

// ReplyResp 返回给reply-interface的评论ID数组，已按热度排好序
type ReplyResp struct {
	RpIDs []int64
	// 属于哪一个实验组
	TestSetName string
}

// ReplyList 存在redis sorted set中的数据结构
type ReplyList struct {
	RpID []int64
}

// SlotStat slot stat
type SlotStat struct {
	Name      string
	Slot      int
	Algorithm string
	Weight    string
}

// SlotsStat SlotsStat
type SlotsStat struct {
	Name      string
	Slots     []int
	Algorithm string
	Weight    string
}

// SlotsMapping E group slots
type SlotsMapping struct {
	Name  string
	Slots []int
}

// StatisticsStat 实验组或者对照组的各项统计
type StatisticsStat struct {
	// 流量所属槽位 0~99
	Slot int
	// 所属实验组名
	Name string
	// 用户在评论首页看到的热门评论被点赞点踩评论以及举报的次数
	HotLike       int64
	HotHate       int64
	HotChildReply int64
	HotReport     int64
	// 整个评论区
	TotalLike       int64
	TotalHate       int64
	TotalReport     int64
	TotalRootReply  int64
	TotalChildReply int64

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

// Merge merge two statistics
func (stat1 *StatisticsStat) Merge(stat2 *StatisticsStat) (stat3 *StatisticsStat) {
	stat3 = new(StatisticsStat)
	stat3.TotalLike = stat1.TotalLike + stat2.TotalLike
	stat3.TotalHate = stat1.TotalHate + stat2.TotalHate
	stat3.TotalReport = stat1.TotalReport + stat2.TotalReport
	stat3.TotalRootReply = stat1.TotalRootReply + stat2.TotalRootReply
	stat3.TotalChildReply = stat1.TotalChildReply + stat2.TotalChildReply
	stat3.HotLike = stat1.HotLike + stat2.HotLike
	stat3.HotHate = stat1.HotHate + stat2.HotHate
	stat3.HotReport = stat1.HotReport + stat2.HotReport
	stat3.HotChildReply = stat1.HotChildReply + stat2.HotChildReply

	stat3.HotLikeUV = stat1.HotLikeUV + stat2.HotLikeUV
	stat3.HotHateUV = stat1.HotHateUV + stat2.HotHateUV
	stat3.HotReportUV = stat1.HotReportUV + stat2.HotReportUV
	stat3.HotChildUV = stat1.HotChildUV + stat2.HotChildUV

	stat3.TotalLikeUV = stat1.TotalLikeUV + stat2.TotalLikeUV
	stat3.TotalHateUV = stat1.TotalHateUV + stat2.TotalHateUV
	stat3.TotalReportUV = stat1.TotalReportUV + stat2.TotalReportUV
	stat3.TotalRootUV = stat1.TotalRootUV + stat2.TotalRootUV
	stat3.TotalChildUV = stat1.TotalChildUV + stat2.TotalChildUV
	return
}

// StrategyStat 实验组所使用算法，以及各个参数情况
type StrategyStat struct {
	Name      string             `json:"name"`
	Percent   int                `json:"percent"`
	Algorithm string             `json:"algorithm"`
	Args      map[string]float64 `json:"args"`
}

// RefreshChecker 刷新热门评论的触发条件，用来对同一个评论区的所有请求进行聚合
type RefreshChecker struct {
	Oid           int64
	Type          int
	LastTimeStamp int64
}

// WilsonLHRRWeight wilson score interval weight
type WilsonLHRRWeight struct {
	Like   float64
	Hate   float64
	Reply  float64
	Report float64
}

// WilsonLHRRFluidWeight wilson
type WilsonLHRRFluidWeight struct {
	Like   float64
	Hate   float64
	Reply  float64
	Report float64
	Slope  float64
}

// EventMsg event message
type EventMsg struct {
	Action string `json:"action"`
	Oid    int64  `json:"oid"`
	Tp     int    `json:"tp"`
}

// StatsMsg stats message
type StatsMsg struct {
	Action  string   `json:"action"`
	Mid     int64    `json:"mid"`
	Subject *Subject `json:"subject"`
	Reply   *Reply   `json:"reply"`
	Report  *Report  `json:"report,omitempty"`
}

// Sharding 返回该用户属于哪一个组
// 将流量划分为100份
func (r *StatsMsg) Sharding() int64 {
	return r.Mid % SlotsNum
}

// HotCondition return if should check exists in hot reply
func (r *StatsMsg) HotCondition() bool {
	if r.Action == DatabusActionReply && !r.Reply.IsRoot() {
		return true
	}
	if r.Reply.IsRoot() && r.Reply.Like >= MinLikeCount &&
		(r.Action == DatabusActionLike || r.Action == DatabusActionHate ||
			r.Action == DatabusActionCancelLike || r.Action == DatabusActionCancelHate || r.Action == DatabusActionReport) {
		return true
	}
	return false
}

// Reply define reply object
type Reply struct {
	RpID   int64      `json:"rpid"`
	Mid    int64      `json:"mid"`
	Root   int64      `json:"root"`
	Parent int64      `json:"parent"`
	RCount int        `json:"rcount"`
	Floor  int        `json:"floor"`
	State  int8       `json:"state"`
	Attr   uint32     `json:"attr"`
	CTime  xtime.Time `json:"ctime"`
	Like   int        `json:"like"`
	Hate   int        `json:"hate"`
}

// Legal return a reply legal
func (r *Reply) Legal() bool {
	// 0,1,2,5,6 所有需要显示给用户的评论state
	return r.State == 0 || r.State == 1 || r.State == 2 || r.State == 5 || r.State == 6
}

// ShowAfterAudit ShowAfterAudit
func (r *Reply) ShowAfterAudit() bool {
	return r.State == 11
}

// AuditButShow AuditButShow
func (r *Reply) AuditButShow() bool {
	return r.State == 5
}

// IsRoot IsRoot
func (r *Reply) IsRoot() bool {
	return r.Root == 0
}

// Qualified Qualified
func (r *Reply) Qualified() bool {
	return r.Like >= MinLikeCount
}

// Report define reply report
type Report struct {
	RpID  int64      `json:"rpid"`
	Mid   int64      `json:"mid"`
	Count int        `json:"count"`
	Score int        `json:"score"`
	State int8       `json:"state"`
	CTime xtime.Time `json:"ctime"`
	Attr  uint32     `json:"attr"`
}

// Subject is subject of reply
type Subject struct {
	Oid    int64      `json:"oid"`
	Type   int        `json:"type"`
	Mid    int64      `json:"mid"`
	RCount int        `json:"rcount"`
	State  int8       `json:"state"`
	Attr   uint32     `json:"attr"`
	CTime  xtime.Time `json:"ctime"`
}

// ShowHotReply if show
func (s *Subject) ShowHotReply() bool {
	return s.RCount >= MinRootReplyCount
}
