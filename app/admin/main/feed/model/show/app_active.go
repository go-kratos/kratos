package show

//AppActive db show app_active table
type AppActive struct {
	Name string `json:"name,omitempty"`
}

// TableName .
func (a AppActive) TableName() string {
	return "app_active"
}
