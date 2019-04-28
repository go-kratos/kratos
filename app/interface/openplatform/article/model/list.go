package model

import xtime "go-common/library/time"

// sort type
const (
	ListSortPtime = 0
	ListSortView  = 1
)

// CreativeList creative list
type CreativeList struct {
	*List
	Total int `json:"total"`
}

// ListArtMeta .
type ListArtMeta struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	State       int         `json:"state"`
	PublishTime xtime.Time  `json:"publish_time"`
	Position    int         `json:"-"`
	Words       int64       `json:"words"`
	ImageURLs   []string    `json:"image_urls"`
	Category    *Category   `json:"category"`
	Categories  []*Category `json:"categories"`
	Summary     string      `json:"summary"`
}

// Strong fill
func (a *ListArtMeta) Strong() {
	if a == nil {
		return
	}
	if a.ImageURLs == nil {
		a.ImageURLs = []string{}
	}
	if a.Category == nil {
		a.Category = &Category{}
	}
	if a.Categories == nil {
		a.Categories = []*Category{}
	}
}

// FullListArtMeta .
type FullListArtMeta struct {
	*ListArtMeta
	Stats     Stats `json:"stats"`
	LikeState int8  `json:"like_state"`
}

// IsNormal judge whether article's state is normal.
func (a *ListArtMeta) IsNormal() bool {
	return (a != nil) && (a.State >= StateOpen)
}

// ListArticles list articles
type ListArticles struct {
	List      *List          `json:"list"`
	Articles  []*ListArtMeta `json:"articles"`
	Author    *Author        `json:"author"`
	Last      ListArtMeta    `json:"last"`
	Attention bool           `json:"attention"`
}

// WebListArticles .
type WebListArticles struct {
	List      *List              `json:"list"`
	Articles  []*FullListArtMeta `json:"articles"`
	Author    *Author            `json:"author"`
	Last      ListArtMeta        `json:"last"`
	Attention bool               `json:"attention"`
}

// ListInfo list info
type ListInfo struct {
	List  *List        `json:"list"`
	Last  *ListArtMeta `json:"last"`
	Next  *ListArtMeta `json:"next"`
	Now   int          `json:"now"`
	Total int          `json:"total"`
}

// UpLists .
type UpLists struct {
	Lists []*List `json:"lists"`
	Total int     `json:"total"`
}
