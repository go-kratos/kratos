package v1

// NoticeRequest .
type NoticeRequest struct {
	RegID    string `json:"register_id" form:"register_id" validate:"required"`
	Platform int32  `json:"platform" form:"platform" validate:"required"`
	SDK      int32  `json:"sdk" form:"sdk" validate:"required"`
	Title    string `json:"title" form:"title" validate:"required"`
	Content  string `json:"content" form:"content" validate:"required"`
	Schema   string `json:"schema" form:"schema" validate:"required"`
	Callback string `json:"callback" form:"callback" validate:"required"`
}

// NoticeResponse .
type NoticeResponse struct{}

// MessageRequest .
type MessageRequest struct {
	RegID       string `json:"register_id" form:"register_id" validate:"required"`
	Platform    int32  `json:"platform" form:"platform" validate:"required"`
	SDK         int32  `json:"sdk" form:"sdk" validate:"required"`
	Title       string `json:"title" form:"title" validate:"required"`
	Content     string `json:"content" form:"content" validate:"required"`
	ContentType string `json:"content_type" form:"content_type" validate:"required"`
	Schema      string `json:"schema" form:"schema" validate:"required"`
	Callback    string `json:"callback" form:"callback" validate:"required"`
}

// MessageResponse .
type MessageResponse struct{}
