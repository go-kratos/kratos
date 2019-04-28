package reply

import (
	"go-common/app/interface/main/reply/model/reply"
	xtime "go-common/library/time"
)

// Reply str
type Reply struct {
	RpID      int64      `json:"rpid"`
	Oid       int64      `json:"oid"`
	Type      int8       `json:"type"`
	Mid       int64      `json:"mid"`
	Root      int64      `json:"root"`
	Parent    int64      `json:"parent"`
	Count     int        `json:"count"`
	RCount    int        `json:"rcount"`
	Floor     int        `json:"floor"`
	State     int8       `json:"state"`
	Attr      int8       `json:"attr"`
	CTime     xtime.Time `json:"ctime"`
	MTime     xtime.Time `json:"-"`
	RpIDStr   string     `json:"rpid_str,omitempty"`
	RootStr   string     `json:"root_str,omitempty"`
	ParentStr string     `json:"parent_str,omitempty"`
	// action count, from ReplyAction count
	Like   int  `json:"like"`
	Hate   int  `json:"-"`
	Action int8 `json:"action"`
	// member info
	Member *reply.Info `json:"member"`
	// other
	Content *CreativeReplyCont `json:"content"`
	Replies []*Reply           `json:"replies"`
}

// CreativeReplyCont str
type CreativeReplyCont struct {
	RpID    int64      `json:"-"`
	Message string     `json:"message"`
	Ats     Ints       `json:"ats,omitempty"`
	IP      uint32     `json:"ipi,omitempty"`
	Plat    int8       `json:"plat"`
	Device  string     `json:"device"`
	Version string     `json:"version,omitempty"`
	CTime   xtime.Time `json:"-"`
	MTime   xtime.Time `json:"-"`
	// ats member info
	Members []*reply.Info `json:"members"`
}
