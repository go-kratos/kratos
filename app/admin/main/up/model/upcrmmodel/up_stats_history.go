package upcrmmodel

import "go-common/library/time"

// UpStatsHistory is table up_stats_history mapping
type UpStatsHistory struct {
	ID           uint32    `gorm:"column:id"`
	Type         int       `gorm:"column:type"`
	SubType      int       `gorm:"column:sub_type"`
	Value1       int64     `gorm:"column:value1"`
	Value2       int64     `gorm:"column:value2"`
	GenerateDate time.Time `gorm:"column:generate_date"`
	Ctime        time.Time `gorm:"column:ctime"`
	Mtime        time.Time `gorm:"column:mtime"`
}

const (
	//ActivityType 活跃Up主数量
	ActivityType = iota + 1
	//IncrType 新增
	IncrType
	//TotalType 总数
	TotalType
)

//Activity 活跃度
type Activity int

const (
	//HighActivity 1
	HighActivity Activity = iota + 1
	//MediumActivity 2
	MediumActivity
	//LowActivity 3
	LowActivity
	//LostActivity 4
	LostActivity
)
