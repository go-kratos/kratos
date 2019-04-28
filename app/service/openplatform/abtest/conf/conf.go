package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

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
	// redis
	Redis *Redis
	// http
	BM *HTTPServers
	// http client
	HTTPClient HTTPClient
	// tracer
	Tracer *trace.Config
	//stat
	Stat int
	Stra int
}

// DB db config.
type DB struct {
	Ab *sql.Config
}

// Redis conf.
type Redis struct {
	*redis.Config
	Expire        xtime.Duration
	VerifyCdTimes xtime.Duration
}

// HTTPServers Http Servers
type HTTPServers struct {
	Inner *bm.ServerConfig
	Outer *bm.ServerConfig
}

// HTTPClient config
type HTTPClient struct {
	Read  *bm.ClientConfig
	Write *bm.ClientConfig
}

func outer() (err error) {
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
	if s, ok = client.Value("open-abtest.toml"); !ok {
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
		return outer()
	}
	return remote()
}
