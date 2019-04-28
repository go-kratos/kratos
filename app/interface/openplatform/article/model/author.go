package model

// RecommendAuthor .
type RecommendAuthor struct {
	UpID      int64  `json:"up_id"`
	RecReason string `json:"rec_reason"`
	RecType   int    `json:"rec_type"`
	Tid       int    `json:"tid"`
	Stid      int    `json:"second_tid"`
}

// RecommendAuthors .
type RecommendAuthors struct {
	Count   int          `json:"count"`
	Authors []*RecAuthor `json:"authors"`
}

// RecAuthor .
type RecAuthor struct {
	*AccountCard
	RecReason string `json:"rec_reason"`
}
