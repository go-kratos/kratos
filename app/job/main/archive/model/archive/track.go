package archive

import (
	"go-common/library/time"
)

type Track struct {
	// common values
	Timestamp time.Time `json:"timestamp"`
	// archive stat
	State int `json:"state"`
	Round int `json:"round"`
	// AID    int64  `json:"aid,omitempty"`
	Remark    string `json:"remark,omitempty"`
	Attribute int32  `json:"attribute"`
}

type VideoTrack struct {
	// common values
	Timestamp  time.Time `json:"timestamp"`
	XCodeState int8      `json:"xcode_state"`
	// video status
	Status int16  `json:"status"`
	AID    int64  `json:"aid,omitempty"`
	Remark string `json:"remark,omitempty"`
}
