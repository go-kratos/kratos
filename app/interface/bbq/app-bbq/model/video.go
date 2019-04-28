package model

//VideoStatistics .
type VideoStatistics struct {
	SVID      int64 `json:"svid"`
	Play      int64 `json:"play"`
	Subtitles int64 `json:"subtitles"`
	Like      int64 `json:"like"`
	Share     int64 `json:"share"`
	Report    int64 `json:"report"`
}
