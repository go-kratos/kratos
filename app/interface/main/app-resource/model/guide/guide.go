package guide

type Interest struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
}

type InterestTM struct {
	Interests []*Interest `json:"interests,omitempty"`
	FeedType  int         `json:"feed_type"`
}
