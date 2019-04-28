package model

import (
	"strconv"

	"go-common/library/log"
)

// const .
const (
	// pendant status
	PendantStatusON  = 1
	PendantStatusOFF = 0

	// group status
	GroupStatusON  = 1
	GroupStatusOFF = 0

	// packpage status
	InvalidPendantPKG = int32(0)
	ValidPendantPKG   = int32(1)
	EquipPendantPKG   = int32(2)

	// pendant equip
	PendantEquipOFF = int8(1)
	PendantEquipON  = int8(2)

	// pendant source
	UnknownEquipSource = 0
	EquipFromPackage   = 1
	EquipFromVIP       = 2
)

// IsValidSource 挂件来源是否合法  合法：true,无效：false
func IsValidSource(source int64) bool {
	if source != EquipFromPackage && source != EquipFromVIP && source != UnknownEquipSource {
		log.Error("IsValidSource souce=%v is not correct value", source)
		return false
	}
	return true
}

// ParseSource c处理挂件来源
func ParseSource(sourceStr string) int64 {
	// 没有传值，则设置为未知挂件
	if sourceStr == "" {
		return UnknownEquipSource
	}
	// 有传递参数,但是没有按照要求传值，也设置为未知挂件
	source, err := strconv.ParseInt(sourceStr, 10, 64)
	if err != nil {
		log.Error("ParseSource err(%+v)", err)
		return UnknownEquipSource
	}
	// 没有按照要求传值，也设置为未知挂件
	if source != EquipFromPackage && source != EquipFromVIP && source != UnknownEquipSource {
		log.Error("ParseSource souce=%v is not correct value", source)
		return UnknownEquipSource
	}
	return source
}
