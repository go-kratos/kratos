package archive

// BIZPOI 1
const (
	BIZPOI = int64(1)
)

// PoiObj str
type PoiObj struct {
	POI           string      `json:"poi"`
	Type          int32       `json:"type"`
	Addr          string      `json:"address"`
	Title         string      `json:"title"`
	ShowTitle     string      `json:"show_title"`
	AdInfo        *AdInfo     `json:"ad_info"`
	Ancestors     []*Ancestor `json:"ancestors"`
	Distance      float64     `json:"distance"`
	ShowDistrance string      `json:"show_distance"`
	Location      *Location   `json:"location"`
}

// AdInfo str
type AdInfo struct {
	Nation string `json:"nation"`
	Provin string `json:"province"`
	Distri string `json:"district"`
	City   string `json:"city"`
}

// Ancestor str
type Ancestor struct {
	POI  string `json:"poi"`
	Type int32  `json:"type"`
}

// Location str
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
