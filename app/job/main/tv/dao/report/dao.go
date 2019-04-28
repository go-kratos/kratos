package report

import (
	"go-common/app/job/main/tv/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao .
type Dao struct {
	conf  *conf.Config
	httpR *bm.Client
	mc    *memcache.Pool
	DB    *sql.DB
}

// New create a instance of Dao and return .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:  c,
		httpR: bm.NewClient(c.DpClient),
		mc:    memcache.NewPool(c.Memcache.Config),
		DB:    sql.NewMySQL(c.Mysql),
	}
	return
}
