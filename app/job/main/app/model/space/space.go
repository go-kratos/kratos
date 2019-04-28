package space

import xtime "go-common/library/time"

type Attrs struct {
	Archive bool `json:"archive,omitempty"`
	Article bool `json:"article,omitempty"`
	Clip    bool `json:"clip,omitempty"`
	Album   bool `json:"album,omitempty"`
	Audio   bool `json:"audio,omitempty"`
}

type Item struct {
	ID    int64      `json:"id,omitempty"`
	Goto  string     `json:"goto,omitempty"`
	CTime xtime.Time `json:"ctime,omitempty"`
}

type Clip struct {
	VideoID int64      `json:"video_id,omitempty"`
	CTime   xtime.Time `json:"ctime,omitempty"`
}

type Album struct {
	DocID int64      `json:"doc_id,omitempty"`
	CTime xtime.Time `json:"ctime,omitempty"`
}

type Audio struct {
	ID    int64      `json:"audioId,omitempty"`
	CTime xtime.Time `json:"cTime,omitempty"`
}
