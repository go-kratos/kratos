package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// CommonConfig .
type CommonConfig struct {
	ExpireTime int64
}

// LRUConfig .
type LRUConfig struct {
	InstCnt  int
	Bucket   int
	Capacity int
	Timeout  int64
}

type AreaInfo struct {
	Name string
}

// Config .
type Config struct {
	Log                   *log.Config
	BM                    *bm.ServerConfig
	Verify                *verify.Config
	Tracer                *trace.Config
	Redis                 *redis.Config
	MySQL                 *sql.Config
	LiveAppMySQL          *sql.Config
	Ecode                 *ecode.Config
	LiveDanmuSub          *databus.Config //发送弹幕回调配置
	LiveGiftSendByPaySub  *databus.Config
	LiveGiftSendByFreeSub *databus.Config
	LiveGuardBuySub       *databus.Config
	LivePopularitySub     *databus.Config
	LiveValidLiveDaysSub  *databus.Config
	LiveRoomTagSub        *databus.Config
	LiveRankListSub       *databus.Config
	Common                *CommonConfig
	LRUCache              *LRUConfig
	FirstAreas            map[string]*AreaInfo
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
