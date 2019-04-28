package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

// is
var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

// Config is
type Config struct {
	// Env
	Env string
	// interface XLog
	XLog *log.Config
	// host
	Host *Host
	// httpClinet
	HTTPClient *bm.ClientConfig
	// databus
	VideoupSub       *databus.Config
	DmSub            *databus.Config
	ArchiveResultPub *databus.Config
	DmPub            *databus.Config
	CacheSub         *databus.Config
	AccountNotifySub *databus.Config
	// rpc
	ArchiveServices []*rpc.ClientConfig
	Dm2RPC          *rpc.ClientConfig
	// mail
	Mail *Mail
	// DB
	DB *DB
	// BM
	BM *bm.ServerConfig
	// Redis
	Redis *redis.Config
	// ChanSize aid%ChanSize
	ChanSize    int
	PGCAsync    int
	UGCAsync    int
	MonitorSize int
	// qiye wechat
	WeChatToken  string
	WeChatSecret string
	WeChantUsers string
}

// Mail is
type Mail struct {
	Host               string
	Port               int
	Username, Password string
	Bangumi, Movie     []string
}

// Host is
type Host struct {
	APICo string
}

// DB is db config.
type DB struct {
	Archive *sql.Config
	Result  *sql.Config
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config.
func Init() (err error) {
	if confPath != "" {
		_, err = toml.DecodeFile(confPath, &Conf)
		return
	}
	err = remote()
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
		tmpConf = &Config{}
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
