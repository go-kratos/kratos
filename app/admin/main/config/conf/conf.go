package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf conf
	Conf   = &Config{}
	client *conf.Client
)

// Config config
type Config struct {
	// identify
	Auth *permit.Config
	// log
	Log *log.Config
	// orm
	ORM *orm.Config
	// http client
	HTTPClient *bm.ClientConfig
	//ConfSvr
	ConfSvr *rpc.ClientConfig
	//tree
	Tree *ServiceTree
	// db apm
	ORMApm *orm.Config
	//BM
	BM *bm.ServerConfig
}

// ServiceTree service tree.
type ServiceTree struct {
	Platform string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config file")
}

// Init init.
func Init() (err error) {
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
			if err := load(); err != nil {
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
