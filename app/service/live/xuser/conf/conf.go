package conf

import (
	"errors"
	"flag"
	"go-common/library/net/rpc/warden"

	"go-common/library/queue/databus"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	eCode "go-common/library/ecode/tip"
	"go-common/library/log"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	liverRPC "go-common/library/net/rpc/liverpc"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// TraceInit weather need init trace
	TraceInit bool
	client    *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log                 *log.Config
	BM                  *bm.ServerConfig
	BMClient            *bm.ClientConfig
	Verify              *verify.Config
	Tracer              *trace.Config
	VipRedis            *redis.Config
	GuardRedis          *redis.Config
	Redis               *redis.Config
	Memcache            *memcache.Config
	ExpMemcache         *memcache.Config
	LiveUserMysql       *sql.Config
	LiveAppMySQL        *sql.Config
	LiveAppORM          *orm.Config
	Ecode               *eCode.Config
	LiveVipChangePub    *databus.Config
	UserExpMySQL        *sql.Config
	LiveRPC             map[string]*liverRPC.ClientConfig
	LiveEntryEffectPub  *databus.Config
	GuardCfg            *GuardCfg
	AccountRPC          *rpc.ClientConfig
	Switch              *ConfigSwitch
	UserExpExpire       *UserExpExpireConf
	UserDaHangHaiExpire *UserDhhExpireConf
	// report
	Report        *databus.Config
	XanchorClient *warden.ClientConfig
}

// GuardCfg config for guard
type GuardCfg struct {
	OpenEntryEffectDatabus bool
	EnableGuardBroadcast   bool
	DanmuHost              string
}

// ConfigSwitch config for query
type ConfigSwitch struct {
	QueryExp int
}

// UserExpExpireConf config for cache expire
type UserExpExpireConf struct {
	ExpireTime int32
}

// UserDhhExpireConf config for cache expire
type UserDhhExpireConf struct {
	ExpireTime int32
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
	flag.BoolVar(&TraceInit, "traceInit", true, "default trace init")
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
