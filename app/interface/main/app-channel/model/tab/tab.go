package tab

type Menu struct {
	ID       int64  `json:"id,omitempty"`
	TagID    int64  `json:"tag_id,omitempty"`
	TabID    int64  `json:"tab_id,omitempty"`
	Title    string `json:"title,omitempty"`
	Name     string `json:"name,omitempty"`
	Priority int    `json:"priority,omitempty"`
}
