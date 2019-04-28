package model

import (
	accv1 "go-common/app/service/main/account/api"
	feedmdl "go-common/app/service/main/feed/model"
)

// Feed feed
type Feed struct {
	*feedmdl.Feed
	OfficialVerify *accv1.OfficialInfo `json:"official_verify,omitempty"`
}
