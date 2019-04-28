package model

// Tag tag struct.
type Tag struct {
	ID    int64  `json:"tag_id"`
	Name  string `json:"tag_name"`
	Cover string `json:"cover"`
}
