package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/app/service/bbq/video/conf"
	"go-common/app/service/bbq/video/model/grpc"
	account "go-common/app/service/main/account/api"
	archive "go-common/app/service/main/archive/api"
	"go-common/library/cache/redis"
	"go-common/library/conf/env"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
	"net/url"
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -batch_err=break -nullcache=&v1.VideoBase{Svid:-1} -check_null_code=$==nil||$.Svid==-1
	VideoBase(c context.Context, svid []int64) (map[int64]*v1.VideoBase, error)
}

// Dao dao
type Dao struct {
	c              *conf.Config
	redis          *redis.Pool
	cache          *fanout.Fanout
	db             *xsql.DB
	cmsdb          *xsql.DB
	httpClient     *bm.Client
	AccountClient  account.AccountClient
	cmsPub         *databus.Databus
	archiveSub     *databus.Databus
	archiveFilters *ArchiveFilters
	ArchiveClient  archive.ArchiveClient
	bvcPlayClient  grpc.PlayurlServiceClient
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:             c,
		redis:         redis.NewPool(c.Redis),
		cache:         fanout.New("cache"),
		db:            xsql.NewMySQL(c.MySQL),
		cmsdb:         xsql.NewMySQL(c.CMSMySQL),
		httpClient:    bm.NewClient(c.BM.Client),
		cmsPub:        databus.New(conf.Conf.Databus["cms"]),
		archiveSub:    newArchiveSub(c),
		AccountClient: newAccountClient(c.GRPCClient["account"]),
		ArchiveClient: newArchiveClient(c.GRPCClient["archive"]),
		bvcPlayClient: newBVCPlayClient(c.GRPCClient["bvcplay"]),
	}
	dao.newArchiveFilters(conf.ArchiveRules)
	return
}

// newBVCPlayClient .
func newBVCPlayClient(cfg *conf.GRPCConf) grpc.PlayurlServiceClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return grpc.NewPlayurlServiceClient(cc)
}

// newAccountClient .
func newAccountClient(cfg *conf.GRPCConf) account.AccountClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return account.NewAccountClient(cc)
}

// newArchiveClient .
func newArchiveClient(cfg *conf.GRPCConf) archive.ArchiveClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return archive.NewArchiveClient(cc)
}

// newArchiveFilters .
func (d *Dao) newArchiveFilters(c *conf.Rules) {
	response, err := d.ArchiveClient.Types(context.Background(), &archive.NoArgRequest{})
	if err != nil {
		panic(err)
	}

	detailFilter := &ArchiveDetailFilter{
		rule:        c.Archive,
		archiveType: response.Types,
	}
	upFilter := &ArchiveUpFilter{
		rule:          c.Up,
		accountClient: d.AccountClient,
	}
	dimensionFilter := &ArchivePageFilter{
		rule:          c.Dimension,
		archiveClient: d.ArchiveClient,
	}
	d.archiveFilters = NewArchiveFilters(detailFilter, upFilter, dimensionFilter)
}

func newArchiveSub(c *conf.Config) *databus.Databus {
	if env.DeployEnv != env.DeployEnvProd {
		return nil
	}
	if _, ok := c.Databus["archive"]; !ok {
		return nil
	}
	return databus.New(c.Databus["archive"])
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}

// BeginTran begin mysql transaction
func (d *Dao) BeginTran(c context.Context) (*xsql.Tx, error) {
	return d.db.Begin(c)
}

// ReplyHTTPCommon 评论公用请求
func replyHTTPCommon(c context.Context, httpClient *bm.Client, path string, method string, data map[string]interface{}, ip string) (r []byte, err error) {
	params := url.Values{}
	t := reflect.TypeOf(data).Kind()
	if t == reflect.Map {
		for k, v := range data {
			// params.Set(k, v.(string))
			switch reflect.TypeOf(v).Kind() {
			case reflect.Int64:
				params.Set(k, strconv.FormatInt(v.(int64), 10))
			case reflect.Int16:
				params.Set(k, strconv.FormatInt(int64(v.(int16)), 10))
			case reflect.String:
				params.Set(k, v.(string))
			case reflect.Int:
				params.Set(k, strconv.FormatInt(int64(v.(int)), 10))
			}
		}
	}
	log.V(5).Infov(c, log.KV("log", fmt.Sprintf("reply req url(%s)", path+"?"+params.Encode())))
	req, err := httpClient.NewRequest(method, path, ip, params)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("reply url(%s) error(%v)", path+"?"+params.Encode(), err)))
		return
	}
	var res struct {
		Code int             `json:"code"`
		Msg  string          `json:"message"`
		Data json.RawMessage `json:"data"`
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = httpClient.Do(c, req, &res); err != nil {
		str, _ := json.Marshal(res)
		log.Errorv(c, log.KV("log", fmt.Sprintf("reply ret data(%s) err[%v]", str, err)))
		return
	}
	str, _ := json.Marshal(res)
	log.V(5).Infov(c, log.KV("log", fmt.Sprintf("reply ret data(%s)", str)))
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Warnv(c, log.KV("log", fmt.Sprintf("reply url(%s) error(%v)", path+"?"+params.Encode(), err)))
	}
	r = res.Data
	return
}
