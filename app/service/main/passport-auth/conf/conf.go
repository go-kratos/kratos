package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Config login conf
type Config struct {
	// tracer
	Tracer *trace.Config
	// xlog
	Xlog *log.Config
	// Identify
	VerifyConfig *verify.Config
	// db
	Mysql *sql.Config
	// mc
	Memcache *Memcache
	// app
	App *bm.App
	// BM
	BM *bm.ServerConfig
	// RPC
	RPCServer *rpc.ServerConfig
	// grpc
	WardenServer *warden.ServerConfig
	DC           *DC
	ServiceConf  *ServiceConf
}

// ServiceConf Switch
type ServiceConf struct {
	SupportOld bool
	Permit     map[string]string
}

// DC DC
type DC struct {
	Num  int
	Desc string
}

// Memcache cache
type Memcache struct {
	*memcache.Config
	Expire xtime.Duration
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
	if s, ok = client.Value("passport-auth-service.toml"); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init int config
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}
