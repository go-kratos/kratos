package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
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
	Xlog *log.Config
	// Tracer tracer
	Tracer *trace.Config
	// DB db
	DB *DB
	// Compare compare
	Compare *Compare
	// InitCloud init cloud.
	InitCloud *InitCloud
	// BM
	BM *httpx.ServerConfig
}

// InitCloud init cloud conf.
type InitCloud struct {
	OffsetFilePath string
	UseOldOffset   bool

	Start, End int64

	Batch int

	Sleep time.Duration
}

// Compare compare
type Compare struct {
	Cloud2Local *CompareConfig
	Local2Cloud *CompareConfig
}

// CompareConfig compare proc config.
type CompareConfig struct {
	On    bool
	Debug bool

	OffsetFilePath string
	UseOldOffset   bool

	End bool

	StartTime     string
	EndTime       string
	DelayDuration time.Duration
	StepDuration  time.Duration
	LoopDuration  time.Duration

	BatchSize           int
	BatchMissRetryCount int

	Fix bool
}

// DB db config.
type DB struct {
	Local *sql.Config
	Cloud *sql.Config
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
