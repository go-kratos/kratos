package model

const (
	_defaultWmPaddingX = 10
	_defaultWmPaddingY = 10
	_defaultWmScale    = float64(1) / 24

	// delete status .please read document.
	// http://info.bilibili.co/pages/viewpage.action?pageId=8718262#bfs%E7%AE%A1%E7%90%86%E5%90%8E%E5%8F%B0%E7%9B%B8%E5%85%B3%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3-db%E8%AE%BE%E8%AE%A1

	// PassStatus express pass status.
	PassStatus = 2
	// DeleteStatus .
	DeleteStatus = 3
)

// AddParam describe add api param
type AddParam struct {
	Bucket   string `json:"bucket" form:"bucket" validate:"required"`
	FileName string `json:"filename" form:"filename" validate:"required"`
	URL      string `json:"url" form:"url"`
	Sex      int    `json:"sex" form:"sex"`
	Politics int    `json:"politics" form:"politics"`
}

// ListParam describe list api param
type ListParam struct {
	Bucket string `json:"bucket" form:"bucket" validate:"required"`
	State  int    `json:"state" form:"state" validate:"required,min=0"`
	PN     int    `json:"pn" form:"pn" validate:"min=1"`
	PS     int    `json:"ps" form:"ps" validate:"min=1"`
}

// MultiListParam describe list api param
type MultiListParam struct {
	Bucket []string `json:"bucket" form:"bucket"`
	State  int      `json:"state" form:"state" validate:"min=0"`
	PN     int      `json:"pn" form:"pn" validate:"min=1" default:"1"`
	PS     int      `json:"ps" form:"ps" validate:"min=1" default:"50"`
}

// DeleteParam describe list api param
type DeleteParam struct {
	Rid      int    `json:"rid" form:"rid" validate:"required"`
	Bucket   string `json:"bucket"`
	FileName string `json:"filename"`
	AdminID  int64  `json:"admin_id"`
}

// DeleteV2Param describe list api param
type DeleteV2Param struct {
	Rid      int    `json:"rid" form:"rid" validate:"required"`
	Status   int    `json:"status" form:"status" validate:"required"`
	Bucket   string `json:"bucket"`
	FileName string `json:"filename"`
	AdminID  int64  `json:"admin_id"`
}

// DeleteRawParam describe list api param
type DeleteRawParam struct {
	Bucket   string `json:"bucket" form:"bucket" validate:"required"`
	FileName string `json:"filename" form:"filename" validate:"required"`
}

// AddBucketParam .
type AddBucketParam struct {
	Name         string `form:"name" json:"name" validate:"required"`
	Property     int    `form:"property" json:"property" validate:"min=0,max=3"`
	KeyID        string `form:"key_id" json:"key_id" validate:"required"`
	KeySecret    string `form:"key_secret" json:"key_secret" validate:"required"`
	PurgeCDN     bool   `form:"purge_cdn" json:"purge_cdn"`
	CacheControl int    `form:"cache_control" json:"cache_control"`
	Domain       string `form:"domain" json:"domain"`
}

// AddDirParam .
type AddDirParam struct {
	BucketName string `form:"bucket_name" validate:"required"`
	DirName    string `form:"dir_name" validate:"required"`
	Pic        string `form:"pic"`
	Rate       string `form:"rate"`
}

// ListBucketParam .
type ListBucketParam struct {
	PN int `form:"pn" validate:"min=1"`
	PS int `form:"ps" validate:"min=1"`
}

// UploadParam .
type UploadParam struct {
	Bucket      string  `form:"bucket" json:"bucket" validate:"required" `
	ContentType string  `form:"content_type" json:"content_type"`
	Auth        string  `form:"auth" json:"-"`
	Dir         string  `form:"dir" json:"dir"`
	FileName    string  `form:"file_name" json:"file_name"`
	WmKey       string  `form:"wm_key" json:"wm_key"`
	WmText      string  `form:"wm_text" json:"wm_text"`
	WmPaddingX  int     `form:"wm_padding_x" json:"wm_padding_x"`
	WmPaddingY  int     `form:"wm_padding_y" json:"wm_padding_y"`
	WmScale     float64 `form:"wm_scale" json:"wm_scale"`
}

// WMInit init UploadParam
func (up *UploadParam) WMInit() {
	if up.WmKey != "" || up.WmText != "" {
		if up.WmPaddingX < 0 {
			up.WmPaddingX = _defaultWmPaddingX
		}
		if up.WmPaddingY < 0 {
			up.WmPaddingY = _defaultWmPaddingY
		}
		if up.WmScale <= 0 {
			up.WmScale = _defaultWmScale
		}
	}
}

// UploadResult .
type UploadResult struct {
	Location string `json:"location"`
	Etag     string `json:"etag"`
}

// MultiListResult .
type MultiListResult struct {
	Bucket string    `json:"bucket"`
	Imgs   []*Record `json:"imgs"`
}
