package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
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
	// env
	Env string
	// elk
	Log *log.Config
	// HTTPClient .
	HTTPClient *bm.ClientConfig
	// http
	BM *HTTPServers
	// tracer
	Tracer *trace.Config
	// MySQL
	MySQL *sql.Config
	// hbase
	HBase          *HBaseConfig
	BlackListHBase *HBaseConfig
	// databuse
	LiveRoomSub   *databus.Config // 开播提醒
	LiveCommonSub *databus.Config // 直播通用
	// push
	Push *push
	// redis
	Redis *Redis
}

// HBaseConfig extra hbase config for compatible
type HBaseConfig struct {
	*hbase.Config
	WriteTimeout xtime.Duration
	ReadTimeout  xtime.Duration
}

// HTTPServers Http Servers
type HTTPServers struct {
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

type push struct {
	MultiAPI           string
	AppID              int
	BusinessID         int
	BusinessToken      string
	LinkType           int
	PushRetryTimes     int
	PushOnceLimit      int
	DefaultCopyWriting string
	SpecialCopyWriting string
	ConsumerProcNum    int
	IntervalLimit      int
	PushFilterIgnores  struct {
		Smooth, Limit []int
	}
}

// Redis Redis.PushInterval config
type Redis struct {
	PushInterval *struct {
		*redis.Config
		Expire xtime.Duration
	}
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
