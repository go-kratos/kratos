package goblin

import (
	"go-common/app/interface/main/tv/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao .
type Dao struct {
	conf   *conf.Config
	client *bm.Client
	db     *sql.DB
	mc     *memcache.Pool
}

// New .
func New(c *conf.Config) *Dao {
	return &Dao{
		conf:   c,
		client: bm.NewClient(c.PlayurlClient),
		db:     sql.NewMySQL(c.Mysql),
		mc:     memcache.NewPool(c.Memcache.Config),
	}
}
