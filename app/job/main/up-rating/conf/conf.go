package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	// ConfPath local config path
	confPath string
	client   *conf.Client
	// Conf is global config object.
	Conf = &Config{}
)

// Config is project all config
type Config struct {
	// log
	Log *log.Config
	// Mysql
	MySQL *MySQL
	// tracer
	Tracer *trace.Config
	// http client
	HTTPClient *bm.ClientConfig
	// bm
	BM *bm.ServerConfig
	// concurrent
	Con *Concurrent
}

// MySQL mysql config
type MySQL struct {
	Rating *sql.Config
}

// Concurrent concurrent compute
type Concurrent struct {
	Concurrent int
	Limit      int
}

// MailAddr mail send addr.
type MailAddr struct {
	Type int
	Addr []string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config.
func Init() (err error) {
	if confPath != "" {
		_, err = toml.DecodeFile(confPath, &Conf)
		return
	}
	err = remote()
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
