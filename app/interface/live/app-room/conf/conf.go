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
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/trace"

	"go-common/library/net/rpc/warden"

	"go-common/library/net/rpc/liverpc"

	"github.com/BurntSushi/toml"
)

var (
	// ConfPath 配置地址
	ConfPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

type host struct {
	PayCenter string
	LiveRpc   string
}

type httpClient struct {
	PayCenter *bm.ClientConfig
	LiveRpc   *bm.ClientConfig
}
type gift struct {
	RechargeTip *rechargeTip
}

type rechargeTip struct {
	SilverTipDays []int
}

// Config .
type Config struct {
	Log            *log.Config
	BM             *bm.ServerConfig
	Verify         *verify.Config
	Tracer         *trace.Config
	Redis          *redis.Config
	Memcache       *memcache.Config
	MySQL          *sql.Config
	Ecode          *ecode.Config
	ResourceClient *warden.ClientConfig
	Auth           *auth.Config
	Warden         *warden.ClientConfig
	DM             *warden.ClientConfig
	Risk           *warden.ClientConfig
	VerifyConf     *warden.ClientConfig
	Host           host
	HTTPClient     *httpClient
	Gift           *gift
	LiveRpc        map[string]*liverpc.ClientConfig
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init init conf
func Init() error {
	if ConfPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(ConfPath, &Conf)
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
