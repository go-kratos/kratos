package model

const (
	// ConfigCategoryLogDel ConfigCategoryLogDel
	ConfigCategoryLogDel int32 = 1
)

// Config reply config info.
type Config struct {
	ID        int64  `json:"id"`
	Type      int32  `json:"type"`
	Oid       int64  `json:"oid"`
	AdminID   int64  `json:"adminid"`
	Operator  string `json:"operator"`
	Category  int32  `json:"category"`
	Config    string `json:"config"`
	ShowEntry int32  `json:"showentry"`
	ShowAdmin int32  `json:"showadmin"`
	CTime     string `json:"ctime"`
	MTime     string `json:"mtime"`
}
