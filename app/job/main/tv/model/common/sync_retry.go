package common

import "fmt"

const (
	_ctypeSN = "retry_sn"
	_ctypeEP = "retry_ep"
)

// SyncRetry is the struct used for retry info storage
type SyncRetry struct {
	Ctype string
	Retry int
	CID   int64
}

// FromSn def.
func (v *SyncRetry) FromSn(retry int, sid int64) {
	v.Ctype = _ctypeSN
	v.Retry = retry
	v.CID = sid
}

// FromEp def.
func (v *SyncRetry) FromEp(retry int, epid int64) {
	v.Ctype = _ctypeEP
	v.Retry = retry
	v.CID = epid
}

// MCKey def.
func (v *SyncRetry) MCKey() (key string) {
	return v.Ctype + "_" + fmt.Sprintf("%d", v.CID)
}
