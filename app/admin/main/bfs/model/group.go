package model

// Group .
type Group struct {
	Stores     []string          `json:"stores"`
	StoreDatas map[string]*Store `json:"store_datas"`
	Total      struct {
		Space     int64 `json:"space"`
		FreeSpace int64 `json:"free_space"`
		Volumes   int64 `json:"volumes"`
	} `json:"total"`
}
