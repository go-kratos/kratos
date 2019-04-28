package archive

import (
	"context"
	"runtime"
	"time"

	hisrpc "go-common/app/interface/main/history/rpc/client"
	"go-common/app/interface/main/tv/conf"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Dao is archive dao.
type Dao struct {
	// cfg
	db        *sql.DB
	conf      *conf.Config
	relateURL string
	// memory
	arcTypes    map[int32]*arcwar.Tp   // map for arc types
	arcTypesRel map[int32][]*arcwar.Tp // map for relation between arc type
	// http client
	client *bm.Client
	// rpc
	arcClient arcwar.ArchiveClient
	hisRPC    *hisrpc.Service
	// memcache
	arcMc      *memcache.Pool
	expireRlt  int32
	expireArc  int32
	expireView int32
	// chan
	mCh chan func()
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client: bm.NewClient(c.HTTPClient),
		// rpc
		hisRPC: hisrpc.New(c.HisRPC),
		// cfg
		relateURL: c.Host.Data + _realteURL,
		conf:      c,
		db:        sql.NewMySQL(c.Mysql),
		// memorry
		arcTypes:    make(map[int32]*arcwar.Tp),
		arcTypesRel: make(map[int32][]*arcwar.Tp),
		// memcache
		arcMc:      memcache.NewPool(c.Memcache.Config),
		expireRlt:  int32(time.Duration(c.Memcache.RelateExpire) / time.Second),
		expireView: int32(time.Duration(c.Memcache.ViewExpire) / time.Second),
		expireArc:  int32(time.Duration(c.Memcache.ArcExpire) / time.Second),
		// mc proc
		mCh: make(chan func(), 10240),
	}
	var err error
	if d.arcClient, err = arcwar.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	// video db
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go d.cacheproc()
	}
	d.loadTypes(context.Background())
	return
}

// addCache add archive to mc or redis
func (d *Dao) addCache(f func()) {
	select {
	case d.mCh <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc write memcache and stat redis use goroutine
func (d *Dao) cacheproc() {
	for {
		f := <-d.mCh
		f()
	}
}
