package pgc

import "go-common/library/time"

// SearPgcCon is used for setting search pgc content
type SearPgcCon struct {
	ID       int       `json:"id"`
	Category int       `json:"category"`
	Cover    string    `json:"cover"`
	Title    string    `json:"title"`
	PlayTime time.Time `json:"pubtime"`
	Role     string    `json:"cv"`
	Staff    string    `json:"staff"`
	Desc     string    `json:"description"`
}

// SearUgcCon is used for setting search ugc content
type SearUgcCon struct {
	AID     int       `json:"id"`
	Title   string    `json:"title"`
	Cover   string    `json:"cover"`
	Content string    `json:"description"`
	Pubtime time.Time `json:"pubtime"`
	Typeid  int       `json:"category"`
}

// SearchSug represents the search suggestion unit
type SearchSug struct {
	Term string `json:"term"`
	ID   int    `json:"id"`
	Type string `json:"type"`
}
