package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// global var
var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config config set
type Config struct {
	// base
	// elk
	Log *log.Config
	// http
	BM *HTTPServers
	// tracer
	Tracer *trace.Config
	// ecode
	Ecode *ecode.Config

	DonwloadHost    string
	UploadHost      string
	UploadAdminHost string
	RoutineCount    int
	HTTPClient      *bm.ClientConfig // call upload admin
	HTTPTimeout     xtime.Duration
	Databus         *databus.Config
	AIYellowing     AIYellowing
	Threshold       *Threshold
}

// AIYellowing describe AI upload config
type AIYellowing struct {
	ExceptBuckets []string
	Producer      *databus.Config
	Consumer      *databus.Config
}

// Threshold .
type Threshold struct {
	Sex      int
	Politics int
	Blood    int
	Violent  int
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
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
