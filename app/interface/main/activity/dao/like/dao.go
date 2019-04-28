package like

import (
	"context"
	"time"

	"go-common/app/interface/main/activity/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

const (
	_lotteryIndex    = "/matsuri/api/mission"
	_lotteryAddTimes = "/matsuri/api/add/times"
	_likeItemURI     = "/activity/likes/list/%d"
	_sourceItemURI   = "/activity/web/view/data/%d"
	_tagsURI         = "/x/internal/tag/archive/multi/tags"
)

// Dao struct
type Dao struct {
	db                 *xsql.DB
	subjectStmt        *xsql.Stmt
	voteLogStmt        *xsql.Stmt
	mc                 *memcache.Pool
	mcLikeExpire       int32
	mcLikeIPExpire     int32
	mcPerpetualExpire  int32
	mcItemExpire       int32
	mcSubStatExpire    int32
	mcViewRankExpire   int32
	mcSourceItemExpire int32
	mcProtocolExpire   int32
	redis              *redis.Pool
	redisExpire        int32
	matchExpire        int32
	followExpire       int32
	hotDotExpire       int32
	randomExpire       int32
	lotteryIndexURL    string
	addLotteryTimesURL string
	likeItemURL        string
	sourceItemURL      string
	tagURL             string
	client             *httpx.Client
	cacheCh            chan func()
	cache              *fanout.Fanout
	es                 *elastic.Elastic
}

// New init
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		db:                 xsql.NewMySQL(c.MySQL.Like),
		mc:                 memcache.NewPool(c.Memcache.Like),
		mcLikeExpire:       int32(time.Duration(c.Memcache.LikeExpire) / time.Second),
		mcLikeIPExpire:     int32(time.Duration(c.Memcache.LikeIPExpire) / time.Second),
		mcPerpetualExpire:  int32(time.Duration(c.Memcache.PerpetualExpire) / time.Second),
		mcItemExpire:       int32(time.Duration(c.Memcache.ItemExpire) / time.Second),
		mcSubStatExpire:    int32(time.Duration(c.Memcache.SubStatExpire) / time.Second),
		mcViewRankExpire:   int32(time.Duration(c.Memcache.ViewRankExpire) / time.Second),
		mcSourceItemExpire: int32(time.Duration(c.Memcache.SourceItemExpire) / time.Second),
		mcProtocolExpire:   int32(time.Duration(c.Memcache.ProtocolExpire) / time.Second),
		redis:              redis.NewPool(c.Redis.Config),
		cacheCh:            make(chan func(), 1024),
		cache:              fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		redisExpire:        int32(time.Duration(c.Redis.Expire) / time.Second),
		matchExpire:        int32(time.Duration(c.Redis.MatchExpire) / time.Second),
		followExpire:       int32(time.Duration(c.Redis.FollowExpire) / time.Second),
		hotDotExpire:       int32(time.Duration(c.Redis.HotDotExpire) / time.Second),
		randomExpire:       int32(time.Duration(c.Redis.RandomExpire) / time.Second),
		lotteryIndexURL:    c.Host.Activity + _lotteryIndex,
		addLotteryTimesURL: c.Host.Activity + _lotteryAddTimes,
		likeItemURL:        c.Host.Activity + _likeItemURI,
		sourceItemURL:      c.Host.Activity + _sourceItemURI,
		tagURL:             c.Host.APICo + _tagsURI,
		client:             httpx.NewClient(c.HTTPClient),
		es:                 elastic.NewElastic(c.Elastic),
	}
	dao.subjectStmt = dao.db.Prepared(_selSubjectSQL)
	dao.voteLogStmt = dao.db.Prepared(_votLogSQL)
	go dao.cacheproc()
	return
}

// CVoteLog chan Vote Log
func (dao *Dao) CVoteLog(c context.Context, sid int64, aid int64, mid int64, stage int64, vote int64) {
	dao.cacheCh <- func() {
		dao.VoteLog(c, sid, aid, mid, stage, vote)
	}
}

// Close Dao
func (dao *Dao) Close() {
	if dao.db != nil {
		dao.db.Close()
	}
	if dao.redis != nil {
		dao.redis.Close()
	}
	if dao.mc != nil {
		dao.mc.Close()
	}
	close(dao.cacheCh)
}

// Ping Dao
func (dao *Dao) Ping(c context.Context) error {
	return dao.db.Ping(c)
}

func (dao *Dao) cacheproc() {
	for {
		f, ok := <-dao.cacheCh
		if !ok {
			return
		}
		f()
	}
}

// PromError stat and log.
func PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}
