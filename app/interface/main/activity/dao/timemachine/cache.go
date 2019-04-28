package timemachine

import (
	"context"
	"fmt"

	"go-common/app/interface/main/activity/model/timemachine"
)

func timemachineKey(mid int64) string {
	return fmt.Sprintf("tm_%d", mid)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache
	Timemachine(c context.Context, mid int64) (*timemachine.Item, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// mc: -key=timemachineKey -expire=d.mcTmExpire -encode=pb
	AddCacheTimemachine(c context.Context, mid int64, data *timemachine.Item) error
	// mc: -key=timemachineKey
	CacheTimemachine(c context.Context, mid int64) (*timemachine.Item, error)
}
