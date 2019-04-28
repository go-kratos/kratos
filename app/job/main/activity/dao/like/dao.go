package like

import (
	"context"
	"time"

	"go-common/app/job/main/activity/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	"go-common/library/net/http/blademaster"
)

const _activity = "activity"

// Dao  dao
type Dao struct {
	db                 *sql.DB
	subjectStmt        *sql.Stmt
	inOnlineLog        *sql.Stmt
	mcLike             *memcache.Pool
	mcLikeExpire       int32
	redis              *redis.Pool
	redisExpire        int32
	httpClient         *blademaster.Client
	es                 *elastic.Elastic
	setObjStatURL      string
	setViewRankURL     string
	setLikeContentURL  string
	addLotteryTimesURL string
}

// New init
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:                 sql.NewMySQL(c.MySQL.Like),
		mcLike:             memcache.NewPool(c.Memcache.Like),
		mcLikeExpire:       int32(time.Duration(c.Memcache.LikeExpire) / time.Second),
		redis:              redis.NewPool(c.Redis.Config),
		redisExpire:        int32(time.Duration(c.Redis.Expire) / time.Second),
		httpClient:         blademaster.NewClient(c.HTTPClient),
		es:                 elastic.NewElastic(c.Elastic),
		setObjStatURL:      c.Host.APICo + _setObjStatURI,
		setViewRankURL:     c.Host.APICo + _setViewRankURI,
		setLikeContentURL:  c.Host.APICo + _setLikeContentURI,
		addLotteryTimesURL: c.Host.Activity + _addLotteryTimesURI,
	}
	d.subjectStmt = d.db.Prepared(_selSubjectSQL)
	d.inOnlineLog = d.db.Prepared(_inOnlineLogSQL)
	return
}

// Close close
func (d *Dao) Close() {
	d.db.Close()
}

// Ping ping
func (d *Dao) Ping(c context.Context) error {
	return d.db.Ping(c)
}
