package dao

import (
	"context"
	"time"

	"go-common/app/job/main/figure/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
)

// Dao figure job dao
type Dao struct {
	c              *conf.Config
	hbase          *hbase.Client
	db             *sql.DB
	redis          *redis.Pool
	redisExpire    int
	waiteMidExpire int
}

// New new a figure DAO
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:              c,
		hbase:          hbase.NewClient(c.HBase.Config),
		db:             sql.NewMySQL(c.Mysql),
		redis:          redis.NewPool(c.Redis.Config),
		redisExpire:    int(time.Duration(c.Redis.Expire) / time.Second),
		waiteMidExpire: int(time.Duration(c.Redis.WaiteMidExpire) / time.Second),
	}
	return
}

// Ping check service health
func (d *Dao) Ping(c context.Context) (err error) {
	return d.PingRedis(c)
}

// Close close all dao.
func (d *Dao) Close() {
	if d.hbase != nil {
		d.hbase.Close()
	}
}

//Version get ever monday start time ts.
func (d *Dao) Version(now time.Time) (ts int64) {
	var (
		n int8
	)
	y, m, day := now.Date()
	w := now.Weekday()
	switch w {
	case time.Sunday:
		n = 6
	default:
		n = int8(w) - 1
	}
	t := time.Date(y, m, day, 0, 0, 0, 0, time.Local).Add(-time.Duration(n) * 24 * time.Hour)
	return t.Unix()
}
