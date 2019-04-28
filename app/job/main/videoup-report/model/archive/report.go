package archive

import (
	"encoding/json"
	"time"
)

//  1 耗时 2 耗时(30分钟) 3 视频进审/过审分布
var (
	ReportArchiveRound     = map[int8]string{30: "30", 40: "40", 90: "90"}
	ReportTypeTookMinute   = int8(1)
	ReportTypeTookHalfHour = int8(2)
	ReportTypeVideoAudit   = int8(3)
	ReportTypeArcMoveType  = int8(4)
	ReportTypeArcRoundFlow = int8(5)
	ReportTypeXcode        = int8(6) //video sd_finish,hd_finish,dispatch take time
	ReportTypeTraffic      = int8(7) //视频审核耗时统计。10分钟聚合的一转、一审、二转、分发耗时结果
)

// Report struct
type Report struct {
	ID      int64           `json:"-"`
	TypeID  int8            `json:"type"`
	Content json.RawMessage `json:"content"`
	CTime   time.Time       `json:"ctime"`
	MTime   time.Time       `json:"mtime"`
}
