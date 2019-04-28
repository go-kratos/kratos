package match

import "encoding/json"

// ActUpdate .
const (
	ActUpdate  = "update"
	ActInsert  = "insert"
	ActDelete  = "delete"
	ResultNo   = 0
	ResultHome = 1
	ResultDraw = 2
	ResultAway = 3
)

// Message canal binlog message.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// ActMatchObj match object struct.
type ActMatchObj struct {
	ID      int64 `json:"id"`
	MatchID int64 `json:"match_id"`
	SID     int64 `json:"sid"`
	Result  int   `json:"result"`
}

// ActMatchUser match user.
type ActMatchUser struct {
	Mid    int64
	Result int
	Stake  int64
}
