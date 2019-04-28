package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// global var
var (
	ConfPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config config set
type Config struct {
	// base
	// elk
	Log *log.Config

	BM *HTTPServers
	// tracer
	Tracer *trace.Config
	// memcache
	Memcache     *Memcache
	WalletExpire int32
	// redis
	Redis *Redis
	// db
	DB *DB
	// http client
	HTTPClient *bm.ClientConfig
	//DataBus
	DataBus *DataBus
}

type DataBus struct {
	Change *databus.Config
}

type HTTPServers struct {
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

// DB db config.
type DB struct {
	Wallet *sql.Config
}

type Memcache struct {
	Wallet       *memcache.Config
	WalletExpire time.Duration
}

type Redis struct {
	Wallet *redis.Config
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init init conf
func Init() error {
	if ConfPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(ConfPath, &Conf)
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
	if s, ok = client.Toml(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
