package manager

import "encoding/json"

const (
	UpTypeExcitationWhite = 19 //激励回查白名单分组
	TableUps              = "ups"
)

//User user info
type User struct {
	ID         int64  `json:"uid"`
	Username   string `json:"username"`
	Department string `json:"department"`
}

// BinMsg manager binlog消息结构
type BinMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Ups UP主分组关联结构
type Ups struct {
	ID   int64 `json:"id"`
	MID  int64 `json:"mid"`
	Type int64 `json:"type"`
	//Note string `json:"note"`
	//CTime string `json:"ctime"`
	//MTime string `json:"mtime"`
}

// UpGroup UP主分组关联关系
type UpGroup struct {
	ID        int64  `json:"id"`
	MID       int64  `json:"mid"`
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
}
