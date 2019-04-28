package bottom

// Bottom bottom
type Bottom struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Rank   int64  `json:"rank"`
	Action int    `json:"action"`
	Param  string `json:"param"`
	State  int    `json:"state"`
}

// Param param
type Param struct {
	ID     int64  `form:"id"`
	Name   string `form:"name"`
	Logo   string `form:"logo"`
	Rank   int64  `form:"rank"`
	Action int    `form:"action"`
	Param  string `form:"param"`
	State  int    `form:"state"`
	IDs    string `form:"ids"`
}

// TableName return table name
func (*Bottom) TableName() string {
	return "bottom_entry"
}
