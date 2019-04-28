package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	v "go-common/library/net/http/blademaster/middleware/verify"
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

	Log           *log.Config
	BM            *bm.ServerConfig
	HTTPClient    *bm.ClientConfig
	RPCServer     *rpc.ServerConfig
	WardenServer  *warden.ServerConfig
	AccMysql      *sql.Config
	Mysql         *sql.Config
	Memcache      *memcache.Config
	CacheTTL      *CacheTTL
	Tracer        *trace.Config
	Databus       *databus.Config
	Redis         *redis.Config
	AccountNotify *databus.Config
	ReportUser    *databus.Config
	ReportManager *databus.Config
	Host          *Host
	Verify        *v.Config
	// realname
	RealnameProperty *RealnameProperty
	// block
	BlockMySQL    *sql.Config
	BlockMemcache *memcache.Config
	BlockProperty *BlockProperty
	BlockCacheTTL *BlockCacheTTL
}

// Host is
type Host struct {
	Search string
}

// CacheTTL cache live time.
type CacheTTL struct {
	BaseTTL            time.Duration
	MoralTTL           time.Duration
	CaptureTimesTTL    time.Duration
	CaptureCodeTTL     time.Duration
	CaptureErrTimesTTL time.Duration
	ApplyInfoTTL       time.Duration
}

// RealnameProperty .
type RealnameProperty struct {
	IMGURLTemplate string
}

// BlockProperty .
type BlockProperty struct {
	MSGURL    string
	WhiteList []int64
}

// BlockCacheTTL is
type BlockCacheTTL struct {
	UserTTL     time.Duration
	UserMaxRate float64
	UserT       float64
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
	// go func() {
	// 	for range client.Event() {
	// 		log.Info("config reload")
	// 		if load() != nil {
	// 			log.Error("config reload error (%v)", err)
	// 		}
	// 	}
	// }()
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

// RsaPub is realname rsa pub key
func RsaPub() (key string) {
	if client == nil {
		return ""
	}
	key, _ = client.Value("realname.rsa.pub")
	return
}

// RsaPriv is realname rsa priv key
func RsaPriv() (key string) {
	if client == nil {
		return ""
	}
	key, _ = client.Value("realname.rsa.priv")
	return
}
