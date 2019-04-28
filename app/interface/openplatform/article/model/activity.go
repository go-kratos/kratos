package model

// ActInfo .
type ActInfo struct {
	Activities []*Activity `json:"activities"`
	Banners    []*Banner   `json:"banners"`
}

// AnniversaryInfo .
type AnniversaryInfo struct {
	Mid        int64              `json:"mid"`
	Uname      string             `json:"uname"`
	Face       string             `json:"face"`
	ReaderInfo *AnniversaryReader `json:"reader_info"`
	AuthorInfo *AnniversaryAuthor `json:"author_info"`
}

// AnniversaryAuthor .
type AnniversaryAuthor struct {
	Articles    int32  `json:"articles"`
	Words       int64  `json:"words"`
	Views       int64  `json:"views"`
	Coins       int64  `json:"coins"`
	Title       string `json:"title"`
	Publish     string `json:"publish"`
	Rank        string `json:"rank"`
	ReaderMid   int64  `json:"reader"`
	ReaderUname string `json:"reader_name"`
	ReaderFace  string `json:"reader_face"`
}

// AnniversaryReader .
type AnniversaryReader struct {
	Words        int64  `json:"words"`
	Views        int64  `json:"views"`
	Coins        int64  `json:"coins"`
	Comments     int64  `json:"comments"`
	Title        string `json:"title"`
	AuthorMid    int64  `json:"author"`
	AuthorUname  string `json:"author_name"`
	Rank         string `json:"rank"`
	FirstComment string `json:"first_comment"`
	CommentDate  string `json:"comment_date"`
}
