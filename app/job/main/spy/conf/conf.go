package conf

import (
	"errors"
	"flag"
	"time"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"go-common/library/database/hbase.v2"

	"github.com/BurntSushi/toml"
)

// Conf global config.
var (
	confPath string
	client   *conf.Client
	// config
	Conf = &Config{}
)

// Config info.
type Config struct {
	// base
	// customized property
	Property *Property
	// log
	Xlog *log.Config
	// db
	DB *sql.Config
	// databus
	Databus *DataSource
	// rpc to spy-service
	SpyRPC   *rpc.ClientConfig
	HBase    *HBaseConfig
	Memcache *memcache.Config
	// redis
	Redis *Redis
	// http client
	HTTPClient *bm.ClientConfig
	// HTTPServer
	HTTPServer *bm.ServerConfig
}

// HBaseConfig extra hbase config
type HBaseConfig struct {
	*hbase.Config
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DataSource config all spy job dataSource
type DataSource struct {
	EventData   *databus.Config
	SpyStatData *databus.Config
	SecLogin    *databus.Config
}

// Property config for biz logic.
type Property struct {
	TaskTimer       time.Duration
	UserInfoShard   int64
	Debug           bool
	ConfigLoadTick  xtime.Duration
	BlockTick       xtime.Duration
	BlockWaitTick   xtime.Duration
	LoadEventTick   xtime.Duration
	BlockAccountURL string
	HistoryShard    int64
	BlockEvent      int64
	Block           *struct {
		CycleTimes int64 // unit per seconds
		CycleCron  string
	}
	ReportCron     string
	ActivityEvents []string
}

// Redis conf.
type Redis struct {
	*redis.Config
	Expire        xtime.Duration
	MsgUUIDExpire xtime.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init init conf.
func Init() (err error) {
	if confPath == "" {
		return configCenter()
	}
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func configCenter() (err error) {
	if client, err = conf.New(); err != nil {
		panic(err)
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
