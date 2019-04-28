package tag

// ActivityInfo activity info
type ActivityInfo struct {
	ActivityID int64  `json:"mission_id"`
	AvID       int64  `json:"id"`
	MID        int64  `json:"mid"`
	TypeID     int64  `json:"typeid"`
	CDate      string `json:"cdate"`
}

// ColumnAct column activity ret
type ColumnAct struct {
	List []*CmActInfo `json:"list"`
}

// CmActInfo arcitle info
type CmActInfo struct {
	ID       int64      `json:"id"`
	SID      int64      `json:"sid"`
	MID      int64      `json:"mid"`
	CTime    string     `json:"ctime"`
	Category CmCategory `json:"category"`
}

// CmCategory column category
type CmCategory struct {
	ID int64 `json:"id"`
}

// TypesInfo category info.
type TypesInfo struct {
	PID int64 `json:"pid"`
	ID  int64 `json:"id"`
}

// ColumnType column type
type ColumnType struct {
	ID       int64        `json:"id"`
	ParentID int64        `json:"parent_id"`
	Children []ColumnType `json:"children"`
}
