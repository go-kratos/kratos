package notice

import "time"

// Notice notice
type Notice struct {
	ID         int64     `json:"id"`
	Plat       int       `json:"plat"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	URL        string    `json:"url"`
	Eftime     time.Time `json:"ef_time"`
	Extime     time.Time `json:"ex_time"`
	Build      int       `json:"build"`
	Conditions string    `json:"conditions"`
	Area       string    `json:"area"`
	State      int       `json:"state"`
	Type       int       `json:"type"`
}

// Param param
type Param struct {
	ID         int64  `form:"id"`
	Plat       int    `form:"plat"`
	Title      string `form:"title"`
	Content    string `form:"content"`
	URL        string `form:"url"`
	EftimeStr  string `form:"ef_time"`
	ExtimeStr  string `form:"ex_time"`
	Eftime     time.Time
	Extime     time.Time
	Build      int    `form:"build"`
	Conditions string `form:"conditions"`
	Area       string `form:"area"`
	State      int    `form:"state"`
	Type       int    `form:"type"`
}

// TableName return table name
func (*Notice) TableName() string {
	return "notice"
}
