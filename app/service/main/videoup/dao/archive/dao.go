package archive

import (
	"context"
	"fmt"

	"go-common/app/service/main/videoup/conf"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/sync/pipeline/fanout"
)

// Dao is redis dao.
type Dao struct {
	c *conf.Config
	// db
	db      *sql.DB
	rddb    *sql.DB
	slaveDB *sql.DB
	redis   *redis.Pool
	cache   *fanout.Fanout
	mc      *memcache.Pool
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:       c,
		db:      sql.NewMySQL(c.DB.Archive),
		rddb:    sql.NewMySQL(c.DB.ArchiveRead),
		slaveDB: sql.NewMySQL(c.DB.ArchiveSlave),
		redis:   redis.NewPool(c.Redis.Track.Config),
		cache:   fanout.New("cache"),
		mc:      memcache.NewPool(c.Memcache.Archive.Config),
	}
	return d
}

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}

// Close close dao.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

func staffKey(aid int64) string {
	return fmt.Sprintf("staff_aid_%d", aid)
}

func (d *Dao) cacheSFStaffData(aid int64) string {
	return fmt.Sprintf("staff_aid_sf_%d", aid)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -singleflight=true -nullcache=[]*archive.Staff{{ID:-1}} -check_null_code=len($)==1&&$[0].ID==-1
	StaffData(c context.Context, aid int64) ([]*archive.Staff, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// mc: -key=staffKey
	CacheStaffData(c context.Context, key int64) ([]*archive.Staff, error)
	// 这里也支持自定义注释 会替换默认的注释
	// mc: -key=staffKey -expire=3 -encode=json|gzip
	AddCacheStaffData(c context.Context, key int64, value []*archive.Staff) error
	// mc: -key=staffKey
	DelCacheStaffData(c context.Context, key int64) error
}
