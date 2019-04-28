package v1

import "go-common/app/interface/bbq/app-bbq/model"

// ReportConfigResponse .
type ReportConfigResponse struct {
	Report  []*model.ReportConfig `json:"report_config,omitempty"`
	Reasons []*model.ReasonConfig `json:"reason_config,omitempty"`
}

// ReportRequest .
type ReportRequest struct {
	Type   int16 `json:"type" form:"type" validate:"required"`
	UpMID  int64 `json:"up_mid" form:"up_mid"`
	SVID   int64 `json:"svid" form:"svid"`
	RpID   int64 `json:"rpid" form:"rpid"`
	Danmu  int64 `json:"danmu" form:"danmu"`
	Reason int16 `json:"reason" form:"reason"`
}
