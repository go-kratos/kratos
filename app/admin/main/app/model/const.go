package model

import (
	"time"
)

// NoticeChangeTime Change Time
func NoticeChangeTime(timeStr string) (theTime time.Time, err error) {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation(timeLayout, timeStr, loc)
}

// IsCoverType file type  is image
func IsCoverType(fileType string) bool {
	return fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/webp"
}
