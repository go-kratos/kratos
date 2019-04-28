package model

import (
	"time"
)

// answer constants
const (
	LangZhCN = "zh-CN"
	LangZhTW = "zh-TW"
	LangZhHK = "zh-HK"
)

// answer constants
const (
	UserInfoRank = 5000
	PenDantDays  = 7 //答题优秀设置挂件的天数
	ExtraAnsA    = "符合规范"
	ExtraAnsB    = "不符合规范"
)

// Score info.
const (
	FullScore = 100
	Score85   = 85
	Score60   = 60
	Score0    = 0
)

// Rank info.
const (
	RankTop int = 122
)

// question type.
const (
	Q                int8 = iota
	BaseExtraNoPassQ      // 1 extra no pass
	BaseExtraPassQ        // 2 extra pass
)

// extra question ans.
const (
	UnKownQ int8 = iota
	ViolationQ
	NormalQ
)

// answer captcha pass
const (
	CaptchaNopass int8 = iota
	CaptchaPass
)

// BaseQues question record
type BaseQues struct {
	Question string
	Check    int8
	Ctime    time.Time
}

// MyQues my question
type MyQues struct {
	Count int64
	List  []*BaseQues
}

// RankInfo rank
type RankInfo struct {
	Mid       int64          `json:"mid"`
	Face      string         `json:"face"`
	Uname     string         `json:"uname"`
	Num       int64          `json:"num"`
	Nameplate *NameplateInfo `json:"nameplate"`
}

// NameplateInfo .
type NameplateInfo struct {
	Nid        int    `json:"nid"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ImageSmall string `json:"image_small"`
	Level      string `json:"level"`
	Condition  string `json:"condition"`
}

// TypeInfo type info
type TypeInfo struct {
	ID, Parentid int64
	Name         string
	LabelName    string
	Subs         []*SubType
}

// ProTypes .
type ProTypes struct {
	List                 []*TypeInfo
	CurrentTime, EndTime time.Time
	Repro                bool
}

// SubType sub type info
type SubType struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	LabelName string `json:"-"`
}

// AnsQue .
type AnsQue struct {
	ID        int64
	Img       string
	Height    float64
	PositionY float64 // background-position-y
	Ans       []*AnsPosition
}

// AnsPosition .
type AnsPosition struct {
	AnsHash   string
	Height    float64 // height
	PositionY float64 // background-position-y
}

// AnsQuesList .
type AnsQuesList struct {
	CurrentTime, EndTime time.Time
	QuesList             []*AnsQue
}

// AnsCheck .
type AnsCheck struct {
	QidList   []int64
	HistoryID int64
	Pass      bool
}

// CaptchaReq Captcha request.
type CaptchaReq struct {
	Mid        int64
	IP         string
	ClientType string
	NewCaptcha int
}

// CaptchaCheckReq Captcha check request.
type CaptchaCheckReq struct {
	Mid        int64
	IP         string
	Challenge  string
	ClientType string
	Validate   string
	Seccode    string
	Success    int
	Cookie     string
	Comargs    map[string]string
}

// QueReq request
type QueReq struct {
	ID int64
}

// AnsCool .
type AnsCool struct {
	Hid            int64        `json:"hid,omitempty"`
	URL            string       `json:"url,omitempty"`
	Name           string       `json:"uname"`
	Face           string       `json:"face"`
	Powers         []*CoolPower `json:"power_result"`
	Score          int8         `json:"score"`
	Rank           *CoolRank    `json:"rank"`
	Share          *CoolShare   `json:"share"`
	CanShowRankBtn bool         `json:"can_show_rank_btn"`
	IsSameUser     bool         `json:"is_same_user"`
	ViewMore       string       `json:"view_more"`
	VideoInfo      *CoolVideo   `json:"video_info"`
	IsFirstPass    int8         `json:"is_first_pass"`
	Level          int8         `json:"level"`
	MainTids       []int        `json:"main_tids"`
	SubTids        []int        `json:"sub_tids"`
}

// CoolPower .
type CoolPower struct {
	Num  int64  `json:"num"`
	Name string `json:"name"`
}

// CoolRank .
type CoolRank struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

// CoolShare .
type CoolShare struct {
	Content      string `json:"content"`
	ShortContent string `json:"short_content"`
}

// CoolVideo .
type CoolVideo struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	Img      string `json:"img"`
	WatchNum string `json:"watch_num"`
	UpNum    string `json:"up_num"`
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

// AnswerTime info.
type AnswerTime struct {
	Stime  time.Time `json:"stime"`  // answer start time
	Etimes int8      `json:"etimes"` // base answer error times
}

// AnsHash .
type AnsHash struct {
	AnsHashName string
	AnsHashVal  string
}

// ExtraScoreReply .
type ExtraScoreReply struct {
	Score int64 `json:"score"`
}
