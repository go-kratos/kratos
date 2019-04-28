package model

type MusicResult struct {
	Code int              `json:"code"`
	Data map[int64]*Music `json:"data"`
}

type Music struct {
	ID    int64  `json:"song_id"`
	Cover string `json:"cover_url"`
	Title string `json:"title"`
}
