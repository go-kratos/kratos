package model

import xtime "go-common/library/time"

// LabourQs labour question.
type LabourQs struct {
	ID       int64      `json:"id"`
	Question string     `json:"question"`
	Ans      int64      `json:"-"`
	TrueAns  int64      `json:"-"` // 真实答案 答案0:未知 1:违规 2:不违规
	AvID     int64      `json:"av_id"`
	AvTitle  string     `json:"av_title"`
	Status   int64      `json:"-"`
	Source   int64      `json:"-"`
	Ctime    xtime.Time `json:"-"`
	Mtime    xtime.Time `json:"-"`
}

// LabourAns labour answer.
type LabourAns struct {
	ID     []int64
	Answer []int64
}

//AIQsID AI give question id.
type AIQsID struct {
	Pend []int64 `json:"pend"` // 未审核
	Done []int64 `json:"done"` // 已审核
}

// DataBusResult databus结果
type DataBusResult struct {
	Mid   int64  `json:"mid"`   // 用户 ID
	Buvid string `json:"buvid"` // 设备标识符 前端传入
	IP    string `json:"ip"`    // 用户 IP 地址
	Ua    string `json:"ua"`    // 客户端 User Agent
	Refer string `json:"refer"` // 页面跳转来源地址 Refer
	Score int64  `json:"score"` // 答题总分数
	Rs    []Rs
}

// Rs struct
type Rs struct {
	ID       int64      `json:"id"`       // 题目自增 ID
	Question string     `json:"question"` // 问题内容
	Ans      int64      `json:"ans"`      // 用户答案
	TrueAns  int64      `json:"trueAns"`  // 真实答案 答案0:未知 1:违规 2:不违规
	AvID     int64      `json:"av_id"`    // 相关视频id
	Status   int64      `json:"status"`   // 问题状态 1:未申核 2:已审核
	Source   int64      `json:"source"`   // 问题来源 0:未知1:评论 2:弹幕
	Ctime    xtime.Time `json:"ctime"`    // 创建时间
	Mtime    xtime.Time `json:"mtime"`    // 修改时间
}

// BlockAndMoralStatus blocked status and moral.
type BlockAndMoralStatus struct {
	MID    int64      `json:"mid"`
	Status int8       `json:"status"`
	STime  xtime.Time `json:"start_time"`
	ETime  xtime.Time `json:"end_time"`
}

// CommitRs struct
type CommitRs struct {
	Score int64 `json:"score"`
	Day   int64 `json:"day"`
}

// QsCache struct
type QsCache struct {
	Stime xtime.Time
	QsStr string
}
