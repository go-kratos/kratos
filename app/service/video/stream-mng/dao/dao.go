package dao

import (
	"context"

	"go-common/app/service/video/stream-mng/conf"
	"go-common/library/cache"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"

	"github.com/bluele/gcache"
)

// Dao dao
type Dao struct {
	c          *conf.Config
	mc         *memcache.Pool
	redis      *redis.Pool
	db         *xsql.DB
	tidb       *xsql.DB
	httpClient *bm.Client
	cache      *cache.Cache
	liveAside  *fanout.Fanout
	localCache gcache.Cache

	// DB
	// 新版本主流
	stmtMainStreamCreate              *xsql.Stmt
	stmtMainStreamChangeDefaultVendor *xsql.Stmt
	stmtMainStreamChangeOptions       *xsql.Stmt
	stmtMainStreamClearAllStreaming   *xsql.Stmt
	// 备用流
	stmtBackupStreamCreate *xsql.Stmt
	// 旧版本流
	stmtLegacyStreamCreate            *xsql.Stmt
	stmtLegacyStreamEnableNewUpRank   *xsql.Stmt
	stmtLegacyStreamDisableUpRank     *xsql.Stmt
	stmtLegacyStreamClearStreamFoward *xsql.Stmt
	stmtLegacyStreamNotify            *xsql.Stmt
	// tidb
	stmtUpStreamDispatch *xsql.Stmt
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:          c,
		mc:         memcache.NewPool(c.Memcache),
		redis:      redis.NewPool(c.Redis),
		db:         xsql.NewMySQL(c.MySQL),
		tidb:       xsql.NewMySQL(c.TIDB),
		httpClient: bm.NewClient(c.HTTPClient),
		cache:      cache.New(1, 10240),
		liveAside:  fanout.New("stream-mng"),
		localCache: gcache.New(10240).Simple().Build(),
	}

	// 新版本主流
	dao.stmtMainStreamCreate = dao.db.Prepared(_insertMainStream)
	dao.stmtMainStreamChangeDefaultVendor = dao.db.Prepared(_changeDefaultVendor)
	dao.stmtMainStreamChangeOptions = dao.db.Prepared(_changeOptions)
	dao.stmtMainStreamClearAllStreaming = dao.db.Prepared(_clearAllStreaming)
	// 备用流
	dao.stmtBackupStreamCreate = dao.db.Prepared(_insertBackupStream)
	// 旧版本流
	dao.stmtLegacyStreamCreate = dao.db.Prepared(_insertOfficialStream)
	dao.stmtLegacyStreamEnableNewUpRank = dao.db.Prepared(_updateUpOfficialStreamStatus)
	dao.stmtLegacyStreamDisableUpRank = dao.db.Prepared(_updateForwardOfficialStreamStatus)
	dao.stmtLegacyStreamClearStreamFoward = dao.db.Prepared(_updateOfficalStreamUpRankStatus)
	dao.stmtLegacyStreamNotify = dao.db.Prepared(_setOriginStreamingStatus)
	// tidb
	dao.stmtUpStreamDispatch = dao.tidb.Prepared(_insertUpStreamInfo)
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
	d.liveAside.Close()
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}
