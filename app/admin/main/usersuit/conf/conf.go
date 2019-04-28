package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
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
	// tracer
	Tracer *trace.Config
	// http
	BM *bm.ServerConfig
	// Auth auth
	Auth *permit.Config
	// MySQL mysql
	DB *DB
	// BFS
	BFS *BFS
	// redis
	Redis *redis.Config
	// memcache
	Memcache *Memcache
	// AccountGRPC account grpc
	AccountGRPC *warden.ClientConfig
	// http client
	HTTPClient *bm.ClientConfig
	// databus
	AccountNotify *databus.Config
	// host
	Host *Host
	// ecodes
	Ecode *ecode.Config
}

// Host host config .
type Host struct {
	Bfs     string
	Msg     string
	Manager string
}

// BFS bfs config
type BFS struct {
	Bucket string
	Key    string
	Secret string
}

// DB .
type DB struct {
	Usersuit *sql.Config
}

// Memcache config.
type Memcache struct {
	*memcache.Config
	Expire      time.Duration
	PointExpire time.Duration
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
