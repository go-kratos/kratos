package conf

import (
	"errors"
	"flag"

	"go-common/app/interface/main/app-feed/model/feed"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

type Config struct {
	// env
	Env string
	// infoc log2
	ShowInfoc2     *infoc.Config
	TagInfoc2      *infoc.Config
	RedirectInfoc2 *infoc.Config
	// show  XLog
	XLog *log.Config
	// tick time
	Tick xtime.Duration
	// tracer
	Tracer *trace.Config
	// httpClinet
	HTTPClient *bm.ClientConfig
	// httpAsyn
	HTTPClientAsyn *bm.ClientConfig
	// httpData
	HTTPData *bm.ClientConfig
	// httpTag
	HTTPTag *bm.ClientConfig
	// httpAd
	HTTPAd *bm.ClientConfig
	// httpActivity
	HTTPActivity *bm.ClientConfig
	// httpBangumi
	HTTPBangumi *bm.ClientConfig
	// httpShow
	HTTPShow *bm.ClientConfig
	// httpDynamic
	HTTPDynamic *bm.ClientConfig
	// httpClinet
	HTTPSearch *bm.ClientConfig
	// rpc Location
	LocationRPC *rpc.ClientConfig
	// http
	BM *HTTPServers
	// host
	Host *Host
	// db
	MySQL *MySQL
	// redis
	Redis *Redis
	// mc
	Memcache *Memcache
	// rpc client
	AccountRPC  *rpc.ClientConfig
	RelationRPC *rpc.ClientConfig
	ArchiveRPC  *rpc.ClientConfig
	FeedRPC     *rpc.ClientConfig
	TagRPC      *rpc.ClientConfig
	ArticleRPC  *rpc.ClientConfig
	ResourceRPC *rpc.ClientConfig
	// databus
	DislikeDatabus *databus.Config
	// ecode
	Ecode *ecode.Config
	// feed
	Feed *Feed
	// bnj2018
	Bnj *BnjConfig
	// BroadcastRPC grpc
	PGCRPC *warden.ClientConfig
	// autoplay mids
	AutoPlayMids []int64
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
}

// BnjConfig 2018拜年祭配置
type BnjConfig struct {
	TabImg    string
	TabID     int64
	BeginTime string
}

type Host struct {
	LiveAPI   string
	Bangumi   string
	Data      string
	Hetongzi  string
	APICo     string
	Ad        string
	Activity  string
	Rank      string
	Show      string
	Dynamic   string
	DynamicCo string
	BigData   string
	Search    string
}

type MySQL struct {
	Show    *sql.Config
	Manager *sql.Config
}

type Redis struct {
	Feed *struct {
		*redis.Config
		ExpireRecommend xtime.Duration
		ExpireBlack     xtime.Duration
	}
	Upper *struct {
		*redis.Config
		ExpireUpper xtime.Duration
	}
}

type Memcache struct {
	Feed *struct {
		*memcache.Config
		ExpireArchive xtime.Duration
	}
	Cache *struct {
		*memcache.Config
		ExpireCache xtime.Duration
	}
}

type Feed struct {
	// feed
	FeedCacheCount int
	LiveFeedCount  int
	// index
	Index *Index
	// ad
	CMResource map[string]int64
}

type Index struct {
	Count          int
	IPadCount      int
	MoePosition    int
	FollowPosition int
	// only archive for data disaster recovery
	Abnormal   bool
	Interest   []string
	FollowMode *feed.FollowMode
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init conf
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	if err = load(); err != nil {
		return
	}
	client.Watch("app-feed.toml")
	go func() {
		for range client.Event() {
			log.Info("config reload")
			if load() != nil {
				log.Error("config reload error (%v)", err)
			}
		}
	}()
	return
}

func load() (err error) {
	var (
		s       string
		ok      bool
		tmpConf *Config
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
