package model

//Apply apply model
type Apply struct {
	ID        int64  `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	Path      string `json:"path"`
	From      string `json:"from" form:"from"`
	To        string `json:"to" form:"to"`
	Status    int32  `json:"status" form:"status"`
	StartTime string `json:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" form:"end_time"`
	Active    int32  `json:"active"`
}

//const definition
const (
	ApplyValid   = 1  // active=1 有效
	ApplyInvalid = -1 // active=-1 无效
)

// QueryApplyResponse response model
type QueryApplyResponse struct {
	ApplyList []*Apply `json:"apply_list"`
	Pagination
}

// QueryApplyRequest request model
type QueryApplyRequest struct {
	Apply
	Pagination
}

// TableName get table name model
func (w Apply) TableName() string {
	return "apply"
}
