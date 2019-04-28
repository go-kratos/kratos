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
	// memcache
	Memcache *Memcache
	// MySQL
	MySQL *sql.Config
	// ecode
	Ecode *ecode.Config
	// biz property.
	Property *Property
	// rpc server
	RPCServer *rpc.ServerConfig
	// http client
	HTTPClient *bm.ClientConfig
	// vipinfo grpc
	VipinfoRPC  *warden.ClientConfig
	NewYearConf *NewYearConf
	// grpc server
	WardenServer *warden.ServerConfig
	Platform     map[string]string
}

// NewYearConf .
type NewYearConf struct {
	ActID               int64
	StartTime           int64
	EndTime             int64
	RandNum             int64
	NoVipBatchToken1    string
	NoVipBatchToken3    string
	NoVipBatchToken12   string
	More180BatchToken1  string
	More180BatchToken3  string
	More180BatchToken12 string
	Less180BatchToken1  string
	Less180BatchToken3  string
	Less180BatchToken12 string
	MonthBatchToken1    string
	MonthBatchToken3    string
	MonthBatchToken12   string
}

// Property def.
type Property struct {
	MessageURL       string
	CaptchaTokenURL  string
	CaptchaVerifyURL string
	CaptchaBID       string
}

// Memcache memcache
type Memcache struct {
	*memcache.Config
	Expire      xtime.Duration
	PrizeExpire xtime.Duration
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
