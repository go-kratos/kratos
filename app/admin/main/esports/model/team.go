package model

// Team .
type Team struct {
	ID         int64  `json:"id" form:"id"`
	Title      string `json:"title" form:"title" validate:"required"`
	SubTitle   string `json:"sub_title" form:"sub_title"`
	ETitle     string `json:"e_title" form:"e_title"`
	CreateTime int64  `json:"create_time" form:"create_time"`
	Area       string `json:"area" form:"area"`
	Logo       string `json:"logo" form:"logo" validate:"required"`
	UID        int64  `json:"uid" form:"uid" gorm:"column:uid"`
	Members    string `json:"members" form:"members"`
	Dic        string `json:"dic" form:"dic"`
	IsDeleted  int    `json:"is_deleted" form:"is_deleted"`
}

// TeamInfo .
type TeamInfo struct {
	*Team
	Games []*Game `json:"games"`
}

// TableName .
func (t Team) TableName() string {
	return "es_teams"
}
