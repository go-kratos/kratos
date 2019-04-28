package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	Conf     = &Config{}
)

type Config struct {
	// Env
	Env          string
	WeChantUsers string
	WeChatToken  string
	WeChatSecret string
	// host
	Host *Host
	// interface XLog
	XLog *log.Config
	// databus
	ArchiveNotifySub *databus.Config
	AccountNotifySub *databus.Config
	StatViewSub      *databus.Config
	StatDMSub        *databus.Config
	StatReplySub     *databus.Config
	StatFavSub       *databus.Config
	StatCoinSub      *databus.Config
	StatShareSub     *databus.Config
	StatRankSub      *databus.Config
	StatLikeSub      *databus.Config
	ContributeSub    *databus.Config
	// http
	BM *HTTPServers
	// httpClinet
	HTTPClient     *bm.ClientConfig
	HTTPClientAsyn *bm.ClientConfig
	// mc
	Memcache *Memcache
	// rpc client
	ArchiveRPC *rpc.ClientConfig
	AccountRPC *rpc.ClientConfig
	ArticleRPC *rpc.ClientConfig
	// tick time
	Tick xtime.Duration
	// db
	MySQL *MySQL
	// redis
	Redis *Redis
	View  *View
}

// HTTPServers Http Servers
type HTTPServers struct {
	Inner *bm.ServerConfig
}

type Host struct {
	APP      string
	Config   string
	Hetongzi string
	APICo    string
	VC       string
}

type Memcache struct {
	Feed *struct {
		*memcache.Config
		ExpireMaxAid xtime.Duration
	}
}

type MySQL struct {
	Show *sql.Config
}

type Redis struct {
	Feed *struct {
		*redis.Config
	}
	Contribute *struct {
		*redis.Config
	}
}

type View struct {
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
