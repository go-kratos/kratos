package wall

// Wall wall
type Wall struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Logo     string `json:"logo"`
	Package  string `json:"package"`
	Size     string `json:"size"`
	Download string `json:"download"`
	Remark   string `json:"remark"`
	Rank     int    `json:"rank"`
	State    int    `json:"state"`
}

// Param param
type Param struct {
	ID       int64  `form:"id"`
	Name     string `form:"name"`
	Title    string `form:"title"`
	Logo     string `form:"logo"`
	Package  string `form:"package"`
	Size     string `form:"size"`
	Download string `form:"download"`
	Remark   string `form:"remark"`
	Rank     int    `form:"rank"`
	State    int    `form:"state"`
	IDs      string `form:"ids"`
}

// TableName return table name
func (*Wall) TableName() string {
	return "wall"
}
