package model

// UpBaseStat for up base.
type UpBaseStat struct {
	View  int64 `json:"view"`
	Reply int64 `json:"reply"`
	Dm    int64 `json:"dm"`
	Fans  int64 `json:"fans"`
	Fav   int64 `json:"fav"`
	Like  int64 `json:"like"`
	Share int64 `json:"share"`
}
