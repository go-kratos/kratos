package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf .
	Conf   = &Config{}
	client *conf.Client
)

// Pagination .
type Pagination struct {
	PageNum     int
	PageSize    int
	MaxPageNum  int
	MaxPageSize int
}

// Config .
type Config struct {
	// base
	// xlog
	Xlog *log.Config
	// rpc server2
	//RPCServer2 *conf.RPCServer2
	// tracer
	Tracer *trace.Config
	// xhttp
	HTTPServer *bm.ServerConfig
	// es cluster
	Es map[string]*EsInfo
	// ecode
	Ecode *ecode.Config
	// pagination
	Pagination *Pagination
	// httpClinet
	HTTPClient *bm.ClientConfig
	SMS        *SMS
}

// EsInfo .
type EsInfo struct {
	Addr []string
}

// SMS config
type SMS struct {
	Phone    string
	Token    string
	Interval int64
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
		tmpConf = &Config{}
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}