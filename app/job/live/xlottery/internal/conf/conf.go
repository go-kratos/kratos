package conf

import (
	"errors"
	"flag"

	"go-common/library/net/rpc/liverpc"

	"go-common/library/database/sql"

	"go-common/library/queue/databus"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log      *log.Config
	BM       *bm.ServerConfig
	Verify   *verify.Config
	Tracer   *trace.Config
	Redis    *Redis
	Database *Database
	Ecode    *ecode.Config
	Cfg      *Cfg
	// databus
	GiftPaySub    *databus.Config
	GiftFreeSub   *databus.Config
	AddCapsuleSub *databus.Config
	UserReport    *databus.Config
	LiveRpc       map[string]*liverpc.ClientConfig
	HTTPClient    *bm.ClientConfig
	CouponConf    *CouponConfig
}

// CouponConfig .
type CouponConfig struct {
	Url    string
	Coupon map[string]string
}

// Database mysql
type Database struct {
	Lottery *sql.Config
}

// Redis redis
type Redis struct {
	Lottery *redis.Config
}

// Cfg def
type Cfg struct {
	// ExpireCountFrequency crontab frequency
	ExpireCountFrequency string
	CouponRetryFrequency string
	ConsumerProcNum      int64
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
