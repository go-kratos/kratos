package model

import "time"

// PushMsg is used to push to the mobile clients to indicate the mod name to request
type PushMsg struct {
	ResID   int    `json:"res_id" gorm:"column:id"`
	ModID   int    `json:"mod_id" gorm:"column:pool_id"`
	ModName string `json:"mod_name" gorm:"column:name"`
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
	FileType   int8      `json:"file_type"`
	FromVer    int64     `json:"from_ver"`
	IsDeleted  int8      `json:"is_deleted"`
}

//FileInfo : the uploaded file information
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
	Md5  string `json:"md5"`
	URL  string `json:"url"`
}

// Resource reprensents the resource table
type Resource struct {
	ID      int64     `json:"id" params:"id"`
	Name    string    `json:"name" params:"name"`
	Version int64     `json:"version" params:"version"`
	PoolID  int64     `json:"pool_id" params:"pool_id"`
	Ctime   time.Time `json:"ctime" params:"ctime"`
	Mtime   time.Time `json:"mtime" params:"mtime"`
}
