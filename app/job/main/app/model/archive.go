package model

// Archive model
type Archive struct {
	Aid       int64  `json:"aid"`
	Mid       int64  `json:"mid"`
	TypeID    int16  `json:"typeid"`
	Duration  int    `json:"duration"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Content   string `json:"content"`
	Attribute int32  `json:"attribute"`
	Copyright int8   `json:"copyright"`
	State     int    `json:"state"`
	Access    int    `json:"access"`
	PubTime   string `json:"pubtime"`
}
