package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
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
	// HTTPClient
	HTTPClient *HTTPClient
	// limit
	Limit *Limit
	// geetest
	Geetest *Geetest
	// rule
	Rule map[string]*Limit
	// url
	URL *URL
	// base
	Base *Base
}

// URL .
type URL struct {
	Shield string
}

// Base .
type Base struct {
	ShieldListTime int64
}

//Limit 限制
type Limit struct {
	Name             string
	SaleTimeOut      int64
	MIDCreateTimeOut int64
	MIDCreateMax     int64
	IPCreateTimeOut  int64
	IPCreateMax      int64
	IPChangeInterval int64
	IPWhiteList      []string
}

//Geetest 极验
type Geetest struct {
	Count int64
}

// DB db config.
type DB struct {
	AntiFraud *sql.Config
	PayShield *sql.Config
}

// Redis conf.
type Redis struct {
	*redis.Config
	Expire        xtime.Duration
	VerifyCdTimes xtime.Duration
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
	if s, ok = client.Value("anti-fraud.toml"); !ok {
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
