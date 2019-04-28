package tag

// Tag struct
type Tag struct {
	ID      int64  `json:"tag_id,omitempty"`
	Name    string `json:"tag_name,omitempty"`
	IsAtten int8   `json:"is_atten,omitempty"`
	Count   *struct {
		Atten int `json:"atten,omitempty"`
	} `json:"count,omitempty"`
	Cover     string `json:"cover,omitempty"`
	Likes     int64  `json:"likes,omitempty"`
	Hates     int64  `json:"hates,omitempty"`
	Liked     int8   `json:"liked,omitempty"`
	Hated     int8   `json:"hated,omitempty"`
	Attribute int8   `json:"attribute,omitempty"`
}

// Hot struct
type Hot struct {
	Rid  int16  `json:"rid"`
	Tags []*Tag `json:"tags"`
}

// SubTag struct
type SubTag struct {
	Count   int    `json:"count"`
	SubTags []*Tag `json:"subscribe"`
}
