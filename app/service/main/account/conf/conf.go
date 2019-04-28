package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/elastic"
	"go-common/library/log"
	"go-common/library/naming/livezk"
	bm "go-common/library/net/http/blademaster"
	v "go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf common conf
	Conf   = &Config{}
	client *conf.Client
)

//Config config struct
type Config struct {
	Host *Host
	// log
	Log *log.Config
	// http
	BM        *bm.ServerConfig
	RPCServer *rpc.ServerConfig
	// http client
	HTTPClient HTTPClient
	// identify
	// Identify *identify.Config
	// mc
	Memcache *Memcache
	// tracer
	Tracer *trace.Config
	// rpc
	MemberRPC   *rpc.ClientConfig
	MemberGRPC  *warden.ClientConfig
	RelationRPC *rpc.ClientConfig
	CoinRPC     *rpc.ClientConfig
	SuitRPC     *rpc.ClientConfig
	BlockRPC    *rpc.ClientConfig
	// warden
	WardenServer *warden.ServerConfig
	LiveZK       *livezk.Zookeeper
	// Elastic config
	Elastic      *elastic.Config
	Verify       *v.Config
	AppkeyFilter *AppkeyFilter
}

// AppkeyFilter is.
type AppkeyFilter struct {
	Privacy []string
}

// Host host.
type Host struct {
	AccountURI  string
	VipURI      string
	PassportURI string
}

// HTTPClient config
type HTTPClient struct {
	Read    *bm.ClientConfig
	Write   *bm.ClientConfig
	Privacy *bm.ClientConfig
}

// Memcache config
type Memcache struct {
	Account       *memcache.Config
	AccountExpire time.Duration
}

func configCenter() (err error) {
	if client, err = conf.New(); err != nil {
		panic(err)
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
		return errors.New("<account-service.toml> is not exists")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config, maybe <account-service.toml> file err")
	}
	*Conf = *tmpConf
	return
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config
func Init() (err error) {
	if confPath == "" {
		return configCenter()
	}
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
