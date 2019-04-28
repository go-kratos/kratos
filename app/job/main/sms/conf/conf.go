package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
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
	Log        *log.Config
	Tracer     *trace.Config
	Databus    *databus.Config
	SmsGRPC    *warden.ClientConfig
	HTTPClient *bm.ClientConfig
	HTTPServer *bm.ServerConfig
	UserReport *databus.Config
	Wechat     *wechat
	Provider   *Provider
	Speedup    *speedup
	Sms        *sms
}

type sms struct {
	// PassportMobileURL 从passport获取用户手机号
	PassportMobileURL string
	// CallbackProc 处理回执的并发数
	CallbackProc int
	// SingleSendGoroutines 单发短信的并发数
	SingleSendProc int
	// BatchSendGoroutines 批量发送短信的并发数
	BatchSendProc int
	// MonitorProcDuration 定期监控databus有没有消费
	MonitorProcDuration xtime.Duration
	// Blacklist 黑名单手机号，用于压测
	Blacklist []string
}

type wechat struct {
	Token    string
	Secret   string
	Username string
}

// Provider provider conf
type Provider struct {
	Providers []int32
	// meng wang
	MengWangSmsURL          string
	MengWangSmsUser         string
	MengWangSmsPwd          string
	MengWangActURL          string
	MengWangBatchURL        string
	MengWangActUser         string
	MengWangActPwd          string
	MengWangInternationURL  string
	MengWangInternationUser string
	MengWangInternationPwd  string
	// chaung lan
	ChuangLanSmsURL          string
	ChuangLanSmsUser         string
	ChuangLanSmsPwd          string
	ChuangLanActURL          string
	ChuangLanActUser         string
	ChuangLanActPwd          string
	ChuangLanInternationURL  string
	ChuangLanInternationUser string
	ChuangLanInternationPwd  string
	// chuang lan callback
	ChuangLanSmsCallbackURL           string
	ChuangLanActCallbackURL           string
	ChuangLanInternationalCallbackURL string
	// meng wang callback
	MengWangSmsCallbackURL           string
	MengWangActCallbackURL           string
	MengWangInternationalCallbackURL string
}

// speedup network
type speedup struct {
	Switch bool
	// meng wang
	MengWangSmsURL         string
	MengWangActURL         string
	MengWangBatchURL       string
	MengWangInternationURL string
	// chaung lan
	ChuangLanSmsURL         string
	ChuangLanInternationURL string
	ChuangLanActURL         string
	// meng wang callback
	MengWangSmsCallbackURL           string
	MengWangActCallbackURL           string
	MengWangInternationalCallbackURL string
	// chaung lan callback
	ChuangLanSmsCallbackURL           string
	ChuangLanActCallbackURL           string
	ChuangLanInternationalCallbackURL string
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
	err = load()
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
