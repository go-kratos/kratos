package model

import "go-common/library/time"

// RatingStatis rating statistics
type RatingStatis struct {
	Ups             int64
	Section         int64
	Tips            string
	TotalScore      int64
	CreativityScore int64
	InfluenceScore  int64
	CreditScore     int64
	Fans            int64
	Avs             int64
	Coin            int64
	Play            int64
	CDate           time.Time
	TagID           int64
	CType           int
}

// Diff rating diff
type Diff struct {
	MID             int64
	MagneticScore   int64
	CreativityScore int64
	InfluenceScore  int64
	CreditScore     int64
	MagneticDiff    int
	CreativityDiff  int
	InfluenceDiff   int
	CreditDiff      int
	TotalAvs        int64
	Fans            int64
	TagID           int64
	CType           int
	Section         int
	Tips            string
	Date            time.Time
}

// TopRating top rating
type TopRating struct {
	MID   int64
	CType int
	TagID int64
	Score int64
	Play  int64
	Fans  int64
}

const (
	// MAGNETIC magnetic ctype
	MAGNETIC = iota
	// CREATIVITY creativity ctype
	CREATIVITY
	// INFLUENCE influence ctype
	INFLUENCE
	// CREDIT influence ctype
	CREDIT
)

// GetScore get score
func (a *Diff) GetScore(ctype int) (score int64) {
	switch ctype {
	case MAGNETIC:
		return a.MagneticScore
	case CREATIVITY:
		return a.CreativityScore
	case INFLUENCE:
		return a.InfluenceScore
	case CREDIT:
		return a.CreditScore
	}
	return
}
