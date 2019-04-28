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

var (
	// ConfPath dm-job config file path
	ConfPath string
	client   *conf.Client
	// Conf config export var
	Conf = &Config{}
)

// Config dm-job config struct
type Config struct {
	// base
	// log
	Xlog *log.Config
	// tracer
	Tracer *trace.Config
	// http
	HTTPServer *bm.ServerConfig
	// databus
	Databus *Databus
	// database
	DB *DB
	// redis
	Redis *Redis
	// memcache
	Memcache *Memcache
}

// Databus databus.
type Databus struct {
	DMMetaCsmr *databus.Config
}

// DB bilibili_dm
type DB struct {
	DMReader *sql.Config
	DMWriter *sql.Config
}

// Redis dm redis
type Redis struct {
	*redis.Config
	Expire time.Duration
}

// Memcache dm memcache
type Memcache struct {
	*memcache.Config
	Expire time.Duration
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "config path")
}

//Init int config
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
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
