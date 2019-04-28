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
	"go-common/library/net/http/blademaster/middleware/permit"
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
	// base
	// elk
	Log *log.Config
	// http
	BM *bm.ServerConfig
	// MySQL
	MySQL *sql.Config
	// ecode
	Ecode *ecode.Config
	//auth
	Auth *permit.Config
	Prop *Properties
	// rpc client
	RPCClient2 *RPC
	// memcache
	Memcache *Memcache
	// BroadcastRPC grpc
	PGCRPC *warden.ClientConfig
	// http client
	HTTPClient *bm.ClientConfig
}

// Properties sysconf.
type Properties struct {
	MessageURL          string
	AllowanceTableCount int64
	SalarySleepTime     xtime.Duration
	SalaryMsgOpen       bool
	SalaryNormalMsgOpen bool
	MsgSysnSize         int
}

// Memcache memcache
type Memcache struct {
	*memcache.Config
	Expire xtime.Duration
}

// RPC config
type RPC struct {
	Coupon *rpc.ClientConfig
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
	if s, ok = client.Value("coupon-admin.toml"); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
