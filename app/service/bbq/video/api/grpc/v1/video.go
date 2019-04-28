package v1

// 表示video表中的limit字段，每一位代表的含义
// limit中0表示通过，1表示禁止
const (
	VideoLimitBitAll    uint64 = iota
	VideoLimitBitBullet        // 表示弹幕不显示，发布失败
	VideoLimitBitMax
)

// IsLimitSet 根据输入的limits值，返回相应limitType是否被设置
// example: 例如需要判断该视频是否可以发布弹幕的时候，如果(sql_limit_value & (1<<VideoLimitBitBullet)) > 0 则该视频被禁止发布弹幕
func IsLimitSet(limits uint64, limitType uint64) bool {
	return limits&(1<<limitType) > 0
}
