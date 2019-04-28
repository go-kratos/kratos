package model

// all const variable used in dm monitor
const (
	// 监控状态
	MonitorClosed = int32(0)
	MonitorBefore = int32(1) // 先审后发
	MonitorAfter  = int32(2) // 先发后审
)

// MonitorResult dm monitor result
type MonitorResult struct {
	Order    string     `json:"order"`
	Sort     string     `json:"sort"`
	Page     int64      `json:"page"`
	PageSize int64      `json:"pagesize"`
	Total    int64      `json:"total"`
	Result   []*Monitor `json:"result"`
}

// Monitor dm monitors
type Monitor struct {
	ID     int64  `json:"id"`
	Type   int32  `json:"type"`
	Pid    int64  `json:"pid"`
	Oid    int64  `json:"oid"`
	State  int32  `json:"state"`
	MCount int64  `json:"mcount"`
	Ctime  string `json:"ctime"`
	Mtime  string `json:"mtime"`
	Mid    int64  `json:"mid"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// SearchMonitor dm monitor struct
type SearchMonitor struct {
	ID     int64  `json:"id"`
	Type   int32  `json:"type"`
	Pid    int64  `json:"pid"`
	Oid    int64  `json:"oid"`
	State  int32  `json:"state"`
	Attr   int32  `json:"attr"`
	MCount int64  `json:"mcount"`
	Ctime  string `json:"ctime"`
	Mtime  string `json:"mtime"`
	Mid    int64  `json:"mid"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Page search page info
type Page struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// SearchMonitorResult result from search
type SearchMonitorResult struct {
	Order  string           `json:"order"`
	Sort   string           `json:"sort"`
	Page   *Page            `json:"page"`
	Result []*SearchMonitor `json:"result"`
}
