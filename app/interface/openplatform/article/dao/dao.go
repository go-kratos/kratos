package dao

import (
	"context"
	xhttp "net/http"
	"runtime"
	"time"

	"go-common/app/interface/openplatform/article/conf"
	"go-common/library/cache"
	"go-common/library/cache/memcache"
	xredis "go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"

	hbase "go-common/library/database/hbase.v2"
)

var (
	errorsCount = prom.BusinessErrCount
	infosCount  = prom.BusinessInfoCount
	cachedCount = prom.CacheHit
	missedCount = prom.CacheMiss
)

// PromError prom error
func PromError(name string) {
	errorsCount.Incr(name)
}

// promErrorCheck check prom error
func promErrorCheck(err error) {
	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		name := "d." + runtime.FuncForPC(pc).Name()
		PromError(name)
		log.Error("%s err: %+v", name, err)
	}
}

// PromInfo add prom info
func PromInfo(name string) {
	infosCount.Incr(name)
}

// Dao dao
type Dao struct {
	// config
	c *conf.Config
	// db
	articleDB *sql.DB
	// http client
	httpClient        *bm.Client
	messageHTTPClient *bm.Client
	bfsClient         *xhttp.Client
	// memcache
	mc *memcache.Pool
	//redis
	redis               *xredis.Pool
	mcArticleExpire     int32
	mcStatsExpire       int32
	mcLikeExpire        int32
	mcCardsExpire       int32
	mcListArtsExpire    int32
	mcListExpire        int32
	mcArtListExpire     int32
	mcUpListsExpire     int32
	mcSubExp            int32
	mcListReadExpire    int32
	mcHotspotExpire     int32
	mcAuthorExpire      int32
	mcArticlesIDExpire  int32
	mcArticleTagExpire  int32
	mcUpStatDailyExpire int32
	redisUpperExpire    int32
	redisSortExpire     int64
	redisSortTTL        int64
	redisArtLikesExpire int32
	redisRankExpire     int64
	redisRankTTL        int64
	redisMaxLikeExpire  int64
	redisHotspotExpire  int64
	redisReadPingExpire int64
	redisReadSetExpire  int64
	// stmt
	categoriesStmt                *sql.Stmt
	authorsStmt                   *sql.Stmt
	applyStmt                     *sql.Stmt
	addAuthorStmt                 *sql.Stmt
	applyCountStmt                *sql.Stmt
	articleMetaStmt               *sql.Stmt
	allArticleMetaStmt            *sql.Stmt
	articleContentStmt            *sql.Stmt
	articleUpperCountStmt         *sql.Stmt
	updateArticleStateStmt        *sql.Stmt
	upPassedStmt                  *sql.Stmt
	recommendCategoryStmt         *sql.Stmt
	delRecommendStmt              *sql.Stmt
	allRecommendStmt              *sql.Stmt
	allRecommendCountStmt         *sql.Stmt
	newestArtsMetaStmt            *sql.Stmt
	upperArtCntCreationStmt       *sql.Stmt
	articleMetaCreationStmt       *sql.Stmt
	articleUpCntTodayStmt         *sql.Stmt
	addComplaintStmt              *sql.Stmt
	complaintExistStmt            *sql.Stmt
	complaintProtectStmt          *sql.Stmt
	addComplaintCountStmt         *sql.Stmt
	settingsStmt                  *sql.Stmt
	authorStmt                    *sql.Stmt
	noticeStmt                    *sql.Stmt
	userNoticeStmt                *sql.Stmt
	updateUserNoticeStmt          *sql.Stmt
	creativeListsStmt             *sql.Stmt
	creativeListAddStmt           *sql.Stmt
	creativeListDelStmt           *sql.Stmt
	creativeListUpdateStmt        *sql.Stmt
	creativeListUpdateTimeStmt    *sql.Stmt
	creativeListArticlesStmt      *sql.Stmt
	listStmt                      *sql.Stmt
	creativeListDelAllArticleStmt *sql.Stmt
	creativeListAddArticleStmt    *sql.Stmt
	creativeListDelArticleStmt    *sql.Stmt
	allListStmt                   *sql.Stmt
	hotspotsStmt                  *sql.Stmt
	searchArtsStmt                *sql.Stmt
	addCheatStmt                  *sql.Stmt
	delCheatStmt                  *sql.Stmt
	// databus
	statDbus *databus.Databus
	// inteval
	UpdateRecommendsInterval int64
	UpdateBannersInterval    int64
	// hbase
	hbase *hbase.Client
	//cache
	cache *cache.Cache
}

