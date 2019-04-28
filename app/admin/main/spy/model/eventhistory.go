package model

import (
	"time"
)

// EventHistory def.
type EventHistory struct {
	ID         int64
	Mid        int64     // 用户ID
	EventID    int64     // 事件ID
	Score      int8      // 用户真实分
	BaseScore  int8      // 基础信息得分
	EventScore int8      // 事件得分
	Remark     string    // 备注
	Reason     string    // 原因
	FactorVal  float32   // 风险因子
	Ctime      time.Time // 创建时间
	TargetID   int64     // 目标id
	TargetMid  int64     // 目标mid
	SpyTime    time.Time // 作弊时间
}

// EventHistoryDto dto.
type EventHistoryDto struct {
	ID         int64  `json:"id"`
	Score      int8   `json:"score"`       // 用户真实分
	BaseScore  int8   `json:"base_score"`  // 基础信息得分
	EventScore int8   `json:"event_score"` // 事件得分
	Reason     string `json:"reason"`      // 原因
	Ctime      int64  `json:"ctime"`       // 创建时间
	TargetID   int64  `json:"target_id"`   // 目标id
	TargetMid  int64  `json:"target_mid"`  // 目标mid
	SpyTime    int64  `json:"spy_time"`    // 作弊时间
}

// HisParamReq def.
type HisParamReq struct {
	Mid    int64
	Pn, Ps int
}

// HistoryPage def.
type HistoryPage struct {
	TotalCount int                `json:"total_count"`
	Pn         int                `json:"pn"`
	Ps         int                `json:"ps"`
	Items      []*EventHistoryDto `json:"items"`
}
