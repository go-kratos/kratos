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
	"go-common/library/net/rpc/warden"
	"go-common/library/queue/databus"

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
	Log            *log.Config
	BM             *bm.ServerConfig
	Redis          *redis.Config
	Memcache       *memcache.Config
	MemcacheRank   *memcache.Config
	MySQL          *sql.Config
	MySQLRank      *sql.Config
	MySQLRankOld   *sql.Config
	GRPCUGCPayRank *warden.ClientConfig
	Ecode          *ecode.Config
	BinlogMQ       *databus.Config
	ElecBinlogMQ   *databus.Config
	CacheTTL       *CacheTTL
	Biz            *Biz
}

// Biz .
type Biz struct {
	RunCASTimes    int64
	AccountUserMin int64
	Cron           *struct {
		TaskDailyBill     string
		TaskAccountUser   string
		TaskAccountBiz    string
		TaskMonthlyBill   string
		TaskRechargeShell string
	}
	Tax *struct {
		AssetRate float64
	}
	Pay struct {
		ID                  string
		Token               string
		CheckOrderURL       string
		CheckRefundOrderURL string
		RechargeShellURL    string
		RechargeCallbackURL string
		OrderQueryURL       string
	}
	Task *struct {
		DailyBillPrefix     string
		DailyBillOffset     int
		AccountUserPrefix   string
		AccountBizPrefix    string
		MonthBillPrefix     string
		MonthBillOffset     int
		RechargeShellPrefix string
		RechargeShellOffset int
	}
}

// CacheTTL .
type CacheTTL struct {
	ElecOrderIDTTL int32
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
