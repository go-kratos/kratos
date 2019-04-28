package dao

import (
	"context"
	"go-common/app/service/bbq/user/api"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"

	notice "go-common/app/service/bbq/notice-service/api/v1"
	"go-common/app/service/bbq/user/internal/conf"
	acc "go-common/app/service/main/account/api"
	filter "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/net/rpc/warden"
)

// Dao dao
type Dao struct {
	c             *conf.Config
	cache         *fanout.Fanout
	redis         *redis.Pool
	db            *xsql.DB
	accountClient acc.AccountClient
	noticeClient  notice.NoticeClient
	filterClient  filter.FilterClient
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -batch=10 -max_group=10 -batch_err=break -nullcache=&api.UserBase{Mid:-1} -check_null_code=$==nil||$.Mid==-1
	UserBase(c context.Context, mid []int64) (map[int64]*api.UserBase, error)
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:             c,
		cache:         fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		redis:         redis.NewPool(c.Redis),
		db:            xsql.NewMySQL(c.MySQL),
		accountClient: newAccountClient(c.GRPCClient["account"]),
		noticeClient:  newNoticeClient(c.GRPCClient["notice"]),
		filterClient:  newFilterClient(c.GRPCClient["filter"]),
	}
	return
}

// newNoticeClient .
func newFilterClient(cfg *conf.GRPCConf) filter.FilterClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return filter.NewFilterClient(cc)
}

// newNoticeClient .
func newNoticeClient(cfg *conf.GRPCConf) notice.NoticeClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return notice.NewNoticeClient(cc)
}

//newAccountClient .
func newAccountClient(cfg *conf.GRPCConf) acc.AccountClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return acc.NewAccountClient(cc)
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

// BeginTran begin mysql transaction
func (d *Dao) BeginTran(c context.Context) (*xsql.Tx, error) {
	return d.db.Begin(c)
}

// CreateNotice 创建通知
func (d *Dao) CreateNotice(ctx context.Context, notice *notice.NoticeBase) (err error) {
	_, err = d.noticeClient.CreateNotice(ctx, notice)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "create notice fail: notice="+notice.String()))
		return
	}

	log.V(1).Infov(ctx, log.KV("log", "create notice: notice="+notice.String()))
	return
}

// Filter .
func (d *Dao) Filter(ctx context.Context, content string, area string) (level int32, err error) {
	req := new(filter.FilterReq)
	req.Message = content
	req.Area = area
	reply, err := d.filterClient.Filter(ctx, req)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "filter fail : req="+req.String()))
		return
	}
	level = reply.Level
	log.V(1).Infov(ctx, log.KV("log", "get filter reply="+reply.String()))
	return
}
