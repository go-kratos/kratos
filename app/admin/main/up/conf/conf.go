package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
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
	HTTPServer *blademaster.ServerConfig
	// db
	DB *DB
	// base
	// elk
	XLog *log.Config
	// report log
	LogCli *log.AgentConfig
	// httpClinet
	HTTPClient *HTTPClient
	// tracer
	Tracer *trace.Config
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
	// auth
	Auth *permit.Config
	// hbase
	HBase   *HBaseConfig
	HBase2  *HBaseConfig
	BfsConf *Bfs
	Debug   bool

	TimeConf *TimeConfig
	MailConf *Mail
	// manager log config
	ManagerLog *databus.Config
	//  GRPCClient
	GRPCClient *GRPC
}

// GRPC .
type GRPC struct {
	Archive *warden.ClientConfig
	Account *warden.ClientConfig
}

// HBaseConfig for new hbase client.
type HBaseConfig struct {
	*hbase.Config
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

//UpSub upsub config
type upSub struct {
	*databus.Config
	UpChanSize   int
	ConsumeLimit int
	RoutineLimit int
}

// DB conf.
type DB struct {
	// Creative db
	Creative *sql.Config
	Manager  *sql.Config
	Upcrm    *orm.Config
}

// Redis conf.
type Redis struct {
	Databus *struct {
		*redis.Config
		Expire time.Duration
	}
}

// HTTPClient conf.
type HTTPClient struct {
	Normal *blademaster.ClientConfig
	Slow   *blademaster.ClientConfig
}

// Host conf.
type Host struct {
	API      string
	Live     string
	Search   string
	Manager  string
	Data     string
	Tag      string
	Coverrec string
	Videoup  string
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
	UpExpire time.Duration
	Up       *memcache.Config
}

// Bfs struct.
type Bfs struct {
	Addr        string
	Bucket      string
	Key         string
	Secret      string
	MaxFileSize int
}

//TimeConfig 定期任务时间
type TimeConfig struct {
	TaskScheduleTime     string // 每天定时检查完成的task情况，format "10:59:59"
	CheckDueScheduleTime string // 每天定时检查快要过期的任务，format "10:59:59"
	RefreshUpRankTime    string // 每天定时检查upRank表的任务，format "10:59:59"
}

//Mail 邮件配置
type Mail struct {
	Host               string
	Port               int
	Username, Password string
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
	if tomlStr, ok = client.Value("up-admin.toml"); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(tomlStr, &tmpConf); err != nil {
		return errors.New("could not decode toml config")
	}
	*Conf = *tmpConf
	return
}
