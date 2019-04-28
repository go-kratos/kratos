package model

import (
	"encoding/json"
)

// consts
const (
	ActUpdateExp      = "updateExp"
	ActUpdateLevel    = "updateLevel"
	ActUpdateFace     = "updateFace"
	ActUpdateMoral    = "updateMoral"
	ActUpdateUname    = "updateUname"
	ActUpdateRealname = "updateRealname"
	ActUpdateByAdmin  = "updateByAdmin"
	ActBlockUser      = "blockUser"
)

// consts
const (
	CacheKeyBase  = "bs_%d"    // key of baseInfo
	CacheKeyMoral = "moral_%d" // key of detail
	CacheKeyInfo  = "i_"
)

// Binlog is
type Binlog struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// NewExp userexp for mysql scan.
type NewExp struct {
	Mid  int64 `json:"mid"`
	Exp  int64 `json:"exp"`
	Flag int32 `json:"flag"`
}

// NeastMid is
type NeastMid struct {
	Mid int64 `json:"mid"`
}

// NotifyInfo notify info.
type NotifyInfo struct {
	Uname   string `json:"uname"`
	Mid     int64  `json:"mid"`
	Type    string `json:"type"`
	NewName string `json:"newName"`
	Action  string `json:"action"`
}

// ExpMessage exp msg
type ExpMessage struct {
	Mid int64 `json:"mid"`
	Exp int64 `json:"exp"`
}

// MemberBase is
type MemberBase struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
	Sex  int64  `json:"sex"`
	Face string `json:"face"`
	Sign string `json:"sign"`
	Rank int64  `json:"rank"`
}
