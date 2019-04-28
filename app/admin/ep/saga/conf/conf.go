package conf

import (
	"flag"
	"strings"

	"go-common/app/admin/ep/saga/model"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

const (
	_configKey = "saga-admin.toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf store the global config
	Conf   = &Config{}
	reload chan bool
)

// Config def.
type Config struct {
	Tracer     *trace.Config
	BM         *bm.ServerConfig
	HTTPClient *bm.ClientConfig
	Memcache   *Memcache
	Redis      *redis.Config
	Log        *log.Config
	ORM        *orm.Config
	Permit     *permit.Config2
	Property   *Property
}

// Memcache config.
type Memcache struct {
	MC             *memcache.Config
	MCRecordExpire xtime.Duration
}

// Property config for biz logic.
type Property struct {
	Gitlab *struct {
		API   string // gitlab api host
		Token string // saga 账户 access token
	}
	Git *struct {
		API           string // gitlab api host
		Token         string // saga 账户 access token
		CheckCron     string
		UserList      []string
		AlertPipeline []*model.AlertPipeline
	}
	SyncProject *struct {
		CheckCron string
	}
	SyncData *struct {
		SyncAllTime     bool
		DefaultSyncDays int
		CheckCron       string
		CheckCronAll    string
		CheckCronWeek   string
		WechatUser      []string
	}

	Department *model.PairKey
	Business   *model.PairKey
	DeInfo     []*model.PairKey
	BuInfo     []*model.PairKey
	Mail       *struct {
		Host    string
		Port    int
		Address string
		Pwd     string
		Name    string
	}
	Wechat  *model.AppConfig
	Contact *model.AppConfig
	Group   *struct {
		Name       string
		Department string
		Business   string
	}
	DefaultProject *struct {
		ProjectIDs []int
		Status     []string
		Types      []int
	}
	Sven *struct {
		ConfigValue      string
		Configs          string
		ConfigUpdate     string
		TagUpdate        string
		ConfigsParam     *model.ConfigsParam
		SagaConfigsParam *model.SagaConfigsParam
	}
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
	reload = make(chan bool, 10)
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
	Conf = parseTeamInfo(Conf)
	return
}

func configCenter() (err error) {
	if client, err = conf.New(); err != nil {
		panic(err)
	}
	if err = load(); err != nil {
		return
	}
	client.WatchAll()
	go func() {
		for range client.Event() {
			log.Info("config reload")
			if load() != nil {
				log.Error("config reload error (%v)", err)
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
	if s, ok = client.Value(_configKey); !ok {
		err = errors.Errorf("load config center error [%s]", _configKey)
		return
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		err = errors.Wrapf(err, "could not decode config err(%+v)", err)
		return
	}
	Conf = parseTeamInfo(tmpConf)
	return
}

func parseTeamInfo(c *Config) *Config {

	DeLabel := strings.Fields(c.Property.Department.Label)
	DeValue := strings.Fields(c.Property.Department.Value)
	for i := 0; i < len(DeLabel); i++ {

		info := &model.PairKey{
			Label: DeLabel[i],
			Value: DeValue[i],
		}
		c.Property.DeInfo = append(c.Property.DeInfo, info)
	}

	buLabel := strings.Fields(c.Property.Business.Label)
	buValue := strings.Fields(c.Property.Business.Value)
	for i := 0; i < len(buLabel); i++ {

		info := &model.PairKey{
			Label: buLabel[i],
			Value: buValue[i],
		}
		c.Property.BuInfo = append(c.Property.BuInfo, info)
	}

	/*for _, r := range c.Property.Developer {
		r.Total = r.Android + r.Ios + r.Service + r.Web
	}*/

	return c
}
