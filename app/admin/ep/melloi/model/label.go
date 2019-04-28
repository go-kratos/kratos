package model

//Label label
type Label struct {
	ID          int64  `json:"id" form:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	Name        string `json:"label_name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Active      int    `json:"active"`
}

//LabelName db table name for label
func (l Label) LabelName() string {
	return "label"
}

//LabelRelation label relation
type LabelRelation struct {
	ID          int64  `json:"id" form:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	LabelID     int64  `json:"label_id" form:"label_id"`
	LabelName   string `json:"label_name" form:"label_name"`
	Color       string `json:"color" form:"color"`
	Description string `json:"description" form:"description"`
	TargetID    int64  `json:"target_id" form:"target_id"`
	Type        int    `json:"type"`
	Active      int    `json:"active"`
}

// LabelRelation type const
const (
	DefaultType = iota
	ScriptType
	ReportType
)

//LabelRelationName db table name of label relation
func (l LabelRelation) LabelRelationName() string {
	return "label_relation"
}
