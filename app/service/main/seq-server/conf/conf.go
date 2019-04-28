package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

type Config struct {
	// base
	SeqSvrs []string
	Tick    time.Duration
	// rpc server
	RPCServer *rpc.ServerConfig
	// http
	BM *bm.ServerConfig
	// log
	Log *log.Config
	// tracer
	Tracer *trace.Config
	// db
	DB *DB
}

type DB struct {
	Number *sql.Config
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
