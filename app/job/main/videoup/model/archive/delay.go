package archive

import "time"

//定时发布类型
const (
	DelayTypeForAdmin = int8(1)
	DelayTypeForUser  = int8(2)
)

//Delay 定时发布结构
type Delay struct {
	ID    int64
	Aid   int64
	DTime time.Time
	Type  int8
	State int8
}
