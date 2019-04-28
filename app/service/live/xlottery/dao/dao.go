package dao

import (
	"context"
	"time"

	"math/rand"

	"go-common/app/service/live/xlottery/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

// Dao dao
type Dao struct {
	c *conf.Config
	//mc    *memcache.Pool
	redis *redis.Pool
	db    *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		redis: redis.NewPool(c.Redis.Lottery),
		db:    xsql.NewMySQL(c.Database.Lottery),
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
	return d.db.Ping(c)
}

//　通过提供的sql和bind来更新，没有实际业务意义，只是为了少写重复代码
func (d *Dao) execSqlWithBindParams(c context.Context, sql *string, bindParams ...interface{}) (affect int64, err error) {
	res, err := d.db.Exec(c, *sql, bindParams...)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", *sql, err)
		return
	}
	return res.RowsAffected()
}

func randomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
