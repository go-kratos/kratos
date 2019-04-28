package language

// Language language
type Language struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
}

// Param param
type Param struct {
	ID     int64  `form:"id"`
	Name   string `form:"name"`
	Remark string `form:"remark"`
}

// TableName return table name
func (*Language) TableName() string {
	return "language"
}
