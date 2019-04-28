package tag

type Tag struct {
	TagID      int64  `json:"tag_id"`
	Name       string `json:"tag_name"`
	Cover      string `json:"cover"`
	Likes      int64  `json:"likes"`
	Hates      int64  `json:"hates"`
	Liked      int8   `json:"liked"`
	Hated      int8   `json:"hated"`
	Attribute  int8   `json:"attribute"`
	IsActivity int8   `json:"is_activity"`
}
