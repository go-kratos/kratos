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
	// Xlog log
	Xlog *log.Config
	// Tracer tracer
	Tracer *trace.Config
	// DB db
	DB *DB
	// Memcache memcache
	Memcache *Memcache
	// Game
	Game *Game
	// HTTP
	BM *bm.ServerConfig
	// Group
	Group *Group
	// DataBus databus
	DataBus *DataBus
}

// Game game notify conf.
type Game struct {
	AppIDs      []int32
	DelCacheURI string
	Client      *bm.ClientConfig
}

// Group multi group config collection.
type Group struct {
	BinLog       *GroupConfig
	EncryptTrans *GroupConfig
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

// DataBus databus infomation
type DataBus struct {
	BinLogSub       *databus.Config
	EncryptTransSub *databus.Config
}

// DB db config.
type DB struct {
	Cloud *sql.Config
}

// Memcache general memcache config.
type Memcache struct {
	*memcache.Config
	Expire time.Duration
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
