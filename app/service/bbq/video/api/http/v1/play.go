package v1

// VideoPlay playinfo
type VideoPlay struct {
	SVID           int64       `json:"svid"`
	ExpireTime     int64       `json:"expire_time"`     //过期时间
	FileInfo       []*FileInfo `json:"file_info"`       //分片信息
	Quality        int64       `json:"quality"`         //清晰度
	SupportQuality []int64     `json:"support_quality"` //支持清晰度
	URL            string      `json:"url"`             //基础url
	CurrentTime    int64       `json:"current_time"`    //当前时间戳
}

// FileInfo bvc fileinfo
type FileInfo struct {
	Ahead      string `json:"ahead"`
	FileSize   int64  `json:"filesize"`
	TimeLength int64  `json:"timelength"`
	Vhead      string `json:"vhead"`
	Path       string `json:"path"`
	URL        string `json:"url"`
	URLBc      string `json:"url_bc"`
}

//VideoPlayRequest ..
type VideoPlayRequest struct {
	SIVD []int64 `json:"svid" form:"svid" validate:"required"`
}
