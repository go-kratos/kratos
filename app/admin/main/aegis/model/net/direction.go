package net

import (
	"time"
)

const (
	// TableDirection .
	TableDirection = "net_direction"
	//DirInput 从flow指向transition
	DirInput = int8(1)
	//DirOutput 从transition指向flow
	DirOutput = int8(2)

	//DirOrderSequence 下游顺序执行
	DirOrderSequence = int8(0)
	//DirOrderOrSplit 下游根据条件拆分，各分支若为transition,可允许操作个数>=1
	DirOrderOrSplit = int8(1)
	//DirOrderOrResultSplit 下游transition操作拆分，每个transition只有一个可允许操作，enable均默认为true，但只能操作一个
	//先不做, 如果做了，在详情页提交后，怎么知道提交的是哪个transition（有多个）呢？
	DirOrderOrResultSplit = int8(2)
	//todo --- 其他顺序 v2
)

// DirDirectionDesc .
var DirDirectionDesc = map[int8]string{
	DirInput:  "从节点指向变化",
	DirOutput: "从变化指向节点",
}

// DirOrderDesc 有向线下游顺序描述
var DirOrderDesc = map[int8]string{
	DirOrderSequence: "直序",
}

// Direction 有向线，连接flow和transition
type Direction struct {
	ID           int64     `gorm:"primary_key" json:"id" form:"id" validate:"omitempty,gt=0"`
	NetID        int64     `gorm:"column:net_id" json:"net_id" form:"net_id" validate:"omitempty,gt=0"`
	FlowID       int64     `gorm:"column:flow_id" json:"flow_id" form:"flow_id" validate:"required,gt=0"`
	TransitionID int64     `gorm:"column:transition_id" json:"transition_id" form:"transition_id" validate:"required,gt=0"`
	Direction    int8      `gorm:"column:direction" json:"direction" form:"direction" validate:"required,min=1,max=2"`
	Order        int8      `gorm:"column:order" json:"order" form:"order" validate:"omitempty,min=0,max=2"`
	Guard        string    `gorm:"column:guard" json:"guard"`
	Output       string    `gorm:"column:output" json:"output"`
	UID          int64     `gorm:"column:uid" json:"uid"`
	DisableTime  time.Time `gorm:"column:disable_time" json:"disable_time"`
	Ctime        time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime        time.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName .
func (d *Direction) TableName() string {
	return TableDirection
}

// IsAvailable .
func (d *Direction) IsAvailable() bool {
	return d.DisableTime.IsZero()
}
