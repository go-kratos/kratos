package dao

import (
	"context"

	"go-common/app/job/bbq/recall/internal/conf"
	recall "go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/net/rpc/warden"
)

// Dao dao
type Dao struct {
	c            *conf.Config
	redis        *redis.Pool
	bfredis      *redis.Pool
	db           *xsql.DB
	dbOffline    *xsql.DB
	dbCms        *xsql.DB
	recallClient recall.RecsysRecallClient
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:            c,
		redis:        redis.NewPool(c.Redis),
		bfredis:      redis.NewPool(c.BfRedis),
		db:           xsql.NewMySQL(c.MySQL),
		dbOffline:    xsql.NewMySQL(c.OfflineMySQL),
		dbCms:        xsql.NewMySQL(c.CmsMySQL),
		recallClient: newRecallClient(c.GRPCClient["recall"]),
	}
	return
}

func newRecallClient(cfg *conf.GRPCConfig) recall.RecsysRecallClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return recall.NewRecsysRecallClient(cc)
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return d.db.Ping(ctx)
}
