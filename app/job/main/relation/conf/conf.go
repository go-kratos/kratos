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
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

// Config struct of conf.
type Config struct {
	// base
	// app
	App *bm.App
	// Env
	Env string
	// goroutine sleep
	Tick time.Duration
	// log
	Xlog *log.Config
	// db
	Mysql *sql.Config
	// databus
	DataBus *databus.Config
	// httpClinet
	HTTPClient *bm.ClientConfig
	// clearPath
	ClearPath *ClearPath
	// apiPath
	ApiPath *ApiPath
	// sms monitor
	Sms *Sms
	// redis
	// Redis    *Redis
	RelRedis *Redis
	Memcache *Memcache
	Relation *Relation
	// bm
	BM *bm.ServerConfig
}

// ClearPath clear cache path
type ClearPath struct {
	Following string
	Follower  string
	Stat      string
}

// ApiPath api path collections
type ApiPath struct {
	FollowersNotify string
}

// Sms is sms monitor config.
type Sms struct {
	Phone string
	Token string
}

type Redis struct {
	*redis.Config
	Expire time.Duration
}

type Memcache struct {
	*memcache.Config
	Expire         time.Duration
	FollowerExpire time.Duration
}

// Relation relation related config
type Relation struct {
	// followers unread duration
	// FollowersUnread time.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

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
