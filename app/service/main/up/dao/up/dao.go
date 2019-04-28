package up

import (
	"context"
	"time"

	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/dao/global"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"

	article "go-common/app/interface/openplatform/article/rpc/client"
)

// DateLayout date layout
const DateLayout = "2006-01-02"

// Dao is creative dao.
type Dao struct {
	// config
	c *conf.Config
	// db
	db *sql.DB
	//cache tool
	cache *fanout.Fanout
	//memcache pool
	mcPool *memcache.Pool
	mc     *memcache.Pool
	//up expiration
	upExpire int32
	// http client
	client *httpx.Client
	//api url
	picUpInfoURL   string
	blinkUpInfoURL string
	// rpc
	art *article.Service
}

// New init db.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:        c,
		db:       sql.NewMySQL(c.DB.Creative),
		cache:    global.GetWorker(),
		mcPool:   memcache.NewPool(c.Memcache.Up),
		upExpire: int32(time.Duration(c.Memcache.UpExpire) / time.Second),
		// http client
		client:         httpx.NewClient(c.HTTPClient.Normal),
		picUpInfoURL:   c.Host.Live + _picUpInfoURI,
		blinkUpInfoURL: c.Host.Live + _blinkUpInfoURI,
		// rpc
		art: article.New(c.ArticleRPC),
	}
	d.mc = d.mcPool
	return
}

// Ping creativeDb
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close db
func (d *Dao) Close() (err error) {
	if d.db != nil {
		d.db.Close()
	}
	if d.mcPool != nil {
		d.mcPool.Close()
	}
	return
}

//GetHTTPClient get http client
func (d *Dao) GetHTTPClient() *httpx.Client {
	return d.client
}
