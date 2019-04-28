package model

// Attachment struct
type Attachment struct {
	ID   int32  `gorm:"column:id" json:"id"`
	Cid  int32  `gorm:"column:cid" json:"cid"`
	Path string `gorm:"column:path" json:"path"`
}

// TableName by Attachment
func (*Attachment) TableName() string {
	return "workflow_attachment"
}
