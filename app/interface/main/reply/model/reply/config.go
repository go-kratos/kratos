package reply

// Config reply config info.
type Config struct {
	ID        int64  `json:"id"`
	Type      int8   `json:"type"`
	Oid       int64  `json:"oid"`
	AdminID   int64  `json:"adminid"`
	Operator  string `json:"operator"`
	Category  int8   `json:"category"`
	Config    string `json:"config"`
	ShowEntry int8   `json:"showentry"`
	ShowAdmin int8   `json:"showadmin"`
	CTime     string `json:"ctime"`
	MTime     string `json:"mtime"`
}
