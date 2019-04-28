package audit

import (
	"go-common/app/interface/main/tv/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
)

// Dao is account dao.
type Dao struct {
	mc   *memcache.Pool
	conf *conf.Config
	db   *sql.DB
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf: c,
		mc:   memcache.NewPool(c.Memcache.Config),
		db:   sql.NewMySQL(c.Mysql),
	}
	return
}
