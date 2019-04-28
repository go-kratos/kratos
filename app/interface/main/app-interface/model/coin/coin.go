package coin

type Arc struct {
	Aid         int64  `json:"aid,omitempty"`
	Tid         int64  `json:"tid,omitempty"`
	Name        string `json:"tname,omitempty"`
	Copyright   int    `json:"copyright,omitempty"`
	Title       string `json:"title,omitempty"`
	Pic         string `json:"pic,omitempty"`
	Play        int    `json:"play,omitempty"`
	VideoReview int    `json:"video_review,omitempty"`
}
