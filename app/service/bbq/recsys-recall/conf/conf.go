package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log          *log.Config
	BM           *bm.ServerConfig
	Verify       *verify.Config
	Tracer       *trace.Config
	Redis        *redis.Config
	BFRedis      *redis.Config
	MySQL        *sql.Config
	Ecode        *ecode.Config
	WorkPool     *WorkPoolConfig
	ForwardIndex *ForwardIndexConfig
	LocalCache   *LocalCacheConfig
}

// WorkPoolConfig .
type WorkPoolConfig struct {
	Capacity       uint64
	MaxWorkers     uint64
	MaxIdleWorkers uint64
	MinIdleWorkers uint64
	KeepAlive      xtime.Duration
}

// ForwardIndexConfig .
type ForwardIndexConfig struct {
	LocalPath      string
	RemotePath     string
	MD5Path        string
	Protocol       string
	ReloadDucation xtime.Duration
}

// LocalCacheConfig .
type LocalCacheConfig struct {
	L1Tags []string
	Level1 xtime.Duration
	L2Tags []string
	Level2 xtime.Duration
	Level3 xtime.Duration
	MaxAge xtime.Duration
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
