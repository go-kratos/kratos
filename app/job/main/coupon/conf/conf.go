package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
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
	// base
	// elk
	Log *log.Config
	// http
	BM *bm.ServerConfig
	// memcache
	Memcache *memcache.Config
	// MySQL
	MySQL *sql.Config
	// Databus
	DataBus *DataSource
	// Properties
	Properties *Properties
	// Backoff retries config
	Backoff *netutil.BackoffConfig
	// http client
	HTTPClient  *bm.ClientConfig
	NewYearConf *NewYearConf
}

// NewYearConf .
type NewYearConf struct {
	ActID int64
}

// DataSource databus source
type DataSource struct {
	CouponBinlog *databus.Config
}

// Properties def.
type Properties struct {
	MaxRetries                  int
	BangumiNotifyURL            string
	CheckInUseCouponCron        string
	CheckInUseCouponCartoonCron string
	NotifyTimeInterval          xtime.Duration
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
