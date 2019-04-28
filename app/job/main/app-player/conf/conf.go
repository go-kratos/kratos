package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

// is
var (
	confPath string
	client   *conf.Client
	Conf     = &Config{}
)

// Config is
type Config struct {
	// Env
	Env string
	// interface XLog
	XLog *log.Config
	// databus
	ArchiveNotifySub *databus.Config
	// http
	BM *bm.ServerConfig
	// mc
	Memcache *memcache.Config
	// rpc client
	ArchiveRPC *rpc.ClientConfig
	// redis
	Redis *redis.Config
	// Custom 自定义启动参数
	Custom *Custom
}

// Custom is
type Custom struct {
	Flush bool
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
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
	client.Watch("app-player-job.toml")
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
