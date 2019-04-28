package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"go-common/library/database/hbase.v2"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// HBaseConfig ...
type HBaseConfig struct {
	*hbase.Config
	WriteTimeout xtime.Duration
	ReadTimeout  xtime.Duration
}

// Config struct of conf.
type Config struct {
	// base
	// log
	Xlog *log.Config
	// tracer
	Tracer *trace.Config
	// identify
	Identify *verify.Config
	// BM
	BM *blademaster.ServerConfig
	// Switch switch
	Switch *Switch
	// RPCServer rpc server2
	RPCServer *rpc.ServerConfig
	// HBase
	HBase *HBase
}

// Switch switch.
type Switch struct {
	LoginLogHBase bool
	RPC           bool
}

// HBase multi hbase.
type HBase struct {
	FaceApply *HBaseConfig
	LoginLog  *HBaseConfig
	PwdLog    *HBaseConfig
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
