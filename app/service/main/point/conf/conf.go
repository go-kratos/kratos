package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"

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
	// http
	BM *bm.ServerConfig
	// Verify
	Verify *verify.Config
	// Auth
	Auth *auth.Config
	// memcache
	Memcache *Memcache
	// MySQL
	MySQL *sql.Config
	// ecode
	Ecode *ecode.Config
	// property
	Property *Property
	// rpc server
	RPCServer *rpc.ServerConfig
	// grpc server
	WardenServer *warden.ServerConfig
}

// Property config for biz logic.
type Property struct {
	PointActiveStartDate string
	PointActiveEndDate   string
	ActivityAllowAppID   []int64
	ConfigLoadTick       xtime.Duration
	PointGetRule         map[string]int64
	PointWhiteAppkeys    string
}

// Memcache memcache
type Memcache struct {
	*memcache.Config
	Expire xtime.Duration
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
