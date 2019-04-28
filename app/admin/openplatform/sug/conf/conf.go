package conf

import (
	"errors"
	"flag"

	"github.com/BurntSushi/toml"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/trace"
	xtime "go-common/library/time"
)

// global var
var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config config set
type Config struct {
	Log        *log.Config
	BM         *bm.ServerConfig
	Tracer     *trace.Config
	Redis      *redis.Config
	Ecode      *ecode.Config
	Auth       *permit.Config
	HTTPClient *bm.ClientConfig
	Es         *ElasticSearch
	DB         *DB
	URL        *URL
	Env        string
	SourcePath string
	Bfs        *Bfs
}

// URL search items
type URL struct {
	ItemSearch string
}

// DB db config
type DB struct {
	MallDB   *sql.Config
	TicketDB *sql.Config
}

// ElasticSearch config
type ElasticSearch struct {
	Addr []string
}

// Bfs config
type Bfs struct {
	Timeout     xtime.Duration
	MaxFileSize int
	Bucket      string
	Addr        string
	Key         string
	Secret      string
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
