package model

import "go-common/library/time"

// VipPushData .
type VipPushData struct {
	ID              int64               `json:"id" form:"id"`
	DisableType     int8                `json:"disable_type"`
	GroupName       string              `json:"group_name" form:"group_name"`
	Title           string              `json:"title" form:"title"`
	Content         string              `json:"content" form:"content" validate:"required"`
	PushTotalCount  int32               `json:"-"`
	PushedCount     int32               `json:"-"`
	PushProgress    string              `json:"push_progress"`
	ProgressStatus  int8                `json:"progress_status"`
	Operator        string              `json:"operator"`
	Status          int8                `json:"status"`
	Platform        string              `json:"platform" form:"platform"`
	LinkType        int32               `json:"link_type" form:"link_type" validate:"required"`
	ErrorCode       int32               `json:"error_code"`
	LinkURL         string              `json:"link_url" form:"link_url" validate:"required"`
	ExpiredDayStart int32               `json:"expired_day_start" form:"expired_day_start"`
	ExpiredDayEnd   int64               `json:"expired_day_end" form:"expired_day_end"`
	EffectStartDate time.Time           `json:"effect_start_date" form:"effect_start_date" validate:"required"`
	EffectEndDate   time.Time           `json:"effect_end_date" form:"effect_end_date" validate:"required"`
	PushStartTime   string              `json:"push_start_time" form:"push_start_time" validate:"required"`
	PushEndTime     string              `json:"push_end_time" form:"push_end_time" validate:"required"`
	PlatformArr     []*PushDataPlatform `json:"platform_arr"`
}

// PushDataPlatform .
type PushDataPlatform struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
	Build     int64  `json:"build"`
}
