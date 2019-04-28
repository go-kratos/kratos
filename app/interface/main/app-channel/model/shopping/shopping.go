package shopping

type Card struct {
	ID               int64  `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	PerformanceImage string `json:"performance_image,omitempty"`
	STime            string `json:"stime,omitempty"`
	ETime            string `json:"etime,omitempty"`
	Tags             []*struct {
		TagID   int64  `json:"tag_id,omitempty"`
		TagName string `json:"tag_name,omitempty"`
	} `json:"tags,omitempty"`
	CityName string `json:"city_name,omitempty"`
	URL      string `json:"url,omitempty"`
	Subname  string `json:"subname,omitempty"`
	Pricelt  string `json:"pricelt,omitempty"`
	Want     string `json:"want,omitempty"`
	Type     int8   `json:"type,omitempty"`
}
