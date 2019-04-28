package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/rate"
	v "go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
	"go-common/library/net/rpc/warden"
)

// Conf global variable.
var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

// Config struct of conf.
type Config struct {
	Host *Host
	// bm
	BM *bm.ServerConfig
	// log
	Log *log.Config
	// rpc server2
	RPCServer *rpc.ServerConfig
	// db
	Mysql *sql.Config
	// mc
	Memcache *Memcache
	// redis
	Redis *Redis
	// tracer
	Tracer *trace.Config
	// realtion
	Relation *Relation
	// rpc clients
	RPCClient2 *RPC
	// Infoc
	Infoc *infoc.Config
	// Antispam
	Antispam *antispam.Config
	// statCache
	StatCache *StatCache
	// httpClinet
	HTTPClient *bm.ClientConfig
	// Report
	Report *databus.Config
	// Verify
	Verify *v.Config
	// Rate
	AddFollowingRate *rate.Config
	// WardenServer
	WardenServer *warden.ServerConfig
}

// Host host.
type Host struct {
	Passport string
}

// RPC clients config.
type RPC struct {
	// member rpc client
	Member *rpc.ClientConfig
}

// Redis redis
type Redis struct {
	*redis.Config
	Expire time.Duration
}

// Memcache memcache
type Memcache struct {
	*memcache.Config
	Expire         time.Duration
	FollowerExpire time.Duration
}

// Relation relation related config
type Relation struct {
	MaxFollowingCached int
	MaxFollowerCached  int
	MaxWhisperCached   int
	// max following limit
	MaxFollowingLimit int
	// max black limit
	MaxBlackLimit int
	// monitor switch: true, user cannot be following.
	Monitor bool
	// prompt
	Period time.Duration // prompt count flush period
	Bcount int64         // business prompt count
	Ucount int64         // up prompt count
	// followers unread duration
	FollowersUnread time.Duration
	// achieve key
	AchieveKey string
}

// StatCache is
type StatCache struct {
	Size          int
	Expire        time.Duration
	LeastFollower int
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init conf.
func Init() (err error) {
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
