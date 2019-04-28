package monitor

import (
	"context"
	"go-common/app/job/main/aegis/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

type Dao struct {
	c           *conf.Config
	redis       *redis.Pool
	db          *xsql.DB
	http        *bm.Client
	URLArcAddit string
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:           c,
		redis:       redis.NewPool(c.Redis),
		db:          xsql.NewMySQL(c.MySQL.Fast),
		http:        bm.NewClient(c.HTTP.Fast),
		URLArcAddit: c.Host.Videoup + _arcAdditURL,
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}
