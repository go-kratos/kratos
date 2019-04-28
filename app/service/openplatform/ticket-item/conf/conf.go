package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	"go-common/library/log"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"

	"go-common/library/net/rpc/warden"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// Conf common conf
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Config config struct
type Config struct {
	// base
	// 数据库配置
	DB *DB
	// redis
	Redis *Redis
	// http client
	HTTPClient HTTPClient
	// http
	BM *HTTPServers
	// grpc server
	RPCServer *warden.ServerConfig
	// log
	Log *log.Config

	// auth
	Auth   *auth.Config
	Verify *verify.Config

	// orm
	ORM *orm.Config
	// UT
	UT *UT
	// URL
	URL *URL
	// BASECenter
	BASECenter *BASECenter
	// Tag
	Tag *Tag
}

// HTTPClient config
type HTTPClient struct {
	Read  *bm.ClientConfig
	Write *bm.ClientConfig
}

// URL external resources
type URL struct {
	ElasticHost string
	DefaultHead string
}

// BASECenter config
type BASECenter struct {
	AppID    string
	AppToken string
	URL      string
}

// HTTPServers Http Servers
type HTTPServers struct {
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

// Redis config
type Redis struct {
	Master *redis.Config
	Expire time.Duration
}

// DB config
type DB struct {
	Master *sql.Config
}

// UT config
type UT struct {
	DistPrefix string
}

// Tag config
type Tag struct {
	Tags string
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
	if s, ok = client.Value("item.toml"); !ok {
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

// Init init config
func Init() error {
	if confPath == "" {
		return remote()
	}
	return local()
}
