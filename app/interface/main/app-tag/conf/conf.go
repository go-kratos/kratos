package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
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
	// redis
	Redis *Redis
	// show  XLog
	XLog *log.Config
	// tracer
	Tracer *trace.Config
	// bm http
	BM *HTTPServers
	// tick time
	Tick xtime.Duration
	// rpc client
	TagRPC      *rpc.ClientConfig
	ArchiveRPC  *rpc.ClientConfig
	ArchiveRPC2 *rpc.ClientConfig
	// host
	Host *Host
	// httpClinet
	HTTPClient *bm.ClientConfig
	// mc
	Memcache *Memcache
	// db
	MySQL *MySQL
	// infoc2
	FeedInfoc2 *infoc.Config
}

type Host struct {
	Data    string
	ApiCo   string
	Bangumi string
}

type Redis struct {
	Stat *struct {
		*redis.Config
		Expire xtime.Duration
	}
}

type HTTPServers struct {
	Outer *bm.ServerConfig
}

type Memcache struct {
	Archive *struct {
		*memcache.Config
		Expire xtime.Duration
	}
}

type MySQL struct {
	Show *sql.Config
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
