package operation

// Operation tool.
type Operation struct {
	ID       int64  `form:"id" json:"id"`
	Type     string `form:"type" json:"type"`
	Ads      int8   `form:"ads" json:"ads"`
	Platform int8   `form:"platform" json:"platform"`
	Rank     int8   `form:"rank" json:"rank"`
	Pic      string `form:"pic" json:"pic"`
	Link     string `form:"link" json:"link"`
	Content  string `form:"content" json:"content"`
	Username string `form:"username" json:"username"`
	Remark   string `form:"remark" json:"remark"`
	Note     string `form:"note" json:"note"`
	AppPic   string `form:"app_pic" json:"app_pic"`
	Stime    string `form:"stime" json:"stime" gorm:"column:stime"`
	Etime    string `form:"etime" json:"etime" gorm:"column:etime"`
	Ctime    string `form:"ctime" json:"ctime" gorm:"column:ctime"`
	Mtime    string `form:"mtime" json:"mtime" gorm:"column:mtime"`
	Dtime    string `form:"dtime" json:"dtime" gorm:"column:dtime"`
}

// TableName fn
func (Operation) TableName() string {
	return "operations"
}

// Banner for app index.
type Banner struct {
	Ty      string `json:"-"`
	Rank    string `json:"rank"`
	Pic     string `json:"pic"`
	Link    string `json:"link"`
	Content string `json:"content"`
}

// ViewOperation tool.
type ViewOperation struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Ads      int8   `json:"ads"`
	Platform int8   `json:"platform"`
	Rank     int8   `json:"rank"`
	Pic      string `json:"pic"`
	Link     string `json:"link"`
	Content  string `json:"content"`
	Username string `json:"username"`
	Remark   string `json:"remark"`
	Note     string `json:"note"`
	AppPic   string `json:"app_pic"`
	Stime    string `json:"stime"`
	Etime    string `json:"etime"`
	Ctime    string `json:"ctime"`
	Mtime    string `json:"mtime"`
	Dtime    string `json:"dtime"`
	Status   string `json:"status"`
}
