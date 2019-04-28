package model

const (
	//LogBlacklist blacklist action log type id
	LogBlacklist = 1
)

//Page pager
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}
