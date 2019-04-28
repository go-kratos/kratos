package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Config struct of conf.
type Config struct {
	// base
	// log
	Log *log.Config
	// mysql
	Mysql *sql.Config
	// http
	BM *bm.ServerConfig
	// hbase
	Hbase *HBaseConfig
	// redis
	Redis *Redis
	// extra property
	Property *Property
}

// HBaseConfig is.
type HBaseConfig struct {
	*hbase.Config
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

// Redis redis.
type Redis struct {
	*redis.Config
	Expire xtime.Duration
}

// Property figure conf
type Property struct {
	PendingMidStart int64
	PendingMidShard int64
	PendingMidRetry int64
	CalcWeekOffset  int64
	CycleCron       string
	ConcurrencySize int64
	FixRecord       bool
	CycleAll        bool
	CycleAllCron    string
	Calc            *Calc
}

// Calc figure calc config
type Calc struct {
	InitLawfulScore     int64
	InitWideScore       int64
	InitFriendlyScore   int64
	InitBountyScore     int64
	InitCreativityScore int64
	K1                  float64
	K2                  float64
	K3                  float64
	K4                  float64
	K5                  float64
	LawfulNegMax        float64
	LawfulPosMax        float64
	LawfulPosK          float64
	LawfulNegK1         float64
	LawfulNegK2         float64
	LawfulPosL          float64
	LawfulNegL          float64
	LawfulPosC1         float64
	LawfulPosC2         float64
	LawfulPosC3         float64
	LawfulNegC1         float64
	LawfulNegC2         float64
	LawfulNegC3         float64
	LawfulPosQ1         float64
	LawfulPosQ2         float64
	LawfulPosQ3         float64
	LawfulNegQ1         float64
	LawfulNegQ2         float64
	LawfulNegQ3         float64
	WidePosMax          float64
	WidePosK            float64
	WideC1              float64
	WideQ1              float64
	WideC2              float64
	WideQ2              float64
	FriendlyPosMax      float64
	FriendlyNegMax      float64
	FriendlyPosK        float64
	FriendlyNegK        float64
	FriendlyPosL        float64
	FriendlyNegL        float64
	FriendlyPosQ1       float64
	FriendlyPosC1       float64
	FriendlyPosQ2       float64
	FriendlyPosC2       float64
	FriendlyPosQ3       float64
	FriendlyPosC3       float64
	FriendlyNegQ1       float64
	FriendlyNegC1       float64
	FriendlyNegQ2       float64
	FriendlyNegC2       float64
	FriendlyNegQ3       float64
	FriendlyNegC3       float64
	FriendlyNegQ4       float64
	FriendlyNegC4       float64
	BountyMax           float64
	BountyPosL          float64
	BountyK             float64
	BountyQ1            float64
	BountyC1            float64
	BountyQ2            float64
	BountyC2            float64
	BountyQ3            float64
	BountyC3            float64
	CreativityPosMax    float64
	CreativityPosK      float64
	CreativityPosL1     float64
	CreativityQ1        float64
	CreativityC1        float64
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
