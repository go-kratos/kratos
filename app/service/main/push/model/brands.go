package model

import "strings"

const (
	_brandOhters  = 0
	_brandXiaomi  = 1
	_brandHuawei  = 2
	_brandOppo    = 3
	_brandVivo    = 4
	_brandMeizu   = 5
	_brandSamsung = 6
)

// mapping 映射可以解决一个品牌对应多个品牌标识的问题
var brandMapping = map[string]int{
	"xiaomi":  _brandXiaomi,
	"huawei":  _brandHuawei,
	"honor":   _brandHuawei,
	"oppo":    _brandOppo,
	"vivo":    _brandVivo,
	"meizu":   _brandMeizu,
	"samsung": _brandSamsung,
}

// DeviceBrand .
func DeviceBrand(s string) int {
	s = strings.Trim(s, " ")
	if s == "" {
		return _brandOhters
	}
	s = strings.ToLower(s)
	if v, ok := brandMapping[s]; ok {
		return v
	}
	return _brandOhters
}
