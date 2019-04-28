package channel

import (
	"go-common/app/admin/main/feed/model/common"
)

var (
	//NotDelete not delete
	NotDelete uint8
	//Delete delete
	Delete uint8 = 1
	//LogBusPgcsRcmd log business id
	LogBusPgcsRcmd = 203
	//LogBusRcmdNew log business id
	LogBusRcmdNew = 204
	//ActUpCsPgcRcmd log
	ActUpCsPgcRcmd = "ActUpCsPgcRcmd"
	//ActUpCsRcmdNew log
	ActUpCsRcmdNew = "ActUpCsRcmdNew"
	//ActDelCsPgcRcmd log
	ActDelCsPgcRcmd = "ActDelCsPgcRcmd"
	//ActDelCsRcmdNew log
	ActDelCsRcmdNew = "ActDelCsRcmdNew"
	//ActAddCsPgcRcmd log
	ActAddCsPgcRcmd = "ActAddCsPgcRcmd"
	//ActAddCsRcmdNew log
	ActAddCsRcmdNew = "ActAddCsRcmdNew"
)

//AddCardSetup 复合卡片 需要首先单独创建 然后再在频道首页创建
type AddCardSetup struct {
	Type      string `form:"type" validate:"required"`
	Value     string `form:"value"`
	Title     string `form:"title"`
	LongTitle string `form:"longtitle"`
	Content   string `form:"content" validate:"required"`
	UID       int64  `gorm:"column:uid"`
	Person    string
}

//UpdateCardSetup 复合卡片 需要首先单独创建 然后再在频道首页创建
type UpdateCardSetup struct {
	ID        int    `form:"id" validate:"required"`
	Type      string `form:"type" validate:"required"`
	Value     string `form:"value"`
	Title     string `form:"title" validate:"required"`
	LongTitle string `form:"longtitle"`
	Content   string `form:"content" validate:"required"`
}

//Setup 复合卡片 需要首先单独创建 然后再在频道首页创建
type Setup struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Title     string `json:"title"`
	LongTitle string `json:"longtitle"`
	Content   string `json:"content"`
	Deleted   int    `json:"deleted"`
	Person    string `json:"person"`
}

//SetupPager return values
type SetupPager struct {
	Item []*Setup    `json:"item"`
	Page common.Page `json:"page"`
}

// TableName DarkPubLog dark word publish log
func (a Setup) TableName() string {
	return "card_set"
}

// TableName DarkPubLog dark word publish log
func (a AddCardSetup) TableName() string {
	return "card_set"
}
