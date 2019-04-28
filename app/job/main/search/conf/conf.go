package conf

import (
	"errors"
	"flag"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	hbase "go-common/library/database/hbase.v2"

	"github.com/BurntSushi/toml"
)

var (
	// ConfPath .
	ConfPath string
	client   *conf.Client
	// Conf .
	Conf = &Config{}
)

// Config .
type Config struct {
	// log
	XLog *log.Config
	// tracer
	Tracer *trace.Config
	// hbase
	HBase *HBaseConfig
	// business
	Business *Business
	// xhttp
	HTTPServer *bm.ServerConfig
	// http client
	HTTPClient *bm.ClientConfig
	// database
	DB map[string]*sql.Config
	// es cluster
	Es map[string]EsInfo
	// databus
	Databus map[string]*databus.Config
	// infoc
	InfoC map[string]*infoc.Config
	// sms
	SMS *SMS
}

// HBaseConfig combine with hbase.Config add ReadTimeout, WriteTimeout
type HBaseConfig struct {
	*hbase.Config
	// extra config
	ReadTimeout   xtime.Duration
	ReadsTimeout  xtime.Duration
	WriteTimeout  xtime.Duration
	WritesTimeout xtime.Duration
}

// Consumer .
type Consumer struct {
	GroupID string
	Topic   []string
	Offset  string
	Addrs   []string
}

// Business .
type Business struct {
	Env   string
	Index bool
}

// Redis search redis.
type Redis struct {
	*redis.Config
	Expire time.Duration
}

// EsInfo (deprecated).
type EsInfo struct {
	Addr []string
}

// SMS config
type SMS struct {
	Phone    string
	Token    string
	Interval int64
}

// init .
func init() {
	flag.StringVar(&ConfPath, "conf", "", "config path")
}

// Init .
func Init() (err error) {
	if ConfPath != "" {
		return local()
	}
	return remote()
}

// local .
func local() (err error) {
	_, err = toml.DecodeFile(ConfPath, &Conf)
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