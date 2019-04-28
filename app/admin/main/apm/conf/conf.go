package conf

import (
	"errors"
	"flag"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// config
	confPath string
	client   *conf.Client
	// Conf global config.
	Conf = &Config{}
)

// Config def.
type Config struct {
	// base
	Superman []string
	AppToken string
	//utBaseLine
	UTBaseLine *UTBaseLine
	// Host
	Host *Host
	// Canal
	Canal *Canal
	// log
	Log *log.Config
	// ecode
	Ecode *ecode.Config
	// client
	HTTPClient *bm.ClientConfig
	// db
	ORM *orm.Config
	// db databus
	ORMDatabus *orm.Config
	// db canal
	ORMCanal *orm.Config
	// kafka
	Kafka map[string]*Kafka
	// identify
	Auth *permit.Config
	// discovery
	Discovery *Discovery
	// tree
	Tree *Tree
	// databus kafka create topic config
	DatabusConfig *Databus
	//BM
	BM *bm.ServerConfig
	// ManagerReport
	ManagerReport *databus.Config
	// pprof
	Pprof *Pprof
	//Bfs
	Bfs        *Bfs
	Prometheus *Prometheus
	Apps       *Apps
	Cron       *Cron
	BroadCast  *BroadCast
	Gitlab     *Gitlab
	WeChat     *WeChat
	//Alarm
	Alarm *Alarm
	// Redis
	// Redis *redis.Config
	Redis *Redis
}

// Alarm .
type Alarm struct {
	DatabusURL string
	DatabusKey string
}

// WeChat .
type WeChat struct {
	Users  []string
	ChatID string
}

// BroadCast .
type BroadCast struct {
	TenCent  []string
	KingSoft []string
}

// Cron .
type Cron struct {
	Crontab     string
	CrontabRepo string
}

// Prometheus .
type Prometheus struct {
	URL    string
	Key    string
	Secret string
}

// Apps .
type Apps struct {
	Name []string
	Max  int64
}

//UTBaseLine .
type UTBaseLine struct {
	Coverage int
	Passrate int
}

// Bfs bfs config
type Bfs struct {
	Addr        string
	Bucket      string
	Key         string
	Secret      string
	MaxFileSize int
}

// Pprof dir path
type Pprof struct {
	Dir    string
	GoPath string
}

// Databus config
type Databus struct {
	Partitions int32
	Factor     int16
}

// Tree PlatformID
type Tree struct {
	PlatformID    string
	MsmPlatformID string
}

// Host hosts
type Host struct {
	APICo     string
	SVENCo    string
	MANAGERCo string
	DapperCo  string
}

//Canal canal
type Canal struct {
	CANALSVENCo string
	BUILD       string
	Reviewer    []string
}

// Discovery discovery
type Discovery struct {
	API []string
}

// Kafka kafka config
type Kafka struct {
	Brokers []string
}

// Gitlab gitlab config
type Gitlab struct {
	API   string // gitlab api host
	Token string // saga 账户 access token
}

// Redis .
type Redis struct {
	*redis.Config
	ExpireTime xtime.Duration
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

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init int config
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}
