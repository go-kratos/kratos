package feed

type Tag struct {
	ID   int64  `json:"tag_id,omitempty"`
	Name string `json:"tag_name,omitempty"`
}
