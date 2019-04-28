package tag

type Tag struct {
	ID      int64  `json:"tag_id,omitempty"`
	Name    string `json:"tag_name,omitempty"`
	IsAtten int8   `json:"is_atten,omitempty"`
	Count   *struct {
		Atten int `json:"atten,omitempty"`
	} `json:"count,omitempty"`
}

type Hot struct {
	Rid  int16  `json:"rid"`
	Tags []*Tag `json:"tags"`
}

type SubTag struct {
	Count   int    `json:"count"`
	SubTags []*Tag `json:"subscribe"`
}
