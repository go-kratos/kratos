package show

import "go-common/app/admin/main/feed/model/common"

//ChannelTab channel tab
type ChannelTab struct {
	ID       int64  `json:"id" form:"id"`
	TagID    int64  `json:"tag_id" form:"tag_id" validate:"required"`
	TabID    int64  `json:"tab_id" form:"tab_id" validate:"required"`
	Title    string `json:"title" form:"title" validate:"required"`
	Stime    int64  `json:"stime" form:"stime" validate:"required"`
	Etime    int64  `json:"etime" form:"etime" validate:"required"`
	Check    int    `json:"check" form:"check"`
	Priority int    `json:"priority" form:"priority" validate:"required"`
	UID      int64  `json:"uid" form:"uid"`
	Person   string `json:"person" form:"person"`
	IsDelete int    `json:"is_delete" form:"is_delete"`
	Status   int    `json:"status" form:"status"`
}

//ChannelTabPager .
type ChannelTabPager struct {
	Item []*ChannelTab `json:"item"`
	Page common.Page   `json:"page"`
}

// TableName .
func (a ChannelTab) TableName() string {
	return "channel_tab"
}

/*
---------------------------
 struct param
---------------------------
*/

//ChannelTabAP channel tab add param
type ChannelTabAP struct {
	TagID    int64  `form:"tag_id" validate:"required"`
	TabID    int64  `form:"tab_id" validate:"required"`
	Title    string `form:"title" validate:"required"`
	Stime    int64  `form:"stime" validate:"required"`
	Etime    int64  `form:"etime" validate:"required"`
	Priority int    `form:"priority" validate:"required"`
	UID      int64  `form:"uid" gorm:"column:uid"`
	Person   string `form:"person"`
}

//ChannelTabUP channel tab update param
type ChannelTabUP struct {
	ID       int64  `form:"id" validate:"required"`
	TagID    int64  `form:"tag_id" validate:"required"`
	TabID    int64  `form:"tab_id" validate:"required"`
	Title    string `form:"title" validate:"required"`
	Stime    int64  `form:"stime" validate:"required"`
	Etime    int64  `form:"etime" validate:"required"`
	Priority int    `form:"priority" validate:"required"`
	UID      int64  `form:"uid" gorm:"column:uid"`
	Person   string `form:"person"`
}

//ChannelTabLP channel tab list param
type ChannelTabLP struct {
	TagID  int    `form:"tag_id"`
	TabID  int    `form:"tab_id"`
	Stime  int64  `form:"stime"`
	Etime  int64  `form:"etime"`
	Status int    `form:"status"`
	Person string `form:"person"`
	Order  int    `form:"order" default:"2"`
	Ps     int    `form:"ps" default:"20"` // 分页大小
	Pn     int    `form:"pn" default:"1"`  // 第几个分页
}

// TableName .
func (a ChannelTabAP) TableName() string {
	return "channel_tab"
}

// TableName .
func (a ChannelTabUP) TableName() string {
	return "channel_tab"
}
