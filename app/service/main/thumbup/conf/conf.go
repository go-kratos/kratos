package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/tidb"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/rate"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// global var
var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config config set
type Config struct {
	// elk
	Log *log.Config
	// BM
	BM *bm.ServerConfig
	// rpc server
	RPCServer *rpc.ServerConfig
	GRPC      *warden.ServerConfig
	// tracer
	Tracer *trace.Config
	// verify
	Verify *verify.Config
	Rate   *rate.Config
	// redis
	Redis *Redis
	// memcache
	Memcache *Memcache
	// Tidb
	Tidb *tidb.Config
	// ecode
	Ecode       *ecode.Config
	StatDatabus *databus.Config
	LikeDatabus *databus.Config
	ItemDatabus *databus.Config
	UserDatabus *databus.Config
	StatMerge   *StatMerge
	// ThumbUp
	ThumbUp ThumbUp
}

// StatMerge .
type StatMerge struct {
	Business string
	Target   int64
	Sources  []int64
}

// Memcache config
type Memcache struct {
	*memcache.Config
	StatsExpire time.Duration
}

// Redis config
type Redis struct {
	*redis.Config
	StatsExpire     time.Duration
	UserLikesExpire time.Duration
	ItemLikesExpire time.Duration
}

// ThumbUp thumb up config
type ThumbUp struct {
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init conf
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
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
