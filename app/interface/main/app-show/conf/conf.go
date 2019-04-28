package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	xlog "go-common/library/log"
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
	// Env
	Env string
	//show log
	ShowLog string
	// show  XLog
	XLog *xlog.Config
	// tick time
	Tick xtime.Duration
	// tracer
	Tracer *trace.Config
	// httpClinet
	HTTPClient *bm.ClientConfig
	// httpClinetAsyn
	HTTPClientAsyn *bm.ClientConfig
	// httpData
	HTTPData *bm.ClientConfig
	// bm http
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
	ArchiveRPC *rpc.ClientConfig
	// dynamicRPC client
	DynamicRPC *rpc.ClientConfig
	// rpc account
	AccountRPC *rpc.ClientConfig
	// resource
	ResourceRPC *rpc.ClientConfig
	// relationRPC
	RelationRPC *rpc.ClientConfig
	// location rpc
	LocationRPC *rpc.ClientConfig
	// rec host
	Recommend *Recommend
	// Infoc2
	Infoc2       *infoc.Config
	FeedInfoc2   *infoc.Config
	FeedTabInfoc *infoc.Config
	// databus
	DislikeDataBus *databus.Config
	// duration
	Duration *Duration
	// BroadcastRPC grpc
	PGCRPC *warden.ClientConfig
}

type Duration struct {
	// splash
	Splash string
	// search time_from
	Search string
}

type Host struct {
	ApiLiveCo    string
	Bangumi      string
	Hetongzi     string
	HetongziRank string
	Data         string
	ApiCo        string
	ApiCoX       string
	Ad           string
	Search       string
	Activity     string
	Dynamic      string
}

type HTTPServers struct {
	Outer *bm.ServerConfig
}

type MySQL struct {
	Show     *sql.Config
	Resource *sql.Config
}

type Redis struct {
	Recommend *struct {
		*redis.Config
		Expire xtime.Duration
	}
	Stat *struct {
		*redis.Config
		Expire xtime.Duration
	}
}

type Memcache struct {
	Archive *struct {
		*memcache.Config
		Expire xtime.Duration
	}
	Cards *struct {
		*memcache.Config
		Expire xtime.Duration
	}
}

type Recommend struct {
	Host  map[string][]string
	Group map[string]int
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config.
func Init() (err error) {
	if confPath != "" {
		_, err = toml.DecodeFile(confPath, &Conf)
		return
	}
	err = remote()
	return
}

func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	if err = load(); err != nil {
		return
	}
	go func() {
		for range client.Event() {
			xlog.Info("config reload")
			if load() != nil {
				xlog.Error("config reload error (%v)", err)
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
