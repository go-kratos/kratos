package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	xtime "go-common/library/time"

	"go-common/library/database/hbase.v2"

	"github.com/BurntSushi/toml"
)

// global var
var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// HTTPClient http client
type HTTPClient struct {
	Read  *bm.ClientConfig
	Write *bm.ClientConfig
}

// Config config set
type Config struct {
	// base
	Log *log.Config
	// http
	BM *bm.ServerConfig
	// auth
	Auth *permit.Config
	// MySQL
	MySQL *sql.Config
	//httpClient
	HTTPClient *HTTPClient
	//ORM
	ORM *orm.Config
	//bfs hosts
	BfsDownloadHost string
	BfsUpdateHost   string
	BfsDeleteHost   string
	// bfs
	Bfs *Bfs
	// Hbase
	Hbase *HBaseConfig
}

// Bfs .
type Bfs struct {
	BfsURL          string
	WaterMarkURL    string
	ImageGenURL     string
	TimeOut         xtime.Duration
	WmTimeOut       xtime.Duration
	ImageGenTimeOut xtime.Duration
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

// Item describe bfs bucket accessKey and accessSecret
type Item struct {
	Name      string // bucket name
	KeyID     string // accessKey
	KeySecret string // accessSecret
}

// HBaseConfig ...
type HBaseConfig struct {
	*hbase.Config
	WriteTimeout xtime.Duration
	ReadTimeout  xtime.Duration
}
