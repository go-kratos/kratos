package medal

import (
	"go-common/app/job/main/usersuit/conf"
	"go-common/app/service/main/usersuit/rpc/client"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

var (
	_updateinfo = "/mingpai/api/updateinfo/%s"
)

// Dao struct info of Dao.
type Dao struct {
	db         *sql.DB
	c          *conf.Config
	client     *bm.Client
	suitRPC    *client.Service2
	updateInfo string
	// memcache
	mc *memcache.Pool
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		db:         sql.NewMySQL(c.Mysql),
		client:     bm.NewClient(c.HTTPClient),
		updateInfo: c.Properties.UpInfoURL + _updateinfo,
		suitRPC:    client.New(c.SuitRPC),
		// memcache
		mc: memcache.NewPool(c.Memcache.Config),
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
