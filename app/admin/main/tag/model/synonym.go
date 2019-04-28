package model

import "go-common/library/time"

// SynonymTag SynonymTag.
type SynonymTag struct {
	ID     int64       `json:"id"` //主键
	Ptid   int64       `json:"ptid"`
	Tid    int64       `json:"tid"`
	TName  string      `json:"tname"`
	UName  string      `json:"uname"`
	Adverb []*BasicTag `json:"adverb"` //副词
	CTime  time.Time   `json:"ctime"`
	MTime  time.Time   `json:"mtime"`
}

// SynonymInfo SynonymInfo.
type SynonymInfo struct {
	Tag    *Tag   `json:"tag"`
	Adverb []*Tag `json:"adverb"`
}
