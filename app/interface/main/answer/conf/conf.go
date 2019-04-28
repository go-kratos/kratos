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
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/netutil"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Config struct of conf.
type Config struct {
	App        *bm.App
	Host       *Host
	Log        *log.Config
	Tracer     *trace.Config
	Infoc2     *infoc.Config
	Ecode      *ecode.Config
	BM         *bm.ServerConfig
	HTTPClient *HTTPClient
	RPCClient  *RPC
	DataBus    *DataSource
	Mysql      *sql.Config
	Memcache   *Memcache
	Redis      *Redis
	AuthN      *auth.Config
	Antispam   *antispam.Config
	Geetest    *Geetest
	Answer     *Answer
	Question   *Question
	Backoff    *netutil.BackoffConfig
	Report     *databus.Config
	AccountRPC *warden.ClientConfig
	Captcha    *bm.ClientConfig
}

// RPC config
type RPC struct {
	Member  *rpc.ClientConfig
	Account *rpc.ClientConfig
}

// Answer conf.
type Answer struct {
	Captcha            bool // true:only use bili captcha
	Debug              bool
	Duration           int64
	BlockedTimestamp   int64
	BaseNum            int
	BaseExtraPassNum   int
	BaseExtraNoPassNum int
	ProNum             int
	BaseExtraScore     int
	BaseExtraPassCount int
	ExtraNum           int
	MaxRetries         int
	CaptchaTokenURL    string
	CaptchaVerifyURL   string
}

// Redis conf.
type Redis struct {
	*redis.Config
	Expire                time.Duration
	AnsCountExpire        time.Duration
	AnsAddFlagCountExpire time.Duration
}

// Memcache conf.
type Memcache struct {
	*memcache.Config
	Expire            time.Duration
	AnswerBolckExpire time.Duration
}

// Question conf.
type Question struct {
	// question total count tick
	TcQestTick   time.Duration
	RankQestTick time.Duration
}

// HTTPClient conf.
type HTTPClient struct {
	Normal *bm.ClientConfig
	Slow   *bm.ClientConfig
}

// Host conf.
type Host struct {
	Geetest  string
	Account  string
	ExtraIds string
	API      string
}

// Geetest geetest id & key
type Geetest struct {
	PC GeetestConfig
	H5 GeetestConfig
}

// GeetestConfig conf.
type GeetestConfig struct {
	CaptchaID  string
	PrivateKEY string
}

// App bilibili intranet authorization.
type App struct {
	Key    string
	Secret string
}

// DataSource .
type DataSource struct {
	ExtraAnswer *databus.Config
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
