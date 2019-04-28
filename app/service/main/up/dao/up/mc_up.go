package up

import (
	"context"
	"fmt"

	"go-common/app/service/main/up/model"
)

const (
	_upKey        = "up_srv_%d"
	_upSwitchKey  = "up_sw_%d"
	_upInfoActive = "up_info_active_%d"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	//cache: -nullcache=&model.Up{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	Up(c context.Context, mid int64) (up *model.Up, err error)
	//cache: -nullcache=&model.UpSwitch{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	UpSwitch(c context.Context, mid int64) (up *model.UpSwitch, err error)
	//cache: -nullcache=&model.UpInfoActiveReply{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	UpInfoActive(c context.Context, mid int64) (up *model.UpInfoActiveReply, err error)
	// cache: -batch=100 -max_group=1 -batch_err=break -nullcache=&model.UpInfoActiveReply{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	UpsInfoActive(c context.Context, mids []int64) (res map[int64]*model.UpInfoActiveReply, err error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	//mc: -key=upCacheKey -expire=d.upExpire -encode=json
	AddCacheUp(c context.Context, mid int64, up *model.Up) (err error)
	//mc: -key=upCacheKey
	CacheUp(c context.Context, mid int64) (up *model.Up, err error)
	//mc: -key=upCacheKey
	DelCacheUp(c context.Context, mid int64) (err error)
	//mc: -key=upSwitchKey -expire=d.upExpire -encode=json
	AddCacheUpSwitch(c context.Context, mid int64, up *model.UpSwitch) (err error)
	//mc: -key=upSwitchKey
	CacheUpSwitch(c context.Context, mid int64) (res *model.UpSwitch, err error)
	//mc: -key=upSwitchKey
	DelCacheUpSwitch(c context.Context, mid int64) (err error)
	//mc: -key=upInfoActiveKey -expire=d.upExpire -encode=json
	AddCacheUpInfoActive(c context.Context, mid int64, up *model.UpInfoActiveReply) (err error)
	//mc: -key=upInfoActiveKey
	CacheUpInfoActive(c context.Context, mid int64) (res *model.UpInfoActiveReply, err error)
	//mc: -key=upInfoActiveKey
	DelCacheUpInfoActive(c context.Context, mid int64) (err error)
	// mc: -key=upInfoActiveKey -expire=d.upExpire -encode=json
	AddCacheUpsInfoActive(c context.Context, res map[int64]*model.UpInfoActiveReply) (err error)
	// mc: -key=upInfoActiveKey
	CacheUpsInfoActive(c context.Context, mids []int64) (res map[int64]*model.UpInfoActiveReply, err error)
	// mc: -key=upInfoActiveKey
	DelCacheUpsInfoActive(c context.Context, mids []int64) (err error)
}

//upCacheKey 缓存key
func upCacheKey(mid int64) string {
	return fmt.Sprintf(_upKey, mid)
}

//upSwitchCacheKey 缓存key
func upSwitchKey(mid int64) string {
	return fmt.Sprintf(_upSwitchKey, mid)
}

//upInfoActiveCacheKey 缓存key
func upInfoActiveKey(mid int64) string {
	return fmt.Sprintf(_upInfoActive, mid)
}

//DelCpUp 异步删除缓存
func (d *Dao) DelCpUp(c context.Context, mid int64) {
	d.cache.Do(c, func(c context.Context) {
		d.DelCacheUp(c, mid)
	})
}
