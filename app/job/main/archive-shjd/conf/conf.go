package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

type Config struct {
	// Env
	Env string
	// interface XLog
	Log *log.Config
	// databus
	Databus       *databus.Config
	ViewSub       *databus.Config
	DmSub         *databus.Config
	ReplySub      *databus.Config
	FavSub        *databus.Config
	CoinSub       *databus.Config
	ShareSub      *databus.Config
	RankSub       *databus.Config
	LikeSub       *databus.Config
	NotifyPub     *databus.Config
	AccountNotify *databus.Config
	// rpc
	ArchiveRPCs []*rpc.ClientConfig
	// http
	BM *bm.ServerConfig
	// redis
	Redis *redis.Config
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
