package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/naming/livezk"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf is global config
	Conf *Config
)

// Config service config
type Config struct {
	// base
	// db
	DB *DB
	// redis
	Redis *Redis
	// xlog
	Log *log.Config
	// verify
	Verify *verify.Config
	// http
	BM      *BM
	Account *warden.ClientConfig
	// rpc server2
	RPCServer *rpc.ServerConfig
	// tracer
	Tracer *trace.Config
	// file
	FilePath       string
	AnonymFileName string
	// filter ip
	FilterZone []string
	// grpc server
	WardenServer *warden.ServerConfig
	LiveZK       *livezk.Zookeeper
	// Host
	Host *Host
	// AnonymKey
	AnonymKey string
	// new library
	IPv4Name string
	IPv6Name string
	// httpClinet
	HTTPClient *bm.ClientConfig
}

// BM http
type BM struct {
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

// DB define MySQL config
type DB struct {
	Zlimit *sql.Config
}

// Redis define Redis config
type Redis struct {
	Zlimit *Zlimit
}

// Zlimit struct about zlimit
type Zlimit struct {
	*redis.Config
	Expire time.Duration
}

// Host url
type Host struct {
	Maxmind string
	Bvcip   string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config
func Init() (err error) {
	if confPath != "" {
		_, err = toml.DecodeFile(confPath, &Conf)
		return
	}
	err = configCenter()
	return
}

// configCenter ugc
func configCenter() (err error) {
	var (
		client *conf.Client
		c      string
		ok     bool
	)
	if client, err = conf.New(); err != nil {
		panic(err)
	}
	if c, ok = client.Toml2(); !ok {
		err = errors.New("load config center error")
		return
	}
	_, err = toml.Decode(c, &Conf)
	go func() {
		for e := range client.Event() {
			log.Error("get config from config center error(%v)", e)
		}
	}()
	return
}
