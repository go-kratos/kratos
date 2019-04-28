package conf

import (
	"flag"

	"errors"
	"github.com/BurntSushi/toml"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	xtime "go-common/library/time"
)

// ConfPath str.
var (
	ConfPath string
	Conf     = &Config{}
	client   *conf.Client
)

// Config app meta config.
type Config struct {
	// base
	// ELK
	Log *log.Config
	// Auth
	Auth *permit.Config
	// http
	BM *bm.ServerConfig
	// mysql
	Mysql *sql.Config
	// http client
	HTTPClient *bm.ClientConfig
	// host
	Host *Host
	// Memcache
	Memcache *Memcache
}

// Memcache conf.
type Memcache struct {
	Laser struct {
		*memcache.Config
		Expire xtime.Duration
	}
}

// Host conf.
type Host struct {
	Manager string
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init fn.
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
		cf      string
		ok      bool
		tmpConf *Config
	)
	if cf, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(cf, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
