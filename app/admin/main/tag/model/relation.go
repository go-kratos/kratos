package model

import "go-common/library/time"

// const const value.
const (
	ResTypeArticle  = int32(1) // 文章类型
	ResTypeMusic    = int32(2) // 音乐类型
	ResTypeArchive  = int32(3) // 稿件类型
	ResTypeOpenMall = int32(4) // 开放平台电商

	ResRoleALL   = int32(-1) // ALL
	ResRoleUpper = int32(0)  // up主
	ResRoleUser  = int32(1)  // 用户
	ResRoleAdmin = int32(2)  // 管理员

	ResTagALL     = int32(-1) //全部
	ResTagAdd     = int32(0)  //增加
	ResTagDel     = int32(1)  //删除
	ResTagRestore = int32(2)  //恢复

	RelationStateNormal = int32(0)
	RelationStateDelete = int32(1)

	AttrLockNone = int32(0) // 未锁定
	AttrLocked   = int32(1) // 已锁定

	QueryTypeTName = int32(0) //按照tag名称查询
	QueryTypeOid   = int32(1) //按照资源ID查询

	ArchiveStateOpen        = int32(0)
	ArchiveStateForbidFixed = int32(-6)
)

// Relation res-tag and tag-res.
type Relation struct {
	ID    int64 `json:"id"`
	Oid   int64 `json:"oid"`
	Type  int32 `json:"type"`
	Tid   int64 `json:"tid"`
	Mid   int64 `json:"mid"`
	Role  int32 `json:"role"`
	Enjoy int64 `json:"like"`
	Hate  int64 `json:"hate"`
	Attr  int32 `json:"attr"`
	State int32 `json:"state"`
}

// Resource Resource tag info
type Resource struct {
	ID       int64     `json:"id"`
	Oid      int64     `json:"oid"`
	Type     int64     `json:"type"`
	Title    string    `json:"title"`
	Tid      int64     `json:"tid"`
	Mid      int64     `json:"mid"`
	Author   string    `json:"author"`
	Role     int64     `json:"role"`
	Enjoy    int64     `json:"like"`
	Hate     int64     `json:"hate"`
	Attr     int64     `json:"attr"`
	State    int64     `json:"state"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
	Tag      *Tag      `json:"tag"`
	TagCount *TagCount `json:"tag_count"`
}

// SearchRes seach res.
type SearchRes struct {
	ID        int64  `json:"id"`
	Oid       int64  `json:"oid"`
	Type      int64  `json:"type"`
	TypeID    int64  `json:"typeid"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	MissionID int64  `json:"mission_id"`
	Mid       int64  `json:"mid"`
	PubDate   string `json:"pubtime"`
	CTime     string `json:"ctime"`
	Copyright int32  `json:"copyright"`
	State     int32  `json:"state"`
}

// IsNormal archive state.
func (t *SearchRes) IsNormal() bool {
	return t.State >= ArchiveStateOpen || t.State == ArchiveStateForbidFixed
}

// UserInfo user info.
type UserInfo struct {
	Mid  int64
	Name string
}
