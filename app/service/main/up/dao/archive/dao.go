package archive

import (
	"time"

	"go-common/app/service/main/up/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/stat/prom"
)

// Dao is archive dao.
type Dao struct {
	c         *conf.Config
	resultDB  *sql.DB
	archiveDB *sql.DB
	// redis
	upRds    *redis.Pool
	upExpire int32
	errProm  *prom.Prom
	infoProm *prom.Prom
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// db
		resultDB:  sql.NewMySQL(c.DB.ArcResult),
		archiveDB: sql.NewMySQL(c.DB.Archive),
		// redis
		upRds:    redis.NewPool(c.Redis.Up.Config),
		upExpire: int32(time.Duration(c.Redis.Up.UpExpire) / time.Second),
		errProm:  prom.BusinessErrCount,
		infoProm: prom.BusinessInfoCount,
	}
	return
}
