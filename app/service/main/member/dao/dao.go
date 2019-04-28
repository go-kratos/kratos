package dao

import (
	"context"
	"time"

	"go-common/app/service/main/member/conf"
	"go-common/app/service/main/member/dao/block"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
	xtime "go-common/library/time"
)

const (
	_searchLogURI = "/x/admin/search/log"
	_deleteLogURI = "/x/admin/search/log/delete"
)

// Dao struct info of Dao.
type Dao struct {
	*cacheTTL
	accdb  *sql.DB
	db     *sql.DB
	mc     *memcache.Pool
	c      *conf.Config
	client *bm.Client
	redis  *redis.Pool
	cache  *fanout.Fanout
	// databus
	logDatabus *databus.Databus
	accNotify  *databus.Databus

	block *block.Dao
}

type cacheTTL struct {
	baseTTL            int32
	moralTTL           int32
	captureTimesTTL    int32
	captureCodeTTL     int32
	captureErrTimesTTL int32
	applyInfoTTL       int32
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		accdb:      sql.NewMySQL(c.AccMysql),
		db:         sql.NewMySQL(c.Mysql),
		mc:         memcache.NewPool(c.Memcache),
		client:     bm.NewClient(c.HTTPClient),
		redis:      redis.NewPool(c.Redis),
		cache:      fanout.New("memberServiceCache", fanout.Worker(1), fanout.Buffer(10240)),
		logDatabus: databus.New(c.Databus),
		accNotify:  databus.New(c.AccountNotify),
	}
	d.block = block.New(
		c,
		sql.NewMySQL(c.BlockMySQL),
		memcache.NewPool(c.BlockMemcache),
		d.client,
		d.NotifyPurgeCache,
	)
	d.cacheTTL = newCacheTTL(c.CacheTTL)
	return
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) error {
	if err := d.db.Ping(c); err != nil {
		log.Error("Failed to ping database: %+v", err)
		return err
	}
	return d.pingMC(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.mc != nil {
		d.mc.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
	if d.block != nil {
		d.block.Close()
	}
}

func durationToSeconds(expire xtime.Duration) int32 {
	return int32(time.Duration(expire) / time.Second)
}

func newCacheTTL(c *conf.CacheTTL) *cacheTTL {
	return &cacheTTL{
		baseTTL:            durationToSeconds(c.BaseTTL),
		moralTTL:           durationToSeconds(c.MoralTTL),
		captureTimesTTL:    durationToSeconds(c.CaptureErrTimesTTL),
		captureCodeTTL:     durationToSeconds(c.CaptureCodeTTL),
		captureErrTimesTTL: durationToSeconds(c.CaptureErrTimesTTL),
		applyInfoTTL:       durationToSeconds(c.ApplyInfoTTL),
	}
}

// BlockImpl is
func (d *Dao) BlockImpl() *block.Dao {
	return d.block
}
