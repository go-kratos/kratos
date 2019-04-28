package monitor

const (
	// RedisPrefix 参数:business。参数1：bid；参数2：监控ID
	RedisPrefix = "monitor_stats_%d"
	// RedisDelArcInfo 稿件删除监控key
	RedisDelArcInfo = "monitor_stats_del_arc"
	// BusVideo 视频业务
	BusVideo = 1
	// BusArc 稿件业务
	BusArc = 2

	// NotifyTypeEmail 邮件通知
	NotifyTypeEmail = 1
	// NotityTypeSms 短信通知
	NotityTypeSms = 2

	// 稿件业务常量
	// ArchiveBitPGC 稿件PGC属性位
	ArchiveBitPGC = 9
	// ArchiveStateDel 稿件删除状态
	ArchiveStateDel = -100
	// ArchiveOriginal 自制稿件
	ArchiveOriginal = 1

	// RuleHighUpDelArc 高能联盟UP主大量删除稿件监控
	RuleHighUpDelArc = 1
	// RuleFamUpDelArc 大UP主大量删除稿件监控
	RuleFamUpDelArc = 17

	// CompGT 大于
	CompGT = ">"
	// CompLT 小于
	CompLT = "<"
	// CompGET 大于等于
	CompGET = ">="
	// CompLET 小于等于
	CompLET = "<="
	// CompNE 不等于
	CompNE = "!="
	// CompE 等于
	CompE = "="
)

var (
	// SpecialTypeIDs 特殊分区（变更非常不频繁）
	SpecialTypeIDs = map[int64]int8{
		15: 1, 34: 1, 32: 1, 82: 1, 33: 1, 83: 1, 145: 1, 146: 1,
		147: 1, 153: 1, 185: 1, 186: 1, 187: 1, 37: 1, 178: 1, 179: 1,
		180: 1, 128: 1, 85: 1, 86: 1, 183: 1,
	}
)

// BinlogArchive 稿件 binlog 结构
type BinlogArchive struct {
	ID        int64         `json:"id"`
	State     int64         `json:"state"`
	Round     int64         `json:"round"`
	MID       int64         `json:"mid"`
	Attr      int64         `json:"attribute"`
	TypeID    int64         `json:"typeid"`
	IsSpecTID int8          `json:"is_special_tid"`
	HumanRank int           `json:"humanrank"`
	Duration  int           `json:"duration"`
	Desc      string        `json:"desc"`
	Title     string        `json:"title"`
	Cover     string        `json:"cover"`
	Content   string        `json:"content"`
	Tag       string        `json:"tag"`
	Copyright int8          `json:"copyright"`
	AreaLimit int8          `json:"arealimit"`
	Author    string        `json:"author"`
	Access    int           `json:"access"`
	Forward   int           `json:"forward"`
	PubTime   string        `json:"pubtime"`
	Reason    string        `json:"reject_reason"`
	CTime     string        `json:"ctime"`
	MTime     string        `json:"mtime"`
	PTime     string        `json:"ptime"`
	Addit     *ArchiveAddit `json:"_"`
}

// BinlogVideo 视频binlog结构
type BinlogVideo struct {
	ID          int64  `json:"id"`
	Filename    string `json:"filename"`
	Cid         int64  `json:"cid"`
	Aid         int64  `json:"aid"`
	Title       string `json:"eptitle"`
	Desc        string `json:"description"`
	SrcType     string `json:"src_type"`
	Duration    int64  `json:"duration"`
	Filesize    int64  `json:"filesize"`
	Resolutions string `json:"resolutions"`
	Playurl     string `json:"playurl"`
	FailCode    int8   `json:"failinfo"`
	Index       int    `json:"index_order"`
	Attribute   int32  `json:"attribute"`
	XcodeState  int8   `json:"xcode_state"`
	State       int8   `json:"state"`
	Status      int16  `json:"status"`
	CTime       string `json:"ctime"`
	MTime       string `json:"mtime"`
}

// ArchiveAddit 稿件附加属性
type ArchiveAddit struct {
	Aid           int64  `json:"aid"`
	MissionID     int64  `json:"mission_id"`
	UpFrom        int8   `json:"up_from"`
	FromIP        int64  `json:"from_ip"`
	IPv6          []byte `json:"ipv6"`
	Source        string `json:"source"`
	OrderID       int64  `json:"order_id"`
	RecheckReason string `json:"recheck_reason"`
	RedirectURL   string `json:"redirect_url"`
	FlowID        int64  `json:"flow_id"`
	Advertiser    string `json:"advertiser"`
	FlowRemark    string `json:"flow_remark"`
	DescFormatID  int64  `json:"desc_format_id"`
	Desc          string `json:"desc"`
	Dynamic       string `json:"dynamic"`
}

// RuleResultRes 监控结果
type RuleResultRes struct {
	Code int               `json:"code"`
	Data []*RuleResultData `json:"data"`
}

// RuleResultData 监控结果
type RuleResultData struct {
	Rule  *Rule  `json:"rule"`
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
		Comp string `json:"comparison"`
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

// FieldsConf 监控字段配置
type FieldsConf struct {
	Comparison string
}

// DelArcInfo UP主删稿信息
type DelArcInfo struct {
	AID   int64  `json:"aid"`
	MID   int64  `json:"mid"`
	Time  string `json:"time"`
	Title string `json:"title"`
}
