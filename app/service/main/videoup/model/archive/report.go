package archive

import xtime "go-common/library/time"

// const ArcReport
const (
	ArcReportNew    = int8(0)
	ArcReportReject = int8(1)
	ArcReportAccept = int8(2)
)

// ArcReport  archive_report
type ArcReport struct {
	Mid    int64      `json:"mid"`
	Aid    int64      `json:"aid"`
	Type   int8       `json:"type"`
	Reason string     `json:"reason"`
	Pics   string     `json:"pics"`
	State  int8       `json:"state"`
	CTime  xtime.Time `json:"ctime"`
	MTime  xtime.Time `json:"-"`
}
