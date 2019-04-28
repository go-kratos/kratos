package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf conf.
	Conf   = &Config{}
	client *conf.Client
)

// Config config.
type Config struct {
	// log
	Xlog *log.Config
	// Tracer tracer
	Tracer *trace.Config
	// HTTP
	BM *bm.ServerConfig
	// Group group
	Group *Group
	//Databus databus
	DataBus *DataBus
}

// Group multi group config collection.
type Group struct {
	AsoBinLog *GroupConfig
}

// GroupConfig group config.
type GroupConfig struct {
	// Size merge size
	Size int
	// Num merge goroutine num
	Num int
	// Ticker duration of submit merges when no new message
	Ticker time.Duration
	// Chan size of merge chan and done chan
	Chan int
}

// DataBus databus.
type DataBus struct {
	AsoBinLogSub    *databus.Config
	EncryptTransPub *databus.Config
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config.
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
