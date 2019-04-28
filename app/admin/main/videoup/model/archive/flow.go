package archive

// pool .
const (
	PoolArc          = int8(0)
	PoolUp           = int8(1)
	PoolPrivateOrder = int8(2)
	PoolArticle      = int8(3)
	PoolArcForbid    = int8(4)

	FlowOpen   = int8(0)
	FlowDelete = int8(1)

	FlowLogAdd    = int8(1)
	FlowLogUpdate = int8(2)
	FlowLogDel    = int8(3)

	FlowGroupNoChannel = int64(23)
	FlowGroupNoHot     = int64(24)
)

var (
	//FlowAttrMap archive submit with flow attr
	FlowAttrMap = map[string]int64{
		"nochannel": FlowGroupNoChannel,
		"nohot":     FlowGroupNoHot,
	}
)

//FlowData flow_design data
type FlowData struct {
	ID         int64  `json:"id"`
	Pool       int8   `json:"pool"`
	OID        int64  `json:"oid"`
	GroupID    int64  `json:"group_id"`
	Parent     int8   `json:"parent"`
	State      int8   `json:"state"`
	GroupValue []byte `json:"group_value"`
}
