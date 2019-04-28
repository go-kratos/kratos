package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	xlog "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Config .
type Config struct {
	// Env
	Env string
	// App
	App *bm.App
	// Xlog is go-common log.
	Xlog *xlog.Config
	// rpc
	FeedRPC    *rpc.ClientConfig
	AccountRPC *rpc.ClientConfig
	// multihttp
	BM   *bm.ServerConfig
	Auth *auth.Config
	// tracer
	Tracer *trace.Config
	// feed
	Feed *feed
	// memcache
	Memcache *mc
}

// Feed feed controls
type feed struct {
	DefaultSize int
	MaxSize     int
}

type mc struct {
	*memcache.Config
	FeedExpire xtime.Duration
}

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init .
func Init() (err error) {
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
			xlog.Info("config reload")
			if load() != nil {
				xlog.Error("config reload err")
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
