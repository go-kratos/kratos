package show

import (
	"context"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
)

// Dao is show dao.
type Dao struct {
	// mysql
	db         *sql.DB
	getHead    *sql.Stmt
	getItem    *sql.Stmt
	getHeadTmp *sql.Stmt
	getItemTmp *sql.Stmt
	// redis
	rcmmndRds *redis.Pool
	rcmmndExp int
}

// New new a show dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// mysql
		db: sql.NewMySQL(c.MySQL.Show),
		// redis
		rcmmndRds: redis.NewPool(c.Redis.Recommend.Config),
		rcmmndExp: int(time.Duration(c.Redis.Recommend.Expire) / time.Second),
	}
	d.getHead = d.db.Prepared(_headSQL)
	d.getItem = d.db.Prepared(_itemSQL)
	d.getHeadTmp = d.db.Prepared(_headTmpSQL)
	d.getItemTmp = d.db.Prepared(_itemTmpSQL)
	return d
}

// Close close memcache resource.
func (d *Dao) Close() (err error) {
	if d.rcmmndRds != nil {
		return d.rcmmndRds.Close()
	}
	return nil
}

func (d *Dao) Ping(c context.Context) (err error) {
	conn := d.rcmmndRds.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
