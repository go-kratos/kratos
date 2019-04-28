package xreply

import "go-common/app/interface/main/reply/model/reply"

// const
const (
	MaxPageSize = 50

	ModeOrigin = 1 // origin
	ModeTime   = 2 // sort by time
	ModeHot    = 3 // sort by hot

	FolderKindSub  = "s"
	FolderKindRoot = "r"

	CursorModePage   = 1 // pn ps翻页的
	CursorModeCursor = 2 // 按游标翻页的
)

type ReplyReq struct {
	CommonReq
	ReplyCommonReq
	Cursor CursorReq
}

var (
	_SupportModeAll    = []int{ModeOrigin, ModeTime, ModeHot}
	_SupportModeOrigin = []int{ModeOrigin}
)

//  ...
func (req *ReplyReq) ModeInfo(hotMap map[int64]int8, floorMap map[int64]int8) (mode int, supportMode []int) {
	supportMode = _SupportModeAll
	switch req.Cursor.Mode {
	case ModeHot:
		mode = ModeHot
	case ModeTime:
		mode = ModeTime
	case ModeOrigin:
		mode = ModeOrigin
		supportMode = _SupportModeOrigin
	default:
		if tp, ok := hotMap[req.Oid]; ok && tp == req.Type {
			mode = ModeHot
		} else if tp, ok := floorMap[req.Oid]; ok && tp == req.Type {
			mode = ModeTime
		} else {
			mode = ModeOrigin
			supportMode = _SupportModeOrigin
		}
	}
	return
}

type ReplyRes struct {
	Cursor  CursorRes      `json:"cursor"`
	Hots    []*reply.Reply `json:"hots"`
	Notice  *reply.Notice  `json:"notice"`
	Replies []*reply.Reply `json:"replies"`
	Top     TopReply       `json:"top"`
	Folder  reply.Folder   `json:"folder"`
	CommonRes
}

type CommonRes struct {
	Assist    int         `json:"assist"`
	Blacklist int         `json:"blacklist"`
	Config    ReplyConfig `json:"config"`
	Upper     Upper       `json:"upper"`
}

type TopReply struct {
	Admin *reply.Reply `json:"admin"`
	Upper *reply.Reply `json:"upper"`
}

type Upper struct {
	Mid int64 `json:"mid"`
}

type ReplyConfig struct {
	ShowAdmin int8 `json:"showadmin"`
	ShowEntry int8 `json:"showentry"`
	ShowFloor int8 `json:"showfloor"`
}

// CommonReq ...
type CommonReq struct {
	Plat    int8   `form:"plat"`
	Build   int64  `form:"build"`
	Buvid   string `form:"buvid"`
	MobiApp string `form:"mobi_app"`
	Mid     int64  `form:"mid"`
	IP      string `form:"ip`
}

// ReplyCommonReq ...
type ReplyCommonReq struct {
	Oid  int64 `form:"oid" validate:"required"`
	Type int8  `form:"type" validate:"required"`
}

// Cursor Common Cursor
type Cursor struct {
	IsBegin bool `json:"is_begin"`
	Prev    int  `json:"prev"`
	Next    int  `json:"next"`
	IsEnd   bool `json:"is_end"`
	Ps      int  `json:"ps"`
}

// Latest ...
func (c *Cursor) Latest() bool {
	return c.Next == 0 && c.Prev == 0
}

// Forward ...
func (c *Cursor) Forward() bool {
	return c.Next != 0
}

// Backward ...
func (c *Cursor) Backward() bool {
	return c.Prev != 0
}

// CursorRes ...
type CursorRes struct {
	AllCount    int   `json:"all_count,omitempty"`
	IsBegin     bool  `json:"is_begin"`
	Prev        int   `json:"prev"`
	Next        int   `json:"next"`
	IsEnd       bool  `json:"is_end"`
	Ps          int   `json:"ps,omitempty"`
	Mode        int   `json:"mode,omitempty"`
	SupportMode []int `json:"support_mode,omitempty"`
}

// CursorReq ...
type CursorReq struct {
	Ps   int `form:"ps" validate:"omitempty,min=1,max=50" default:"20"`
	Prev int `form:"prev"`
	Next int `form:"next"`
	Mode int `form:"mode"`
}

// Legal ...
func (cq *CursorReq) Legal() bool {
	if cq.Next != 0 && cq.Prev != 0 {
		return false
	}
	return true
}

func (cq *CursorReq) Forward() bool {
	return cq.Next != 0
}

func (cq *CursorReq) Backward() bool {
	return cq.Prev != 0
}

// Latest ...
func (cq *CursorReq) Latest() bool {
	if cq.Next == 0 && cq.Prev == 0 {
		return true
	}
	return false
}

type SubFolderReq struct {
	CommonReq
	ReplyCommonReq
	Cursor CursorReq
}

type RootFolderReq struct {
	CommonReq
	ReplyCommonReq
	Cursor CursorReq
	Root   int64 `form:"root" validate:"required"`
}

type SubFolderRes struct {
	Cursor  CursorRes      `json:"cursor"`
	Replies []*reply.Reply `json:"replies"`
	CommonRes
}

type RootFolderRes struct {
	Cursor  CursorRes      `json:"cursor"`
	Replies []*reply.Reply `json:"replies"`
	CommonRes
}
