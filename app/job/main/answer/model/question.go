package model

import "time"

// LabourQs labour question.
type LabourQs struct {
	ID       int64  `json:"id"`
	Mid      int64  `json:"mid"`
	Question string `json:"question"`
	Ans      int8   `json:"ans"`
	AvID     int64  `json:"av_id"`
	Status   int8   `json:"status"`
	Source   int8   `json:"source"`
	Isdel    int8   `json:"isdel"`
	State    int8   `json:"state"`
	Ctime    string `json:"ctime"`
	Mtime    string `json:"mtime"`
}

// Question question info .
type Question struct {
	ID        int64     `json:"id"`
	Mid       int64     `json:"mid"`
	IP        string    `json:"ip"`
	TypeID    int8      `json:"type"`
	MediaType int8      `json:"media_type"`
	Check     int8      `json:"check"`
	Source    int8      `json:"source"`
	Question  string    `json:"question"`
	Ans1      string    `json:"ans1"`
	Ans2      string    `json:"ans2"`
	Ans3      string    `json:"ans3"`
	Ans4      string    `json:"ans4"`
	Ans       []string  `json:"-"`
	Tips      string    `json:"tips"`
	AvID      int32     `json:"av_id"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
	Operator  string    `json:"operator"`
}

// answer constants
const (
	ExtraAnsA    = "符合规范"
	ExtraAnsB    = "不符合规范"
	HadCreateImg = 1
	LimitSize    = 100
)

// Formal user formal info.
type Formal struct {
	Mid      int64     `json:"mid"`        // 用户 ID
	Hid      int64     `json:"history_id"` // 答题历史 ID
	Cookie   string    `json:"cookie"`     // cookie
	IP       string    `json:"ip"`         // cookie
	PassTime time.Time `json:"pass_time"`  // 通过时间
}
