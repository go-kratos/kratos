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
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// ConfPath local config path
	ConfPath string
	// Conf config
	Conf   = &Config{}
	client *conf.Client
)

// Config str
type Config struct {
	// base
	// channal len
	ChanSize int64
	// log
	Log *log.Config
	// http
	BM *bm.ServerConfig
	// identify
	App *bm.App
	// tracer
	Tracer *trace.Config
	// tick load pgc
	Tick time.Duration
	// orm
	DB *DB
	// http client of search
	HTTPClient *HTTPClient
	// host
	Host *Host
	// redis
	Redis *Redis
}

// DB def db struct
type DB struct {
	Main  *sql.Config
	Slave *sql.Config
}

// Redis .
type Redis struct {
	*redis.Config
	UpRatingExpire time.Duration
}

// HTTPClient http client
type HTTPClient struct {
	Read *bm.ClientConfig
}

// Host http host
type Host struct {
	AccountURI string
	ArchiveURI string
	UperURI    string
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init init conf
func Init() (err error) {
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
