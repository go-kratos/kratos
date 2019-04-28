package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/orm"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
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

// Memcache memcache.
type Memcache struct {
	*memcache.Config
	Expire time.Duration
}

// Config config set
type Config struct {
	// elk
	Log *log.Config
	// http
	BM *bm.ServerConfig
	// ecode
	Ecode *ecode.Config

	HTTPClient *bm.ClientConfig

	Memcache *Memcache

	// orm
	ORM *orm.Config

	Mail *Mail

	Scheduler *Scheduler

	Auth *permit.Config

	Tapd *Tapd
}

// Mail mail
type Mail struct {
	Host        string
	Port        int
	Username    string
	Password    string
	NoticeOwner []string
}

// Scheduler Scheduler.
type Scheduler struct {
	UpdateHookURLCacheTask string
	Active                 bool
}

// Tapd Tapd.
type Tapd struct {
	CallbackToken string
	UseCache      bool
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
	return load()
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
