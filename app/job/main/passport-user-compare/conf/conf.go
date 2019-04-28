package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
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
	// log
	Log *log.Config
	// Tracer tracer
	Tracer *trace.Config
	// DB db
	DB *DB
	// BM
	BM *httpx.ServerConfig
	// FullTask
	FullTask *FullTask
	// IncTask
	IncTask *IncTask
	// WeChat
	WeChat *WeChat
	// data fix switch
	DataFixSwitch     bool
	IncrDataFixSwitch bool
	// httpClient
	HTTPClient *bm.ClientConfig
	// DuplicateTask
	DuplicateTask          *DuplicateTask
	FixEmailVerifiedSwitch bool
}

// FullTask full task compare
type FullTask struct {
	Switch         bool
	Step           int64
	AccountEnd     int64
	AccountInfoEnd int64
	AccountSnsEnd  int64
	CronFullStr    string
}

// WeChat wechat basic info
type WeChat struct {
	Token    string
	Secret   string
	Username string
}

// IncTask incr task compare
type IncTask struct {
	Switch       bool
	StartTime    string
	StepDuration time.Duration
	CronIncStr   string
}

// DB db config.
type DB struct {
	User   *sql.Config
	Origin *sql.Config
	Secret *sql.Config
}

// DuplicateTask check duplicatet ask
type DuplicateTask struct {
	Switch        bool
	DuplicateCron string
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
