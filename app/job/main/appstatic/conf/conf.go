package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	xlog "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Config .
type Config struct {
	Log          *xlog.Config
	Tracer       *trace.Config
	HTTPServer   *bm.ServerConfig
	HTTPClient   *bm.ClientConfig
	AppresClient *warden.ClientConfig
	MySQL        *sql.Config
	Cfg          *Cfg // push cfg.
	Redis        *redis.Config
	Bfs          *Bfs
}

// Cfg .
type Cfg struct {
	Push *PushCfg
	Diff *DiffCfg
	Grpc *GrpcCfg
}

// DiffCfg represents the diff calc config
type DiffCfg struct {
	FreDiff xtime.Duration // diff calculation frequency
	Folder  string
	Retry   string
}

// Bfs represents  the bfs config
type Bfs struct {
	Key     string
	Secret  string
	Host    string
	Timeout int
	OldURL  string
	NewURL  string
}

// PushCfg def.
type PushCfg struct {
	QPS       int            // qps limit
	Operation int            // operation number
	URL       string         // push url
	Timeout   int64          // timeout
	Fre       xtime.Duration // frequency
	Pause     xtime.Duration // pause between call app-resource and broadcast
}

// GrpcCfg def.
type GrpcCfg struct {
	ApiAppID string
	Method   string
}

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
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
