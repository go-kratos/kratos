package archive

import (
	"encoding/json"
	"time"
)

const (
	//FlowPoolRecheck 回查pool（含热门回查、频道回查）
	FlowPoolRecheck = 4
	//FLowGroupIDChannel 频道回查的流量控制分组id
	FLowGroupIDChannel = 23
	//FlowGroupIDHot 热门回查的流量控制分组id
	FlowGroupIDHot = 24

	//FlowOpen 开启
	FlowOpen = int8(0)
	//FlowDelete 取消
	FlowDelete = int8(1)

	//FlowLogAdd 流量添加日志
	FlowLogAdd = int8(1)
	//FlowLogUpdate 流量更新日志
	FlowLogUpdate = int8(2)
	//FlowLogDel 流量删除日志
	FlowLogDel = int8(3)

	//PoolArc 稿件流量
	PoolArc = int8(0)
	//PoolUp up主流量
	PoolUp = int8(1)
	//PoolPrivateOrder 私单流量
	PoolPrivateOrder = int8(2)
	//PoolArticle 专栏流量
	PoolArticle = int8(3)
	//PoolArcForbid 稿件禁止流量
	PoolArcForbid = int8(4)
)

// Flow info
type Flow struct {
	ID     int64           `json:"id"`
	Remark string          `json:"remark"`
	Rank   int64           `json:"rank"`
	Type   int8            `json:"type"`
	Value  json.RawMessage `json:"value"`
	CTime  time.Time       `json:"ctime"`
	Pool   int8            `json:"pool"`
	State  int8            `json:"state"`
}

//FlowData Flow data
type FlowData struct {
	ID         int64     `json:"id"`
	Pool       int8      `json:"pool"`
	OID        int64     `json:"oid"`
	UID        int64     `json:"uid"`
	Parent     int8      `json:"parent"`
	GroupID    int64     `json:"group_id"`
	Remark     string    `json:"remark"`
	State      int8      `json:"state"`
	CTime      time.Time `json:"ctime"`
	MTime      time.Time `json:"mtime"`
	GroupValue []byte    `json:"group_value"`
}
