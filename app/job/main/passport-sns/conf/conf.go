package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
	xtime "go-common/library/time"

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
	// http
	BM *bm.ServerConfig
	// DB
	DB       *DB
	Memcache *Memcache
	//Databus databus
	DataBus *DataBus
	// DataUtil config
	DatabusUtil *databusutil.Config
	// SyncConf
	SyncConf *SyncConf
}

// DB db config
type DB struct {
	Aso *sql.Config
	Sns *sql.Config
}

// Memcache memcache
type Memcache struct {
	*memcache.Config
	Expire xtime.Duration
}

// DataBus databus.
type DataBus struct {
	SnsLogSub    *databus.Config
	AsoBinLogSub *databus.Config
}

// SyncConf sync conf
type SyncConf struct {
	IncSwitch   bool
	FullSwitch  bool
	CheckSwitch bool
	ChanNum     int
	ChanSize    int
	CheckTicker xtime.Duration
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
