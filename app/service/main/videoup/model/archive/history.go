package archive

import "go-common/library/time"

// ArcHistory archive edit history.
type ArcHistory struct {
	ID           int64           `json:"id,omitempty"`
	Aid          int64           `json:"aid,omitempty"`
	Title        string          `json:"title,omitempty"`
	Tag          string          `json:"tag,omitempty"`
	Content      string          `json:"content,omitempty"`
	Cover        string          `json:"cover,omitempty"`
	Mid          int64           `json:"mid,omitempty"`
	CTime        time.Time       `json:"ctime,omitempty"`
	VideoHistory []*VideoHistory `json:"videos,omitempty"`
}

// VideoHistory video edit history.
type VideoHistory struct {
	ID       int64  `json:"-"`
	Aid      int64  `json:"-"`
	Cid      int64  `json:"cid"`
	Hid      int64  `json:"-"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Filename string `json:"filename"`
}
