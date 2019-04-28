package model

const (
	// BusinessArchiveComplain 稿件投诉
	BusinessArchiveComplain = int8(1)
	// BusinessArchiveAppeal 稿件申诉
	BusinessArchiveAppeal = int8(2)
	// BusinessBlackedAppeal 小黑屋申诉
	BusinessBlackedAppeal = int8(5)
	// BusinessAudit 稿件审核
	BusinessAudit = int8(6)

	// Disbaled 禁用
	Disbaled = int8(0)
	// Enabled 启用
	Enabled = int8(1)
)

// Business struct
type Business struct {
	ID       int32  `gorm:"column:id" json:"id"`
	Cid      int32  `gorm:"column:cid" json:"cid"`
	Oid      int64  `gorm:"column:oid" json:"oid"`
	Business int8   `gorm:"column:business" json:"business"`
	Typeid   int16  `gorm:"column:typeid" json:"business_typeid"`
	Mid      int64  `gorm:"column:mid" json:"business_mid"`
	Title    string `gorm:"column:title" json:"business_title"`
	Content  string `gorm:"column:content" json:"business_content"`
	Extra    string `gorm:"column:extra" json:"business_extra"`
}

// TableName by business
func (*Business) TableName() string {
	return "workflow_business"
}
