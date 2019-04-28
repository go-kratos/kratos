package conf

import (
	"errors"
	"flag"
	"github.com/BurntSushi/toml"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/liverpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log       *log.Config
	BM        *bm.ServerConfig
	Verify    *verify.Config
	Tracer    *trace.Config
	Redis     *redis.Config
	Memcache  *memcache.Config
	MySQL     *sql.Config
	Ecode     *ecode.Config
	LiveRpc   map[string]*liverpc.ClientConfig
	DataBus   *DataBus
	GeeTest   *GeeTestConfig
	LogStream *LogStream
}

// DataBus ...
type DataBus struct {
	CaptchaAnti *databus.Config
}

// GeeTestConfig ...
type GeeTestConfig struct {
	On    int64      //是否开启极验
	Id    string     //公钥
	Key   string     //秘钥
	Qps   int64      //限制qps
	Slice int64      //限制qps的key分几份
	Get   HttpMethod //get配置
	Post  HttpMethod //post配置
}

// HttpMethod config
type HttpMethod struct {
	Timeout   int64
	KeepAlive int64
}

// LogStream config
type LogStream struct {
	LogId    string
	Token    string
	Address  string
	Capacity int
	Timeout  time.Duration
}

// init
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
