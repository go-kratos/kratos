package conf

import (
	"errors"
	"flag"

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
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// confPath is.
	confPath string
	// Conf is.
	Conf = &Config{}
	// client is.
	client *conf.Client
)

// Config struct
type Config struct {
	ShowInfoc     *infoc.Config
	TagInfoc      *infoc.Config
	RedirectInfoc *infoc.Config
	CoinInfoc     *infoc.Config
	ViewInfoc     *infoc.Config
	RelateInfoc   *infoc.Config
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
	// httpBangumi
	HTTPBangumi *bm.ClientConfig
	// HTTPSearch
	HTTPSearch *bm.ClientConfig
	// HTTPAudio
	HTTPAudio *bm.ClientConfig
	// HTTPWrite
	HTTPWrite *bm.ClientConfig
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
	ArchiveRPC2 *rpc.ClientConfig
	TagRPC      *rpc.ClientConfig
	FavoriteRPC *rpc.ClientConfig
	CoinRPC     *rpc.ClientConfig
	AssistRPC   *rpc.ClientConfig
	ThumbupRPC  *rpc.ClientConfig
	DMRPC       *rpc.ClientConfig
	ResourceRPC *rpc.ClientConfig
	HisRPC      *rpc.ClientConfig
	ArticleRPC  *rpc.ClientConfig
	LocationRPC *rpc.ClientConfig
	// BroadcastRPC grpc
	PGCRPC       *warden.ClientConfig
	UGCpayClient *warden.ClientConfig
	// ecode
	Ecode *ecode.Config
	// feed
	Feed *Feed
	// view
	View *View
	// search
	Search *Search
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
}

// Host struct
type Host struct {
	Bangumi   string
	Data      string
	Hetongzi  string
	APICo     string
	Rank      string
	BigData   string
	Search    string
	AI        string
	Bvcvod    string
	VIP       string
	Playurl   string
	PlayurlBk string
}

// MySQL struct
type MySQL struct {
	Show    *sql.Config
	Manager *sql.Config
}

// Redis struct
type Redis struct {
	Feed *struct {
		*redis.Config
		ExpireRecommend xtime.Duration
		ExpireBlack     xtime.Duration
	}
}

// Memcache struct
type Memcache struct {
	Feed *struct {
		*memcache.Config
		Expire xtime.Duration
	}
	Cache *struct {
		*memcache.Config
		Expire xtime.Duration
	}
	Archive *struct {
		*memcache.Config
		RelateExpire xtime.Duration
	}
}

// Feed struct
type Feed struct {
	// index
	Index *Index
	// ad
	CMResource map[string]int64
}

// Index struct
type Index struct {
	Count          int
	IPadCount      int
	MoePosition    int
	FollowPosition int
	// only archive for data disaster recovery
	Abnormal bool
}

// View struct
type View struct {
	VipTick xtime.Duration
	// 相关推荐秒开个数
	RelateCnt int
}

// Search struct
type Search struct {
	SeasonNum          int
	MovieNum           int
	SeasonMore         int
	MovieMore          int
	UpUserNum          int
	UVLimit            int
	UserNum            int
	UserVideoLimit     int
	BiliUserNum        int
	BiliUserVideoLimit int
	OperationNum       int
	IPadSearchBangumi  int
	IPadSearchFt       int
}

// init is.
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

// local is.
func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

// reomte is.
func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	if err = load(); err != nil {
		return
	}
	client.Watch("app-intl.toml")
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

// load is.
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
