package databus

// TableName case tablename
func (*Filter) TableName() string {
	return "filters"
}

// Filter apply model
type Filter struct {
	ID        int    `gorm:"column:id" json:"id"`
	Nid       int    `gorm:"column:nid" json:"nid"`
	Filters   string `gorm:"column:filters" json:"-"`
	Field     string `gorm:"-" json:"field"`
	Condition int8   `gorm:"-" json:"condition"`
	Value     string `gorm:"-" json:"value"`
}
