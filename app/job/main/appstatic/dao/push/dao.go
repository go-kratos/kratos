package push

import (
	appres "go-common/app/interface/main/app-resource/api/v1"
	"go-common/app/job/main/appstatic/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao .
type Dao struct {
	c         *conf.Config
	db        *xsql.DB
	client    *bm.Client
	redis     *redis.Pool
	appresCli appres.AppResourceClient
}

// New creates a dao instance.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		db:     xsql.NewMySQL(c.MySQL),
		client: bm.NewClient(c.HTTPClient),
		redis:  redis.NewPool(c.Redis),
	}
	var err error
	if d.appresCli, err = appres.NewClient(c.AppresClient); err != nil {
		panic(err)
	}
	return
}
