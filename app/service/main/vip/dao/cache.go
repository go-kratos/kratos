package dao

import (
	"context"

	"go-common/app/service/main/vip/model"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -nullcache=&model.OpenBindInfo{Mid:-1} -check_null_code=$!=nil&&$.Mid==-1
	BindInfoByMid(c context.Context, mid int64, appID int64) (*model.OpenBindInfo, error)
	// cache: -nullcache=&model.OpenInfo{Mid:-1} -check_null_code=$!=nil&&$.Mid==-1
	OpenInfoByOpenID(c context.Context, openID string, appID int64) (*model.OpenInfo, error)
}
