package model

// Archive .
type Archive struct {
	ID      int64  `json:"id" gorm:"column:id"`
	MID     int64  `json:"mid" gorm:"column:mid"`
	TypeID  int32  `json:"typeid" gorm:"column:typeid"`
	Title   string `json:"title" gorm:"column:title"`
	Content string `json:"content" gorm:"column:content"`
}

// TableName .
func (a *Archive) TableName() string {
	return "archive"
}
