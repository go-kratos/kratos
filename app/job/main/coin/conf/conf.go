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
	"go-common/library/net/rpc/warden"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf config.
	Conf   = &Config{}
	client *conf.Client
)

// Config conf.
type Config struct {
	Env          string
	Sms          *Sms
	Databus      *databus.Config
	LoginDatabus *databus.Config
	ExpDatabus   *databus.Config
	Xlog         *log.Config
	DB           *DB
	CoinJob      *CoinJob
	AccountRPC   *warden.ClientConfig
	MemRPC       *warden.ClientConfig
	ArchiveRPC   *rpc.ClientConfig
	CoinRPC      *rpc.ClientConfig
	// BM
	BM *bm.ServerConfig
	// redis
	Redis       *redis.Config
	Databusutil *databusutil.Config
}

// CoinJob job conf.
type CoinJob struct {
	// award conf
	StartTime   int64
	Start       bool
	LoginExpire xtime.Duration
}

// Sms sms conf.
type Sms struct {
	Phone string
	Token string
}

// DB db conf
type DB struct {
	Coin *sql.Config
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
	err = load()
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
