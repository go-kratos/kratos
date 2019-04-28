package conf

import (
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xtime "go-common/library/time"
)

// Redis redis
type Redis struct {
	*redis.Config
	Expire xtime.Duration
}

// Memcache config
type Memcache struct {
	*memcache.Config
	RelateExpire xtime.Duration
	ViewExpire   xtime.Duration
	ArcExpire    xtime.Duration
	CmsExpire    xtime.Duration
	HisExpire    xtime.Duration
	MangoExpire  xtime.Duration
}
