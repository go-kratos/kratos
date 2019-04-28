package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/job/bbq/video/conf"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	searchv1 "go-common/app/service/bbq/search/api/grpc/v1"
	videov1 "go-common/app/service/bbq/video/api/grpc/v1"
	account "go-common/app/service/main/account/api"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"net/url"
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	gomail "gopkg.in/gomail.v2"
)

// Dao dao
type Dao struct {
	c             *conf.Config
	redis         *redis.Pool
	bfredis       *redis.Pool
	db            *xsql.DB
	dbCms         *xsql.DB
	HTTPClient    *bm.Client
	SearchClient  searchv1.SearchClient
	VideoClient   videov1.VideoClient
	AccountClient account.AccountClient
	email         *gomail.Dialer
	noticeClient  notice.NoticeClient
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:             c,
		redis:         redis.NewPool(c.Redis),
		bfredis:       redis.NewPool(c.BfRedis),
		db:            xsql.NewMySQL(c.MySQL),
		dbCms:         xsql.NewMySQL(c.MySQLCms),
		HTTPClient:    bm.NewClient(c.BM.Client),
		SearchClient:  newSearchClient(c.GRPCClient["search"]),
		VideoClient:   newVideoClient(c.GRPCClient["video"]),
		AccountClient: newAccountClient(c.GRPCClient["account"]),
		email:         gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.From, c.Mail.Password),
		noticeClient:  newNoticeClient(c.GRPCClient["notice"]),
	}
	return
}

// newNoticeClient .
func newNoticeClient(cfg *conf.GRPCConf) notice.NoticeClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return notice.NewNoticeClient(cc)
}

//newSearchClient .
func newSearchClient(cfg *conf.GRPCConf) searchv1.SearchClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return searchv1.NewSearchClient(cc)
}

//newAccountClient .
func newAccountClient(cfg *conf.GRPCConf) account.AccountClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return account.NewAccountClient(cc)
}

//newVideoClient
func newVideoClient(cfg *conf.GRPCConf) videov1.VideoClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return videov1.NewVideoClient(cc)
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
