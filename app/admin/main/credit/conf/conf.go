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
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	// config
	confPath string
	client   *conf.Client
	Conf     = &Config{}
)

// Config struct of conf.
type Config struct {
	// base
	App *bm.App
	BFS *BFS
	// http
	BM *bm.ServerConfig
	// db
	ORM *ORM
	// host
	Host *Host
	// log
	Xlog *log.Config
	// tracer
	Tracer *trace.Config
	// rpc client
	RPCClient *RPC
	AccClient *warden.ClientConfig
	// http client
	HTTPClient *bm.ClientConfig
	// chan
	ChanSize   *ChanSize
	Auth       *permit.Config
	CoinClient *warden.ClientConfig
}

// ORM orm write and read config.
type ORM struct {
	Write *orm.Config
	Read  *orm.Config
}

// Host host config .
type Host struct {
	Bfs     string
	Manager string
	Msg     string
}

// ChanSize sysmsg channel size.
type ChanSize struct {
	SysMsg int64
}

// BFS bfs config
type BFS struct {
	Key    string
	Secret string
}

// RPC rpc client config.
type RPC struct {
	Account  *rpc.ClientConfig
	Relation *rpc.ClientConfig
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
