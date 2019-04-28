package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/orm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf .
	Conf   = &Config{}
	client *conf.Client
)

// Config struct
type Config struct {
	// base
	// log
	Log *log.Config
	// http
	BM *bm.ServerConfig
	// verify
	Verify *verify.Config
	// db
	ORM *ORM
	// host
	Host *Host
	// http client test
	HTTPClient *HTTPClient
	// archive rpc
	ArchiveRPC *warden.ClientConfig
	// es
	Elastic *elastic.Config
}

// ORM struct
type ORM struct {
	Write *orm.Config
}

// Host struct
type Host struct {
	ServiceURI string
	ManagerURI string
}

// HTTPClient struct
type HTTPClient struct {
	Sobot *bm.ClientConfig
	Read  *bm.ClientConfig
	Write *bm.ClientConfig
	Audit *bm.ClientConfig
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init .
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
