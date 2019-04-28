package model

// FopVideoCovers video cloud fop=videocovers_format response
type FopVideoCovers []struct {
	Count     int    `json:"count"`
	URLFormat string `json:"url_format"`
}
