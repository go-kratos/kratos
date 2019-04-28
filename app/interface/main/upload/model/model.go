package model

const (
	_defaultWmPaddingX = 10
	_defaultWmPaddingY = 10
	_defaultWmScale    = float64(1) / 24
)

// Result upload result
type Result struct {
	Location string `json:"location"`
	Etag     string `json:"etag"`
}

// ResultWm watermark result
type ResultWm struct {
	Location string `json:"location"`
	Md5      string `json:"md5"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// UploadParam upload params
type UploadParam struct {
	Bucket      string  `form:"bucket" json:"bucket" validate:"required" `
	ContentType string  `form:"content_type" json:"content_type"`
	Dir         string  `form:"dir" json:"dir"`
	FileName    string  `form:"file_name" json:"file_name"`
	WmKey       string  `form:"wm_key" json:"wm_key"`
	WmText      string  `form:"wm_text" json:"wm_text"`
	WmPaddingX  int     `form:"wm_padding_x" json:"wm_padding_x"`
	WmPaddingY  int     `form:"wm_padding_y" json:"wm_padding_y"`
	WmScale     float64 `form:"wm_scale" json:"wm_scale"`
}

// WMInit init watermark default value.
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
