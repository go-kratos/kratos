package model

import (
	"github.com/satori/go.uuid"
)

// UserLog log info.
type UserLog struct {
	Mid     int64             `json:"mid"`
	IP      string            `json:"ip"`
	TS      int64             `json:"ts"`
	LogID   string            `json:"log_id"`
	Content map[string]string `json:"content"`
}

// UUID4 is generate uuid
func UUID4() string {
	return uuid.NewV4().String()
}
