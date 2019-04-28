package model

import (
	"database/sql"
	"go-common/library/time"
)

// 常量
const (
	DeledStatus    = 1  // 已删除
	PAGESIZE       = 20 // 每页条数
	STARTINDEX     = 20 // 开始页码
	MULTIPLECHOICE = 2  //多选
)

// Question 题目返回数据
type Question struct {
	QsID       int64  `json:"qid"`
	QsType     int8   `json:"question_type"`
	AnswerType int8   `json:"answer_type"`
	QsName     string `json:"question_name"`
	QsDif      int8   `json:"difficulty"`
	QsBId      int64  `json:"qb_id"`
	IsDeleted  uint8  `json:"is_deleted"`
	//Ctime         string `json:"ctime"`
	//Mtime         string `json:"mtime"`
}

// GetQuestionItem 获取题目接口返回数据
type GetQuestionItem struct {
	*Question
	Answers    []*Answer   `json:"answers"`
	QuestBkPic *QuestBkPic `json:"qspic"`
	AllCnt     int64       `json:"total"`
	AnTime     int64       `json:"answer_cnt"`
}

// QuestBkPic 坐标
type QuestBkPic struct {
	X   int    `json:"x"`
	Y   int    `json:"y"`
	Src string `json:"src"`
}

// QuestionAll 题目所有
type QuestionAll struct {
	Question
	AnswersList []*Answer `json:"answers"`
}

//QuestionBank stuct
type QuestionBank struct {
	//ID           int64  `json:"id"`
	QsBId        int64  `json:"qb_id"`
	QBName       string `json:"qb_name"`
	CdTime       int64  `json:"cd_time"`
	MaxRetryTime int64  `json:"max_retry_time"`
	IsDeleted    int8   `json:"is_deleted"`
}

// QusBankSt 返回
type QusBankSt struct {
	QuestionBank
	ID        int64 `json:"id"`
	TotalCnt  int64 `json:"total_cnt"`
	EasyCnt   int64 `json:"easy_cnt"`
	NormalCnt int64 `json:"normal_cnt"`
	HardCnt   int64 `json:"hard_cnt"`
}

// QusBankCnt 统计类
type QusBankCnt struct {
	ID        int64         `json:"id"`
	TotalCnt  int64         `json:"total_cnt"`
	EasyCnt   sql.NullInt64 `json:"easy_cnt"`
	NormalCnt sql.NullInt64 `json:"normal_cnt"`
	HardCnt   sql.NullInt64 `json:"hard_cnt"`
}

//Answer stuct
type Answer struct {
	QsID          int64  `json:"qid"`
	AnswerContent string `json:"answer_content"`
	IsCorrect     int8   `json:"is_correct"`
	AnswerID      int64  `json:"answer_id"`
}

// AnswerAdd add
type AnswerAdd struct {
	Answer
}

// AddReturn return
type AddReturn struct {
	ID int64 `json:"id"`
}

// Page page
type Page struct {
	Total    int64 `json:"total"`
	PageNo   int   `json:"page_no" default:"1"`
	PageSize int   `json:"page_size" default:"20"`
}

// QuestionBankBind 绑定题库字段
type QuestionBankBind struct {
	ID             int64     `json:"id"`
	TargetItem     string    `json:"target_item" validate:"required"`
	TargetItemType int8      `json:"target_item_type" validate:"required"`
	QsBId          int64     `json:"bank_id" validate:"required"`
	UseInTime      int64     `json:"use_in_time" validate:"required"`
	Source         int8      `json:"source" validate:"required"`
	IsDeleted      int8      `json:"is_deleted"`
	Ctime          time.Time `json:"ctime"`
	Mtime          time.Time `json:"mtime"`

	QuestionBank *QuestionBank `json:"question_bank,omitempty"`
}

// RespList 返回
type RespList struct {
	Page
	Items interface{} `json:"items"`
}

// AddLog log
type AddLog struct {
	UID       string
	QsID      int64
	Platform  int8
	Source    int8
	Ids       string
	IsCorrect int8
}
