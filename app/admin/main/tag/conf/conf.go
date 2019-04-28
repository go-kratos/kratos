package conf

import (
	"errors"
	"flag"

	"go-common/app/admin/main/tag/model"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// Conf Conf.
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Tag Tag.
type Tag struct {
	MaxSubmitNum  int
	HotSyncTime   time.Duration
	TraversalTime time.Duration
	SelectTime    int
}

// Redis Redis.
type Redis struct {
	Tag     *redis.Config
	TagRank *redis.Config
}

// Memcache Memcache.
type Memcache struct {
	Tag           *memcache.Config
	ChannelExpire time.Duration
}

// Config config.
type Config struct {
	Log        *log.Config
	HTTPServer *bm.ServerConfig
	Perms      map[string]string
	Verify     *verify.Config
	Auth       *permit.Config
	Tracer     *trace.Config
	Mysql      *sql.Config
	Ecode      *ecode.Config
	HTTPClient *bm.ClientConfig
	ES         *elastic.Config
	Redis      *Redis
	Memcache   *Memcache
	Hosts      *model.DependServiceHost
	Tag        *Tag
	AccClient  *warden.ClientConfig
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init conf.
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
		tomlStr string
		ok      bool
		tmpConf *Config
	)
	if tomlStr, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(tomlStr, &tmpConf); err != nil {
		return errors.New("could not decode toml config")
	}
	*Conf = *tmpConf
	return
}
