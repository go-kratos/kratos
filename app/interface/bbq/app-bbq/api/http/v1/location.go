package v1

import (
	"go-common/app/interface/bbq/app-bbq/model"
)

// LocationRequest .
type LocationRequest struct {
	PID int32 `form:"pid"`
}

// LocationResponse .
type LocationResponse struct {
	List []*model.Location
}
