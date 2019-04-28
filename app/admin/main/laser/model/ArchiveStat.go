package model

import (
	xtime "go-common/library/time"
)

// ArchiveStat is table archive_stat
type ArchiveStat struct {
	ID        int64      `json:"id"`
	Business  int        `json:"business"`
	StatType  int        `json:"stat_type"`
	TypeID    int        `json:"typeid"`
	GroupID   int        `json:"group_id"`
	UID       int64      `json:"uid"`
	StatDate  xtime.Time `json:"stat_date"`
	Content   string     `json:"content"`
	Ctime     xtime.Time `json:"ctime"`
	Mtime     xtime.Time `json:"mtime"`
	State     int        `json:"state"`
	StatValue int64      `json:"stat_value"`
}

// StatNode is Minimum dimension stat value.
type StatNode struct {
	StatDate  xtime.Time `json:"stat_date"`
	Business  int        `json:"business"`
	StatType  int        `json:"stat_type"`
	TypeID    int        `json:"typeid"`
	UID       int64      `json:"uid"`
	StatValue int64      `json:"stat_value"`
}

// CsvMetaNode is
type CsvMetaNode struct {
	Index    int
	Name     string
	DataCode int
}

// StatItem is element of stat view json model.
type StatItem struct {
	DataCode int   `json:"data_code"`
	Value    int64 `json:"value"`
}

// StatView is common stat view json model.
type StatView struct {
	Date  int64       `json:"date"`
	Stats []*StatItem `json:"stats"`
}

// StatItemExt is StatItem extension.
type StatItemExt struct {
	Uname string      `json:"uname"`
	Stats []*StatItem `json:"stat"`
}

// StatViewExt is StatView extension
type StatViewExt struct {
	Date  int64          `json:"date"`
	Wraps []*StatItemExt `json:"stats"`
}

const (
	// business字段枚举值

	// ArchiveRecheck is 稿件回查项目
	ArchiveRecheck = 1
	// TagRecheck is 稿件频道tag回查项目
	TagRecheck = 2
	// RandomVideoAudit is 视频非定时审核操作数据
	RandomVideoAudit = 3
	// FixedVideoAudit is 视频定时审核操作数据
	FixedVideoAudit = 4

	// stat_type字段枚举值

	// 统计指标枚举值

	// TotalArchive is 总稿件量
	TotalArchive = 1
	// TotalOper is 总操作量
	TotalOper = 2
	// ReCheck is 打回量
	ReCheck = 3
	// Lock is 锁定量
	Lock = 4
	// ThreeLimit is 三限量
	ThreeLimit = 5
	// FirstCheck is 一查稿件量
	FirstCheck = 6
	// SecondCheck is 二查稿件量
	SecondCheck = 7
	// ThirdCheck is 三查稿件量
	ThirdCheck = 8
	// TotalOperFrequency is 总操作次数
	TotalOperFrequency = 9
	// FirstCheckOper is 一查次数
	FirstCheckOper = 10
	// SecondCheckOper is 二查次数
	SecondCheckOper = 11
	// ThirdCheckOper is 三查次数
	ThirdCheckOper = 12
	// FirstCheckTime is 一查响应总时间
	FirstCheckTime = 13
	// SecondCheckTime is 二查响应总时间
	SecondCheckTime = 14
	// ThirdCheckTime is 三查响应总时间
	ThirdCheckTime = 15
	// FirstAvgTime is 一查响应平均耗时
	FirstAvgTime = 16
	// SecondAvgTime is 二查响应平均耗时
	SecondAvgTime = 17
	// ThirdAvgTime is 三查响应平均耗时
	ThirdAvgTime = 18
	// NoRankArchive is 排行禁止
	NoRankArchive = 19
	// NoIndexArchive is 动态禁止
	NoIndexArchive = 20
	// NoRecommendArchive is 推荐禁止
	NoRecommendArchive = 21
	// NoPushArchive is 粉丝动态禁止
	NoPushArchive = 22
	// TagRecheckTotalTime is tag回查总时间
	TagRecheckTotalTime = 23
	// TagRecheckTotalCount is 频道回查操作总量
	TagRecheckTotalCount = 24
	// TagChangeCount is tag变更的稿件量
	TagChangeCount = 25
	// TagRecheckAvgTime is tag保存操作平均耗时
	TagRecheckAvgTime = 26
	// TotalVideo is 总操视频量
	TotalVideo = 27
	// TotalVideoOper is 总操作次数
	TotalVideoOper = 28
	// OpenVideo is 开放浏视频量
	OpenVideo = 29
	// OpenVideoOper is 开放浏览操作次数
	OpenVideoOper = 30
	// VipAccessVideo is 会员可视频量
	VipAccessVideo = 31
	// VipAccessVideoOper is 会员可见操作次数
	VipAccessVideoOper = 32
	// RejectVideo is 打视频量
	RejectVideo = 33
	// RejectVideoOper is 打回操作次数
	RejectVideoOper = 34
	// LockVideo is 锁视频量
	LockVideo = 35
	// LockVideoOper is 锁定操作次数
	LockVideoOper = 36
	// PassVideoTotalDuration is 通过视频总时长
	PassVideoTotalDuration = 37
	// FailVideoTotalDuration is 未通过视频总时长
	FailVideoTotalDuration = 38
	// WaitAuditAvgTime is 视频提交到进入待审平均响应时间
	WaitAuditAvgTime = 39
	// WaitAuditDuration is 视频提交到进入待审时间
	WaitAuditDuration = 40
	// WaitAuditOper is 视频提交到进入待审次数
	WaitAuditOper = 41

	//valueType

	// NumValue is num unit
	NumValue = 1
	// TimeValue is second unit
	TimeValue = 2
)
