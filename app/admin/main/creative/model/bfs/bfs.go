package bfs

//FileInfo : the uploaded file information
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
	Md5  string `json:"md5"`
	URL  string `json:"url"`
}
