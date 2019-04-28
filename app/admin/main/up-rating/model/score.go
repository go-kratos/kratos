package model

import (
	"time"
)

// ScoreType .
type ScoreType int8

// ScoreType enums
const (
	Magnetic ScoreType = iota
	Creativity
	Influence
	Credit
)

// RatingListArg .
type RatingListArg struct {
	ScoreDate string    `form:"score_date"`                        // 年月 "2006-01"
	Mid       int64     `form:"mid"`                               // up id
	Tags      []int64   `form:"tag_ids,split" validate:"required"` // 分区
	ScoreType ScoreType `form:"score_type" default:"0"`            // 分数段类型
	ScoreMin  int64     `form:"score_min"`                         // 左闭右开
	ScoreMax  int64     `form:"score_max"`                         // 左闭右开
	From      int64     `form:"from" default:"0" validate:"min=0"`
	Limit     int64     `form:"limit" default:"20" validate:"min=1"`
}

// RatingListResp .
type RatingListResp struct {
	Result []*RatingInfo `json:"result"`
}

// RatingInfo .
type RatingInfo struct {
	Mid             int64     `json:"mid"`
	TagID           int       `json:"tag_id"`
	ScoreDate       time.Time `json:"-"`
	Date            string    `json:"date"`
	NickName        string    `json:"nickname"`
	TotalFans       int64     `json:"total_fans"`
	TotalAvs        int64     `json:"total_avs"`
	CreativityScore int64     `json:"creativity_score"`
	InfluenceScore  int64     `json:"influence_score"`
	CreditScore     int64     `json:"credit_score"`
	MagneticScore   int64     `json:"magnetic_score"`
}

// Paging .
type Paging struct {
	Ps    int64 `json:"page_size"`
	Total int64 `json:"total"`
}

// UpRatingHistoryArg .
type UpRatingHistoryArg struct {
	Mid       int64     `form:"mid" validate:"required"`
	Month     int       `form:"month" default:"0" validate:"min=0"`
	ScoreType ScoreType `form:"score_type" default:"0"`
}

// UpRatingHistoryResp .
type UpRatingHistoryResp struct {
	Data []*UpScoreHistory `json:"score_data"`
}

// UpScoreHistory .
type UpScoreHistory struct {
	ScoreType ScoreType `json:"type"`
	Date      []int64   `json:"date"`
	Score     []int64   `json:"score"`
}

// ScoreCurrentResp .
type ScoreCurrentResp struct {
	Date       int64         `json:"date"`
	Credit     *ScoreCurrent `json:"credit_score"`
	Influence  *ScoreCurrent `json:"influence_score"`
	Creativity *ScoreCurrent `json:"creativity_score"`
}

// ScoreCurrent .
type ScoreCurrent struct {
	Current int64 `json:"current"`
	Diff    int64 `json:"diff"`
}
