package conf

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf info.
var (
	ConfPath     string
	Conf         = &Config{}
	client       *conf.Client
	CreditConfig = &CreditConf{}
	IsMaster     = true
)

const (
	//ServiceName service name
	ServiceName = "upcredit-service"
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
	// report log
	LogCli *log.AgentConfig
	// httpClinet
	HTTPClient *HTTPClient
	// tracer
	Tracer *trace.Config

	// Redis
	Redis *Redis

	// rpc server
	RPCServer *rpc.ServerConfig
	// auth
	Auth   *permit.Config
	IsTest bool

	CreditLogSub      *databus.Config
	BusinessBinLogSub *databus.Config
	RunStatJobConf    *RunStatJob

	MiscConf          *MiscConfig
	ElectionZooKeeper *Zookeeper
}

//UpSub upsub config
//type upSub struct {
//	*databus.Config
//	UpChanSize   int
//	ConsumeLimit int
//	RoutineLimit int
//}

//MiscConfig other config set
type MiscConfig struct {
	CreditLogWriteRoutineNum int
	BusinessBinLogLimitRate  float64 // 每秒多少个，business bin log 消费速度
}

//HTTPServers for http server.
type HTTPServers struct {
	Inner *blademaster.ServerConfig
}

// DB conf.
type DB struct {
	Upcrm       *orm.Config
	UpcrmReader *orm.Config
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
	Normal *bm.ClientConfig
	Slow   *bm.ClientConfig
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

//RunStatJob 定时任务时间
type RunStatJob struct {
	// 启动时间，比如 12:00:00，每天定时运行
	StartTime string
	// 起的计算线程数
	WorkerNumber int
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

//CreditLog 需要记录日志的那些稿件状态，在配置文件中配置。只有这些状态，才会记录信用日志
type CreditLog struct {
	NeedLogState map[int]CreditLogStateInfo
}

//CreditLogStateInfo nothing
type CreditLogStateInfo struct {
}

// Zookeeper Server&Client settings.
type Zookeeper struct {
	Root    string
	Addrs   []string
	Timeout time.Duration
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
	if err != nil {
		return
	}
	ConfPath = strings.Replace(ConfPath, string(os.PathSeparator), "/", -1)
	var dir = path.Dir(ConfPath)
	var articleConfPath = path.Join(dir, "credit_score_conf.toml")
	_, err = toml.DecodeFile(articleConfPath, &CreditConfig)
	CreditConfig.AfterLoad()
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

	if tomlStr, ok = client.Value("upcredit-service.toml"); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(tomlStr, &tmpConf); err != nil {
		return errors.New("could not decode toml config")
	}
	*Conf = *tmpConf

	fmt.Printf("loading credit_score_conf.toml from remoate...")
	if tomlStr, ok = client.Value("credit_score_conf.toml"); !ok {
		return errors.New("load config center error for credit_score_conf.toml")
	}

	var tmpConf2 *CreditConf
	if _, err = toml.Decode(tomlStr, &tmpConf2); err != nil {
		return errors.New("could not decode toml config for credit_score_conf.toml")
	}
	*CreditConfig = *tmpConf2
	CreditConfig.AfterLoad()
	return
}
