package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// var .
var (
	confPath string
	client   *conf.Client
	// config
	Conf = &Config{}
)

// Config def.
type Config struct {
	// base
	// log
	Log *log.Config
	// tracer
	Tracer *trace.Config
	//app
	Verify *verify.Config
	// http
	BM *bm.ServerConfig
	// db
	Mysql *sql.Config
	// redis
	Redis *Redis
	// RPC
	RPCServer *rpc.ServerConfig
	// property
	Property *Property
}

// Redis redis.
type Redis struct {
	*redis.Config
	Expire xtime.Duration
}

// Property .
type Property struct {
	LoadRankPeriod xtime.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init init conf.
func Init() (err error) {
	if confPath == "" {
		return configCenter()
	}
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func configCenter() (err error) {
	if client, err = conf.New(); err != nil {
		panic(err)
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
