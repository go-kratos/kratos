package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// config
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config def.
type Config struct {
	// base
	// db
	DB *DB
	// spy rpc client
	SpyRPC     *rpc.ClientConfig
	AccountRPC *rpc.ClientConfig
	// memcache
	Memcache *Memcache
	// log
	Log *log.Config
	// customized property
	Property   *Property
	HTTPServer *bm.ServerConfig
}

// DB config.
type DB struct {
	Spy *sql.Config
}

// Memcache config.
type Memcache struct {
	User       *memcache.Config
	UserExpire time.Duration
}

// Property config for biz logic.
type Property struct {
	TelValidateURL  string
	BlockAccountURL string
	UserInfoShard   int64
	HistoryShard    int64
	LoadEventTick   time.Duration
	Score           *struct {
		BaseInit  int8
		EventInit int8
	}
	Punishment *struct {
		ScoreThreshold int8
		Times          int8
	}
	Event *struct {
		ServiceName           string
		BindMailAndTelLowRisk string
		BindMailOnly          string
		BindNothing           string
		BindTelLowRiskOnly    string
		BindTelMediumRisk     string
		BindTelHighRisk       string
		BindTelUnknownRisk    string
	}
	// activity events
	ActivityEvents map[int32]struct{}
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
