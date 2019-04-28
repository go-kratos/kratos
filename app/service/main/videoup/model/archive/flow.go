package archive

import (
	"encoding/json"
	"go-common/library/time"
)

// pool .
const (
	PoolArc       = int8(0)
	PoolUp        = int8(1)
	PoolPorder    = int8(2)
	PoolArticle   = int8(3)
	PoolArcForbid = int8(4)
	PoolArcPGC    = int8(5)

	FlowOpen   = int8(0)
	FlowDelete = int8(1)

	FlowLogAdd    = int8(1)
	FlowLogUpdate = int8(2)
	FlowLogDel    = int8(3)

	FlowGroupNoChannel   = int64(23)
	FlowGroupNoHot       = int64(24)
	FlowGroupNoTimeline  = int64(25)
	FlowGroupNoOtt       = int64(26)
	FlowGroupNoRecommend = int64(27)
	FlowGroupNoRank      = int64(28)
)

var (
	//FlowAttrMap archive submit with flow attr
	FlowAttrMap = map[string]int64{
		"nochannel":   FlowGroupNoChannel,
		"nohot":       FlowGroupNoHot,
		"notimeline":  FlowGroupNoTimeline,
		"noott":       FlowGroupNoOtt,
		"norecommend": FlowGroupNoRecommend,
		"norank":      FlowGroupNoRank,
	}
)

// Flow info
type Flow struct {
	ID     int64           `json:"id"`
	Remark string          `json:"remark"`
	Rank   int64           `json:"rank"`
	Type   int8            `json:"type"`
	Value  json.RawMessage `json:"value"`
	CTime  time.Time       `json:"ctime"`
}

//FlowData Flow data
type FlowData struct {
	ID         int64     `json:"id"`
	Pool       int8      `json:"pool"`
	OID        int64     `json:"oid"`
	UID        int64     `json:"uid"`
	Parent     int8      `json:"parent"`
	State      int8      `json:"state"`
	GroupID    int64     `json:"group_id"`
	Remark     string    `json:"remark"`
	GroupValue []byte    `json:"group_value"`
	CTime      time.Time `json:"ctime"`
	MTime      time.Time `json:"mtime"`
}

//FlowPagerData .
type FlowPagerData struct {
	Items []*FlowData `json:"items"`
	Pager *Pager      `json:"pager,omitempty"`
}

//Pager .
type Pager struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}
