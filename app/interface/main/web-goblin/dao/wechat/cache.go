package wechat

import (
	"context"

	"go-common/app/interface/main/web-goblin/model/wechat"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache
	AccessToken(c context.Context) (*wechat.AccessToken, error)
}
