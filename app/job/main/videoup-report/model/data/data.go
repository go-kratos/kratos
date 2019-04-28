package data

// HotArchiveRes hot recheck archive response
type HotArchiveRes struct {
	Code       int    `json:"code"`
	Note       bool   `json:"note"`
	SourceDate string `json:"source_date"`
	Num        int    `json:"num"`
	List       []struct {
		Aid   int64 `json:"aid"`
		Score int   `json:"score"`
	} `json:"list"`
}
