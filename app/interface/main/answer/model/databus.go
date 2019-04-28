package model

import "time"

// const .
const (
	AnswerLogID  = 15
	AnswerUpdate = "answer_update"

	BasePass       = "basePass"
	ExtraStartTime = "extraStartTime"
	ExtraCheck     = "extraCheck"
	ProQues        = "proQues"
	ProCheck       = "proCheck"
	Captcha        = "captchaPass"
	Level          = "level"
)

// DataBusResult databus.
type DataBusResult struct {
	Mid   int64  `json:"mid"`   // 用户 ID
	Buvid string `json:"buvid"` // 设备标识符 前端传入
	IP    string `json:"ip"`    // 用户 IP 地址
	Ua    string `json:"ua"`    // 客户端 User Agent
	Refer string `json:"refer"` // 页面跳转来源地址 Refer
	Score int8   `json:"score"` // 答题总分数
	Hid   int64  `json:"hid"`   // hid
	Rs    []*Rs
}

// Rs def.
type Rs struct {
	ID       int64     `json:"id"`       // 题目自增 ID
	Question string    `json:"question"` // 问题内容
	Ans      int8      `json:"ans"`      // 用户答案
	TrueAns  int8      `json:"trueAns"`  // 真实答案 答案0:未知 1:违规 2:不违规
	AvID     int64     `json:"av_id"`    // 相关视频id
	Status   int8      `json:"status"`   // 问题状态 1:未申核 2:已审核
	Source   int8      `json:"source"`   // 问题来源 0:未知1:评论 2:弹幕
	Ctime    time.Time `json:"ctime"`    // 创建时间
	Mtime    time.Time `json:"mtime"`    // 修改时间
}

// Formal user formal info.
type Formal struct {
	Mid      int64     `json:"mid"`        // 用户 ID
	Hid      int64     `json:"history_id"` // 答题历史 ID
	Cookie   string    `json:"cookie"`     // cookie
	IP       string    `json:"ip"`         // cookie
	PassTime time.Time `json:"pass_time"`  // 通过时间
}
