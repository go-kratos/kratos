package model

// const value.
const (
	RptStateUnhandledFirst  = int32(0)
	RptStateUnhandledSecond = int32(3)
	RptStateHandledFirst    = int32(4)
	RptStateHandledSecond   = int32(1)
	RptStateIgnoreSecond    = int32(2)
)

// Completed judge report id completed.
func (r *Report) Completed() bool {
	return r.State == RptStateHandledFirst || r.State == RptStateHandledSecond || r.State == RptStateIgnoreSecond
}

// SetManager set manager.
func (r *ReportUser) SetManager() {
	r.Attr = r.Attr | 2
}

// AddReportReq add report request params.
type AddReportReq struct {
	Oid      int64
	Tid      int64
	Type     int32
	Mid      int64
	PartID   int32
	ReasonID int32
	Score    int32
}
