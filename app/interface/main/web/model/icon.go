package model

// IndexIcon index icon struct.
type IndexIcon struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	Links  []string `json:"links"`
	Icon   string   `json:"icon"`
	Weight int      `json:"weight"`
}
