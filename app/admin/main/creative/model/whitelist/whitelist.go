package whitelist

// Whitelist tool.
type Whitelist struct {
	ID           int64  `form:"id" json:"id"`
	MID          int64  `form:"mid" json:"mid" gorm:"column:mid"`
	AdminMID     int64  `form:"admin_mid" json:"admin_mid" gorm:"column:admin_mid"`
	Comment      string `form:"comment" json:"comment"`
	State        int8   `form:"state" json:"state"`
	Type         int8   `form:"type" json:"type"`
	Fans         int64  `form:"fans" json:"fans" gorm:"-"`
	CurrentLevel int32  `form:"current_level" json:"current_level" gorm:"-"`
	Name         string `form:"name" json:"name" gorm:"-"`
	Ctime        string `form:"ctime" json:"ctime" gorm:"column:ctime"`
	Mtime        string `form:"mtime" json:"mtime" gorm:"column:mtime"`
}

// TableName fn
func (Whitelist) TableName() string {
	return "whitelist"
}
