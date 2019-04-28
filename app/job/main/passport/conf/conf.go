package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"go-common/library/database/hbase.v2"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Config struct of conf.
type Config struct {
	// base
	// Tracer tracer
	Tracer *trace.Config
	// Xlog log
	Xlog *log.Config
	// BM
	BM *blademaster.ServerConfig
	// URI uri
	URI *URI
	// Game game
	Game *Game
	// RPC rpc
	RPC *RPC
	// DB db
	DB *DB
	// HTTPClient httpx client
	HTTPClient *httpx.ClientConfig
	// Group group.
	Group *Group
	// DataBus databus
	DataBus *DataBus
	// HBase hbase
	HBase  *HBase
	Encode *Encode
	Sync   *Sync
}

// Encode encode
type Encode struct {
	AesKey string
	Salt   string
}

// HBase multi hbase.
type HBase struct {
	LoginLog *HBaseConfig
	PwdLog   *HBaseConfig
}

// HBaseConfig .
type HBaseConfig struct {
	*hbase.Config
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// URI multi uri.
type URI struct {
	SetToken string
	DelCache string
}

// RPC multi rpc conf collection.
type RPC struct {
	IdentifyGame *rpc.ClientConfig
}

// Game game notify conf.
type Game struct {
	AppIDs      []int32
	DelCacheURI string
	Client      *httpx.ClientConfig
}

// Group multi group config collection.
type Group struct {
	AsoBinLog      *GroupConfig
	User           *GroupConfig
	Log            *GroupConfig
	ContactBindLog *GroupConfig
	PwdLog         *GroupConfig
	AuthBinLog     *GroupConfig
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

// DataBus multi databus collection.
type DataBus struct {
	AsoBinLog      *databus.Config
	User           *databus.Config
	Log            *databus.Config
	ContactBindLog *databus.Config
	UserLog        *databus.Config
	PwdLog         *databus.Config
	AuthBinLog     *databus.Config
}

// DB db config.
type DB struct {
	Log *sql.Config
	ASO *sql.Config
}

// Sync config.
type Sync struct {
	SyncPwdID int64
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
