package model

const (
	// MonitorClose 取消监控
	MonitorClose = int32(0)
	// MonitorOpen 打开监控
	MonitorOpen = int32(1)
	// MonitorAudit 先审后发
	MonitorAudit = int32(2)

	// MonitorStatsAll all
	MonitorStatsAll = 1
	// MonitorStatsUser all
	MonitorStatsUser = 2
)

// SearchMonitorParams search params.
type SearchMonitorParams struct {
	Mode     int8
	Type     int8
	Oid      int64
	UID      int64
	NickName string
	Keyword  string
	Sort     string
	Order    string
}

// SearchMonitor search monitor.
type SearchMonitor struct {
	ID          int64  `json:"id"`
	Oid         int64  `json:"oid"`
	OidStr      string `json:"oid_str"`
	Type        int8   `json:"type"`
	Mid         int64  `json:"mid"`
	State       int8   `json:"state"`
	Attr        int32  `json:"attr"`
	Ctime       string `json:"ctime"`
	Mtime       string `json:"mtime"`
	Title       string `json:"title"`
	Uname       string `json:"uname"`
	UnverifyNmu int    `json:"unverify_num"`
	MCount      int32  `json:"mcount"`
	DocID       string `json:"doc_id"`
	Remark      string `json:"remark"`
}

// SearchMonitorResult search result.
type SearchMonitorResult struct {
	Code      int              `json:"code,omitempty"`
	Page      int64            `json:"page"`
	PageSize  int64            `json:"pagesize"`
	PageCount int64            `json:"pagecount"`
	Total     int64            `json:"total"`
	Order     string           `json:"order"`
	Result    []*SearchMonitor `json:"result"`
	Message   string           `json:"msg,omitempty"`
}

// StatsMonitor stats monitor.
type StatsMonitor struct {
	Date           string `json:"date"`
	AdminID        int64  `json:"adminid"`
	MonitorTotal   int64  `json:"monitor_total"`
	MonitorPending int64  `json:"monitor_pending"`
	MonitorPass    int64  `json:"monitor_pass"`
	MonitorDel     int64  `json:"monitor_del"`
	MonitorAvgCost string `json:"monitor_avg_cost"`
}

// StatsMonitorResult search result.
type StatsMonitorResult struct {
	Code      int             `json:"code,omitempty"`
	Page      int             `json:"page"`
	PageSize  int             `json:"pagesize"`
	PageCount int             `json:"pagecount"`
	Total     int             `json:"total"`
	Order     string          `json:"order"`
	Message   string          `json:"msg,omitempty"`
	Result    []*StatsMonitor `json:"result"`
}

// MonitorLogResult MonitorLogResult
type MonitorLogResult struct {
	Logs  []*MonitorLog `json:"logs"`
	Page  Page          `json:"page"`
	Order string        `json:"order"`
	Sort  string        `json:"sort"`
}

// Page Page
type Page struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// MonitorLog MonitorLog
type MonitorLog struct {
	AdminID     int64  `json:"adminid"`
	AdminName   string `json:"admin_name"`
	Oid         int64  `json:"oid"`
	OidStr      string `json:"oid_str"`
	Type        int32  `json:"type"`
	Title       string `json:"title"`
	RedirectURL string `json:"redirect_url"`
	Remark      string `json:"remark"`
	UserName    string `json:"username"`
	Mid         int64  `json:"mid"`
	CTime       string `json:"ctime"`
	LogState    int64  `json:"log_state"`
	State       int64  `json:"state"`
}
