package model

// User user model for login
type User struct {
	ID     int64  `json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Accept int32  `json:"accept"`
	Active string `json:"active"`
}

// TableName get user model name
func (u User) TableName() string {
	return "user"
}
