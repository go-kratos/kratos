package dao

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go-common/library/sync/pipeline/fanout"
	"net/url"
	"reflect"
	"strconv"

	"go-common/app/interface/bbq/app-bbq/conf"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/app/interface/bbq/app-bbq/model/grpc"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	user "go-common/app/service/bbq/user/api"
	image "go-common/app/service/bbq/video-image/api/grpc/v1"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	acc "go-common/app/service/main/account/api"
	filter "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"

	jsoniter "github.com/json-iterator/go"
)

var (
	videoPath string
)

// Dao dao
type Dao struct {
	c                *conf.Config
	cache            *fanout.Fanout
	redis            *redis.Pool
	db               *xsql.DB
	dbDM             *xsql.DB
	httpClient       *bm.Client
	httpslowClient   *bm.Client
	accountClient    acc.AccountClient
	recsysClient     recsys.RecsysClient
	imageClient      image.VideoImageClient
	bvcPlayClient    grpc.PlayurlServiceClient
	redundanceVideos []*model.RVideo
	noticeClient     notice.NoticeClient
	userClient       user.UserClient
	videoClient      video.VideoClient
	filterClient     filter.FilterClient
}

func init() {
	flag.StringVar(&videoPath, "video_json", "./video.json", "接口冗余降级video数据")
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:                c,
		cache:            fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		redis:            redis.NewPool(c.Redis),
		db:               xsql.NewMySQL(c.MySQL),
		dbDM:             xsql.NewMySQL(c.DMMySQL),
		httpClient:       bm.NewClient(c.HTTPClient.Normal),
		httpslowClient:   bm.NewClient(c.HTTPClient.Slow),
		accountClient:    newAccountClient(c.GRPCClient["account"]),
		recsysClient:     newRecsysClient(c.GRPCClient["recsys"]),
		imageClient:      newVideoImageClient(c.GRPCClient["videoimage"]),
		bvcPlayClient:    newBVCPlayClient(c.GRPCClient["bvcplay"]),
		redundanceVideos: model.RedundanceVideo(),
		noticeClient:     newNoticeClient(c.GRPCClient["notice"]),
		userClient:       newUserClient(c.GRPCClient["user"]),
		videoClient:      newVideoClient(c.GRPCClient["video"]),
		filterClient:     newFilterClient(c.GRPCClient["filter"]),
	}
	return
}

// newVideoClient .
func newVideoClient(cfg *conf.GRPCConf) video.VideoClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return video.NewVideoClient(cc)
}

// newUserClient .
func newUserClient(cfg *conf.GRPCConf) user.UserClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return user.NewUserClient(cc)
}

// newNoticeClient .
func newNoticeClient(cfg *conf.GRPCConf) notice.NoticeClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return notice.NewNoticeClient(cc)
}

// newBVCPlayClient .
func newBVCPlayClient(cfg *conf.GRPCConf) grpc.PlayurlServiceClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return grpc.NewPlayurlServiceClient(cc)
}

// newVideoImageClient .
func newVideoImageClient(cfg *conf.GRPCConf) image.VideoImageClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return image.NewVideoImageClient(cc)
}

//newAccountClient .
func newAccountClient(cfg *conf.GRPCConf) acc.AccountClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return acc.NewAccountClient(cc)
}

//newRecsysClient .
func newRecsysClient(cfg *conf.GRPCConf) recsys.RecsysClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return recsys.NewRecsysClient(cc)
}

// newUserClient .
func newFilterClient(cfg *conf.GRPCConf) filter.FilterClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return filter.NewFilterClient(cc)
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
