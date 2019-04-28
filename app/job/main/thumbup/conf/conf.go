package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/tidb"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
	"go-common/library/sync/pipeline"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// Conf global variable.
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Config .
type Config struct {
	Redis                *Redis
	Tidb                 *tidb.Config
	ItemTidb             *tidb.Config
	Log                  *log.Config
	Databus              *Databus
	Thumbup              *Thumbup
	Memcache             *Memcache
	StatMerge            *StatMerge
	Merge                *pipeline.Config
	LikeDatabusutil      *databusutil.Config
	ItemLikesDatabusutil *databusutil.Config
	UserLikesDatabusutil *databusutil.Config
}

// StatMerge .
type StatMerge struct {
	Business string
	Target   int64
	Sources  []int64
}

// Redis .
type Redis struct {
	*redis.Config
	StatsExpire     xtime.Duration
	UserLikesExpire xtime.Duration
	ItemLikesExpire xtime.Duration
}

// Memcache .
type Memcache struct {
	*memcache.Config
	StatsExpire xtime.Duration
}

// Databus .
type Databus struct {
	Stat      *databus.Config
	Like      *databus.Config
	ItemLikes *databus.Config
	UserLikes *databus.Config
}

// Thumbup .
type Thumbup struct {
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init create config instance.
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
	str, ok := client.Toml2()
	if !ok {
		return errors.New("load config center error")
	}
	var tmpConf *Config
	if _, err = toml.Decode(str, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
