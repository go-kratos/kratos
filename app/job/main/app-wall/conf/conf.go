package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
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
	// redis
	Redis *Redis
	// memcache
	Memcache *Memcache
	// db
	MySQL *MySQL
	// databus
	ReportDatabus *databus.Config
	// chan number
	ChanNum int64
	// chan db number
	ChanDBNum int64
	// ecode
	Ecode *ecode.Config
	// Report
	Report *databus.Config
	// client
	Consumer *Consumer
	// HTTPClient
	HTTPClient *bm.ClientConfig
	// HTTPUnicom
	HTTPUnicom *bm.ClientConfig
	// host
	Host *Host
	// unicom
	Unicom *Unicom
	// infoc2
	UnicomUserInfoc2 *infoc.Config
	UnicomPackInfoc  *infoc.Config
	// tick time
	Tick xtime.Duration
	// monthly
	Monthly bool
	// seq
	SeqRPC *rpc.ClientConfig
	// Seq
	Seq *Seq
}

type Seq struct {
	BusinessID int64
	Token      string
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

type MySQL struct {
	Show *sql.Config
}

type Unicom struct {
	PackKeyExpired xtime.Duration
	KeyExpired     xtime.Duration
}

type Memcache struct {
	Operator *struct {
		*memcache.Config
		Expire xtime.Duration
	}
}

type Redis struct {
	Feed *struct {
		*redis.Config
	}
}

type Consumer struct {
	Group   string
	Topic   string
	Offset  string
	Brokers []string
}

type Host struct {
	APP        string
	UnicomFlow string
	Unicom     string
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
