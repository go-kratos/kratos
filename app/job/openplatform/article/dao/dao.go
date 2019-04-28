package dao

import (
	"context"
	"time"

	"go-common/app/job/openplatform/article/conf"
	"go-common/library/cache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

// Dao .
type Dao struct {
	c                             *conf.Config
	db                            *xsql.DB
	redis                         *redis.Pool
	artRedis                      *redis.Pool
	httpClient                    *bm.Client
	gameHTTPClient                *bm.Client
	viewCacheTTL, gameCacheExpire int64
	dupViewCacheTTL               int64
	redisSortExpire               int64
	redisSortTTL                  int64
	// stmt
	updateSearchStmt      *xsql.Stmt
	delSearchStmt         *xsql.Stmt
	updateSearchStatsStmt *xsql.Stmt
	gameStmt              *xsql.Stmt
	cheatStmt             *xsql.Stmt
	newestArtsMetaStmt    *xsql.Stmt
	searchArtsStmt        *xsql.Stmt
	updateRecheckStmt     *xsql.Stmt
	getRecheckStmt        *xsql.Stmt
	settingsStmt          *xsql.Stmt
	midByPubtimeStmt      *xsql.Stmt
	statByMidStmt         *xsql.Stmt
	dynamicDbus           *databus.Databus
	cache                 *cache.Cache
}

var (
	errorsCount = prom.BusinessErrCount
	cacheLen    = prom.BusinessInfoCount
	infosCount  = prom.BusinessInfoCount
)

// New creates a dao instance.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:               c,
		db:              xsql.NewMySQL(c.DB),
		redis:           redis.NewPool(c.Redis),
		artRedis:        redis.NewPool(c.ArtRedis),
		httpClient:      bm.NewClient(c.HTTPClient),
		gameHTTPClient:  bm.NewClient(c.GameHTTPClient),
		viewCacheTTL:    int64(time.Duration(c.Job.ViewCacheTTL) / time.Second),
		dupViewCacheTTL: int64(time.Duration(c.Job.DupViewCacheTTL) / time.Second),
		gameCacheExpire: int64(time.Duration(c.Job.GameCacheExpire) / time.Second),
		redisSortExpire: int64(time.Duration(c.Job.ExpireSortArts) / time.Second),
		redisSortTTL:    int64(time.Duration(c.Job.TTLSortArts) / time.Second),
		dynamicDbus:     databus.New(c.DynamicDbus),
		cache:           cache.New(1, 1024),
	}
	d.updateSearchStmt = d.db.Prepared(_updateSearch)
	d.delSearchStmt = d.db.Prepared(_delSearch)
	d.updateSearchStatsStmt = d.db.Prepared(_updateSearchStats)
	d.newestArtsMetaStmt = d.db.Prepared(_newestArtsMetaSQL)
	d.gameStmt = d.db.Prepared(_gameList)
	d.cheatStmt = d.db.Prepared(_allCheat)
	d.searchArtsStmt = d.db.Prepared(_searchArticles)
	d.updateRecheckStmt = d.db.Prepared(_updateCheckState)
	d.getRecheckStmt = d.db.Prepared(_checkStateSQL)
	d.settingsStmt = d.db.Prepared(_settingsSQL)
	d.midByPubtimeStmt = d.db.Prepared(_midsByPublishTimeSQL)
	d.statByMidStmt = d.db.Prepared(_statByMidSQL)
	return
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	GameList(c context.Context) (res []int64, err error)
}

// PromError prometheus error count.
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo prometheus info count.
func PromInfo(name string) {
	infosCount.Incr(name)
}

// Ping reports the health of the db/cache etc.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		PromError("db:Ping")
		return
	}
	err = d.pingRedis(c)
	return
}
