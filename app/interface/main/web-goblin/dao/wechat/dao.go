package wechat

import (
	"go-common/app/interface/main/web-goblin/conf"
	"go-common/library/cache"
	"go-common/library/cache/redis"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao struct.
type Dao struct {
	// config
	c *conf.Config
	// redis
	redis *redis.Pool
	// httpClient
	httpClient *bm.Client
	// url
	wxAccessTokenURL string
	wxQrcodeURL      string
	cache            *cache.Cache
}

// New new dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// config
		c:          c,
		redis:      redis.NewPool(c.Redis.Config),
		httpClient: bm.NewClient(c.HTTPClient),
		cache:      cache.New(1, 1024),
	}
	d.wxAccessTokenURL = d.c.Host.Wechat + _accessTokenURI
	d.wxQrcodeURL = d.c.Host.Wechat + _qrcodeURI
	return
}
