package model

// HotTags .
type HotTags struct {
	Rid  int64     `json:"rid"`
	Tags []*HotTag `json:"tags"`
}

// HotTag .
type HotTag struct {
	Rid       int64  `json:"-"`
	Tid       int64  `json:"tag_id"`
	Tname     string `json:"tag_name"`
	HighLight int64  `json:"highlight"`
	IsAtten   int8   `json:"is_atten"`
}

// SimilarTag .
type SimilarTag struct {
	Rid    int64  `json:"rid"`
	Rname  string `json:"rname"`
	Tid    int64  `json:"tid"`
	TCover string `json:"cover"`
	Tatten int    `json:"atten"`
	Tname  string `json:"tname"`
}
