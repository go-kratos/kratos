package creative

import (
	"context"
	"fmt"
	"go-common/app/interface/main/creative/conf"
	"go-common/library/cache"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
)

// Dao is creative dao.
type Dao struct {
	// config
	c     *conf.Config
	cache *cache.Cache
	mc    *memcache.Pool
	// db
	creativeDb *sql.DB
	// select
	getTypeStmt *sql.Stmt
	getToolStmt *sql.Stmt
}

type BgmData struct {
	Data string `json:"data"`
}

// New init api url
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		creativeDb: sql.NewMySQL(c.DB.Creative),
		cache:      cache.New(10240, 10),
		mc:         memcache.NewPool(c.Memcache.Archive.Config),
	}
	d.getTypeStmt = d.creativeDb.Prepared(_getTypeSQL)
	d.getToolStmt = d.creativeDb.Prepared(_getToolByTypeSQL)
	return
}

// Ping creativeDb
func (d *Dao) Ping(c context.Context) (err error) {
	return d.creativeDb.Ping(c)
}

// Close creativeDb
func (d *Dao) Close() (err error) {
	d.mc.Close()
	return d.creativeDb.Close()
}

func bgmKey(aid, cid, mtype int64) string {
	return fmt.Sprintf("bgm_oid_%d_%d_%d", aid, cid, mtype)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -sync=true -nullcache=&BgmData{Data:""} -check_null_code=$.Data==""
	BgmData(c context.Context, aid, cid, mtype int64) (*BgmData, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// mc: -key=bgmKey
	CacheBgmData(c context.Context, key int64) (*BgmData, error)
	// 这里也支持自定义注释 会替换默认的注释
	// mc: -key=bgmKey -expire=3 -encode=json|gzip
	AddCacheBgmData(c context.Context, key int64, value *BgmData) error
	// mc: -key=bgmKey
	DelCacheBgmData(c context.Context, key int64) error
}
