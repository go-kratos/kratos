package track

import (
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/time"
)

const ()

var ()

type Archive struct {
	// common values
	Timestamp time.Time `json:"timestamp"`
	// archive stat
	State int `json:"state"`
	Round int `json:"round"`
	// AID    int64  `json:"aid,omitempty"`
	Remark    string `json:"remark,omitempty"`
	Attribute int32  `json:"attribute"`
}

type Video struct {
	// common values
	Timestamp  time.Time `json:"timestamp"`
	XCodeState int8      `json:"xcode_state"`
	// video status
	Status    int16  `json:"status"`
	AID       int64  `json:"aid,omitempty"`
	Remark    string `json:"remark,omitempty"`
	Attribute int32  `json:"attribute"`
}

// Archives Archive sorted.
type Archives []*Archive

func (a Archives) Len() int           { return len(a) }
func (a Archives) Less(i, j int) bool { return int64(a[i].Timestamp) > int64(a[j].Timestamp) }
func (a Archives) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//ArcTrackInfo 稿件追踪信息
type ArcTrackInfo struct {
	EditHistory []*archive.EditHistory `json:"edit_history"`
	Track       []*Archive             `json:"track"`
	Relation    [][]int                `json:"relation"`
}
