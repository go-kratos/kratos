package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
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

// Config .
type Config struct {
	// elk
	Log *log.Config
	// Verify
	Verify *verify.Config
	// app
	App *bm.App
	// http client
	HTTPClient httpClient
	// rpc
	ArchiveRPC *rpc.ClientConfig
	// rpc server
	RPCServer *rpc.ServerConfig
	// tracer
	Tracer *trace.Config
	// mc
	Memcache *Memcache
	// Rule
	Rule *Rule
	// Host
	Host *Host
	// HTTPServer
	HTTPServer *bm.ServerConfig
}

// Host hosts.
type Host struct {
	BigDataURI string
	LiveURI    string
	APIURI     string
}

// Rule config.
type Rule struct {
	// region tick.
	TickRegion time.Duration
	// tag tick.
	TickTag time.Duration
	// default num of dynamic archives.
	NumArcs int
	//default num of index dynamic archives.
	NumIndexArcs int
	//min region count.
	MinRegionCount int
}
type httpClient struct {
	Read *bm.ClientConfig
}

// Memcache memcache
type Memcache struct {
	*memcache.Config
	Expire time.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init init conf
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
