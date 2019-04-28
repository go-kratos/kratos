package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf init Config struct.
	Conf   = &Config{}
	client *conf.Client
)

// Config struct.
type Config struct {
	// Env
	Env string
	// tracer
	Tracer *trace.Config
	// xlog
	XLog *log.Config
	// http
	BM *HTTPServers
	// ecode
	Ecode *ecode.Config
	// db
	MySQL *MySQL
	// Bfs
	Bfs *Bfs
	// Feedback
	Feedback *Feedback
	// rpc
	LocationRPC *rpc.ClientConfig
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
	Local *bm.ServerConfig
}

// MySQL struct.
type MySQL struct {
	Master *sql.Config
}

// Bfs struct.
type Bfs struct {
	Addr        string
	Bucket      string
	Key         string
	Secret      string
	MaxFileSize int
}

// Feedback struct.
type Feedback struct {
	ReplysNum      int
	MaxContentSize int
	ImgLimit       int
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init init conf
func Init() error {
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
