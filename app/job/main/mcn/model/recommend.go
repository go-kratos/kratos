package model

import xtime "go-common/library/time"

// MCNUPRecommendSource .
type MCNUPRecommendSource int8

// const .
const (
	// MCNUPRecommendSourceUnKnown 未知来源
	MCNUPRecommendSourceUnKnown MCNUPRecommendSource = iota
	// MCNUPRecommendSourceAuto 自动添加(大数据)
	MCNUPRecommendSourceAuto
	// MCNUPRecommendStateManual  手动添加
	MCNUPRecommendStateManual
)

// MCNUPRecommendState .
type MCNUPRecommendState int8

// const .
const (
	// MCNUPRecommendStateUnKnown 未知状态
	MCNUPRecommendStateUnKnown MCNUPRecommendState = 0
	// MCNUPRecommendStateOff 未推荐
	MCNUPRecommendStateOff MCNUPRecommendState = 1
	// MCNUPRecommendStateOn  推荐中
	MCNUPRecommendStateOn MCNUPRecommendState = 2
	// MCNUPRecommendStateBan 禁止推荐
	MCNUPRecommendStateBan MCNUPRecommendState = 3
	// MCNUPRecommendStateDel 移除中
	MCNUPRecommendStateDel MCNUPRecommendState = 100
)

// McnUpRecommendPool .
type McnUpRecommendPool struct {
	ID                     int64                `json:"id"`
	UpMid                  int64                `json:"up_mid"`
	UpName                 string               `json:"up_name"`
	FansCount              int64                `json:"fans_count"`
	FansCountIncreaseMonth int64                `json:"fans_count_increase_month"`
	ArchiveCount           int64                `json:"archive_count"`
	PlayCountAccumulate    int64                `json:"play_count_accumulate"`
	PlayCountAverage       int64                `json:"play_count_average"`
	ActiveTid              int16                `json:"active_tid"`
	LastArchiveTime        xtime.Time           `json:"last_archive_time"`
	State                  MCNUPRecommendState  `json:"state"`
	Source                 MCNUPRecommendSource `json:"source"`
	GenerateTime           xtime.Time           `json:"generate_time"`
	Ctime                  xtime.Time           `json:"ctime"`
	Mtime                  xtime.Time           `json:"mtime"`
}
