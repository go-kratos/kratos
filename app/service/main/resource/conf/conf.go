package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	xlog "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
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
	Version string `toml:"version"`
	// reload
	Reload *ReloadInterval
	// rpc server2
	RPCServer *rpc.ServerConfig
	// verify
	Verify *verify.Config
	// http
	BM *BM
	// tracer
	Tracer *trace.Config
	// db
	DB *DB
	// httpClient
	HTTPClient *bm.ClientConfig
	// Host
	Host *Host
	// XLog
	XLog *xlog.Config
	// rpc
	LocationRPC *rpc.ClientConfig
	ArchiveRPC  *rpc.ClientConfig
	// redis
	Redis *Redis
	// hash number
	HashNum int64
	// databus
	ArchiveSub *databus.Config
	// qiye wechat
	WeChatToken  string
	WeChatSecret string
	WeChantUsers []string
	// kai guan off line
	MonitorArchive bool
	MonitorURL     bool
	// sp limit
	SpLimit time.Duration
}

// BM http
type BM struct {
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

// ReloadInterval define reolad config
type ReloadInterval struct {
	Ad time.Duration
}

// Host defeine host info
type Host struct {
	DataPlat string
	Ad       string
}

// DB define MySQL config
type DB struct {
	Res     *sql.Config
	Ads     *sql.Config
	Show    *sql.Config
	Manager *sql.Config
}

// Redis define Redis config
type Redis struct {
	Ads *struct {
		*redis.Config
		Expire time.Duration
	}
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
			xlog.Error("get config from config center error(%v)", e)
		}
	}()
	return
}
