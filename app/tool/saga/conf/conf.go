package conf

import (
	"flag"
	"regexp"
	"sync"
	"time"

	"go-common/app/tool/saga/model"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/orm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

const (
	configKey = "saga.toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf store the global config
	Conf   = &Config{}
	reload chan bool
)

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
	reload = make(chan bool, 10)
}

// Config def.
type Config struct {
	Tracer     *trace.Config
	BM         *bm.ServerConfig
	HTTPClient *bm.ClientConfig
	Memcache   *Memcache
	Redis      *redis.Config
	HBase      *HBaseConfig
	Log        *log.Config
	Property   *Property
	sync.RWMutex
	// orm
	ORM *orm.Config
}

// Memcache config.
type Memcache struct {
	MR             *memcache.Config
	MRRecordExpire xtime.Duration
}

// HBaseConfig for new hbase client.
type HBaseConfig struct {
	*hbase.Config
	WriteTimeout xtime.Duration
	ReadTimeout  xtime.Duration
}

// Property config for biz logic.
type Property struct {
	TaskInterval xtime.Duration // 任务轮询时间
	TaskTimeout  xtime.Duration // 任务超时时间
	PollPipeline xtime.Duration //pipeline轮询间隔时间
	Gitlab       *struct {
		API   string // gitlab api host
		Token string // saga 账户 access token
	}
	WebHooks []*model.WebHook
	Mail     *struct {
		Host    string
		Port    int
		Address string
		Pwd     string
		Name    string
	}
	HealthCheck *struct {
		CheckCron  string
		AlertAddrs []*model.MailAddress
	}
	ReportRequiredVisible *struct {
		CheckCron  string
		AlertAddrs []*model.MailAddress
	}
	SyncContact *struct {
		CheckCron string
	}
	UT *struct {
		Rate float64
	}
	Wechat  *model.AppConfig
	Contact *model.AppConfig
	Repos   []*model.RepoConfig
}

// Init init conf.
func Init() (err error) {
	if confPath == "" {
		return configCenter()
	}
	if _, err = toml.DecodeFile(confPath, &Conf); err != nil {
		log.Error("toml.DecodeFile(%s) err(%+v)", confPath, err)
		return
	}
	Conf = doDefault(Conf)
	return
}

func configCenter() (err error) {
	if client, err = conf.New(); err != nil {
		panic(err)
	}
	if err = load(); err != nil {
		return
	}
	//client.WatchAll()
	client.Watch(configKey)
	go func() {
		for range client.Event() {
			log.Info("config reload")
			if err = load(); err != nil {
				log.Error("config reload error (+%v)", err)
			} else {
				reload <- true
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
	if s, ok = client.Value(configKey); !ok {
		err = errors.Errorf("load config center error [%s]", configKey)
		return
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		err = errors.Wrapf(err, "could not decode config err(%+v)", err)
		return
	}
	Conf = doDefault(tmpConf)
	return
}

func doDefault(c *Config) *Config {
	if int64(c.Property.TaskInterval) == 0 {
		c.Property.TaskInterval = xtime.Duration(3 * time.Second)
	}
	if int64(c.Property.TaskTimeout) == 0 {
		c.Property.TaskTimeout = xtime.Duration(time.Minute)
	}
	if int64(c.Property.PollPipeline) == 0 {
		c.Property.PollPipeline = xtime.Duration(10 * time.Second)
	}
	for _, r := range c.Property.Repos {
		if r.GName == "" {
			r.GName = r.Name
		}
		if r.Language == "" {
			r.Language = "any"
		}
		if r.MinReviewer < 0 {
			r.MinReviewer = 0
		}
		if r.LockTimeout == 0 {
			r.LockTimeout = 600
		}
		if len(r.AuthBranches) == 0 {
			r.AuthBranches = []string{"master"}
		}
		if len(r.TargetBranches) == 0 {
			r.TargetBranches = r.AuthBranches
		}
		for _, b := range r.TargetBranches {
			r.TargetBranchRegexes = append(r.TargetBranchRegexes, regexp.MustCompile(b))
		}
	}
	return c
}

// ReloadEvents return the reload chan
func ReloadEvents() <-chan bool {
	return reload
}
