package model

import "time"

//Comment model for performance test job comment
type Comment struct {
	ID         int       `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ReportID   int       `json:"report_id" form:"report_id"`
	Content    string    `json:"content" gorm:"content"`
	UserName   string    `json:"user_name" form:"user_name" gorm:"user_name"`
	Status     int       `json:"status" form:"status"`
	SubmitDate time.Time `json:"submit_date"`
}

//QueryCommentResponse model for QueryCommentResponse
type QueryCommentResponse struct {
	Total    int        `json:"total"`
	Comments []*Comment `json:"comment_list"`
}

//TableName db table name of Comment
func (w Comment) TableName() string {
	return "comment"
}
