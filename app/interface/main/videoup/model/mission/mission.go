package mission

import "time"

type Mission struct {
	ID    int       `json:"id"`
	Name  string    `json:"name"`
	Tags  string    `json:"tags"`
	ETime time.Time `json:"etime"`
}

// ActInfo act proctocol & tag.
type ActInfo struct {
	ID       string `json:"id"`
	SID      string `json:"sid"`
	Protocol string `json:"protocol"`
	Types    string `json:"types"`
	Tag      string `json:"tags"`
	CTime    string `json:"ctime"`
	MTime    string `json:"mtime"`
}
