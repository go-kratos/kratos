package model

import (
	"encoding/json"
)

const (
	// MaxRetryTimes 最大重试次数.
	MaxRetryTimes = int(3)
	// TNameMaxLen tag name length .
	TNameMaxLen = 32
	// TagBatchNumMax .
	TagBatchNumMax = 10

	// TagStateNormal .
	TagStateNormal = int32(0)

	// ResTagBind .
	ResTagBind = "bind"
	// ResTagDelete .
	ResTagDelete = "delete"

	// ResTagStateNormal .
	ResTagStateNormal = int32(0)
	// ResTagStateDelete .
	ResTagStateDelete = int32(1)

	// ResTagRoleUp .
	ResTagRoleUp = "up"
	// ResTagRoleUser .
	ResTagRoleUser = "user"
	// ResTagRoleAdmin .
	ResTagRoleAdmin = "admin"

	// RoleUp .
	RoleUp = int32(0)
	// RoleUser .
	RoleUser = int32(1)
	// RoleAdmin .
	RoleAdmin = int32(2)

	// BusinessStateNormal .
	BusinessStateNormal = int32(0)

	// IsChannel tag is channel
	IsChannel = int32(1)
)

// Message .
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Archive .
type Archive struct {
	ID        int64  `json:"id"`
	Aid       int64  `json:"aid"`
	Mid       int64  `json:"mid"`
	TypeID    int16  `json:"typeid"`
	HumanRank int    `json:"humanrank"`
	Duration  int    `json:"duration"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Content   string `json:"content"`
	Tag       string `json:"tag"`
	Attribute int32  `json:"attribute"`
	Copyright int8   `json:"copyright"`
	AreaLimit int8   `json:"arealimit"`
	State     int    `json:"state"`
	Author    string `json:"author"`
	Access    int    `json:"access"`
	Forward   int    `json:"forward"`
	PubTime   string `json:"pubtime"`
	Round     int8   `json:"round"`
	CTime     string `json:"ctime"`
	MTime     string `json:"mtime"`
}

// PlatformTagInfo .
type PlatformTagInfo struct {
	ID           int64  `json:"id"`
	Name         string `json:"title"`
	Cover        string `json:"cover"`
	HeadCover    string `json:"banner"`
	Sub          int64  `json:"use_count"`
	Bind         int64  `json:"bgmcount"`
	Channel      int32  `json:"verify_state"`
	ShortContent string `json:"brief"`
	CommonList   string `json:"commonlist"`
}

// ResTagMessage .
type ResTagMessage struct {
	Oid    int64    `json:"oid"`
	Type   string   `json:"type"`
	Tids   []int64  `json:"tids"`
	TNames []string `json:"tnames"`
	Mid    int64    `json:"mid"`
	Action string   `json:"action"`
	Role   string   `json:"role"`
	MTime  int64    `json:"mtime"`
	Appkey string   `json:"appkey"`
}

// ResTag .
type ResTag struct {
	Oid   int64
	Type  int32
	Mid   int64
	Role  int32
	State int32
	MTime int64
	Tids  []int64
}

// Tag tag.
type Tag struct {
	ID     int64
	Name   string
	Type   int32
	Verify int32
	Attr   int32
	State  int32
}

// Business  tag business.
type Business struct {
	Appkey string
	Remark string
	Name   string
	Type   int32
	State  int32
	Alias  string
}

// ChannelRule channel rule.
type ChannelRule struct {
	ID        int64
	Tid       int64
	ATid      int64
	BTid      int64
	InRule    string
	NotInRule string
	Rule      string
}
