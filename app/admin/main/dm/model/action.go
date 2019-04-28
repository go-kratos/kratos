package model

import (
	"encoding/json"
)

// All const variable used in admin
const (
	ActFlushDM   = "flush"
	ActReportDel = "report_del"
)

// Action job msg.
type Action struct {
	Oid    int64           `json:"oid"`
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}
