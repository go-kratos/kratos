package model

//ConfValues 配置
type ConfValues map[int64]string

const (
	// Empty 空
	Empty = ""
	// Set 存在
	Set = "1"
	// GoldTarget 金瓜子
	GoldTarget = 1010
	// SilverTarget 银瓜子
	SilverTarget = 1011
)

// IsSet 是否存在
func (vs ConfValues) IsSet(key int64) bool {
	v, ok := vs[key]
	if !ok {
		return false
	}
	return v == Set
}
