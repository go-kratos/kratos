package model

import "go-common/library/time"

// BaseInfo up rating info
type BaseInfo struct {
	ID        int64
	MID       int64
	TagID     int64
	PlayIncr  int64
	CoinIncr  int64
	Avs       int64
	MAAFans   int64
	MAHFans   int64
	OpenAvs   int64
	LockedAvs int64
	Date      time.Time
	TotalFans int64
	TotalAvs  int64
	TotalCoin int64
	TotalPlay int64
}

// RatingParameter rating parameter
type RatingParameter struct {
	WDP      int64 // dp weight
	WDC      int64 // dc weight
	WDV      int64 // dv weight
	WMDV     int64 // mdv weight
	WCS      int64
	WCSR     int64
	WMAAFans int64
	WMAHFans int64
	WIS      int64
	WISR     int64
	// 信用分
	HBASE int64
	HR    int64
	HV    int64
	HVM   int64
	HL    int64
	HLM   int64
}

// Rating rating
type Rating struct {
	MID                 int64
	TagID               int64
	MetaCreativityScore int64
	CreativityScore     int64
	MetaInfluenceScore  int64
	InfluenceScore      int64
	CreditScore         int64
	MagneticScore       int64
	Score               int64
	Date                time.Time
}

// Past past stat
type Past struct {
	MID                 int64
	MetaCreativityScore int64
	MetaInfluenceScore  int64
	CreditScore         int64
}
