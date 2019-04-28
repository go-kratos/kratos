package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"
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
	Log      *log.Config
	BM       *bm.ServerConfig
	Verify   *verify.Config
	Tracer   *trace.Config
	Memcache *memcache.Config
	MySQL    *sql.Config
	GRPC     *warden.ServerConfig

	HTTPClient *bm.ClientConfig

	MVIPClient *warden.ClientConfig
	ACCClient  *warden.ClientConfig

	YST *YstConfig
	PAY *PayConfig

	Ticker *TickerConfig

	MVIP *MVIPConfig

	CacheTTL *CacheTTL
}

// CacheTTL contains ttl configs.
type CacheTTL struct {
	UserInfoTTL int32
	PayParamTTL int32
	LockTTL     int32
}

// TickerConfig contains durations configs of ticker proc.
type TickerConfig struct {
	PanelRefreshDuration   string
	UnpaidDurationStime    string
	UnpaidDurationEtime    string
	UnpaidRefreshDuratuion string
}

// MVIPConfig contains mvip configs.
type MVIPConfig struct {
	BatchIdsMap      map[string]int
	BatchUserInfoUrl string
}

// YstConfig contains yst configs.
type YstConfig struct {
	Domain string
	Key    string
}

// PayConfig contains pay configs.
type PayConfig struct {
	PayExpireDuration     string
	RenewFromDuration     string
	RenewToDuration       string
	OrderRateFromDuration xtime.Duration
	OrderRateMaxNumber    int

	QrURL      string
	GuestQrURL string
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