// New dao new
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// config
		c: c,
		// http client
		httpClient:        bm.NewClient(c.HTTPClient),
		messageHTTPClient: bm.NewClient(c.MessageHTTPClient),
		bfsClient:         &xhttp.Client{Timeout: time.Duration(c.BFS.Timeout)},
		// mc
		mc:                  memcache.NewPool(c.Memcache.Config),
		mcArticleExpire:     int32(time.Duration(c.Memcache.ArticleExpire) / time.Second),
		mcStatsExpire:       int32(time.Duration(c.Memcache.StatsExpire) / time.Second),
		mcLikeExpire:        int32(time.Duration(c.Memcache.LikeExpire) / time.Second),
		mcCardsExpire:       int32(time.Duration(c.Memcache.CardsExpire) / time.Second),
		mcSubExp:            int32(time.Duration(c.Memcache.SubmitExpire) / time.Second),
		mcListArtsExpire:    int32(time.Duration(c.Memcache.ListArtsExpire) / time.Second),
		mcListExpire:        int32(time.Duration(c.Memcache.ListExpire) / time.Second),
		mcArtListExpire:     int32(time.Duration(c.Memcache.ArtListExpire) / time.Second),
		mcUpListsExpire:     int32(time.Duration(c.Memcache.UpListsExpire) / time.Second),
		mcListReadExpire:    int32(time.Duration(c.Memcache.ListReadExpire) / time.Second),
		mcHotspotExpire:     int32(time.Duration(c.Memcache.HotspotExpire) / time.Second),
		mcAuthorExpire:      int32(time.Duration(c.Memcache.AuthorExpire) / time.Second),
		mcArticlesIDExpire:  int32(time.Duration(c.Memcache.ArticlesIDExpire) / time.Second),
		mcArticleTagExpire:  int32(time.Duration(c.Memcache.ArticleTagExpire) / time.Second),
		mcUpStatDailyExpire: int32(time.Duration(c.Memcache.UpStatDailyExpire) / time.Second),
		//redis
		redis:               xredis.NewPool(c.Redis),
		redisUpperExpire:    int32(time.Duration(c.Article.ExpireUpper) / time.Second),
		redisSortExpire:     int64(time.Duration(c.Article.ExpireSortArts) / time.Second),
		redisSortTTL:        int64(time.Duration(c.Article.TTLSortArts) / time.Second),
		redisArtLikesExpire: int32(time.Duration(c.Article.ExpireArtLikes) / time.Second),
		redisRankExpire:     int64(time.Duration(c.Article.ExpireRank) / time.Second),
		redisRankTTL:        int64(time.Duration(c.Article.TTLRank) / time.Second),
		redisMaxLikeExpire:  int64(time.Duration(c.Article.ExpireMaxLike) / time.Second),
		redisHotspotExpire:  int64(time.Duration(c.Article.ExpireHotspot) / time.Second),
		redisReadPingExpire: int64(time.Duration(c.Article.ExpireReadPing) / time.Second),
		redisReadSetExpire:  int64(time.Duration(c.Article.ExpireReadSet) / time.Second),
		// db
		articleDB: sql.NewMySQL(c.MySQL.Article),
		// prom
		statDbus:                 databus.New(c.StatDatabus),
		UpdateRecommendsInterval: int64(time.Duration(c.Article.UpdateRecommendsInteval) / time.Second),
		UpdateBannersInterval:    int64(time.Duration(c.Article.UpdateBannersInteval) / time.Second),
		// hbase
		hbase: hbase.NewClient(c.HBase),
		cache: cache.New(1, 1024),
	}
	d.categoriesStmt = d.articleDB.Prepared(_categoriesSQL)
	d.authorsStmt = d.articleDB.Prepared(_authorsSQL)
	d.applyStmt = d.articleDB.Prepared(_applySQL)
	d.addAuthorStmt = d.articleDB.Prepared(_addAuthorSQL)
	d.applyCountStmt = d.articleDB.Prepared(_applyCountSQL)
	d.articleMetaStmt = d.articleDB.Prepared(_articleMetaSQL)
	d.allArticleMetaStmt = d.articleDB.Prepared(_allArticleMetaSQL)
	d.articleContentStmt = d.articleDB.Prepared(_articleContentSQL)
	d.updateArticleStateStmt = d.articleDB.Prepared(_updateArticleStateSQL)
	d.upPassedStmt = d.articleDB.Prepared(_upperPassedSQL)
	d.recommendCategoryStmt = d.articleDB.Prepared(_recommendCategorySQL)
	d.allRecommendStmt = d.articleDB.Prepared(_allRecommendSQL)
	d.allRecommendCountStmt = d.articleDB.Prepared(_allRecommendCountSQL)
	d.delRecommendStmt = d.articleDB.Prepared(_deleteRecommendSQL)
	d.newestArtsMetaStmt = d.articleDB.Prepared(_newestArtsMetaSQL)
	d.upperArtCntCreationStmt = d.articleDB.Prepared(_upperArticleCountCreationSQL)
	d.articleUpperCountStmt = d.articleDB.Prepared(_articleUpperCountSQL)
	d.articleMetaCreationStmt = d.articleDB.Prepared(_articleMetaCreationSQL)
	d.articleUpCntTodayStmt = d.articleDB.Prepared(_articleUpCntTodaySQL)
	d.addComplaintStmt = d.articleDB.Prepared(_addComplaintsSQL)
	d.complaintExistStmt = d.articleDB.Prepared(_complaintExistSQL)
	d.complaintProtectStmt = d.articleDB.Prepared(_complaintProtectSQL)
	d.addComplaintCountStmt = d.articleDB.Prepared(_addComplaintCountSQL)
	d.settingsStmt = d.articleDB.Prepared(_settingsSQL)
	d.authorStmt = d.articleDB.Prepared(_authorSQL)
	d.noticeStmt = d.articleDB.Prepared(_noticeSQL)
	d.userNoticeStmt = d.articleDB.Prepared(_userNoticeSQL)
	d.updateUserNoticeStmt = d.articleDB.Prepared(_updateUserNoticeSQL)
	d.creativeListsStmt = d.articleDB.Prepared(_creativeListsSQL)
	d.creativeListAddStmt = d.articleDB.Prepared(_creativeListAddSQL)
	d.creativeListDelStmt = d.articleDB.Prepared(_creativeListDelSQL)
	d.creativeListUpdateStmt = d.articleDB.Prepared(_creativeListUpdateSQL)
	d.creativeListUpdateTimeStmt = d.articleDB.Prepared(_creativeListUpdateTimeSQL)
	d.creativeListArticlesStmt = d.articleDB.Prepared(_creativeListArticlesSQL)
	d.creativeListDelAllArticleStmt = d.articleDB.Prepared(_creativeListDelAllArticleSQL)
	d.creativeListAddArticleStmt = d.articleDB.Prepared(_creativeListAddArticleSQL)
	d.listStmt = d.articleDB.Prepared(_listSQL)
	d.creativeListDelArticleStmt = d.articleDB.Prepared(_creativeListDelArticleSQL)
	d.allListStmt = d.articleDB.Prepared(_allListsSQL)
	d.hotspotsStmt = d.articleDB.Prepared(_hotspotsSQL)
	d.searchArtsStmt = d.articleDB.Prepared(_searchArticles)
	d.addCheatStmt = d.articleDB.Prepared(_addCheatSQL)
	d.delCheatStmt = d.articleDB.Prepared(_delCheatSQL)
	return d
}

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.articleDB.Begin(c)
}

// Ping check connection success.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		PromError("mc:Ping")
		log.Error("d.pingMC error(%+v)", err)
		return
	}
	if err = d.pingRedis(c); err != nil {
		PromError("redis:Ping")
		log.Error("d.pingRedis error(%+v)", err)
		return
	}
	if err = d.articleDB.Ping(c); err != nil {
		PromError("db:Ping")
		log.Error("d.articleDB.Ping error(%+v)", err)
	}
	return
}

// Close close  resource.
func (d *Dao) Close() {
	d.articleDB.Close()
	d.mc.Close()
	d.redis.Close()
}
