package archive

import "time"

// StatsPoint points struct for stats
type StatsPoint struct {
	ID      int64     `json:"id"`
	Type    int8      `json:"type"`
	Content string    `json:"content"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// ArcStayCount .
type ArcStayCount struct {
	Round int8  `json:"round"`
	State int8  `json:"state"`
	Count int64 `json:"count"`
}
