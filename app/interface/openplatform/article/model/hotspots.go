package model

// hotspot type
const (
	HotspotTypeView  = 0
	HotspotTypePtime = 1
)

// HotspotTypes types
var HotspotTypes = [...]int8{HotspotTypeView, HotspotTypePtime}

// Hotspot model
type Hotspot struct {
	ID          int64        `json:"id"`
	Tag         string       `json:"tag"`
	Title       string       `json:"title"`
	TopArticles []int64      `json:"top_articles"`
	Icon        bool         `json:"icon"`
	Stats       HotspotStats `json:"stats"`
}

// HotspotStats .
type HotspotStats struct {
	Read  int64 `json:"read"`
	Reply int64 `json:"reply"`
	Count int64 `json:"count"`
}

// SearchArt search article model
type SearchArt struct {
	ID          int64
	PublishTime int64
	Tags        []string
	StatsView   int64
	StatsReply  int64
}

// HotspotResp model
type HotspotResp struct {
	Hotspot  *Hotspot        `json:"hotspot"`
	Articles []*MetaWithLike `json:"articles"`
}
