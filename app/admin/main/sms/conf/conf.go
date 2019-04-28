package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

// these are config global variables
var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

// Config is the model for parse workflow config
type Config struct {
	Log        *log.Config
	HTTPServer *bm.ServerConfig
	HTTPClient *bm.ClientConfig
	Auth       *permit.Config
	Tracer     *trace.Config
	DB         *orm.Config
	Sms        *sms
}

type sms struct {
	TplPsMax int
	MountDir string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init config
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
	err = load()
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
