package conf

import (
	"errors"
	"flag"
	"path"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log      *log.Config
	BM       *bm.ServerConfig
	Tracer   *trace.Config
	Memcache *memcache.Config
	MySQL    *sql.Config
	Ecode    *ecode.Config
	// Property
	Property *Property
	// mail
	MailConf         *Mail
	MailTemplateConf *MailTemplateConfig
	// rpc client
	GRPCClient *RPC
}

// RPC rpc client config.
type RPC struct {
	Account *warden.ClientConfig
}

// Property config for biz logic.
type Property struct {
	UpMcnSignStateCron    string
	UpMcnUpStateCron      string
	UpExpirePayCron       string
	UpMcnDataSummaryCron  string
	McnRecommendCron      string
	DealFailRecommendCron string
	CheckMcnSignUpDueCron string
}

// Mail 邮件配置
type Mail struct {
	Host               string
	Port               int
	Username, Password string
	DueMailReceivers   []string //  []adminname, send to adminname@bilibili.com
	DueAuthorityGroups []string
}

//MailTemplateConfig mail template conf
type MailTemplateConfig struct {
	SignTmplTitle   string
	SignTmplContent string
	PayTmplTitle    string
	PayTmplContent  string
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
	if err != nil {
		return
	}
	var templateConfPath = path.Join(path.Dir(confPath), "mail-template.toml")
	_, err = toml.DecodeFile(templateConfPath, &Conf)
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
