package model

import "go-common/library/time"

//FileInfo : the uploaded file information
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
	Md5  string `json:"md5"`
	URL  string `json:"url"`
}

// ResourceFile represents the table structure
type ResourceFile struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Md5        string    `json:"md5"`
	Size       int       `json:"size"`
	URL        string    `json:"url"`
	ResourceID int       `json:"resource_id"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}

// TableName gives the table name of the model
func (*ResourceFile) TableName() string {
	return "resource_file"
}
