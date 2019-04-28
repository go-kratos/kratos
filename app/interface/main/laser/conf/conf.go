package conf

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/time"
)

var (
	// ConfPath local config path
	ConfPath string
	// Conf config
	Conf   = &Config{}
	client *conf.Client
)

// Config str
type Config struct {
	// base config
	// Elk config
	Log *log.Config
	// http config
	BM *bm.ServerConfig
	// db config
	Mysql *sql.Config
	// mc
	Memcache *Memcache
	// mail
	Mail *Mail
}

// Memcache conf.
type Memcache struct {
	Laser struct {
		*memcache.Config
		LaserExpire time.Duration
	}
}

// Mail conf.
type Mail struct {
	Host               string
	Port               int
	Username, Password string
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init init conf
func Init() (err error) {
	if ConfPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(ConfPath, &Conf)
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
		s        string
		ok       bool
		tempConf *Config
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tempConf); err != nil {
		return errors.New("could not decode toml config")
	}
	*Conf = *tempConf
	return
}
