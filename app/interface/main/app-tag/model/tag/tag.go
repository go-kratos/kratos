package tag

import "go-common/app/service/main/archive/api"

type Tag struct {
	ID      int64  `json:"tag_id,omitempty"`
	Name    string `json:"tag_name,omitempty"`
	IsAtten int8   `json:"is_atten,omitempty"`
	Count   *Count `json:"count,omitempty"`
}

type Count struct {
	Atten int `json:"atten"`
}

type Hot struct {
	Rid  int16  `json:"rid"`
	Tags []*Tag `json:"tags"`
}

type TagArc struct {
	Tag   *Tag
	Aid   int64
	TagID int64
}

type SubTag struct {
	Count   int    `json:"count"`
	SubTags []*Tag `json:"subscribe"`
}

type SimilarTag struct {
	TagId   int64  `json:"tid"`
	TagName string `json:"tname"`
	Rid     int    `json:"rid,omitempty"`
	Rname   string `json:"rname,omitempty"`
	Reid    int    `json:"reid,omitempty"`
	Rename  string `json:"rename,omitempty"`
}

type News struct {
	Archives []*api.Arc `json:"archives"`
}
