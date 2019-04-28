package dao

import (
	"context"

	account "go-common/app/service/main/account/api"
	archive "go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf/paladin"
	"go-common/library/database/elastic"
	"go-common/library/database/orm"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"

	"github.com/jinzhu/gorm"
)

// Dao is the appeal database access object
type Dao struct {
	ReadORM *gorm.DB
	ORM     *gorm.DB
	// search
	httpRead  *bm.Client
	httpWrite *bm.Client

	// memcache
	mc *memcache.Pool

	// redis
	redis *redis.Pool

	// es
	es *elastic.Elastic

	// account-service rpc
	accRPC account.AccountClient
	// archive-service rpc
	arcRPC archive.ArchiveClient

	c *paladin.Map // application.toml can reload

	writeConf *bm.ClientConfig
}

// New will create a new appeal Dao instance
func New() (d *Dao) {
	var (
		db struct {
			ReadORM *orm.Config
			ORM     *orm.Config
		}
		http struct {
			HTTPClientRead *bm.ClientConfig
			HTTPClient     *bm.ClientConfig
			Elastic        *elastic.Config
		}
		grpc struct {
			Account *warden.ClientConfig
			Archive *warden.ClientConfig
		}
		mc struct {
			Workflow *memcache.Config
		}
		rds struct {
			Workflow *redis.Config
		}
		ac = new(paladin.TOML)
	)
	checkErr(paladin.Watch("application.toml", ac))
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&db))
	checkErr(paladin.Get("http.toml").UnmarshalTOML(&http))
	checkErr(paladin.Get("memcache.toml").UnmarshalTOML(&mc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rds))
	checkErr(paladin.Get("grpc.toml").UnmarshalTOML(&grpc))
	d = &Dao{
		c:         ac,
		ReadORM:   orm.NewMySQL(db.ReadORM),
		ORM:       orm.NewMySQL(db.ORM),
		httpRead:  bm.NewClient(http.HTTPClientRead),
		httpWrite: bm.NewClient(http.HTTPClient),
		// memcache
		mc: memcache.NewPool(mc.Workflow),
		// redis
		redis: redis.NewPool(rds.Workflow),
		// es
		//es: elastic.NewElastic(nil),
		es:        elastic.NewElastic(http.Elastic),
		writeConf: http.HTTPClient,
	}
	// account-service rpc
	var err error
	if d.accRPC, err = account.NewClient(grpc.Account); err != nil {
		panic(err)
	}
	// archive-service rpc
	if d.arcRPC, err = archive.NewClient(grpc.Archive); err != nil {
		panic(err)
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	d.ORM.LogMode(true)
	d.ReadORM.LogMode(true)
}

// Close close dao.
func (d *Dao) Close() {
	if d.ORM != nil {
		d.ORM.Close()
	}
	if d.ReadORM != nil {
		d.ReadORM.Close()
	}
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.ORM.DB().PingContext(c); err != nil {
		return
	}
	if err = d.ReadORM.DB().PingContext(c); err != nil {
		return
	}

	if err = d.pingMC(c); err != nil {
		return
	}

	return d.pingRedis(c)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
