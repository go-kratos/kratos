package operation

// Operation tool.
type Operation struct {
	ID       int64  `json:"id"`
	Ty       string `json:"-"`
	Rank     string `json:"rank"`
	Pic      string `json:"pic"`
	Link     string `json:"link"`
	Content  string `json:"content"`
	Remark   string `json:"remark"`
	Note     string `json:"note"`
	Stime    string `json:"start_time"`
	Etime    string `json:"end_time"`
	AppPic   string `json:"-"`
	Platform int8   `json:"-"`
}

// Banner for app index.
type Banner struct {
	Ty      string `json:"-"`
	Rank    string `json:"rank"`
	Pic     string `json:"pic"`
	Link    string `json:"link"`
	Content string `json:"content"`
}

// BannerCreator for creator index.
type BannerCreator struct {
	Ty      string `json:"-"`
	Rank    int    `json:"rank"`
	Pic     string `json:"pic"`
	Link    string `json:"link"`
	Content string `json:"content"`
	Stime   int64  `json:"start_time"`
	Etime   int64  `json:"end_time"`
}

// BannerList for operation list.
type BannerList struct {
	BannerCreator []*BannerCreator `json:"operations"`
	Pn            int              `json:"pn"`
	Ps            int              `json:"ps"`
	Total         int              `json:"total"`
}

// FullTypes  get full operations.
func FullTypes() (tys []string) {
	return []string{"'play'", "'notice'", "'road'", "'creative'", "'collect_arc'"}
}
