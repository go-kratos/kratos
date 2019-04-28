package v1

// ImgUploadRequest .
type ImgUploadRequest struct {
	Type int `json:"type" form:"type"`
}

//PreUploadRequest ...
type PreUploadRequest struct {
	Title     string `json:"title" form:"title" validate:"required"`
	Extension string `json:"extension" form:"extension"`
	FileExt   string `json:"file_ext" form:"file_ext" validate:"required"`
}

//CallBackRequest ..
type CallBackRequest struct {
	Svid     int64  `json:"biz_id" form:"biz_id" validate:"required"`
	URL      string `json:"url" form:"url" validate:"required"`
	Profile  string `json:"profile" form:"profile" validate:"required"`
	UploadID string `json:"upload_id" form:"upload_id" validate:"required"`
	Auth     string `json:"auth" form:"auth" validate:"required"`
}

// UploadCheckResponse 创作中心上传过滤
type UploadCheckResponse struct {
	Msg     string `json:"msg"`
	IsAllow bool   `json:"is_allow"`
}

//HomeImgRequest ...
type HomeImgRequest struct {
	SVID   int64  `json:"biz_id" form:"biz_id" validate:"required"`
	URL    string `json:"url" form:"url" validate:"required"`
	Width  int64  `json:"width" form:"width" validate:"required"`
	Height int64  `json:"height" form:"height" validate:"required"`
}
