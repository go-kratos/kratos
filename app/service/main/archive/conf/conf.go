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
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// config
var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

// Config struct
type Config struct {
	// base
	Tick      time.Duration
	Videoshot *Videoshot
	// xlog
	Xlog *log.Config
	// tracer
	Tracer *trace.Config
	// http
	BM *BM
	// http client
	PlayerClient *bm.ClientConfig
	// switch get player
	PlayerSwitch bool
	PlayerNum    int64
	// player qn config
	PlayerQn    []int
	PlayerVipQn []int
	BnjList     []int64
	// PlayerAPI path
	PlayerAPI    string
	PGCPlayerAPI string
	RPCServer    *rpc.ServerConfig
	// db
	DB *DB
	// ecode
	Ecode *ecode.Config
	// rpc client
	AccountRPC *rpc.ClientConfig
	// grpc client
	AccClient *warden.ClientConfig
	// mc
	Memcache *Memcache
	// redis
	Redis *Redis
	// databus
	Databus      *databus.Config
	StatDatabus  *databus.Config
	ShareDatabus *databus.Config
	CacheDatabus *databus.Config
}

// BM http
type BM struct {
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

// Videoshot videoshot uri and key
type Videoshot struct {
	URI string
	Key string
}

// DB db config
type DB struct {
	// archive db
	Arc       *sql.Config
	ArcRead   *sql.Config
	ArcResult *sql.Config
	Stat      *sql.Config
	Click     *sql.Config
}

// Memcache memcache config
type Memcache struct {
	Archive *struct {
		*memcache.Config
		ArchiveExpire time.Duration
		VideoExpire   time.Duration
		StatExpire    time.Duration
	}
}

// Redis redis config
type Redis struct {
	Archive *struct {
		*redis.Config
	}
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config.
func Init() (err error) {
	if confPath != "" {
		_, err = toml.DecodeFile(confPath, &Conf)
		return
	}
	err = remote()
	return
}

func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	if err = load(); err != nil {
		return
	}
	client.Watch("archive-service.toml")
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
		tmpConf = &Config{}
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
