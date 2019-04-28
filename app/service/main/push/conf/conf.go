package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
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
	Env         string
	Ecode       *ecode.Config
	Log         *log.Config
	HTTPServer  *bm.ServerConfig
	HTTPClient  *bm.ClientConfig
	RPCServer   *rpc.ServerConfig
	GRPC        *warden.ServerConfig
	FilterRPC   *rpc.ClientConfig
	Tracer      *trace.Config
	Verify      *verify.Config
	App         *bm.App
	Redis       *rds
	Memcache    *mc
	MySQL       *sql.Config
	ReportPub   *databus.Config
	CallbackPub *databus.Config
	Android     *android
	Apns        *apns
	Push        *push
}

// mc config
type mc struct {
	*memcache.Config
	SettingExpire xtime.Duration
	ReportExpire  xtime.Duration
	UUIDExpire    xtime.Duration
}

type rds struct {
	*redis.Config
	TokenExpire xtime.Duration
	LaterExpire xtime.Duration
	MidsExpire  xtime.Duration
}

type android struct {
	PoolSize       int
	Timeout        xtime.Duration
	PushHuaweiPart int
	MiUseVip       int
}

type apns struct {
	PoolSize    int
	Proxy       int
	ProxySocket string
	Timeout     xtime.Duration
	Deadline    xtime.Duration
}

type push struct {
	PickUpTask                               bool
	LoadBusinessInteval                      xtime.Duration
	LoadTaskInteval                          xtime.Duration
	UpdateTaskProgressInteval                xtime.Duration
	PushChanSizeAPNS, PushGoroutinesAPNS     int
	PushChanSizeMi, PushGoroutinesMi         int
	PushChanSizeHuawei, PushGoroutinesHuawei int
	PushChanSizeOppo, PushGoroutinesOppo     int
	PushChanSizeJpush, PushGoroutinesJpush   int
	PushChanSizeFCM, PushGoroutinesFCM       int
	PassThrough                              int
	RetryTimes                               int
	PushPartInterval                         xtime.Duration
	PushPartChanSize                         int
	PushPartSize                             int
	CallbackSize, CallbackChanLen            int
	UpimgURL                                 string
	UpimgMaxSize                             int64
	UpdateTaskProgressProc                   int
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
