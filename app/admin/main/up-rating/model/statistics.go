package model

import (
	"go-common/library/time"
)

// RatingStatis rating statitics
type RatingStatis struct {
	Ups             int64     `json:"ups"`
	Section         int64     `json:"-"`
	Tips            string    `json:"tips"`
	Score           int64     `json:"score"`
	TotalScore      int64     `json:"-"`
	CreativityScore int64     `json:"creativity_score"`
	InfluenceScore  int64     `json:"influence_score"`
	CreditScore     int64     `json:"credit_score"`
	Fans            int64     `json:"fans"`
	Avs             int64     `json:"avs"`
	Coin            int64     `json:"coin"`
	Play            int64     `json:"play"`
	CDate           time.Time `json:"-"`
	TagID           int64     `json:"-"`
	CType           int64     `json:"-"`
	Proportion      string    `json:"proportion"`
	Compare         int64     `json:"compare"`
	ComparePropor   string    `json:"compare_propor"`
}

// Trend rating trend
type Trend struct {
	MID             int64  `json:"mid"`
	Nickname        string `json:"nickname"`
	DValue          int    `json:"d_value"`
	MagneticScore   int64  `json:"magnetic_score"`
	CreativityScore int64  `json:"creativity_score"`
	InfluenceScore  int64  `json:"influence_score"`
	CreditScore     int64  `json:"credit_score"`
	Avs             int64  `json:"avs"`
	Fans            int64  `json:"fans"`
}
