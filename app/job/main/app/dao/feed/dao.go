package feed

import (
	"go-common/app/job/main/app/conf"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/library/cache/memcache"
	httpx "go-common/library/net/http/blademaster"
)

type Dao struct {
	c *conf.Config
	// rpc
	arcRPC *arcrpc.Service2
	// http client
	client     *httpx.Client
	clientAsyn *httpx.Client
	// hetongzi
	hot string
	// tag
	tags string
	// rcmdUp
	rcmdUp string
	// mc
	mc *memcache.Pool
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// rpc
		arcRPC: arcrpc.New2(c.ArchiveRPC),
		// http client
		client:     httpx.NewClient(c.HTTPClient),
		clientAsyn: httpx.NewClient(c.HTTPClientAsyn),
		// hetongzi
		hot: c.Host.Hetongzi + _hot,
		// tag
		tags: c.Host.APICo + _tags,
		// rcmdUp
		rcmdUp: c.Host.APP + _rcmdUp,
		// mc
		mc: memcache.NewPool(c.Memcache.Feed.Config),
	}
	return d
}
