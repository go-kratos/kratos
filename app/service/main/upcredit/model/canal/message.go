package canal

import "encoding/json"

// Msg canal databus msg.
type Msg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}
