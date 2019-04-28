package monitor

const (
	// RedisPrefix 参数:business。参数1：bid；参数2：监控ID
	RedisPrefix = "monitor_stats_%d"

	// CompGT 大于
	CompGT = ">"
	// CompLT 小于
	CompLT = "<"
)

// RuleResultData 返回的结果数据
type RuleResultData struct {
	Rule  *Rule  `json:"rule"`
	User  *User  `json:"user"`
	Stats *Stats `json:"stats"`
}

// Rule 监控规则信息
type Rule struct {
	ID       int64     `json:"id"`
	Type     int8      `json:"type"`
	BID      int8      `json:"bid"`
	Name     string    `json:"name"`
	State    int8      `json:"state"`
	STime    string    `json:"stime"`
	ETime    string    `json:"etime"`
	CTime    string    `json:"ctime"`
	MTime    string    `json:"mtime"`
	UID      int64     `json:"uid"`
	RuleConf *RuleConf `json:"rule"`
}

// RuleConf 监控方案配置结构体
type RuleConf struct {
	Name    string              `json:"name"`
	MoniCdt map[string]struct { //监控方案的监控条件
		Comp  string `json:"comparison"`
		Value int64  `json:"value"`
	} `json:"moni_cdt"`
	NotifyCdt map[string]struct { //达到发送通知的条件
		Comp  string `json:"comparison"`
		Value int64  `json:"value"`
	} `json:"notify_cdt"`
	Notify struct { //通知类型配置
		Way    int8     `json:"way"`
		Member []string `json:"member"`
	} `json:"notify"`
}

// Stats 监控统计
type Stats struct {
	TotalCount int `json:"total_count"`
	MoniCount  int `json:"moni_count"`
	MaxTime    int `json:"max_time"`
}

// User manager user struct
type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	NickName string `json:"nickname"`
}
