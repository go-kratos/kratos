package model

import "time"

// UserResource def
type UserResource struct {
	ID       int32     `json:"id" gorm:"id"`
	ResType  int32     `json:"res_type" form:"res_type"`
	CustomID int32     `json:"custom_id" form:"custom_id"`
	Title    string    `json:"title" form:"title"`
	URL      string    `json:"url" form:"url"`
	Weight   int32     `json:"weight" form:"weight"`
	Status   int32     `json:"status" form:"status"`
	Creator  string    `json:"creator" form:"creator"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// TableName user_resource
func (c UserResource) TableName() string {
	return "user_resource"
}
