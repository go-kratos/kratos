package bplus

import xtime "go-common/library/time"

const (
	CLIPS  = 1
	ALBUMS = 2
)

type Clip struct {
	ID       int64      `json:"id,omitempty"`
	Duration int64      `json:"duration,omitempty"`
	CTime    xtime.Time `json:"ctime,omitempty"`
	View     int        `json:"view,omitempty"`
	Damaku   int        `json:"damaku,omitempty"`
	Title    string     `json:"title,omitempty"`
	Cover    string     `json:"cover,omitempty"`
	Tag      string     `json:"tag,omitempty"`
}

type Album struct {
	ID       int64       `json:"doc_id,omitempty"`
	CTime    xtime.Time  `json:"ctime,omitempty"`
	Count    int         `json:"count,omitempty"`
	View     int         `json:"view,omitempty"`
	Comment  int         `json:"comment,omitempty"`
	Title    string      `json:"title,omitempty"`
	Desc     string      `json:"description,omitempty"`
	Pictures []*Pictures `json:"pictures,omitempty"`
}

type Pictures struct {
	ImgSrc    string `json:"img_src,omitempty"`
	ImgWidth  int    `json:"img_width,omitempty"`
	ImgHeight int    `json:"img_height,omitempty"`
	ImgSize   int    `json:"img_size,omitempty"`
}
