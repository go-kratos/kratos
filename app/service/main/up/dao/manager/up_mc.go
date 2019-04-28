package manager

import (
	"context"
	"fmt"

	upgrpc "go-common/app/service/main/up/api/v1"
)

const (
	_upSpecialKey = "up_special_%d"
)

// upSpecialCacheKey 缓存key
func upSpecialCacheKey(mid int64) string {
	return fmt.Sprintf(_upSpecialKey, mid)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -nullcache=&upgrpc.UpSpecial{GroupIDs:[]int64{-1}} -check_null_code=$!=nil&&len($.GroupIDs)>0&&$.GroupIDs[0]==-1
	UpSpecial(c context.Context, mid int64) (us *upgrpc.UpSpecial, err error)
	// cache: -batch=100 -max_group=1 -batch_err=break -nullcache=&upgrpc.UpSpecial{GroupIDs:[]int64{-1}} -check_null_code=$!=nil&&len($.GroupIDs)>0&&$.GroupIDs[0]==-1
	UpsSpecial(c context.Context, mids []int64) (map[int64]*upgrpc.UpSpecial, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// mc: -key=upSpecialCacheKey -expire=d.upSpecialExpire -encode=pb
	AddCacheUpSpecial(c context.Context, mid int64, us *upgrpc.UpSpecial) (err error)
	// mc: -key=upSpecialCacheKey
	CacheUpSpecial(c context.Context, mid int64) (res *upgrpc.UpSpecial, err error)
	// mc: -key=upSpecialCacheKey
	DelCacheUpSpecial(c context.Context, mid int64) (err error)
	// mc: -key=upSpecialCacheKey -expire=d.upSpecialExpire -encode=pb
	AddCacheUpsSpecial(c context.Context, mu map[int64]*upgrpc.UpSpecial) (err error)
	// mc: -key=upSpecialCacheKey
	CacheUpsSpecial(c context.Context, mid []int64) (res map[int64]*upgrpc.UpSpecial, err error)
	// mc: -key=upSpecialCacheKey
	DelCacheUpsSpecial(c context.Context, mids []int64) (err error)
}
