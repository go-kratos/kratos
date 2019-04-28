package dao

import (
	"context"
	"time"

	"go-common/app/service/main/thumbup/conf"
	"go-common/app/service/main/thumbup/model"
	"go-common/library/cache/memcache"
	xredis "go-common/library/cache/redis"
	"go-common/library/database/tidb"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

// PromError prom error
func PromError(name string) {
	prom.BusinessErrCount.Incr(name)
}

// Dao dao
type Dao struct {
	// config
	c *conf.Config
	// db
	tidb *tidb.DB
	// memcache
	mc            *memcache.Pool
	mcStatsExpire int32
	//redis
	redis                *xredis.Pool
	redisStatsExpire     int64
	redisUserLikesExpire int64
	redisItemLikesExpire int64
	// redisSortExpire int64
	// stmt
	businessesStmt        *tidb.Stmts
	likeStateStmt         *tidb.Stmts
	userLikeCountStmt     *tidb.Stmts
	itemLikeListStmt      *tidb.Stmts
	userLikeListStmt      *tidb.Stmts
	statsOriginStmt       *tidb.Stmts
	statStmt              *tidb.Stmts
	updateLikeStmt        *tidb.Stmts
	updateCountChangeStmt *tidb.Stmts
	statDbus              *databus.Databus
	likeDbus              *databus.Databus
	itemDbus              *databus.Databus
	userDbus              *databus.Databus
	cache                 *fanout.Fanout
	async                 *fanout.Fanout
	tidbAsync             *fanout.Fanout
	BusinessMap           map[string]*model.Business
	BusinessIDMap         map[int64]*model.Business
}

// New dao new
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// config
		c: c,
		// mc
		mc:            memcache.NewPool(c.Memcache.Config),
		mcStatsExpire: int32(time.Duration(c.Memcache.StatsExpire) / time.Second),
		// redis
		redis:                xredis.NewPool(c.Redis.Config),
		redisStatsExpire:     int64(time.Duration(c.Redis.StatsExpire) / time.Second),
		redisUserLikesExpire: int64(time.Duration(c.Redis.UserLikesExpire) / time.Second),
		redisItemLikesExpire: int64(time.Duration(c.Redis.ItemLikesExpire) / time.Second),
		// db
		tidb:      tidb.NewTiDB(c.Tidb),
		statDbus:  databus.New(c.StatDatabus),
		likeDbus:  databus.New(c.LikeDatabus),
		itemDbus:  databus.New(c.ItemDatabus),
		userDbus:  databus.New(c.UserDatabus),
		cache:     fanout.New("cache", fanout.Worker(1), fanout.Buffer(10240)),
		async:     fanout.New("async", fanout.Worker(2), fanout.Buffer(10240)),
		tidbAsync: fanout.New("tidb-async", fanout.Worker(10), fanout.Buffer(10240)),
	}
	d.businessesStmt = d.tidb.Prepared(_tidbBusinessesSQL)
	d.likeStateStmt = d.tidb.Prepared(_tidbLikeMidSQL)
	d.userLikeCountStmt = d.tidb.Prepared(_tidbUserLikeCountSQL)
	d.itemLikeListStmt = d.tidb.Prepared(_tidbItemLikeListSQL)
	d.userLikeListStmt = d.tidb.Prepared(_tidbUserLikeListSQL)
	d.statsOriginStmt = d.tidb.Prepared(_tidbStatsOriginSQL)
	d.statStmt = d.tidb.Prepared(_tidbStatSQL)
	d.updateLikeStmt = d.tidb.Prepared(_tidbUpdateLikeSQL)
	d.updateCountChangeStmt = d.tidb.Prepared(_tidbupdateCountChange)
	d.loadBusiness()
	go d.loadBusinessproc()
	return d
}

// Ping check connection success.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		PromError("mc:Ping")
		log.Error("d.pingMC error(%v)", err)
		return
	}
	if err = d.pingRedis(c); err != nil {
		PromError("redis:Ping")
		log.Error("d.pingRedis error(%v)", err)
		return
	}
	return
}

// Close close  resource.
func (d *Dao) Close() {
	d.async.Close()
	d.tidbAsync.Close()
	d.tidb.Close()
	d.mc.Close()
	d.redis.Close()
}

// pingMc ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 100}
	err = conn.Set(&item)
	return
}

// pingRedis ping redis.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		PromError("redis: ping remote")
		log.Error("remote redis: conn.Do(SET,PING,PONG) error(%v)", err)
	}
	conn.Close()
	return
}

// LoadBusiness .
func (d *Dao) loadBusiness() {
	var business []*model.Business
	var err error
	businessMap := make(map[string]*model.Business)
	businessIDMap := make(map[int64]*model.Business)
	for {
		if business, err = d.Businesses(context.TODO()); err != nil {
			time.Sleep(time.Second)
			continue
		}
		for _, b := range business {
			businessMap[b.Name] = b
			businessIDMap[b.ID] = b
		}
		d.BusinessMap = businessMap
		d.BusinessIDMap = businessIDMap
		return
	}
}

func (d *Dao) loadBusinessproc() {
	for {
		time.Sleep(time.Minute * 5)
		d.loadBusiness()
	}
}
