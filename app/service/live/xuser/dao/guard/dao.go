package guard

import (
	"context"

	"go-common/app/service/live/xuser/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
)

// GuardDao vip dao
type GuardDao struct {
	c     *conf.Config
	db    *xsql.DB
	redis *redis.Pool
}

// NewGuardDao init mysql db
func NewGuardDao(c *conf.Config) (dao *GuardDao) {
	dao = &GuardDao{
		c:     c,
		db:    xsql.NewMySQL(c.LiveAppMySQL),
		redis: redis.NewPool(c.GuardRedis),
	}
	return
}

// Close close the resource.
func (d *GuardDao) Close() {
	d.db.Close()
	d.redis.Close()
}

// Ping dao ping
func (d *GuardDao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return nil
}

func (d *GuardDao) getExpire() (respExpire int32) {
	if t := conf.Conf.UserDaHangHaiExpire; t != nil {
		respExpire = t.ExpireTime
	} else {
		respExpire = _emptyExpire
	}
	return
}
