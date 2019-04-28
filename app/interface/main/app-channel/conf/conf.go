package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
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
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

type Config struct {
	// Env
	Env string
	// show  XLog
	Log *log.Config
	// tick time
	Tick xtime.Duration
	// tracer
	Tracer *trace.Config
	// httpClinet
	HTTPClient *bm.ClientConfig
	// httpClinetAsyn
	HTTPClientAsyn *bm.ClientConfig
	// HTTPShopping
	HTTPShopping *bm.ClientConfig
	// bm http
	BM *HTTPServers
	// host
	Host *Host
	// db
	MySQL *MySQL
	// rpc client
	ArchiveRPC *rpc.ClientConfig
	// rpc account
	AccountRPC *rpc.ClientConfig
	// relationRPC
	RelationRPC *rpc.ClientConfig
	// rpc client
	TagRPC *rpc.ClientConfig
	// rpc Article
	ArticleRPC *rpc.ClientConfig
	// rpc Location
	LocationRPC *rpc.ClientConfig
	// Infoc2
	FeedInfoc2    *infoc.Config
	ChannelInfoc2 *infoc.Config
	// memcache
	Memcache *Memcache
	// BroadcastRPC grpc
	PGCRPC *warden.ClientConfig
	// Square Count
	SquareCount int
}

type Host struct {
	Bangumi  string
	Data     string
	APICo    string
	Activity string
	LiveAPI  string
	Shopping string
}

type HTTPServers struct {
	Outer *bm.ServerConfig
}

type MySQL struct {
	Show    *sql.Config
	Manager *sql.Config
}

type Memcache struct {
	Channels *struct {
		*memcache.Config
		Expire xtime.Duration
	}
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
	client.Watch("app-channel.toml")
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
