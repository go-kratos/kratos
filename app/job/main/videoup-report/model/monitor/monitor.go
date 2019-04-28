package monitor

const (
	// RedisPrefix 参数:business。
	RedisPrefix = "monitor_stats_%d_"
	// SuffixVideo 视频停留统计。参数:state
	SuffixVideo = "%d"
	// SuffixArc 稿件停留统计。参数:round。参数:state。
	SuffixArc = "%d_%d"
	BusVideo  = 1
	BusArc    = 2

	NotifyTypeEmail = 1
	NotityTypeSms   = 2

	RuleStateOK      = 1
	RuleStateDisable = 0
)

type RuleResultRes struct {
	Code int               `json:"code"`
	Data []*RuleResultData `json:"data"`
}
type RuleResultData struct {
	Rule  *Rule  `json:"rule"`
	Stats *Stats `json:"stats"`
}

// Rule 监控规则信息
type Rule struct {
	ID       int64     `json:"id"`
	Type     int8      `json:"type"`
	Business int8      `json:"business"`
	Name     string    `json:"name"`
	State    int8      `json:"state"`
	STime    string    `json:"s_time"`
	ETime    string    `json:"e_time"`
	RuleConf *RuleConf `json:"rule_conf"`
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

type Stats struct {
	TotalCount int `json:"total_count"`
	MoniCount  int `json:"moni_count"`
	MaxTime    int `json:"max_time"`
}
