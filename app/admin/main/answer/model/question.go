package model

import "time"

// check info.
const (
	WaitCheck   = int8(0)
	PassCheck   = int8(1)
	NoPassCheck = int8(2)
)

// count.
const (
	ArgsCount   = 6
	MaxCount    = 1000
	FileMaxSize = 2 * (1024 * 1024) //  FileMaxSize max 2M
)

// Medal info
const (
	PassNum50  = 50
	PassNum100 = 100
	PassNum200 = 200

	Nid53 = 53
	Nid54 = 54
	Nid55 = 55

	StageDisable int8 = 2
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
)

//QuestionPage admin page
type QuestionPage struct {
	Total int64         `json:"total"`
	Items []*QuestionDB `json:"items"`
}

//HistoryPage .
type HistoryPage struct {
	Total int64              `json:"total"`
	Items []*AnswerHistoryDB `json:"items"`
}

// QuestionDB question info.
type QuestionDB struct {
	ID        int64     `gorm:"column:id" json:"id" form:"id" validate:"required"`
	Mid       int64     `gorm:"column:mid" json:"mid"`
	IP        string    `gorm:"column:ip" json:"ip"`
	TypeID    int8      `gorm:"column:type_id" json:"type_id"`
	Question  string    `gorm:"column:question" json:"question" form:"question" validate:"required"`
	Ans1      string    `gorm:"column:ans1" json:"ans1" form:"ans1" validate:"required"`
	Ans2      string    `gorm:"column:ans2" json:"ans2" form:"ans2" validate:"required"`
	Ans3      string    `gorm:"column:ans3" json:"ans3" form:"ans3" validate:"required"`
	Ans4      string    `gorm:"column:ans4" json:"ans4" form:"ans4" validate:"required"`
	State     int8      `gorm:"column:state" json:"state"`
	Tips      string    `gorm:"column:tips" json:"tips"`
	AvID      int32     `gorm:"column:avid" json:"avid"`
	MediaType int8      `gorm:"column:media_type" json:"media_type"`
	Source    int8      `gorm:"column:source" json:"source"`
	Ctime     time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime     time.Time `gorm:"column:mtime" json:"mtime"`
	Operator  string    `gorm:"column:operator" json:"operator"`
}

// TableName for gorm.
func (b QuestionDB) TableName() string {
	return "ans_v3_question"
}

// Question question info.
type Question struct {
	*QuestionDB
	Ans []string
}

// ArgQue admin question query param.
type ArgQue struct {
	Question string `form:"question"`
	TypeID   int8   `form:"type_id"`
	State    int8   `form:"state" default:"-1"`
	Ps       int    `form:"ps" default:"20"`
	Pn       int    `form:"pn" default:"1"`
}

// ArgHistory .
type ArgHistory struct {
	Mid int64 `form:"mid" validate:"required"`
	Ps  int   `form:"ps" default:"20"`
	Pn  int   `form:"pn" default:"1"`
}

// Sizer .
type Sizer interface {
	Size() int64
}

// AnswerHistory info.
type AnswerHistory struct {
	ID                    int64     `json:"id"`
	Hid                   int64     `json:"hid"`
	Mid                   int64     `json:"mid"`
	StartTime             time.Time `json:"start_time"`
	StepOneErrTimes       int8      `json:"step_one_err_times"`
	StepOneCompleteTime   int64     `json:"step_one_complete_time"`
	StepExtraStartTime    time.Time `json:"step_extra_start_time"`
	StepExtraCompleteTime int64     `json:"step_extra_complete_time"`
	StepExtraScore        int64     `json:"step_extra_score"`
	StepTwoStartTime      time.Time `json:"step_two_start_time"`
	CompleteTime          time.Time `json:"complete_time"`
	CompleteResult        string    `json:"complete_result"`
	Score                 int8      `json:"score"`
	IsFirstPass           int8      `json:"is_first_pass"`
	IsPassCaptcha         int8      `json:"is_pass_captcha"`
	PassedLevel           int8      `json:"passed_level"`
	RankID                int       `json:"rank_id"`
	Ctime                 time.Time `json:"ctime"`
	Mtime                 time.Time `json:"mtime"`
}

//List .
type List struct {
	Total int              `json:"total"`
	Items []*AnswerHistory `json:"items"`
}

// Histories history sorted.
type Histories []*AnswerHistory

func (h Histories) Len() int { return len(h) }
func (h Histories) Less(i, j int) bool {
	return h[i].Ctime.Unix() > h[j].Ctime.Unix()
}
func (h Histories) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
