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
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	Conf     = &Config{}
)

// Config def.
type Config struct {
	// base
	// log
	Log *log.Config
	// http
	BM     *bm.ServerConfig
	Figure figure
	// databus
	DataSource *DataSource
	// hbase
	HBase *HBaseConfig
	// redis
	Redis *Redis
	// mysql
	Mysql *sql.Config
}

type figure struct {
	Sync       bool
	SpyPath    string
	VipPath    string
	Lawful     int32
	Wide       int32
	Friendly   int32
	Bounty     int32
	Creativity int32
}

// DataSource config all figure job dataSource
type DataSource struct {
	AccountExp *databus.Config
	AccountReg *databus.Config
	Vip        *databus.Config
	Spy        *databus.Config
	Coin       *databus.Config
	ReplyInfo  *databus.Config
	Pay        *databus.Config
	Blocked    *databus.Config
	Danmaku    *databus.Config
}

// Redis conf.
type Redis struct {
	*redis.Config
	Expire         xtime.Duration
	WaiteMidExpire xtime.Duration
}

// HBaseConfig extra hbase config
type HBaseConfig struct {
	*hbase.Config
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
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
