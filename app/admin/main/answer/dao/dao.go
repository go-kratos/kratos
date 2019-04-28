package dao

import (
	"time"

	"go-common/app/admin/main/answer/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	"go-common/library/database/orm"

	"github.com/jinzhu/gorm"
)

// Dao struct info of Dao.
type Dao struct {
	c           *conf.Config
	db          *gorm.DB
	es          *elastic.Elastic
	redis       *redis.Pool
	redisExpire int32
}

// TextImgConf text img conf.
type TextImgConf struct {
	Fontsize    int
	Length      int
	Ansfontsize int
	Spacing     float64
	Ansspacing  float64
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		db:          orm.NewMySQL(c.Mysql),
		es:          elastic.NewElastic(nil),
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
}
