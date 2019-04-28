package v1

// PushRegisterRequest .
type PushRegisterRequest struct {
	RegID    string `json:"register_id" form:"register_id" validate:"required"`
	Platform string `json:"platform" form:"platform" validate:"required"`
	SDK      uint8  `json:"sdk" form:"sdk"`
}

// PushRegisterResponse .
type PushRegisterResponse struct{}

// PushCallbackRequest .
type PushCallbackRequest struct {
	Base
	TID string `json:"tid" form:"tid"`
	NID string `json:"nid" form:"nid"`
}

// PushCallbackResponse .
type PushCallbackResponse struct{}
