package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
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
	// base
	// elk
	Log *log.Config
	// App
	App *bm.App
	// rpc server2
	RPCServer *rpc.ServerConfig
	// tracer
	Tracer *trace.Config
	// auth
	Auth *auth.Config
	// verify
	Verify *verify.Config
	// HTTPServer
	HTTPServer *bm.ServerConfig
	// Ecode
	Ecode *ecode.Config
	// rpc
	FavoriteRPC *rpc.ClientConfig
	ArchiveRPC  *rpc.ClientConfig
	AccountRPC  *rpc.ClientConfig
	FilterRPC   *rpc.ClientConfig
	// databus
	ViewDatabus  *databus.Config
	ShareDatabus *databus.Config
	// Mysql
	Mysql *sql.Config
	// Redis
	Redis *Redis
	// HTTP client
	HTTPClient *bm.ClientConfig
	// Rule
	Rule *Rule
	// Host
	Host *Host
	// Warden Client
	ArcClient *warden.ClientConfig
	AccClient *warden.ClientConfig
}

// Host hosts.
type Host struct {
	Search string
}

// Redis redis struct
type Redis struct {
	*redis.Config
	PlExpire   time.Duration
	StatExpire time.Duration
}

// Rule playlist config
type Rule struct {
	MaxNameLimit      int
	MaxPlDescLimit    int
	MaxVideoDescLimit int
	MaxArcChangeLimit int
	MaxVideoCnt       int
	MaxPlCnt          int
	MaxPlArcsPs       int
	MaxPlsPageSize    int
	SortStep          int64
	MinSort           int64
	MaxSearchArcPs    int
	MaxSearchLimit    int
	BeginSort         int64
	PowerMids         []int64
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
