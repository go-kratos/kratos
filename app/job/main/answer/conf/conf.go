package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf service config
	Conf = &Config{}
)

// Config def.
type Config struct {
	Log        *log.Config
	Databus    *DataSource
	Mysql      *sql.Config
	BFS        *BFS
	HTTPClient *bm.ClientConfig
	Properties *Properties
	Backoff    *netutil.BackoffConfig
	BM         *bm.ServerConfig
}

// BFS bfs config
type BFS struct {
	Timeout     xtime.Duration
	MaxFileSize int
	Bucket      string
	URL         string
	Method      string
	Key         string
	Secret      string
	Host        string
}

// Properties def.
type Properties struct {
	UploadInterval     xtime.Duration
	AccountIntranetURI string
	MaxRetries         int
	FontFilePath       string
}

// DataSource databus source
type DataSource struct {
	Labour  *databus.Config
	Account *databus.Config
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
