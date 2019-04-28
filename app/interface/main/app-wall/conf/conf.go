package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

type Config struct {
	// Env
	Env string
	// db
	MySQL *MySQL
	// show  XLog
	Log *log.Config
	// tracer
	Tracer *trace.Config
	// tick time
	Tick xtime.Duration
	// httpClinet
	HTTPClient *bm.ClientConfig
	// HTTPTelecom
	HTTPTelecom *bm.ClientConfig
	// HTTPBroadband
	HTTPBroadband *bm.ClientConfig
	// HTTPUnicom
	HTTPUnicom *bm.ClientConfig
	// HTTPUnicom
	HTTPActive *bm.ClientConfig
	// bm http
	BM *HTTPServers
	// rpc account
	AccountRPC *rpc.ClientConfig
	// seq
	SeqRPC *rpc.ClientConfig
	// host
	Host *Host
	// ecode
	Ecode *ecode.Config
	// Report
	Report *databus.Config
	// iplimit
	IPLimit *IPLimit
	// infoc2
	UnicomUserInfoc2 *infoc.Config
	UnicomIpInfoc2   *infoc.Config
	UnicomPackInfoc  *infoc.Config
	// Seq
	Seq *Seq
	// Telecom
	Telecom *Telecom
	// Redis
	Redis *Redis
	// mc
	Memcache *Memcache
	// reddot
	Reddot *Reddot
	// unicom
	Unicom *Unicom
	ES     *elastic.Config
	// databus
	UnicomDatabus *databus.Config
}

type Host struct {
	APICo               string
	Dotin               string
	Live                string
	APILive             string
	Telecom             string
	Unicom              string
	UnicomFlow          string
	Broadband           string
	Sms                 string
	Mall                string
	TelecomReturnURL    string
	TelecomCancelPayURL string
}

type HTTPServers struct {
	Outer *bm.ServerConfig
}

type Seq struct {
	BusinessID int64
	Token      string
}

// App bilibili intranet authorization.
type App struct {
	Key    string
	Secret string
}

type MySQL struct {
	Show *sql.Config
}

type IPLimit struct {
	MobileIPFile string
	Addrs        map[string][]string
}

type Reddot struct {
	StartTime string
	EndTime   string
}

type Unicom struct {
	KeyExpired xtime.Duration
	FlowWait   xtime.Duration
}

type Telecom struct {
	KeyExpired         xtime.Duration
	PayKeyExpired      xtime.Duration
	SMSTemplate        string
	SMSMsgTemplate     string
	SMSFlowTemplate    string
	SMSOrderTemplateOK string
	FlowPercentage     int
	Area               map[string][]string
}

type Redis struct {
	Recommend *struct {
		*redis.Config
		Expire xtime.Duration
	}
}

type Memcache struct {
	Operator *struct {
		*memcache.Config
		Expire xtime.Duration
	}
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config.
func Init() (err error) {
	if confPath != "" {
		_, err = toml.DecodeFile(confPath, &Conf)
		return
	}
	err = remote()
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
