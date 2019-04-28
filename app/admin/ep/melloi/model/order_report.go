package model

// OrderReport order(manual) report model
type OrderReport struct {
	ID       int32  `json:"id"`
	OrderID  int64  `json:"order_id"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	UpdateBy string `json:"update_by" column:"update_by"`
	Active   int32  `json:"active"`
}

// TableName get table name model
func (w OrderReport) TableName() string {
	return "order_report"
}
