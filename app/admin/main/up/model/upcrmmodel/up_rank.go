package upcrmmodel

import "go-common/library/time"

const (
	//UpRankTypeFans30day1k 7
	UpRankTypeFans30day1k = 7 // 第一次投稿时间在30天内的UP主最快达到1k粉的top500UP列表
	//UpRankTypeFans30day1w 8
	UpRankTypeFans30day1w = 8 // 第一次投稿时间在30天内的UP主最快达到1w粉的top500UP列表
	//UpRankTypePlay30day1k 9
	UpRankTypePlay30day1k = 9 // 第一次投稿时间在30天内的UP主最快达到1k播放量的top500UP列表
	//UpRankTypePlay30day1w 10
	UpRankTypePlay30day1w = 10 // 第一次投稿时间在30天内的UP主最快达到1w播放量的top500UP列表
	//UpRankTypePlay30day10k 11
	UpRankTypePlay30day10k = 11 // 第一次投稿时间在30天内的UP主最快达到10w播放量的top500UP列表
	//UpRankTypeFans30dayIncreaseCount 12
	UpRankTypeFans30dayIncreaseCount = 12 //30天内粉丝增长绝对值最多的top500列表
	//UpRankTypeFans30dayIncreasePercent 13
	UpRankTypeFans30dayIncreasePercent = 13 // 30天内粉丝增长百分比最多（30天前粉丝量超过100）的top500列表
)

//UpRank  struct
type UpRank struct {
	ID           uint64 `gorm:"column:id"`
	Mid          int64
	Type         int16 // 排行榜类型
	Value        uint
	Value2       int
	GenerateDate string
	CTime        time.Time `gorm:"column:ctime"`
	MTime        time.Time `gorm:"column:mtime"`
}
