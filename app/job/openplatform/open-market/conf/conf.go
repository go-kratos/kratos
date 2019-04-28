package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"

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
	// Xlog xlog conf
	Xlog *log.Config
	// DB multiple db conf
	DB *DB
	// HTTPClient httpClient conf
	HTTPClient *bm.ClientConfig
	//ES
	ElasticSearch    *ElasticSearch
	ElasticSearchUgc *ElasticSearch
	//comment
	Comment *Comment
	BM      *HTTPServers
	//berserker
	Berserker *Berserker
}

// HTTPServers bm inner config
type HTTPServers struct {
	Inner *bm.ServerConfig
}

// Berserker api conf
type Berserker struct {
	Appkey string
	Secret string
	URL    string
}

// ElasticSearch elasticSearch.
type ElasticSearch struct {
	Addr    []string
	Check   xtime.Duration
	Timeout string
}

// Comment config with url and type
type Comment struct {
	URL  string
	Type int
}

// DB config
type DB struct {
	TicketDB *sql.Config
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init conf.
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
	if s, ok = client.Value("open-market-job.toml"); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
