package oplog

import (
	"go-common/library/log"
)

// Infoc operation log for administrator
type Infoc struct {
	Oid           int64        `json:"oid"`
	Type          int          `json:"type"`
	DMIds         []int64      `json:"dmids"`
	Subject       string       `json:"subject"`
	OriginVal     string       `json:"origin_val"`
	CurrentVal    string       `json:"current_val"`
	OperationTime string       `json:"optime"`
	OperatorType  OperatorType `json:"operator_type"`
	Operator      int64        `json:"operator"`
	Source        Source       `json:"source"`
	Remark        string       `json:"remark"`
}

// Source enum integer value
type Source int

// Source enum definition list
const (
	_ Source = iota
	SourceManager
	SourceUp
	SourcePlayer
)

// String returns the Source enmu description
func (source Source) String() string {
	var text string
	switch source {
	case SourceManager:
		text = "运营后台"
	case SourceUp:
		text = "创作中心"
	case SourcePlayer:
		text = "播放器"
	default:
		log.Warn("String() Unknow Source, warn(%v)", int(source))
		text = "未知来源"
	}
	return text
}

// OperatorType enum integer value
type OperatorType int

// OperatorType enum definition list
const (
	_ OperatorType = iota
	OperatorAdmin
	OperatorMember
	OperatorUp
	OperatorSystem
)

// String returns the Source enmu description
func (opType OperatorType) String() string {
	var text string
	switch opType {
	case OperatorAdmin:
		text = "管理员"
	case OperatorMember:
		text = "用户"
	case OperatorUp:
		text = "UP主"
	case OperatorSystem:
		text = "系统"
	default:
		log.Warn("String() Unknow Source, warn(%v)", int(opType))
		text = "未知来源"
	}
	return text
}
