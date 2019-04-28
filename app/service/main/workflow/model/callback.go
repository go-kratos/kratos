package model

// Callback struct
type Callback struct {
	ID       int32  `gorm:"column:id"`
	URL      string `gorm:"column:url"`
	Business int8   `gorm:"column:business"`
	State    int8   `gorm:"column:state"`
}

// TableName by Callback
func (*Callback) TableName() string {
	return "workflow_callback"
}
