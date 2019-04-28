package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/orm"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf Config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log *log.Config

	Bugly *BuglyConf

	BM *bm.ServerConfig

	Ecode *ecode.Config

	ORM *orm.Config

	HTTPClient *bm.ClientConfig

	Mail *Mail

	Auth *permit.Config

	Memcache *Memcache

	Scheduler *Scheduler

	Tapd *TapdConf
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Scheduler scheduler
type Scheduler struct {
	BatchRunEnableVersion string

	BatchRunUpdateTapdBug string

	DisableBatchRunOverTime string
	BatchRunOverHourTime    int

	SyncWechatContact string

	Active bool
}

// Memcache memcache
type Memcache struct {
	*memcache.Config
	Expire xtime.Duration
}

// Mail mail
type Mail struct {
	Host        string
	Port        int
	Username    string
	Password    string
	NoticeOwner []string
}

// BuglyConf Bugly Conf.
type BuglyConf struct {
	Host             string
	UrlRetryCount    int
	CookieUsageUpper int
	IssuePageSize    int
	IssueCountUpper  int

	SuperOwner []string
}

// TapdConf Tapd Conf.
type TapdConf struct {
	BugOperateAuth bool
}

// Tapd Tapd info
type Tapd struct {
	IterationWorkspaceIDs []string
	StoryWorkspaceIDs     []string
	BugWorkspaceIDs       []string
	IPS                   int
	SPS                   int
	SCPS                  int
	CPS                   int
	StoryFilePath         string
	ChangeFilePath        string
	IterationFilePath     string
	BugFilePath           string
	RetryTime             int
	WaitTime              xtime.Duration
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
