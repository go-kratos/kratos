package model

// Version Cookie中保存的版本信息
type Version struct {
	VersionID int64       `json:"v"`
	Data      map[int]int `json:"d"`
}
