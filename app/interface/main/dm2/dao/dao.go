package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/dm2/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/bfs"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

var (
	errorsCount = prom.BusinessErrCount
	missedCount = prom.CacheMiss
	cachedCount = prom.CacheHit
)

// Dao dm dao.
type Dao struct {
	conf *conf.Config
	// mysql
	dmWriter *sql.DB
	dmReader *sql.DB
	dbDM     *sql.DB
	// redis
	dmRds       *redis.Pool
	dmRdsExpire int32
	// recent dm redis
	dmRctRds    *redis.Pool
	dmRctExpire int32
	// segment dm redis
	dmSegRds    *redis.Pool
	dmSegExpire int32
	// memcache
	dmMC              *memcache.Pool
	dmExpire          int32
	subjectExpire     int32
	historyExpire     int32
	ajaxExpire        int32
	dmMaskExpire      int32
	filterMC          *memcache.Pool
	filterMCExpire    int32
	dmSegMC           *memcache.Pool
	dmSegMCExpire     int32
	dmLimiterMCExpire int32
	// http
	httpCli *bm.Client
	// databus
	databus   *databus.Databus
	actionPub *databus.Databus
	// elastic
	elastic *elastic.Elastic
	// bfsCli
	bfsCli *bfs.BFS
	// subtitle mc
	subtitleMc       *memcache.Pool
	subtitleMcExpire int32
	subtitleCheckPub *databus.Databus
}

// New new a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf: c,
		// mysql
		dmWriter: sql.NewMySQL(c.DB.DMWriter),
		dmReader: sql.NewMySQL(c.DB.DMReader),
		dbDM:     sql.NewMySQL(c.DB.DM),
		// redis
		dmRds:       redis.NewPool(c.Redis.DM.Config),
		dmRdsExpire: int32(time.Duration(c.Redis.DM.Expire) / time.Second),
		// recent dm redis
		dmRctRds:    redis.NewPool(c.Redis.DMRct.Config),
		dmRctExpire: int32(time.Duration(c.Redis.DMRct.Expire) / time.Second),
		// segment dm redis
		dmSegRds:    redis.NewPool(c.Redis.DMSeg.Config),
		dmSegExpire: int32(time.Duration(c.Redis.DMSeg.Expire) / time.Second),
		// memcache
		dmMC:              memcache.NewPool(c.Memcache.DM.Config),
		dmExpire:          int32(time.Duration(c.Memcache.DM.DMExpire) / time.Second),
		subjectExpire:     int32(time.Duration(c.Memcache.DM.SubjectExpire) / time.Second),
		historyExpire:     int32(time.Duration(c.Memcache.DM.HistoryExpire) / time.Second),
		ajaxExpire:        int32(time.Duration(c.Memcache.DM.AjaxExpire) / time.Second),
		dmMaskExpire:      int32(time.Duration(c.Memcache.DM.DMMaskExpire) / time.Second),
		filterMC:          memcache.NewPool(c.Memcache.Filter.Config),
		filterMCExpire:    int32(time.Duration(c.Memcache.Filter.Expire) / time.Second),
		dmSegMC:           memcache.NewPool(c.Memcache.DMSeg.Config),
		dmSegMCExpire:     int32(time.Duration(c.Memcache.DMSeg.DMExpire) / time.Second),
		dmLimiterMCExpire: int32(time.Duration(c.Memcache.DMSeg.DMLimiterExpire) / time.Second),
		// http
		httpCli: bm.NewClient(c.HTTPCli),
		// databus
		databus:   databus.New(c.Databus),
		actionPub: databus.New(c.ActionPub),
		// elastic
		elastic: elastic.NewElastic(c.Elastic),
		// bfscli
		bfsCli: bfs.New(c.Bfs.Client),
		// subtitle MC
		subtitleMc:       memcache.NewPool(c.Memcache.Subtitle.Config),
		subtitleMcExpire: int32(time.Duration(c.Memcache.Subtitle.Expire) / time.Second),
		subtitleCheckPub: databus.New(c.SubtitleCheckPub),
	}
	return
}

func (d *Dao) hitSubject(oid int64) int64 {
	return oid % _subjectSharding
}

func (d *Dao) hitIndex(oid int64) int64 {
	return oid % _indexSharding
}

func (d *Dao) hitContent(dmid int64) int64 {
	return dmid % _contentSharding
}

// BeginBiliDMTrans begin db transaction.
func (d *Dao) BeginBiliDMTrans(c context.Context) (*sql.Tx, error) {
	return d.dbDM.Begin(c)
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.dmReader.Ping(c); err != nil {
		log.Error("d.dmReader error(%v)", err)
		return
	}
	if err = d.dmWriter.Ping(c); err != nil {
		log.Error("d.dmWriter error(%v)", err)
		return
	}
	// mc
	dmMC := d.dmMC.Get(c)
	defer dmMC.Close()
	if err = dmMC.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("dmMC.Set error(%v)", err)
		return
	}
	filterMC := d.filterMC.Get(c)
	defer filterMC.Close()
	if err = filterMC.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("filterMC.Set error(%v)", err)
		return
	}
	dmSegMC := d.dmSegMC.Get(c)
	defer dmSegMC.Close()
	if err = dmSegMC.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("dmSegMC.Set error(%v)", err)
		return
	}
	// dm redis
	dmConn := d.dmRds.Get(c)
	defer dmConn.Close()
	if _, err = dmConn.Do("SET", "ping", "pong"); err != nil {
		log.Error("dmConn.Do(SET) error(%v)", err)
		return
	}
	rctRdsConn := d.dmRctRds.Get(c)
	defer rctRdsConn.Close()
	if _, err = rctRdsConn.Do("SET", "ping", "pong"); err != nil {
		log.Error("rctRdsConn.Do(SET) error(%v)", err)
		return
	}
	// segment dm redis
	segRdsConn := d.dmSegRds.Get(c)
	defer segRdsConn.Close()
	if _, err = segRdsConn.Do("SET", "ping", "pong"); err != nil {
		log.Error("segRdsConn.Do(SET) error(%v)", err)
	}
	return
}

// PromError prom error
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromCacheHit prom cache hit
func PromCacheHit(name string, v int64) {
	cachedCount.Add(name, v)
}

// PromCacheMiss prom cache hit
func PromCacheMiss(name string, v int64) {
	missedCount.Add(name, v)
}
