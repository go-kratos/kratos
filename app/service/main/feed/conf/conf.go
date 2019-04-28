package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
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
	// log
	Xlog *log.Config
	// BM
	BM *bm.ServerConfig
	// rpc server
	RPCServer *rpc.ServerConfig
	// redis
	MultiRedis *MultiRedis
	// memcache
	Memcache *Memcache
	// tracer
	Tracer *trace.Config
	// rpc client
	ArchiveRPC *rpc.ClientConfig
	AccountRPC *rpc.ClientConfig
	ArticleRPC *rpc.ClientConfig
	// httpClient
	HTTPClient *bm.ClientConfig
	// feed
	Feed *Feed
}

// Feed .
type Feed struct {
	AppLength         int
	WebLength         int
	ArchiveFeedLength int
	ArticleFeedLength int
	ArchiveFeedExpire time.Duration
	BangumiFeedExpire time.Duration
	AppPullInterval   time.Duration
	WebPullInterval   time.Duration
	ArtPullInterval   time.Duration
	BulkSize          int
	MinUpCnt          int
	MaxTotalCnt       int
}

// MultiRedis .
type MultiRedis struct {
	MaxArcsNum  int
	TTLUpper    time.Duration
	ExpireUpper time.Duration
	ExpireFeed  time.Duration
	Local       *redis.Config
	Cache       *redis.Config
}

// Memcache .
type Memcache struct {
	*memcache.Config
	Expire        time.Duration
	BangumiExpire time.Duration
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
