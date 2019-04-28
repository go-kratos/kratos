package user

import (
	"go-common/library/time"
)

// TableName case tablename
func (*User) TableName() string {
	return "user"
}

// User user model
type User struct {
	ID        int64     `gorm:"column:id" json:"id" params:"id;Min(1)"`
	UserName  string    `gorm:"column:username" json:"username" params:"username"`
	NickName  string    `gorm:"column:nickname" json:"nickname" params:"nickname"`
	Email     string    `gorm:"column:email" json:"email" params:"email"`
	Phone     string    `gorm:"column:phone" json:"phone" params:"phone"`
	Status    int8      `gorm:"column:status" json:"status" params:"status"`
	AvatarURL string    `gorm:"-" json:"avatar_url"`
	Ctime     time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime     time.Time `gorm:"column:mtime" json:"-"`
}

// Pager user pager
type Pager struct {
	Total int64   `json:"total"`
	Pn    int     `params:"pn" default:"1"`
	Ps    int     `params:"ps" default:"20"`
	Items []*User `json:"items"`
}

// Result contains user and modules and rules.
type Result struct {
	Super bool     `json:"superman"`
	Env   string   `json:"env"`
	User  *User    `json:"user"`
	Rules []string `json:"rules"`
}
