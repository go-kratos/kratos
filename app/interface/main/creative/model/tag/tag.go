package tag

const (
	// UserTag 普通tag
	UserTag = int8(0)
	// UpTag up主tag
	UpTag = int8(1)
	// OfficailClassifyTag  官方-分类tag
	OfficailClassifyTag = int8(2)
	// OfficailContentTag 官方-内容tag
	OfficailContentTag = int8(3)
	// OfficailActiveTag 官方-活动tag
	OfficailActiveTag = int8(4)
	// TagStateNormal normal
	TagStateNormal = 0
	// TagStateDel del
	TagStateDel = 1
	// TagStateHide hide
	TagStateHide = 2
)

// Meta for tag info.
type Meta struct {
	TagID   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
}

// Tag str
type Tag struct {
	ID      int64  `json:"tag_id"`
	Name    string `json:"tag_name"`
	Cover   string `json:"cover"`
	Content string `json:"content"`
	Type    int8   `json:"type"`
	State   int8   `json:"state"`
}

// StaffTitle 联合投稿职能
type StaffTitle struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
