package notice

import xtime "go-common/library/time"

// Notice is notice type.
type Notice struct {
	ID        int        `json:"id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Content   string     `json:"content,omitempty"`
	Start     xtime.Time `json:"start_time,omitempty"`
	End       xtime.Time `json:"end_time,omitempty"`
	URI       string     `json:"uri,omitempty"`
	Type      int        `json:"-"`
	Plat      int8       `json:"-"`
	Build     int        `json:"-"`
	Condition string     `json:"-"`
	Area      string     `json:"-"`
}
