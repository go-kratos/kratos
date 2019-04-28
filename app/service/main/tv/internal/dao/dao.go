package dao

import (
	"context"

	acclgrpc "go-common/app/service/main/account/api"
	"go-common/app/service/main/tv/internal/conf"
	"go-common/app/service/main/tv/internal/model"
	"go-common/app/service/main/tv/internal/pkg"
	mvipgrpc "go-common/app/service/main/vipinfo/api"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"
)

// Dao dao
type Dao struct {
	c           *conf.Config
	mc          *memcache.Pool
	db          *xsql.DB
	httpCli     *bm.Client
	mvipCli     mvipgrpc.VipInfoClient
	mvipHttpCli *bm.Client
	accCli      acclgrpc.AccountClient
	ystCli      *YstClient
	cache       *fanout.Fanout
	cacheTTL    *conf.CacheTTL
	signer      *pkg.Signer
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	httpCli := bm.NewClient(c.HTTPClient)
	mvipHttpCli := bm.NewClient(c.HTTPClient)
	mvipCli, err := mvipgrpc.NewClient(c.MVIPClient)
	if err != nil {
		panic(err)
	}
	accCli, err := acclgrpc.NewClient(c.ACCClient)
	if err != nil {
		panic(err)
	}
	dao = &Dao{
		c:           c,
		mc:          memcache.NewPool(c.Memcache),
		db:          xsql.NewMySQL(c.MySQL),
		mvipCli:     mvipCli,
		mvipHttpCli: mvipHttpCli,
		accCli:      accCli,
		httpCli:     httpCli,
		ystCli:      NewYstClient(httpCli),
		cache:       fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		cacheTTL:    c.CacheTTL,
		signer:      &pkg.Signer{Key: c.YST.Key},
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	return d.db.Ping(ctx)
}

// BeginTran begins transaction.
func (d *Dao) BeginTran(c context.Context) (*xsql.Tx, error) {
	return d.db.Begin(c)
}

// EndTran ends transaction.
func (d *Dao) EndTran(tx *xsql.Tx, err error) error {
	if err != nil {
		log.Info("d.EndTran.Rollback(%+v) err(%+v)", tx, err)
		tx.Rollback()
	} else {
		err = tx.Commit()
	}
	return err
}

// Signer returns yst signer.
func (d *Dao) Signer() *pkg.Signer {
	return d.signer
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -nullcache=&model.UserInfo{ID:-1} -check_null_code=$.ID==-1
	UserInfoByMid(c context.Context, key int64) (*model.UserInfo, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// mc: -key=userInfoKey
	CacheUserInfoByMid(c context.Context, key int64) (*model.UserInfo, error)
	// mc: -key=userInfoKey -expire=d.cacheTTL.UserInfoTTL -encode=json
	AddCacheUserInfoByMid(c context.Context, key int64, value *model.UserInfo) error
	// mc: -key=userInfoKey
	DelCacheUserInfoByMid(c context.Context, key int64) error
	// mc: -key=payParamKey
	CachePayParamByToken(c context.Context, token string) (*model.PayParam, error)
	// mc: -key=payParamKey
	CachePayParamsByTokens(c context.Context, tokens []string) (map[string]*model.PayParam, error)
	// mc: -key=payParamKey -expire=d.cacheTTL.PayParamTTL -encode=json
	AddCachePayParam(c context.Context, key string, value *model.PayParam) error
	// mc: -type=replace -key=payParamKey -expire=d.cacheTTL.PayParamTTL -encode=json
	UpdateCachePayParam(c context.Context, key string, value *model.PayParam) error
}
