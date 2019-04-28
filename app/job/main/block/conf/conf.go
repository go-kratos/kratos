package conf

import (
	"errors"
	"flag"

	"go-common/app/job/main/block/model"
	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
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
	Log           *log.Config
	Memcache      *memcache.Config
	DB            *sql.Config
	BM            *bm.ServerConfig
	HTTPClient    *bm.ClientConfig
	Databus       *Databus
	AccountNotify *databus.Config

	Property *Property
	// manager log config
	ManagerLog *databus.Config
}

// Databus .
type Databus struct {
	Credit *databus.Config
}

// Property .
type Property struct {
	LimitExpireCheckLimit  int
	LimitExpireCheckTick   xtime.Duration
	CreditExpireCheckLimit int
	CreditExpireCheckTick  xtime.Duration
	MSGURL                 string
	MSG                    *MSG
	Flag                   *struct {
		ExpireCheck bool
		CreditSub   bool
	}
}

// MSG .
type MSG struct {
	BlockRemove model.MSG
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
