package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"

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
	Redis    *redis.Config
	Memcache *memcache.Config
	MySQL    *sql.Config
	Ecode    *ecode.Config
	CacheTTL *CacheTTL
	Biz      *Biz
}

// CacheTTL .
type CacheTTL struct {
	OrderTTL                 int32
	AssetTTL                 int32
	AssetRelationTTL         int32
	AggrIncomeUserTTL        int32
	AggrIncomeUserMonthlyTTL int32
}

// Biz .
type Biz struct {
	RunCASTimes int64
	Pay         struct {
		ID                string
		Token             string
		OrderTTL          int32
		URLQuery          string
		URLRefund         string
		URLCancel         string
		URLPayCallback    string
		URLRefundCallback string
	}
	Price struct {
		PlatformTax map[string]float64
	}
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
