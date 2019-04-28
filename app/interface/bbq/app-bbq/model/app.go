package model

import (
	"go-common/library/time"
)

// AppVersion .
type AppVersion struct {
	ID       int    `json:"id"`
	Platform int    `json:"platform"`
	Name     string `json:"name"`
	Code     int    `json:"version"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Download string `json:"download"`
	MD5      string `json:"md5"`
	Size     int    `json:"file_size"`
	Force    int    `json:"force"`
	Status   int    `json:"status"`
}

// AppResource .
type AppResource struct {
	ID        int       `json:"id"`
	Platform  int       `json:"platform"`
	Name      string    `json:"name"`
	Code      int       `json:"version"`
	Download  string    `json:"download"`
	MD5       string    `json:"md5"`
	Status    int       `json:"status"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// DynamicEffect .
type DynamicEffect struct {
	ID   int    `json:"rid"`
	Name string `json:"rname"`
}
