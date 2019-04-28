package dao

import (
	"context"

	"go-common/app/interface/main/space/model"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -nullcache=&model.Notice{Notice:"ff2364a0be3d20e46cc69efb36afe9a5"} -check_null_code=$.Notice=="ff2364a0be3d20e46cc69efb36afe9a5"
	Notice(c context.Context, mid int64) (*model.Notice, error)
	// cache: -nullcache=&model.AidReason{Aid:-1} -check_null_code=$!=nil&&$.Aid==-1
	TopArc(c context.Context, mid int64) (*model.AidReason, error)
	// cache: -nullcache=&model.AidReasons{List:[]*model.AidReason{{Aid:-1}}} -check_null_code=len($.List)==1&&$.List[0].Aid==-1
	Masterpiece(c context.Context, mid int64) (*model.AidReasons, error)
	// cache: -nullcache=&model.ThemeDetails{List:[]*model.ThemeDetail{{ID:-1}}} -check_null_code=len($.List)==1&&$.List[0].ID==-1
	Theme(c context.Context, mid int64) (*model.ThemeDetails, error)
	// cache: -nullcache=-1 -check_null_code=$==-1
	TopDynamic(c context.Context, mid int64) (int64, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// get notice data from mc cache.
	// mc: -key=noticeKey
	CacheNotice(c context.Context, mid int64) (*model.Notice, error)
	// set notice data to mc cache.
	// mc: -key=noticeKey -expire=d.mcNoticeExpire -encode=pb
	AddCacheNotice(c context.Context, mid int64, data *model.Notice) error
	// mc: -key=noticeKey
	DelCacheNotice(c context.Context, mid int64) error
	// get top archive data from mc cache.
	// mc: -key=topArcKey
	CacheTopArc(c context.Context, mid int64) (*model.AidReason, error)
	// set top archive data to mc cache.
	// mc: -key=topArcKey -expire=d.mcTopArcExpire -encode=pb
	AddCacheTopArc(c context.Context, mid int64, data *model.AidReason) error
	// get top archive data from mc cache.
	// mc: -key=masterpieceKey
	CacheMasterpiece(c context.Context, mid int64) (*model.AidReasons, error)
	// set top archive data to mc cache.
	// mc: -key=masterpieceKey -expire=d.mcMpExpire -encode=pb
	AddCacheMasterpiece(c context.Context, mid int64, data *model.AidReasons) error
	// get theme data from mc cache.
	// mc: -key=themeKey
	CacheTheme(c context.Context, mid int64) (*model.ThemeDetails, error)
	// set theme data to mc cache.
	// mc: -key=themeKey -expire=d.mcThemeExpire -encode=pb
	AddCacheTheme(c context.Context, mid int64, data *model.ThemeDetails) error
	// mc: -key=themeKey
	DelCacheTheme(c context.Context, mid int64) error
	// get top dynamic id cache.
	// mc: -key=topDyKey
	CacheTopDynamic(c context.Context, key int64) (int64, error)
	// set top dynamic id cache.
	// mc: -key=topDyKey -expire=d.mcTopDyExpire -encode=raw
	AddCacheTopDynamic(c context.Context, key int64, value int64) error
}
