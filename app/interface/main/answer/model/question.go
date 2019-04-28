package model

import (
	"time"
)

// state
const (
	WaitCheck   = int8(0)
	PassCheck   = int8(1)
	NoPassCheck = int8(2)
)

// size
const (
	MaxQuestion    = 120
	MinQuestion    = 6
	MaxAns         = 100
	MinAns         = 2
	MaxTips        = 100
	MinTips        = 2
	MaxLoadQueSize = 100000
)

// media type
const (
	TextMediaType  = int8(1)
	ImageMediaType = int8(2)

	TopRankSize    = 10
	MyQuestionSize = 12

	AttrUnknown = -1
)

// QuestionRPC stion question info.
type QuestionRPC struct {
	ID        int64     `json:"id"`
	Mid       int64     `json:"mid"`
	IP        uint32    `json:"ip"`
	TypeID    int8      `json:"type"`
	MediaType int8      `json:"media_type"`
	Check     int8      `json:"check"`
	Source    int8      `json:"source"`
	Question  string    `json:"question"`
	Ans1      string    `json:"ans1"`
	Ans2      string    `json:"ans2"`
	Ans3      string    `json:"ans3"`
	Ans4      string    `json:"ans4"`
	Tips      string    `json:"tips"`
	AvID      int32     `json:"av_id"`
	Ctime     time.Time `json:"ctime"`
}

// Question question info rpc.
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
}

// ExtraQst etc.
type ExtraQst struct {
	ID       int64     `json:"id"`
	Question string    `json:"question"`
	Ans      int8      `json:"ans"`
	Status   int8      `json:"status"`
	OriginID int64     `json:"origin_id"`
	AvID     int64     `json:"av_id"`
	Source   int8      `json:"source"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// ImgPosition .
type ImgPosition struct {
	Ans0H float64 `json:"ans_0_h"`
	Ans0Y float64 `json:"ans_0_y"`
	Ans1H float64 `json:"ans_1_h"`
	Ans1Y float64 `json:"ans_1_y"`
	Ans2H float64 `json:"ans_2_h"`
	Ans2Y float64 `json:"ans_2_y"`
	Ans3H float64 `json:"ans_3_h"`
	Ans3Y float64 `json:"ans_3_y"`
	Ans4H float64 `json:"ans_4_h"`
	Ans4Y float64 `json:"ans_4_y"`
	QsH   float64 `json:"qs_h"`
	QsY   float64 `json:"qs_y"`
}

// ExtraBigData ret.
type ExtraBigData struct {
	Done []int64 `json:"done"`
	Pend []int64 `json:"pend"`
}
