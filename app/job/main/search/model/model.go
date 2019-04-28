package model

import (
	"encoding/json"
	"time"
)

// Stat all data statistics
type Stat struct {
	Counts int `json:"counts"`
}

// Message canal binlog message.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// DatabusPool poll from db.
var DatabusPool = []string{"dm", "dmreport_new"}

// JSONTime .
type JSONTime time.Time

// UnmarshalJSON .
func (p *JSONTime) UnmarshalJSON(data []byte) error {
	local, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*p = JSONTime(local)
	return err
}

// MarshalJSON .
func (p JSONTime) MarshalJSON() ([]byte, error) {
	data := make([]byte, 0)
	data = append(data, '"')
	data = time.Time(p).AppendFormat(data, "2006-01-02 15:04:05")
	data = append(data, '"')
	return data, nil
}

// String .
func (p JSONTime) String() string {
	return time.Time(p).Format("2006-01-02 15:04:05")
}
