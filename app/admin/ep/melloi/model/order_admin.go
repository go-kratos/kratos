package model

// OrderAdmin perf administrator model for perf order
type OrderAdmin struct {
	ID       int64  `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	UserName string `json:"user_name" form:"user_name"`
}

// TableName get table name
func (w OrderAdmin) TableName() string {
	return "order_admin"
}
