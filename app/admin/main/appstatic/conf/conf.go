package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	// config
	confPath string
	client   *conf.Client
	// Conf .
	Conf = &Config{}
)

// Config def.
type Config struct {
	// base
	// auth
	Auth *permit.Config
	// http
	BM *bm.ServerConfig
	// db
	ORM *orm.Config
	// log
	XLog *log.Config
	// tracer
	Tracer *trace.Config
	// cfg
	Cfg *Cfg
	// bfs
	Bfs *Bfs
	// nas
	Nas *Bfs
	// HTTPClient .
	HTTPClient *bm.ClientConfig
	// Redis
	Redis *Redis
}

// Redis redis
type Redis struct {
	*redis.Config
}

// Bfs reprensents the bfs config
type Bfs struct {
	Key     string
	Secret  string
	Host    string
	Timeout int
	OldURL  string
	NewURL  string
}

// Cfg def.
type Cfg struct {
	HistoryVer int
	Storage    string   // NAS or BFS
	Filetypes  []string // allowed file type to upload
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
