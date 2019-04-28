package archive

import (
	"go-common/library/time"
)

// ArcHistory  archive of history
type ArcHistory struct {
	ID      int64     `json:"id,omitempty"`
	Aid     int64     `json:"aid,omitempty"`
	Mid     int64     `json:"mid,omitempty"`
	Tag     string    `json:"tag,omitempty"`
	Title   string    `json:"title,omitempty"`
	Content string    `json:"content,omitempty"`
	Cover   string    `json:"cover,omitempty"`
	CTime   time.Time `json:"ctime,omitempty"`
	Video   []*Video  `json:"videos,omitempty"`
}
