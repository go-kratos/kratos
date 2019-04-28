package conf

import (
	"errors"
	"flag"
	"go-common/library/queue/databus"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Config struct of conf.
type Config struct {
	// base
	// log
	Log *log.Config
	// db
	DB *DB
	// mc
	Memcache    *Memcache
	LevelExpire int32
	// http
	BM *HTTPServers
	// http client
	HTTPClient HTTPClient
	// app
	App *bm.App
	// tracer
	Tracer *trace.Config
	// switch
	Switch *ConfigSwitch
	// report
	Report *databus.Config
}

// ConfigSwitch switch config.
type ConfigSwitch struct {
	QueryExp uint64
}

// DB db config.
type DB struct {
	Exp *sql.Config
}

// Memcache config
type Memcache struct {
	Exp       *memcache.Config
	ExpExpire time.Duration
}

// HTTPServers Http Servers
type HTTPServers struct {
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

// HTTPClient config
type HTTPClient struct {
	Read  *bm.ClientConfig
	Write *bm.ClientConfig
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
			log.Info("config event")
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

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init int config
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}
