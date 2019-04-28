package topic

type Topic struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	H5URL   string `json:"h5_url"`
	H5Cover string `json:"h5_cover"`
	Desc    string `json:"desc"`
}
