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
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"go-common/library/database/hbase.v2"

	"github.com/BurntSushi/toml"
)

// Conf info.
var (
	ConfPath string
	Conf     = &Config{}
	client   *conf.Client
)

// Config struct.
type Config struct {
	// bm
	BM *HTTPServers
	// db
	DB *DB
	// base
	// elk
	Xlog *log.Config
	// httpClinet
	HTTPClient *HTTPClient
	// tracer
	Tracer *trace.Config
	// Ecode
	Ecode *ecode.Config
	// host
	Host *Host
	// ArchStatus
	ArchStatus map[string]string
	// Redis
	Redis *Redis
	// Databus
	Env      string
	Consume  bool
	IsTest   bool
	UpSub    *upSub
	ChanSize int64
	Monitor  *Monitor
	// rpc client2
	ArticleRPC *rpc.ClientConfig
	// identify
	App      *blademaster.App
	Memcache *MC
	//API HOST
	API  string
	Live string
	// rpc server
	RPCServer  *rpc.ServerConfig
	GRPCServer *warden.ServerConfig
	// auth
	Auth *permit.Config
	// hbase
	HBase *HBaseConfig
	// Manager
	ManagerReport *databus.Config
	//  GRPCClient
	GRPCClient *GRPC
}

// HBaseConfig combine with hbase.Config add ReadTimeout, WriteTimeout
type HBaseConfig struct {
	hbase.Config
	// extra config
	ReadTimeout   time.Duration
	ReadsTimeout  time.Duration
	WriteTimeout  time.Duration
	WritesTimeout time.Duration
}

//UpSub upsub config
type upSub struct {
	*databus.Config
	UpChanSize        int
	ConsumeLimit      int
	RoutineLimit      int
	SpecialAddDBLimit int
}

// GRPC .
type GRPC struct {
	Archive *warden.ClientConfig
	Account *warden.ClientConfig
}

//HTTPServers for http server.
type HTTPServers struct {
	Inner *blademaster.ServerConfig
}

// DB conf.
type DB struct {
	// Creative db
	Creative  *sql.Config
	Manager   *sql.Config
	UpCRM     *sql.Config
	ArcResult *sql.Config
	Archive   *sql.Config
}

// Redis redis config
type Redis struct {
	Up *struct {
		*redis.Config
		UpExpire time.Duration
	}
}

// HTTPClient conf.
type HTTPClient struct {
	Normal *blademaster.ClientConfig
	Slow   *blademaster.ClientConfig
}

// Host conf.
type Host struct {
	API     string
	Live    string
	Search  string
	Manager string
}

// Monitor conf.
type Monitor struct {
	Host          string
	Moni          string
	UserName      string
	AppSecret     string
	AppToken      string
	IntervalAlarm time.Duration
}

//App for key secret.
type App struct {
	Key    string
	Secret string
}

//MC memcache
type MC struct {
	UpExpire        time.Duration
	UpSpecialExpire time.Duration
	Up              *memcache.Config
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init conf.
func Init() (err error) {
	if ConfPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(ConfPath, &Conf)
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
